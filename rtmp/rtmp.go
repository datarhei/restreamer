// Package rtmp provides an RTMP server
package rtmp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/session"

	"github.com/datarhei/joy4/av/avutil"
	"github.com/datarhei/joy4/av/pktque"
	"github.com/datarhei/joy4/av/pubsub"
	"github.com/datarhei/joy4/format"
	"github.com/datarhei/joy4/format/flv/flvio"
	"github.com/datarhei/joy4/format/rtmp"
)

// ErrServerClosed is returned by ListenAndServe (or ListenAndServeTLS) if the server
// has been closed regularly with the Close() function.
var ErrServerClosed = rtmp.ErrServerClosed

func init() {
	format.RegisterAll()
}

type client struct {
	conn      *rtmp.Conn
	id        string
	createdAt time.Time

	txbytes uint64
	rxbytes uint64

	collector session.Collector

	cancel context.CancelFunc
}

func newClient(conn *rtmp.Conn, id string, collector session.Collector) *client {
	c := &client{
		conn:      conn,
		id:        id,
		createdAt: time.Now(),

		collector: collector,
	}

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

	go c.ticker(ctx)

	return c
}

func (c *client) ticker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			txbytes := c.conn.TxBytes()
			rxbytes := c.conn.RxBytes()

			c.collector.Ingress(c.id, int64(rxbytes-c.rxbytes))
			c.collector.Egress(c.id, int64(txbytes-c.txbytes))

			c.txbytes = txbytes
			c.rxbytes = rxbytes
		}
	}
}

func (c *client) Close() {
	c.cancel()
}

// channel represents a stream that is sent to the server
type channel struct {
	// The packet queue for the stream
	queue *pubsub.Queue

	// The metadata of the stream
	metadata flvio.AMFMap

	// Whether the stream has an audio track
	hasAudio bool

	// Whether the stream has a video track
	hasVideo bool

	collector session.Collector
	path      string
	reference string

	publisher  *client
	subscriber map[string]*client
	lock       sync.RWMutex
}

func newChannel(conn *rtmp.Conn, reference string, collector session.Collector) *channel {
	ch := &channel{
		path:       conn.URL.Path,
		reference:  reference,
		publisher:  newClient(conn, conn.URL.Path, collector),
		subscriber: make(map[string]*client),
		collector:  collector,
	}

	addr := conn.NetConn().RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(addr)

	if collector.IsCollectableIP(ip) {
		collector.RegisterAndActivate(ch.path, ch.reference, "publish:"+ch.path, addr)
	}

	return ch
}

func (ch *channel) Close() {
	if ch.publisher == nil {
		return
	}

	ch.publisher.Close()
	ch.publisher = nil

	ch.queue.Close()
}

func (ch *channel) AddSubscriber(conn *rtmp.Conn) string {
	addr := conn.NetConn().RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(addr)

	client := newClient(conn, addr, ch.collector)

	if ch.collector.IsCollectableIP(ip) {
		ch.collector.RegisterAndActivate(addr, ch.reference, "play:"+ch.path, addr)
	}

	ch.lock.Lock()
	ch.subscriber[addr] = client
	ch.lock.Unlock()

	return addr
}

func (ch *channel) RemoveSubscriber(id string) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	client := ch.subscriber[id]
	if client != nil {
		delete(ch.subscriber, id)
		client.Close()
	}
}

// Config for a new RTMP server
type Config struct {
	// Logger. Optional.
	Logger log.Logger

	Collector session.Collector

	// The address the RTMP server should listen on, e.g. ":1935"
	Addr string

	// The address the RTMPS server should listen on, e.g. ":1936"
	TLSAddr string

	// The app path for the streams, e.g. "/live". Optional. Defaults
	// to "/".
	App string

	// A token that needs to be added to the URL as query string
	// in order push or pull a stream. The key for the query
	// parameter is "token". Optional. By default no token is
	// required.
	Token string

	// TLSConfig optionally provides a TLS configuration for use
	// by ListenAndServe. Note that this value is cloned by
	// ListenAndServe, so it's not possible to modify the configuration
	// with methods like tls.Config.SetSessionTicketKeys.
	TLSConfig *tls.Config

	// ConnectionIdleTimeout is the timeout in seconds after which idle
	// connection will be closes. Default is no timeout.
	ConnectionIdleTimeout time.Duration
}

