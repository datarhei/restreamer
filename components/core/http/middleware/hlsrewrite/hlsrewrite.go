package hlsrewrite

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HLSRewriteConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper    middleware.Skipper
	PathPrefix string
}

var DefaultHLSRewriteConfig = HLSRewriteConfig{
	Skipper: func(c echo.Context) bool {
		req := c.Request()

		return !strings.HasSuffix(req.URL.Path, ".m3u8")
	},
	PathPrefix: "",
}

// NewHTTP returns a new HTTP session middleware with default config
func NewHLSRewrite() echo.MiddlewareFunc {
	return NewHLSRewriteWithConfig(DefaultHLSRewriteConfig)
}

type hlsrewrite struct {
	pathPrefix string
}

func NewHLSRewriteWithConfig(config HLSRewriteConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultHLSRewriteConfig.Skipper
	}

	pathPrefix := config.PathPrefix
	if len(pathPrefix) != 0 {
		if !strings.HasSuffix(pathPrefix, "/") {
			pathPrefix += "/"
		}
	}

	hls := hlsrewrite{
		pathPrefix: pathPrefix,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			if req.Method == "GET" || req.Method == "HEAD" {
				return hls.rewrite(c, next)
			}

			return next(c)
		}
	}
}

func (h *hlsrewrite) rewrite(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()

	path := req.URL.Path

	isM3U8 := strings.HasSuffix(path, ".m3u8")

	rewrite := false

	if isM3U8 {
		rewrite = true
	}

	var rewriter *hlsRewriter

	// Keep the current writer for later
	writer := res.Writer

	if rewrite {
		// Put the session rewriter in the middle. This will collect
		// the data that we need to rewrite.
		rewriter = &hlsRewriter{
			ResponseWriter: res.Writer,
		}

		res.Writer = rewriter
	}

	if err := next(c); err != nil {
		c.Error(err)
	}

	// Restore the original writer
	res.Writer = writer

	if rewrite {
		if res.Status != 200 {
			res.Write(rewriter.buffer.Bytes())
			return nil
		}

		// Rewrite the data befor sending it to the client
		rewriter.rewrite(h.pathPrefix)

		res.Header().Set("Cache-Control", "private")
		res.Write(rewriter.buffer.Bytes())
	}

	return nil
}

type hlsRewriter struct {
	http.ResponseWriter
	buffer bytes.Buffer
}

func (g *hlsRewriter) Write(data []byte) (int, error) {
	// Write the data into internal buffer for later rewrite
	w, err := g.buffer.Write(data)

	return w, err
}

func (g *hlsRewriter) rewrite(pathPrefix string) {
	var buffer bytes.Buffer

	// Find all URLS in the .m3u8 and add the session ID to the query string
	scanner := bufio.NewScanner(&g.buffer)
	for scanner.Scan() {
		line := scanner.Text()

		// Write empty lines unmodified
		if len(line) == 0 {
			buffer.WriteString(line + "\n")
			continue
		}

		// Write comments unmodified
		if strings.HasPrefix(line, "#") {
			buffer.WriteString(line + "\n")
			continue
		}

		// Rewrite
		line = strings.TrimPrefix(line, pathPrefix)
		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return
	}

	g.buffer = buffer
}
