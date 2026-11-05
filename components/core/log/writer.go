package log

import (
	"container/ring"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/mattn/go-isatty"
)

type Writer interface {
	Write(e *Event) error
}

type jsonWriter struct {
	writer    io.Writer
	level     Level
	formatter Formatter
}

func NewJSONWriter(w io.Writer, level Level) Writer {
	writer := &jsonWriter{
		writer:    w,
		level:     level,
		formatter: NewJSONFormatter(),
	}

	return NewSyncWriter(writer)
}

func (w *jsonWriter) Write(e *Event) error {
	if w.level < e.Level || e.Level == Lsilent {
		return nil
	}

	_, err := w.writer.Write(w.formatter.Bytes(e))

	return err
}

type consoleWriter struct {
	writer    io.Writer
	level     Level
	formatter Formatter
}

func NewConsoleWriter(w io.Writer, level Level, useColor bool) Writer {
	writer := &consoleWriter{
		writer: w,
		level:  level,
	}

	color := useColor

	if color {
		if w, ok := w.(*os.File); ok {
			if !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()) {
				color = false
			}
		} else {
			color = false
		}
	}

	writer.formatter = NewConsoleFormatter(color)

	return NewSyncWriter(writer)
}

func (w *consoleWriter) Write(e *Event) error {
	if w.level < e.Level || e.Level == Lsilent {
		return nil
	}

	_, err := w.writer.Write(w.formatter.Bytes(e))

	return err
}

type topicWriter struct {
	writer Writer
	topics map[string]struct{}
}

func NewTopicWriter(writer Writer, topics []string) Writer {
	w := &topicWriter{
		writer: writer,
		topics: make(map[string]struct{}),
	}

	for _, topic := range topics {
		topic = strings.ToLower(topic)
		w.topics[topic] = struct{}{}
	}

	return w
}

func (w *topicWriter) Write(e *Event) error {
	if len(w.topics) > 0 {
		topic := strings.ToLower(e.Component)
		if _, ok := w.topics[topic]; !ok {
			return nil
		}
	}

	err := w.writer.Write(e)

	return err
}

type levelRewriter struct {
	writer Writer
	rules  []levelRewriteRule
}

type LevelRewriteRule struct {
	Level     Level
	Component string
	Match     map[string]string
}

type levelRewriteRule struct {
	level     Level
	component string
	match     map[string]*regexp.Regexp
}

func NewLevelRewriter(writer Writer, rules []LevelRewriteRule) Writer {
	w := &levelRewriter{
		writer: writer,
	}

	for _, rule := range rules {
		r := levelRewriteRule{
			level:     rule.Level,
			component: rule.Component,
			match:     make(map[string]*regexp.Regexp),
		}

		for k, v := range rule.Match {
			re, err := regexp.Compile(v)
			if err != nil {
				continue
			}
			r.match[k] = re
		}

		w.rules = append(w.rules, r)
	}

	return w
}

func (w *levelRewriter) Write(e *Event) error {
rules:
	for _, r := range w.rules {
		if e.Component != r.component {
			continue
		}

		for k, re := range r.match {
			value, ok := e.Data[k]
			if !ok {
				continue rules
			}

			switch value := value.(type) {
			case string:
				if !re.MatchString(value) {
					continue rules
				}
			}
		}

		e.Level = r.level
	}

	return w.writer.Write(e)
}

type syncWriter struct {
	mu     sync.Mutex
	writer Writer
}

func NewSyncWriter(writer Writer) Writer {
	return &syncWriter{
		writer: writer,
	}
}

func (s *syncWriter) Write(e *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.writer.Write(e)
}

type multiWriter struct {
	writer []Writer
}

func NewMultiWriter(writer ...Writer) Writer {
	mw := &multiWriter{}

	mw.writer = append(mw.writer, writer...)

	return mw
}

func (m *multiWriter) Write(e *Event) error {
	for _, w := range m.writer {
		if err := w.Write(e); err != nil {
			return err
		}
	}

	return nil
}

type BufferWriter interface {
	Writer
	Events() []*Event
}

type bufferWriter struct {
	lines *ring.Ring
	lock  sync.RWMutex
	level Level
}

func NewBufferWriter(level Level, lines int) BufferWriter {
	b := &bufferWriter{
		level: level,
	}

	if lines > 0 {
		b.lines = ring.New(lines)
	}

	return b
}

func (b *bufferWriter) Write(e *Event) error {
	if b.level < e.Level || e.Level == Lsilent {
		return nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.lines != nil {
		b.lines.Value = e.clone()
		b.lines = b.lines.Next()
	}

	return nil
}

func (b *bufferWriter) Events() []*Event {
	var lines = []*Event{}

	if b.lines == nil {
		return lines
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	b.lines.Do(func(l interface{}) {
		if l == nil {
			return
		}

		lines = append(lines, l.(*Event).clone())
	})

	return lines
}
