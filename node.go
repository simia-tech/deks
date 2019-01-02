package deks

import (
	"context"
	"log"
	"time"

	"github.com/simia-tech/errx"
)

// Node defines the node.
type Node struct {
	Store  *Store
	server *Server
	cancel context.CancelFunc
}

// NewNode returns a new node.
func NewNode(o Options, m Metric) (*Node, error) {
	store := NewStore(m)
	server, err := NewServer(store, o.ListenURL, m)
	if err != nil {
		return nil, errx.Annotatef(err, "new server")
	}
	for _, peerURL := range o.PeerURLs {
		_, err := server.Reconcilate(peerURL)
		if err != nil {
			log.Printf("reconsilate: %v", err)
		}
		if err := server.AddPeer(peerURL, o.PeerPingInterval, o.PeerReconnectInterval); err != nil {
			return nil, errx.Annotatef(err, "peer add")
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(o.TidyInterval)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := store.Tidy(); err != nil {
					log.Printf("tidy: %v", err)
				}
			}
		}
	}()

	return &Node{
		Store:  store,
		server: server,
		cancel: cancel,
	}, nil
}

// ListenURL returns the listen url.
func (n *Node) ListenURL() string {
	return n.server.ListenURL()
}

// Close tears down the node.
func (n *Node) Close() error {
	n.cancel()
	return n.server.Close()
}
