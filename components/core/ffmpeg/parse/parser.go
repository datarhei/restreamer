package parse

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/ffmpeg/prelude"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/net/url"
	"github.com/datarhei/core/v16/process"
	"github.com/datarhei/core/v16/restream/app"
	"github.com/datarhei/core/v16/session"
)

// Parser is an extension to the process.Parser interface
type Parser interface {
	process.Parser

	// Progress returns the current progress information of the process
	Progress() app.Progress

	// Prelude returns an array of the lines before the progress information started
	Prelude() []string

	// Report returns the current logs
	Report() Report

	// ReportHistory returns an array of previews logs
	ReportHistory() []Report

	// TransferReportHistory transfers the report history to another parser
	TransferReportHistory(Parser) error
}

// Config is the config for the Parser implementation
type Config struct {
	LogHistory       int
	LogLines         int
	PreludeHeadLines int
	PreludeTailLines int
	Logger           log.Logger
	Collector        session.Collector
}

type parser struct {
	re struct {
		frame     *regexp.Regexp
		quantizer *regexp.Regexp
		size      *regexp.Regexp
		time      *regexp.Regexp
		speed     *regexp.Regexp
		drop      *regexp.Regexp
		dup       *regexp.Regexp
	}

	prelude struct {
		headLines      int
		tailLines      int
		truncatedLines uint64
		data           []string
		tail           *ring.Ring
		done           bool
	}

	log      *ring.Ring
	logLines int
	logStart time.Time

	logHistory       *ring.Ring
	logHistoryLength int

	progress struct {
		ffmpeg   ffmpegProgress
		avstream map[string]ffmpegAVstream
	}

	process ffmpegProcess

	stats struct {
		initialized bool
		main        stats
		input       []stats
		output      []stats
	}

	averager struct {
		initialized bool
		window      time.Duration
		granularity time.Duration
		main        averager
		input       []averager
		output      []averager
	}

	collector session.Collector

	logger log.Logger

	lock struct {
		progress sync.RWMutex
		prelude  sync.RWMutex
		log      sync.RWMutex
	}
}

// New returns a Parser that satisfies the Parser interface
func New(config Config) Parser {
	p := &parser{
		logHistoryLength: config.LogHistory,
		logLines:         config.LogLines,
		logger:           config.Logger,
		collector:        config.Collector,
	}

	if p.logger == nil {
		p.logger = log.New("Parser")
	}

	if p.logLines <= 0 {
		p.logLines = 1
	}

	p.averager.window = 30 * time.Second
	p.averager.granularity = time.Second

	p.re.frame = regexp.MustCompile(`frame=\s*([0-9]+)`)
	p.re.quantizer = regexp.MustCompile(`q=\s*([0-9\.]+)`)
	p.re.size = regexp.MustCompile(`size=\s*([0-9]+)kB`)
	p.re.time = regexp.MustCompile(`time=\s*([0-9]+):([0-9]{2}):([0-9]{2}).([0-9]{2})`)
	p.re.speed = regexp.MustCompile(`speed=\s*([0-9\.]+)x`)
	p.re.drop = regexp.MustCompile(`drop=\s*([0-9]+)`)
	p.re.dup = regexp.MustCompile(`dup=\s*([0-9]+)`)

	p.lock.prelude.Lock()
	p.prelude.headLines = config.PreludeHeadLines
	if p.prelude.headLines <= 0 {
		p.prelude.headLines = 100
	}
	p.prelude.tailLines = config.PreludeTailLines
	if p.prelude.tailLines <= 0 {
		p.prelude.tailLines = 50
	}
	p.prelude.tail = ring.New(p.prelude.tailLines)
	p.lock.prelude.Unlock()

	p.lock.log.Lock()
	p.log = ring.New(config.LogLines)

	if p.logHistoryLength > 0 {
		p.logHistory = ring.New(p.logHistoryLength)
	}

	if p.collector == nil {
		p.collector = session.NewNullCollector()
	}

	p.logStart = time.Now()
	p.lock.log.Unlock()

	p.ResetStats()

	return p
}

