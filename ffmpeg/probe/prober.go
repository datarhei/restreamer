package probe

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/datarhei/core/v16/ffmpeg/prelude"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/process"
	"github.com/datarhei/core/v16/restream/app"
)

type Parser interface {
	process.Parser

	Probe() app.Probe
}

type Config struct {
	Logger log.Logger
}

type prober struct {
	data   []process.Line
	inputs []probeIO

	logger log.Logger
}

func New(config Config) Parser {
	p := &prober{
		logger: config.Logger,
	}

	if p.logger == nil {
		p.logger = log.New("Parser")
	}

	return p
}

func (p *prober) Probe() app.Probe {
	probe := app.Probe{}

	for _, io := range p.inputs {
		probe.Streams = append(probe.Streams, io.export())
	}

	probe.Log = make([]string, len(p.data))
	for i, line := range p.data {
		probe.Log[i] = line.Data
	}

	return probe
}

func (p *prober) Parse(line string) uint64 {
	if strings.HasPrefix(line, "avstream.progress:") {
		return 0
	}

	p.data = append(p.data, process.Line{
		Timestamp: time.Now(),
		Data:      line,
	})

	return 0
}

func (p *prober) parseJSON(line string) {
	inputs := []probeIO{}

	err := json.Unmarshal([]byte(line), &inputs)
	if err != nil {
		p.logger.WithFields(log.Fields{
			"line":  line,
			"error": err,
		}).Error().Log("Failed parsing inputs")
	}

	p.inputs = inputs
}

func (p *prober) parseDefault() {
	lines := make([]string, len(p.data))

	for i, line := range p.data {
		lines[i] = line.Data
	}

	inputs, _, _ := prelude.Parse(lines)

	p.inputs = make([]probeIO, len(inputs))

	for i, input := range inputs {
		p.inputs[i].Address = input.Address
		p.inputs[i].Format = input.Format
		p.inputs[i].Duration = input.Duration
		p.inputs[i].Index = input.Index
		p.inputs[i].Stream = input.Stream
		p.inputs[i].Language = input.Language
		p.inputs[i].Type = input.Type
		p.inputs[i].Codec = input.Codec
		p.inputs[i].Bitrate = input.Bitrate
		p.inputs[i].Pixfmt = input.Pixfmt
		p.inputs[i].Width = input.Width
		p.inputs[i].Height = input.Height
		p.inputs[i].FPS = input.FPS
		p.inputs[i].Sampling = input.Sampling
		p.inputs[i].Layout = input.Layout
	}
}

func (p *prober) Log() []process.Line {
	return p.data
}

func (p *prober) ResetStats() {
	hasJSON := false
	prefix := "ffmpeg.inputs:"

	for _, line := range p.data {
		isFFmpegInputs := strings.HasPrefix(line.Data, prefix)

		if isFFmpegInputs {
			p.parseJSON(strings.TrimPrefix(line.Data, prefix))
			hasJSON = true

			break
		}
	}

	if !hasJSON {
		p.parseDefault()
	}
}

func (p *prober) ResetLog() {
	p.data = []process.Line{}
	p.inputs = []probeIO{}
}
