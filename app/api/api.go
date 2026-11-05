package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	golog "log"
	"math"
	gonet "net"
	gohttp "net/http"
	"net/url"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/datarhei/core/v16/app"
	"github.com/datarhei/core/v16/config"
	configstore "github.com/datarhei/core/v16/config/store"
	configvars "github.com/datarhei/core/v16/config/vars"
	"github.com/datarhei/core/v16/ffmpeg"
	"github.com/datarhei/core/v16/http"
	"github.com/datarhei/core/v16/http/cache"
	httpfs "github.com/datarhei/core/v16/http/fs"
	"github.com/datarhei/core/v16/http/jwt"
	"github.com/datarhei/core/v16/http/router"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/math/rand"
	"github.com/datarhei/core/v16/monitor"
	"github.com/datarhei/core/v16/net"
	"github.com/datarhei/core/v16/prometheus"
	"github.com/datarhei/core/v16/restream"
	restreamapp "github.com/datarhei/core/v16/restream/app"
	"github.com/datarhei/core/v16/restream/replace"
	restreamstore "github.com/datarhei/core/v16/restream/store"
	"github.com/datarhei/core/v16/rtmp"
	"github.com/datarhei/core/v16/service"
	"github.com/datarhei/core/v16/session"
	"github.com/datarhei/core/v16/srt"
	"github.com/datarhei/core/v16/update"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
)

// The API interface is the implementation for the restreamer API.
type API interface {
	// Start starts the API. This is blocking until the app has
	// been ended with Stop() or Destroy(). In this case a nil error
	// is returned. An ErrConfigReload error is returned if a
	// configuration reload has been requested.
	Start() error

	// Stop stops the API, some states may be kept intact such
	// that they can be reused after starting the API again.
	Stop()

	// Destroy is the same as Stop() but no state will be kept intact.
	Destroy()

	// Reload the configuration for the API. If there's an error the
	// previously loaded configuration is not altered.
	Reload() error
}

type api struct {
	restream      restream.Restreamer
	ffmpeg        ffmpeg.FFmpeg
	diskfs        fs.Filesystem
	memfs         fs.Filesystem
	s3fs          map[string]fs.Filesystem
	rtmpserver    rtmp.Server
	srtserver     srt.Server
	metrics       monitor.HistoryMonitor
	prom          prometheus.Metrics
	service       service.Service
	sessions      session.Registry
	cache         cache.Cacher
	mainserver    *gohttp.Server
	sidecarserver *gohttp.Server
	httpjwt       jwt.JWT
	update        update.Checker
	replacer      replace.Replacer

	errorChan chan error

	gcTickerStop context.CancelFunc

	log struct {
		writer io.Writer
		buffer log.BufferWriter
		logger struct {
			core    log.Logger
			main    log.Logger
			sidecar log.Logger
			rtmp    log.Logger
			rtmps   log.Logger
			srt     log.Logger
		}
	}

	config struct {
		path   string
		store  configstore.Store
		config *config.Config
	}

	lock   sync.Mutex
	wgStop sync.WaitGroup
	state  string
}

// ErrConfigReload is an error returned to indicate that a reload of
// the configuration has been requested.
var ErrConfigReload = fmt.Errorf("configuration reload")

