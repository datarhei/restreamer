package v2

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/datarhei/core/v16/config/copy"
	v1 "github.com/datarhei/core/v16/config/v1"
	"github.com/datarhei/core/v16/config/value"
	"github.com/datarhei/core/v16/io/fs"
)

type Data struct {
	CreatedAt       time.Time `json:"created_at"`
	LoadedAt        time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	Version         int64     `json:"version" jsonschema:"minimum=2,maximum=2" format:"int64"`
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	CheckForUpdates bool      `json:"update_check"`
	Log             struct {
		Level    string   `json:"level" enums:"debug,info,warn,error,silent" jsonschema:"enum=debug,enum=info,enum=warn,enum=error,enum=silent"`
		Topics   []string `json:"topics"`
		MaxLines int      `json:"max_lines" format:"int"`
	} `json:"log"`
	DB struct {
		Dir string `json:"dir"`
	} `json:"db"`
	Host struct {
		Name []string `json:"name"`
		Auto bool     `json:"auto"`
	} `json:"host"`
	API struct {
		ReadOnly bool `json:"read_only"`
		Access   struct {
			HTTP struct {
				Allow []string `json:"allow"`
				Block []string `json:"block"`
			} `json:"http"`
			HTTPS struct {
				Allow []string `json:"allow"`
				Block []string `json:"block"`
			} `json:"https"`
		} `json:"access"`
		Auth struct {
			Enable           bool   `json:"enable"`
			DisableLocalhost bool   `json:"disable_localhost"`
			Username         string `json:"username"`
			Password         string `json:"password"`
			JWT              struct {
				Secret string `json:"secret"`
			} `json:"jwt"`
			Auth0 struct {
				Enable  bool                `json:"enable"`
				Tenants []value.Auth0Tenant `json:"tenants"`
			} `json:"auth0"`
		} `json:"auth"`
	} `json:"api"`
	TLS struct {
		Address  string `json:"address"`
		Enable   bool   `json:"enable"`
		Auto     bool   `json:"auto"`
		CertFile string `json:"cert_file"`
		KeyFile  string `json:"key_file"`
	} `json:"tls"`
	Storage struct {
		Disk struct {
			Dir   string `json:"dir"`
			Size  int64  `json:"max_size_mbytes" format:"int64"`
			Cache struct {
				Enable   bool     `json:"enable"`
				Size     uint64   `json:"max_size_mbytes" format:"uint64"`
				TTL      int64    `json:"ttl_seconds" format:"int64"`
				FileSize uint64   `json:"max_file_size_mbytes" format:"uint64"`
				Types    []string `json:"types"`
			} `json:"cache"`
		} `json:"disk"`
		Memory struct {
			Auth struct {
				Enable   bool   `json:"enable"`
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"auth"`
			Size  int64 `json:"max_size_mbytes" format:"int64"`
			Purge bool  `json:"purge"`
		} `json:"memory"`
		CORS struct {
			Origins []string `json:"origins"`
		} `json:"cors"`
		MimeTypes string `json:"mimetypes_file"`
	} `json:"storage"`
	RTMP struct {
		Enable     bool   `json:"enable"`
		EnableTLS  bool   `json:"enable_tls"`
		Address    string `json:"address"`
		AddressTLS string `json:"address_tls"`
		App        string `json:"app"`
		Token      string `json:"token"`
	} `json:"rtmp"`
	SRT struct {
		Enable     bool   `json:"enable"`
		Address    string `json:"address"`
		Passphrase string `json:"passphrase"`
		Token      string `json:"token"`
		Log        struct {
			Enable bool     `json:"enable"`
			Topics []string `json:"topics"`
		} `json:"log"`
	} `json:"srt"`
	FFmpeg struct {
		Binary       string `json:"binary"`
		MaxProcesses int64  `json:"max_processes" format:"int64"`
		Access       struct {
			Input struct {
				Allow []string `json:"allow"`
				Block []string `json:"block"`
			} `json:"input"`
			Output struct {
				Allow []string `json:"allow"`
				Block []string `json:"block"`
			} `json:"output"`
		} `json:"access"`
		Log struct {
			MaxLines   int `json:"max_lines" format:"int"`
			MaxHistory int `json:"max_history" format:"int"`
		} `json:"log"`
	} `json:"ffmpeg"`
	Playout struct {
		Enable  bool `json:"enable"`
		MinPort int  `json:"min_port" format:"int"`
		MaxPort int  `json:"max_port" format:"int"`
	} `json:"playout"`
	Debug struct {
		Profiling bool `json:"profiling"`
		ForceGC   int  `json:"force_gc" format:"int"`
	} `json:"debug"`
	Metrics struct {
		Enable           bool  `json:"enable"`
		EnablePrometheus bool  `json:"enable_prometheus"`
		Range            int64 `json:"range_sec" format:"int64"`    // seconds
		Interval         int64 `json:"interval_sec" format:"int64"` // seconds
	} `json:"metrics"`
	Sessions struct {
		Enable          bool     `json:"enable"`
		IPIgnoreList    []string `json:"ip_ignorelist"`
		SessionTimeout  int      `json:"session_timeout_sec" format:"int"`
		Persist         bool     `json:"persist"`
		PersistInterval int      `json:"persist_interval_sec" format:"int"`
		MaxBitrate      uint64   `json:"max_bitrate_mbit" format:"uint64"`
		MaxSessions     uint64   `json:"max_sessions" format:"uint64"`
	} `json:"sessions"`
	Service struct {
		Enable bool   `json:"enable"`
		Token  string `json:"token"`
		URL    string `json:"url"`
	} `json:"service"`
	Router struct {
		BlockedPrefixes []string          `json:"blocked_prefixes"`
		Routes          map[string]string `json:"routes"`
		UIPath          string            `json:"ui_path"`
	} `json:"router"`
}

func UpgradeV1ToV2(d *v1.Data, fs fs.Filesystem) (*Data, error) {
	cfg := New(fs)

	return MergeV1ToV2(&cfg.Data, d)
}

// Migrate will migrate some settings, depending on the version it finds. Migrations
// are only going upwards, i.e. from a lower version to a higher version.
func MergeV1ToV2(data *Data, d *v1.Data) (*Data, error) {
	data.CreatedAt = d.CreatedAt
	data.LoadedAt = d.LoadedAt
	data.UpdatedAt = d.UpdatedAt

	data.ID = d.ID
	data.Name = d.Name
	data.Address = d.Address
	data.CheckForUpdates = d.CheckForUpdates

	data.Log = d.Log
	data.DB = d.DB
	data.Host = d.Host
	data.API = d.API
	data.TLS = d.TLS
	data.Storage = d.Storage
	data.SRT = d.SRT
	data.FFmpeg = d.FFmpeg
	data.Playout = d.Playout
	data.Debug = d.Debug
	data.Metrics = d.Metrics
	data.Sessions = d.Sessions
	data.Service = d.Service
	data.Router = d.Router

	data.Log.Topics = copy.Slice(d.Log.Topics)

	data.Host.Name = copy.Slice(d.Host.Name)

	data.API.Access.HTTP.Allow = copy.Slice(d.API.Access.HTTP.Allow)
	data.API.Access.HTTP.Block = copy.Slice(d.API.Access.HTTP.Block)
	data.API.Access.HTTPS.Allow = copy.Slice(d.API.Access.HTTPS.Allow)
	data.API.Access.HTTPS.Block = copy.Slice(d.API.Access.HTTPS.Block)

	data.API.Auth.Auth0.Tenants = copy.TenantSlice(d.API.Auth.Auth0.Tenants)

	data.Storage.CORS.Origins = copy.Slice(d.Storage.CORS.Origins)

	data.FFmpeg.Access.Input.Allow = copy.Slice(d.FFmpeg.Access.Input.Allow)
	data.FFmpeg.Access.Input.Block = copy.Slice(d.FFmpeg.Access.Input.Block)
	data.FFmpeg.Access.Output.Allow = copy.Slice(d.FFmpeg.Access.Output.Allow)
	data.FFmpeg.Access.Output.Block = copy.Slice(d.FFmpeg.Access.Output.Block)

	data.Sessions.IPIgnoreList = copy.Slice(d.Sessions.IPIgnoreList)

	data.SRT.Log.Topics = copy.Slice(d.SRT.Log.Topics)

	data.Router.BlockedPrefixes = copy.Slice(d.Router.BlockedPrefixes)
	data.Router.Routes = copy.StringMap(d.Router.Routes)

	// Actual changes
	data.RTMP.Enable = d.RTMP.Enable
	data.RTMP.EnableTLS = d.RTMP.EnableTLS
	data.RTMP.Address = d.RTMP.Address
	data.RTMP.App = d.RTMP.App
	data.RTMP.Token = d.RTMP.Token

	if !strings.HasPrefix(data.RTMP.App, "/") {
		data.RTMP.App = "/" + data.RTMP.App
	}

	if d.RTMP.EnableTLS {
		data.RTMP.Enable = true
		data.RTMP.AddressTLS = data.RTMP.Address
		host, sport, err := net.SplitHostPort(data.RTMP.Address)
		if err != nil {
			return nil, fmt.Errorf("migrating rtmp.address to rtmp.address_tls failed: %w", err)
		}

		port, err := strconv.Atoi(sport)
		if err != nil {
			return nil, fmt.Errorf("migrating rtmp.address to rtmp.address_tls failed: %w", err)
		}

		data.RTMP.Address = net.JoinHostPort(host, strconv.Itoa(port-1))
	}

	data.Version = 2

	return data, nil
}

func DowngradeV2toV1(d *Data) (*v1.Data, error) {
	data := &v1.Data{}

	data.CreatedAt = d.CreatedAt
	data.LoadedAt = d.LoadedAt
	data.UpdatedAt = d.UpdatedAt

	data.ID = d.ID
	data.Name = d.Name
	data.Address = d.Address
	data.CheckForUpdates = d.CheckForUpdates

	data.Log = d.Log
	data.DB = d.DB
	data.Host = d.Host
	data.API = d.API
	data.TLS = d.TLS
	data.Storage = d.Storage
	data.SRT = d.SRT
	data.FFmpeg = d.FFmpeg
	data.Playout = d.Playout
	data.Debug = d.Debug
	data.Metrics = d.Metrics
	data.Sessions = d.Sessions
	data.Service = d.Service
	data.Router = d.Router

	data.Log.Topics = copy.Slice(d.Log.Topics)

	data.Host.Name = copy.Slice(d.Host.Name)

	data.API.Access.HTTP.Allow = copy.Slice(d.API.Access.HTTP.Allow)
	data.API.Access.HTTP.Block = copy.Slice(d.API.Access.HTTP.Block)
	data.API.Access.HTTPS.Allow = copy.Slice(d.API.Access.HTTPS.Allow)
	data.API.Access.HTTPS.Block = copy.Slice(d.API.Access.HTTPS.Block)

	data.API.Auth.Auth0.Tenants = copy.TenantSlice(d.API.Auth.Auth0.Tenants)

	data.Storage.CORS.Origins = copy.Slice(d.Storage.CORS.Origins)

	data.FFmpeg.Access.Input.Allow = copy.Slice(d.FFmpeg.Access.Input.Allow)
	data.FFmpeg.Access.Input.Block = copy.Slice(d.FFmpeg.Access.Input.Block)
	data.FFmpeg.Access.Output.Allow = copy.Slice(d.FFmpeg.Access.Output.Allow)
	data.FFmpeg.Access.Output.Block = copy.Slice(d.FFmpeg.Access.Output.Block)

	data.Sessions.IPIgnoreList = copy.Slice(d.Sessions.IPIgnoreList)

	data.SRT.Log.Topics = copy.Slice(d.SRT.Log.Topics)

	data.Router.BlockedPrefixes = copy.Slice(d.Router.BlockedPrefixes)
	data.Router.Routes = copy.StringMap(d.Router.Routes)

	// Actual changes
	data.RTMP.Enable = d.RTMP.Enable
	data.RTMP.EnableTLS = d.RTMP.EnableTLS
	data.RTMP.Address = d.RTMP.Address
	data.RTMP.App = d.RTMP.App
	data.RTMP.Token = d.RTMP.Token

	data.Version = 1

	return data, nil
}
