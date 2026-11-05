package iplimit

import (
	"net/http"

	"github.com/datarhei/core/v16/net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
	Limiter net.IPLimiter
}

var DefaultConfig = Config{
	Skipper: middleware.DefaultSkipper,
	Limiter: net.NewNullIPLimiter(),
}

func New() echo.MiddlewareFunc {
	return NewWithConfig(DefaultConfig)
}

func NewWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	if config.Limiter == nil {
		config.Limiter = DefaultConfig.Limiter
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			ip := c.RealIP()

			if !config.Limiter.IsAllowed(ip) {
				return echo.NewHTTPError(http.StatusForbidden)
			}

			return next(c)
		}
	}
}
