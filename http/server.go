// @title datarhei Core API
// @version 3.0
// @description Expose REST API for the datarhei Core

// @contact.name datarhei Core Support
// @contact.url https://www.datarhei.com
// @contact.email hello@datarhei.com

// @license.name Apache 2.0
// @license.url https://github.com/datarhei/core/v16/blob/main/LICENSE

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @securityDefinitions.apikey ApiRefreshKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey Auth0KeyAuth
// @in header
// @name Authorization

// @securityDefinitions.basic BasicAuth

package http

import (
	"fmt"
	"net/http"
	"strings"

	cfgstore "github.com/datarhei/core/v16/config/store"
	"github.com/datarhei/core/v16/http/cache"
	"github.com/datarhei/core/v16/http/errorhandler"
	"github.com/datarhei/core/v16/http/fs"
	"github.com/datarhei/core/v16/http/graph/resolver"
	"github.com/datarhei/core/v16/http/handler"
	api "github.com/datarhei/core/v16/http/handler/api"
	"github.com/datarhei/core/v16/http/jwt"
	"github.com/datarhei/core/v16/http/router"
	"github.com/datarhei/core/v16/http/validator"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/monitor"
	"github.com/datarhei/core/v16/net"
	"github.com/datarhei/core/v16/prometheus"
	"github.com/datarhei/core/v16/restream"
	"github.com/datarhei/core/v16/rtmp"
	"github.com/datarhei/core/v16/session"
	"github.com/datarhei/core/v16/srt"

	mwcache "github.com/datarhei/core/v16/http/middleware/cache"
	mwcors "github.com/datarhei/core/v16/http/middleware/cors"
	mwgzip "github.com/datarhei/core/v16/http/middleware/gzip"
	mwhlsrewrite "github.com/datarhei/core/v16/http/middleware/hlsrewrite"
	mwiplimit "github.com/datarhei/core/v16/http/middleware/iplimit"
	mwlog "github.com/datarhei/core/v16/http/middleware/log"
	mwmime "github.com/datarhei/core/v16/http/middleware/mime"
	mwredirect "github.com/datarhei/core/v16/http/middleware/redirect"
	mwsession "github.com/datarhei/core/v16/http/middleware/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware

	// Expose the API docs
	_ "github.com/datarhei/core/v16/docs"
)

var ListenAndServe = http.ListenAndServe

type Config struct {
	Logger        log.Logger
	LogBuffer     log.BufferWriter
	Restream      restream.Restreamer
	Metrics       monitor.HistoryReader
	Prometheus    prometheus.Reader
	MimeTypesFile string
	Filesystems   []fs.FS
	IPLimiter     net.IPLimiter
	Profiling     bool
	Cors          CorsConfig
	RTMP          rtmp.Server
	SRT           srt.Server
	JWT           jwt.JWT
	Config        cfgstore.Store
	Cache         cache.Cacher
	Sessions      session.RegistryReader
	Router        router.Router
	ReadOnly      bool
}

type CorsConfig struct {
	Origins []string
}

type Server interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type server struct {
	logger log.Logger

	handler struct {
		about      *api.AboutHandler
		prometheus *handler.PrometheusHandler
		profiling  *handler.ProfilingHandler
		ping       *handler.PingHandler
		graph      *api.GraphHandler
		jwt        jwt.JWT
	}

	v3handler struct {
		log       *api.LogHandler
		restream  *api.RestreamHandler
		playout   *api.PlayoutHandler
		rtmp      *api.RTMPHandler
		srt       *api.SRTHandler
		config    *api.ConfigHandler
		session   *api.SessionHandler
		widget    *api.WidgetHandler
		resources *api.MetricsHandler
	}

	middleware struct {
		iplimit    echo.MiddlewareFunc
		log        echo.MiddlewareFunc
		accessJWT  echo.MiddlewareFunc
		refreshJWT echo.MiddlewareFunc
		cors       echo.MiddlewareFunc
		cache      echo.MiddlewareFunc
		session    echo.MiddlewareFunc
		hlsrewrite echo.MiddlewareFunc
	}

	gzip struct {
		mimetypes []string
	}

	filesystems map[string]*filesystem

	router        *echo.Echo
	mimeTypesFile string
	profiling     bool

	readOnly bool
}

