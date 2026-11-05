package api

import (
	"encoding/json"
	"strconv"

	"github.com/datarhei/core/v16/restream/app"
	"github.com/lithammer/shortuuid/v4"
)

// Process represents all information on a process
type Process struct {
	ID        string         `json:"id" jsonschema:"minLength=1"`
	Type      string         `json:"type" jsonschema:"enum=ffmpeg"`
	Reference string         `json:"reference"`
	CreatedAt int64          `json:"created_at" jsonschema:"minimum=0" format:"int64"`
	UpdatedAt int64          `json:"updated_at" jsonschema:"minimum=0" format:"int64"`
	Config    *ProcessConfig `json:"config,omitempty"`
	State     *ProcessState  `json:"state,omitempty"`
	Report    *ProcessReport `json:"report,omitempty"`
	Metadata  Metadata       `json:"metadata,omitempty"`
}

// ProcessConfigIO represents an input or output of an ffmpeg process config
type ProcessConfigIO struct {
	ID      string                   `json:"id"`
	Address string                   `json:"address" validate:"required" jsonschema:"minLength=1"`
	Options []string                 `json:"options"`
	Cleanup []ProcessConfigIOCleanup `json:"cleanup,omitempty"`
}

type ProcessConfigIOCleanup struct {
	Pattern       string `json:"pattern" validate:"required"`
	MaxFiles      uint   `json:"max_files" format:"uint"`
	MaxFileAge    uint   `json:"max_file_age_seconds" format:"uint"`
	PurgeOnDelete bool   `json:"purge_on_delete"`
}

type ProcessConfigLimits struct {
	CPU     float64 `json:"cpu_usage" jsonschema:"minimum=0,maximum=100"`
	Memory  uint64  `json:"memory_mbytes" jsonschema:"minimum=0" format:"uint64"`
	WaitFor uint64  `json:"waitfor_seconds" jsonschema:"minimum=0" format:"uint64"`
}

// ProcessConfig represents the configuration of an ffmpeg process
type ProcessConfig struct {
	ID             string              `json:"id"`
	Type           string              `json:"type" validate:"oneof='ffmpeg' ''" jsonschema:"enum=ffmpeg,enum="`
	Reference      string              `json:"reference"`
	Input          []ProcessConfigIO   `json:"input" validate:"required"`
	Output         []ProcessConfigIO   `json:"output" validate:"required"`
	Options        []string            `json:"options"`
	Reconnect      bool                `json:"reconnect"`
	ReconnectDelay uint64              `json:"reconnect_delay_seconds" format:"uint64"`
	Autostart      bool                `json:"autostart"`
	StaleTimeout   uint64              `json:"stale_timeout_seconds" format:"uint64"`
	Limits         ProcessConfigLimits `json:"limits"`
}

// Marshal converts a process config in API representation to a restreamer process config
func (cfg *ProcessConfig) Marshal() *app.Config {
	p := &app.Config{
		ID:             cfg.ID,
		Reference:      cfg.Reference,
		Options:        cfg.Options,
		Reconnect:      cfg.Reconnect,
		ReconnectDelay: cfg.ReconnectDelay,
		Autostart:      cfg.Autostart,
		StaleTimeout:   cfg.StaleTimeout,
		LimitCPU:       cfg.Limits.CPU,
		LimitMemory:    cfg.Limits.Memory * 1024 * 1024,
		LimitWaitFor:   cfg.Limits.WaitFor,
	}

	cfg.generateInputOutputIDs(cfg.Input)

	for _, x := range cfg.Input {
		p.Input = append(p.Input, app.ConfigIO{
			ID:      x.ID,
			Address: x.Address,
			Options: x.Options,
		})
	}

	cfg.generateInputOutputIDs(cfg.Output)

	for _, x := range cfg.Output {
		output := app.ConfigIO{
			ID:      x.ID,
			Address: x.Address,
			Options: x.Options,
		}

		for _, c := range x.Cleanup {
			output.Cleanup = append(output.Cleanup, app.ConfigIOCleanup{
				Pattern:       c.Pattern,
				MaxFiles:      c.MaxFiles,
				MaxFileAge:    c.MaxFileAge,
				PurgeOnDelete: c.PurgeOnDelete,
			})
		}

		p.Output = append(p.Output, output)

	}

	return p
}

func (cfg *ProcessConfig) generateInputOutputIDs(ioconfig []ProcessConfigIO) {
	ids := map[string]struct{}{}

	for _, io := range ioconfig {
		if len(io.ID) == 0 {
			continue
		}

		ids[io.ID] = struct{}{}
	}

	for i, io := range ioconfig {
		if len(io.ID) != 0 {
			continue
		}

		for {
			id := shortuuid.New()
			if _, ok := ids[id]; !ok {
				ioconfig[i].ID = id
				break
			}
		}
	}
}

