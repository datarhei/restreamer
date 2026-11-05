package models

import (
	"time"

	"github.com/datarhei/core/v16/http/graph/scalars"
	"github.com/datarhei/core/v16/playout"
	"github.com/datarhei/core/v16/restream/app"
)

func (s *RawAVstream) UnmarshalPlayout(status playout.Status) {
	s.ID = status.ID
	s.URL = status.Address
	s.Stream = scalars.Uint64(status.Stream)
	s.Queue = scalars.Uint64(status.Queue)
	s.Aqueue = scalars.Uint64(status.AQueue)
	s.Dup = scalars.Uint64(status.Dup)
	s.Drop = scalars.Uint64(status.Drop)
	s.Enc = scalars.Uint64(status.Enc)
	s.Looping = status.Looping
	s.Duplicating = status.Duplicating
	s.Gop = status.GOP
	s.Debug = status.Debug
	s.Input = &RawAVstreamIo{}
	s.Output = &RawAVstreamIo{}
	s.Swap = &RawAVstreamSwap{}

	s.Input.UnmarshalPlayout(status.Input)
	s.Output.UnmarshalPlayout(status.Output)
	s.Swap.UnmarshalPlayout(status)
}

func (i *RawAVstreamIo) UnmarshalPlayout(io playout.StatusIO) {
	i.State = State(io.State)
	i.Packet = scalars.Uint64(io.Packet)
	i.Time = scalars.Uint64(io.Time)
	i.SizeKb = scalars.Uint64(io.Size)
}

func (s *RawAVstreamSwap) UnmarshalPlayout(status playout.Status) {
	s.URL = status.Swap.Address
	s.Status = status.Swap.Status
	s.Lasturl = status.Swap.LastAddress
	s.Lasterror = status.Swap.LastError
}

func (p *Process) UnmarshalRestream(process *app.Process, state *app.State, report *app.Log, metadata map[string]interface{}) {
	p.ID = process.ID
	p.Type = "ffmpeg"
	p.Reference = process.Reference
	p.CreatedAt = time.Unix(process.CreatedAt, 0)
	p.Config = &ProcessConfig{}
	p.State = &ProcessState{}
	p.Report = &ProcessReport{}
	p.Metadata = metadata

	p.Config.UnmarshalRestream(process.Config)
	p.State.UnmarshalRestream(state)
	p.Report.UnmarshalRestream(report)
}

func (c *ProcessConfig) UnmarshalRestream(config *app.Config) {
	c.ID = config.ID
	c.Type = "ffmpeg"
	c.Reference = config.Reference
	c.Input = []*ProcessConfigIo{}
	c.Output = []*ProcessConfigIo{}
	c.Options = config.Options
	c.Reconnect = config.Reconnect
	c.ReconnectDelaySeconds = scalars.Uint64(config.ReconnectDelay)
	c.Autostart = config.Autostart
	c.StaleTimeoutSeconds = scalars.Uint64(config.StaleTimeout)

	c.Limits = &ProcessConfigLimits{
		CPUUsage:       config.LimitCPU,
		MemoryBytes:    scalars.Uint64(config.LimitMemory),
		WaitforSeconds: scalars.Uint64(config.LimitWaitFor),
	}

	for _, io := range config.Input {
		c.Input = append(c.Input, &ProcessConfigIo{
			ID:      io.ID,
			Address: io.Address,
			Options: io.Options,
		})
	}

	for _, io := range config.Output {
		c.Output = append(c.Output, &ProcessConfigIo{
			ID:      io.ID,
			Address: io.Address,
			Options: io.Options,
		})
	}
}

func (s *ProcessState) UnmarshalRestream(state *app.State) {
	s.Order = state.Order
	s.State = state.State
	s.RuntimeSeconds = scalars.Uint64(state.Duration)
	s.ReconnectSeconds = int(state.Reconnect)
	s.LastLogline = state.LastLog
	s.Progress = &Progress{}
	s.MemoryBytes = scalars.Uint64(state.Memory)
	s.CPUUsage = state.CPU
	s.Command = state.Command

	s.Progress.UnmarshalRestream(&state.Progress)
}

func (p *Progress) UnmarshalRestream(progress *app.Progress) {
	p.Input = []*ProgressIo{}
	p.Output = []*ProgressIo{}
	p.Frame = scalars.Uint64(progress.Frame)
	p.Packet = scalars.Uint64(progress.Packet)
	p.Fps = progress.FPS
	p.Q = progress.Quantizer
	p.SizeKb = scalars.Uint64(progress.Size)
	p.Time = progress.Time
	p.BitrateKbit = progress.Bitrate / 1024
	p.Speed = progress.Speed
	p.Drop = scalars.Uint64(progress.Drop)
	p.Dup = scalars.Uint64(progress.Dup)

	for _, io := range progress.Input {
		input := &ProgressIo{}
		input.UnmarshalRestream(&io)

		p.Input = append(p.Input, input)
	}

	for _, io := range progress.Output {
		output := &ProgressIo{}
		output.UnmarshalRestream(&io)

		p.Output = append(p.Output, output)
	}
}