type filesystem struct {
	fs.FS

	handler    *handler.FSHandler
	middleware echo.MiddlewareFunc
}

func NewServer(config Config) (Server, error) {
	s := &server{
		logger:        config.Logger,
		mimeTypesFile: config.MimeTypesFile,
		profiling:     config.Profiling,
		readOnly:      config.ReadOnly,
	}

	s.filesystems = map[string]*filesystem{}

	corsPrefixes := map[string][]string{
		"/api": {"*"},
	}

	for _, fs := range config.Filesystems {
		if _, ok := s.filesystems[fs.Name]; ok {
			return nil, fmt.Errorf("the filesystem name '%s' is already in use", fs.Name)
		}

		if !strings.HasPrefix(fs.Mountpoint, "/") {
			fs.Mountpoint = "/" + fs.Mountpoint
		}

		if !strings.HasSuffix(fs.Mountpoint, "/") {
			fs.Mountpoint = strings.TrimSuffix(fs.Mountpoint, "/")
		}

		if _, ok := corsPrefixes[fs.Mountpoint]; ok {
			return nil, fmt.Errorf("the mount point '%s' is already in use (%s)", fs.Mountpoint, fs.Name)
		}

		corsPrefixes[fs.Mountpoint] = config.Cors.Origins

		filesystem := &filesystem{
			FS:      fs,
			handler: handler.NewFS(fs),
		}

		if fs.Filesystem.Type() == "disk" {
			filesystem.middleware = mwhlsrewrite.NewHLSRewriteWithConfig(mwhlsrewrite.HLSRewriteConfig{
				PathPrefix: fs.Filesystem.Metadata("base"),
			})
		}

		s.filesystems[filesystem.Name] = filesystem
	}

	if _, ok := corsPrefixes["/"]; !ok {
		return nil, fmt.Errorf("one filesystem must be mounted at /")
	}

	if config.Logger == nil {
		s.logger = log.New("HTTP")
	}

	if config.JWT == nil {
		s.handler.about = api.NewAbout(
			config.Restream,
			[]string{},
		)
	} else {
		s.handler.about = api.NewAbout(
			config.Restream,
			config.JWT.Validators(),
		)
	}

	s.v3handler.log = api.NewLog(
		config.LogBuffer,
	)

	if config.Restream != nil {
		s.v3handler.restream = api.NewRestream(
			config.Restream,
		)

		s.v3handler.playout = api.NewPlayout(
			config.Restream,
		)
	}

	if config.Prometheus != nil {
		s.handler.prometheus = handler.NewPrometheus(
			config.Prometheus.HTTPHandler(),
		)
	}

	if config.Profiling {
		s.handler.profiling = handler.NewProfiling()
	}

	if config.IPLimiter != nil {
		s.middleware.iplimit = mwiplimit.NewWithConfig(mwiplimit.Config{
			Limiter: config.IPLimiter,
		})
	}

	s.handler.ping = handler.NewPing()

	if config.RTMP != nil {
		s.v3handler.rtmp = api.NewRTMP(
			config.RTMP,
		)
	}

	if config.SRT != nil {
		s.v3handler.srt = api.NewSRT(
			config.SRT,
		)
	}

	if config.Config != nil {
		s.v3handler.config = api.NewConfig(
			config.Config,
		)
	}

	if config.JWT != nil {
		s.handler.jwt = config.JWT
		s.middleware.accessJWT = config.JWT.AccessMiddleware()
		s.middleware.refreshJWT = config.JWT.RefreshMiddleware()
	}

	if config.Sessions == nil {
		config.Sessions, _ = session.New(session.Config{})
	}

	s.v3handler.session = api.NewSession(
		config.Sessions,
	)

	s.middleware.session = mwsession.NewHLSWithConfig(mwsession.HLSConfig{
		EgressCollector:  config.Sessions.Collector("hls"),
		IngressCollector: config.Sessions.Collector("hlsingress"),
	})

	s.middleware.log = mwlog.NewWithConfig(mwlog.Config{
		Logger: s.logger,
	})

	s.v3handler.widget = api.NewWidget(api.WidgetConfig{
		Restream: config.Restream,
		Registry: config.Sessions,
	})

	s.v3handler.resources = api.NewMetrics(api.MetricsConfig{
		Metrics: config.Metrics,
	})

	if middleware, err := mwcors.NewWithConfig(mwcors.Config{
		Prefixes: corsPrefixes,
	}); err != nil {
		return nil, err
	} else {
		s.middleware.cors = middleware
	}

	s.handler.graph = api.NewGraph(resolver.Resolver{
		Restream:  config.Restream,
		Monitor:   config.Metrics,
		LogBuffer: config.LogBuffer,
	}, "/api/graph/query")

	s.gzip.mimetypes = []string{
		"text/plain",
		"text/html",
		"text/javascript",
		"application/json",
		"application/x-mpegurl",
		"application/vnd.apple.mpegurl",
		"image/svg+xml",
	}

	s.router = echo.New()
	s.router.HTTPErrorHandler = errorhandler.HTTPErrorHandler
	s.router.Validator = validator.New()
	s.router.Use(s.middleware.log)
	s.router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			rows := strings.Split(string(stack), "\n")
			s.logger.Error().WithField("stack", rows).Log("recovered from a panic")
			return nil
		},
	}))
	s.router.Use(mwsession.NewHTTPWithConfig(mwsession.HTTPConfig{
		Collector: config.Sessions.Collector("http"),
	}))

	s.router.HideBanner = true
	s.router.HidePort = true

	s.router.Logger.SetOutput(newLogwrapper(s.logger))

	if s.middleware.cors != nil {
		s.router.Use(s.middleware.cors)
	}

	// Add static routes
	if path, target := config.Router.StaticRoute(); len(target) != 0 {
		group := s.router.Group(path)
		group.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
			Skipper: func(c echo.Context) bool {
				return path != c.Request().URL.Path
			},
			RedirectCode: 301,
		}))
		group.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper:    middleware.DefaultSkipper,
			Root:       target,
			Index:      "index.html",
			IgnoreBase: true,
		}))
	}

	s.router.Use(mwredirect.NewWithConfig(mwredirect.Config{
		Redirects: config.Router.FileRoutes(),
	}))

	for prefix, target := range config.Router.DirRoutes() {
		group := s.router.Group(prefix)
		group.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
			Skipper: func(prefix string) func(c echo.Context) bool {
				return func(c echo.Context) bool {
					return prefix != c.Request().URL.Path
				}
			}(prefix),
			RedirectCode: 301,
		}))
		group.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper:    middleware.DefaultSkipper,
			Root:       target,
			Index:      "index.html",
			IgnoreBase: true,
		}))
	}

	s.setRoutes()

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) setRoutes() {
	gzipMiddleware := mwgzip.NewWithConfig(mwgzip.Config{
		Level:     mwgzip.BestSpeed,
		MinLength: 1000,
		Skipper:   mwgzip.ContentTypeSkipper(nil),
	})

	// API router grouo
	api := s.router.Group("/api")

	if s.middleware.iplimit != nil {
		api.Use(s.middleware.iplimit)
	}

	if s.middleware.accessJWT != nil {
		// Enable JWT auth
		api.Use(s.middleware.accessJWT)

		// The login endpoint should not be blocked by auth
		s.router.POST("/api/login", s.handler.jwt.LoginHandler)
		s.router.GET("/api/login/refresh", s.handler.jwt.RefreshHandler, s.middleware.refreshJWT)
	}

	api.GET("", s.handler.about.About)

	// Swagger API documentation router group
	doc := s.router.Group("/api/swagger/*")
	doc.Use(gzipMiddleware)
	doc.GET("", echoSwagger.WrapHandler)

	// Mount filesystems
	for _, filesystem := range s.filesystems {
		// Define a local variable because later in the loop we have a closure
		filesystem := filesystem

		mountpoint := filesystem.Mountpoint + "/*"
		if filesystem.Mountpoint == "/" {
			mountpoint = "/*"
		}

		fs := s.router.Group(mountpoint)
		fs.Use(mwmime.NewWithConfig(mwmime.Config{
			MimeTypesFile:      s.mimeTypesFile,
			DefaultContentType: filesystem.DefaultContentType,
		}))

		if filesystem.Gzip {
			fs.Use(mwgzip.NewWithConfig(mwgzip.Config{
				Skipper:   mwgzip.ContentTypeSkipper(s.gzip.mimetypes),
				Level:     mwgzip.BestSpeed,
				MinLength: 1000,
			}))
		}

		if filesystem.Cache != nil {
			mwcache := mwcache.NewWithConfig(mwcache.Config{
				Cache: filesystem.Cache,
			})
			fs.Use(mwcache)
		}

		if filesystem.middleware != nil {
			fs.Use(filesystem.middleware)
		}

		if s.middleware.session != nil {
			fs.Use(s.middleware.session)
		}

		fs.GET("", filesystem.handler.GetFile)
		fs.HEAD("", filesystem.handler.GetFile)

		if filesystem.AllowWrite {
			if filesystem.EnableAuth {
				authmw := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					if username == filesystem.Username && password == filesystem.Password {
						return true, nil
					}

					return false, nil
				})

				fs.POST("", filesystem.handler.PutFile, authmw)
				fs.PUT("", filesystem.handler.PutFile, authmw)
				fs.DELETE("", filesystem.handler.DeleteFile, authmw)
			} else {
				fs.POST("", filesystem.handler.PutFile)
				fs.PUT("", filesystem.handler.PutFile)
				fs.DELETE("", filesystem.handler.DeleteFile)
			}
		}
	}

	// Prometheus metrics
	if s.handler.prometheus != nil {
		metrics := s.router.Group("/metrics")

		if s.middleware.iplimit != nil {
			metrics.Use(s.middleware.iplimit)
		}

		metrics.GET("", s.handler.prometheus.Metrics)
	}

	// Health check
	s.router.GET("/ping", s.handler.ping.Ping)

	// Profiling routes
	if s.profiling {
		prof := s.router.Group("/profiling")

		if s.middleware.iplimit != nil {
			prof.Use(s.middleware.iplimit)
		}

		s.handler.profiling.Register(prof)
	}

	// GraphQL
	graphql := api.Group("/graph")
	graphql.Use(gzipMiddleware)

	graphql.GET("", s.handler.graph.Playground)
	graphql.POST("/query", s.handler.graph.Query)

	// APIv3 router group
	v3 := api.Group("/v3")

	if s.handler.jwt != nil {
		v3.Use(s.middleware.accessJWT)
	}

	v3.Use(gzipMiddleware)

	s.setRoutesV3(v3)
}

