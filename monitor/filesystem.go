package monitor

import (
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/monitor/metric"
)

type filesystemCollector struct {
	fs         fs.Filesystem
	name       string
	limitDescr *metric.Description
	usageDescr *metric.Description
	filesDescr *metric.Description
}

func NewFilesystemCollector(name string, fs fs.Filesystem) metric.Collector {
	c := &filesystemCollector{
		fs:   fs,
		name: name,
	}

	c.limitDescr = metric.NewDesc("filesystem_limit", "Total size of the filesystem in bytes, negative if unlimited", []string{"name"})
	c.usageDescr = metric.NewDesc("filesystem_usage", "Number of used bytes on the filesystem", []string{"name"})
	c.filesDescr = metric.NewDesc("filesystem_files", "Number of files on the filesystem (excluding directories)", []string{"name"})

	return c
}

func (c *filesystemCollector) Prefix() string {
	return "filesystem"
}

func (c *filesystemCollector) Describe() []*metric.Description {
	return []*metric.Description{
		c.limitDescr,
		c.usageDescr,
		c.filesDescr,
	}
}

func (c *filesystemCollector) Collect() metric.Metrics {
	size, limit := c.fs.Size()
	files := c.fs.Files()

	metrics := metric.NewMetrics()

	metrics.Add(metric.NewValue(c.limitDescr, float64(limit), c.name))
	metrics.Add(metric.NewValue(c.usageDescr, float64(size), c.name))
	metrics.Add(metric.NewValue(c.filesDescr, float64(files), c.name))

	return metrics
}

func (c *filesystemCollector) Stop() {}
