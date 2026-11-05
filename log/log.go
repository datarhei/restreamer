// Package log provides an opiniated logging facility as it provides only 4 log levels.
package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// LogLevel represents a log level
type Level uint

const (
	Lsilent Level = 0
	Lerror  Level = 1
	Lwarn   Level = 2
	Linfo   Level = 3
	Ldebug  Level = 4
)

// String returns a string representing the log level.
func (level Level) String() string {
	names := []string{
		"SILENT",
		"ERROR",
		"WARN",
		"INFO",
		"DEBUG",
	}

	if level > Ldebug {
		return `¯\_(ツ)_/¯`
	}

	return names[level]
}

func (level *Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(level.String())
}

type Fields map[string]interface{}

// Logger is an interface that provides means for writing log messages.
//
// There are 4 log levels available (debug, info, warn, error) with increasing
// severity. A message will be written to an output if the log level of the message
// has the same or a higher severity than the output. Otherwise it will be
// discarded.
//
// Different outputs should be added with the AddOutput function.
//
// The context is a string that represents who wrote the message.
type Logger interface {
	// WithOutput adds an output to the Logger. The messages are written to the
	// provided writer if the log level of the message is more or equally critical
	// than level. Pass true for useColoe if colored output is desired. If the
	// writer doesn't support colored output, it will be automatically disabled.
	// The returned value implements the LoggerOutput interface which allows to
	// change the log level at any later point in time.
	WithOutput(w Writer) Logger

	// With returns a new Logger with the given context. The context may
	// printed along the message. This is up to the implementation.
	WithComponent(component string) Logger

	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger

	WithError(err error) Logger

	Log(format string, args ...interface{})

	// Debug writes a message with the debug log level to all registered outputs.
	// The message will be written according to fmt.Printf(). The detail field will
	// be reset to nil.
	Debug() Logger

	// Info writes a message with the info log level to all registered outputs.
	// The message will be written according to fmt.Printf(). The detail field will
	// be reset to nil.
	Info() Logger

	// Warn writes a message with the warn log level to all registered outputs.
	// The message will be written according to fmt.Printf(). The detail field will
	// be reset to nil.
	Warn() Logger

	// Error writes a message with the error log level to all registered outputs.
	// The message will be written according to fmt.Printf(). The detail field will
	// be reset to nil.
	Error() Logger

	// WithLevel writes a message with the given level to all registered outputs.
	// The message will be written according to fmt.Printf(). The detail field will
	// be reset to nil.
	WithLevel(level Level) Logger

	// Write implements the io.Writer interface such that it can be used in e.g. the
	// the log/Logger facility. Messages will be printed with debug level.
	Write(p []byte) (int, error)
}

// logger is an implementation of the Logger interface.
type logger struct {
	output     Writer
	component  string
	modulePath string
}

// New returns an implementation of the Logger interface.
func New(component string) Logger {
	l := &logger{
		component: component,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		l.modulePath = info.Path
	}

	return l
}

func (l *logger) clone() *logger {
	clone := &logger{
		output:     l.output,
		component:  l.component,
		modulePath: l.modulePath,
	}

	return clone
}

func (l *logger) WithOutput(w Writer) Logger {
	clone := l.clone()
	clone.output = w

	return clone
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return newEvent(l).WithField(key, value)
}

func (l *logger) WithFields(f Fields) Logger {
	return newEvent(l).WithFields(f)
}

func (l *logger) WithError(err error) Logger {
	return newEvent(l).WithError(err)
}

func (l *logger) WithComponent(component string) Logger {
	clone := l.clone()
	clone.component = component

	return clone
}

func (l *logger) Log(format string, args ...interface{}) {
	e := newEvent(l)

	e.Log(format, args...)
}

func (l *logger) Debug() Logger {
	return newEvent(l).Debug()
}

func (l *logger) Info() Logger {
	return newEvent(l).Info()
}

func (l *logger) Warn() Logger {
	return newEvent(l).Warn()
}

func (l *logger) Error() Logger {
	return newEvent(l).Error()
}

func (l *logger) WithLevel(level Level) Logger {
	return newEvent(l).WithLevel(level)
}

func (l *logger) Write(p []byte) (int, error) {
	return newEvent(l).Write(p)
}

type Event struct {
	logger *logger

	Time      time.Time
	Level     Level
	Component string
	Caller    string
	Message   string

	err string

	Data Fields
}

func newEvent(l *logger) Logger {
	e := &Event{
		logger:    l,
		Component: l.component,
		Data:      map[string]interface{}{},
	}

	return e
}

func (e *Event) WithOutput(w Writer) Logger {
	return e.logger.WithOutput(w)
}

func (e *Event) WithComponent(component string) Logger {
	clone := e.clone()

	clone.Component = component

	return clone
}

func (e *Event) Log(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = strings.TrimPrefix(file, e.logger.modulePath)

	n := e.clone()

	n.logger = nil
	n.Time = time.Now()
	n.Caller = fmt.Sprintf("%s:%d", file, line)

	if n.Level == Lsilent {
		n.Level = Ldebug
	}

	if len(format) != 0 {
		if len(args) == 0 {
			n.Message = format
		} else {
			n.Message = fmt.Sprintf(format, args...)
		}
	}

	if e.logger.output != nil {
		e.logger.output.Write(n)
	}
}

func (e *Event) clone() *Event {
	data := make(Fields, len(e.Data))
	for k, v := range e.Data {
		data[k] = v
	}

	return &Event{
		Time:      e.Time,
		Caller:    e.Caller,
		logger:    e.logger,
		Level:     e.Level,
		Component: e.Component,
		Message:   e.Message,
		err:       e.err,
		Data:      data,
	}
}

func (e *Event) WithField(key string, value interface{}) Logger {
	return e.WithFields(Fields{
		key: value,
	})
}

const maxFields = 1024

func (e *Event) WithFields(f Fields) Logger {
	if maxFields-len(e.Data)-len(f) < 0 {
		return e
	}

	data := make(Fields, len(e.Data)+len(f))
	for k, v := range e.Data {
		data[k] = v
	}

	fieldErr := e.err
	for k, v := range f {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = e.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}

	return &Event{
		logger:    e.logger,
		Component: e.Component,
		Level:     e.Level,
		err:       fieldErr,
		Data:      data,
	}
}

func (e *Event) WithError(err error) Logger {
	return e.WithFields(Fields{
		"error": err,
	})
}

func (e *Event) Debug() Logger {
	return e.WithLevel(Ldebug)
}

func (e *Event) Info() Logger {
	return e.WithLevel(Linfo)
}

func (e *Event) Warn() Logger {
	return e.WithLevel(Lwarn)
}

func (e *Event) Error() Logger {
	return e.WithLevel(Lerror)
}

func (e *Event) WithLevel(level Level) Logger {
	clone := e.clone()
	clone.Level = level

	return clone
}

func (l *Event) Write(p []byte) (int, error) {
	l.Log("%s", strings.TrimSpace(string(p)))

	return len(p), nil
}

type Eventx struct {
	Time      time.Time   `json:"ts"`
	Level     Level       `json:"level"`
	Component string      `json:"component"`
	Reference string      `json:"ref"`
	Message   string      `json:"message"`
	Caller    string      `json:"caller"`
	Detail    interface{} `json:"detail"`
}
