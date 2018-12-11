package deks

import "log"

// MetricLog defines a metric log.
type MetricLog struct{}

// NewMetricLog returns a new metric log.
func NewMetricLog() *MetricLog {
	return &MetricLog{}
}

// CountChanged is called if the number of value or deleted values has changed.
func (ml *MetricLog) CountChanged(valueCount, deletedCount int) {
	log.Printf("count changed: values = %d / deleted = %d", valueCount, deletedCount)
}

// ClientConnected is called if a new client connects.
func (ml *MetricLog) ClientConnected(clientURL string) {
	log.Printf("client [%s] connected", clientURL)
}

// ClientDisconnected is called if a client disconnects.
func (ml *MetricLog) ClientDisconnected(clientURL string) {
	log.Printf("client [%s] disconnected", clientURL)
}

// PeerConnected is called if a new peer connects.
func (ml *MetricLog) PeerConnected(peerURL string) {
	log.Printf("peer [%s] connected", peerURL)
}

// PeerDisconnected is called if a peer disconnects.
func (ml *MetricLog) PeerDisconnected(peerURL string) {
	log.Printf("peer [%s] disconnected", peerURL)
}
