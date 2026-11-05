package net

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnonymizeIPString(t *testing.T) {
	_, err := AnonymizeIPString("127.987.475.21")
	require.Error(t, err)

	_, err = AnonymizeIPString("bbd1:xxxx")
	require.Error(t, err)

	_, err = AnonymizeIPString("hello-world")
	require.Error(t, err)

	ipv4 := "192.168.1.42"
	ipv6 := "bbd1:e95a:adbb:b29a:e38b:577f:6f9a:1fa7"

	anonymizedIPv4, err := AnonymizeIPString(ipv4)
	require.NoError(t, err)
	require.Equal(t, "192.168.1.0", anonymizedIPv4)

	anonymizedIPv6, err := AnonymizeIPString(ipv6)
	require.NoError(t, err)
	require.Equal(t, "bbd1:e95a:adbb:b29a::", anonymizedIPv6)
}
