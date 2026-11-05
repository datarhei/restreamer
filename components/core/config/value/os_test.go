package value

import (
	"testing"

	"github.com/datarhei/core/v16/io/fs"
	"github.com/stretchr/testify/require"
)

func TestMustDirValue(t *testing.T) {
	memfs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	_, err = memfs.Stat("/foobar")
	require.Error(t, err)

	var x string

	val := NewMustDir(&x, "./foobar", memfs)

	require.Equal(t, "./foobar", val.String())
	require.NoError(t, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	info, err := memfs.Stat("/foobar")
	require.NoError(t, err)
	require.True(t, info.IsDir())

	x = "/bar/foo"

	require.Equal(t, "/bar/foo", val.String())

	_, err = memfs.Stat("/bar/foo")
	require.Error(t, err)

	require.NoError(t, val.Validate())

	info, err = memfs.Stat("/bar/foo")
	require.NoError(t, err)
	require.True(t, info.IsDir())

	memfs.WriteFile("/foo/bar", []byte("hello"))

	val.Set("/foo/bar")

	require.Error(t, val.Validate())
}

func TestDirValue(t *testing.T) {
	memfs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	var x string

	val := NewDir(&x, "/foobar", memfs)

	require.Equal(t, "/foobar", val.String())
	require.Error(t, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	err = memfs.MkdirAll("/foobar", 0755)
	require.NoError(t, err)

	require.NoError(t, val.Validate())

	_, _, err = memfs.WriteFile("/foo/bar", []byte("hello"))
	require.NoError(t, err)

	val.Set("/foo/bar")

	require.Error(t, val.Validate())
}

func TestFileValue(t *testing.T) {
	memfs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	var x string

	val := NewFile(&x, "/foobar", memfs)

	require.Equal(t, "/foobar", val.String())
	require.Error(t, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	_, _, err = memfs.WriteFile("/foobar", []byte("hello"))
	require.NoError(t, err)

	require.NoError(t, val.Validate())

	err = memfs.MkdirAll("/foo/bar", 0755)
	require.NoError(t, err)

	val.Set("/foo/bar")

	require.Error(t, val.Validate())
}

func TestExecValue(t *testing.T) {
	memfs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	var x string

	val := NewExec(&x, "/foobar", memfs)

	require.Equal(t, "/foobar", val.String())
	require.Error(t, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	_, _, err = memfs.WriteFile("/foobar", []byte("hello"))
	require.NoError(t, err)

	require.NoError(t, val.Validate())

	err = memfs.MkdirAll("/foo/bar", 0755)
	require.NoError(t, err)

	val.Set("/foo/bar")

	require.Error(t, val.Validate())
}

func TestAbsolutePathValue(t *testing.T) {
	var x string

	val := NewAbsolutePath(&x, "foobar")

	require.Equal(t, "foobar", val.String())
	require.Error(t, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = "/foobaz"

	require.Equal(t, "/foobaz", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("/fooboz")

	require.Equal(t, "/fooboz", x)
}
