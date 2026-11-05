package api

import (
	"encoding/json"
	"fmt"

	"github.com/datarhei/core/v16/restream/app"
)

type ProgressIOFramerate struct {
	Min     json.Number `json:"min" swaggertype:"number" jsonschema:"type=number"`
	Max     json.Number `json:"max" swaggertype:"number" jsonschema:"type=number"`
	Average json.Number `json:"avg" swaggertype:"number" jsonschema:"type=number"`
}

// ProgressIO represents the progress of an ffmpeg input or output
type ProgressIO struct {
	ID      string `json:"id" jsonschema:"minLength=1"`
	Address string `json:"address" jsonschema:"minLength=1"`

	// General
	Index     uint64              `json:"index" format:"uint64"`
	Stream    uint64              `json:"stream" format:"uint64"`
	Format    string              `json:"format"`
	Type      string              `json:"type"`
	Codec     string              `json:"codec"`
	Coder     string              `json:"coder"`
	Frame     uint64              `json:"frame" format:"uint64"`
	Keyframe  uint64              `json:"keyframe" format:"uint64"`
	Framerate ProgressIOFramerate `json:"framerate"`
	FPS       json.Number         `json:"fps" swaggertype:"number" jsonschema:"type=number"`
	Packet    uint64              `json:"packet" format:"uint64"`
	PPS       json.Number         `json:"pps" swaggertype:"number" jsonschema:"type=number"`
	Size      uint64              `json:"size_kb" format:"uint64"`                                    // kbytes
	Bitrate   json.Number         `json:"bitrate_kbit" swaggertype:"number" jsonschema:"type=number"` // kbit/s
	Extradata uint64              `json:"extradata_size_bytes" format:"uint64"`                       // bytes

	// Video
	Pixfmt    string      `json:"pix_fmt,omitempty"`
	Quantizer json.Number `json:"q,omitempty" swaggertype:"number" jsonschema:"type=number"`
	Width     uint64      `json:"width,omitempty" format:"uint64"`
	Height    uint64      `json:"height,omitempty" format:"uint64"`

	// Audio
	Sampling uint64 `json:"sampling_hz,omitempty" format:"uint64"`
	Layout   string `json:"layout,omitempty"`
	Channels uint64 `json:"channels,omitempty" format:"uint64"`

	// avstream
	AVstream *AVstream `json:"avstream"`
}

// Unmarshal converts a restreamer ProgressIO to a ProgressIO in API representation
func (i *ProgressIO) Unmarshal(io *app.ProgressIO) {
	if io == nil {
		return
	}

	i.ID = io.ID
	i.Address = io.Address
	i.Index = io.Index
	i.Stream = io.Stream
	i.Format = io.Format
	i.Type = io.Type
	i.Codec = io.Codec
	i.Coder = io.Coder
	i.Frame = io.Frame
	i.Keyframe = io.Keyframe
	i.Framerate.Min = json.Number(fmt.Sprintf("%.3f", io.Framerate.Min))
	i.Framerate.Max = json.Number(fmt.Sprintf("%.3f", io.Framerate.Max))
	i.Framerate.Average = json.Number(fmt.Sprintf("%.3f", io.Framerate.Average))
	i.FPS = json.Number(fmt.Sprintf("%.3f", io.FPS))
	i.Packet = io.Packet
	i.PPS = json.Number(fmt.Sprintf("%.3f", io.PPS))
	i.Size = io.Size / 1024
	i.Bitrate = json.Number(fmt.Sprintf("%.3f", io.Bitrate/1024))
	i.Extradata = io.Extradata
	i.Pixfmt = io.Pixfmt
	i.Quantizer = json.Number(fmt.Sprintf("%.3f", io.Quantizer))
	i.Width = io.Width
	i.Height = io.Height
	i.Sampling = io.Sampling
	i.Layout = io.Layout
	i.Channels = io.Channels

	if io.AVstream != nil {
		i.AVstream = &AVstream{}
		i.AVstream.Unmarshal(io.AVstream)
	}
}

// Progress represents the progress of an ffmpeg process
type Progress struct {
	Input     []ProgressIO `json:"inputs"`
	Output    []ProgressIO `json:"outputs"`
	Frame     uint64       `json:"frame" format:"uint64"`
	Packet    uint64       `json:"packet" format:"uint64"`
	FPS       json.Number  `json:"fps" swaggertype:"number" jsonschema:"type=number"`
	Quantizer json.Number  `json:"q" swaggertype:"number" jsonschema:"type=number"`
	Size      uint64       `json:"size_kb" format:"uint64"` // kbytes
	Time      json.Number  `json:"time" swaggertype:"number" jsonschema:"type=number"`
	Bitrate   json.Number  `json:"bitrate_kbit" swaggertype:"number" jsonschema:"type=number"` // kbit/s
	Speed     json.Number  `json:"speed" swaggertype:"number" jsonschema:"type=number"`
	Drop      uint64       `json:"drop" format:"uint64"`
	Dup       uint64       `json:"dup" format:"uint64"`
}

// Unmarshal converts a restreamer Progress to a Progress in API representation
func (progress *Progress) Unmarshal(p *app.Progress) {
	progress.Input = []ProgressIO{}
	progress.Output = []ProgressIO{}

	if p == nil {
		return
	}

	progress.Input = make([]ProgressIO, len(p.Input))
	progress.Output = make([]ProgressIO, len(p.Output))
	progress.Frame = p.Frame
	progress.Packet = p.Packet
	progress.FPS = toNumber(p.FPS)
	progress.Quantizer = toNumber(p.Quantizer)
	progress.Size = p.Size / 1024
	progress.Time = toNumber(p.Time)
	progress.Bitrate = toNumber(p.Bitrate / 1024)
	progress.Speed = toNumber(p.Speed)
	progress.Drop = p.Drop
	progress.Dup = p.Dup

	for i, io := range p.Input {
		progress.Input[i].Unmarshal(&io)
	}

	for i, io := range p.Output {
		progress.Output[i].Unmarshal(&io)
	}
}
