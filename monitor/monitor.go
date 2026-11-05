package monitor

import (
	"container/ring"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/datarhei/core/v16/monitor/metric"
)

type Reader interface {
	Collect(patterns []metric.Pattern) metric.Metrics
	Describe() []*metric.Description
}

type Monitor interface {
	Reader
	Register(c metric.Collector)
	UnregisterAll()
}

type HistoryReader interface {
	Reader
	History(timerange, interval time.Duration, patterns []metric.Pattern) []HistoryMetrics
	Resolution() (timerange, interval time.Duration)
}

type HistoryMonitor interface {
	HistoryReader
	Register(c metric.Collector)
	UnregisterAll()
}

type Config struct{}

type monitor struct {
	lock       sync.RWMutex
	collectors map[string]metric.Collector
}

func New(config Config) Monitor {
	m := &monitor{
		collectors: map[string]metric.Collector{},
	}

	return m
}

func (m *monitor) Register(c metric.Collector) {
	if c == nil {
		return
	}

	descriptors := c.Describe()

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, d := range descriptors {
		m.collectors[d.Name()] = c
	}
}

func (m *monitor) Collect(patterns []metric.Pattern) metric.Metrics {
	metrics := metric.NewMetrics()
	prefixes := make(map[metric.Collector][]metric.Pattern)

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, pattern := range patterns {
		c, ok := m.collectors[pattern.Name()]
		if !ok {
			continue
		}

		prefixes[c] = append(prefixes[c], pattern)
	}

	for c, patterns := range prefixes {
		vs := c.Collect()

		for _, v := range vs.All() {
			if v.Match(patterns) {
				metrics.Add(v)
			}
		}
	}

	return metrics
}

func (m *monitor) Describe() []*metric.Description {
	descriptors := []*metric.Description{}
	collectors := map[metric.Collector]struct{}{}

	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, c := range m.collectors {
		if _, ok := collectors[c]; ok {
			continue
		}

		collectors[c] = struct{}{}

		descriptors = append(descriptors, c.Describe()...)
	}

	return descriptors
}

func (m *monitor) UnregisterAll() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, collector := range m.collectors {
		collector.Stop()
	}

	m.collectors = make(map[string]metric.Collector)
}

type historyMonitor struct {
	monitor Monitor

	enable    bool
	timerange time.Duration
	interval  time.Duration

	metrics *ring.Ring
	lock    sync.RWMutex

	stopTicker chan struct{}

	patterns []metric.Pattern
}

type HistoryMetrics struct {
	TS      time.Time
	Metrics metric.Metrics
}

type HistoryConfig struct {
	Config    Config
	Enable    bool
	Timerange time.Duration
	Interval  time.Duration
}

func NewHistory(config HistoryConfig) (HistoryMonitor, error) {
	m := &historyMonitor{
		monitor: New(config.Config),
		enable:  config.Enable,
	}

	if m.enable {
		if config.Interval <= 0 {
			return nil, fmt.Errorf("the interval must be greater than 0")
		}

		if config.Timerange <= 0 {
			return nil, fmt.Errorf("the timeframe must be greater than 0")
		}

		if config.Interval > config.Timerange {
			return nil, fmt.Errorf("the interval has to be shorter than the frame")
		}

		n := config.Timerange / config.Interval
		if n > math.MaxInt {
			return nil, fmt.Errorf("too many intervals")
		}

		if n == 0 {
			n = 1
		}

		m.timerange = config.Timerange
		m.interval = config.Interval

		m.metrics = ring.New(int(n))
		m.stopTicker = make(chan struct{})

	}

	return m, nil
}

func (m *historyMonitor) tick() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopTicker:
			return
		case ts := <-ticker.C:
			m.collectAll(ts)
		}
	}
}

func (m *historyMonitor) collectAll(ts time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()

	metrics := m.Collect(m.patterns)

	m.metrics.Value = &HistoryMetrics{
		TS:      ts,
		Metrics: metrics,
	}
	m.metrics = m.metrics.Next()
}

func (m *historyMonitor) Register(c metric.Collector) {
	m.monitor.Register(c)

	if !m.enable {
		return
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, d := range c.Describe() {
		m.patterns = append(m.patterns, metric.NewPattern(d.Name()))
	}

	if len(m.patterns) == 1 {
		// start collecting metrics with the first registered collector
		go m.tick()
	}
}

func (m *historyMonitor) Collect(patterns []metric.Pattern) metric.Metrics {
	return m.monitor.Collect(patterns)
}

func (m *historyMonitor) Describe() []*metric.Description {
	return m.monitor.Describe()
}

func (m *historyMonitor) UnregisterAll() {
	m.monitor.UnregisterAll()

	if !m.enable {
		return
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.patterns = nil

	// stop collecting metrics if all collectors are unregisterd
	close(m.stopTicker)
}

func (m *historyMonitor) History(timerange, interval time.Duration, patterns []metric.Pattern) []HistoryMetrics {
	metricsList := []HistoryMetrics{}

	if !m.enable {
		return metricsList
	}

	notBefore := time.Now().Add(-timerange - m.interval)

	m.lock.Lock()

	m.metrics.Do(func(l interface{}) {
		if l == nil {
			return
		}

		historyMetrics := l.(*HistoryMetrics)

		if historyMetrics.TS.Before(notBefore) {
			return
		}

		metrics := metric.NewMetrics()

		for _, v := range historyMetrics.Metrics.All() {
			if v.Match(patterns) {
				metrics.Add(v)
			}
		}

		metricsList = append(metricsList, HistoryMetrics{
			TS:      historyMetrics.TS,
			Metrics: metrics,
		})
	})

	m.lock.Unlock()

	metricsList = m.resample(metricsList, timerange, interval)

	return metricsList
}

func (m *historyMonitor) Resolution() (timerange, interval time.Duration) {
	return m.timerange, m.interval
}

func (m *historyMonitor) resample(values []HistoryMetrics, timerange, interval time.Duration) []HistoryMetrics {
	v := []HistoryMetrics{}
	nvalues := len(values)

	if nvalues == 0 || timerange == 0 || interval == 0 {
		return v
	}

	to := time.Now()
	from := to.Add(-timerange - m.interval)

	start := values[0].TS
	end := values[nvalues-1].TS
	//startValue := values[0].Metrics
	endValue := values[nvalues-1].Metrics

	steps := int(timerange / interval)

	lastJ := 0
	for i := 0; i < steps; i++ {
		now := from.Add(time.Duration(i) * interval)

		if now.Before(start) {
			v = append(v, HistoryMetrics{
				TS:      now,
				Metrics: nil,
			})
			continue
		}

		if now.After(end) {
			v = append(v, HistoryMetrics{
				TS:      now,
				Metrics: endValue,
			})
			continue
		}

		for j := lastJ; j < nvalues-1; j++ {
			x := values[j].TS
			y := values[j+1].TS

			if (now.Equal(x) || now.After(x)) && now.Before(y) {
				v = append(v, HistoryMetrics{
					TS:      now,
					Metrics: values[j].Metrics,
				})
				lastJ = j
				break
			}
		}
	}

	return v
}
