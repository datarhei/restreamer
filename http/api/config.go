package api

import (
	"time"

	"github.com/datarhei/core/v16/config"
	v1config "github.com/datarhei/core/v16/config/v1"
	v2config "github.com/datarhei/core/v16/config/v2"
)

// ConfigVersion is used to only unmarshal the version field in order
// find out which SetConfig should be used.
type ConfigVersion struct {
	Version int64 `json:"version"`
}

// ConfigData embeds config.Data
type ConfigData struct {
	config.Data
}

// Config is the config returned by the API
type Config struct {
	CreatedAt time.Time `json:"created_at"`
	LoadedAt  time.Time `json:"loaded_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Config ConfigData `json:"config"`

	Overrides []string `json:"overrides"`
}

type SetConfigV1 struct {
	v1config.Data
}

// NewSetConfigV1 creates a new SetConfigV1 based on the current
// config with downgrading.
func NewSetConfigV1(cfg *config.Config) SetConfigV1 {
	v2data, _ := config.DowngradeV3toV2(&cfg.Data)
	v1data, _ := v2config.DowngradeV2toV1(v2data)

	data := SetConfigV1{
		Data: *v1data,
	}

	return data
}

// MergeTo merges the v1 config into the current config.
func (s *SetConfigV1) MergeTo(cfg *config.Config) {
	v2data, _ := config.DowngradeV3toV2(&cfg.Data)

	v2config.MergeV1ToV2(v2data, &s.Data)
	config.MergeV2toV3(&cfg.Data, v2data)
}

type SetConfigV2 struct {
	v2config.Data
}

// NewSetConfigV2 creates a new SetConfigV2 based on the current
// config with downgrading.
func NewSetConfigV2(cfg *config.Config) SetConfigV2 {
	v2data, _ := config.DowngradeV3toV2(&cfg.Data)

	data := SetConfigV2{
		Data: *v2data,
	}

	return data
}

// MergeTo merges the v2 config into the current config.
func (s *SetConfigV2) MergeTo(cfg *config.Config) {
	config.MergeV2toV3(&cfg.Data, &s.Data)
}

// SetConfig embeds config.Data. It is used to send a new config to the server.
type SetConfig struct {
	config.Data
}

// NewSetConfig converts a config.Config into a SetConfig in order to prepopulate
// a SetConfig with the current values. The uploaded config can have missing fields that
// will be filled with the current values after unmarshalling the JSON.
func NewSetConfig(cfg *config.Config) SetConfig {
	data := SetConfig{
		cfg.Data,
	}

	return data
}

// MergeTo merges a sent config into a config.Config
func (rscfg *SetConfig) MergeTo(cfg *config.Config) {
	cfg.ID = rscfg.ID
	cfg.Name = rscfg.Name
	cfg.Address = rscfg.Address
	cfg.CheckForUpdates = rscfg.CheckForUpdates

	cfg.Log = rscfg.Log
	cfg.DB = rscfg.DB
	cfg.Host = rscfg.Host
	cfg.API = rscfg.API
	cfg.TLS = rscfg.TLS
	cfg.Storage = rscfg.Storage
	cfg.RTMP = rscfg.RTMP
	cfg.SRT = rscfg.SRT
	cfg.FFmpeg = rscfg.FFmpeg
	cfg.Playout = rscfg.Playout
	cfg.Debug = rscfg.Debug
	cfg.Metrics = rscfg.Metrics
	cfg.Sessions = rscfg.Sessions
	cfg.Service = rscfg.Service
	cfg.Router = rscfg.Router
}

// Unmarshal converts a config.Config to a Config.
func (c *Config) Unmarshal(cfg *config.Config) {
	if cfg == nil {
		return
	}

	configData := ConfigData{
		cfg.Data,
	}

	c.CreatedAt = cfg.CreatedAt
	c.LoadedAt = cfg.LoadedAt
	c.UpdatedAt = cfg.UpdatedAt
	c.Config = configData
	c.Overrides = cfg.Overrides()
}

// ConfigError is used to return error messages when uploading a new config
type ConfigError map[string][]string
