package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/monitor/metric"
	"github.com/datarhei/core/v16/service/api"
)

// Config is the configuration for the service
type Config struct {
	ID      string
	Version string
	Domain  string
	URL     string
	Token   string
	Monitor metric.Reader
	Logger  log.Logger
}

// Service is an interface
type Service interface {
	Start()
	Stop()
}

type service struct {
	id      string
	version string
	domain  string
	api     api.API
	monitor metric.Reader

	startOnce sync.Once
	stopOnce  sync.Once

	stopTicker context.CancelFunc

	logger log.Logger
}

// New creates a new service instance that implements the Service interface
func New(config Config) (Service, error) {
	s := &service{
		id:      config.ID,
		version: config.Version,
		domain:  config.Domain,
		monitor: config.Monitor,
		logger:  config.Logger,
	}

	if s.logger == nil {
		s.logger = log.New("")
	}

	s.logger = s.logger.WithField("url", config.URL)

	if s.monitor == nil {
		return nil, fmt.Errorf("no monitor provided")
	}

	a, err := api.New(api.Config{
		URL:   config.URL,
		Token: config.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service API: %w", err)
	}

	s.api = a

	// drain stop once, so it can't be called before startOnce has been called
	s.stopOnce.Do(func() {})

	return s, nil
}

func (s *service) tick(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	nerrors := 0
	next := time.Time{}

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			if t.After(next) {
				n, err := s.collect()
				if err != nil {
					nerrors++
					n = time.Duration(nerrors*10) * time.Second
					if nerrors > 3 {
						n = 15 * time.Minute
						s.logger.WithError(err).Warn().Log("Failed to send data")
					}
				} else {
					nerrors = 0
				}

				next = time.Now().Add(n - time.Since(t))
			}
		}
	}
}

func (s *service) collect() (time.Duration, error) {
	data := api.MonitorData{
		Version: s.version,
	}

	metrics := s.monitor.Collect([]metric.Pattern{
		metric.NewPattern("uptime_uptime"),
		metric.NewPattern("cpu_ncpu"),
		metric.NewPattern("cpu_idle"),
		metric.NewPattern("mem_total"),
		metric.NewPattern("mem_free"),
		metric.NewPattern("disk_total"),
		metric.NewPattern("disk_usage"),
		metric.NewPattern("filesystem_limit"),
		metric.NewPattern("filesystem_usage"),
		metric.NewPattern("ffmpeg_process"),
		metric.NewPattern("restream_process"),
		metric.NewPattern("session_limit"),
		metric.NewPattern("session_active"),
		metric.NewPattern("session_txbitrate"),
		metric.NewPattern("session_maxtxbitrate"),
	})

	data.Uptime = uint64(metrics.Value("uptime_uptime").Val())

	data.SysCPU = []json.Number{
		toNumber(metrics.Value("cpu_ncpu").Val()),
		toNumber(100.0 - metrics.Value("cpu_idle").Val()),
	}

	data.SysMemory = []json.Number{
		toNumber(metrics.Value("mem_total").Val()),
		toNumber(metrics.Value("mem_total").Val() - metrics.Value("mem_free").Val()),
	}

	data.SysDisk = []json.Number{
		toNumber(metrics.Value("disk_total").Val()),
		toNumber(metrics.Value("disk_usage").Val()),
	}

	data.FSMem = []json.Number{
		toNumber(metrics.Value("filesystem_limit", "name", "memfs").Val()),
		toNumber(metrics.Value("filesystem_usage", "name", "memfs").Val()),
	}

	data.FSDisk = []json.Number{
		toNumber(metrics.Value("filesystem_limit", "name", "diskfs").Val()),
		toNumber(metrics.Value("filesystem_usage", "name", "diskfs").Val()),
	}

	data.NetTX = []json.Number{
		toNumber((metrics.Value("session_maxtxbitrate", "collector", "hls").Val()) / 1024),
		toNumber((metrics.Value("session_txbitrate", "collector", "hls").Val() +
			metrics.Value("session_txbitrate", "collector", "rtmp").Val() +
			metrics.Value("session_txbitrate", "collector", "ffmpeg").Val()) / 1024),
	}

	data.Session = []json.Number{
		toNumber(metrics.Value("session_limit", "collector", "hls").Val()),
		toNumber(metrics.Value("session_active", "collector", "hls").Val() + metrics.Value("session_active", "collector", "rtmp").Val()),
	}

	data.ProcessStates[0] = uint64(metrics.Value("ffmpeg_process", "state", "finished").Val())
	data.ProcessStates[4] = uint64(metrics.Value("ffmpeg_process", "state", "failed").Val())
	data.ProcessStates[5] = uint64(metrics.Value("ffmpeg_process", "state", "killed").Val())

	for _, processid := range metrics.Labels("restream_process", "processid") {
		pid := "^" + processid + "$"
		switch metrics.Value("restream_process", "processid", pid).L("state") {
		case "starting":
			data.ProcessStates[1]++
		case "running":
			data.ProcessStates[2]++
		case "finishing":
			data.ProcessStates[3]++
		default:
		}

		/*
			proc := api.MonitorProcessData{
				ID:    processid,
				RefID: "",
				CPU: []float64{
					g.Value("processid", pid, "name", "cpu_limit").Val(),
					g.Value("processid", pid, "name", "cpu").Val(),
				},
				Mem: []uint64{
					uint64(g.Value("processid", pid, "name", "memory_limit").Val()),
					uint64(g.Value("processid", pid, "name", "memory").Val()),
				},
				Uptime: uint64(g.Value("processid", pid, "name", "uptime").Val()),
				Output: map[string][]uint64{},
			}

			data.Processes = append(data.Processes, proc)
		*/
	}

	r, err := s.api.Monitor(s.id, data)
	if err != nil {
		return 15 * time.Minute, fmt.Errorf("failed to send monitor data to service: %w", err)
	}

	s.logger.Debug().WithFields(log.Fields{
		"next": r.Next,
		"data": data,
	}).Log("Sent monitor data")

	if r.Next == 0 {
		r.Next = 5 * 60
	}

	return time.Duration(r.Next) * time.Second, nil
}

func (s *service) Start() {
	s.startOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		s.stopTicker = cancel
		go s.tick(ctx, time.Second)

		s.stopOnce = sync.Once{}

		s.logger.Info().Log("Connected")
	})
}

func (s *service) Stop() {
	s.stopOnce.Do(func() {
		s.stopTicker()
		s.startOnce = sync.Once{}

		s.logger.Info().Log("Disconnected")
	})
}

func toNumber(f float64) json.Number {
	var s string

	if f == float64(int64(f)) {
		s = fmt.Sprintf("%.0f", f) // 0 decimal if integer
	} else {
		s = fmt.Sprintf("%.3f", f) // max. 3 decimal if float
	}

	return json.Number(s)
}
