package session

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/net"

	"github.com/prep/average"
)

// Session represents an active session
type Session struct {
	ID           string
	Reference    string
	CreatedAt    time.Time
	Location     string
	Peer         string
	Extra        string
	RxBytes      uint64
	RxBitrate    float64 // bit/s
	TopRxBitrate float64 // bit/s
	TxBytes      uint64
	TxBitrate    float64 // bit/s
	TopTxBitrate float64 // bit/s
}

// Summary is a summary over all current and past sessions.
// The past sessions are grouped over the Peers/Locations and
// the Locations.
type Summary struct {
	MaxSessions  uint64
	MaxRxBitrate float64 // bit/s
	MaxTxBitrate float64 // bit/s

	CurrentSessions  uint64
	CurrentRxBitrate float64 // bit/s
	CurrentTxBitrate float64 // bit/s

	Active []Session

	Summary struct {
		Peers      map[string]Peers
		Locations  map[string]Stats
		References map[string]Stats
		Stats
	}
}

// Peers is a group of the same peer grouped of the locations.
type Peers struct {
	Stats

	Locations map[string]Stats
}

// Stats holds the basic accumulated values like the number of sessions,
// total transmitted and received bytes.
type Stats struct {
	TotalSessions uint64
	TotalRxBytes  uint64
	TotalTxBytes  uint64
}

// The Collector interface
type Collector interface {
	// Register registers a new session. A session has to be activated in order
	// not to be dropped. A different id distinguishes different sessions.
	Register(id, reference, location, peer string)

	// Activate activates the session with the id. Returns true if the session
	// has been activated, false if the session was already activated.
	Activate(id string) bool

	// RegisterAndActivate registers and an activates a session.
	RegisterAndActivate(id, reference, location, peer string)

	// Add arbitrary extra data to a session
	Extra(id, extra string)

	// Unregister cancels a session prematurely.
	Unregister(id string)

	// Ingress adds size bytes of ingress traffic to a session.
	Ingress(id string, size int64)

	// Egress adds size bytes of egress traffic to a session.
	Egress(id string, size int64)

	// IngressBitrate returns the current bitrate of ingress traffic.
	IngressBitrate() float64

	// EgressBitrate returns the current bitrate of egress traffic.
	EgressBitrate() float64

	// MaxIngressBitrate return the defined maximum ingress bitrate. All values <= 0
	// mean no limit.
	MaxIngressBitrate() float64

	// MaxEgressBitrate return the defined maximum egress bitrate. All values <= 0
	// mean no limit.
	MaxEgressBitrate() float64

	// TopIngressBitrate returns the summed current top bitrates of all ingress sessions.
	TopIngressBitrate() float64

	// TopEgressBitrate returns the summed current top bitrates of all egress sessions.
	TopEgressBitrate() float64

	// IsIngressBitrateExceeded returns whether the defined maximum ingress bitrate has
	// been exceeded.
	IsIngressBitrateExceeded() bool

	// IsEgressBitrateExceeded returns whether the defined maximum egress bitrate has
	// been exceeded.
	IsEgressBitrateExceeded() bool

	// IsSessionsExceeded return whether the maximum number of session have been exceeded.
	IsSessionsExceeded() bool

	// IsKnowsession returns whether a session with the given id exists.
	IsKnownSession(id string) bool

	// IsAllowedIP returns whether traffic from/to the given IP should be considered.
	IsCollectableIP(ip string) bool

	// Summary returns the summary of all currently active sessions and the session history.
	Summary() Summary

	// Active returns a list of currently active sessions.
	Active() []Session

	// SessionIngressTopBitrate returns the top ingress bitrate of a specific session.
	SessionTopIngressBitrate(id string) float64

	// SessionIngressTopBitrate returns the top egress bitrate of a specific session.
	SessionTopEgressBitrate(id string) float64

	// SessionSetIngressTopBitrate sets the current top ingress bitrate of a session.
	SessionSetTopIngressBitrate(id string, bitrate float64)

	// SessionSetEgressTopBitrate sets the current top egress bitrate of a session.
	SessionSetTopEgressBitrate(id string, bitrate float64)

	// Sessions returns the number of currently active sessions.
	Sessions() uint64

	AddCompanion(collector Collector)

	// IngressBitrate returns the current bitrate of ingress traffic.
	CompanionIngressBitrate() float64

	// EgressBitrate returns the current bitrate of egress traffic.
	CompanionEgressBitrate() float64

	// TopIngressBitrate returns the summed current top bitrates of all ingress sessions.
	CompanionTopIngressBitrate() float64

	// TopEgressBitrate returns the summed current top bitrates of all egress sessions.
	CompanionTopEgressBitrate() float64

	// Stop stops the collector to calculate rates
	Stop()
}

