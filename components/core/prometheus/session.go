package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type sessionCollector struct {
	core      string
	collector metric.Reader

	totalDesc  *prometheus.Desc
	activeDesc *prometheus.Desc
	rxDesc     *prometheus.Desc
	txDesc     *prometheus.Desc
}

func NewSessionCollector(core string, c metric.Reader) prometheus.Collector {
	return &sessionCollector{
		core:      core,
		collector: c,
		totalDesc: prometheus.NewDesc(
			"session_total",
			"Total number of sessions by collector",
			[]string{"core", "collector"}, nil),
		activeDesc: prometheus.NewDesc(
			"session_active",
			"Current number of active sessions by collector",
			[]string{"core", "collector"}, nil),
		rxDesc: prometheus.NewDesc(
			"session_rx_bytes",
			"Total received bytes by collector",
			[]string{"core", "collector"}, nil),
		txDesc: prometheus.NewDesc(
			"session_tx_bytes",
			"Total sent bytes by collector",
			[]string{"core", "collector"}, nil),
	}
}

func (c *sessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalDesc
	ch <- c.activeDesc
	ch <- c.rxDesc
	ch <- c.txDesc
}

func (c *sessionCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("session_total"),
		metric.NewPattern("session_active"),
		metric.NewPattern("session_rxbytes"),
		metric.NewPattern("session_txbytes"),
	})

	for _, m := range metrics.Values("session_total") {
		ch <- prometheus.MustNewConstMetric(c.totalDesc, prometheus.CounterValue, m.Val(), c.core, m.L("collector"))
	}

	for _, m := range metrics.Values("session_active") {
		ch <- prometheus.MustNewConstMetric(c.activeDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("collector"))
	}

	for _, m := range metrics.Values("session_rxbytes") {
		ch <- prometheus.MustNewConstMetric(c.rxDesc, prometheus.CounterValue, m.Val(), c.core, m.L("collector"))
	}

	for _, m := range metrics.Values("session_txbytes") {
		ch <- prometheus.MustNewConstMetric(c.txDesc, prometheus.CounterValue, m.Val(), c.core, m.L("collector"))
	}
}
