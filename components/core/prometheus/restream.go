package prometheus

import (
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/prometheus/client_golang/prometheus"
)

type restreamCollector struct {
	core      string
	collector metric.Reader

	ffmpegProcessDesc       *prometheus.Desc
	ffmpegProcessStatesDesc *prometheus.Desc
	ffmpegProcessIODesc     *prometheus.Desc
	ffmpegStatesDesc        *prometheus.Desc
	ffmpegStatesTotalDesc   *prometheus.Desc
}

func NewRestreamCollector(core string, c metric.Reader) prometheus.Collector {
	return &restreamCollector{
		core:      core,
		collector: c,
		ffmpegProcessDesc: prometheus.NewDesc(
			"ffmpeg_process",
			"General stats per process",
			[]string{"core", "process", "state", "order", "name"}, nil),
		ffmpegProcessStatesDesc: prometheus.NewDesc(
			"ffmpeg_process_states",
			"Accumulated states per process",
			[]string{"core", "process", "state"}, nil),
		ffmpegProcessIODesc: prometheus.NewDesc(
			"ffmpeg_process_io",
			"Stats per input and output of a process",
			[]string{"core", "process", "type", "id", "index", "stream", "media", "name"}, nil),
		ffmpegStatesDesc: prometheus.NewDesc(
			"ffmpeg_states",
			"Current process states",
			[]string{"core", "state"}, nil),
		ffmpegStatesTotalDesc: prometheus.NewDesc(
			"ffmpeg_states_total",
			"Accumulated process states",
			[]string{"core", "state"}, nil),
	}
}

func (c *restreamCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ffmpegProcessDesc
	ch <- c.ffmpegProcessStatesDesc
	ch <- c.ffmpegProcessIODesc
	ch <- c.ffmpegStatesDesc
	ch <- c.ffmpegStatesTotalDesc
}

func (c *restreamCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.collector.Collect([]metric.Pattern{
		metric.NewPattern("restream_process"),
		metric.NewPattern("restream_process_states"),
		metric.NewPattern("restream_io"),
		metric.NewPattern("ffmpeg_process"),
	})

	for _, m := range metrics.Values("restream_process") {
		ch <- prometheus.MustNewConstMetric(c.ffmpegProcessDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("processid"), m.L("state"), m.L("order"), m.L("name"))
	}

	for _, m := range metrics.Values("restream_process_states") {
		ch <- prometheus.MustNewConstMetric(c.ffmpegProcessStatesDesc, prometheus.GaugeValue, m.Val(), c.core, m.L("processid"), m.L("state"))
	}

	for _, m := range metrics.Values("restream_io") {
		ch <- prometheus.MustNewConstMetric(c.ffmpegProcessIODesc, prometheus.GaugeValue, m.Val(), c.core, m.L("processid"), m.L("type"), m.L("id"), m.L("index"), m.L("stream"), m.L("media"), m.L("name"))
	}

	states := map[string]float64{
		"failed":    0,
		"finished":  0,
		"finishing": 0,
		"killed":    0,
		"running":   0,
		"starting":  0,
	}

	for _, processid := range metrics.Labels("restream_process", "processid") {
		s := metrics.Value("restream_process", "processid", "^"+processid+"$").L("state")

		if _, ok := states[s]; !ok {
			continue
		}

		states[s]++
	}

	for state, value := range states {
		ch <- prometheus.MustNewConstMetric(c.ffmpegStatesDesc, prometheus.GaugeValue, value, c.core, state)
	}

	for _, m := range metrics.Values("ffmpeg_process") {
		ch <- prometheus.MustNewConstMetric(c.ffmpegStatesTotalDesc, prometheus.CounterValue, m.Val(), c.core, m.L("state"))
	}
}
