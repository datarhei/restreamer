package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type uptimeCollector struct {
	core      string
	collector metric.Reader

	uptimeDesc *prometheus.Desc
}

func NewUptimeCollector(core string, c metric.Reader) prometheus.Collector {
	return &uptimeCollector{
		core:      core,
		collector: c,
		uptimeDesc: prometheus.NewDesc(
			"uptime_seconds",
			"Number of seconds the core is up",
			[]string{"core"}, nil),
	}
}

func (c *uptimeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptimeDesc
}

func (c *uptimeCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("uptime_uptime"),
	})

	for _, m := range metrics.Values("uptime_uptime") {
		ch <- prometheus.MustNewConstMetric(c.uptimeDesc, prometheus.CounterValue, m.Val(), c.core)
	}
}
