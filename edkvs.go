package edkvs

import (
	"log"

	"github.com/simia-tech/errx"
)

// EDKVS defines the Embedded Distributed Key-Value Store.
type EDKVS struct {
	store *Store
	node  *Node
}

// NewEDKVS returns a new EDKVS.
func NewEDKVS(o Options) (*EDKVS, error) {
	store := NewStore()
	node, err := NewNode(store, o.ListenURL)
	if err != nil {
		return nil, errx.Annotatef(err, "new node")
	}
	for _, peerURL := range o.PeerURLs {
		count, err := node.Reconcilate(peerURL)
		if err != nil {
			log.Printf("reconsilate: %v", err)
		}
		log.Printf("reconsilated %d values from %s", count, peerURL)
		node.AddPeer(peerURL, o.PeerReconnectInterval)
	}
	return &EDKVS{
		store: store,
		node:  node,
	}, nil
}

// ListenURL returns the listen url.
func (e *EDKVS) ListenURL() string {
	return e.node.ListenURL()
}

// Close tears down the EDKVS.
func (e *EDKVS) Close() error {
	return e.node.Close()
}
