package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/datarhei/core/v16/config"
	"github.com/datarhei/core/v16/config/store"
	v1 "github.com/datarhei/core/v16/config/v1"
	"github.com/datarhei/core/v16/http/mock"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func getDummyConfigRouter(t *testing.T) (*echo.Echo, store.Store) {
	router := mock.DummyEcho()

	memfs, err := fs.NewMemFilesystem(fs.MemConfig{})
	require.NoError(t, err)

	_, _, err = memfs.WriteFileReader("./mime.types", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	_, _, err = memfs.WriteFileReader("/bin/ffmpeg", strings.NewReader("xxxxx"))
	require.NoError(t, err)

	config, err := store.NewJSON(memfs, "/config.json", nil)
	require.NoError(t, err)

	handler := NewConfig(config)

	router.Add("GET", "/", handler.Get)
	router.Add("PUT", "/", handler.Set)

	return router, config
}

func TestConfigGet(t *testing.T) {
	router, _ := getDummyConfigRouter(t)

	mock.Request(t, http.StatusOK, router, "GET", "/", nil)

	//validate(t, &api.RestreamerConfig{}, response.Data)
}

func TestConfigSetConflict(t *testing.T) {
	router, _ := getDummyConfigRouter(t)

	cfg := config.New(nil)
	cfg.Storage.MimeTypes = "/path/to/mime.types"

	var data bytes.Buffer

	encoder := json.NewEncoder(&data)
	encoder.Encode(cfg)

	mock.Request(t, http.StatusConflict, router, "PUT", "/", &data)
}

func TestConfigSet(t *testing.T) {
	router, store := getDummyConfigRouter(t)

	storedcfg := store.Get()

	require.Equal(t, []string{}, storedcfg.Host.Name)

	var data bytes.Buffer
	encoder := json.NewEncoder(&data)

	// Setting a new v3 config
	cfg := config.New(nil)
	cfg.DB.Dir = "."
	cfg.Storage.Disk.Dir = "."
	cfg.Storage.Disk.Cache.Types.Allow = []string{".aaa"}
	cfg.Storage.Disk.Cache.Types.Block = []string{".zzz"}
	cfg.Host.Name = []string{"foobar.com"}

	encoder.Encode(cfg)

	mock.Request(t, http.StatusOK, router, "PUT", "/", &data)

	storedcfg = store.Get()

	require.Equal(t, []string{"foobar.com"}, storedcfg.Host.Name)
	require.Equal(t, []string{".aaa"}, cfg.Storage.Disk.Cache.Types.Allow)
	require.Equal(t, []string{".zzz"}, cfg.Storage.Disk.Cache.Types.Block)
	require.Equal(t, "cert@datarhei.com", cfg.TLS.Email)

	// Setting a complete v1 config
	cfgv1 := v1.New(nil)
	cfgv1.DB.Dir = "."
	cfgv1.Storage.Disk.Dir = "."
	cfgv1.Storage.Disk.Cache.Types = []string{".bbb"}
	cfgv1.Host.Name = []string{"foobar.com"}

	data.Reset()

	encoder.Encode(cfgv1)

	mock.Request(t, http.StatusOK, router, "PUT", "/", &data)

	storedcfg = store.Get()

	require.Equal(t, []string{"foobar.com"}, storedcfg.Host.Name)
	require.Equal(t, []string{".bbb"}, storedcfg.Storage.Disk.Cache.Types.Allow)
	require.Equal(t, []string{".zzz"}, storedcfg.Storage.Disk.Cache.Types.Block)
	require.Equal(t, "cert@datarhei.com", cfg.TLS.Email)

	// Setting a partial v1 config
	type customconfig struct {
		Version int `json:"version"`
		Storage struct {
			Disk struct {
				Cache struct {
					Types []string `json:"types"`
				} `json:"cache"`
			} `json:"disk"`
		} `json:"storage"`
	}

	customcfg := customconfig{
		Version: 1,
	}

	customcfg.Storage.Disk.Cache.Types = []string{".ccc"}

	data.Reset()

	encoder.Encode(customcfg)

	mock.Request(t, http.StatusOK, router, "PUT", "/", &data)

	storedcfg = store.Get()

	require.Equal(t, []string{"foobar.com"}, storedcfg.Host.Name)
	require.Equal(t, []string{".ccc"}, storedcfg.Storage.Disk.Cache.Types.Allow)
	require.Equal(t, []string{".zzz"}, storedcfg.Storage.Disk.Cache.Types.Block)
	require.Equal(t, "cert@datarhei.com", cfg.TLS.Email)
}
