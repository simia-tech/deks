package edkvs

import "log"

// MetricLog defines a metric log.
type MetricLog struct{}

// NewMetricLog returns a new metric log.
func NewMetricLog() *MetricLog {
	return &MetricLog{}
}

// CountChanged is called if the number of value or deleted values has changed.
func (mm *MetricLog) CountChanged(valueCount, deletedCount int) {
	log.Printf("count changed: values = %d / deleted = %d", valueCount, deletedCount)
}
