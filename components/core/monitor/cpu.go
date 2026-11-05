package monitor

import (
	"github.com/datarhei/core/v16/monitor/metric"
	"github.com/datarhei/core/v16/psutil"
)

type cpuCollector struct {
	ncpuDescr   *metric.Description
	systemDescr *metric.Description
	userDescr   *metric.Description
	idleDescr   *metric.Description
	otherDescr  *metric.Description

	ncpu float64
}

func NewCPUCollector() metric.Collector {
	c := &cpuCollector{
		ncpu: 1,
	}

	c.ncpuDescr = metric.NewDesc("cpu_ncpu", "Number of logical CPUs in the system", nil)
	c.systemDescr = metric.NewDesc("cpu_system", "Percentage of CPU used for the system", nil)
	c.userDescr = metric.NewDesc("cpu_user", "Percentage of CPU used for the user", nil)
	c.idleDescr = metric.NewDesc("cpu_idle", "Percentage of idle CPU", nil)
	c.otherDescr = metric.NewDesc("cpu_other", "Percentage of CPU used for other subsystems", nil)

	if ncpu, err := psutil.CPUCounts(true); err == nil {
		c.ncpu = ncpu
	}

	return c
}

func (c *cpuCollector) Stop() {}

func (c *cpuCollector) Prefix() string {
	return "cpu"
}

func (c *cpuCollector) Describe() []*metric.Description {
	return []*metric.Description{
		c.ncpuDescr,
		c.systemDescr,
		c.userDescr,
		c.idleDescr,
		c.otherDescr,
	}
}

func (c *cpuCollector) Collect() metric.Metrics {
	metrics := metric.NewMetrics()

	metrics.Add(metric.NewValue(c.ncpuDescr, c.ncpu))

	stat, err := psutil.CPUPercent()
	if err != nil {
		return metrics
	}

	metrics.Add(metric.NewValue(c.systemDescr, stat.System))
	metrics.Add(metric.NewValue(c.userDescr, stat.User))
	metrics.Add(metric.NewValue(c.idleDescr, stat.Idle))
	metrics.Add(metric.NewValue(c.otherDescr, stat.Other))

	return metrics
}