func (s *server) setRoutesV3(v3 *echo.Group) {
	if s.v3handler.widget != nil {
		// The widget endpoint should not be blocked by auth
		s.router.GET("/api/v3/widget/process/:id", s.v3handler.widget.Get)
	}

	// v3 Restreamer
	if s.v3handler.restream != nil {
		v3.GET("/skills", s.v3handler.restream.Skills)
		v3.GET("/skills/reload", s.v3handler.restream.ReloadSkills)

		v3.GET("/process", s.v3handler.restream.GetAll)
		v3.GET("/process/:id", s.v3handler.restream.Get)

		v3.GET("/process/:id/config", s.v3handler.restream.GetConfig)
		v3.GET("/process/:id/state", s.v3handler.restream.GetState)
		v3.GET("/process/:id/report", s.v3handler.restream.GetReport)
		v3.GET("/process/:id/probe", s.v3handler.restream.Probe)

		v3.GET("/process/:id/metadata", s.v3handler.restream.GetProcessMetadata)
		v3.GET("/process/:id/metadata/:key", s.v3handler.restream.GetProcessMetadata)

		v3.GET("/metadata", s.v3handler.restream.GetMetadata)
		v3.GET("/metadata/:key", s.v3handler.restream.GetMetadata)

		if !s.readOnly {
			v3.POST("/process", s.v3handler.restream.Add)
			v3.PUT("/process/:id", s.v3handler.restream.Update)
			v3.DELETE("/process/:id", s.v3handler.restream.Delete)
			v3.PUT("/process/:id/command", s.v3handler.restream.Command)
			v3.PUT("/process/:id/metadata/:key", s.v3handler.restream.SetProcessMetadata)
			v3.PUT("/metadata/:key", s.v3handler.restream.SetMetadata)
		}

		// v3 Playout
		if s.v3handler.playout != nil {
			v3.GET("/process/:id/playout/:inputid/status", s.v3handler.playout.Status)
			v3.GET("/process/:id/playout/:inputid/reopen", s.v3handler.playout.ReopenInput)
			v3.GET("/process/:id/playout/:inputid/keyframe/*", s.v3handler.playout.Keyframe)
			v3.GET("/process/:id/playout/:inputid/errorframe/encode", s.v3handler.playout.EncodeErrorframe)

			if !s.readOnly {
				v3.PUT("/process/:id/playout/:inputid/errorframe/*", s.v3handler.playout.SetErrorframe)
				v3.POST("/process/:id/playout/:inputid/errorframe/*", s.v3handler.playout.SetErrorframe)

				v3.PUT("/process/:id/playout/:inputid/stream", s.v3handler.playout.SetStream)
			}
		}
	}

	// v3 Filesystems
	fshandlers := map[string]api.FSConfig{}
	for _, fs := range s.filesystems {
		fshandlers[fs.Name] = api.FSConfig{
			Type:       fs.Filesystem.Type(),
			Mountpoint: fs.Mountpoint,
			Handler:    fs.handler,
		}
	}

	handler := api.NewFS(fshandlers)

	v3.GET("/fs", handler.List)

	v3.GET("/fs/:name", handler.ListFiles)
	v3.GET("/fs/:name/*", handler.GetFile, mwmime.NewWithConfig(mwmime.Config{
		MimeTypesFile:      s.mimeTypesFile,
		DefaultContentType: "application/data",
	}))
	v3.HEAD("/fs/:name/*", handler.GetFile, mwmime.NewWithConfig(mwmime.Config{
		MimeTypesFile:      s.mimeTypesFile,
		DefaultContentType: "application/data",
	}))

	if !s.readOnly {
		v3.PUT("/fs/:name/*", handler.PutFile)
		v3.DELETE("/fs/:name/*", handler.DeleteFile)
	}

	// v3 RTMP
	if s.v3handler.rtmp != nil {
		v3.GET("/rtmp", s.v3handler.rtmp.ListChannels)
	}

	// v3 SRT
	if s.v3handler.srt != nil {
		v3.GET("/srt", s.v3handler.srt.ListChannels)
	}

	// v3 Config
	if s.v3handler.config != nil {
		v3.GET("/config", s.v3handler.config.Get)

		if !s.readOnly {
			v3.PUT("/config", s.v3handler.config.Set)
			v3.GET("/config/reload", s.v3handler.config.Reload)
		}
	}

	// v3 Session
	if s.v3handler.session != nil {
		v3.GET("/session", s.v3handler.session.Summary)
		v3.GET("/session/active", s.v3handler.session.Active)
	}

	// v3 Log
	v3.GET("/log", s.v3handler.log.Log)

	// v3 Metrics
	v3.GET("/metrics", s.v3handler.resources.Describe)
	v3.POST("/metrics", s.v3handler.resources.Metrics)
}
