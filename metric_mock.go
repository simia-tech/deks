package edkvs

// MetricMock defines a metric mock.
type MetricMock struct{}

// NewMetricMock returns a new metric mock.
func NewMetricMock() *MetricMock {
	return &MetricMock{}
}

// CountChanged is called if the number of value or deleted values has changed.
func (mm *MetricMock) CountChanged(valueCount, deletedCount int) {}
