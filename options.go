package edkvs

import "time"

// Options defines all edkvs options.
type Options struct {
	// Listener address in format 'tcp://localhost:5000'.
	ListenURL string

	// Peer addreses in the format `tcp://localhost:5000`.
	PeerURLs []string

	// PeerReconnectInterval defines a duration after which a failing peer is reconnected.
	PeerReconnectInterval time.Duration
}
