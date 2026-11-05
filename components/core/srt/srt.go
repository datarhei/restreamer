package srt

import (
	"container/ring"
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/session"
	srt "github.com/datarhei/gosrt"
)

// ErrServerClosed is returned by ListenAndServe if the server
// has been closed regularly with the Close() function.
var ErrServerClosed = srt.ErrServerClosed

type client struct {
	conn      srt.Conn
	id        string
	createdAt time.Time

	txbytes uint64
	rxbytes uint64

	collector session.Collector

	cancel context.CancelFunc
}

func newClient(conn srt.Conn, id string, collector session.Collector) *client {
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

	stats := &srt.Statistics{}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.conn.Stats(stats)

			rxbytes := stats.Accumulated.ByteRecv
			txbytes := stats.Accumulated.ByteSent

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
	pubsub    srt.PubSub
	collector session.Collector
	path      string

	publisher  *client
	subscriber map[string]*client
	lock       sync.RWMutex
}

func newChannel(conn srt.Conn, resource string, collector session.Collector) *channel {
	ch := &channel{
		pubsub:     srt.NewPubSub(srt.PubSubConfig{}),
		path:       resource,
		publisher:  newClient(conn, resource, collector),
		subscriber: make(map[string]*client),
		collector:  collector,
	}

	addr := conn.RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(addr)

	if collector.IsCollectableIP(ip) {
		collector.RegisterAndActivate(resource, resource, "publish:"+resource, addr)
	}

	return ch
}

func (ch *channel) Close() {
	if ch.publisher == nil {
		return
	}

	ch.publisher.Close()
	ch.publisher = nil
}

func (ch *channel) AddSubscriber(conn srt.Conn, resource string) string {
	addr := conn.RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(addr)

	client := newClient(conn, addr, ch.collector)

	if ch.collector.IsCollectableIP(ip) {
		ch.collector.RegisterAndActivate(addr, resource, "play:"+resource, addr)
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

// Config for a new SRT server
type Config struct {
	// The address the SRT server should listen on, e.g. ":1935"
	Addr string

	// A token that needs to be added to the URL as query string
	// in order push or pull a stream. The key for the query
	// parameter is "token". Optional. By default no token is
	// required.
	Token string

	Passphrase string

	// Logger. Optional.
	Logger log.Logger

	Collector session.Collector

	SRTLogTopics []string
}

// Server represents a SRT server
type Server interface {
	// ListenAndServe starts the SRT server
	ListenAndServe() error

	// Close stops the RTMP server and closes all connections
	Close()

	// Channels return a list of currently publishing streams
	Channels() Channels
}

// server implements the Server interface
type server struct {
	addr       string
	token      string
	passphrase string

	collector session.Collector

	server srt.Server

	// Map of publishing channels and a lock to serialize
	// access to the map.
	channels map[string]*channel
	lock     sync.RWMutex

	logger log.Logger

	srtlogger       srt.Logger
	srtloggerCancel context.CancelFunc
	srtlog          map[string]*ring.Ring
	srtlogLock      sync.RWMutex
}

func New(config Config) (Server, error) {
	s := &server{
		addr:       config.Addr,
		token:      config.Token,
		passphrase: config.Passphrase,
		collector:  config.Collector,
		logger:     config.Logger,
	}

	if s.collector == nil {
		s.collector = session.NewNullCollector()
	}

	if s.logger == nil {
		s.logger = log.New("")
	}

	s.srtlogger = srt.NewLogger(config.SRTLogTopics)

	s.srtlogLock.Lock()
	s.srtlog = make(map[string]*ring.Ring)
	s.srtlogLock.Unlock()

	s.channels = make(map[string]*channel)

	srtconfig := srt.DefaultConfig()

	srtconfig.Passphrase = config.Passphrase
	srtconfig.Logger = s.srtlogger

	s.server.Addr = config.Addr
	s.server.Config = &srtconfig
	s.server.HandleConnect = s.handleConnect
	s.server.HandlePublish = s.handlePublish
	s.server.HandleSubscribe = s.handleSubscribe

	return s, nil
}

func (s *server) ListenAndServe() error {
	var ctx context.Context
	ctx, s.srtloggerCancel = context.WithCancel(context.Background())

	go s.srtlogListener(ctx)

	return s.server.ListenAndServe()
}

func (s *server) Close() {
	s.server.Shutdown()

	s.srtloggerCancel()
}

type Log struct {
	Timestamp time.Time
	Message   []string
}

type Connection struct {
	Log   map[string][]Log
	Stats srt.Statistics
}

type Channels struct {
	Publisher   map[string]uint32
	Subscriber  map[string][]uint32
	Connections map[uint32]Connection
	Log         map[string][]Log
}

func (s *server) Channels() Channels {
	st := Channels{
		Publisher:   map[string]uint32{},
		Subscriber:  map[string][]uint32{},
		Connections: map[uint32]Connection{},
		Log:         map[string][]Log{},
	}

	s.lock.RLock()
	for id, ch := range s.channels {
		socketId := ch.publisher.conn.SocketId()
		st.Publisher[id] = socketId

		stats := &srt.Statistics{}
		ch.publisher.conn.Stats(stats)

		st.Connections[socketId] = Connection{
			Stats: *stats,
			Log:   map[string][]Log{},
		}

		for _, c := range ch.subscriber {
			socketId := c.conn.SocketId()
			st.Subscriber[id] = append(st.Subscriber[id], socketId)

			stats := &srt.Statistics{}
			c.conn.Stats(stats)

			st.Connections[socketId] = Connection{
				Stats: *stats,
				Log:   map[string][]Log{},
			}
		}
	}
	s.lock.RUnlock()

	s.srtlogLock.RLock()
	for topic, buf := range s.srtlog {

		buf.Do(func(l interface{}) {
			if l == nil {
				return
			}

			ll := l.(srt.Log)

			log := Log{
				Timestamp: ll.Time,
				Message:   strings.Split(ll.Message, "\n"),
			}

			if ll.SocketId != 0 {
				if _, ok := st.Connections[ll.SocketId]; ok {
					st.Connections[ll.SocketId].Log[topic] = append(st.Connections[ll.SocketId].Log[topic], log)
				}
			} else {
				st.Log[topic] = append(st.Log[topic], log)
			}
		})
	}
	s.srtlogLock.RUnlock()

	return st
}

func (s *server) srtlogListener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case l := <-s.srtlogger.Listen():
			s.srtlogLock.Lock()
			if buf := s.srtlog[l.Topic]; buf == nil {
				s.srtlog[l.Topic] = ring.New(100)
			}
			s.srtlog[l.Topic].Value = l
			s.srtlog[l.Topic] = s.srtlog[l.Topic].Next()
			s.srtlogLock.Unlock()
		}
	}
}

