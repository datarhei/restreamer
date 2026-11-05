// Package process is a wrapper of exec.Cmd for controlling a ffmpeg process.
// It could be used to run other executables but it is tailored to the specifics
// of ffmpeg, e.g. only stderr is captured, and some exit codes != 0 plus certain
// signals are still considered as a non-error exit condition.
package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/psutil"
)

// Process represents a process and ways to control it
// and to extract information.
type Process interface {
	// Status returns the current status of this process
	Status() Status

	// Start starts the process. If the process stops by itself
	// it will restart automatically if it is defined to do so.
	Start() error

	// Stop stops the process and will not let it restart
	// automatically.
	Stop(wait bool) error

	// Kill stops the process such that it will restart
	// automatically if it is defined to do so.
	Kill(wait bool) error

	// IsRunning returns whether the process is currently
	// running or not.
	IsRunning() bool
}

// Config is the configuration of a process
type Config struct {
	Binary         string                // Path to the ffmpeg binary
	Args           []string              // List of arguments for the binary
	Reconnect      bool                  // Whether to restart the process if it exited
	ReconnectDelay time.Duration         // Duration to wait before restarting the process
	StaleTimeout   time.Duration         // Kill the process after this duration if it doesn't produce any output
	LimitCPU       float64               // Kill the process if the CPU usage in percent is above this value
	LimitMemory    uint64                // Kill the process if the memory consumption in bytes is above this value
	LimitDuration  time.Duration         // Kill the process if the limits are exceeded for this duration
	Parser         Parser                // A parser for the output of the process
	OnStart        func()                // A callback which is called after the process started
	OnExit         func()                // A callback which is called after the process exited
	OnStateChange  func(from, to string) // A callback which is called after a state changed
	Logger         log.Logger
}

// Status represents the current status of a process
type Status struct {
	State    string        // State is the current state of the process. See stateType for the known states.
	States   States        // States is the cumulative history of states the process had.
	Order    string        // Order is the wanted condition of process, either "start" or "stop"
	Duration time.Duration // Duration is the time since the last change of the state
	Time     time.Time     // Time is the time of the last change of the state
	CPU      struct {
		Current float64 // Used CPU in percent
		Limit   float64 // Limit in percent
	}
	Memory struct {
		Current uint64 // Used memory in bytes
		Limit   uint64 // Limit in bytes
	}
}

// States
//
// finished - Process has been stopped
//
//	starting - if process has been actively started or has been waiting for reconnect (order=start, reconnect=any)
//	finished - if process shall not reconnect (order=stop, reconnect=any)
//
// starting - Process is about to start
//
//	finishing - if process should be immediately stopped (order=stop, reconnect=any)
//	running - if process could be started (order=start, reconnect=any)
//	failed - if process couldn't be started (e.g. binary not found) (order=start, reconnect=any)
//
// running - Process is running
//
//	finished - if process exited normally (order=any, reconnect=any)
//	finishing - if process has been actively stopped (order=stop, reconnect=any)
//	failed - if process exited abnormally (order=any, reconnect=any)
//	killed - if process has been actively killed with SIGKILL (order=any, reconnect=any)
//
// finishing - Process has been actively stopped and will be killed
//
//	finished - if process has been actively killed with SIGINT and ffmpeg exited normally (order=stop, reconnect=any)
//	killed - if process has been actively killed with SIGKILL (order=stop, reconnect=any)
//
// failed - Process has been failed either by starting or during running
//
//	starting - if process has been waiting for reconnect (order=start, reconnect=true)
//	failed - if process shall not reconnect (order=any, reconnect=false)
//
// killed - Process has been stopped
//
//	starting - if process has been waiting for reconnect (order=start, reconnect=true)
//	killed - if process shall not reconnect (order=start, reconnect=false)
type stateType string

const (
	stateFinished  stateType = "finished"
	stateStarting  stateType = "starting"
	stateRunning   stateType = "running"
	stateFinishing stateType = "finishing"
	stateFailed    stateType = "failed"
	stateKilled    stateType = "killed"
)

// String returns a string representation of the state
func (s stateType) String() string {
	return string(s)
}

// IsRunning returns whether the state is representing a running state
func (s stateType) IsRunning() bool {
	if s == stateStarting || s == stateRunning || s == stateFinishing {
		return true
	}

	return false
}

type States struct {
	Finished  uint64
	Starting  uint64
	Running   uint64
	Finishing uint64
	Failed    uint64
	Killed    uint64
}