func (p *parser) Parse(line string) uint64 {
	isDefaultProgress := strings.HasPrefix(line, "frame=")
	isFFmpegInputs := strings.HasPrefix(line, "ffmpeg.inputs:")
	isFFmpegOutputs := strings.HasPrefix(line, "ffmpeg.outputs:")
	isFFmpegProgress := strings.HasPrefix(line, "ffmpeg.progress:")
	isAVstreamProgress := strings.HasPrefix(line, "avstream.progress:")

	if p.logStart.IsZero() {
		p.lock.log.Lock()
		p.logStart = time.Now()
		p.lock.log.Unlock()
	}

	p.lock.prelude.Lock()
	preludeDone := p.prelude.done
	p.lock.prelude.Unlock()

	if !preludeDone {
		if isAVstreamProgress {
			return 0
		}

		if isFFmpegProgress {
			return 0
		}

		if isFFmpegInputs {
			if err := p.parseIO("input", strings.TrimPrefix(line, "ffmpeg.inputs:")); err != nil {
				p.logger.WithFields(log.Fields{
					"line":  line,
					"error": err,
				}).Error().Log("Failed parsing inputs")
			}

			return 0
		}

		if isFFmpegOutputs {
			if err := p.parseIO("output", strings.TrimPrefix(line, "ffmpeg.outputs:")); err != nil {
				p.logger.WithFields(log.Fields{
					"line":  line,
					"error": err,
				}).Error().Log("Failed parsing outputs")
			} else {
				p.logger.WithField("prelude", p.Prelude()).Debug().Log("")

				p.lock.prelude.Lock()
				p.prelude.done = true
				p.lock.prelude.Unlock()
			}

			return 0
		}

		if isDefaultProgress {
			if !p.parsePrelude() {
				return 0
			}

			p.logger.WithField("prelude", p.Prelude()).Debug().Log("")

			p.lock.prelude.Lock()
			p.prelude.done = true
			p.lock.prelude.Unlock()
		}
	}

	if !isDefaultProgress && !isFFmpegProgress && !isAVstreamProgress {
		// Write the current non-progress line to the log
		p.addLog(line)

		p.lock.prelude.Lock()
		if !p.prelude.done {
			if len(p.prelude.data) < p.prelude.headLines {
				p.prelude.data = append(p.prelude.data, line)
			} else {
				p.prelude.tail.Value = line
				p.prelude.tail = p.prelude.tail.Next()
				p.prelude.truncatedLines++
			}
		}
		p.lock.prelude.Unlock()

		return 0
	}

	if !p.prelude.done {
		return 0
	}

	p.lock.progress.Lock()
	defer p.lock.progress.Unlock()

	// Initialize the averagers

	if !p.averager.initialized {
		p.averager.main.init(p.averager.window, p.averager.granularity)

		p.averager.input = make([]averager, len(p.process.input))
		for i := range p.averager.input {
			p.averager.input[i].init(p.averager.window, p.averager.granularity)
		}

		p.averager.output = make([]averager, len(p.process.output))
		for i := range p.averager.output {
			p.averager.output[i].init(p.averager.window, p.averager.granularity)
		}

		p.averager.initialized = true
	}

	// Initialize the stats

	if !p.stats.initialized {
		p.stats.input = make([]stats, len(p.process.input))
		p.stats.output = make([]stats, len(p.process.output))

		p.collector.Register("", "", "", "")

		p.stats.initialized = true
	}

	// Update the progress

	if isAVstreamProgress {
		if err := p.parseAVstreamProgress(strings.TrimPrefix(line, "avstream.progress:")); err != nil {
			p.logger.WithFields(log.Fields{
				"line":  line,
				"error": err,
			}).Error().Log("Failed parsing AVStream progress")
		}

		return 0
	}

	if isDefaultProgress {
		if err := p.parseDefaultProgress(line); err != nil {
			p.logger.WithFields(log.Fields{
				"line":  line,
				"error": err,
			}).Error().Log("Failed parsing default progress")
			return 0
		}
	} else if isFFmpegProgress {
		if err := p.parseFFmpegProgress(strings.TrimPrefix(line, "ffmpeg.progress:")); err != nil {
			p.logger.WithFields(log.Fields{
				"line":  line,
				"error": err,
			}).Error().Log("Failed parsing progress")
			return 0
		}
	} else {
		return 0
	}

	// Update the averages

	p.stats.main.updateFromProgress(&p.progress.ffmpeg)

	if len(p.stats.input) != 0 && len(p.stats.input) == len(p.progress.ffmpeg.Input) {
		for i := range p.progress.ffmpeg.Input {
			p.stats.input[i].updateFromProgressIO(&p.progress.ffmpeg.Input[i])
		}
	}

	if len(p.stats.output) != 0 && len(p.stats.output) == len(p.progress.ffmpeg.Output) {
		for i := range p.progress.ffmpeg.Output {
			p.stats.output[i].updateFromProgressIO(&p.progress.ffmpeg.Output[i])
		}
	}

	p.averager.main.fps.Add(int64(p.stats.main.diff.frame))
	p.averager.main.pps.Add(int64(p.stats.main.diff.packet))
	p.averager.main.bitrate.Add(int64(p.stats.main.diff.size) * 8)

	p.progress.ffmpeg.FPS = p.averager.main.fps.Average(p.averager.window)
	p.progress.ffmpeg.PPS = p.averager.main.pps.Average(p.averager.window)
	p.progress.ffmpeg.Bitrate = p.averager.main.bitrate.Average(p.averager.window)

	if len(p.averager.input) != 0 && len(p.averager.input) == len(p.progress.ffmpeg.Input) {
		for i := range p.progress.ffmpeg.Input {
			p.averager.input[i].fps.Add(int64(p.stats.input[i].diff.frame))
			p.averager.input[i].pps.Add(int64(p.stats.input[i].diff.packet))
			p.averager.input[i].bitrate.Add(int64(p.stats.input[i].diff.size) * 8)

			p.progress.ffmpeg.Input[i].FPS = p.averager.input[i].fps.Average(p.averager.window)
			p.progress.ffmpeg.Input[i].PPS = p.averager.input[i].pps.Average(p.averager.window)
			p.progress.ffmpeg.Input[i].Bitrate = p.averager.input[i].bitrate.Average(p.averager.window)

			if p.collector.IsCollectableIP(p.process.input[i].IP) {
				p.collector.Activate("")
				p.collector.Ingress("", int64(p.stats.input[i].diff.size))
			}
		}
	}

	if len(p.averager.output) != 0 && len(p.averager.output) == len(p.progress.ffmpeg.Output) {
		for i := range p.progress.ffmpeg.Output {
			p.averager.output[i].fps.Add(int64(p.stats.output[i].diff.frame))
			p.averager.output[i].pps.Add(int64(p.stats.output[i].diff.packet))
			p.averager.output[i].bitrate.Add(int64(p.stats.output[i].diff.size) * 8)

			p.progress.ffmpeg.Output[i].FPS = p.averager.output[i].fps.Average(p.averager.window)
			p.progress.ffmpeg.Output[i].PPS = p.averager.output[i].pps.Average(p.averager.window)
			p.progress.ffmpeg.Output[i].Bitrate = p.averager.output[i].bitrate.Average(p.averager.window)

			if p.collector.IsCollectableIP(p.process.output[i].IP) {
				p.collector.Activate("")
				p.collector.Egress("", int64(p.stats.output[i].diff.size))
			}
		}
	}

	// Calculate if any of the processed frames staled.
	// If one number of frames in an output is the same as before, then pFrames becomes 0.
	pFrames := p.stats.main.diff.frame

	if isFFmpegProgress {
		// Only consider the outputs
		pFrames = 1
		for i := range p.stats.output {
			pFrames *= p.stats.output[i].diff.frame
		}
	}

	return pFrames
}

