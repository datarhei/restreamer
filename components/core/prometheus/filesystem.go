package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type filesystemCollector struct {
	core      string
	collector metric.Reader

	fsLimitDesc *prometheus.Desc
	fsUsageDesc *prometheus.Desc
	fsFilesDesc *prometheus.Desc
}

func NewFilesystemCollector(core string, c metric.Reader) prometheus.Collector {
	return &filesystemCollector{
		core:      core,
		collector: c,
		fsLimitDesc: prometheus.NewDesc(
			"filesystem_limit_bytes",
			"filesystem_limit_bytes",
			[]string{"core", "name"}, nil),
		fsUsageDesc: prometheus.NewDesc(
			"filesystem_usage_bytes",
			"filesystem_usage_bytes",
			[]string{"core", "name"}, nil),
		fsFilesDesc: prometheus.NewDesc(
			"filesystem_files",
			"filesystem_files",
			[]string{"core", "name"}, nil),
	}
}

func (c *filesystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.fsLimitDesc
	ch <- c.fsUsageDesc
	ch <- c.fsFilesDesc
}

func (c *filesystemCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("filesystem_limit"),
		metric.NewPattern("filesystem_usage"),
		metric.NewPattern("filesystem_files"),
	})

	for _, m := range metrics.Values("filesystem_limit") {
		ch <- prometheus.MustNewConstMetric(c.fsLimitDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("name"))
	}

	for _, m := range metrics.Values("filesystem_usage") {
		ch <- prometheus.MustNewConstMetric(c.fsUsageDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("name"))
	}

	for _, m := range metrics.Values("filesystem_files") {
		ch <- prometheus.MustNewConstMetric(c.fsFilesDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("name"))
	}
}
