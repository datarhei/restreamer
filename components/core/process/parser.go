package process

import (
	"time"
)

// Parser is an interface that is accepted by a process for parsing
// the process' output.
type Parser interface {
	// Parse parses the given line and returns an indicator
	// for progress (e.g. based on the contents of the line,
	// or previous line, ...)
	Parse(line string) uint64

	// Reset resets any collected statistics or temporary data.
	// This is called before the process starts and after the
	// process stopped. The stats are meant to be collected
	// during the runtime of the process.
	ResetStats()

	// ResetLog resets any collected logs. This is called
	// before the process starts.
	ResetLog()

	// Log returns a slice of collected log lines
	Log() []Line
}

// Line represents a line from the output with its timestamp. The
// line doesn't include any newline character.
type Line struct {
	Timestamp time.Time
	Data      string
}

type nullParser struct{}

// NewNullParser returns a dummy parser that isn't doing anything
// and always returns progress.
func NewNullParser() Parser {
	return &nullParser{}
}

var _ Parser = &nullParser{}

func (p *nullParser) Parse(line string) uint64 { return 1 }

func (p *nullParser) Log() []Line { return []Line{} }

func (p *nullParser) ResetStats() {}

func (p *nullParser) ResetLog() {}