// Process represents a ffmpeg process
type process struct {
	binary   string
	args     []string
	cmd      *exec.Cmd
	pid      int32
	stdout   io.ReadCloser
	lastLine string
	state    struct {
		state  stateType
		time   time.Time
		states States
		lock   sync.Mutex
	}
	order struct {
		order string
		lock  sync.Mutex
	}
	parser Parser
	stale  struct {
		last    time.Time
		timeout time.Duration
		cancel  context.CancelFunc
		lock    sync.Mutex
	}
	reconn struct {
		enable bool
		delay  time.Duration
		timer  *time.Timer
		lock   sync.Mutex
	}
	killTimer     *time.Timer
	killTimerLock sync.Mutex
	logger        log.Logger
	debuglogger   log.Logger
	callbacks     struct {
		onStart       func()
		onExit        func()
		onStateChange func(from, to string)
		lock          sync.Mutex
	}
	limits Limiter
}

var _ Process = &process{}

// New creates a new process wrapper
func New(config Config) (Process, error) {
	p := &process{
		binary: config.Binary,
		args:   config.Args,
		cmd:    nil,
		parser: config.Parser,
		logger: config.Logger,
	}

	// This is a loose check on purpose. If the e.g. the binary
	// doesn't exist or it is not executable, it will be
	// reflected in the resulting state.
	if len(p.binary) == 0 {
		return nil, fmt.Errorf("no valid binary given")
	}

	if p.parser == nil {
		p.parser = NewNullParser()
	}

	if p.logger == nil {
		p.logger = log.New("Process")
	}

	p.debuglogger = p.logger.WithFields(log.Fields{
		"binary": p.binary,
		"args":   p.args,
	})

	p.order.order = "stop"

	p.initState(stateFinished)

	p.reconn.enable = config.Reconnect
	p.reconn.delay = config.ReconnectDelay

	p.stale.last = time.Now()
	p.stale.timeout = config.StaleTimeout

	p.callbacks.onStart = config.OnStart
	p.callbacks.onExit = config.OnExit
	p.callbacks.onStateChange = config.OnStateChange

	p.limits = NewLimiter(LimiterConfig{
		CPU:     config.LimitCPU,
		Memory:  config.LimitMemory,
		WaitFor: config.LimitDuration,
		OnLimit: func(cpu float64, memory uint64) {
			p.logger.WithFields(log.Fields{
				"cpu":    cpu,
				"memory": memory,
			}).Warn().Log("Stopping because limits are exceeded")
			p.Kill(false)
		},
	})

	p.logger.Info().Log("Created")
	p.debuglogger.Debug().Log("Created")

	return p, nil
}

func (p *process) initState(state stateType) {
	p.state.lock.Lock()
	defer p.state.lock.Unlock()

	p.state.state = state
	p.state.time = time.Now()
}

// setState sets a new state. It also checks if the transition
// of the current state to the new state is allowed. If not,
// the current state will not be changed.
func (p *process) setState(state stateType) error {
	p.state.lock.Lock()
	defer p.state.lock.Unlock()

	prevState := p.state.state
	failed := false

	if p.state.state == stateFinished {
		switch state {
		case stateStarting:
			p.state.state = state
			p.state.states.Starting++
		default:
			failed = true
		}
	} else if p.state.state == stateStarting {
		switch state {
		case stateFinishing:
			p.state.state = state
			p.state.states.Finishing++
		case stateRunning:
			p.state.state = state
			p.state.states.Running++
		case stateFailed:
			p.state.state = state
			p.state.states.Failed++
		default:
			failed = true
		}
	} else if p.state.state == stateRunning {
		switch state {
		case stateFinished:
			p.state.state = state
			p.state.states.Finished++
		case stateFinishing:
			p.state.state = state
			p.state.states.Finishing++
		case stateFailed:
			p.state.state = state
			p.state.states.Failed++
		case stateKilled:
			p.state.state = state
			p.state.states.Killed++
		default:
			failed = true
		}
	} else if p.state.state == stateFinishing {
		switch state {
		case stateFinished:
			p.state.state = state
			p.state.states.Finished++
		case stateFailed:
			p.state.state = state
			p.state.states.Failed++
		case stateKilled:
			p.state.state = state
			p.state.states.Killed++
		default:
			failed = true
		}
	} else if p.state.state == stateFailed {
		switch state {
		case stateStarting:
			p.state.state = state
			p.state.states.Starting++
		default:
			failed = true
		}
	} else if p.state.state == stateKilled {
		switch state {
		case stateStarting:
			p.state.state = state
			p.state.states.Starting++
		default:
			failed = true
		}
	} else {
		return fmt.Errorf("current state is unhandled: %s", p.state.state)
	}

	if failed {
		return fmt.Errorf("can't change from state %s to %s", p.state.state, state)
	}

	p.state.time = time.Now()

	if p.callbacks.onStateChange != nil {
		go p.callbacks.onStateChange(prevState.String(), p.state.state.String())
	}

	return nil
}