// CollectorConfig is the configuration for registering a new collector
type CollectorConfig struct {
	// MaxRxBitrate is the maximum ingress bitrate. It is used to query whether
	// the maximum bitrate is reached, based on the actucal bitrate.
	MaxRxBitrate uint64

	// MaxTxBitrate is the maximum egress bitrate. It is used to query whether
	// the maximum bitrate is reached, based on the actucal bitrate.
	MaxTxBitrate uint64

	// MaxSessions is the maximum number of session. It is used to query whether
	// the maximum number of sessions is reached, based on the actual number
	// of pending and active sessions.
	MaxSessions uint64

	// Limiter is an IPLimiter. It is used to query whether a session for an IP
	// should be created.
	Limiter net.IPLimiter

	// InactiveTimeout is the duration of how long a not yet activated session is kept.
	// A session gets activated with the first ingress or egress bytes.
	InactiveTimeout time.Duration

	// SessionTimeout is the duration of how long an idle active session is kept. A
	// session is idle if there are no ingress or egress bytes.
	SessionTimeout time.Duration

	// PersistInterval is the duration between persisting the
	// history. Can be 0. Then the history will only be persisted
	// at stopping the collector.
	PersistInterval time.Duration
}

type totals struct {
	Location      string `json:"location"`
	Peer          string `json:"peer"`
	Reference     string `json:"reference"`
	TotalSessions uint64 `json:"total_sessions"`
	TotalRxBytes  uint64 `json:"total_rxbytes"`
	TotalTxBytes  uint64 `json:"total_txbytes"`
}

type history struct {
	Sessions map[string]totals `json:"sessions"` // key = `${session.location}:${session.peer}`
}

type collector struct {
	id     string
	logger log.Logger

	sessions    map[string]*session
	sessionPool sync.Pool
	sessionsWG  sync.WaitGroup

	staleCallback func(*session)

	currentPendingSessions uint64
	currentActiveSessions  uint64

	totalSessions uint64
	rxBytes       uint64
	txBytes       uint64

	maxRxBitrate float64
	maxTxBitrate float64
	maxSessions  uint64

	rxBitrate *average.SlidingWindow
	txBitrate *average.SlidingWindow

	history history

	persist struct {
		enable   bool
		fs       fs.Filesystem
		path     string
		interval time.Duration
		done     context.CancelFunc
	}

	inactiveTimeout time.Duration
	sessionTimeout  time.Duration

	limiter net.IPLimiter

	companions []Collector

	lock struct {
		session   sync.RWMutex
		history   sync.RWMutex
		persist   sync.Mutex
		companion sync.RWMutex
	}

	startOnce sync.Once
	stopOnce  sync.Once
}

const (
	averageWindow      = 10 * time.Second
	averageGranularity = time.Second
)

// NewCollector returns a new collector according to the provided configuration. If such a
// collector can't be created, a NullCollector is returned.
func NewCollector(config CollectorConfig) Collector {
	collector, err := newCollector("", nil, nil, config)
	if err != nil {
		return NewNullCollector()
	}

	collector.start()

	return collector
}

