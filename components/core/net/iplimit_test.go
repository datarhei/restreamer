package net

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPLimiterNew(t *testing.T) {
	var err error

	_, err = NewIPLimiter([]string{}, []string{})
	require.Nil(t, err)

	_, err = NewIPLimiter([]string{"::1/128", "127.0.0.1/32", ""}, []string{})
	require.Nil(t, err)

	_, err = NewIPLimiter([]string{}, []string{"::1/128", "127.0.0.1/32", ""})
	require.Nil(t, err)
}

func TestIPLimiterError(t *testing.T) {
	var err error

	_, err = NewIPLimiter([]string{}, []string{})
	require.Nil(t, err)

	_, err = NewIPLimiter([]string{"::1"}, []string{})
	require.NotNil(t, err, "Should not accept invalid IP")

	_, err = NewIPLimiter([]string{}, []string{"::1"})
	require.NotNil(t, err, "Should not accept invalid IP")
}

func TestIPLimiterInvalidIPs(t *testing.T) {
	limiter, _ := NewIPLimiter([]string{}, []string{})

	require.False(t, limiter.IsAllowed(""), "Invalid IP shouldn't be allowed")
}

func TestIPLimiterNoIPs(t *testing.T) {
	limiter, _ := NewIPLimiter([]string{}, []string{})

	require.True(t, limiter.IsAllowed("127.0.0.1"), "IP should be allowed")
}

func TestIPLimiterAllowlist(t *testing.T) {
	limiter, _ := NewIPLimiter([]string{}, []string{"::1/128"})

	require.False(t, limiter.IsAllowed("127.0.0.1"), "Unallowed IP shouldn't be allowed")
	require.True(t, limiter.IsAllowed("::1"), "Allowed IP should be allowed")
}

func TestIPLimiterBlocklist(t *testing.T) {
	limiter, _ := NewIPLimiter([]string{"::1/128"}, []string{})

	require.True(t, limiter.IsAllowed("127.0.0.1"), "Allowed IP should be allowed")
	require.False(t, limiter.IsAllowed("::1"), "Unallowed IP shouldn't be allowed")
}

func TestNullIPLimiter(t *testing.T) {
	limiter := NewNullIPLimiter()

	require.True(t, limiter.IsAllowed("foobar"))
}