func (p *process) getState() stateType {
	p.state.lock.Lock()
	defer p.state.lock.Unlock()

	return p.state.state
}

func (p *process) isRunning() bool {
	p.state.lock.Lock()
	defer p.state.lock.Unlock()

	return p.state.state.IsRunning()
}

func (p *process) getStateString() string {
	p.state.lock.Lock()
	defer p.state.lock.Unlock()

	return p.state.state.String()
}

// Status returns the current status of the process
func (p *process) Status() Status {
	cpu, memory := p.limits.Current()
	cpuLimit, memoryLimit := p.limits.Limits()

	p.state.lock.Lock()
	stateTime := p.state.time
	stateString := p.state.state.String()
	states := p.state.states
	p.state.lock.Unlock()

	p.order.lock.Lock()
	order := p.order.order
	p.order.lock.Unlock()

	s := Status{
		State:    stateString,
		States:   states,
		Order:    order,
		Duration: time.Since(stateTime),
		Time:     stateTime,
	}

	s.CPU.Current = cpu
	s.CPU.Limit = cpuLimit

	s.Memory.Current = memory
	s.Memory.Limit = memoryLimit

	return s
}

// IsRunning returns whether the process is considered running
func (p *process) IsRunning() bool {
	return p.isRunning()
}

// Start will start the process and sets the order to "start". If the
// process has alread the "start" order, nothing will be done. Returns
// an error if start failed.
func (p *process) Start() error {
	p.order.lock.Lock()
	defer p.order.lock.Unlock()

	if p.order.order == "start" {
		return nil
	}

	p.order.order = "start"

	err := p.start()
	if err != nil {
		p.debuglogger.WithFields(log.Fields{
			"state": p.getStateString(),
			"order": p.order.order,
			"error": err,
		}).Debug().Log("Starting failed")
	}

	return err
}

// start will start the process considering the current order. Returns an
// error in case something goes wrong, and it will try to restart the process.
func (p *process) start() error {
	var err error

	// Bail out if the process is already running
	if p.isRunning() {
		return nil
	}

	p.logger.Info().Log("Starting")
	p.debuglogger.WithFields(log.Fields{
		"state": p.getStateString(),
		"order": p.order.order,
	}).Debug().Log("Starting")

	// Stop any restart timer in order to start the process immediately
	p.unreconnect()

	p.setState(stateStarting)

	p.cmd = exec.Command(p.binary, p.args...)
	p.cmd.Env = []string{}

	p.stdout, err = p.cmd.StderrPipe()
	if err != nil {
		p.setState(stateFailed)

		p.parser.Parse(err.Error())
		p.logger.WithError(err).Error().Log("Command failed")
		p.reconnect()

		return err
	}
	if err := p.cmd.Start(); err != nil {
		p.setState(stateFailed)

		p.parser.Parse(err.Error())
		p.logger.WithError(err).Error().Log("Command failed")
		p.reconnect()

		return err
	}

	p.pid = int32(p.cmd.Process.Pid)

	if proc, err := psutil.NewProcess(p.pid); err == nil {
		p.limits.Start(proc)
	}

	p.setState(stateRunning)

	p.logger.Info().Log("Started")
	p.debuglogger.Debug().Log("Started")

	if p.callbacks.onStart != nil {
		go p.callbacks.onStart()
	}

	// Start the reader
	go p.reader()

	// Start the stale timeout if enabled
	if p.stale.timeout != 0 {
		var ctx context.Context

		p.stale.lock.Lock()
		ctx, p.stale.cancel = context.WithCancel(context.Background())
		p.stale.lock.Unlock()

		go p.staler(ctx)
	}

	return nil
}

// Stop will stop the process and set the order to "stop"
func (p *process) Stop(wait bool) error {
	p.order.lock.Lock()
	defer p.order.lock.Unlock()

	if p.order.order == "stop" {
		return nil
	}

	p.order.order = "stop"

	err := p.stop(wait)
	if err != nil {
		p.debuglogger.WithFields(log.Fields{
			"state": p.getStateString(),
			"order": p.order.order,
			"error": err,
		}).Debug().Log("Stopping failed")
	}

	return err
}

