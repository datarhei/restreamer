// Package redirect is an echo middleware that applies defined redirects
package redirect

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper   middleware.Skipper
	Redirects map[string]string
}

var DefaultConfig = Config{
	Skipper:   middleware.DefaultSkipper,
	Redirects: map[string]string{},
}

func New() echo.MiddlewareFunc {
	return NewWithConfig(DefaultConfig)
}

// New returns a new router middleware handler
func NewWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			if target, ok := config.Redirects[req.URL.Path]; ok {
				return c.Redirect(http.StatusMovedPermanently, target)
			}

			return next(c)
		}
	}
}
