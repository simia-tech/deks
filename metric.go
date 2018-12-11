package deks

// Metric defines the metric interface.
type Metric interface {
	CountChanged(int, int)
	ClientConnected(string)
	ClientDisconnected(string)
	PeerConnected(string)
	PeerDisconnected(string)
}
