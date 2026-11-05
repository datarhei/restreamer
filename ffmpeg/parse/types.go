package parse

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/datarhei/core/v16/restream/app"
)

// Duration represents a time.Duration
type Duration struct {
	time.Duration
}

// MarshalJSON marshals a time.Duration to JSON
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Seconds())
}

// UnmarshalJSON unmarshals a JSON value to time.Duration. The JSON value
// can be either a float which is interpreted as seconds or a string
// that is interpreted as a formatted duration.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type ffmpegAVstreamIO struct {
	State  string `json:"state"`
	Packet uint64 `json:"packet"` // counter
	Time   uint64 `json:"time"`
	Size   uint64 `json:"size_kb"` // kbytes
}

func (avio *ffmpegAVstreamIO) export() app.AVstreamIO {
	return app.AVstreamIO{
		State:  avio.State,
		Packet: avio.Packet,
		Time:   avio.Time,
		Size:   avio.Size * 1024,
	}
}

type ffmpegAVstream struct {
	Input          ffmpegAVstreamIO `json:"input"`
	Output         ffmpegAVstreamIO `json:"output"`
	Address        string           `json:"id"`
	URL            string           `json:"url"`
	Stream         uint64           `json:"stream"`
	Aqueue         uint64           `json:"aqueue"`
	Queue          uint64           `json:"queue"`
	Dup            uint64           `json:"dup"`
	Drop           uint64           `json:"drop"`
	Enc            uint64           `json:"enc"`
	Looping        bool             `json:"looping"`
	LoopingRuntime uint64           `json:"looping_runtime"`
	Duplicating    bool             `json:"duplicating"`
	GOP            string           `json:"gop"`
}

func (av *ffmpegAVstream) export() *app.AVstream {
	return &app.AVstream{
		Aqueue:         av.Aqueue,
		Queue:          av.Queue,
		Drop:           av.Drop,
		Dup:            av.Dup,
		Enc:            av.Enc,
		Looping:        av.Looping,
		LoopingRuntime: av.LoopingRuntime,
		Duplicating:    av.Duplicating,
		GOP:            av.GOP,
		Input:          av.Input.export(),
		Output:         av.Output.export(),
	}
}

type ffmpegProgressIO struct {
	// common
	Index     uint64  `json:"index"`
	Stream    uint64  `json:"stream"`
	SizeKB    uint64  `json:"size_kb"`    // kbytes
	Size      uint64  `json:"size_bytes"` // bytes
	Bitrate   float64 `json:"-"`          // bit/s
	Frame     uint64  `json:"frame"`      // counter
	Keyframe  uint64  `json:"keyframe"`   // counter
	Framerate struct {
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Average float64 `json:"avg"`
	} `json:"framerate"`
	Packet    uint64  `json:"packet"`               // counter
	Extradata uint64  `json:"extradata_size_bytes"` // bytes
	FPS       float64 `json:"-"`                    // rate, frames per second
	PPS       float64 `json:"-"`                    // rate, packets per second

	// video
	Quantizer float64 `json:"q"`
}

func (io *ffmpegProgressIO) exportTo(progress *app.ProgressIO) {
	progress.Index = io.Index
	progress.Stream = io.Stream
	progress.Frame = io.Frame
	progress.Keyframe = io.Keyframe
	progress.Framerate.Min = io.Framerate.Min
	progress.Framerate.Max = io.Framerate.Max
	progress.Framerate.Average = io.Framerate.Average
	progress.Packet = io.Packet
	progress.FPS = io.FPS
	progress.PPS = io.PPS
	progress.Quantizer = io.Quantizer
	progress.Bitrate = io.Bitrate
	progress.Extradata = io.Extradata

	if io.Size == 0 {
		progress.Size = io.SizeKB * 1024
	} else {
		progress.Size = io.Size
	}
}

type ffmpegProgress struct {
	Input     []ffmpegProgressIO `json:"inputs"`
	Output    []ffmpegProgressIO `json:"outputs"`
	Frame     uint64             `json:"frame"`  // counter
	Packet    uint64             `json:"packet"` // counter
	FPS       float64            `json:"-"`      // rate, frames per second
	PPS       float64            `json:"-"`      // rate, packets per second
	Quantizer float64            `json:"q"`
	SizeKB    uint64             `json:"size_kb"`    // kbytes
	Size      uint64             `json:"size_bytes"` // bytes
	Bitrate   float64            `json:"-"`          // bit/s
	Time      Duration           `json:"time"`
	Speed     float64            `json:"speed"`
	Drop      uint64             `json:"drop"` // counter
	Dup       uint64             `json:"dup"`  // counter
}

func (p *ffmpegProgress) exportTo(progress *app.Progress) {
	progress.Frame = p.Frame
	progress.Packet = p.Packet
	progress.FPS = p.FPS
	progress.PPS = p.PPS
	progress.Quantizer = p.Quantizer
	progress.Time = p.Time.Seconds()
	progress.Bitrate = p.Bitrate
	progress.Speed = p.Speed
	progress.Drop = p.Drop
	progress.Dup = p.Dup

	if p.Size == 0 {
		progress.Size = p.SizeKB * 1024
	} else {
		progress.Size = p.Size
	}

	for i := range p.Input {
		if len(progress.Input) <= i {
			break
		}

		p.Input[i].exportTo(&progress.Input[i])
	}

	for i := range p.Output {
		if len(progress.Output) <= i {
			break
		}

		p.Output[i].exportTo(&progress.Output[i])
	}
}

type ffmpegProcessIO struct {
	// common
	Address string `json:"url"`
	IP      string `json:"-"`
	Format  string `json:"format"`
	Index   uint64 `json:"index"`
	Stream  uint64 `json:"stream"`
	Type    string `json:"type"`
	Codec   string `json:"codec"`
	Coder   string `json:"coder"`

	// video
	Pixfmt string `json:"pix_fmt"`
	Width  uint64 `json:"width"`
	Height uint64 `json:"height"`

	// audio
	Sampling uint64 `json:"sampling_hz"`
	Layout   string `json:"layout"`
	Channels uint64 `json:"channels"`
}

func (io *ffmpegProcessIO) export() app.ProgressIO {
	return app.ProgressIO{
		Address:  io.Address,
		Format:   io.Format,
		Index:    io.Index,
		Stream:   io.Stream,
		Type:     io.Type,
		Codec:    io.Codec,
		Coder:    io.Coder,
		Pixfmt:   io.Pixfmt,
		Width:    io.Width,
		Height:   io.Height,
		Sampling: io.Sampling,
		Layout:   io.Layout,
		Channels: io.Channels,
	}
}

type ffmpegProcess struct {
	input  []ffmpegProcessIO
	output []ffmpegProcessIO
}

func (p *ffmpegProcess) export() app.Progress {
	progress := app.Progress{}

	for _, io := range p.input {
		aio := io.export()

		progress.Input = append(progress.Input, aio)
	}

	for _, io := range p.output {
		aio := io.export()

		progress.Output = append(progress.Output, aio)
	}

	return progress
}
