package log

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type Formatter interface {
	Bytes(e *Event) []byte
	String(e *Event) string
}

type jsonFormatter struct{}

func NewJSONFormatter() Formatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) Bytes(e *Event) []byte {
	e.Data["ts"] = e.Time
	e.Data["component"] = e.Component

	if len(e.Caller) != 0 {
		e.Data["caller"] = e.Caller
	}

	if len(e.Message) != 0 {
		e.Data["message"] = e.Message
	}

	data, _ := json.Marshal(&e)

	return data
}

func (f *jsonFormatter) String(e *Event) string {
	return string(f.Bytes(e))
}

type consoleFormatter struct {
	color bool
}

func NewConsoleFormatter(useColor bool) Formatter {
	return &consoleFormatter{
		color: useColor,
	}
}

func (f *consoleFormatter) Bytes(e *Event) []byte {
	return []byte(f.String(e))
}

func (f *consoleFormatter) String(e *Event) string {
	datetime := e.Time.UTC().Format(time.RFC3339)
	level := e.Level.String()

	if f.color {
		switch e.Level {
		case Ldebug:
			level = fmt.Sprintf("\033[35m%s\033[0m", level)
		case Linfo:
			level = fmt.Sprintf("\033[34m%s\033[0m", level)
		case Lwarn:
			level = fmt.Sprintf("\033[33m%s\033[0m", level)
		case Lerror:
			level = fmt.Sprintf("\033[31m\033[5m%s\033[0m", level)
		default:
		}
	}

	message := fmt.Sprintf("%s %s %s", f.writeKV("ts", datetime), f.writeKV("level", level), f.writeKV("component", f.quote(e.Component)))

	if len(e.Message) != 0 {
		message += fmt.Sprintf(" %s", f.writeKV("msg", f.quote(e.Message)))
	}

	// Sort the map keys
	keys := make([]string, len(e.Data))
	i := 0

	for key := range e.Data {
		keys[i] = key
		i++
	}

	if i > 1 {
		sort.Strings(keys)
	}

	v := ""
	for _, key := range keys {
		value := e.Data[key]

		switch val := value.(type) {
		case string:
			v = f.quote(val)
		case error:
			v = f.quote(val.Error())
		default:
			if jsonvalue, err := json.Marshal(value); err == nil {
				v = string(jsonvalue)
			} else {
				v = f.quote(err.Error())
			}
		}

		message += fmt.Sprintf(" %s", f.writeKV(key, v))
	}

	message += "\n"

	return message
}

func (f *consoleFormatter) writeKV(key string, value string) string {
	if !f.color {
		return fmt.Sprintf("%s=%s", key, value)
	}

	if key == "error" {
		value = "\033[31m" + value + "\033[0m"
	}

	return fmt.Sprintf("\033[90m%s=\033[0m%s", key, value)
}

func (f *consoleFormatter) quote(s string) string {
	return strconv.Quote(s)
}