func (p *parser) parseDefaultProgress(line string) error {
	var matches []string

	if matches = p.re.frame.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			p.progress.ffmpeg.Frame = x
		}
	}

	if matches = p.re.quantizer.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
			p.progress.ffmpeg.Quantizer = x
		}
	}

	if matches = p.re.size.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			p.progress.ffmpeg.Size = x * 1024
		}
	}

	if matches = p.re.time.FindStringSubmatch(line); matches != nil {
		s := fmt.Sprintf("%sh%sm%ss%s0ms", matches[1], matches[2], matches[3], matches[4])
		if x, err := time.ParseDuration(s); err == nil {
			p.progress.ffmpeg.Time.Duration = x
		}
	}

	if matches = p.re.speed.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
			p.progress.ffmpeg.Speed = x
		}
	}

	if matches = p.re.drop.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			p.progress.ffmpeg.Drop = x
		}
	}

	if matches = p.re.dup.FindStringSubmatch(line); matches != nil {
		if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			p.progress.ffmpeg.Dup = x
		}
	}

	return nil
}

func (p *parser) parseIO(kind, line string) error {
	processIO := []ffmpegProcessIO{}

	err := json.Unmarshal([]byte(line), &processIO)
	if err != nil {
		return err
	}

	if len(processIO) == 0 {
		return fmt.Errorf("the %s length must not be 0", kind)
	}

	for i := range processIO {
		if ip, _ := url.Lookup(processIO[i].Address); len(ip) != 0 {
			processIO[i].IP = ip
		}
	}

	if kind == "input" {
		p.process.input = processIO
	} else if kind == "output" {
		p.process.output = processIO
	}

	return nil
}