func newCollector(id string, persistFS fs.Filesystem, logger log.Logger, config CollectorConfig) (*collector, error) {
	c := &collector{
		maxRxBitrate:    float64(config.MaxRxBitrate),
		maxTxBitrate:    float64(config.MaxTxBitrate),
		maxSessions:     config.MaxSessions,
		inactiveTimeout: config.InactiveTimeout,
		sessionTimeout:  config.SessionTimeout,
		limiter:         config.Limiter,
		logger:          logger,
		id:              id,
	}

	if c.logger == nil {
		c.logger = log.New("Session")
	}

	if c.limiter == nil {
		c.limiter, _ = net.NewIPLimiter(nil, nil)
	}

	if c.sessionTimeout <= 0 {
		c.sessionTimeout = 5 * time.Second
	}

	if c.inactiveTimeout <= 0 {
		c.inactiveTimeout = c.sessionTimeout
	}

	c.sessionPool = sync.Pool{
		New: func() interface{} {
			return &session{}
		},
	}

	c.staleCallback = func(sess *session) {
		defer func() {
			c.sessionsWG.Done()
		}()

		c.lock.session.Lock()
		defer c.lock.session.Unlock()

		delete(c.sessions, sess.id)

		if !sess.active {
			c.currentPendingSessions--

			sess.logger.Debug().Log("Closed pending")

			return
		}

		logger = sess.logger.WithFields(log.Fields{
			"rx_bytes":           sess.rxBytes,
			"rx_bitrate_kbit":    sess.RxBitrate() / 1024,
			"rx_maxbitrate_kbit": sess.MaxRxBitrate() / 1024,
			"tx_bytes":           sess.txBytes,
			"tx_bitrate_kbit":    sess.TxBitrate() / 1024,
			"tx_maxbitrate_kbit": sess.MaxTxBitrate() / 1024,
		})

		// Only log session that have been active
		logger.Info().Log("Closed")
		logger.Debug().Log("Closed")

		c.lock.history.Lock()

		key := sess.location + ":" + sess.peer + ":" + sess.reference

		// Update history totals per key
		t, ok := c.history.Sessions[key]
		t.TotalSessions++
		t.TotalRxBytes += sess.rxBytes
		t.TotalTxBytes += sess.txBytes

		if !ok {
			t.Location = sess.location
			t.Peer = sess.peer
			t.Reference = sess.reference
		}

		c.history.Sessions[key] = t

		c.lock.history.Unlock()

		c.sessionPool.Put(sess)

		c.currentActiveSessions--
	}

	c.sessions = make(map[string]*session)

	c.history.Sessions = make(map[string]totals)

	c.persist.enable = persistFS != nil
	c.persist.fs = persistFS
	c.persist.path = "/" + id + ".json"
	c.persist.interval = config.PersistInterval

	c.loadHistory(c.persist.fs, c.persist.path, &c.history)

	c.stopOnce.Do(func() {})

	c.start()

	return c, nil
}

func (c *collector) start() {
	c.startOnce.Do(func() {
		if c.persist.enable && c.persist.interval > 0 {
			ctx, cancel := context.WithCancel(context.Background())
			c.persist.done = cancel
			go c.persister(ctx, c.persist.interval)
		}

		c.rxBitrate, _ = average.New(averageWindow, averageGranularity)
		c.txBitrate, _ = average.New(averageWindow, averageGranularity)

		c.stopOnce = sync.Once{}
	})
}

func (c *collector) Stop() {
	c.stopOnce.Do(func() {
		if c.persist.enable && c.persist.interval > 0 {
			c.persist.done()
		}

		c.lock.session.RLock()
		for _, sess := range c.sessions {
			// Cancel all current sessions
			sess.Cancel()
		}
		c.lock.session.RUnlock()

		// Wait for all current sessions to finish
		c.sessionsWG.Wait()

		c.Persist()

		c.startOnce = sync.Once{}
	})
}

func (c *collector) Persist() {
	c.lock.history.RLock()
	defer c.lock.history.RUnlock()

	c.saveHistory(c.persist.fs, c.persist.path, &c.history)
}

func (c *collector) persister(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.Persist()
		}
	}
}

func (c *collector) loadHistory(fs fs.Filesystem, path string, data *history) {
	if fs == nil {
		return
	}

	c.logger.WithComponent("SessionStore").WithFields(log.Fields{
		"base": fs.Metadata("base"),
		"path": path,
	}).Debug().Log("Loading history")

	c.lock.persist.Lock()
	defer c.lock.persist.Unlock()

	jsondata, err := fs.ReadFile(path)
	if err != nil {
		return
	}

	if err = json.Unmarshal(jsondata, data); err != nil {
		return
	}
}

