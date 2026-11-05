// Package session is a HLS session middleware for Gin
package session

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/url"
	urlpath "path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/datarhei/core/v16/net"
	"github.com/datarhei/core/v16/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/shortuuid/v4"
)

type HLSConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper          middleware.Skipper
	EgressCollector  session.Collector
	IngressCollector session.Collector
}

var DefaultHLSConfig = HLSConfig{
	Skipper:          middleware.DefaultSkipper,
	EgressCollector:  session.NewNullCollector(),
	IngressCollector: session.NewNullCollector(),
}

// NewHTTP returns a new HTTP session middleware with default config
func NewHLS() echo.MiddlewareFunc {
	return NewHLSWithConfig(DefaultHLSConfig)
}

type hls struct {
	egressCollector  session.Collector
	ingressCollector session.Collector
	reSessionID      *regexp.Regexp

	rxsegments map[string]int64
	lock       sync.Mutex
}

// NewHLS returns a new HLS session middleware
func NewHLSWithConfig(config HLSConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultHLSConfig.Skipper
	}

	if config.EgressCollector == nil && config.IngressCollector == nil {
		config.Skipper = func(c echo.Context) bool { return true }
	}

	if config.EgressCollector == nil {
		config.EgressCollector = DefaultHLSConfig.EgressCollector
	}

	if config.IngressCollector == nil {
		config.IngressCollector = DefaultHLSConfig.IngressCollector
	}

	hls := hls{
		egressCollector:  config.EgressCollector,
		ingressCollector: config.IngressCollector,
		reSessionID:      regexp.MustCompile(`^[` + regexp.QuoteMeta(shortuuid.DefaultAlphabet) + `]{22}$`),
		rxsegments:       make(map[string]int64),
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			if req.Method == "PUT" || req.Method == "POST" {
				return hls.handleIngress(c, next)
			} else if req.Method == "GET" || req.Method == "HEAD" {
				return hls.handleEgress(c, next)
			}

			return next(c)
		}
	}
}

func (h *hls) handleIngress(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()

	path := req.URL.Path

	if strings.HasSuffix(path, ".m3u8") {
		// Read out the path of the .ts files and look them up in the ts-map.
		// Add it as ingress for the respective "sessionId". The "sessionId" is the .m3u8 file name.
		reader := req.Body
		r := &bodyReader{
			reader: req.Body,
		}
		req.Body = r

		defer func() {
			req.Body = reader

			if r.size == 0 {
				return
			}

			if !h.ingressCollector.IsKnownSession(path) {
				// Register a new session
				reference := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
				h.ingressCollector.RegisterAndActivate(path, reference, path, "")
				h.ingressCollector.Extra(path, req.Header.Get("User-Agent"))
			}

			h.ingressCollector.Ingress(path, headerSize(req.Header))
			h.ingressCollector.Ingress(path, r.size)

			segments := r.getSegments(urlpath.Dir(path))

			if len(segments) != 0 {
				h.lock.Lock()
				for _, s := range segments {
					if size, ok := h.rxsegments[s]; ok {
						// Update ingress value
						h.ingressCollector.Ingress(path, size)
						delete(h.rxsegments, s)
					}
				}
				h.lock.Unlock()
			}
		}()
	} else if strings.HasSuffix(path, ".ts") {
		// Get the size of the .ts file and store it in the ts-map for later use.
		reader := req.Body
		r := &bodysizeReader{
			reader: req.Body,
		}
		req.Body = r

		defer func() {
			req.Body = reader

			if r.size != 0 {
				h.lock.Lock()
				h.rxsegments[path] = r.size + headerSize(req.Header)
				h.lock.Unlock()
			}
		}()
	}

	return next(c)
}

