package probe

import (
	"github.com/datarhei/core/v16/restream/app"
)

type probeIO struct {
	// common
	Address  string  `json:"url"`
	Format   string  `json:"format"`
	Index    uint64  `json:"index"`
	Stream   uint64  `json:"stream"`
	Language string  `json:"language"`
	Type     string  `json:"type"`
	Codec    string  `json:"codec"`
	Coder    string  `json:"coder"`
	Bitrate  float64 `json:"bitrate_kbps"`
	Duration float64 `json:"duration_sec"`

	// video
	FPS    float64 `json:"fps"`
	Pixfmt string  `json:"pix_fmt"`
	Width  uint64  `json:"width"`
	Height uint64  `json:"height"`

	// audio
	Sampling uint64 `json:"sampling_hz"`
	Layout   string `json:"layout"`
	Channels uint64 `json:"channels"`
}

func (io *probeIO) export() app.ProbeIO {
	return app.ProbeIO{
		Address:  io.Address,
		Format:   io.Format,
		Index:    io.Index,
		Stream:   io.Stream,
		Language: io.Language,
		Type:     io.Type,
		Codec:    io.Codec,
		Coder:    io.Coder,
		Bitrate:  io.Bitrate,
		Duration: io.Duration,
		FPS:      io.FPS,
		Pixfmt:   io.Pixfmt,
		Width:    io.Width,
		Height:   io.Height,
		Sampling: io.Sampling,
		Layout:   io.Layout,
		Channels: io.Channels,
	}
}
