package ffmpeg

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/datarhei/core/v16/ffmpeg/parse"
	"github.com/datarhei/core/v16/ffmpeg/probe"
	"github.com/datarhei/core/v16/ffmpeg/skills"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/net"
	"github.com/datarhei/core/v16/process"
	"github.com/datarhei/core/v16/session"
)

type FFmpeg interface {
	New(config ProcessConfig) (process.Process, error)
	NewProcessParser(logger log.Logger, id, reference string) parse.Parser
	NewProbeParser(logger log.Logger) probe.Parser
	ValidateInputAddress(address string) bool
	ValidateOutputAddress(address string) bool
	Skills() skills.Skills
	ReloadSkills() error
	GetPort() (int, error)
	PutPort(port int)
	States() process.States
}

type ProcessConfig struct {
	Reconnect      bool
	ReconnectDelay time.Duration
	StaleTimeout   time.Duration
	LimitCPU       float64
	LimitMemory    uint64
	LimitDuration  time.Duration
	Command        []string
	Parser         process.Parser
	Logger         log.Logger
	OnExit         func()
	OnStart        func()
	OnStateChange  func(from, to string)
}

// Config is the configuration for ffmpeg that is part of the configuration
// for the restreamer instance.
type Config struct {
	Binary           string
	MaxProc          int64
	MaxLogLines      int
	LogHistoryLength int
	ValidatorInput   Validator
	ValidatorOutput  Validator
	Portrange        net.Portranger
	Collector        session.Collector
}

type ffmpeg struct {
	binary       string
	validatorIn  Validator
	validatorOut Validator
	portrange    net.Portranger
	skills       skills.Skills

	logLines      int
	historyLength int

	collector session.Collector

	states     process.States
	statesLock sync.RWMutex
}

func New(config Config) (FFmpeg, error) {
	f := &ffmpeg{}

	binary, err := exec.LookPath(config.Binary)
	if err != nil {
		return nil, fmt.Errorf("invalid ffmpeg binary given: %w", err)
	}

	f.binary = binary
	f.historyLength = config.LogHistoryLength
	f.logLines = config.MaxLogLines

	f.portrange = config.Portrange
	if f.portrange == nil {
		f.portrange = net.NewDummyPortrange()
	}

	f.validatorIn = config.ValidatorInput
	if f.validatorIn == nil {
		f.validatorIn, _ = NewValidator(nil, nil)
	}

	f.validatorOut = config.ValidatorOutput
	if f.validatorOut == nil {
		f.validatorOut, _ = NewValidator(nil, nil)
	}

	f.collector = config.Collector
	if f.collector == nil {
		f.collector = session.NewNullCollector()
	}

	s, err := skills.New(f.binary)
	if err != nil {
		return nil, fmt.Errorf("invalid ffmpeg binary given: %w", err)
	}
	f.skills = s

	return f, nil
}

func (f *ffmpeg) New(config ProcessConfig) (process.Process, error) {
	ffmpeg, err := process.New(process.Config{
		Binary:         f.binary,
		Args:           config.Command,
		Reconnect:      config.Reconnect,
		ReconnectDelay: config.ReconnectDelay,
		StaleTimeout:   config.StaleTimeout,
		LimitCPU:       config.LimitCPU,
		LimitMemory:    config.LimitMemory,
		LimitDuration:  config.LimitDuration,
		Parser:         config.Parser,
		Logger:         config.Logger,
		OnStart:        config.OnStart,
		OnExit:         config.OnExit,
		OnStateChange: func(from, to string) {
			f.statesLock.Lock()
			switch to {
			case "finished":
				f.states.Finished++
			case "starting":
				f.states.Starting++
			case "running":
				f.states.Running++
			case "finishing":
				f.states.Finishing++
			case "failed":
				f.states.Failed++
			case "killed":
				f.states.Killed++
			default:
			}
			f.statesLock.Unlock()

			if config.OnStateChange != nil {
				config.OnStateChange(from, to)
			}
		},
	})

	return ffmpeg, err
}

func (f *ffmpeg) NewProcessParser(logger log.Logger, id, reference string) parse.Parser {
	p := parse.New(parse.Config{
		LogHistory: f.historyLength,
		LogLines:   f.logLines,
		Logger:     logger,
		Collector:  NewWrappedCollector(id, reference, f.collector),
	})

	return p
}

func (f *ffmpeg) NewProbeParser(logger log.Logger) probe.Parser {
	p := probe.New(probe.Config{
		Logger: logger,
	})

	return p
}

func (f *ffmpeg) ValidateInputAddress(address string) bool {
	return f.validatorIn.IsValid(address)
}

func (f *ffmpeg) ValidateOutputAddress(address string) bool {
	return f.validatorOut.IsValid(address)
}

func (f *ffmpeg) Skills() skills.Skills {
	return f.skills
}

func (f *ffmpeg) ReloadSkills() error {
	s, err := skills.New(f.binary)
	if err != nil {
		return fmt.Errorf("invalid ffmpeg binary given: %w", err)
	}

	f.skills = s

	return nil
}

func (f *ffmpeg) GetPort() (int, error) {
	return f.portrange.Get()
}

func (f *ffmpeg) PutPort(port int) {
	f.portrange.Put(port)
}

func (f *ffmpeg) States() process.States {
	f.statesLock.RLock()
	defer f.statesLock.RUnlock()

	return f.states
}