func (h *hls) handleEgress(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()

	if !h.egressCollector.IsCollectableIP(c.RealIP()) {
		return next(c)
	}

	path := req.URL.Path
	sessionID := c.QueryParam("session")

	isM3U8 := strings.HasSuffix(path, ".m3u8")
	isTS := strings.HasSuffix(path, ".ts")

	rewrite := false

	if isM3U8 {
		if !h.egressCollector.IsKnownSession(sessionID) {
			if h.egressCollector.IsSessionsExceeded() {
				return echo.NewHTTPError(509, "Number of sessions exceeded")
			}

			streamBitrate := h.ingressCollector.SessionTopIngressBitrate(path) * 2.0 // Multiply by 2 to cover the initial peak
			maxBitrate := h.egressCollector.MaxEgressBitrate()

			if maxBitrate > 0.0 {
				currentBitrate := h.egressCollector.CompanionTopEgressBitrate() * 1.15

				// Add the new session's top bitrate to the ingress top bitrate
				resultingBitrate := currentBitrate + streamBitrate

				if resultingBitrate >= maxBitrate {
					return echo.NewHTTPError(509, "Bitrate limit exceeded")
				}
			}

			if len(sessionID) != 0 {
				if !h.reSessionID.MatchString(sessionID) {
					return echo.NewHTTPError(http.StatusForbidden)
				}

				referrer := req.Header.Get("Referer")
				if u, err := url.Parse(referrer); err == nil {
					referrer = u.Host
				}

				ip, _ := net.AnonymizeIPString(c.RealIP())
				extra := "[" + ip + "] " + req.Header.Get("User-Agent")

				reference := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

				// Register a new session
				h.egressCollector.Register(sessionID, reference, path, referrer)
				h.egressCollector.Extra(sessionID, extra)

				// Give the new session an initial top bitrate
				h.egressCollector.SessionSetTopEgressBitrate(sessionID, streamBitrate)
			}
		}

		rewrite = true
	}

	var rewriter *sessionRewriter

	// Keep the current writer for later
	writer := res.Writer

	if rewrite {
		// Put the session rewriter in the middle. This will collect
		// the data that we need to rewrite.
		rewriter = &sessionRewriter{
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
		rewriter.rewriteHLS(sessionID, c.Request().URL)

		res.Header().Set("Cache-Control", "private")
		res.Write(rewriter.buffer.Bytes())
	}

	if isM3U8 || isTS {
		// Collect how many bytes we've written in this session
		h.egressCollector.Egress(sessionID, headerSize(res.Header()))
		h.egressCollector.Egress(sessionID, res.Size)

		if isTS {
			// Activate the session. If the session is already active, this is a noop
			h.egressCollector.Activate(sessionID)
		}
	}

	return nil
}

func headerSize(header http.Header) int64 {
	var buffer bytes.Buffer

	header.Write(&buffer)

	return int64(buffer.Len())
}

type bodyReader struct {
	reader io.ReadCloser
	buffer bytes.Buffer
	size   int64
}

func (r *bodyReader) Read(b []byte) (int, error) {
	n, err := r.reader.Read(b)
	if n > 0 {
		r.buffer.Write(b[:n])
	}
	r.size += int64(n)

	return n, err
}

func (r *bodyReader) Close() error {
	return r.reader.Close()
}

func (r *bodyReader) getSegments(dir string) []string {
	segments := []string{}

	// Find all segment URLs in the .m3u8
	scanner := bufio.NewScanner(&r.buffer)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty lines
		if len(line) == 0 {
			continue
		}

		// Ignore comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		u, err := url.Parse(line)
		if err != nil {
			// Invalid URL
			continue
		}

		if u.Scheme != "" {
			// Ignore full URLs
			continue
		}

		// Ignore anything that doesn't end in .ts
		if !strings.HasSuffix(u.Path, ".ts") {
			continue
		}

		path := u.Path

		if !strings.HasPrefix(u.Path, "/") {
			path = urlpath.Join(dir, u.Path)
		}

		segments = append(segments, path)
	}

	return segments
}

type bodysizeReader struct {
	reader io.ReadCloser
	size   int64
}

func (r *bodysizeReader) Read(b []byte) (int, error) {
	n, err := r.reader.Read(b)
	r.size += int64(n)

	return n, err
}

func (r *bodysizeReader) Close() error {
	return r.reader.Close()
}

type sessionRewriter struct {
	http.ResponseWriter
	buffer bytes.Buffer
}

func (g *sessionRewriter) Write(data []byte) (int, error) {
	// Write the data into internal buffer for later rewrite
	w, err := g.buffer.Write(data)

	return w, err
}

func (g *sessionRewriter) rewriteHLS(sessionID string, requestURL *url.URL) {
	var buffer bytes.Buffer

	isMaster := false

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

		u, err := url.Parse(line)
		if err != nil {
			buffer.WriteString(line + "\n")
			continue
		}

		// Write anything that doesn't end in .m3u8 or .ts unmodified
		if !strings.HasSuffix(u.Path, ".m3u8") && !strings.HasSuffix(u.Path, ".ts") {
			buffer.WriteString(line + "\n")
			continue
		}

		q := u.Query()

		loop := false

		// If this is a master manifest (i.e. an m3u8 which contains references to other m3u8), then
		// we give each substream an own session ID if they don't have already.
		if strings.HasSuffix(u.Path, ".m3u8") {
			// Check if we're referring to ourselves. This will cause an infinite loop
			// and has to be stopped.
			file := u.Path
			if !strings.HasPrefix(file, "/") {
				dir := urlpath.Dir(requestURL.Path)
				file = filepath.Join(dir, file)
			}

			if requestURL.Path == file {
				loop = true
			}

			q.Set("session", shortuuid.New())

			isMaster = true
		} else {
			q.Set("session", sessionID)
		}

		u.RawQuery = q.Encode()

		if loop {
			buffer.WriteString("# m3u8 is referencing itself: " + u.String() + "\n")
		} else {
			buffer.WriteString(u.String() + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return
	}

	// If this is not a master manifest and there isn't a session ID, we add a new session ID.
	if !isMaster && len(sessionID) == 0 {
		sessionID = shortuuid.New()

		buffer.Reset()

		buffer.WriteString("#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=1024\n")

		// Add the session ID to the query string
		q := requestURL.Query()
		q.Set("session", sessionID)

		buffer.WriteString(urlpath.Base(requestURL.Path) + "?" + q.Encode())
	}

	g.buffer = buffer
}