// Server represents a RTMP server
type Server interface {
	// ListenAndServe starts the RTMP server
	ListenAndServe() error

	// ListenAndServe starts the RTMPS server
	ListenAndServeTLS(certFile, keyFile string) error

	// Close stops the RTMP server and closes all connections
	Close()

	// Channels return a list of currently publishing streams
	Channels() []string
}

// server is an implementation of the Server interface
type server struct {
	// Configuration parameter taken from the Config
	app       string
	token     string
	logger    log.Logger
	collector session.Collector

	// A joy4 RTMP server instance
	server    *rtmp.Server
	tlsServer *rtmp.Server

	// Map of publishing channels and a lock to serialize
	// access to the map.
	channels map[string]*channel
	lock     sync.RWMutex
}

// New creates a new RTMP server according to the given config
func New(config Config) (Server, error) {
	if len(config.App) == 0 {
		config.App = "/"
	}

	if config.Logger == nil {
		config.Logger = log.New("")
	}

	s := &server{
		app:       config.App,
		token:     config.Token,
		logger:    config.Logger,
		collector: config.Collector,
	}

	if s.collector == nil {
		s.collector = session.NewNullCollector()
	}

	s.server = &rtmp.Server{
		Addr:                  config.Addr,
		HandlePlay:            s.handlePlay,
		HandlePublish:         s.handlePublish,
		ConnectionIdleTimeout: config.ConnectionIdleTimeout,
	}

	if len(config.TLSAddr) != 0 {
		s.tlsServer = &rtmp.Server{
			Addr:                  config.TLSAddr,
			TLSConfig:             config.TLSConfig.Clone(),
			HandlePlay:            s.handlePlay,
			HandlePublish:         s.handlePublish,
			ConnectionIdleTimeout: config.ConnectionIdleTimeout,
		}
	}

	s.channels = make(map[string]*channel)

	rtmp.Debug = false

	return s, nil
}

// ListenAndServe starts the RMTP server
func (s *server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *server) ListenAndServeTLS(certFile, keyFile string) error {
	if s.tlsServer == nil {
		return fmt.Errorf("RTMPS server is not configured")
	}

	return s.tlsServer.ListenAndServeTLS(certFile, keyFile)
}

func (s *server) Close() {
	// Stop listening
	s.server.Close()

	if s.tlsServer != nil {
		s.tlsServer.Close()
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// Close all channels
	for _, ch := range s.channels {
		ch.Close()
	}
}

// Channels returns the list of streams that are
// publishing currently
func (s *server) Channels() []string {
	channels := []string{}

	s.lock.RLock()
	defer s.lock.RUnlock()

	for key := range s.channels {
		channels = append(channels, key)
	}

	return channels
}

func (s *server) log(level log.Level, who, action, path, message string, client net.Addr) {
	s.logger.WithLevel(level).WithFields(log.Fields{
		"who":    who,
		"action": action,
		"path":   path,
		"client": client.String(),
	}).Log(message)
}

// getToken returns the path and the token found in the URL. If the token
// was part of the path, the token is removed from the path. The token in
// the query string takes precedence. The token in the path is assumed to
// be the last path element.
func getToken(u *url.URL) (string, string) {
	q := u.Query()
	token := q.Get("token")

	if len(token) != 0 {
		// The token was in the query. Return the unmomdified path and the token
		return u.Path, token
	}

	pathElements := strings.Split(u.EscapedPath(), "/")
	nPathElements := len(pathElements)

	if nPathElements == 0 {
		return u.Path, ""
	}

	// Return the path without the token
	return strings.Join(pathElements[:nPathElements-1], "/"), pathElements[nPathElements-1]
}

// handlePlay is called when a RTMP client wants to play a stream
func (s *server) handlePlay(conn *rtmp.Conn) {
	client := conn.NetConn().RemoteAddr()

	defer conn.Close()

	playPath := conn.URL.Path

	// Check the token in the URL if one is required
	if len(s.token) != 0 {
		path, token := getToken(conn.URL)

		if len(token) == 0 {
			s.log(log.Lwarn, "PLAY", "FORBIDDEN", path, "no streamkey provided", client)
			return
		}

		if s.token != token {
			s.log(log.Lwarn, "PLAY", "FORBIDDEN", path, "invalid streamkey ("+token+")", client)
			return
		}

		playPath = path
	}

	/*
		ip, _, _ := net.SplitHostPort(client.String())
		if s.collector.IsCollectableIP(ip) {
			maxBitrate := s.collector.MaxEgressBitrate()
			if maxBitrate > 0.0 {
				streamBitrate := s.collector.SessionTopIngressBitrate(conn.URL.Path) * 2.0
				currentBitrate := s.collector.CompanionEgressBitrate() * 1.15

				resultingBitrate := currentBitrate + streamBitrate

				if resultingBitrate >= maxBitrate {
					s.log("PLAY", "FORBIDDEN", conn.URL.Path, "bandwidth limit exceeded", client)
					conn.Close()
					return
				}
			}
		}
	*/

	// Look for the stream
	s.lock.RLock()
	ch := s.channels[playPath]
	s.lock.RUnlock()

	if ch != nil {
		// Set the metadata for the client
		conn.SetMetaData(ch.metadata)

		s.log(log.Linfo, "PLAY", "START", playPath, "", client)

		// Get a cursor and apply filters
		cursor := ch.queue.Oldest()

		filters := pktque.Filters{}

		if ch.hasVideo {
			// The first packet has to be a key frame
			filters = append(filters, &pktque.WaitKeyFrame{})
		}

		// Adjust the timestamp such that the stream starts from 0
		filters = append(filters, &pktque.FixTime{StartFromZero: true, MakeIncrement: false})

		demuxer := &pktque.FilterDemuxer{
			Filter:  filters,
			Demuxer: cursor,
		}

		id := ch.AddSubscriber(conn)

		// Transfer the data
		err := avutil.CopyFile(conn, demuxer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.log(log.Lerror, "PLAY", "ERROR", playPath, err.Error(), client)
			}
		}

		ch.RemoveSubscriber(id)

		s.log(log.Linfo, "PLAY", "STOP", playPath, "", client)
	} else {
		s.log(log.Lwarn, "PLAY", "NOTFOUND", playPath, "", client)
	}
}

