package config

import (
	"strings"
	"testing"

	"github.com/datarhei/core/v16/config/vars"
	"github.com/datarhei/core/v16/io/fs"

	"github.com/stretchr/testify/require"
)

func TestConfigCopy(t *testing.T) {
	fs, _ := fs.NewMemFilesystem(fs.MemConfig{})
	config1 := New(fs)

	config1.Version = 42
	config1.DB.Dir = "foo"

	val1, _ := config1.Get("version")
	val2, _ := config1.Get("db.dir")
	val3, _ := config1.Get("host.name")

	require.Equal(t, "42", val1)
	require.Equal(t, "foo", val2)
	require.Equal(t, "(empty)", val3)

	config1.Set("host.name", "foo.com")
	val3, _ = config1.Get("host.name")
	require.Equal(t, "foo.com", val3)

	config2 := config1.Clone()

	require.Equal(t, int64(42), config2.Version)
	require.Equal(t, "foo", config2.DB.Dir)
	require.Equal(t, []string{"foo.com"}, config2.Host.Name)

	config1.Set("version", "77")

	require.Equal(t, int64(77), config1.Version)
	require.Equal(t, int64(42), config2.Version)

	config1.Set("db.dir", "bar")

	require.Equal(t, "bar", config1.DB.Dir)
	require.Equal(t, "foo", config2.DB.Dir)

	config2.DB.Dir = "baz"

	require.Equal(t, "bar", config1.DB.Dir)
	require.Equal(t, "baz", config2.DB.Dir)

	config1.Host.Name[0] = "bar.com"

	require.Equal(t, []string{"bar.com"}, config1.Host.Name)
	require.Equal(t, []string{"foo.com"}, config2.Host.Name)
}

func TestValidateDefault(t *testing.T) {
	fs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	size, fresh, err := fs.WriteFileReader("./mime.types", strings.NewReader("xxxxx"))
	require.Equal(t, int64(5), size)
	require.Equal(t, true, fresh)
	require.NoError(t, err)

	_, _, err = fs.WriteFileReader("/bin/ffmpeg", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	cfg := New(fs)

	cfg.Validate(true)

	errors := []string{}
	cfg.Messages(func(level string, v vars.Variable, message string) {
		if level == "error" {
			errors = append(errors, message)
		}
	})

	require.Equal(t, 0, len(cfg.Overrides()))
	require.Equal(t, false, cfg.HasErrors(), errors)
}