func (c *collector) saveHistory(fs fs.Filesystem, path string, data *history) {
	if fs == nil {
		return
	}

	c.logger.WithComponent("SessionStore").WithFields(log.Fields{
		"base": fs.Metadata("base"),
		"path": path,
	}).Debug().Log("Storing history")

	c.lock.persist.Lock()
	defer c.lock.persist.Unlock()

	jsondata, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return
	}

	_, _, err = fs.WriteFileSafe(path, jsondata)
	if err != nil {
		return
	}
}

func (c *collector) IsCollectableIP(ip string) bool {
	return c.limiter.IsAllowed(ip)
}

func (c *collector) IsKnownSession(id string) bool {
	c.lock.session.RLock()
	_, ok := c.sessions[id]
	c.lock.session.RUnlock()

	return ok
}

func (c *collector) RegisterAndActivate(id, reference, location, peer string) {
	c.Register(id, reference, location, peer)
	c.Activate(id)
}

func (c *collector) Register(id, reference, location, peer string) {
	c.lock.session.Lock()
	defer c.lock.session.Unlock()

	_, ok := c.sessions[id]
	if ok {
		return
	}

	logger := c.logger.WithFields(log.Fields{
		"id":        id,
		"type":      c.id,
		"location":  location,
		"peer":      peer,
		"reference": reference,
	})

	c.sessionsWG.Add(1)
	sess := c.sessionPool.Get().(*session)
	sess.Init(id, reference, c.staleCallback, c.inactiveTimeout, c.sessionTimeout, logger)

	logger.Debug().Log("Pending")

	c.currentPendingSessions++

	sess.Register(location, peer)

	c.sessions[id] = sess
}

func (c *collector) Unregister(id string) {
	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	sess, ok := c.sessions[id]
	if ok {
		sess.Cancel()
	}
}

func (c *collector) Activate(id string) bool {
	if len(id) == 0 {
		return false
	}

	c.lock.session.RLock()
	sess, ok := c.sessions[id]
	c.lock.session.RUnlock()

	if !ok {
		return false
	}

	if sess.Activate() {
		c.currentPendingSessions--
		c.currentActiveSessions++
		c.totalSessions++

		sess.logger.Info().Log("Active")
		sess.logger.Debug().Log("Active")

		return true
	}

	return false
}

func (c *collector) Extra(id, extra string) {
	c.lock.session.RLock()
	sess, ok := c.sessions[id]
	c.lock.session.RUnlock()

	if !ok {
		return
	}

	sess.Extra(extra)
}

func (c *collector) Ingress(id string, size int64) {
	if len(id) == 0 {
		return
	}

	c.lock.session.RLock()
	sess, ok := c.sessions[id]
	c.lock.session.RUnlock()

	if !ok {
		return
	}

	if sess.Ingress(size) {
		c.rxBitrate.Add(size * 8)
		c.rxBytes += uint64(size)
	}
}

func (c *collector) Egress(id string, size int64) {
	if len(id) == 0 {
		return
	}

	c.lock.session.RLock()
	sess, ok := c.sessions[id]
	c.lock.session.RUnlock()

	if !ok {
		return
	}

	if sess.Egress(size) {
		c.txBitrate.Add(size * 8)
		c.txBytes += uint64(size)
	}
}

func (c *collector) IsIngressBitrateExceeded() bool {
	if c.maxRxBitrate <= 0 {
		return false
	}

	if c.IngressBitrate() > c.maxRxBitrate {
		return true
	}

	return false
}

func (c *collector) IsEgressBitrateExceeded() bool {
	if c.maxTxBitrate <= 0 {
		return false
	}

	if c.EgressBitrate() > c.maxTxBitrate {
		return true
	}

	return false
}

func (c *collector) IsSessionsExceeded() bool {
	if c.maxSessions <= 0 {
		return false
	}

	if c.Sessions() >= c.maxSessions {
		return true
	}

	return false
}

func (c *collector) IngressBitrate() float64 {
	return c.rxBitrate.Average(averageWindow)
}

func (c *collector) EgressBitrate() float64 {
	return c.txBitrate.Average(averageWindow)
}

func (c *collector) MaxIngressBitrate() float64 {
	return c.maxRxBitrate
}

func (c *collector) MaxEgressBitrate() float64 {
	return c.maxTxBitrate
}

