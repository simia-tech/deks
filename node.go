package edkvs

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/simia-tech/conflux/recon"
	"github.com/simia-tech/errx"
	redisserver "github.com/tidwall/redcon"
)

const (
	cmdHelp         = "help"
	cmdQuit         = "quit"
	cmdPing         = "ping"
	cmdSet          = "set"
	cmdGet          = "get"
	cmdDelete       = "del"
	cmdKeys         = "keys"
	cmdTidy         = "tidy"
	cmdSetContainer = "cset"        // hidden
	cmdGetContainer = "cget"        // hidden
	cmdReconcilate  = "reconcilate" // hidden

	help = `Supported commands:
help              - prints this help message
set <key> <value> - sets <value> at <key>
get <key>         - returns value at <key>
del <key>         - removes value at <key>
keys              - returns all keys
tidy              - cleans up the store
quit              - closes the connection
`
)

// Node defines a edkvs node.
type Node struct {
	store    *Store
	listener net.Listener
	metric   Metric
	peer     *recon.Peer
	streams  []*stream
}

// NewNode returns a new node.
func NewNode(store *Store, listenURL string, m Metric) (*Node, error) {
	network, address, err := parseURL(listenURL)
	if err != nil {
		return nil, errx.Annotatef(err, "parse listen url [%s]", listenURL)
	}

	l, err := net.Listen(network, address)
	if err != nil {
		return nil, errx.Annotatef(err, "listen [%s %s]", network, address)
	}
	// log.Printf("node is listening at [%s %s]", l.Addr().Network(), l.Addr().String())

	settings := recon.DefaultSettings()
	peer := recon.NewPeer(settings, store.State().prefixTree())

	n := &Node{
		store:    store,
		listener: l,
		metric:   m,
		peer:     peer,
		streams:  make([]*stream, 0),
	}
	store.updateFn = n.update
	go n.acceptLoop()
	return n, nil
}

// ListenURL returns the url of the listener.
func (n *Node) ListenURL() string {
	addr := n.listener.Addr()
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}

// Close tears down the node.
func (n *Node) Close() error {
	for _, stream := range n.streams {
		stream.close()
	}
	if err := n.listener.Close(); err != nil {
		if isClosedNetworkError(err) {
			return nil
		}
		return errx.Annotatef(err, "close listener")
	}
	return nil
}

// AddPeer adds another node as a target for updates.
func (n *Node) AddPeer(
	peerURL string,
	peerPingInterval time.Duration,
	peerReconnectInterval time.Duration,
) {
	n.streams = append(n.streams, newStream(peerURL, peerPingInterval, peerReconnectInterval, n.metric))
}

// Reconcilate performs a reconsiliation with the node at the provided address.
func (n *Node) Reconcilate(url string) (int, error) {
	conn, err := Dial(url)
	if err != nil {
		return 0, errx.Annotatef(err, "dial [%s]", url)
	}
	defer conn.Close()

	netConn, err := conn.Reconsilate()
	if err != nil {
		return 0, errx.Annotatef(err, "reconcilate")
	}

	keyHashes, _, err := n.peer.Reconcilate(netConn, 100)
	if err != nil {
		return 0, errx.Annotatef(err, "reconcilate")
	}

	payloadConn, err := Dial(url)
	if err != nil {
		return 0, errx.Annotatef(err, "dial [%s]", url)
	}
	defer payloadConn.Close()

	for _, keyHash := range keyHashes {
		kh := newKeyHash(keyHash)
		c, err := payloadConn.getContainer(kh)
		if err != nil {
			return 0, errx.Annotatef(err, "get container")
		}
		if err := n.store.setContainer(kh, c); err != nil {
			return 0, errx.Annotatef(err, "set container")
		}
	}

	return len(keyHashes), nil
}

func (n *Node) acceptLoop() {
	done := false
	var err error
	for !done {
		done, err = n.accept()
		if err != nil {
			log.Printf("accept loop: %v", err)
			done = true
		}
	}
}

