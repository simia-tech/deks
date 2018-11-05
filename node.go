package edkvs

import (
	"log"
	"net"

	"github.com/simia-tech/conflux"
	"github.com/simia-tech/conflux/recon"
	"github.com/simia-tech/errx"
)

type Node struct {
	store      *Store
	listener   net.Listener
	prefixTree *recon.MemPrefixTree
	peer       *recon.Peer
}

func NewNode(store *Store, network, address string) (*Node, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, errx.Annotatef(err, "listen [%s %s]", network, address)
	}
	// log.Printf("node is listening at [%s %s]", l.Addr().Network(), l.Addr().String())

	prefixTree := &recon.MemPrefixTree{}
	prefixTree.Init()

	settings := recon.DefaultSettings()
	peer := recon.NewPeer(settings, prefixTree)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); !ok || opErr.Err.Error() != "use of closed network connection" {
					log.Printf("accept: %v", err)
				}
				return
			}

			go func() {
				if err := peer.Accept(conn); err != nil {
					log.Printf("recon accept: %v", err)
					return
				}
			}()
		}
	}()

	return &Node{
		store:      store,
		listener:   l,
		prefixTree: prefixTree,
		peer:       peer,
	}, nil
}

func (n *Node) Addr() net.Addr {
	return n.listener.Addr()
}

func (n *Node) Close() error {
	return n.listener.Close()
}

func (n *Node) Insert(value int) error {
	if err := n.prefixTree.Insert(conflux.Zi(p, value)); err != nil {
		return errx.Annotatef(err, "insert")
	}
	return nil
}

func (n *Node) Reconcilate(addr net.Addr) (int, error) {
	changes, done, err := n.peer.Reconcilate(addr.Network(), addr.String(), 100)
	if err != nil {
		return 0, errx.Annotatef(err, "reconcilate")
	}
	log.Printf("reconcilate / done %v / changes %v", done, changes)

	return len(changes), nil
}

func (n *Node) Elements() []int64 {
	node, err := n.prefixTree.Root()
	if err != nil {
		panic(err)
	}
	elements, err := node.Elements()
	if err != nil {
		panic(err)
	}
	values := []int64{}
	for _, element := range elements {
		values = append(values, element.Int64())
	}
	return values
}
