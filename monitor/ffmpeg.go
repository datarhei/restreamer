package monitor

import (
	"github.com/datarhei/core/v16/ffmpeg"
	"github.com/datarhei/core/v16/monitor/metric"
)

type ffmpegCollector struct {
	prefix       string
	ffmpeg       ffmpeg.FFmpeg
	processDescr *metric.Description
}

func NewFFmpegCollector(f ffmpeg.FFmpeg) metric.Collector {
	c := &ffmpegCollector{
		prefix: "ffmpeg",
		ffmpeg: f,
	}

	c.processDescr = metric.NewDesc("ffmpeg_process", "Accumulated state changes of all ffmpeg processes", []string{"state"})

	return c
}

func (c *ffmpegCollector) Prefix() string {
	return c.prefix
}

func (c *ffmpegCollector) Describe() []*metric.Description {
	return []*metric.Description{
		c.processDescr,
	}
}

func (c *ffmpegCollector) Collect() metric.Metrics {
	metrics := metric.NewMetrics()

	states := c.ffmpeg.States()

	metrics.Add(metric.NewValue(c.processDescr, float64(states.Finished), "finished"))
	metrics.Add(metric.NewValue(c.processDescr, float64(states.Starting), "starting"))
	metrics.Add(metric.NewValue(c.processDescr, float64(states.Running), "running"))
	metrics.Add(metric.NewValue(c.processDescr, float64(states.Finishing), "finishing"))
	metrics.Add(metric.NewValue(c.processDescr, float64(states.Failed), "failed"))
	metrics.Add(metric.NewValue(c.processDescr, float64(states.Killed), "killed"))

	return metrics
}

func (c *ffmpegCollector) Stop() {}
