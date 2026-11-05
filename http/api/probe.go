package api

import (
	"encoding/json"

	"github.com/datarhei/core/v16/restream/app"
)

// ProbeIO represents a stream of a probed file
type ProbeIO struct {
	// common
	Address  string      `json:"url"`
	Format   string      `json:"format"`
	Index    uint64      `json:"index" format:"uint64"`
	Stream   uint64      `json:"stream" format:"uint64"`
	Language string      `json:"language"`
	Type     string      `json:"type"`
	Codec    string      `json:"codec"`
	Coder    string      `json:"coder"`
	Bitrate  json.Number `json:"bitrate_kbps" swaggertype:"number" jsonschema:"type=number"`
	Duration json.Number `json:"duration_sec"  swaggertype:"number" jsonschema:"type=number"`

	// video
	FPS    json.Number `json:"fps" swaggertype:"number" jsonschema:"type=number"`
	Pixfmt string      `json:"pix_fmt"`
	Width  uint64      `json:"width" format:"uint64"`
	Height uint64      `json:"height" format:"uint64"`

	// audio
	Sampling uint64 `json:"sampling_hz" format:"uint64"`
	Layout   string `json:"layout"`
	Channels uint64 `json:"channels" format:"uint64"`
}

func (i *ProbeIO) Unmarshal(io *app.ProbeIO) {
	if io == nil {
		return
	}

	i.Address = io.Address
	i.Format = io.Format
	i.Index = io.Index
	i.Stream = io.Stream
	i.Language = io.Language
	i.Type = io.Type
	i.Codec = io.Codec
	i.Coder = io.Coder
	i.Bitrate = toNumber(io.Bitrate)
	i.Duration = toNumber(io.Duration)

	i.FPS = toNumber(io.FPS)
	i.Pixfmt = io.Pixfmt
	i.Width = io.Width
	i.Height = io.Height

	i.Sampling = io.Sampling
	i.Layout = io.Layout
	i.Channels = io.Channels
}

// Probe represents the result of probing a file. It has a list of detected streams
// and a list of log lone from the probe process.
type Probe struct {
	Streams []ProbeIO `json:"streams"`
	Log     []string  `json:"log"`
}

// Unmarshal converts a restreamer Probe to a Probe in API representation
func (probe *Probe) Unmarshal(p *app.Probe) {
	if p == nil {
		return
	}

	probe.Streams = make([]ProbeIO, len(p.Streams))
	probe.Log = make([]string, len(p.Log))

	for i, io := range p.Streams {
		probe.Streams[i].Unmarshal(&io)
	}

	copy(probe.Log, p.Log)
}
