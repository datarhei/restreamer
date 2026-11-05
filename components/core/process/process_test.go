package process

import (
	"testing"
	"time"

	"github.com/datarhei/core/v16/internal/testhelper"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sleep",
		Args: []string{
			"10",
		},
		Reconnect:    false,
		StaleTimeout: 0,
	})

	require.Equal(t, "finished", p.Status().State)

	p.Start()

	require.Equal(t, "running", p.Status().State)

	time.Sleep(5 * time.Second)

	require.Equal(t, "running", p.Status().State)

	p.Stop(false)

	time.Sleep(2 * time.Second)

	require.Equal(t, "killed", p.Status().State)
}

func TestReconnectProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sleep",
		Args: []string{
			"2",
		},
		Reconnect:      true,
		ReconnectDelay: 2 * time.Second,
		StaleTimeout:   0,
	})

	p.Start()

	time.Sleep(3 * time.Second)

	require.Equal(t, "finished", p.Status().State)

	p.Stop(false)

	require.Equal(t, "finished", p.Status().State)
}

func TestStaleProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sleep",
		Args: []string{
			"10",
		},
		Reconnect:    false,
		StaleTimeout: 2 * time.Second,
	})

	p.Start()

	time.Sleep(5 * time.Second)

	require.Equal(t, "killed", p.Status().State)

	p.Stop(false)

	require.Equal(t, "killed", p.Status().State)
}

func TestStaleReconnectProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sleep",
		Args: []string{
			"10",
		},
		Reconnect:    false,
		StaleTimeout: 2 * time.Second,
	})

	p.Start()

	time.Sleep(10 * time.Second)

	require.Equal(t, "killed", p.Status().State)

	p.Stop(false)

	require.Equal(t, "killed", p.Status().State)
}

func TestNonExistingProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sloop",
		Args: []string{
			"10",
		},
		Reconnect:      false,
		ReconnectDelay: 5 * time.Second,
		StaleTimeout:   0,
	})

	p.Start()

	time.Sleep(3 * time.Second)

	require.Equal(t, "failed", p.Status().State)

	p.Stop(false)

	require.Equal(t, "failed", p.Status().State)
}

func TestNonExistingReconnectProcess(t *testing.T) {
	p, _ := New(Config{
		Binary: "sloop",
		Args: []string{
			"10",
		},
		Reconnect:      true,
		ReconnectDelay: 2 * time.Second,
		StaleTimeout:   0,
	})

	p.Start()

	time.Sleep(5 * time.Second)

	require.Equal(t, "failed", p.Status().State)

	p.Stop(false)

	require.Equal(t, "failed", p.Status().State)
}

func TestProcessFailed(t *testing.T) {
	p, _ := New(Config{
		Binary: "sleep",
		Args: []string{
			"hello",
		},
		Reconnect:    false,
		StaleTimeout: 0,
	})

	p.Start()

	time.Sleep(5 * time.Second)

	p.Stop(false)

	require.Equal(t, "failed", p.Status().State)
}

func TestFFmpegWaitStop(t *testing.T) {
	binary, err := testhelper.BuildBinary("sigintwait", "../internal/testhelper")
	require.NoError(t, err, "Failed to build helper program")

	p, _ := New(Config{
		Binary:       binary,
		Args:         []string{},
		Reconnect:    false,
		StaleTimeout: 0,
		OnExit: func() {
			time.Sleep(2 * time.Second)
		},
	})

	err = p.Start()
	require.NoError(t, err)

	time.Sleep(4 * time.Second)

	p.Stop(true)

	require.Equal(t, "finished", p.Status().State)
}

func TestFFmpegKill(t *testing.T) {
	binary, err := testhelper.BuildBinary("sigint", "../internal/testhelper")
	require.NoError(t, err, "Failed to build helper program")

	p, _ := New(Config{
		Binary:       binary,
		Args:         []string{},
		Reconnect:    false,
		StaleTimeout: 0,
	})

	err = p.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)

	p.Stop(false)

	time.Sleep(3 * time.Second)

	require.Equal(t, "finished", p.Status().State)
}

func TestProcessForceKill(t *testing.T) {
	binary, err := testhelper.BuildBinary("ignoresigint", "../internal/testhelper")
	require.NoError(t, err, "Failed to build helper program")

	p, _ := New(Config{
		Binary:       binary,
		Args:         []string{},
		Reconnect:    false,
		StaleTimeout: 0,
	})

	err = p.Start()
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	p.Stop(false)

	time.Sleep(1 * time.Second)

	require.Equal(t, "finishing", p.Status().State)

	time.Sleep(5 * time.Second)

	require.Equal(t, "killed", p.Status().State)
}
