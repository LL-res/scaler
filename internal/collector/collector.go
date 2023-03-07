package collector

import (
	"time"
)

type MetricSeries struct {
	Metrics   []float64
	TimeStamp []time.Time
	Length    int
}

type Collector interface {
	SetServerAddress(url string) error
	SetMetricType(metricType string) error
	SetCapacity(capacity int)
	ListMetricTypes() []string
	GetMetrics() error
	AddCustomMetrics(name, query string)
	SyncData() error
}
type CollectorBase struct {
	//contain two slices which have the same length
	Data MetricSeries
	//key: the name of  supported metric type,value: the promql to get key metric type
	MetricQL map[string]string
	//if the data length is larger than this value, the data will be synced to database
	Capacity int
	//prometheus server url
	ServerAddress string
	//the metric can be got when GetMetrics is called
	CurrentMetric string
}
