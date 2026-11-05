package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type memCollector struct {
	core      string
	collector metric.Reader

	memLimitDesc *prometheus.Desc
	memFreeDesc  *prometheus.Desc
}

func NewMemCollector(core string, c metric.Reader) prometheus.Collector {
	return &memCollector{
		core:      core,
		collector: c,
		memLimitDesc: prometheus.NewDesc(
			"mem_total_bytes",
			"Total available memory in bytes",
			[]string{"core"}, nil),
		memFreeDesc: prometheus.NewDesc(
			"mem_free_bytes",
			"Free memory in bytes",
			[]string{"core"}, nil),
	}
}

func (c *memCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.memLimitDesc
	ch <- c.memFreeDesc
}

func (c *memCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("mem_total"),
		metric.NewPattern("mem_free"),
	})

	for _, m := range metrics.Values("mem_total") {
		ch <- prometheus.MustNewConstMetric(c.memLimitDesc, prometheus.GaugeValue, m.Val(), c.core)
	}

	for _, m := range metrics.Values("mem_free") {
		ch <- prometheus.MustNewConstMetric(c.memFreeDesc, prometheus.GaugeValue, m.Val(), c.core)
	}
}
