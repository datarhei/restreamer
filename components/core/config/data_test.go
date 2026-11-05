package config

import (
	"testing"

	v2 "github.com/datarhei/core/v16/config/v2"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/stretchr/testify/require"
)

func TestUpgrade(t *testing.T) {
	fs, _ := fs.NewMemFilesystem(fs.MemConfig{})

	v2cfg := v2.New(fs)
	v2cfg.Storage.Disk.Cache.Types = []string{".foo", ".bar"}

	v3cfg, err := UpgradeV2ToV3(&v2cfg.Data, fs)

	require.NoError(t, err)
	require.Equal(t, int64(3), v3cfg.Version)
	require.ElementsMatch(t, []string{".foo", ".bar"}, v3cfg.Storage.Disk.Cache.Types.Allow)
	require.ElementsMatch(t, []string{".m3u8", ".mpd"}, v3cfg.Storage.Disk.Cache.Types.Block)
}

func TestDowngrade(t *testing.T) {
	fs, _ := fs.NewMemFilesystem(fs.MemConfig{})

	v3cfg := New(fs)
	v3cfg.Storage.Disk.Cache.Types.Allow = []string{".foo", ".bar"}

	v2cfg, err := DowngradeV3toV2(&v3cfg.Data)

	require.NoError(t, err)
	require.Equal(t, int64(2), v2cfg.Version)
	require.ElementsMatch(t, []string{".foo", ".bar"}, v2cfg.Storage.Disk.Cache.Types)
}