// Kill will stop the process without changing the order such that it
// will restart automatically if enabled.
func (p *process) Kill(wait bool) error {
	// If the process is currently not running, we don't need
	// to do anything.
	if !p.isRunning() {
		return nil
	}

	p.order.lock.Lock()
	defer p.order.lock.Unlock()

	err := p.stop(wait)

	return err
}

// stop will stop a process considering the current order and state.
func (p *process) stop(wait bool) error {
	// If the process is currently not running, stop the restart timer
	if !p.isRunning() {
		p.unreconnect()
		return nil
	}

	// If the process is already in the finishing state, don't do anything
	if p.getState() == stateFinishing {
		return nil
	}

	p.setState(stateFinishing)

	p.logger.Info().Log("Stopping")
	p.debuglogger.WithFields(log.Fields{
		"state": p.getStateString(),
		"order": p.order.order,
	}).Debug().Log("Stopping")

	wg := sync.WaitGroup{}

	if wait {
		wg.Add(1)

		p.callbacks.lock.Lock()
		if p.callbacks.onExit == nil {
			p.callbacks.onExit = func() {
				wg.Done()

				p.callbacks.lock.Lock()
				p.callbacks.onExit = nil
				p.callbacks.lock.Unlock()
			}
		} else {
			cb := p.callbacks.onExit
			p.callbacks.onExit = func() {
				cb()
				wg.Done()

				p.callbacks.lock.Lock()
				p.callbacks.onExit = cb
				p.callbacks.lock.Unlock()
			}
		}
		p.callbacks.lock.Unlock()
	}

	var err error
	if runtime.GOOS == "windows" {
		// Windows doesn't know the SIGINT
		err = p.cmd.Process.Kill()
	} else {
		// First try to kill the process gracefully. On a SIGINT ffmpeg will exit
		// normally as if "q" has been pressed.
		err = p.cmd.Process.Signal(os.Interrupt)
		if err != nil {
			// If sending the signal fails, try it the hard way, however this will highly
			// likely also fail because it is simply a shortcut for Signal(Kill).
			err = p.cmd.Process.Kill()
		} else {
			// Set up a timer to kill the process with SIGKILL in case SIGINT didn't have
			// an effect.
			p.killTimerLock.Lock()
			p.killTimer = time.AfterFunc(5*time.Second, func() {
				p.cmd.Process.Kill()
			})
			p.killTimerLock.Unlock()
		}
	}

	if err == nil && wait {
		wg.Wait()
	}

	if err != nil {
		p.parser.Parse(err.Error())
		p.debuglogger.WithFields(log.Fields{
			"state": p.getStateString(),
			"order": p.order.order,
			"error": err,
		}).Debug().Log("Stopping failed")

		p.setState(stateFailed)
	}

	return err
}

// reconnect will setup a timer to restart the  process
func (p *process) reconnect() {
	// If restarting a process is not enabled, don't do anything
	if !p.reconn.enable {
		return
	}

	// Stop a currently running timer
	p.unreconnect()

	p.logger.Info().Log("Scheduling restart in %s", p.reconn.delay)

	p.reconn.lock.Lock()
	defer p.reconn.lock.Unlock()

	p.reconn.timer = time.AfterFunc(p.reconn.delay, func() {
		p.order.lock.Lock()
		defer p.order.lock.Unlock()

		p.start()
	})
}

// unreconnect will stop the restart timer
func (p *process) unreconnect() {
	p.reconn.lock.Lock()
	defer p.reconn.lock.Unlock()

	if p.reconn.timer == nil {
		return
	}

	p.reconn.timer.Stop()
	p.reconn.timer = nil
}

// staler checks if the currently running process is stale, i.e. the reader
// didn't update the time of the last read. If the timeout is reached, the
// process will be stopped such that it can restart automatically afterwards.
func (p *process) staler(ctx context.Context) {
	p.stale.lock.Lock()
	p.stale.last = time.Now()
	p.stale.lock.Unlock()

	p.debuglogger.Debug().Log("Starting stale watcher")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.debuglogger.Debug().Log("Stopping stale watcher")
			return
		case t := <-ticker.C:
			p.stale.lock.Lock()
			last := p.stale.last
			timeout := p.stale.timeout
			p.stale.lock.Unlock()

			d := t.Sub(last)
			if d.Seconds() > timeout.Seconds() {
				p.logger.Info().Log("Stale timeout after %s (%.2f).", timeout, d.Seconds())
				p.stop(false)
				return
			}
		}
	}
}

