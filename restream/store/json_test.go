package store

import (
	"testing"

	"github.com/datarhei/core/v16/io/fs"
	"github.com/stretchr/testify/require"
)

func getFS(t *testing.T) fs.Filesystem {
	fs, err := fs.NewRootedDiskFilesystem(fs.RootedDiskConfig{
		Root: ".",
	})
	require.NoError(t, err)

	info, err := fs.Stat("./fixtures/v4_empty.json")
	require.NoError(t, err)
	require.Equal(t, "/fixtures/v4_empty.json", info.Name())

	return fs
}

func TestNew(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
	})
	require.NoError(t, err)
	require.NotEmpty(t, store)
}

func TestLoad(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
		Filepath:   "./fixtures/v4_empty.json",
	})
	require.NoError(t, err)

	_, err = store.Load()
	require.NoError(t, err)
}

func TestLoadFailed(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
		Filepath:   "./fixtures/v4_invalid.json",
	})
	require.NoError(t, err)

	_, err = store.Load()
	require.Error(t, err)
}

func TestIsEmpty(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
		Filepath:   "./fixtures/v4_empty.json",
	})
	require.NoError(t, err)

	data, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, true, data.IsEmpty())
}

func TestNotExists(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
		Filepath:   "./fixtures/v4_notexist.json",
	})
	require.NoError(t, err)

	data, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, true, data.IsEmpty())
}

func TestStore(t *testing.T) {
	fs := getFS(t)
	fs.Remove("./fixtures/v4_store.json")

	store, err := NewJSON(JSONConfig{
		Filesystem: fs,
		Filepath:   "./fixtures/v4_store.json",
	})
	require.NoError(t, err)

	data, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, true, data.IsEmpty())

	data.Metadata.System["somedata"] = "foobar"

	store.Store(data)

	data2, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, data, data2)

	fs.Remove("./fixtures/v4_store.json")
}

func TestInvalidVersion(t *testing.T) {
	store, err := NewJSON(JSONConfig{
		Filesystem: getFS(t),
		Filepath:   "./fixtures/v3_empty.json",
	})
	require.NoError(t, err)

	data, err := store.Load()
	require.Error(t, err)
	require.Equal(t, true, data.IsEmpty())
}
