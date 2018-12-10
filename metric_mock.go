package kea

// MetricMock defines a metric mock.
type MetricMock struct{}

// NewMetricMock returns a new metric mock.
func NewMetricMock() *MetricMock {
	return &MetricMock{}
}

// CountChanged is called if the number of value or deleted values has changed.
func (mm *MetricMock) CountChanged(valueCount, deletedCount int) {}

// ClientConnected is called if a new client connects.
func (mm *MetricMock) ClientConnected(_ string) {}

// ClientDisconnected is called if a client disconnects.
func (mm *MetricMock) ClientDisconnected(_ string) {}

// PeerConnected is called if a new peer connects.
func (mm *MetricMock) PeerConnected(_ string) {}

// PeerDisconnected is called if a peer disconnects.
func (mm *MetricMock) PeerDisconnected(_ string) {}