// handlePublish is called when a RTMP client wants to publish a stream
func (s *server) handlePublish(conn *rtmp.Conn) {
	client := conn.NetConn().RemoteAddr()

	defer conn.Close()

	playPath := conn.URL.Path

	if len(s.token) != 0 {
		path, token := getToken(conn.URL)

		if len(token) == 0 {
			s.log(log.Lwarn, "PUBLISH", "FORBIDDEN", path, "no streamkey provided", client)
			return
		}

		if s.token != token {
			s.log(log.Lwarn, "PUBLISH", "FORBIDDEN", path, "invalid streamkey ("+token+")", client)
			return
		}

		playPath = path
	}

	// Check the app patch
	if !strings.HasPrefix(playPath, s.app) {
		s.log(log.Lwarn, "PUBLISH", "FORBIDDEN", conn.URL.Path, "invalid app", client)
		return
	}

	// Check the stream if it contains any valid/known streams
	streams, _ := conn.Streams()

	if len(streams) == 0 {
		s.log(log.Lwarn, "PUBLISH", "INVALID", playPath, "no streams available", client)
		return
	}

	s.lock.Lock()

	ch := s.channels[conn.URL.Path]
	if ch == nil {
		reference := strings.TrimPrefix(strings.TrimSuffix(playPath, filepath.Ext(playPath)), s.app+"/")

		// Create a new channel
		ch = newChannel(conn, reference, s.collector)
		ch.metadata = conn.GetMetaData()
		ch.queue = pubsub.NewQueue()
		ch.queue.WriteHeader(streams)

		for _, stream := range streams {
			typ := stream.Type()

			switch {
			case typ.IsAudio():
				ch.hasAudio = true
			case typ.IsVideo():
				ch.hasVideo = true
			}
		}

		s.channels[playPath] = ch
	} else {
		ch = nil
	}

	s.lock.Unlock()

	if ch == nil {
		s.log(log.Lwarn, "PUBLISH", "CONFLICT", playPath, "already publishing", client)
		return
	}

	s.log(log.Linfo, "PUBLISH", "START", playPath, "", client)

	for _, stream := range streams {
		s.log(log.Linfo, "PUBLISH", "STREAM", playPath, stream.Type().String(), client)
	}

	// Ingest the data
	err := avutil.CopyPackets(ch.queue, conn)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			s.log(log.Lerror, "PUBLISH", "ERROR", playPath, err.Error(), client)
		}
	}

	s.lock.Lock()
	delete(s.channels, playPath)
	s.lock.Unlock()

	ch.Close()

	s.log(log.Linfo, "PUBLISH", "STOP", playPath, "", client)
}