func (c *collector) TopIngressBitrate() float64 {
	var bitrate float64 = 0

	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	for _, sess := range c.sessions {
		if !sess.active {
			continue
		}

		bitrate += sess.TopRxBitrate()
	}

	return bitrate
}

func (c *collector) TopEgressBitrate() float64 {
	var bitrate float64 = 0

	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	for _, sess := range c.sessions {
		if !sess.active {
			continue
		}

		bitrate += sess.TopTxBitrate()
	}

	return bitrate
}

func (c *collector) SessionTopIngressBitrate(id string) float64 {
	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	if sess, ok := c.sessions[id]; ok {
		return sess.TopRxBitrate()
	}

	return 0.0
}

func (c *collector) SessionTopEgressBitrate(id string) float64 {
	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	if sess, ok := c.sessions[id]; ok {
		return sess.TopTxBitrate()
	}

	return 0.0
}

func (c *collector) SessionSetTopIngressBitrate(id string, bitrate float64) {
	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	if sess, ok := c.sessions[id]; ok {
		sess.SetTopRxBitrate(bitrate)
	}
}

func (c *collector) SessionSetTopEgressBitrate(id string, bitrate float64) {
	c.lock.session.RLock()
	defer c.lock.session.RUnlock()

	if sess, ok := c.sessions[id]; ok {
		sess.SetTopTxBitrate(bitrate)
	}
}

func (c *collector) Sessions() uint64 {
	return c.currentActiveSessions
}

func (c *collector) Summary() Summary {
	summary := Summary{
		MaxSessions:  c.maxSessions,
		MaxRxBitrate: c.maxRxBitrate,
		MaxTxBitrate: c.maxTxBitrate,
	}

	summary.CurrentSessions = c.currentActiveSessions
	summary.CurrentRxBitrate = c.IngressBitrate()
	summary.CurrentTxBitrate = c.EgressBitrate()

	summary.Summary.Peers = make(map[string]Peers)
	summary.Summary.Locations = make(map[string]Stats)
	summary.Summary.References = make(map[string]Stats)

	c.lock.history.RLock()

	for _, v := range c.history.Sessions {
		p := summary.Summary.Peers[v.Peer]

		p.TotalSessions += v.TotalSessions
		p.TotalRxBytes += v.TotalRxBytes
		p.TotalTxBytes += v.TotalTxBytes

		if p.Locations == nil {
			p.Locations = make(map[string]Stats)
		}

		stats := p.Locations[v.Location]

		stats.TotalSessions += v.TotalSessions
		stats.TotalRxBytes += v.TotalRxBytes
		stats.TotalTxBytes += v.TotalTxBytes

		p.Locations[v.Location] = stats

		summary.Summary.Peers[v.Peer] = p

		l := summary.Summary.Locations[v.Location]

		l.TotalSessions += v.TotalSessions
		l.TotalRxBytes += v.TotalRxBytes
		l.TotalTxBytes += v.TotalTxBytes

		summary.Summary.Locations[v.Location] = l

		r := summary.Summary.References[v.Reference]

		r.TotalSessions += v.TotalSessions
		r.TotalRxBytes += v.TotalRxBytes
		r.TotalTxBytes += v.TotalTxBytes

		summary.Summary.References[v.Reference] = r

		summary.Summary.TotalSessions += v.TotalSessions
		summary.Summary.TotalRxBytes += v.TotalRxBytes
		summary.Summary.TotalTxBytes += v.TotalTxBytes
	}

	c.lock.history.RUnlock()

	summary.Active = c.Active()

	return summary
}

func (c *collector) Active() []Session {
	sessions := []Session{}

	c.lock.session.RLock()
	for _, sess := range c.sessions {
		if !sess.active {
			continue
		}

		session := Session{
			ID:           sess.id,
			Reference:    sess.reference,
			CreatedAt:    sess.createdAt,
			Location:     sess.location,
			Peer:         sess.peer,
			Extra:        sess.extra,
			RxBytes:      sess.rxBytes,
			RxBitrate:    sess.RxBitrate(),
			TopRxBitrate: sess.TopRxBitrate(),
			TxBytes:      sess.txBytes,
			TxBitrate:    sess.TxBitrate(),
			TopTxBitrate: sess.TopTxBitrate(),
		}

		sessions = append(sessions, session)
	}
	c.lock.session.RUnlock()

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.Before(sessions[j].CreatedAt)
	})

	return sessions
}

