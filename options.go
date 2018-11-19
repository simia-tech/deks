package edkvs

import "time"

// Options defines all edkvs options.
type Options struct {
	// Listener address in format 'tcp://localhost:5000'.
	ListenURL string

	// Peer addreses in the format `tcp://localhost:5000`.
	PeerURLs []string

	// PeerPingInterval defines the interval in which a peer is pinged in order to test it's availbility.
	PeerPingInterval time.Duration

	// PeerReconnectInterval defines a duration after which a failing peer is reconnected.
	PeerReconnectInterval time.Duration

	// TidyInterval defines the interval in which the store is cleaned up.
	TidyInterval time.Duration
}
