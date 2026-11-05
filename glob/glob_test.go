package glob

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatterns(t *testing.T) {
	ok, err := Match("**/a/b/**", "/s3/a/b/test.m3u8", '/')

	require.NoError(t, err)
	require.True(t, ok)

	ok, err = Match("**/a/b/**", "/a/b/test.m3u8", '/')

	require.NoError(t, err)
	require.True(t, ok)

	ok, err = Match("{/memfs,}/a/b/**", "/a/b/test.m3u8", '/')

	require.NoError(t, err)
	require.True(t, ok)
}