// Unmarshal converts a restream process config to a process config in API representation
func (cfg *ProcessConfig) Unmarshal(c *app.Config) {
	if c == nil {
		return
	}

	cfg.ID = c.ID
	cfg.Reference = c.Reference
	cfg.Type = "ffmpeg"
	cfg.Reconnect = c.Reconnect
	cfg.ReconnectDelay = c.ReconnectDelay
	cfg.Autostart = c.Autostart
	cfg.StaleTimeout = c.StaleTimeout
	cfg.Limits.CPU = c.LimitCPU
	cfg.Limits.Memory = c.LimitMemory / 1024 / 1024
	cfg.Limits.WaitFor = c.LimitWaitFor

	cfg.Options = make([]string, len(c.Options))
	copy(cfg.Options, c.Options)

	for _, x := range c.Input {
		io := ProcessConfigIO{
			ID:      x.ID,
			Address: x.Address,
		}

		io.Options = make([]string, len(x.Options))
		copy(io.Options, x.Options)

		cfg.Input = append(cfg.Input, io)
	}

	for _, x := range c.Output {
		io := ProcessConfigIO{
			ID:      x.ID,
			Address: x.Address,
		}

		io.Options = make([]string, len(x.Options))
		copy(io.Options, x.Options)

		for _, c := range x.Cleanup {
			io.Cleanup = append(io.Cleanup, ProcessConfigIOCleanup{
				Pattern:       c.Pattern,
				MaxFiles:      c.MaxFiles,
				MaxFileAge:    c.MaxFileAge,
				PurgeOnDelete: c.PurgeOnDelete,
			})
		}

		cfg.Output = append(cfg.Output, io)
	}
}

// ProcessReportHistoryEntry represents the logs of a run of a restream process
type ProcessReportHistoryEntry struct {
	CreatedAt int64       `json:"created_at" format:"int64"`
	Prelude   []string    `json:"prelude"`
	Log       [][2]string `json:"log"`
}

// ProcessReport represents the current log and the logs of previous runs of a restream process
type ProcessReport struct {
	ProcessReportHistoryEntry
	History []ProcessReportHistoryEntry `json:"history"`
}

// Unmarshal converts a restream log to a report
func (report *ProcessReport) Unmarshal(l *app.Log) {
	if l == nil {
		return
	}

	report.CreatedAt = l.CreatedAt.Unix()
	report.Prelude = l.Prelude
	report.Log = make([][2]string, len(l.Log))
	for i, line := range l.Log {
		report.Log[i][0] = strconv.FormatInt(line.Timestamp.Unix(), 10)
		report.Log[i][1] = line.Data
	}

	report.History = []ProcessReportHistoryEntry{}

	for _, h := range l.History {
		he := ProcessReportHistoryEntry{
			CreatedAt: h.CreatedAt.Unix(),
			Prelude:   h.Prelude,
			Log:       make([][2]string, len(h.Log)),
		}

		for i, line := range h.Log {
			he.Log[i][0] = strconv.FormatInt(line.Timestamp.Unix(), 10)
			he.Log[i][1] = line.Data
		}

		report.History = append(report.History, he)
	}
}

// ProcessState represents the current state of an ffmpeg process
type ProcessState struct {
	Order     string      `json:"order" jsonschema:"enum=start,enum=stop"`
	State     string      `json:"exec" jsonschema:"enum=finished,enum=starting,enum=running,enum=finishing,enum=killed,enum=failed"`
	Runtime   int64       `json:"runtime_seconds" jsonschema:"minimum=0" format:"int64"`
	Reconnect int64       `json:"reconnect_seconds" format:"int64"`
	LastLog   string      `json:"last_logline"`
	Progress  *Progress   `json:"progress"`
	Memory    uint64      `json:"memory_bytes" format:"uint64"`
	CPU       json.Number `json:"cpu_usage" swaggertype:"number" jsonschema:"type=number"`
	Command   []string    `json:"command"`
}

// Unmarshal converts a restreamer ffmpeg process state to a state in API representation
func (s *ProcessState) Unmarshal(state *app.State) {
	if state == nil {
		return
	}

	s.Order = state.Order
	s.State = state.State
	s.Runtime = int64(state.Duration)
	s.Reconnect = int64(state.Reconnect)
	s.LastLog = state.LastLog
	s.Progress = &Progress{}
	s.Memory = state.Memory
	s.CPU = toNumber(state.CPU)
	s.Command = state.Command

	s.Progress.Unmarshal(&state.Progress)
}
