package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type netCollector struct {
	core      string
	collector metric.Reader

	netRxDesc *prometheus.Desc
	netTxDesc *prometheus.Desc
}

func NewNetCollector(core string, c metric.Reader) prometheus.Collector {
	return &netCollector{
		core:      core,
		collector: c,
		netRxDesc: prometheus.NewDesc(
			"net_rx_bytes",
			"Number of received bytes by interface",
			[]string{"core", "interface"}, nil),
		netTxDesc: prometheus.NewDesc(
			"net_tx_bytes",
			"Number of sent bytes by interface",
			[]string{"core", "interface"}, nil),
	}
}

func (c *netCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.netRxDesc
	ch <- c.netTxDesc
}

func (c *netCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("net_rx"),
		metric.NewPattern("net_tx"),
	})

	for _, m := range metrics.Values("net_rx") {
		ch <- prometheus.MustNewConstMetric(c.netRxDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("interface"))
	}

	for _, m := range metrics.Values("net_tx") {
		ch <- prometheus.MustNewConstMetric(c.netTxDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("interface"))
	}
}
