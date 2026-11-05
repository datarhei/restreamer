package log

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoglevelNames(t *testing.T) {
	require.Equal(t, "DEBUG", Ldebug.String())
	require.Equal(t, "ERROR", Lerror.String())
	require.Equal(t, "WARN", Lwarn.String())
	require.Equal(t, "INFO", Linfo.String())
	require.Equal(t, `SILENT`, Lsilent.String())
}

func TestLogColorToNotTTY(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	w := NewConsoleWriter(writer, Linfo, true).(*syncWriter)
	formatter := w.writer.(*consoleWriter).formatter.(*consoleFormatter)

	require.NotEqual(t, true, formatter.color, "Color should not be used on a buffer logger")
}

func TestLogContext(t *testing.T) {

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("component").WithOutput(NewConsoleWriter(writer, Ldebug, false))

	logger.Debug().Log("debug")
	logger.Info().Log("info")
	logger.Warn().Log("warn")
	logger.Error().Log("error")
	writer.Flush()

	lenWithCtx := buffer.Len()
	buffer.Reset()

	logger = logger.WithComponent("")

	logger.Debug().Log("debug")
	logger.Info().Log("info")
	logger.Warn().Log("warn")
	logger.Error().Log("error")
	writer.Flush()

	lenWithoutCtx := buffer.Len()
	buffer.Reset()

	require.Greater(t, lenWithCtx, lenWithoutCtx, "Log line length without context is not shorter than with context")
}

func TestLogClone(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Linfo, false))

	logger.Info().Log("info")
	writer.Flush()

	require.Contains(t, buffer.String(), `component="test"`)

	buffer.Reset()

	logger2 := logger.WithComponent("tset")

	logger2.Info().Log("info")
	writer.Flush()

	require.Contains(t, buffer.String(), `component="tset"`)
}

func TestLogSilent(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Lsilent, false))

	logger.Debug().Log("debug")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Info().Log("info")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Warn().Log("warn")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Error().Log("error")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()
}

func TestLogDebug(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Ldebug, false))

	logger.Debug().Log("debug")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Info().Log("info")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Warn().Log("warn")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Error().Log("error")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()
}

func TestLogInfo(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Linfo, false))

	logger.Debug().Log("debug")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Info().Log("info")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Warn().Log("warn")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Error().Log("error")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()
}

func TestLogWarn(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Lwarn, false))

	logger.Debug().Log("debug")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Info().Log("info")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Warn().Log("warn")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()

	logger.Error().Log("error")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()
}

func TestLogError(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger := New("test").WithOutput(NewConsoleWriter(writer, Lerror, false))

	logger.Debug().Log("debug")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Info().Log("info")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Warn().Log("warn")
	writer.Flush()
	require.Equal(t, 0, buffer.Len(), "Buffer should be empty")
	buffer.Reset()

	logger.Error().Log("error")
	writer.Flush()
	require.NotEqual(t, 0, buffer.Len(), "Buffer should not be empty")
	buffer.Reset()
}

func TestLogWithField(t *testing.T) {
	bufwriter := NewBufferWriter(Linfo, 10)

	logger := New("test").WithOutput(bufwriter)
	logger = logger.WithField("foo", "bar")
	logger.Info().Log("hello")

	events := bufwriter.Events()

	require.Equal(t, 1, len(events))
	require.Empty(t, events[0].err)
	require.Equal(t, "bar", events[0].Data["foo"])

	logger = logger.WithField("func", func() bool { return true })
	logger.Info().Log("hello")

	events = bufwriter.Events()
	require.Equal(t, 2, len(events))
	require.NotEmpty(t, events[1].err)
	require.Equal(t, "bar", events[0].Data["foo"])
}
