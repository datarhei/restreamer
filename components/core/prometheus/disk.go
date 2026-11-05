package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type diskCollector struct {
	core      string
	collector metric.Reader

	diskTotalDesc *prometheus.Desc
	diskUsageDesc *prometheus.Desc
}

func NewDiskCollector(core string, c metric.Reader) prometheus.Collector {
	return &diskCollector{
		core:      core,
		collector: c,
		diskTotalDesc: prometheus.NewDesc(
			"disk_total_bytes",
			"disk_total_bytes",
			[]string{"core", "path"}, nil),
		diskUsageDesc: prometheus.NewDesc(
			"disk_usage_bytes",
			"disk_usage_bytes",
			[]string{"core", "path"}, nil),
	}
}

func (c *diskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.diskTotalDesc
	ch <- c.diskUsageDesc
}

func (c *diskCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("disk_total"),
		metric.NewPattern("disk_usage"),
	})

	for _, m := range metrics.Values("disk_total") {
		ch <- prometheus.MustNewConstMetric(c.diskTotalDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("path"))
	}

	for _, m := range metrics.Values("disk_usage") {
		ch <- prometheus.MustNewConstMetric(c.diskUsageDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("path"))
	}
}
