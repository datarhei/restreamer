package ffmpeg

import (
	"github.com/datarhei/core/v16/session"
)

type wrappedCollector struct {
	session.Collector

	prefix    string
	reference string
}

func NewWrappedCollector(prefix, reference string, collector session.Collector) session.Collector {
	w := &wrappedCollector{
		prefix:    prefix,
		reference: reference,
		Collector: collector,
	}

	return w
}

func (w *wrappedCollector) Register(id, reference, location, peer string) {
	w.Collector.Register(w.prefix+id, w.reference, location, peer)
}

func (w *wrappedCollector) Activate(id string) bool {
	return w.Collector.Activate(w.prefix + id)
}

func (w *wrappedCollector) RegisterAndActivate(id, reference, location, peer string) {
	w.Collector.RegisterAndActivate(w.prefix+id, w.reference, location, peer)
}

func (w *wrappedCollector) Extra(id, extra string) {
	w.Collector.Extra(w.prefix+id, extra)
}

func (w *wrappedCollector) Unregister(id string) {
	w.Collector.Unregister(w.prefix + id)
}

func (w *wrappedCollector) Ingress(id string, size int64) {
	w.Collector.Ingress(w.prefix+id, size)
}

func (w *wrappedCollector) Egress(id string, size int64) {
	w.Collector.Egress(w.prefix+id, size)
}

func (w *wrappedCollector) IsKnownSession(id string) bool {
	return w.Collector.IsKnownSession(w.prefix + id)
}

func (w *wrappedCollector) SessionTopIngressBitrate(id string) float64 {
	return w.Collector.SessionTopIngressBitrate(w.prefix + id)
}

func (w *wrappedCollector) SessionTopEgressBitrate(id string) float64 {
	return w.Collector.SessionTopEgressBitrate(w.prefix + id)
}

func (w *wrappedCollector) SessionSetTopIngressBitrate(id string, bitrate float64) {
	w.Collector.SessionSetTopIngressBitrate(w.prefix+id, bitrate)
}

func (w *wrappedCollector) SessionSetTopEgressBitrate(id string, bitrate float64) {
	w.Collector.SessionSetTopEgressBitrate(w.prefix+id, bitrate)
}