func (p *parser) parseFFmpegProgress(line string) error {
	progress := ffmpegProgress{}

	err := json.Unmarshal([]byte(line), &progress)
	if err != nil {
		return err
	}

	if len(progress.Input) != len(p.process.input) {
		return fmt.Errorf("input length mismatch (have: %d, want: %d)", len(progress.Input), len(p.process.input))
	}

	if len(progress.Output) != len(p.process.output) {
		return fmt.Errorf("output length mismatch (have: %d, want: %d)", len(progress.Output), len(p.process.output))
	}

	if progress.Size == 0 {
		progress.Size = progress.SizeKB * 1024
	}

	for i, io := range progress.Input {
		if io.Size == 0 {
			io.Size = io.SizeKB * 1024
		}

		progress.Input[i].Size = io.Size
	}

	for i, io := range progress.Output {
		if io.Size == 0 {
			io.Size = io.SizeKB * 1024
		}

		progress.Output[i].Size = io.Size
	}

	p.progress.ffmpeg = progress

	return nil
}

func (p *parser) parseAVstreamProgress(line string) error {
	progress := ffmpegAVstream{}

	err := json.Unmarshal([]byte(line), &progress)
	if err != nil {
		return err
	}

	p.progress.avstream[progress.Address] = progress

	return nil
}

func (p *parser) Progress() app.Progress {
	p.lock.progress.RLock()
	defer p.lock.progress.RUnlock()

	progress := p.process.export()

	p.progress.ffmpeg.exportTo(&progress)

	for i, io := range progress.Input {
		av, ok := p.progress.avstream[io.Address]
		if !ok {
			continue
		}

		progress.Input[i].AVstream = av.export()
	}

	return progress
}

func (p *parser) Prelude() []string {
	p.lock.prelude.RLock()
	if p.prelude.data == nil {
		p.lock.prelude.RUnlock()
		return []string{}
	}

	prelude := make([]string, len(p.prelude.data))
	copy(prelude, p.prelude.data)

	tail := []string{}

	p.prelude.tail.Do(func(l interface{}) {
		if l == nil {
			return
		}

		tail = append(tail, l.(string))
	})

	p.lock.prelude.RUnlock()

	if len(tail) != 0 {
		if p.prelude.truncatedLines > uint64(p.prelude.tailLines) {
			prelude = append(prelude, fmt.Sprintf("... truncated %d lines ...", p.prelude.truncatedLines-uint64(p.prelude.tailLines)))
		}
		prelude = append(prelude, tail...)
	}

	return prelude
}