func (c *collector) AddCompanion(collector Collector) {
	c.lock.companion.Lock()
	c.companions = append(c.companions, collector)
	c.lock.companion.Unlock()
}

func (c *collector) CompanionIngressBitrate() float64 {
	bitrate := c.IngressBitrate()

	c.lock.companion.RLock()
	for _, co := range c.companions {
		bitrate += co.IngressBitrate()
	}
	c.lock.companion.RUnlock()

	return bitrate
}

func (c *collector) CompanionEgressBitrate() float64 {
	bitrate := c.EgressBitrate()

	c.lock.companion.RLock()
	for _, co := range c.companions {
		bitrate += co.EgressBitrate()
	}
	c.lock.companion.RUnlock()

	return bitrate
}

func (c *collector) CompanionTopIngressBitrate() float64 {
	bitrate := c.TopIngressBitrate()

	c.lock.companion.RLock()
	for _, co := range c.companions {
		bitrate += co.TopIngressBitrate()
	}
	c.lock.companion.RUnlock()

	return bitrate
}

func (c *collector) CompanionTopEgressBitrate() float64 {
	bitrate := c.TopEgressBitrate()

	c.lock.companion.RLock()
	for _, co := range c.companions {
		bitrate += co.TopEgressBitrate()
	}
	c.lock.companion.RUnlock()

	return bitrate
}

type nullCollector struct{}

// NewNullCollector returns an implementation of the Collector interface that
// doesn't collect any metrics at all.
func NewNullCollector() Collector                                                 { return &nullCollector{} }
func (n *nullCollector) Register(id, reference, location, peer string)            {}
func (n *nullCollector) Activate(id string) bool                                  { return false }
func (n *nullCollector) RegisterAndActivate(id, reference, location, peer string) {}
func (n *nullCollector) Extra(id, extra string)                                   {}
func (n *nullCollector) Unregister(id string)                                     {}
func (n *nullCollector) Ingress(id string, size int64)                            {}
func (n *nullCollector) Egress(id string, size int64)                             {}
func (n *nullCollector) IngressBitrate() float64                                  { return 0.0 }
func (n *nullCollector) EgressBitrate() float64                                   { return 0.0 }
func (n *nullCollector) MaxIngressBitrate() float64                               { return 0.0 }
func (n *nullCollector) MaxEgressBitrate() float64                                { return 0.0 }
func (n *nullCollector) TopIngressBitrate() float64                               { return 0.0 }
func (n *nullCollector) TopEgressBitrate() float64                                { return 0.0 }
func (n *nullCollector) IsIngressBitrateExceeded() bool                           { return false }
func (n *nullCollector) IsEgressBitrateExceeded() bool                            { return false }
func (n *nullCollector) IsSessionsExceeded() bool                                 { return false }
func (n *nullCollector) IsKnownSession(id string) bool                            { return false }
func (n *nullCollector) IsCollectableIP(ip string) bool                           { return true }
func (n *nullCollector) Summary() Summary                                         { return Summary{} }
func (n *nullCollector) Active() []Session                                        { return []Session{} }
func (n *nullCollector) SessionTopIngressBitrate(id string) float64               { return 0.0 }
func (n *nullCollector) SessionTopEgressBitrate(id string) float64                { return 0.0 }
func (n *nullCollector) SessionSetTopIngressBitrate(id string, bitrate float64)   {}
func (n *nullCollector) SessionSetTopEgressBitrate(id string, bitrate float64)    {}
func (n *nullCollector) Sessions() uint64                                         { return 0 }
func (n *nullCollector) AddCompanion(collector Collector)                         {}
func (n *nullCollector) CompanionIngressBitrate() float64                         { return 0.0 }
func (n *nullCollector) CompanionEgressBitrate() float64                          { return 0.0 }
func (n *nullCollector) CompanionTopIngressBitrate() float64                      { return 0.0 }
func (n *nullCollector) CompanionTopEgressBitrate() float64                       { return 0.0 }
func (n *nullCollector) Stop()                                                    {}
