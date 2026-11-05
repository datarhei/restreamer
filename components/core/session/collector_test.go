package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegisterSession(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Hour,
		SessionTimeout:  time.Hour,
	})
	require.Equal(t, nil, err)

	b := c.IsKnownSession("foobar")
	require.Equal(t, false, b)

	c.Register("foobar", "", "", "")

	b = c.IsKnownSession("foobar")
	require.Equal(t, true, b)

	c.Unregister("foobar")

	time.Sleep(2 * time.Second)

	b = c.IsKnownSession("foobar")
	require.Equal(t, false, b)
}

func TestInactiveSession(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Second,
		SessionTimeout:  time.Hour,
	})
	require.Equal(t, nil, err)

	b := c.IsKnownSession("foobar")
	require.Equal(t, false, b)

	c.Register("foobar", "", "", "")

	b = c.IsKnownSession("foobar")
	require.Equal(t, true, b)

	time.Sleep(3 * time.Second)

	b = c.IsKnownSession("foobar")
	require.Equal(t, false, b)
}

func TestActivateSession(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Second,
		SessionTimeout:  time.Second,
	})
	require.Equal(t, nil, err)

	b := c.IsKnownSession("foobar")
	require.Equal(t, false, b)

	c.RegisterAndActivate("foobar", "", "", "")

	b = c.IsKnownSession("foobar")
	require.Equal(t, true, b)

	time.Sleep(3 * time.Second)

	b = c.IsKnownSession("foobar")
	require.Equal(t, false, b)
}

func TestIngress(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Second,
		SessionTimeout:  time.Hour,
	})
	require.Equal(t, nil, err)

	c.RegisterAndActivate("foobar", "", "", "")

	c.Ingress("foobar", 1024)

	sessions := c.Active()

	require.Equal(t, 1, len(sessions))
	require.Equal(t, uint64(1024), sessions[0].RxBytes)

	c.Stop()
}

func TestEgress(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Second,
		SessionTimeout:  time.Hour,
	})
	require.Equal(t, nil, err)

	c.RegisterAndActivate("foobar", "", "", "")

	c.Egress("foobar", 1024)

	sessions := c.Active()

	require.Equal(t, 1, len(sessions))
	require.Equal(t, uint64(1024), sessions[0].TxBytes)

	c.Stop()
}

func TestNbSessions(t *testing.T) {
	c, err := newCollector("", nil, nil, CollectorConfig{
		InactiveTimeout: time.Hour,
		SessionTimeout:  time.Hour,
	})
	require.Equal(t, nil, err)

	nsessions := c.Sessions()
	require.Equal(t, uint64(0), nsessions)

	c.Register("foo", "", "", "")

	nsessions = c.Sessions()
	require.Equal(t, uint64(0), nsessions)

	c.Activate("foo")

	nsessions = c.Sessions()
	require.Equal(t, uint64(1), nsessions)

	c.RegisterAndActivate("bar", "", "", "")

	nsessions = c.Sessions()
	require.Equal(t, uint64(2), nsessions)

	c.Stop()

	time.Sleep(2 * time.Second)

	nsessions = c.Sessions()
	require.Equal(t, uint64(0), nsessions)
}
