package session

import (
	"io"
	"net"
	"net/http"

	"github.com/datarhei/core/v16/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper   middleware.Skipper
	Collector session.Collector
}

var DefaultHTTPConfig = HTTPConfig{
	Skipper:   middleware.DefaultSkipper,
	Collector: session.NewNullCollector(),
}

// NewHTTP returns a new HTTP session middleware with default config
func NewHTTP() echo.MiddlewareFunc {
	return NewHTTPWithConfig(DefaultHTTPConfig)
}

// NewHTTPWithConfig returns a new HTTP session middleware
func NewHTTPWithConfig(config HTTPConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultHTTPConfig.Skipper
	}

	if config.Collector == nil {
		config.Collector = DefaultHTTPConfig.Collector
		config.Skipper = func(c echo.Context) bool { return true }
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				host = ""
			}

			reader := req.Body
			r := &fakeReader{
				reader: req.Body,
			}
			req.Body = r

			defer func() {
				req.Body = reader

				if config.Collector.IsCollectableIP(host) {
					config.Collector.RegisterAndActivate("HTTP", "", "any", "any")
					config.Collector.Ingress("HTTP", r.size+headerSize(req.Header))
				}
			}()

			writer := res.Writer
			w := &fakeWriter{
				ResponseWriter: res.Writer,
			}
			res.Writer = w

			defer func() {
				res.Writer = writer

				if config.Collector.IsCollectableIP(host) {
					config.Collector.Egress("HTTP", w.size+headerSize(res.Header()))
				}
			}()

			return next(c)
		}
	}
}

type fakeReader struct {
	reader io.ReadCloser
	size   int64
}

func (r *fakeReader) Read(b []byte) (int, error) {
	n, err := r.reader.Read(b)
	r.size += int64(n)

	return n, err
}

func (r *fakeReader) Close() error {
	return r.reader.Close()
}

type fakeWriter struct {
	http.ResponseWriter
	size int64
}

func (w *fakeWriter) Write(body []byte) (int, error) {
	n, err := w.ResponseWriter.Write(body)

	w.size += int64(n)

	return n, err
}

func (w *fakeWriter) Flush() {}