func (n *Node) accept() (bool, error) {
	conn, err := n.listener.Accept()
	if err != nil {
		if isClosedNetworkError(err) {
			return true, nil
		}
		return true, errx.Annotatef(err, "accept")
	}

	clientURL := urlFor(conn.RemoteAddr())

	go func() {
		if err := n.handleConn(conn); err != nil {
			log.Printf("conn %s: %v", conn.RemoteAddr(), err)
		}
		if err := conn.Close(); err != nil {
			if !isClosedNetworkError(err) {
				log.Printf("close conn %s: %v", conn.RemoteAddr(), err)
			}
		}
		n.metric.ClientDisconnected(clientURL)
	}()

	n.metric.ClientConnected(clientURL)

	return false, nil
}

func (n *Node) handleConn(conn net.Conn) error {
	r := redisserver.NewReader(conn)
	w := redisserver.NewWriter(conn)

	done := false
	for !done {
		cmd, err := r.ReadCommand()
		if err == io.EOF {
			done = true
			continue
		}
		if err != nil {
			return errx.Annotatef(err, "read command")
		}

		command := strings.ToLower(string(cmd.Args[0]))
		arguments := cmd.Args[1:]

		switch command {
		case cmdHelp:
			w.WriteBulkString(help)
		case cmdQuit:
			done = true
			w.WriteString("OK")
		case cmdPing:
			w.WriteString("OK")
		case cmdSet:
			if err := n.store.Set(arguments[0], arguments[1]); err != nil {
				return errx.Annotatef(err, "set")
			}
			w.WriteString("OK")
		case cmdGet:
			value, err := n.store.Get(arguments[0])
			if err != nil {
				return errx.Annotatef(err, "get [%s]", arguments[0])
			}
			w.WriteBulk(value)
		case cmdDelete:
			if err := n.store.Delete(arguments[0]); err != nil {
				return errx.Annotatef(err, "delete [%s]", arguments[0])
			}
			w.WriteString("OK")
		case cmdKeys:
			w.WriteArray(n.store.Len())
			n.store.Each(func(key, _ []byte) error {
				w.WriteBulk(key)
				return nil
			})
		case cmdTidy:
			if err := n.store.Tidy(); err != nil {
				return errx.Annotatef(err, "tidy")
			}
			w.WriteString("OK")
		case cmdSetContainer:
			kh := keyHash{}
			copy(kh[:], arguments[0][:keyHashSize])
			if err := n.store.setContainer(kh, arguments[1]); err != nil {
				return errx.Annotatef(err, "set container [%s]", kh)
			}
			w.WriteString("OK")
		case cmdGetContainer:
			kh := keyHash{}
			copy(kh[:], arguments[0][:keyHashSize])
			c, err := n.store.getContainer(kh)
			if err != nil {
				return errx.Annotatef(err, "get container [%s]", kh)
			}
			w.WriteBulk(c)
		case cmdReconcilate:
			w.WriteString("OK")
			if err := w.Flush(); err != nil {
				return errx.Annotatef(err, "flush")
			}
			if err := n.peer.Accept(conn); err != nil {
				return errx.Annotatef(err, "recon accept")
			}
			return nil // exit command loop
		default:
			w.WriteError(fmt.Sprintf("unknown command [%s]", command))
		}

		if err := w.Flush(); err != nil {
			return errx.Annotatef(err, "flush")
		}
	}

	return nil
}

func (n *Node) update(kh keyHash, container *container) {
	for _, stream := range n.streams {
		stream.update(kh, container)
	}
}

func parseURL(u string) (string, string, error) {
	url, err := url.Parse(u)
	if err != nil {
		return "", "", errx.Annotatef(err, "parse url [%s]", u)
	}
	return url.Scheme, url.Host, nil
}

func urlFor(addr net.Addr) string {
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}

func isClosedNetworkError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
		return true
	}
	return false
}
