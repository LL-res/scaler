package prometheus_collector

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	constant "scaler/const"
	"scaler/internal/collector"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type Promc struct {
	collector.CollectorBase
	//prometheus client
	client api.Client
}

func New() collector.Collector {
	metricQL := make(map[string]string, 0)
	metricQL[constant.MetricCPUUsage] = "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"
	res := new(Promc)
	res.MetricQL = metricQL
	res.Data = collector.MetricSeries{
		Metrics:   make([]float64, 0, 100),
		TimeStamp: make([]time.Time, 0, 100),
		Length:    0,
	}
	return res
}
func (p *Promc) SetServerAddress(url string) error {
	p.ServerAddress = url
	client, err := api.NewClient(api.Config{
		Address: p.ServerAddress,
	})
	if err != nil {
		return err
	}
	p.client = client
	// TODO test if can get
	return nil
}

func (p *Promc) SetMetricType(metricType string) error {
	if _, ok := p.MetricQL[metricType]; !ok {
		return errors.New("undefined metric type")
	}
	if p.CurrentMetric != "" {
		//TODO handle current data
		p.Data = collector.MetricSeries{
			Metrics:   make([]float64, 0, 100),
			TimeStamp: make([]time.Time, 0, 100),
		}
	}
	p.CurrentMetric = metricType
	return nil
}

func (p *Promc) SetCapacity(capacity int) {
	if p.Data.Length > capacity {
		// TODO handle data
	}
	p.Capacity = capacity
}

func (p *Promc) ListMetricTypes() []string {
	result := make([]string, 0, len(p.MetricQL))
	for k := range p.MetricQL {
		result = append(result, k)
	}
	return result
}

func (p *Promc) GetMetrics() error {

	// 实例化一个 V1API 客户端
	v1api := v1.NewAPI(p.client)

	// 定义查询表达式
	query := p.MetricQL[p.CurrentMetric]

	// 设置查询的时间范围
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 执行查询
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		log.Log.WithName("prometheus_collector").Info(fmt.Sprintf("warning: %v", warnings))
	}

	// 处理查询结果
	vector := result.(model.Vector)
	for _, sample := range vector {
		if p.Data.Length == p.Capacity {
			// TODO handle data
			p.Data.Length = 0
		}
		p.Data.Metrics = append(p.Data.Metrics, float64(sample.Value))
		p.Data.TimeStamp = append(p.Data.TimeStamp, sample.Timestamp.Time())
		p.Data.Length++
	}
	return nil
}

func (p *Promc) AddCustomMetrics(name, query string) {
	p.MetricQL[name] = query
}
func (p *Promc) SyncData() error {
	return nil
}
