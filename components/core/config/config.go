// Package config implements types for handling the configuation for the app.
package config

import (
	"context"
	"net"
	"time"

	"github.com/datarhei/core/v16/config/copy"
	"github.com/datarhei/core/v16/config/value"
	"github.com/datarhei/core/v16/config/vars"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/math/rand"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/google/uuid"
)

/*
type Config interface {
	// Merge merges the values of the known environment variables into the configuration
	Merge()

	// Validate validates the current state of the Config for completeness and sanity. Errors are
	// written to the log. Use resetLogs to indicate to reset the logs prior validation.
	Validate(resetLogs bool)

	// Messages calls for each log entry the provided callback. The level has the values 'error', 'warn', or 'info'.
	// The name is the name of the configuration value, e.g. 'api.auth.enable'. The message is the log message.
	Messages(logger func(level string, v vars.Variable, message string))

	// HasErrors returns whether there are some error messages in the log.
	HasErrors() bool

	// Overrides returns a list of configuration value names that have been overriden by an environment variable.
	Overrides() []string

	Get(name string) (string, error)
	Set(name, val string) error
}
*/

const version int64 = 3

// Make sure that the config.Config interface is satisfied
//var _ config.Config = &Config{}

// Config is a wrapper for Data
type Config struct {
	fs   fs.Filesystem
	vars vars.Variables

	Data
}

// New returns a Config which is initialized with its default values
func New(f fs.Filesystem) *Config {
	config := &Config{
		fs: f,
	}

	if config.fs == nil {
		config.fs, _ = fs.NewMemFilesystem(fs.MemConfig{})
	}

	config.init()

	return config
}

func (d *Config) Get(name string) (string, error) {
	return d.vars.Get(name)
}

func (d *Config) Set(name, val string) error {
	return d.vars.Set(name, val)
}

// NewConfigFrom returns a clone of a Config
func (d *Config) Clone() *Config {
	data := New(d.fs)

	data.CreatedAt = d.CreatedAt
	data.LoadedAt = d.LoadedAt
	data.UpdatedAt = d.UpdatedAt

	data.Version = d.Version
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
	data.RTMP = d.RTMP
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
	data.Storage.Disk.Cache.Types.Allow = copy.Slice(d.Storage.Disk.Cache.Types.Allow)
	data.Storage.Disk.Cache.Types.Block = copy.Slice(d.Storage.Disk.Cache.Types.Block)
	data.Storage.S3 = copy.Slice(d.Storage.S3)

	data.FFmpeg.Access.Input.Allow = copy.Slice(d.FFmpeg.Access.Input.Allow)
	data.FFmpeg.Access.Input.Block = copy.Slice(d.FFmpeg.Access.Input.Block)
	data.FFmpeg.Access.Output.Allow = copy.Slice(d.FFmpeg.Access.Output.Allow)
	data.FFmpeg.Access.Output.Block = copy.Slice(d.FFmpeg.Access.Output.Block)

	data.Sessions.IPIgnoreList = copy.Slice(d.Sessions.IPIgnoreList)

	data.SRT.Log.Topics = copy.Slice(d.SRT.Log.Topics)

	data.Router.BlockedPrefixes = copy.Slice(d.Router.BlockedPrefixes)
	data.Router.Routes = copy.StringMap(d.Router.Routes)

	data.vars.Transfer(&d.vars)

	return data
}