// New returns a new instance of the API interface
func New(configpath string, logwriter io.Writer) (API, error) {
	a := &api{
		state: "idle",
		s3fs:  map[string]fs.Filesystem{},
	}

	a.config.path = configpath
	a.log.writer = logwriter

	if a.log.writer == nil {
		a.log.writer = io.Discard
	}

	a.errorChan = make(chan error, 1)

	if err := a.Reload(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *api) Reload() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.state == "running" {
		return fmt.Errorf("can't reload config while running")
	}

	if a.errorChan == nil {
		a.errorChan = make(chan error, 1)
	}

	logger := log.New("Core").WithOutput(log.NewConsoleWriter(a.log.writer, log.Lwarn, true))

	rootfs, _ := fs.NewDiskFilesystem(fs.DiskConfig{})
	store, err := configstore.NewJSON(rootfs, a.config.path, func() {
		a.errorChan <- ErrConfigReload
	})
	if err != nil {
		return err
	}

	cfg := store.Get()

	cfg.Merge()

	if len(cfg.Host.Name) == 0 && cfg.Host.Auto {
		cfg.Host.Name = net.GetPublicIPs(5 * time.Second)
	}

	cfg.Validate(false)

	loglevel := log.Linfo

	switch cfg.Log.Level {
	case "silent":
		loglevel = log.Lsilent
	case "error":
		loglevel = log.Lerror
	case "warn":
		loglevel = log.Lwarn
	case "info":
		loglevel = log.Linfo
	case "debug":
		loglevel = log.Ldebug
	default:
		break
	}

	buffer := log.NewBufferWriter(loglevel, cfg.Log.MaxLines)

	logger = logger.WithOutput(log.NewLevelRewriter(
		log.NewMultiWriter(
			log.NewTopicWriter(
				log.NewConsoleWriter(a.log.writer, loglevel, true),
				cfg.Log.Topics,
			),
			buffer,
		),
		[]log.LevelRewriteRule{
			// FFmpeg annoyance, move all warnings about unathorized access to memfs from ffmpeg to debug level
			// ts=2022-04-28T07:24:27Z level=WARN component="HTTP" address=":8080" client="::1" latency_ms=0 method="PUT" path="/memfs/00a10a69-416a-4cd5-9d4f-6d88ed3dd7f5_0917.ts" proto="HTTP/1.1" size_bytes=65 status=401 status_text="Unauthorized" user_agent="Lavf/58.76.100"
			{
				Level:     log.Ldebug,
				Component: "HTTP",
				Match: map[string]string{
					"client":      "^(::1|127.0.0.1)$",
					"method":      "^(PUT|POST|DELETE)$",
					"status_text": "^Unauthorized$",
					"user_agent":  "^Lavf/",
				},
			},
		},
	))

	logfields := log.Fields{
		"application": app.Name,
		"version":     app.Version.String(),
		"repository":  "https://github.com/datarhei/core",
		"license":     "Apache License Version 2.0",
		"arch":        app.Arch,
		"compiler":    app.Compiler,
	}

	if len(app.Commit) != 0 && len(app.Branch) != 0 {
		logfields["commit"] = app.Commit
		logfields["branch"] = app.Branch
	}

	if len(app.Build) != 0 {
		logfields["build"] = app.Build
	}

	logger.Info().WithFields(logfields).Log("")

	logger.Info().WithField("path", a.config.path).Log("Read config file")

	configlogger := logger.WithComponent("Config")
	cfg.Messages(func(level string, v configvars.Variable, message string) {
		configlogger = configlogger.WithFields(log.Fields{
			"variable":    v.Name,
			"value":       v.Value,
			"env":         v.EnvName,
			"description": v.Description,
			"override":    v.Merged,
		})
		configlogger.Debug().Log(message)

		switch level {
		case "warn":
			configlogger.Warn().Log(message)
		case "error":
			configlogger.Error().WithField("error", message).Log("")
		default:
			break
		}
	})

	if cfg.HasErrors() {
		logger.Error().WithField("error", "Not all variables are set or are valid. Check the error messages above. Bailing out.").Log("")
		return fmt.Errorf("not all variables are set or valid")
	}

	cfg.LoadedAt = time.Now()

	store.SetActive(cfg)

	a.config.store = store
	a.config.config = cfg
	a.log.logger.core = logger
	a.log.buffer = buffer

	return nil
}