// reader reads the output from the process line by line and gives
// each line to the parser. The parser returns a postive number to
// indicate progress. If the returned number is zero, then the time
// of the last progress will not be updated thus the stale timeout
// may kick in.
func (p *process) reader() {
	scanner := bufio.NewScanner(p.stdout)
	scanner.Split(scanLine)

	// Reset the parser statistics
	p.parser.ResetStats()

	// Reset the parser logs
	p.parser.ResetLog()

	var n uint64 = 0

	for scanner.Scan() {
		line := scanner.Text()

		p.lastLine = line

		// Parse the output line from ffmpeg
		n = p.parser.Parse(line)

		// Reset the stale progress timer only if the
		// parser reports progress
		if n != 0 {
			p.stale.lock.Lock()
			p.stale.last = time.Now()
			p.stale.lock.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		p.logger.Debug().WithError(err).Log("")
	}

	// Wait for the process to finish
	p.waiter()
}

// waiter waits for the process to finish. If enabled, the process will
// be scheduled for a restart.
func (p *process) waiter() {
	if err := p.cmd.Wait(); err != nil {
		// The process exited abnormally, i.e. the return code is non-zero or a signal
		// has been raised.
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The process exited and the status can be examined
			status := exiterr.Sys().(syscall.WaitStatus)

			p.debuglogger.WithFields(log.Fields{
				"exited":      status.Exited(),
				"signaled":    status.Signaled(),
				"status":      status.ExitStatus(),
				"exit_code":   exiterr.ExitCode(),
				"exit_string": exiterr.String(),
				"signal":      status.Signal(),
			}).Debug().Log("Exited")

			if status.Exited() {
				if status.ExitStatus() == 255 {
					// If ffmpeg has been killed with a SIGINT, SIGTERM, etc., then it exited normally,
					// i.e. closing all stream properly such that all written data is sane.
					p.logger.Info().Log("Finished")
					p.setState(stateFinished)
				} else {
					// The process exited by itself with a non-zero return code
					p.logger.Info().Log("Failed")
					p.setState(stateFailed)
				}
			} else if status.Signaled() {
				// If ffmpeg has been killed the hard way, something went wrong and
				// it can be assumed that any written data is not sane.
				p.logger.Info().Log("Killed")
				p.setState(stateKilled)
			} else {
				// The process exited because of something else (e.g. coredump, ...)
				p.logger.Info().Log("Killed")
				p.setState(stateKilled)
			}
		} else {
			// Some other error regarding I/O triggered during Wait()
			p.logger.Info().Log("Killed")
			p.logger.WithError(err).Debug().Log("Killed")
			p.setState(stateKilled)
		}
	} else {
		// The process exited normally, i.e. the return code is zero and no signal
		// has been raised
		p.setState(stateFinished)
	}

	p.logger.Info().Log("Stopped")
	p.debuglogger.WithField("log", p.parser.Log()).Debug().Log("Stopped")

	p.limits.Stop()

	// Stop the kill timer
	p.killTimerLock.Lock()
	if p.killTimer != nil {
		p.killTimer.Stop()
		p.killTimer = nil
	}
	p.killTimerLock.Unlock()

	// Stop the stale progress timer
	p.stale.lock.Lock()
	if p.stale.cancel != nil {
		p.stale.cancel()
		p.stale.cancel = nil
	}
	p.stale.lock.Unlock()

	// Reset the parser stats
	p.parser.ResetStats()

	// Call the onExit callback
	p.callbacks.lock.Lock()
	if p.callbacks.onExit != nil {
		go p.callbacks.onExit()
	}
	p.callbacks.lock.Unlock()

	p.order.lock.Lock()
	defer p.order.lock.Unlock()

	p.debuglogger.WithFields(log.Fields{
		"state": p.getStateString(),
		"order": p.order.order,
	}).Debug().Log("Waiting")

	// Restart the process
	if p.order.order == "start" {
		p.reconnect()
	}
}

// scanLine splits the data on \r, \n, or \r\n line endings
func scanLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != '\n' && r != '\r' {
			break
		}

	}

	// Scan until new line, marking end of line.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '\n' || r == '\r' {
			return i + width, data[start:i], nil
		}

	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}