func (d *Config) init() {
	d.vars.Register(value.NewInt64(&d.Version, version), "version", "", nil, "Configuration file layout version", true, false)
	d.vars.Register(value.NewTime(&d.CreatedAt, time.Now()), "created_at", "", nil, "Configuration file creation time", false, false)
	d.vars.Register(value.NewString(&d.ID, uuid.New().String()), "id", "CORE_ID", nil, "ID for this instance", true, false)
	d.vars.Register(value.NewString(&d.Name, haikunator.New().Haikunate()), "name", "CORE_NAME", nil, "A human readable name for this instance", false, false)
	d.vars.Register(value.NewAddress(&d.Address, ":8080"), "address", "CORE_ADDRESS", nil, "HTTP listening address", false, false)
	d.vars.Register(value.NewBool(&d.CheckForUpdates, true), "update_check", "CORE_UPDATE_CHECK", nil, "Check for updates and send anonymized data", false, false)

	// Log
	d.vars.Register(value.NewString(&d.Log.Level, "info"), "log.level", "CORE_LOG_LEVEL", nil, "Loglevel: silent, error, warn, info, debug", false, false)
	d.vars.Register(value.NewStringList(&d.Log.Topics, []string{}, ","), "log.topics", "CORE_LOG_TOPICS", nil, "Show only selected log topics", false, false)
	d.vars.Register(value.NewInt(&d.Log.MaxLines, 1000), "log.max_lines", "CORE_LOG_MAX_LINES", []string{"CORE_LOG_MAXLINES"}, "Number of latest log lines to keep in memory", false, false)

	// DB
	d.vars.Register(value.NewMustDir(&d.DB.Dir, "./config", d.fs), "db.dir", "CORE_DB_DIR", nil, "Directory for holding the operational data", false, false)

	// Host
	d.vars.Register(value.NewStringList(&d.Host.Name, []string{}, ","), "host.name", "CORE_HOST_NAME", nil, "Comma separated list of public host/domain names or IPs", false, false)
	d.vars.Register(value.NewBool(&d.Host.Auto, true), "host.auto", "CORE_HOST_AUTO", nil, "Enable detection of public IP addresses", false, false)

	// API
	d.vars.Register(value.NewBool(&d.API.ReadOnly, false), "api.read_only", "CORE_API_READ_ONLY", nil, "Allow only ready only access to the API", false, false)
	d.vars.Register(value.NewCIDRList(&d.API.Access.HTTP.Allow, []string{}, ","), "api.access.http.allow", "CORE_API_ACCESS_HTTP_ALLOW", nil, "List of IPs in CIDR notation (HTTP traffic)", false, false)
	d.vars.Register(value.NewCIDRList(&d.API.Access.HTTP.Block, []string{}, ","), "api.access.http.block", "CORE_API_ACCESS_HTTP_BLOCK", nil, "List of IPs in CIDR notation (HTTP traffic)", false, false)
	d.vars.Register(value.NewCIDRList(&d.API.Access.HTTPS.Allow, []string{}, ","), "api.access.https.allow", "CORE_API_ACCESS_HTTPS_ALLOW", nil, "List of IPs in CIDR notation (HTTPS traffic)", false, false)
	d.vars.Register(value.NewCIDRList(&d.API.Access.HTTPS.Block, []string{}, ","), "api.access.https.block", "CORE_API_ACCESS_HTTPS_BLOCK", nil, "List of IPs in CIDR notation (HTTPS traffic)", false, false)
	d.vars.Register(value.NewBool(&d.API.Auth.Enable, false), "api.auth.enable", "CORE_API_AUTH_ENABLE", nil, "Enable authentication for all clients", false, false)
	d.vars.Register(value.NewBool(&d.API.Auth.DisableLocalhost, false), "api.auth.disable_localhost", "CORE_API_AUTH_DISABLE_LOCALHOST", nil, "Disable authentication for clients from localhost", false, false)
	d.vars.Register(value.NewString(&d.API.Auth.Username, ""), "api.auth.username", "CORE_API_AUTH_USERNAME", []string{"RS_USERNAME"}, "Username", false, false)
	d.vars.Register(value.NewString(&d.API.Auth.Password, ""), "api.auth.password", "CORE_API_AUTH_PASSWORD", []string{"RS_PASSWORD"}, "Password", false, true)

	// Auth JWT
	d.vars.Register(value.NewString(&d.API.Auth.JWT.Secret, rand.String(32)), "api.auth.jwt.secret", "CORE_API_AUTH_JWT_SECRET", nil, "JWT secret, leave empty for generating a random value", false, true)

	// Auth Auth0
	d.vars.Register(value.NewBool(&d.API.Auth.Auth0.Enable, false), "api.auth.auth0.enable", "CORE_API_AUTH_AUTH0_ENABLE", nil, "Enable Auth0", false, false)
	d.vars.Register(value.NewTenantList(&d.API.Auth.Auth0.Tenants, []value.Auth0Tenant{}, ","), "api.auth.auth0.tenants", "CORE_API_AUTH_AUTH0_TENANTS", nil, "List of Auth0 tenants", false, false)

	// TLS
	d.vars.Register(value.NewAddress(&d.TLS.Address, ":8181"), "tls.address", "CORE_TLS_ADDRESS", nil, "HTTPS listening address", false, false)
	d.vars.Register(value.NewBool(&d.TLS.Enable, false), "tls.enable", "CORE_TLS_ENABLE", nil, "Enable HTTPS", false, false)
	d.vars.Register(value.NewBool(&d.TLS.Auto, false), "tls.auto", "CORE_TLS_AUTO", nil, "Enable Let's Encrypt certificate", false, false)
	d.vars.Register(value.NewEmail(&d.TLS.Email, "cert@datarhei.com"), "tls.email", "CORE_TLS_EMAIL", nil, "Email for Let's Encrypt registration", false, false)
	d.vars.Register(value.NewFile(&d.TLS.CertFile, "", d.fs), "tls.cert_file", "CORE_TLS_CERT_FILE", []string{"CORE_TLS_CERTFILE"}, "Path to certificate file in PEM format", false, false)
	d.vars.Register(value.NewFile(&d.TLS.KeyFile, "", d.fs), "tls.key_file", "CORE_TLS_KEY_FILE", []string{"CORE_TLS_KEYFILE"}, "Path to key file in PEM format", false, false)

	// Storage
	d.vars.Register(value.NewFile(&d.Storage.MimeTypes, "./mime.types", d.fs), "storage.mimetypes_file", "CORE_STORAGE_MIMETYPES_FILE", []string{"CORE_MIMETYPES_FILE"}, "Path to file with mime-types", false, false)

	// Storage (Disk)
	d.vars.Register(value.NewMustDir(&d.Storage.Disk.Dir, "./data", d.fs), "storage.disk.dir", "CORE_STORAGE_DISK_DIR", nil, "Directory on disk, exposed on /", false, false)
	d.vars.Register(value.NewInt64(&d.Storage.Disk.Size, 0), "storage.disk.max_size_mbytes", "CORE_STORAGE_DISK_MAX_SIZE_MBYTES", []string{"CORE_STORAGE_DISK_MAXSIZEMBYTES"}, "Max. allowed megabytes for storage.disk.dir, 0 for unlimited", false, false)
	d.vars.Register(value.NewBool(&d.Storage.Disk.Cache.Enable, true), "storage.disk.cache.enable", "CORE_STORAGE_DISK_CACHE_ENABLE", nil, "Enable cache for /", false, false)
	d.vars.Register(value.NewUint64(&d.Storage.Disk.Cache.Size, 0), "storage.disk.cache.max_size_mbytes", "CORE_STORAGE_DISK_CACHE_MAX_SIZE_MBYTES", []string{"CORE_STORAGE_DISK_CACHE_MAXSIZEMBYTES"}, "Max. allowed cache size, 0 for unlimited", false, false)
	d.vars.Register(value.NewInt64(&d.Storage.Disk.Cache.TTL, 300), "storage.disk.cache.ttl_seconds", "CORE_STORAGE_DISK_CACHE_TTL_SECONDS", []string{"CORE_STORAGE_DISK_CACHE_TTLSECONDS"}, "Seconds to keep files in cache", false, false)
	d.vars.Register(value.NewUint64(&d.Storage.Disk.Cache.FileSize, 1), "storage.disk.cache.max_file_size_mbytes", "CORE_STORAGE_DISK_CACHE_MAX_FILE_SIZE_MBYTES", []string{"CORE_STORAGE_DISK_CACHE_MAXFILESIZEMBYTES"}, "Max. file size to put in cache", false, false)
	d.vars.Register(value.NewStringList(&d.Storage.Disk.Cache.Types.Allow, []string{}, " "), "storage.disk.cache.type.allow", "CORE_STORAGE_DISK_CACHE_TYPES_ALLOW", []string{"CORE_STORAGE_DISK_CACHE_TYPES"}, "File extensions to cache, empty for all", false, false)
	d.vars.Register(value.NewStringList(&d.Storage.Disk.Cache.Types.Block, []string{".m3u8", ".mpd"}, " "), "storage.disk.cache.type.block", "CORE_STORAGE_DISK_CACHE_TYPES_BLOCK", nil, "File extensions not to cache, empty for none", false, false)

	// Storage (Memory)
	d.vars.Register(value.NewBool(&d.Storage.Memory.Auth.Enable, true), "storage.memory.auth.enable", "CORE_STORAGE_MEMORY_AUTH_ENABLE", nil, "Enable basic auth for PUT,POST, and DELETE on /memfs", false, false)
	d.vars.Register(value.NewString(&d.Storage.Memory.Auth.Username, "admin"), "storage.memory.auth.username", "CORE_STORAGE_MEMORY_AUTH_USERNAME", nil, "Username for Basic-Auth of /memfs", false, false)
	d.vars.Register(value.NewString(&d.Storage.Memory.Auth.Password, rand.StringAlphanumeric(18)), "storage.memory.auth.password", "CORE_STORAGE_MEMORY_AUTH_PASSWORD", nil, "Password for Basic-Auth of /memfs", false, true)
	d.vars.Register(value.NewInt64(&d.Storage.Memory.Size, 0), "storage.memory.max_size_mbytes", "CORE_STORAGE_MEMORY_MAX_SIZE_MBYTES", []string{"CORE_STORAGE_MEMORY_MAXSIZEMBYTES"}, "Max. allowed megabytes for /memfs, 0 for unlimited", false, false)
	d.vars.Register(value.NewBool(&d.Storage.Memory.Purge, false), "storage.memory.purge", "CORE_STORAGE_MEMORY_PURGE", nil, "Automatically remove the oldest files if /memfs is full", false, false)

	// Storage (S3)
	d.vars.Register(value.NewS3StorageListValue(&d.Storage.S3, []value.S3Storage{}, "|"), "storage.s3", "CORE_STORAGE_S3", nil, "List of S3 storage URLS", false, false)

	// Storage (CORS)
	d.vars.Register(value.NewCORSOrigins(&d.Storage.CORS.Origins, []string{"*"}, ","), "storage.cors.origins", "CORE_STORAGE_CORS_ORIGINS", nil, "Allowed CORS origins for /memfs and /data", false, false)

	// RTMP
	d.vars.Register(value.NewBool(&d.RTMP.Enable, false), "rtmp.enable", "CORE_RTMP_ENABLE", nil, "Enable RTMP server", false, false)
	d.vars.Register(value.NewBool(&d.RTMP.EnableTLS, false), "rtmp.enable_tls", "CORE_RTMP_ENABLE_TLS", nil, "Enable RTMPS server instead of RTMP", false, false)
	d.vars.Register(value.NewAddress(&d.RTMP.Address, ":1935"), "rtmp.address", "CORE_RTMP_ADDRESS", nil, "RTMP server listen address", false, false)
	d.vars.Register(value.NewAddress(&d.RTMP.AddressTLS, ":1936"), "rtmp.address_tls", "CORE_RTMP_ADDRESS_TLS", nil, "RTMPS server listen address", false, false)
	d.vars.Register(value.NewAbsolutePath(&d.RTMP.App, "/"), "rtmp.app", "CORE_RTMP_APP", nil, "RTMP app for publishing", false, false)
	d.vars.Register(value.NewString(&d.RTMP.Token, ""), "rtmp.token", "CORE_RTMP_TOKEN", nil, "RTMP token for publishing and playing", false, true)

	// SRT
	d.vars.Register(value.NewBool(&d.SRT.Enable, false), "srt.enable", "CORE_SRT_ENABLE", nil, "Enable SRT server", false, false)
	d.vars.Register(value.NewAddress(&d.SRT.Address, ":6000"), "srt.address", "CORE_SRT_ADDRESS", nil, "SRT server listen address", false, false)
	d.vars.Register(value.NewString(&d.SRT.Passphrase, ""), "srt.passphrase", "CORE_SRT_PASSPHRASE", nil, "SRT encryption passphrase", false, true)
	d.vars.Register(value.NewString(&d.SRT.Token, ""), "srt.token", "CORE_SRT_TOKEN", nil, "SRT token for publishing and playing", false, true)
	d.vars.Register(value.NewBool(&d.SRT.Log.Enable, false), "srt.log.enable", "CORE_SRT_LOG_ENABLE", nil, "Enable SRT server logging", false, false)
	d.vars.Register(value.NewStringList(&d.SRT.Log.Topics, []string{}, ","), "srt.log.topics", "CORE_SRT_LOG_TOPICS", nil, "List of topics to log", false, false)

	// FFmpeg
	d.vars.Register(value.NewExec(&d.FFmpeg.Binary, "ffmpeg", d.fs), "ffmpeg.binary", "CORE_FFMPEG_BINARY", nil, "Path to ffmpeg binary", true, false)
	d.vars.Register(value.NewInt64(&d.FFmpeg.MaxProcesses, 0), "ffmpeg.max_processes", "CORE_FFMPEG_MAXPROCESSES", nil, "Max. allowed simultaneously running ffmpeg instances, 0 for unlimited", false, false)
	d.vars.Register(value.NewStringList(&d.FFmpeg.Access.Input.Allow, []string{}, " "), "ffmpeg.access.input.allow", "CORE_FFMPEG_ACCESS_INPUT_ALLOW", nil, "List of allowed expression to match against the input addresses", false, false)
	d.vars.Register(value.NewStringList(&d.FFmpeg.Access.Input.Block, []string{}, " "), "ffmpeg.access.input.block", "CORE_FFMPEG_ACCESS_INPUT_BLOCK", nil, "List of blocked expression to match against the input addresses", false, false)
	d.vars.Register(value.NewStringList(&d.FFmpeg.Access.Output.Allow, []string{}, " "), "ffmpeg.access.output.allow", "CORE_FFMPEG_ACCESS_OUTPUT_ALLOW", nil, "List of allowed expression to match against the output addresses", false, false)
	d.vars.Register(value.NewStringList(&d.FFmpeg.Access.Output.Block, []string{}, " "), "ffmpeg.access.output.block", "CORE_FFMPEG_ACCESS_OUTPUT_BLOCK", nil, "List of blocked expression to match against the output addresses", false, false)
	d.vars.Register(value.NewInt(&d.FFmpeg.Log.MaxLines, 50), "ffmpeg.log.max_lines", "CORE_FFMPEG_LOG_MAX_LINES", []string{"CORE_FFMPEG_LOG_MAXLINES"}, "Number of latest log lines to keep for each process", false, false)
	d.vars.Register(value.NewInt(&d.FFmpeg.Log.MaxHistory, 3), "ffmpeg.log.max_history", "CORE_FFMPEG_LOG_MAX_HISTORY", []string{"CORE_FFMPEG_LOG_MAXHISTORY"}, "Number of latest logs to keep for each process", false, false)

	// Playout
	d.vars.Register(value.NewBool(&d.Playout.Enable, false), "playout.enable", "CORE_PLAYOUT_ENABLE", nil, "Enable playout proxy where available", false, false)
	d.vars.Register(value.NewPort(&d.Playout.MinPort, 0), "playout.min_port", "CORE_PLAYOUT_MIN_PORT", []string{"CORE_PLAYOUT_MINPORT"}, "Min. playout server port", false, false)
	d.vars.Register(value.NewPort(&d.Playout.MaxPort, 0), "playout.max_port", "CORE_PLAYOUT_MAX_PORT", []string{"CORE_PLAYOUT_MAXPORT"}, "Max. playout server port", false, false)

	// Debug
	d.vars.Register(value.NewBool(&d.Debug.Profiling, false), "debug.profiling", "CORE_DEBUG_PROFILING", nil, "Enable profiling endpoint on /profiling", false, false)
	d.vars.Register(value.NewInt(&d.Debug.ForceGC, 0), "debug.force_gc", "CORE_DEBUG_FORCE_GC", []string{"CORE_DEBUG_FORCEGC"}, "Number of seconds between forcing GC to return memory to the OS", false, false)
	d.vars.Register(value.NewInt64(&d.Debug.MemoryLimit, 0), "debug.memory_limit_mbytes", "CORE_DEBUG_MEMORY_LIMIT_MBYTES", nil, "Impose a soft memory limit for the core, in megabytes", false, false)

	// Metrics
	d.vars.Register(value.NewBool(&d.Metrics.Enable, false), "metrics.enable", "CORE_METRICS_ENABLE", nil, "Enable collecting historic metrics data", false, false)
	d.vars.Register(value.NewBool(&d.Metrics.EnablePrometheus, false), "metrics.enable_prometheus", "CORE_METRICS_ENABLE_PROMETHEUS", nil, "Enable prometheus endpoint /metrics", false, false)
	d.vars.Register(value.NewInt64(&d.Metrics.Range, 300), "metrics.range_seconds", "CORE_METRICS_RANGE_SECONDS", nil, "Seconds to keep history data", false, false)
	d.vars.Register(value.NewInt64(&d.Metrics.Interval, 2), "metrics.interval_seconds", "CORE_METRICS_INTERVAL_SECONDS", nil, "Interval for collecting metrics", false, false)

	// Sessions
	d.vars.Register(value.NewBool(&d.Sessions.Enable, true), "sessions.enable", "CORE_SESSIONS_ENABLE", nil, "Enable collecting HLS session stats for /memfs", false, false)
	d.vars.Register(value.NewCIDRList(&d.Sessions.IPIgnoreList, []string{"127.0.0.1/32", "::1/128"}, ","), "sessions.ip_ignorelist", "CORE_SESSIONS_IP_IGNORELIST", nil, "List of IP ranges in CIDR notation to ignore", false, false)
	d.vars.Register(value.NewInt(&d.Sessions.SessionTimeout, 30), "sessions.session_timeout_sec", "CORE_SESSIONS_SESSION_TIMEOUT_SEC", nil, "Timeout for an idle session", false, false)
	d.vars.Register(value.NewBool(&d.Sessions.Persist, false), "sessions.persist", "CORE_SESSIONS_PERSIST", nil, "Whether to persist session history. Will be stored as sessions.json in db.dir", false, false)
	d.vars.Register(value.NewInt(&d.Sessions.PersistInterval, 300), "sessions.persist_interval_sec", "CORE_SESSIONS_PERSIST_INTERVAL_SEC", nil, "Interval in seconds in which to persist the current session history", false, false)
	d.vars.Register(value.NewUint64(&d.Sessions.MaxBitrate, 0), "sessions.max_bitrate_mbit", "CORE_SESSIONS_MAXBITRATE_MBIT", nil, "Max. allowed outgoing bitrate in mbit/s, 0 for unlimited", false, false)
	d.vars.Register(value.NewUint64(&d.Sessions.MaxSessions, 0), "sessions.max_sessions", "CORE_SESSIONS_MAX_SESSIONS", []string{"CORE_SESSIONS_MAXSESSIONS"}, "Max. allowed number of simultaneous sessions, 0 for unlimited", false, false)

	// Service
	d.vars.Register(value.NewBool(&d.Service.Enable, false), "service.enable", "CORE_SERVICE_ENABLE", nil, "Enable connecting to the Restreamer Service", false, false)
	d.vars.Register(value.NewString(&d.Service.Token, ""), "service.token", "CORE_SERVICE_TOKEN", nil, "Restreamer Service account token", false, true)
	d.vars.Register(value.NewURL(&d.Service.URL, "https://service.datarhei.com"), "service.url", "CORE_SERVICE_URL", nil, "URL of the Restreamer Service", false, false)

	// Router
	d.vars.Register(value.NewStringList(&d.Router.BlockedPrefixes, []string{"/api"}, ","), "router.blocked_prefixes", "CORE_ROUTER_BLOCKED_PREFIXES", nil, "List of path prefixes that can't be routed", false, false)
	d.vars.Register(value.NewStringMapString(&d.Router.Routes, nil), "router.routes", "CORE_ROUTER_ROUTES", nil, "List of route mappings", false, false)
	d.vars.Register(value.NewDir(&d.Router.UIPath, "", d.fs), "router.ui_path", "CORE_ROUTER_UI_PATH", nil, "Path to a directory holding UI files mounted as /ui", false, false)
}