func (a *api) start() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.errorChan == nil {
		a.errorChan = make(chan error, 1)
	}

	if a.state == "running" {
		return fmt.Errorf("already running")
	}

	a.state = "starting"

	cfg := a.config.store.GetActive()

	if cfg.Sessions.Enable {
		sessionConfig := session.Config{
			Logger: a.log.logger.core.WithComponent("Session"),
		}

		if cfg.Sessions.Persist {
			fs, err := fs.NewRootedDiskFilesystem(fs.RootedDiskConfig{
				Root: filepath.Join(cfg.DB.Dir, "sessions"),
			})
			if err != nil {
				return fmt.Errorf("unable to create filesystem for persisting sessions: %w", err)
			}
			sessionConfig.PersistFS = fs
		}

		sessions, err := session.New(sessionConfig)
		if err != nil {
			return fmt.Errorf("unable to create sessioner: %w", err)
		}

		a.sessions = sessions

		iplimiter, err := net.NewIPLimiter(cfg.Sessions.IPIgnoreList, nil)
		if err != nil {
			return fmt.Errorf("incorret IP ranges for the statistics provided: %w", err)
		}

		config := session.CollectorConfig{
			MaxTxBitrate:    cfg.Sessions.MaxBitrate * 1024 * 1024,
			MaxSessions:     cfg.Sessions.MaxSessions,
			InactiveTimeout: 5 * time.Second,
			SessionTimeout:  time.Duration(cfg.Sessions.SessionTimeout) * time.Second,
			PersistInterval: time.Duration(cfg.Sessions.PersistInterval) * time.Second,
			Limiter:         iplimiter,
		}

		hls, err := sessions.Register("hls", config)
		if err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		_, err = sessions.Register("hlsingress", session.CollectorConfig{
			SessionTimeout: 30 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		rtmp, err := sessions.Register("rtmp", config)
		if err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		srt, err := sessions.Register("srt", config)
		if err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		if _, err := sessions.Register("http", config); err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		if _, err := sessions.Register("https", config); err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		ffmpeg, err := sessions.Register("ffmpeg", config)
		if err != nil {
			return fmt.Errorf("unable to register session collector: %w", err)
		}

		hls.AddCompanion(rtmp)
		hls.AddCompanion(srt)
		hls.AddCompanion(ffmpeg)

		rtmp.AddCompanion(hls)
		rtmp.AddCompanion(ffmpeg)
		rtmp.AddCompanion(srt)

		srt.AddCompanion(hls)
		srt.AddCompanion(ffmpeg)
		srt.AddCompanion(rtmp)

		ffmpeg.AddCompanion(hls)
		ffmpeg.AddCompanion(rtmp)
		ffmpeg.AddCompanion(srt)
	} else {
		sessions, _ := session.New(session.Config{})
		a.sessions = sessions
	}

	diskfs, err := fs.NewRootedDiskFilesystem(fs.RootedDiskConfig{
		Root:   cfg.Storage.Disk.Dir,
		Logger: a.log.logger.core.WithComponent("DiskFS"),
	})
	if err != nil {
		return fmt.Errorf("disk filesystem: %w", err)
	}

	if diskfsRoot, err := filepath.Abs(cfg.Storage.Disk.Dir); err != nil {
		return err
	} else {
		diskfs.SetMetadata("base", diskfsRoot)
	}

	a.diskfs = diskfs

	baseMemFS := url.URL{
		Scheme: "http",
		Path:   "/memfs",
	}

	host, port, _ := gonet.SplitHostPort(cfg.Address)
	if len(host) == 0 {
		baseMemFS.Host = "localhost:" + port
	} else {
		baseMemFS.Host = cfg.Address
	}

	if cfg.Storage.Memory.Auth.Enable {
		baseMemFS.User = url.UserPassword(cfg.Storage.Memory.Auth.Username, cfg.Storage.Memory.Auth.Password)
	}

	if a.memfs == nil {
		memfs, _ := fs.NewMemFilesystem(fs.MemConfig{
			Logger: a.log.logger.core.WithComponent("MemFS"),
		})

		memfs.SetMetadata("base", baseMemFS.String())

		sizedfs, _ := fs.NewSizedFilesystem(memfs, cfg.Storage.Memory.Size*1024*1024, cfg.Storage.Memory.Purge)

		a.memfs = sizedfs
	} else {
		a.memfs.SetMetadata("base", baseMemFS.String())
		if sizedfs, ok := a.memfs.(fs.SizedFilesystem); ok {
			sizedfs.Resize(cfg.Storage.Memory.Size*1024*1024, cfg.Storage.Memory.Purge)
		}
	}

	for _, s3 := range cfg.Storage.S3 {
		if _, ok := a.s3fs[s3.Name]; ok {
			return fmt.Errorf("the name '%s' for a s3 filesystem is already in use", s3.Name)
		}

		baseS3FS := url.URL{
			Scheme: "http",
			Path:   s3.Mountpoint,
		}

		host, port, _ := gonet.SplitHostPort(cfg.Address)
		if len(host) == 0 {
			baseS3FS.Host = "localhost:" + port
		} else {
			baseS3FS.Host = cfg.Address
		}

		if s3.Auth.Enable {
			baseS3FS.User = url.UserPassword(s3.Auth.Username, s3.Auth.Password)
		}

		s3fs, err := fs.NewS3Filesystem(fs.S3Config{
			Name:            s3.Name,
			Endpoint:        s3.Endpoint,
			AccessKeyID:     s3.AccessKeyID,
			SecretAccessKey: s3.SecretAccessKey,
			Region:          s3.Region,
			Bucket:          s3.Bucket,
			UseSSL:          s3.UseSSL,
			Logger:          a.log.logger.core.WithComponent("FS"),
		})
		if err != nil {
			return fmt.Errorf("s3 filesystem (%s): %w", s3.Name, err)
		}

		s3fs.SetMetadata("base", baseS3FS.String())

		a.s3fs[s3.Name] = s3fs
	}

	var portrange net.Portranger

	if cfg.Playout.Enable {
		portrange, err = net.NewPortrange(cfg.Playout.MinPort, cfg.Playout.MaxPort)
		if err != nil {
			return fmt.Errorf("playout port range: %w", err)
		}
	}

	validatorIn, err := ffmpeg.NewValidator(cfg.FFmpeg.Access.Input.Allow, cfg.FFmpeg.Access.Input.Block)
	if err != nil {
		return fmt.Errorf("input address validator: %w", err)
	}

	validatorOut, err := ffmpeg.NewValidator(cfg.FFmpeg.Access.Output.Allow, cfg.FFmpeg.Access.Output.Block)
	if err != nil {
		return fmt.Errorf("output address validator: %w", err)
	}

	ffmpeg, err := ffmpeg.New(ffmpeg.Config{
		Binary:           cfg.FFmpeg.Binary,
		MaxProc:          cfg.FFmpeg.MaxProcesses,
		MaxLogLines:      cfg.FFmpeg.Log.MaxLines,
		LogHistoryLength: cfg.FFmpeg.Log.MaxHistory,
		ValidatorInput:   validatorIn,
		ValidatorOutput:  validatorOut,
		Portrange:        portrange,
		Collector:        a.sessions.Collector("ffmpeg"),
	})
	if err != nil {
		return fmt.Errorf("unable to create ffmpeg: %w", err)
	}

	a.ffmpeg = ffmpeg

	a.replacer = replace.New()

	{
		a.replacer.RegisterTemplateFunc("diskfs", func(config *restreamapp.Config, section string) string {
			return a.diskfs.Metadata("base")
		}, nil)

		a.replacer.RegisterTemplateFunc("fs:disk", func(config *restreamapp.Config, section string) string {
			return a.diskfs.Metadata("base")
		}, nil)

		a.replacer.RegisterTemplateFunc("memfs", func(config *restreamapp.Config, section string) string {
			return a.memfs.Metadata("base")
		}, nil)

		a.replacer.RegisterTemplateFunc("fs:mem", func(config *restreamapp.Config, section string) string {
			return a.memfs.Metadata("base")
		}, nil)

		for name, s3 := range a.s3fs {
			a.replacer.RegisterTemplate("fs:"+name, s3.Metadata("base"), nil)
		}

		a.replacer.RegisterTemplateFunc("rtmp", func(config *restreamapp.Config, section string) string {
			host, port, _ := gonet.SplitHostPort(cfg.RTMP.Address)
			if len(host) == 0 {
				host = "localhost"
			}

			template := "rtmp://" + host + ":" + port
			if cfg.RTMP.App != "/" {
				template += cfg.RTMP.App
			}
			template += "/{name}"

			if len(cfg.RTMP.Token) != 0 {
				template += "?token=" + cfg.RTMP.Token
			}

			return template
		}, nil)

		a.replacer.RegisterTemplateFunc("srt", func(config *restreamapp.Config, section string) string {
			host, port, _ = gonet.SplitHostPort(cfg.SRT.Address)
			if len(host) == 0 {
				host = "localhost"
			}

			template := "srt://" + host + ":" + port + "?mode=caller&transtype=live&latency={latency}&streamid={name}"
			if section == "output" {
				template += ",mode:publish"
			} else {
				template += ",mode:request"
			}
			if len(cfg.SRT.Token) != 0 {
				template += ",token:" + cfg.SRT.Token
			}
			if len(cfg.SRT.Passphrase) != 0 {
				template += "&passphrase=" + cfg.SRT.Passphrase
			}

			return template
		}, map[string]string{
			"latency": "20000", // 20 milliseconds, FFmpeg requires microseconds
		})
	}

	filesystems := []fs.Filesystem{
		a.diskfs,
		a.memfs,
	}

	for _, fs := range a.s3fs {
		filesystems = append(filesystems, fs)
	}

	var store restreamstore.Store = nil

	{
		fs, err := fs.NewRootedDiskFilesystem(fs.RootedDiskConfig{
			Root: cfg.DB.Dir,
		})
		if err != nil {
			return err
		}
		store, err = restreamstore.NewJSON(restreamstore.JSONConfig{
			Filesystem: fs,
			Filepath:   "/db.json",
			Logger:     a.log.logger.core.WithComponent("ProcessStore"),
		})
		if err != nil {
			return err
		}
	}

	restream, err := restream.New(restream.Config{
		ID:           cfg.ID,
		Name:         cfg.Name,
		Store:        store,
		Filesystems:  filesystems,
		Replace:      a.replacer,
		FFmpeg:       a.ffmpeg,
		MaxProcesses: cfg.FFmpeg.MaxProcesses,
		Logger:       a.log.logger.core.WithComponent("Process"),
	})

	if err != nil {
		return fmt.Errorf("unable to create restreamer: %w", err)
	}

	a.restream = restream

	var httpjwt jwt.JWT

	if cfg.API.Auth.Enable {
		secret := rand.String(32)
		if len(cfg.API.Auth.JWT.Secret) != 0 {
			secret = cfg.API.Auth.Username + cfg.API.Auth.Password + cfg.API.Auth.JWT.Secret
		}

		var err error
		httpjwt, err = jwt.New(jwt.Config{
			Realm:         app.Name,
			Secret:        secret,
			SkipLocalhost: cfg.API.Auth.DisableLocalhost,
		})

		if err != nil {
			return fmt.Errorf("unable to create JWT provider: %w", err)
		}

		if validator, err := jwt.NewLocalValidator(cfg.API.Auth.Username, cfg.API.Auth.Password); err == nil {
			if err := httpjwt.AddValidator(app.Name, validator); err != nil {
				return fmt.Errorf("unable to add local JWT validator: %w", err)
			}
		} else {
			return fmt.Errorf("unable to create local JWT validator: %w", err)
		}

		if cfg.API.Auth.Auth0.Enable {
			for _, t := range cfg.API.Auth.Auth0.Tenants {
				if validator, err := jwt.NewAuth0Validator(t.Domain, t.Audience, t.ClientID, t.Users); err == nil {
					if err := httpjwt.AddValidator("https://"+t.Domain+"/", validator); err != nil {
						return fmt.Errorf("unable to add Auth0 JWT validator: %w", err)
					}
				} else {
					return fmt.Errorf("unable to create Auth0 JWT validator: %w", err)
				}
			}
		}
	}

	a.httpjwt = httpjwt

	metrics, err := monitor.NewHistory(monitor.HistoryConfig{
		Enable:    cfg.Metrics.Enable,
		Timerange: time.Duration(cfg.Metrics.Range) * time.Second,
		Interval:  time.Duration(cfg.Metrics.Interval) * time.Second,
	})
	if err != nil {
		return fmt.Errorf("unable to create metrics collector with history: %w", err)
	}

	metrics.Register(monitor.NewUptimeCollector())
	metrics.Register(monitor.NewCPUCollector())
	metrics.Register(monitor.NewMemCollector())
	metrics.Register(monitor.NewNetCollector())
	metrics.Register(monitor.NewDiskCollector(a.diskfs.Metadata("base")))
	metrics.Register(monitor.NewFilesystemCollector("diskfs", a.diskfs))
	metrics.Register(monitor.NewFilesystemCollector("memfs", a.memfs))
	for name, fs := range a.s3fs {
		metrics.Register(monitor.NewFilesystemCollector(name, fs))
	}
	metrics.Register(monitor.NewRestreamCollector(a.restream))
	metrics.Register(monitor.NewFFmpegCollector(a.ffmpeg))
	metrics.Register(monitor.NewSessionCollector(a.sessions, []string{}))

	a.metrics = metrics

	if cfg.Metrics.EnablePrometheus {
		prom := prometheus.New()

		prom.Register(prometheus.NewUptimeCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewCPUCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewMemCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewNetCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewDiskCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewFilesystemCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewRestreamCollector(cfg.ID, metrics))
		prom.Register(prometheus.NewSessionCollector(cfg.ID, metrics))

		a.prom = prom
	}

	if cfg.Service.Enable {
		address := cfg.Address
		domain := "http://"
		if cfg.TLS.Enable {
			address = cfg.TLS.Address
			domain = "https://"
		}

		host, port, _ := gonet.SplitHostPort(address)
		if len(host) == 0 {
			host = "localhost"
		}

		if len(cfg.Host.Name) != 0 {
			host = cfg.Host.Name[0]
		}

		domain += host + ":" + port

		s, err := service.New(service.Config{
			ID:      cfg.ID,
			Version: fmt.Sprintf("%s %s (%s)", app.Name, app.Version, app.Arch),
			Domain:  domain,
			Token:   cfg.Service.Token,
			URL:     cfg.Service.URL,
			Monitor: a.metrics,
			Logger:  a.log.logger.core.WithComponent("Service"),
		})
		if err != nil {
			return fmt.Errorf("unable to connect to service: %w", err)
		}

		a.service = s
	}

	if cfg.CheckForUpdates {
		s, err := update.New(update.Config{
			ID:      cfg.ID,
			Name:    app.Name,
			Version: app.Version.String(),
			Arch:    app.Arch,
			Monitor: a.metrics,
			Logger:  a.log.logger.core.WithComponent("Update"),
		})

		if err != nil {
			return fmt.Errorf("unable to create update checker: %w", err)
		}

		a.update = s
	}

	if cfg.Storage.Disk.Cache.Enable {
		cache, err := cache.NewLRUCache(cache.LRUConfig{
			TTL:             time.Duration(cfg.Storage.Disk.Cache.TTL) * time.Second,
			MaxSize:         cfg.Storage.Disk.Cache.Size * 1024 * 1024,
			MaxFileSize:     cfg.Storage.Disk.Cache.FileSize * 1024 * 1024,
			AllowExtensions: cfg.Storage.Disk.Cache.Types.Allow,
			BlockExtensions: cfg.Storage.Disk.Cache.Types.Block,
			Logger:          a.log.logger.core.WithComponent("HTTPCache"),
		})

		if err != nil {
			return fmt.Errorf("unable to create cache: %w", err)
		}

		a.cache = cache
	}

	var autocertManager *certmagic.Config

	if cfg.TLS.Enable {
		if cfg.TLS.Auto {
			if len(cfg.Host.Name) == 0 {
				return fmt.Errorf("at least one host must be provided in host.name or CORE_HOST_NAME")
			}

			certmagic.Default.Storage = &certmagic.FileStorage{
				Path: cfg.DB.Dir + "/cert",
			}
			certmagic.Default.DefaultServerName = cfg.Host.Name[0]
			certmagic.Default.Logger = zap.NewNop()

			certmagic.DefaultACME.Agreed = true
			certmagic.DefaultACME.Email = cfg.TLS.Email
			certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
			certmagic.DefaultACME.DisableHTTPChallenge = false
			certmagic.DefaultACME.DisableTLSALPNChallenge = true
			certmagic.DefaultACME.Logger = zap.NewNop()

			magic := certmagic.NewDefault()
			acme := certmagic.NewACMEIssuer(magic, certmagic.DefaultACME)
			acme.Logger = zap.NewNop()

			magic.Issuers = []certmagic.Issuer{acme}
			magic.Logger = zap.NewNop()

			autocertManager = magic

			// Start temporary http server on configured port
			tempserver := &gohttp.Server{
				Addr: cfg.Address,
				Handler: acme.HTTPChallengeHandler(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
					w.WriteHeader(gohttp.StatusNotFound)
				})),
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
				tempserver.ListenAndServe()
				wg.Done()
			}()

			var certerror bool

			// For each domain, get the certificate
			for _, host := range cfg.Host.Name {
				logger := a.log.logger.core.WithComponent("Let's Encrypt").WithField("host", host)
				logger.Info().Log("Acquiring certificate ...")

				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))

				err := autocertManager.ManageSync(ctx, []string{host})

				cancel()

				if err != nil {
					logger.Error().WithField("error", err).Log("Failed to acquire certificate")
					certerror = true
					/*
						problems, err := letsdebug.Check(host, letsdebug.HTTP01)
						if err != nil {
							logger.Error().WithField("error", err).Log("Failed to debug certificate acquisition")
						}

						for _, p := range problems {
							logger.Error().WithFields(log.Fields{
								"name":   p.Name,
								"detail": p.Detail,
							}).Log(p.Explanation)
						}
					*/
					break
				}

				logger.Info().Log("Successfully acquired certificate")
			}

			// Shut down the temporary http server
			tempserver.Close()

			wg.Wait()

			if certerror {
				a.log.logger.core.Warn().Log("Continuing with disabled TLS")
				autocertManager = nil
				cfg.TLS.Enable = false
			} else {
				cfg.TLS.CertFile = ""
				cfg.TLS.KeyFile = ""
			}
		} else {
			a.log.logger.core.Info().Log("Enabling TLS with cert and key files")
		}
	}

	if cfg.RTMP.Enable {
		a.log.logger.rtmp = a.log.logger.core.WithComponent("RTMP").WithField("address", cfg.RTMP.Address)

		config := rtmp.Config{
			Addr:                  cfg.RTMP.Address,
			TLSAddr:               cfg.RTMP.AddressTLS,
			App:                   cfg.RTMP.App,
			Token:                 cfg.RTMP.Token,
			Logger:                a.log.logger.rtmp,
			ConnectionIdleTimeout: 10 * time.Second,
			Collector:             a.sessions.Collector("rtmp"),
		}

		if cfg.RTMP.EnableTLS {
			config.Logger = config.Logger.WithComponent("RTMP/S")

			a.log.logger.rtmps = a.log.logger.core.WithComponent("RTMPS").WithField("address", cfg.RTMP.AddressTLS)
			if autocertManager != nil {
				config.TLSConfig = &tls.Config{
					GetCertificate: autocertManager.GetCertificate,
				}
			}
		}

		rtmpserver, err := rtmp.New(config)
		if err != nil {
			return fmt.Errorf("unable to create RMTP server: %w", err)
		}

		a.rtmpserver = rtmpserver
	}

	if cfg.SRT.Enable {
		config := srt.Config{
			Addr:       cfg.SRT.Address,
			Passphrase: cfg.SRT.Passphrase,
			Token:      cfg.SRT.Token,
			Logger:     a.log.logger.core.WithComponent("SRT").WithField("address", cfg.SRT.Address),
			Collector:  a.sessions.Collector("srt"),
		}

		if cfg.SRT.Log.Enable {
			config.SRTLogTopics = cfg.SRT.Log.Topics
		}

		srtserver, err := srt.New(config)
		if err != nil {
			return fmt.Errorf("unable to create SRT server: %w", err)
		}

		a.log.logger.srt = config.Logger
		a.srtserver = srtserver
	}

	logcontext := "HTTP"
	if cfg.TLS.Enable {
		logcontext = "HTTPS"
	}

	var iplimiter net.IPLimiter

	if cfg.TLS.Enable {
		limiter, err := net.NewIPLimiter(cfg.API.Access.HTTPS.Block, cfg.API.Access.HTTPS.Allow)
		if err != nil {
			return fmt.Errorf("incorret IP ranges for the API (HTTPS) provided: %w", err)
		}

		iplimiter = limiter
	} else {
		limiter, err := net.NewIPLimiter(cfg.API.Access.HTTP.Block, cfg.API.Access.HTTP.Allow)
		if err != nil {
			return fmt.Errorf("incorret IP ranges for the API (HTTP) provided: %w", err)
		}

		iplimiter = limiter
	}

	router, err := router.New(cfg.Router.BlockedPrefixes, cfg.Router.Routes, cfg.Router.UIPath)
	if err != nil {
		return fmt.Errorf("incorrect routes provided: %w", err)
	}

	a.log.logger.main = a.log.logger.core.WithComponent(logcontext).WithField("address", cfg.Address)

	httpfilesystems := []httpfs.FS{
		{
			Name:               a.diskfs.Name(),
			Mountpoint:         "",
			AllowWrite:         false,
			EnableAuth:         false,
			Username:           "",
			Password:           "",
			DefaultFile:        "index.html",
			DefaultContentType: "text/html",
			Gzip:               true,
			Filesystem:         a.diskfs,
			Cache:              a.cache,
		},
		{
			Name:               a.memfs.Name(),
			Mountpoint:         "/memfs",
			AllowWrite:         true,
			EnableAuth:         cfg.Storage.Memory.Auth.Enable,
			Username:           cfg.Storage.Memory.Auth.Username,
			Password:           cfg.Storage.Memory.Auth.Password,
			DefaultFile:        "",
			DefaultContentType: "application/data",
			Gzip:               true,
			Filesystem:         a.memfs,
			Cache:              nil,
		},
	}

	for _, s3 := range cfg.Storage.S3 {
		httpfilesystems = append(httpfilesystems, httpfs.FS{
			Name:               s3.Name,
			Mountpoint:         s3.Mountpoint,
			AllowWrite:         true,
			EnableAuth:         s3.Auth.Enable,
			Username:           s3.Auth.Username,
			Password:           s3.Auth.Password,
			DefaultFile:        "",
			DefaultContentType: "application/data",
			Gzip:               true,
			Filesystem:         a.s3fs[s3.Name],
			Cache:              a.cache,
		})
	}

	serverConfig := http.Config{
		Logger:        a.log.logger.main,
		LogBuffer:     a.log.buffer,
		Restream:      a.restream,
		Metrics:       a.metrics,
		Prometheus:    a.prom,
		MimeTypesFile: cfg.Storage.MimeTypes,
		Filesystems:   httpfilesystems,
		IPLimiter:     iplimiter,
		Profiling:     cfg.Debug.Profiling,
		Cors: http.CorsConfig{
			Origins: cfg.Storage.CORS.Origins,
		},
		RTMP:     a.rtmpserver,
		SRT:      a.srtserver,
		JWT:      a.httpjwt,
		Config:   a.config.store,
		Sessions: a.sessions,
		Router:   router,
		ReadOnly: cfg.API.ReadOnly,
	}

	mainserverhandler, err := http.NewServer(serverConfig)

	if err != nil {
		return fmt.Errorf("unable to create server: %w", err)
	}

	var wgStart sync.WaitGroup

	sendError := func(err error) {
		select {
		case a.errorChan <- err:
		default:
		}
	}

	a.mainserver = &gohttp.Server{
		Addr:              cfg.Address,
		Handler:           mainserverhandler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ErrorLog:          golog.New(a.log.logger.main.Debug(), "", 0),
	}

	if cfg.TLS.Enable {
		a.mainserver.Addr = cfg.TLS.Address
		a.log.logger.main = a.log.logger.main.WithField("address", cfg.TLS.Address)

		iplimiter, err := net.NewIPLimiter(cfg.API.Access.HTTP.Block, cfg.API.Access.HTTP.Allow)
		if err != nil {
			return fmt.Errorf("incorret IP ranges for the API (HTTP) provided: %w", err)
		}

		a.log.logger.sidecar = a.log.logger.core.WithComponent("HTTP").WithField("address", cfg.Address)

		serverConfig.Logger = a.log.logger.sidecar
		serverConfig.IPLimiter = iplimiter

		sidecarserverhandler, err := http.NewServer(serverConfig)

		if err != nil {
			return fmt.Errorf("unable to create sidecar HTTP server: %w", err)
		}

		a.sidecarserver = &gohttp.Server{
			Addr:              cfg.Address,
			Handler:           sidecarserverhandler,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20,
			ErrorLog:          golog.New(a.log.logger.sidecar.Debug(), "", 0),
		}

		if cfg.TLS.Auto {
			a.mainserver.TLSConfig = &tls.Config{
				GetCertificate: autocertManager.GetCertificate,
			}

			acme := autocertManager.Issuers[0].(*certmagic.ACMEIssuer)
			a.sidecarserver.Handler = acme.HTTPChallengeHandler(sidecarserverhandler)
		}

		wgStart.Add(1)
		a.wgStop.Add(1)

		go func() {
			logger := a.log.logger.sidecar
			defer func() {
				logger.Info().Log("Sidecar server exited")
				a.wgStop.Done()
			}()

			wgStart.Done()

			logger.Info().Log("Sidecar server started")

			var err error

			err = a.sidecarserver.ListenAndServe()
			if err != nil && err != gohttp.ErrServerClosed {
				err = fmt.Errorf("HTTP sidecar server: %w", err)
			} else {
				err = nil
			}

			sendError(err)
		}()
	}

	if a.rtmpserver != nil {
		wgStart.Add(1)
		a.wgStop.Add(1)

		go func() {
			logger := a.log.logger.rtmp

			defer func() {
				logger.Info().Log("Server exited")
				a.wgStop.Done()
			}()

			wgStart.Done()

			var err error

			logger.Info().Log("Server started")
			err = a.rtmpserver.ListenAndServe()
			if err != nil && err != rtmp.ErrServerClosed {
				err = fmt.Errorf("RTMP server: %w", err)
			} else {
				err = nil
			}

			sendError(err)
		}()

		if cfg.TLS.Enable && cfg.RTMP.EnableTLS {
			wgStart.Add(1)
			a.wgStop.Add(1)

			go func() {
				logger := a.log.logger.rtmps

				defer func() {
					logger.Info().Log("Server exited")
					a.wgStop.Done()
				}()

				wgStart.Done()

				var err error

				logger.Info().Log("Server started")
				err = a.rtmpserver.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
				if err != nil && err != rtmp.ErrServerClosed {
					err = fmt.Errorf("RTMPS server: %w", err)
				} else {
					err = nil
				}

				sendError(err)
			}()
		}
	}

	if a.srtserver != nil {
		wgStart.Add(1)
		a.wgStop.Add(1)

		go func() {
			logger := a.log.logger.srt

			defer func() {
				logger.Info().Log("Server exited")
				a.wgStop.Done()
			}()

			wgStart.Done()

			var err error

			logger.Info().Log("Server started")
			err = a.srtserver.ListenAndServe()
			if err != nil && err != srt.ErrServerClosed {
				err = fmt.Errorf("SRT server: %w", err)
			} else {
				err = nil
			}

			sendError(err)
		}()
	}

	wgStart.Add(1)
	a.wgStop.Add(1)

	go func() {
		logger := a.log.logger.main

		defer func() {
			logger.Info().Log("Server exited")
			a.wgStop.Done()
		}()

		wgStart.Done()

		var err error

		if cfg.TLS.Enable {
			logger.Info().Log("Server started")
			err = a.mainserver.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
			if err != nil && err != gohttp.ErrServerClosed {
				err = fmt.Errorf("HTTPS server: %w", err)
			} else {
				err = nil
			}
		} else {
			logger.Info().Log("Server started")
			err = a.mainserver.ListenAndServe()
			if err != nil && err != gohttp.ErrServerClosed {
				err = fmt.Errorf("HTTP server: %w", err)
			} else {
				err = nil
			}
		}

		sendError(err)
	}()

	// Wait for all servers to be started
	wgStart.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	a.gcTickerStop = cancel

	if cfg.Debug.ForceGC > 0 {
		go func(ctx context.Context) {
			ticker := time.NewTicker(time.Duration(cfg.Debug.ForceGC) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					debug.FreeOSMemory()
				}
			}
		}(ctx)
	}

	if cfg.Debug.MemoryLimit > 0 {
		debug.SetMemoryLimit(cfg.Debug.MemoryLimit * 1024 * 1024)
	} else {
		debug.SetMemoryLimit(math.MaxInt64)
	}

	// Start the restream processes
	restream.Start()

	// Start the service
	if a.service != nil {
		a.service.Start()
	}

	// Start the update checker
	if a.update != nil {
		a.update.Start()
	}

	a.state = "running"

	return nil
}

