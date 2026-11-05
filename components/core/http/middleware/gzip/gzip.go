package gzip

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config defines the config for Gzip middleware.
type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Gzip compression level.
	// Optional. Default value -1.
	Level int

	// Length threshold before gzip compression
	// is used. Optional. Default value 0
	MinLength int
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	wroteHeader       bool
	wroteBody         bool
	minLength         int
	minLengthExceeded bool
	buffer            *bytes.Buffer
	code              int
}

const gzipScheme = "gzip"

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// DefaultConfig is the default Gzip middleware config.
var DefaultConfig = Config{
	Skipper:   middleware.DefaultSkipper,
	Level:     DefaultCompression,
	MinLength: 0,
}

// ContentTypesSkipper returns a Skipper based on the list of content types
// that should be compressed. If the list is empty, all responses will be
// compressed.
func ContentTypeSkipper(contentTypes []string) middleware.Skipper {
	return func(c echo.Context) bool {
		// If no allowed content types are given, compress all
		if len(contentTypes) == 0 {
			return false
		}

		// Iterate through the allowed content types and don't skip if the content type matches
		responseContentType := c.Response().Header().Get(echo.HeaderContentType)

		for _, contentType := range contentTypes {
			if strings.Contains(responseContentType, contentType) {
				return false
			}
		}

		return true
	}
}

// New returns a middleware which compresses HTTP response using gzip compression
// scheme.
func New() echo.MiddlewareFunc {
	return NewWithConfig(DefaultConfig)
}

// NewWithConfig return Gzip middleware with config.
// See: `New()`.
func NewWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	if config.Level == 0 {
		config.Level = DefaultConfig.Level
	}

	if config.MinLength < 0 {
		config.MinLength = DefaultConfig.MinLength
	}

	pool := gzipPool(config)
	bpool := bufferPool()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			res := c.Response()
			res.Header().Add(echo.HeaderVary, echo.HeaderAcceptEncoding)

			if strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), gzipScheme) {
				i := pool.Get()
				w, ok := i.(*gzip.Writer)
				if !ok {
					return echo.NewHTTPError(http.StatusInternalServerError, i.(error).Error())
				}
				rw := res.Writer
				w.Reset(rw)

				buf := bpool.Get().(*bytes.Buffer)
				buf.Reset()

				grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw, minLength: config.MinLength, buffer: buf}

				defer func() {
					if !grw.wroteBody {
						if res.Header().Get(echo.HeaderContentEncoding) == gzipScheme {
							res.Header().Del(echo.HeaderContentEncoding)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						res.Writer = rw
						w.Reset(io.Discard)
					} else if !grw.minLengthExceeded {
						// If the minimum content length hasn't exceeded, write the uncompressed response
						res.Writer = rw
						if grw.wroteHeader {
							grw.ResponseWriter.WriteHeader(grw.code)
						}
						grw.buffer.WriteTo(rw)
						w.Reset(io.Discard)
					}
					w.Close()
					bpool.Put(buf)
					pool.Put(w)
				}()

				res.Writer = grw
			}

			return next(c)
		}
	}
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(echo.HeaderContentEncoding)
	}
	w.Header().Del(echo.HeaderContentLength) // Issue #444

	w.wroteHeader = true

	// Delay writing of the header until we know if we'll actually compress the response
	w.code = code
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(echo.HeaderContentType) == "" {
		w.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}

	w.wroteBody = true

	if !w.minLengthExceeded {
		n, err := w.buffer.Write(b)

		if w.buffer.Len() >= w.minLength {
			w.minLengthExceeded = true

			// The minimum length is exceeded, add Content-Encoding header and write the header
			w.Header().Set(echo.HeaderContentEncoding, gzipScheme) // Issue #806
			if w.wroteHeader {
				w.ResponseWriter.WriteHeader(w.code)
			}

			return w.Writer.Write(w.buffer.Bytes())
		} else {
			return n, err
		}
	}

	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	if !w.minLengthExceeded {
		// Enforce compression
		w.minLengthExceeded = true
		w.Header().Set(echo.HeaderContentEncoding, gzipScheme) // Issue #806
		if w.wroteHeader {
			w.ResponseWriter.WriteHeader(w.code)
		}

		w.Writer.Write(w.buffer.Bytes())
	}

	w.Writer.(*gzip.Writer).Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

func gzipPool(config Config) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(io.Discard, config.Level)
			if err != nil {
				return err
			}
			return w
		},
	}
}

func bufferPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			b := &bytes.Buffer{}
			return b
		},
	}
}
