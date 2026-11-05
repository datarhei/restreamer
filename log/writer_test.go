package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJSONWriter(t *testing.T) {
	buffer := bytes.Buffer{}

	writer := NewJSONWriter(&buffer, Linfo)
	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "test",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	require.Equal(t, `{"Time":"2009-11-10T23:00:00Z","Level":"INFO","Component":"test","Caller":"me","Message":"hello world","Data":{"caller":"me","component":"test","foo":"bar","message":"hello world","ts":"2009-11-10T23:00:00Z"}}`, buffer.String())
}

func TestConsoleWriter(t *testing.T) {
	buffer := bytes.Buffer{}

	writer := NewConsoleWriter(&buffer, Linfo, false)
	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "test",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	require.Equal(t, `ts=2009-11-10T23:00:00Z level=INFO component="test" msg="hello world" foo="bar"`+"\n", buffer.String())
}

func TestTopicWriter(t *testing.T) {
	bufwriter := NewBufferWriter(Linfo, 10)
	writer1 := NewTopicWriter(bufwriter, []string{})
	writer2 := NewTopicWriter(bufwriter, []string{"foobar"})

	writer1.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "test",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	writer2.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "test",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	require.Equal(t, 1, len(bufwriter.Events()))

	writer1.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "foobar",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	writer2.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "foobar",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	require.Equal(t, 3, len(bufwriter.Events()))
}

func TestMultiwriter(t *testing.T) {
	bufwriter1 := NewBufferWriter(Linfo, 10)
	bufwriter2 := NewBufferWriter(Linfo, 10)

	writer := NewMultiWriter(bufwriter1, bufwriter2)

	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "foobar",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	require.Equal(t, 1, len(bufwriter1.Events()))
	require.Equal(t, 1, len(bufwriter2.Events()))
}

func TestLevelRewriter(t *testing.T) {
	bufwriter := NewBufferWriter(Linfo, 10)

	rule := LevelRewriteRule{
		Level:     Lwarn,
		Component: "foobar",
		Match: map[string]string{
			"foo": "bar",
		},
	}

	writer := NewLevelRewriter(bufwriter, []LevelRewriteRule{rule})
	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "foobar",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	events := bufwriter.Events()

	require.Equal(t, 1, len(events))
	require.Equal(t, Lwarn, events[0].Level)

	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "foobar",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"bar": "foo"},
	})

	events = bufwriter.Events()

	require.Equal(t, 2, len(events))
	require.Equal(t, Linfo, events[1].Level)

	writer.Write(&Event{
		logger:    &logger{},
		Time:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Level:     Linfo,
		Component: "test",
		Caller:    "me",
		Message:   "hello world",
		err:       "",
		Data:      map[string]interface{}{"foo": "bar"},
	})

	events = bufwriter.Events()

	require.Equal(t, 3, len(events))
	require.Equal(t, Linfo, events[2].Level)
}