// Validate validates the current state of the Config for completeness and sanity. Errors are
// written to the log. Use resetLogs to indicate to reset the logs prior validation.
func (d *Config) Validate(resetLogs bool) {
	if resetLogs {
		d.vars.ResetLogs()
	}

	if d.Version != version {
		d.vars.Log("error", "version", "unknown configuration layout version (found version %d, expecting version %d)", d.Version, version)

		return
	}

	d.vars.Validate()

	// Individual sanity checks

	// If HTTP Auth is enabled, check that the username and password are set
	if d.API.Auth.Enable {
		if len(d.API.Auth.Username) == 0 || len(d.API.Auth.Password) == 0 {
			d.vars.Log("error", "api.auth.enable", "api.auth.username and api.auth.password must be set")
		}
	}

	// If Auth0 is enabled, check that domain, audience, and clientid are set
	if d.API.Auth.Auth0.Enable {
		if len(d.API.Auth.Auth0.Tenants) == 0 {
			d.vars.Log("error", "api.auth.auth0.enable", "at least one tenants must be set")
		}

		for i, t := range d.API.Auth.Auth0.Tenants {
			if len(t.Domain) == 0 || len(t.Audience) == 0 || len(t.ClientID) == 0 {
				d.vars.Log("error", "api.auth.auth0.tenants", "domain, audience, and clientid must be set (tenant %d)", i)
			}
		}
	}

	// If TLS is enabled and Let's Encrypt is disabled, require certfile and keyfile
	if d.TLS.Enable && !d.TLS.Auto {
		if len(d.TLS.CertFile) == 0 || len(d.TLS.KeyFile) == 0 {
			d.vars.Log("error", "tls.enable", "tls.certfile and tls.keyfile must be set")
		}
	}

	// If TLS and Let's Encrypt certificate is enabled, we require a public hostname
	if d.TLS.Enable && d.TLS.Auto {
		if len(d.Host.Name) == 0 {
			d.vars.Log("error", "host.name", "a hostname must be set in order to get an automatic TLS certificate")
		} else {
			r := &net.Resolver{
				PreferGo:     true,
				StrictErrors: true,
			}

			for _, host := range d.Host.Name {
				// Don't lookup IP addresses
				if ip := net.ParseIP(host); ip != nil {
					d.vars.Log("error", "host.name", "only host names are allowed if automatic TLS is enabled, but found IP address: %s", host)
				}

				// Lookup host name with a timeout
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

				_, err := r.LookupHost(ctx, host)
				if err != nil {
					d.vars.Log("error", "host.name", "the host '%s' can't be resolved and will not work with automatic TLS", host)
				}

				cancel()
			}
		}
	}

	// If TLS and Let's Encrypt certificate is enabled, we require a non-empty email address
	if d.TLS.Enable && d.TLS.Auto {
		if len(d.TLS.Email) == 0 {
			d.vars.SetDefault("tls.email")
		}
	}

	// If TLS for RTMP is enabled, TLS must be enabled
	if d.RTMP.EnableTLS {
		if !d.RTMP.Enable {
			d.vars.Log("error", "rtmp.enable", "RTMP server must be enabled if RTMPS server is enabled")
		}

		if !d.TLS.Enable {
			d.vars.Log("error", "rtmp.enable_tls", "RTMPS server can only be enabled if TLS is enabled")
		}

		if len(d.RTMP.AddressTLS) == 0 {
			d.vars.Log("error", "rtmp.address_tls", "RTMPS server address must be set")
		}

		if d.RTMP.Enable && d.RTMP.Address == d.RTMP.AddressTLS {
			d.vars.Log("error", "rtmp.address", "The RTMP and RTMPS server can't listen on the same address")
		}
	}

	// If CORE_MEMFS_USERNAME and CORE_MEMFS_PASSWORD are set, automatically active/deactivate Basic-Auth for memfs
	if d.vars.IsMerged("storage.memory.auth.username") && d.vars.IsMerged("storage.memory.auth.password") {
		d.Storage.Memory.Auth.Enable = true

		if len(d.Storage.Memory.Auth.Username) == 0 && len(d.Storage.Memory.Auth.Password) == 0 {
			d.Storage.Memory.Auth.Enable = false
		}
	}

	// If Basic-Auth for memfs is enable, check that the username and password are set
	if d.Storage.Memory.Auth.Enable {
		if len(d.Storage.Memory.Auth.Username) == 0 || len(d.Storage.Memory.Auth.Password) == 0 {
			d.vars.Log("error", "storage.memory.auth.enable", "storage.memory.auth.username and storage.memory.auth.password must be set")
		}
	}

	if len(d.Storage.S3) != 0 {
		names := map[string]struct{}{
			"disk": {},
			"mem":  {},
		}

		for _, s3 := range d.Storage.S3 {
			if _, ok := names[s3.Name]; ok {
				d.vars.Log("error", "storage.s3", "the name %s is already in use or reserved", s3.Name)
			}

			names[s3.Name] = struct{}{}
		}
	}

	// If playout is enabled, check that the port range is sane
	if d.Playout.Enable {
		if d.Playout.MinPort >= d.Playout.MaxPort {
			d.vars.Log("error", "playout.min_port", "must be bigger than playout.max_port")
		}
	}

	// If cache is enabled, a valid TTL has to be set to a useful value
	if d.Storage.Disk.Cache.Enable && d.Storage.Disk.Cache.TTL < 0 {
		d.vars.Log("error", "storage.disk.cache.ttl_seconds", "must be equal or greater than 0")
	}

	// If the stats are enabled, the session timeout has to be set to a useful value
	if d.Sessions.Enable && d.Sessions.SessionTimeout < 1 {
		d.vars.Log("error", "stats.session_timeout_sec", "must be equal or greater than 1")
	}

	// If the stats and their persistence are enabled, the persist interval has to be set to a useful value
	if d.Sessions.Enable && d.Sessions.PersistInterval < 0 {
		d.vars.Log("error", "stats.persist_interval_sec", "must be at equal or greater than 0")
	}

	// If the service is enabled, the token and enpoint have to be defined
	if d.Service.Enable {
		if len(d.Service.Token) == 0 {
			d.vars.Log("error", "service.token", "must be non-empty")
		}

		if len(d.Service.URL) == 0 {
			d.vars.Log("error", "service.url", "must be non-empty")
		}
	}

	// If historic metrics are enabled, the timerange and interval have to be valid
	if d.Metrics.Enable {
		if d.Metrics.Range <= 0 {
			d.vars.Log("error", "metrics.range", "must be greater 0")
		}

		if d.Metrics.Interval <= 0 {
			d.vars.Log("error", "metrics.interval", "must be greater 0")
		}

		if d.Metrics.Interval > d.Metrics.Range {
			d.vars.Log("error", "metrics.interval", "must be smaller than the range")
		}
	}
}

// Merge merges the values of the known environment variables into the configuration
func (d *Config) Merge() {
	d.vars.Merge()
}

// Messages calls for each log entry the provided callback. The level has the values 'error', 'warn', or 'info'.
// The name is the name of the configuration value, e.g. 'api.auth.enable'. The message is the log message.
func (d *Config) Messages(logger func(level string, v vars.Variable, message string)) {
	d.vars.Messages(logger)
}

// HasErrors returns whether there are some error messages in the log.
func (d *Config) HasErrors() bool {
	return d.vars.HasErrors()
}

// Overrides returns a list of configuration value names that have been overriden by an environment variable.
func (d *Config) Overrides() []string {
	return d.vars.Overrides()
}
