package monitor

import (
	"github.com/datarhei/core/v16/monitor/metric"
	"github.com/datarhei/core/v16/psutil"
)

type memCollector struct {
	totalDescr *metric.Description
	freeDescr  *metric.Description
}

func NewMemCollector() metric.Collector {
	c := &memCollector{}

	c.totalDescr = metric.NewDesc("mem_total", "Total available memory in bytes", nil)
	c.freeDescr = metric.NewDesc("mem_free", "Free memory in bytes", nil)

	return c
}

func (c *memCollector) Prefix() string {
	return "mem"
}

func (c *memCollector) Describe() []*metric.Description {
	return []*metric.Description{
		c.totalDescr,
		c.freeDescr,
	}
}

func (c *memCollector) Collect() metric.Metrics {
	metrics := metric.NewMetrics()

	stat, err := psutil.VirtualMemory()
	if err != nil {
		return metrics
	}

	metrics.Add(metric.NewValue(c.totalDescr, float64(stat.Total)))
	metrics.Add(metric.NewValue(c.freeDescr, float64(stat.Available)))

	return metrics
}

func (c *memCollector) Stop() {}
