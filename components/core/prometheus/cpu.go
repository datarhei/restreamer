package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type cpuCollector struct {
	core      string
	collector metric.Reader

	cpuSystemTimeDesc *prometheus.Desc
	cpuUserTimeDesc   *prometheus.Desc
	cpuIdleTimeDesc   *prometheus.Desc
	cpuOtherTimeDesc  *prometheus.Desc
}

func NewCPUCollector(core string, c metric.Reader) prometheus.Collector {
	return &cpuCollector{
		core:      core,
		collector: c,
		cpuSystemTimeDesc: prometheus.NewDesc(
			"cpu_system_time_percent",
			"CPU system time in percent",
			[]string{"core"}, nil),
		cpuUserTimeDesc: prometheus.NewDesc(
			"cpu_user_time_percent",
			"CPU user time in percent",
			[]string{"core"}, nil),
		cpuIdleTimeDesc: prometheus.NewDesc(
			"cpu_idle_time_percent",
			"CPU idle time in percent",
			[]string{"core"}, nil),
		cpuOtherTimeDesc: prometheus.NewDesc(
			"cpu_other_time_percent",
			"CPU other time in percent",
			[]string{"core"}, nil),
	}
}

func (c *cpuCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuSystemTimeDesc
	ch <- c.cpuUserTimeDesc
	ch <- c.cpuIdleTimeDesc
	ch <- c.cpuOtherTimeDesc
}

func (c *cpuCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("cpu_system"),
		metric.NewPattern("cpu_user"),
		metric.NewPattern("cpu_idle"),
		metric.NewPattern("cpu_other"),
	})

	ch <- prometheus.MustNewConstMetric(c.cpuSystemTimeDesc, prometheus.GaugeValue, metrics.Value("cpu_system").Val(), c.core)
	ch <- prometheus.MustNewConstMetric(c.cpuUserTimeDesc, prometheus.GaugeValue, metrics.Value("cpu_user").Val(), c.core)
	ch <- prometheus.MustNewConstMetric(c.cpuIdleTimeDesc, prometheus.GaugeValue, metrics.Value("cpu_idle").Val(), c.core)
	ch <- prometheus.MustNewConstMetric(c.cpuOtherTimeDesc, prometheus.GaugeValue, metrics.Value("cpu_other").Val(), c.core)
}
