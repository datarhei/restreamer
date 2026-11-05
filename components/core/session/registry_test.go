package session

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	r, err := New(Config{})
	require.Equal(t, nil, err)

	_, err = r.Register("", CollectorConfig{})
	require.NotEqual(t, nil, err)

	_, err = r.Register("../foo/bar", CollectorConfig{})
	require.NotEqual(t, nil, err)

	_, err = r.Register("foobar", CollectorConfig{})
	require.Equal(t, nil, err)

	_, err = r.Register("foobar", CollectorConfig{})
	require.NotEqual(t, nil, err)
}

func TestUnregister(t *testing.T) {
	r, err := New(Config{})
	require.Equal(t, nil, err)

	_, err = r.Register("foobar", CollectorConfig{})
	require.Equal(t, nil, err)

	err = r.Unregister("foobar")
	require.Equal(t, nil, err)

	err = r.Unregister("foobar")
	require.NotEqual(t, nil, err)
}

func TestCollectors(t *testing.T) {
	r, err := New(Config{})
	require.Equal(t, nil, err)

	c := r.Collectors()
	require.Equal(t, []string{}, c)

	_, err = r.Register("foobar", CollectorConfig{})
	require.Equal(t, nil, err)

	c = r.Collectors()
	require.Equal(t, []string{"foobar"}, c)

	err = r.Unregister("foobar")
	require.Equal(t, nil, err)

	c = r.Collectors()
	require.Equal(t, []string{}, c)
}

func TestGetCollector(t *testing.T) {
	r, err := New(Config{})
	require.Equal(t, nil, err)

	c := r.Collector("foobar")
	require.Equal(t, nil, c)

	_, err = r.Register("foobar", CollectorConfig{})
	require.Equal(t, nil, err)

	c = r.Collector("foobar")
	require.NotEqual(t, nil, c)

	err = r.Unregister("foobar")
	require.Equal(t, nil, err)

	c = r.Collector("foobar")
	require.Equal(t, nil, c)
}

func TestUnregisterAll(t *testing.T) {
	r, err := New(Config{})
	require.Equal(t, nil, err)

	_, err = r.Register("foo", CollectorConfig{})
	require.Equal(t, nil, err)

	_, err = r.Register("bar", CollectorConfig{})
	require.Equal(t, nil, err)

	c := r.Collectors()
	require.ElementsMatch(t, []string{"foo", "bar"}, c)

	r.UnregisterAll()

	c = r.Collectors()
	require.Equal(t, []string{}, c)
}
