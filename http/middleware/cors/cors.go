package cors

import (
	"fmt"
	"strings"
	"time"

	"github.com/datarhei/core/v16/http/cors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper  middleware.Skipper
	Prefixes map[string][]string
}

var DefaultConfig = Config{
	Skipper:  middleware.DefaultSkipper,
	Prefixes: nil,
}

func New() echo.MiddlewareFunc {
	mw, _ := NewWithConfig(DefaultConfig)

	return mw
}

func NewWithConfig(config Config) (echo.MiddlewareFunc, error) {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	prefixes := make(map[string]echo.MiddlewareFunc)

	for prefix, origins := range config.Prefixes {
		if err := cors.Validate(origins); err != nil {
			return nil, fmt.Errorf("CORS config for prefix %s is invalid: %w", prefix, err)
		}

		conf := middleware.CORSConfig{
			AllowOrigins:     origins,
			AllowMethods:     []string{"GET", "HEAD", "PUT", "POST", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           int((24 * time.Hour).Seconds()),
		}

		prefixes[prefix] = middleware.CORSWithConfig(conf)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			//if origin := c.Request.Header.Get("Origin"); len(origin) == 0 {
			//	c.Request.Header.Set("Origin", "*")
			//}

			path := c.Request().URL.Path

			var middlewareFunc echo.MiddlewareFunc
			var maxPrefixLen int

			for prefix, h := range prefixes {
				if strings.HasPrefix(path, prefix) {
					if len(prefix) > maxPrefixLen {
						maxPrefixLen = len(prefix)
						middlewareFunc = h
					}
				}
			}

			if middlewareFunc != nil {
				handler := middlewareFunc(next)
				return handler(c)
			}

			return next(c)
		}
	}, nil
}
