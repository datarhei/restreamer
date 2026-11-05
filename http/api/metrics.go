package api

import (
	"fmt"
	"time"

	"github.com/datarhei/core/v16/monitor"
)

type MetricsDescription struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
}

type MetricsQueryMetric struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type MetricsQuery struct {
	Timerange int64                `json:"timerange_sec" format:"int64"`
	Interval  int64                `json:"interval_sec" format:"int64"`
	Metrics   []MetricsQueryMetric `json:"metrics"`
}

type MetricsResponseMetric struct {
	Name   string                 `json:"name"`
	Labels map[string]string      `json:"labels"`
	Values []MetricsResponseValue `json:"values"`
}

type MetricsResponseValue struct {
	TS    time.Time `json:"ts"`
	Value float64   `json:"value"`
}

// MarshalJSON marshals a MetricsResponseValue to JSON
func (v MetricsResponseValue) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("[%d,", v.TS.Unix())

	if v.Value == float64(int64(v.Value)) {
		s += fmt.Sprintf("%.0f", v.Value) // 0 decimal if integer
	} else {
		s += fmt.Sprintf("%.3f", v.Value) // max. 3 decimal if float
	}

	s += "]"

	return []byte(s), nil
}

type MetricsResponse struct {
	Timerange int64                   `json:"timerange_sec" format:"int64"`
	Interval  int64                   `json:"interval_sec" format:"int64"`
	Metrics   []MetricsResponseMetric `json:"metrics"`
}

func (m *MetricsResponse) Unmarshal(data []monitor.HistoryMetrics, timerange, interval time.Duration) {
	series := make(map[string]MetricsResponseMetric)

	for _, d := range data {
		if d.Metrics == nil {
			continue
		}

		for _, v := range d.Metrics.All() {
			hash := v.Hash()

			if _, ok := series[hash]; !ok {
				series[hash] = MetricsResponseMetric{
					Name:   v.Name(),
					Labels: v.Labels(),
					Values: []MetricsResponseValue{},
				}
			}

			k := series[hash]

			k.Values = append(k.Values, MetricsResponseValue{
				TS:    d.TS,
				Value: v.Val(),
			})

			series[hash] = k
		}
	}

	m.Metrics = make([]MetricsResponseMetric, len(series))

	i := 0
	for _, metric := range series {
		m.Metrics[i] = metric
		i++
	}

	m.Timerange = int64(timerange.Seconds())
	m.Interval = int64(interval.Seconds())
}