func (s *server) log(handler, action, resource, message string, client net.Addr) {
	s.logger.Info().WithFields(log.Fields{
		"handler":  handler,
		"status":   action,
		"resource": resource,
		"client":   client.String(),
	}).Log(message)
}

type streamInfo struct {
	mode     string
	resource string
	token    string
}

func parseStreamId(streamid string) (streamInfo, error) {
	si := streamInfo{}

	if strings.HasPrefix(streamid, "#!:") {
		return parseOldStreamId(streamid)
	}

	re := regexp.MustCompile(`,(token|mode):(.+)`)

	results := map[string]string{}

	idEnd := -1
	value := streamid
	key := ""

	for {
		matches := re.FindStringSubmatchIndex(value)
		if matches == nil {
			break
		}

		if idEnd < 0 {
			idEnd = matches[2] - 1
		}

		if len(key) != 0 {
			results[key] = value[:matches[2]-1]
		}

		key = value[matches[2]:matches[3]]
		value = value[matches[4]:matches[5]]

		results[key] = value
	}

	if idEnd < 0 {
		idEnd = len(streamid)
	}

	si.resource = streamid[:idEnd]
	if token, ok := results["token"]; ok {
		si.token = token
	}

	if mode, ok := results["mode"]; ok {
		si.mode = mode
	} else {
		si.mode = "request"
	}

	return si, nil
}

func parseOldStreamId(streamid string) (streamInfo, error) {
	si := streamInfo{}

	if !strings.HasPrefix(streamid, "#!:") {
		return si, fmt.Errorf("unknown streamid format")
	}

	streamid = strings.TrimPrefix(streamid, "#!:")

	kvs := strings.Split(streamid, ",")

	splitFn := func(s, sep string) (string, string, error) {
		splitted := strings.SplitN(s, sep, 2)

		if len(splitted) != 2 {
			return "", "", fmt.Errorf("invalid key/value pair")
		}

		return splitted[0], splitted[1], nil
	}

	for _, kv := range kvs {
		key, value, err := splitFn(kv, "=")
		if err != nil {
			continue
		}

		switch key {
		case "m":
			si.mode = value
		case "r":
			si.resource = value
		case "token":
			si.token = value
		default:
		}
	}

	return si, nil
}