func (p *parser) parsePrelude() bool {
	p.lock.progress.Lock()
	defer p.lock.progress.Unlock()

	data := p.Prelude()

	inputs, outputs, noutputs := prelude.Parse(data)

	if len(outputs) != noutputs {
		return false
	}

	for _, in := range inputs {
		io := ffmpegProcessIO{
			Address:  in.Address,
			Format:   in.Format,
			Index:    in.Index,
			Stream:   in.Stream,
			Type:     in.Type,
			Codec:    in.Codec,
			Pixfmt:   in.Pixfmt,
			Width:    in.Width,
			Height:   in.Height,
			Sampling: in.Sampling,
			Layout:   in.Layout,
		}

		if ip, _ := url.Lookup(io.Address); len(ip) != 0 {
			io.IP = ip
		}

		p.process.input = append(p.process.input, io)
	}

	for _, out := range outputs {
		io := ffmpegProcessIO{
			Address:  out.Address,
			Format:   out.Format,
			Index:    out.Index,
			Stream:   out.Stream,
			Type:     out.Type,
			Codec:    out.Codec,
			Pixfmt:   out.Pixfmt,
			Width:    out.Width,
			Height:   out.Height,
			Sampling: out.Sampling,
			Layout:   out.Layout,
		}

		if ip, _ := url.Lookup(io.Address); len(ip) != 0 {
			io.IP = ip
		}

		p.process.output = append(p.process.output, io)
	}

	return true
}

func (p *parser) addLog(line string) {
	p.lock.log.Lock()
	defer p.lock.log.Unlock()

	p.log.Value = process.Line{
		Timestamp: time.Now(),
		Data:      line,
	}
	p.log = p.log.Next()
}

func (p *parser) Log() []process.Line {
	var log = []process.Line{}

	p.lock.log.RLock()
	defer p.lock.log.RUnlock()

	p.log.Do(func(l interface{}) {
		if l == nil {
			return
		}

		log = append(log, l.(process.Line))
	})

	return log
}

func (p *parser) ResetStats() {
	p.lock.progress.Lock()
	defer p.lock.progress.Unlock()

	if p.averager.initialized {
		p.averager.main.stop()

		p.averager.main = averager{}

		for i := range p.averager.input {
			p.averager.input[i].stop()
		}

		p.averager.input = []averager{}

		for i := range p.averager.output {
			p.averager.output[i].stop()
		}

		p.averager.output = []averager{}

		p.averager.initialized = false
	}

	if p.stats.initialized {
		p.stats.main = stats{}

		p.stats.input = []stats{}
		p.stats.output = []stats{}

		p.stats.initialized = false
	}

	p.process = ffmpegProcess{}
	p.progress.ffmpeg = ffmpegProgress{}
	p.progress.avstream = make(map[string]ffmpegAVstream)

	p.lock.prelude.Lock()
	p.prelude.done = false
	p.lock.prelude.Unlock()
}

func (p *parser) ResetLog() {
	p.storeLogHistory()

	p.lock.prelude.Lock()
	p.prelude.data = []string{}
	p.prelude.tail = ring.New(p.prelude.tailLines)
	p.prelude.truncatedLines = 0
	p.prelude.done = false
	p.lock.prelude.Unlock()

	p.lock.log.Lock()
	p.log = ring.New(p.logLines)
	p.logStart = time.Now()
	p.lock.log.Unlock()
}

// Report represents a log report, including the prelude and the last log lines
// of the process.
type Report struct {
	CreatedAt time.Time
	Prelude   []string
	Log       []process.Line
}

func (p *parser) storeLogHistory() {
	if p.logHistory == nil {
		return
	}

	h := p.Report()

	if len(h.Prelude) != 0 {
		p.logHistory.Value = h
		p.logHistory = p.logHistory.Next()
	}
}

func (p *parser) Report() Report {
	h := Report{
		Prelude: p.Prelude(),
		Log:     p.Log(),
	}

	p.lock.log.RLock()
	h.CreatedAt = p.logStart
	p.lock.log.RUnlock()

	return h
}

func (p *parser) ReportHistory() []Report {
	var history = []Report{}

	p.logHistory.Do(func(l interface{}) {
		if l == nil {
			return
		}

		history = append(history, l.(Report))
	})

	return history
}

func (p *parser) TransferReportHistory(dst Parser) error {
	pp, ok := dst.(*parser)
	if !ok {
		return fmt.Errorf("the target parser is not of the required type")
	}

	p.logHistory.Do(func(l interface{}) {
		if l == nil {
			return
		}

		pp.logHistory.Value = l
		pp.logHistory = pp.logHistory.Next()
	})

	return nil
}
