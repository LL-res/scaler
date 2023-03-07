package collector

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

type Collector interface {
}

func main() {
	// 连接到 Prometheus 服务器
	client, err := api.NewClient(api.Config{
		Address: "http://prometheus-server:9090",
	})
	if err != nil {
		panic(err)
	}

	// 实例化一个 V1API 客户端
	v1api := v1.NewAPI(client)

	// 定义查询表达式
	query := "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"

	// 设置查询的时间范围
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 执行查询
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		panic(err)
	}
	if len(warnings) > 0 {
		fmt.Println("Warnings:", warnings)
	}

	// 处理查询结果
	vector := result.(model.Vector)
	for _, sample := range vector {
		fmt.Printf("CPU Usage: %v%%\n", sample.Value)
	}
}
