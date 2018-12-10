package kea

import (
	"context"
	"log"
	"time"

	"github.com/simia-tech/errx"
)

// EDKVS defines the Embedded Distributed Key-Value Store.
type EDKVS struct {
	Store  *Store
	server *Server
	cancel context.CancelFunc
}

// NewEDKVS returns a new EDKVS.
func NewEDKVS(o Options, m Metric) (*EDKVS, error) {
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
		server.AddPeer(peerURL, o.PeerPingInterval, o.PeerReconnectInterval)
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

	return &EDKVS{
		Store:  store,
		server: server,
		cancel: cancel,
	}, nil
}

// ListenURL returns the listen url.
func (e *EDKVS) ListenURL() string {
	return e.server.ListenURL()
}

// Close tears down the EDKVS.
func (e *EDKVS) Close() error {
	e.cancel()
	return e.server.Close()
}
