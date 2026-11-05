package fs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadOnly(t *testing.T) {
	mem, err := NewMemFilesystemFromDir(".", MemConfig{})
	require.NoError(t, err)

	ro, err := NewReadOnlyFilesystem(mem)
	require.NoError(t, err)

	err = ro.Symlink("/readonly.go", "/foobar.go")
	require.Error(t, err)

	_, _, err = ro.WriteFile("/readonly.go", []byte("foobar"))
	require.Error(t, err)

	_, _, err = ro.WriteFileReader("/readonly.go", strings.NewReader("foobar"))
	require.Error(t, err)

	_, _, err = ro.WriteFileSafe("/readonly.go", []byte("foobar"))
	require.Error(t, err)

	err = ro.MkdirAll("/foobar/baz", 0700)
	require.Error(t, err)

	res := ro.Remove("/readonly.go")
	require.Equal(t, int64(-1), res)

	res = ro.RemoveAll()
	require.Equal(t, int64(0), res)

	rop, ok := ro.(PurgeFilesystem)
	require.True(t, ok, "must implement PurgeFilesystem")

	size, _ := ro.Size()
	res = rop.Purge(size)
	require.Equal(t, int64(0), res)

	ros, ok := ro.(SizedFilesystem)
	require.True(t, ok, "must implement SizedFilesystem")

	err = ros.Resize(100, false)
	require.Error(t, err)
}
