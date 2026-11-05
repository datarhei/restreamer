package net

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPortrange(t *testing.T) {
	_, err := NewPortrange(1000, 1999)

	require.Nil(t, err, "Valid port range not accepted: %s", err)
}

func TestInvalidPortrange(t *testing.T) {
	_, err := NewPortrange(1999, 1000)

	require.NotNil(t, err, "Invalid port range accepted")
}

func TestOutOfRangePortrange(t *testing.T) {
	p, err := NewPortrange(-1, 70000)

	require.NoError(t, err)

	portrange := p.(*portrange)

	require.Equal(t, 1, portrange.min)
	require.Equal(t, 65535, len(portrange.ports))
}

func TestGetPort(t *testing.T) {
	portrange, _ := NewPortrange(1000, 1999)

	port, err := portrange.Get()

	require.Nil(t, err)
	require.Equal(t, 1000, port)
}

func TestGetPutPort(t *testing.T) {
	portrange, _ := NewPortrange(1000, 1999)

	port, err := portrange.Get()
	require.Nil(t, err)
	require.Equal(t, 1000, port)

	port, err = portrange.Get()
	require.Nil(t, err)
	require.Equal(t, 1001, port)

	portrange.Put(1000)

	port, err = portrange.Get()
	require.Nil(t, err)
	require.Equal(t, 1000, port)
}

func TestPortUnavailable(t *testing.T) {
	portrange, _ := NewPortrange(1000, 1999)

	for i := 0; i < 1000; i++ {
		port, _ := portrange.Get()
		require.Equal(t, 1000+i, port, "at index %d", i)
	}

	port, err := portrange.Get()
	require.NotNil(t, err)
	require.Less(t, port, 0)
}

func TestPutPort(t *testing.T) {
	portrange, _ := NewPortrange(1000, 1999)

	portrange.Put(999)
	portrange.Put(1000)

	portrange.Put(1999)
	portrange.Put(2000)
}

func TestClampRange(t *testing.T) {
	portrange, _ := NewPortrange(65000, 70000)

	port, _ := portrange.Get()

	require.Equal(t, 65000, port)

	portrange.Put(65000)

	for i := 65000; i <= 65535; i++ {
		port, _ := portrange.Get()
		require.Equal(t, i, port, "at index %d", i)
	}

	port, _ = portrange.Get()

	require.Less(t, port, 0)
}

func TestDummyPortranger(t *testing.T) {
	portrange := NewDummyPortrange()

	port, err := portrange.Get()

	require.Error(t, err)
	require.Equal(t, 0, port)

	portrange.Put(42)
}