func (p *ProgressIo) UnmarshalRestream(io *app.ProgressIO) {
	p.ID = io.ID
	p.Address = io.Address
	p.Index = scalars.Uint64(io.Index)
	p.Stream = scalars.Uint64(io.Stream)
	p.Format = io.Format
	p.Type = io.Type
	p.Codec = io.Codec
	p.Coder = io.Coder
	p.Frame = scalars.Uint64(io.Frame)
	p.Fps = io.FPS
	p.Packet = scalars.Uint64(io.Packet)
	p.Pps = io.PPS
	p.SizeKb = scalars.Uint64(io.Size)
	p.BitrateKbit = io.Bitrate / 1024
	p.Pixfmt = io.Pixfmt
	p.Q = io.Quantizer
	p.Width = scalars.Uint64(io.Width)
	p.Height = scalars.Uint64(io.Height)
	p.Sampling = scalars.Uint64(io.Sampling)
	p.Layout = io.Layout
	p.Channels = scalars.Uint64(io.Channels)
	p.Avstream = &AVStream{}

	if io.AVstream != nil {
		p.Avstream.UnmarshalRestream(io.AVstream)
	}
}

func (a *AVStream) UnmarshalRestream(avstream *app.AVstream) {
	a.Input = &AVStreamIo{}
	a.Output = &AVStreamIo{}
	a.Aqueue = scalars.Uint64(avstream.Aqueue)
	a.Queue = scalars.Uint64(avstream.Queue)
	a.Dup = scalars.Uint64(avstream.Dup)
	a.Drop = scalars.Uint64(avstream.Drop)
	a.Enc = scalars.Uint64(avstream.Enc)
	a.Looping = avstream.Looping
	a.Duplicating = avstream.Duplicating
	a.Gop = avstream.GOP

	a.Input.UnmarshalRestream(avstream.Input)
	a.Output.UnmarshalRestream(avstream.Output)
}

func (a *AVStreamIo) UnmarshalRestream(io app.AVstreamIO) {
	a.State = io.State
	a.Packet = scalars.Uint64(io.Packet)
	a.Time = scalars.Uint64(io.Time)
	a.SizeKb = scalars.Uint64(io.Size)
}

func (r *ProcessReport) UnmarshalRestream(report *app.Log) {
	r.CreatedAt = report.CreatedAt
	r.Prelude = report.Prelude
	r.Log = []*ProcessReportLogEntry{}
	r.History = []*ProcessReportHistoryEntry{}

	for _, l := range report.Log {
		r.Log = append(r.Log, &ProcessReportLogEntry{
			Timestamp: l.Timestamp,
			Data:      l.Data,
		})
	}

	for _, h := range report.History {
		entry := &ProcessReportHistoryEntry{}
		entry.UnmarshalRestream(h)

		r.History = append(r.History, entry)
	}
}

func (h *ProcessReportHistoryEntry) UnmarshalRestream(entry app.LogHistoryEntry) {
	h.CreatedAt = entry.CreatedAt
	h.Prelude = entry.Prelude
	h.Log = []*ProcessReportLogEntry{}

	for _, l := range entry.Log {
		h.Log = append(h.Log, &ProcessReportLogEntry{
			Timestamp: l.Timestamp,
			Data:      l.Data,
		})
	}
}

func (p *Probe) UnmarshalRestream(probe app.Probe) {
	p.Streams = []*ProbeIo{}
	p.Log = probe.Log

	for _, io := range probe.Streams {
		i := &ProbeIo{}
		i.UnmarshalRestream(io)

		p.Streams = append(p.Streams, i)
	}
}

func (i *ProbeIo) UnmarshalRestream(io app.ProbeIO) {
	i.URL = io.Address
	i.Index = scalars.Uint64(io.Index)
	i.Stream = scalars.Uint64(io.Stream)
	i.Language = io.Language
	i.Type = io.Type
	i.Codec = io.Codec
	i.Coder = io.Coder
	i.BitrateKbps = io.Bitrate
	i.DurationSeconds = io.Duration
	i.Fps = io.FPS
	i.PixFmt = io.Pixfmt
	i.Width = scalars.Uint64(io.Width)
	i.Height = scalars.Uint64(io.Height)
	i.Sampling = scalars.Uint64(io.Sampling)
	i.Layout = io.Layout
	i.Channels = scalars.Uint64(io.Channels)
}