func (a *api) Start() error {
	if err := a.start(); err != nil {
		a.stop()
		return err
	}

	// Block until is an error from the servers
	err := <-a.errorChan

	return err
}

func (a *api) stop() {
	a.lock.Lock()
	defer a.lock.Unlock()

	logger := a.log.logger.core.WithField("action", "shutdown")

	if a.state == "idle" {
		logger.Info().Log("Complete")
		return
	}

	// Stop JWT authentication
	if a.httpjwt != nil {
		a.httpjwt.ClearValidators()
	}

	if a.update != nil {
		a.update.Stop()
		a.update = nil
	}

	if a.service != nil {
		a.service.Stop()
		a.service = nil
	}

	// Stop all restream processes
	if a.restream != nil {
		logger.Info().Log("Stopping all processes ...")
		a.restream.Stop()
		a.restream = nil
	}

	// Stop the session tracker
	if a.sessions != nil {
		a.sessions.UnregisterAll()
		a.sessions = nil
	}

	// Unregister all collectors
	if a.metrics != nil {
		a.metrics.UnregisterAll()
		a.metrics = nil
	}

	if a.prom != nil {
		a.prom.UnregisterAll()
		a.prom = nil
	}

	// Free the cached objects
	if a.cache != nil {
		a.cache.Purge()
		a.cache = nil
	}

	// Free the S3 mounts
	a.s3fs = map[string]fs.Filesystem{}

	// Stop the SRT server
	if a.srtserver != nil {
		a.log.logger.srt.Info().Log("Stopping ...")

		a.srtserver.Close()
		a.srtserver = nil
	}

	// Stop the RTMP server
	if a.rtmpserver != nil {
		a.log.logger.rtmp.Info().Log("Stopping ...")

		if a.log.logger.rtmps != nil {
			a.log.logger.rtmps.Info().Log("Stopping ...")
		}

		a.rtmpserver.Close()
		a.rtmpserver = nil
	}

	// Shutdown the HTTP/S mainserver
	if a.mainserver != nil {
		logger := a.log.logger.main
		logger.Info().Log("Stopping ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.mainserver.Shutdown(ctx); err != nil {
			logger.Error().WithError(err).Log("")
		}

		a.mainserver = nil
	}

	// Shutdown the HTTP sidecar server
	if a.sidecarserver != nil {
		logger := a.log.logger.sidecar

		logger.Info().Log("Stopping ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.sidecarserver.Shutdown(ctx); err != nil {
			logger.Error().WithError(err).Log("")
		}

		a.sidecarserver = nil
	}

	// Stop the GC ticker
	if a.gcTickerStop != nil {
		a.gcTickerStop()
		a.gcTickerStop = nil
	}

	// Wait for all server goroutines to exit
	logger.Info().Log("Waiting for all servers to stop ...")
	a.wgStop.Wait()

	// Drain error channel
	if a.errorChan != nil {
		close(a.errorChan)
		a.errorChan = nil
	}

	a.state = "idle"

	logger.Info().Log("Complete")
}

func (a *api) Stop() {
	a.log.logger.core.Info().Log("Shutdown requested ...")
	a.stop()
}

func (a *api) Destroy() {
	a.log.logger.core.Info().Log("Shutdown requested ...")
	a.stop()

	// Free the MemFS
	if a.memfs != nil {
		a.memfs.RemoveAll()
		a.memfs = nil
	}
}
