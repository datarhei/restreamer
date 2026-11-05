package rtmp

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	data := [][]string{
		{"/foo/bar", "/foo", "bar"},
		{"/foo/bar?token=abc", "/foo/bar", "abc"},
		{"/foo/bar/abc", "/foo/bar", "abc"},
	}

	for _, d := range data {
		u, err := url.Parse(d[0])
		require.NoError(t, err)

		path, token := getToken(u)

		require.Equal(t, d[1], path, "url=%s", u.String())
		require.Equal(t, d[2], token, "url=%s", u.String())
	}
}
