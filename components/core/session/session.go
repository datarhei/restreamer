package session

import (
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
	"github.com/prep/average"
)

type session struct {
	id        string
	reference string
	createdAt time.Time

	logger log.Logger

	sessionActivate sync.Mutex
	active          bool

	location string
	peer     string
	extra    string

	stale    *time.Timer
	timeout  time.Duration
	callback func(*session)

	rxBitrate *average.SlidingWindow
	rxBytes   uint64

	txBitrate *average.SlidingWindow
	txBytes   uint64

	tickerStop   chan struct{}
	sessionClose sync.Once

	lock         sync.Mutex
	topRxBitrate float64
	topTxBitrate float64
	maxRxBitrate float64
	maxTxBitrate float64
}

func (s *session) Init(id, reference string, closeCallback func(*session), inactive, timeout time.Duration, logger log.Logger) {
	s.id = id
	s.reference = reference
	s.createdAt = time.Now()

	s.logger = logger

	s.active = false

	s.location = ""
	s.peer = ""

	s.rxBitrate, _ = average.New(averageWindow, averageGranularity)
	s.txBitrate, _ = average.New(averageWindow, averageGranularity)

	s.topRxBitrate = 0.0
	s.topTxBitrate = 0.0

	s.maxRxBitrate = 0.0
	s.maxTxBitrate = 0.0

	s.timeout = timeout
	s.callback = closeCallback
	if s.callback == nil {
		s.callback = func(s *session) {}
	}

	s.tickerStop = make(chan struct{})
	s.sessionClose = sync.Once{}

	pendingTimeout := inactive
	if timeout < pendingTimeout {
		pendingTimeout = timeout
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.stale = time.AfterFunc(pendingTimeout, func() {
		s.close()
	})
}

func (s *session) close() {
	s.sessionClose.Do(func() {
		s.lock.Lock()
		s.stale.Stop()
		s.lock.Unlock()

		close(s.tickerStop)
		s.rxBitrate.Stop()
		s.txBitrate.Stop()
		go s.callback(s)
	})
}

func (s *session) Register(location, peer string) {
	if len(location) == 0 {
		location = "unknown"
	}

	if len(peer) == 0 {
		peer = "unknown"
	}

	s.location = location
	s.peer = peer
}

func (s *session) Activate() bool {
	s.sessionActivate.Lock()
	defer s.sessionActivate.Unlock()

	if s.active {
		return false
	}

	s.active = true
	go s.ticker()

	return true
}

func (s *session) Extra(extra string) {
	s.extra = extra

	s.logger = s.logger.WithField("extra", extra)
}

func (s *session) Ingress(size int64) bool {
	if size == 0 {
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.stale.Stop()
	s.stale.Reset(s.timeout)

	s.rxBitrate.Add(size * 8)
	s.rxBytes += uint64(size)

	bitrate := s.rxBitrate.Average(averageWindow)
	if bitrate > s.topRxBitrate {
		s.topRxBitrate = bitrate
	}

	if bitrate > s.maxRxBitrate {
		s.maxRxBitrate = bitrate
	}

	return true
}

func (s *session) Egress(size int64) bool {
	if size == 0 {
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.stale.Stop()
	s.stale.Reset(s.timeout)

	s.txBitrate.Add(size * 8)
	s.txBytes += uint64(size)

	bitrate := s.txBitrate.Average(averageWindow)
	if bitrate > s.topTxBitrate {
		s.topTxBitrate = bitrate
	}

	if bitrate > s.maxTxBitrate {
		s.maxTxBitrate = bitrate
	}

	return true
}

func (s *session) RxBitrate() float64 {
	return s.rxBitrate.Average(averageWindow)
}

func (s *session) TxBitrate() float64 {
	return s.txBitrate.Average(averageWindow)
}

func (s *session) TopRxBitrate() float64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.topRxBitrate
}

func (s *session) MaxTxBitrate() float64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.maxTxBitrate
}

func (s *session) MaxRxBitrate() float64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.maxRxBitrate
}

func (s *session) TopTxBitrate() float64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.topTxBitrate
}

func (s *session) SetTopRxBitrate(bitrate float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.topRxBitrate = bitrate
}

func (s *session) SetTopTxBitrate(bitrate float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.topTxBitrate = bitrate
}

func (s *session) ticker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.tickerStop:
			return
		case <-ticker.C:
			s.lock.Lock()
			s.topRxBitrate *= 0.95
			s.topTxBitrate *= 0.95
			s.lock.Unlock()
		}
	}
}

func (s *session) Cancel() {
	s.close()
}