func (s *server) handleConnect(req srt.ConnRequest) srt.ConnType {
	mode := srt.REJECT
	client := req.RemoteAddr()
	streamId := req.StreamId()

	var si streamInfo
	var err error

	if req.Version() == 4 {
		si.mode = "publish"
		si.resource = client.String()

		if len(s.passphrase) != 0 {
			req.SetPassphrase(s.passphrase)
		}
	} else if req.Version() == 5 {
		si, err = parseStreamId(streamId)
		if err != nil {
			s.log("CONNECT", "INVALID", "", err.Error(), client)
			return srt.REJECT
		}

		if len(si.resource) == 0 {
			s.log("CONNECT", "INVALID", "", "stream resource not provided", client)
			return srt.REJECT
		}

		if len(s.passphrase) != 0 {
			if !req.IsEncrypted() {
				s.log("CONNECT", "FORBIDDEN", si.resource, "connection has to be encrypted", client)
				return srt.REJECT
			}

			if err := req.SetPassphrase(s.passphrase); err != nil {
				s.log("CONNECT", "FORBIDDEN", si.resource, err.Error(), client)
				return srt.REJECT
			}
		} else {
			if req.IsEncrypted() {
				s.log("CONNECT", "INVALID", si.resource, "connection must not be encrypted", client)
				return srt.REJECT
			}
		}
	}

	if si.mode == "publish" {
		mode = srt.PUBLISH
	} else if si.mode == "request" {
		mode = srt.SUBSCRIBE
	} else {
		s.log("CONNECT", "INVALID", si.resource, "invalid connection mode", client)
		return srt.REJECT
	}

	// Check the token
	if len(s.token) != 0 && s.token != si.token {
		if len(si.token) == 0 {
			s.log("CONNECT", "FORBIDDEN", si.resource, "token required", client)
		} else {
			s.log("CONNECT", "FORBIDDEN", si.resource, "invalid token ("+si.token+")", client)
		}
		return srt.REJECT
	}

	s.lock.RLock()
	ch := s.channels[si.resource]
	s.lock.RUnlock()

	if mode == srt.PUBLISH && ch != nil {
		s.log("CONNECT", "CONFLICT", si.resource, "already publishing", client)
		return srt.REJECT
	}

	if mode == srt.SUBSCRIBE && ch == nil {
		s.log("CONNECT", "NOTFOUND", si.resource, "no publisher for this resource found", client)
		return srt.REJECT
	}

	return mode
}

func (s *server) handlePublish(conn srt.Conn) {
	streamId := conn.StreamId()
	client := conn.RemoteAddr()

	var si streamInfo

	if len(streamId) == 0 {
		si.resource = client.String()
	} else {
		si, _ = parseStreamId(streamId)
	}

	// Look for the stream
	s.lock.Lock()
	ch := s.channels[si.resource]
	if ch == nil {
		ch = newChannel(conn, si.resource, s.collector)
		s.channels[si.resource] = ch
	} else {
		ch = nil
	}
	s.lock.Unlock()

	if ch == nil {
		s.log("PUBLISH", "CONFLICT", si.resource, "already publishing", client)
		conn.Close()
		return
	}

	s.log("PUBLISH", "START", si.resource, "", client)

	ch.pubsub.Publish(conn)

	s.lock.Lock()
	delete(s.channels, si.resource)
	s.lock.Unlock()

	ch.Close()

	s.log("PUBLISH", "STOP", si.resource, "", client)

	conn.Close()
}

func (s *server) handleSubscribe(conn srt.Conn) {
	streamId := conn.StreamId()
	client := conn.RemoteAddr()

	si, _ := parseStreamId(streamId)

	// Look for the stream
	s.lock.RLock()
	ch := s.channels[si.resource]
	s.lock.RUnlock()

	if ch == nil {
		s.log("SUBSCRIBE", "NOTFOUND", si.resource, "no publisher for this resource found", client)
		conn.Close()
		return
	}

	s.log("SUBSCRIBE", "START", si.resource, "", client)

	id := ch.AddSubscriber(conn, si.resource)

	ch.pubsub.Subscribe(conn)

	s.log("SUBSCRIBE", "STOP", si.resource, "", client)

	ch.RemoveSubscriber(id)

	conn.Close()
}
