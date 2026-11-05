package monitor

import (
	"time"

	"github.com/datarhei/core/v16/monitor/metric"
)

type uptimeCollector struct {
	t           time.Time
	uptimeDescr *metric.Description
}

func NewUptimeCollector() metric.Collector {
	c := &uptimeCollector{
		t: time.Now(),
	}

	c.uptimeDescr = metric.NewDesc("uptime_uptime", "Current uptime in seconds", nil)

	return c
}

func (c *uptimeCollector) Prefix() string {
	return "uptime"
}

func (c *uptimeCollector) Describe() []*metric.Description {
	return []*metric.Description{
		c.uptimeDescr,
	}
}

func (c *uptimeCollector) Collect() metric.Metrics {
	uptime := time.Since(c.t).Seconds()

	metrics := metric.NewMetrics()

	metrics.Add(metric.NewValue(c.uptimeDescr, uptime))

	return metrics
}

func (c *uptimeCollector) Stop() {}
