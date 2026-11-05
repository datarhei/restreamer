package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	Register(cs prometheus.Collector) error
	UnregisterAll()
	Reader
}

type Reader interface {
	HTTPHandler() http.Handler
}

type metrics struct {
	registry   *prometheus.Registry
	collectors []prometheus.Collector
}

func New() Metrics {
	m := &metrics{
		registry: prometheus.NewRegistry(),
	}

	return m
}

func (m *metrics) Register(cs prometheus.Collector) error {
	if err := m.registry.Register(cs); err != nil {
		return err
	}

	m.collectors = append(m.collectors, cs)

	return nil
}

func (m *metrics) UnregisterAll() {
	for _, cs := range m.collectors {
		m.registry.Unregister(cs)
	}
}

func (m *metrics) HTTPHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
}
