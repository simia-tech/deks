package edkvs

// Metric defines the EDKVS metric interface.
type Metric interface {
	CountChanged(int, int)
}
