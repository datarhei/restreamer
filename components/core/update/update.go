package update

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/monitor/metric"
	"golang.org/x/mod/semver"
)

// Config is the configuration for the update check
type Config struct {
	ID      string
	Name    string
	Version string
	Arch    string
	Monitor metric.Reader
	Logger  log.Logger
}

// UpdateCheck is an interface
type Checker interface {
	Start()
	Stop()
}

type checker struct {
	id      string
	name    string
	version string
	arch    string
	monitor metric.Reader

	startOnce sync.Once
	stopOnce  sync.Once

	stopTicker context.CancelFunc

	logger log.Logger
}

// New creates a new service instance that implements the Service interface
func New(config Config) (Checker, error) {
	s := &checker{
		id:      config.ID,
		name:    config.Name,
		version: "v" + config.Version,
		arch:    config.Arch,
		monitor: config.Monitor,
		logger:  config.Logger,
	}

	if s.logger == nil {
		s.logger = log.New("")
	}

	if s.monitor == nil {
		return nil, fmt.Errorf("no monitor provided")
	}

	// drain stop once, so it can't be called before startOnce has been called
	s.stopOnce.Do(func() {})

	return s, nil
}

func (s *checker) tick(ctx context.Context, interval, delay time.Duration) {
	time.Sleep(delay)

	err := s.check()
	if err != nil {
		s.logger.WithError(err).Warn().Log("Failed to check for updates")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.check()
			if err != nil {
				s.logger.WithError(err).Warn().Log("Failed to check for updates")
			}
		}
	}
}

func (s *checker) Start() {
	s.startOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		s.stopTicker = cancel
		go s.tick(ctx, 24*time.Hour, 10*time.Second)

		s.stopOnce = sync.Once{}
	})
}

func (s *checker) Stop() {
	s.stopOnce.Do(func() {
		s.stopTicker()
		s.startOnce = sync.Once{}
	})
}

type checkRequest struct {
	AppVersion         string `json:"app_version"`
	CoreID             string `json:"core_id"`
	CoreArch           string `json:"core_aarch"`
	CoreUptimeSeconds  uint64 `json:"core_uptime_seconds"`
	CoreProcessRunning uint64 `json:"core_process_running"`
	CoreProcessFailed  uint64 `json:"core_process_failed"`
	CoreProcessKilled  uint64 `json:"core_process_killed"`
	CoreViewer         uint64 `json:"core_viewer"`
}

type checkResponse struct {
	LatestVersion string `json:"latest_version"`
}

func (s *checker) check() error {
	metrics := s.monitor.Collect([]metric.Pattern{
		metric.NewPattern("uptime_uptime"),
		metric.NewPattern("ffmpeg_process"),
		metric.NewPattern("restream_state"),
		metric.NewPattern("session_active"),
	})

	request := checkRequest{
		AppVersion:         s.name + " " + s.version,
		CoreID:             s.id,
		CoreArch:           s.arch,
		CoreUptimeSeconds:  uint64(metrics.Value("uptime_uptime").Val()),
		CoreProcessRunning: uint64(metrics.Value("restream_state", "state", "running").Val()),
		CoreProcessFailed:  uint64(metrics.Value("ffmpeg_process", "state", "failed").Val()),
		CoreProcessKilled:  uint64(metrics.Value("ffmpeg_process", "state", "killed").Val()),
		CoreViewer:         uint64(metrics.Value("session_active", "collector", "hls").Val() + metrics.Value("session_active", "collector", "rtmp").Val()),
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	if err := encoder.Encode(&request); err != nil {
		return err
	}

	s.logger.Debug().WithField("request", data.String()).Log("")

	req, err := http.NewRequest(http.MethodPut, "https://service.datarhei.com/api/v1/app_version", &data)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("request failed: %s", http.StatusText(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	response := checkResponse{}

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	re := regexp.MustCompile(`\s(v\d+\.\d+\.\d+)\s?`)
	matches := re.FindStringSubmatch(response.LatestVersion)
	if matches == nil {
		return fmt.Errorf("no version information detected in response")
	}

	cmp := semver.Compare(matches[1], s.version)

	s.logger.Debug().WithFields(log.Fields{
		"comparison": cmp,
		"current":    s.version,
		"available":  matches[1],
	}).Log("")

	if cmp == 1 {
		s.logger.Info().WithFields(log.Fields{
			"current":   s.version,
			"available": matches[1],
		}).Log("New version available")
	}

	return nil
}
