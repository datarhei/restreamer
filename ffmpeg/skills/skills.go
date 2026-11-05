// Package skills provides an easy way to find out which
// codec, de/muxers, filters, and formats ffmpeg
// supports and which devices are available for ffmpeg on
// the system.
package skills

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Codec represents a codec with its availabe encoders and decoders
type Codec struct {
	Id       string
	Name     string
	Encoders []string
	Decoders []string
}

type ffCodecs struct {
	Audio    []Codec
	Video    []Codec
	Subtitle []Codec
}

// HWDevice represents a hardware device (e.g. USB device)
type HWDevice struct {
	Id    string
	Name  string
	Extra string
	Media string
}

// Device represents a type of device (e.g. V4L2) including connected actual devices
type Device struct {
	Id      string
	Name    string
	Devices []HWDevice
}

type ffDevices struct {
	Demuxers []Device
	Muxers   []Device
}

// Format represents a supported format (e.g. flv)
type Format struct {
	Id   string
	Name string
}

type ffFormats struct {
	Demuxers []Format
	Muxers   []Format
}

// Protocol represents a supported protocol (e.g. rtsp)
type Protocol struct {
	Id   string
	Name string
}

type ffProtocols struct {
	Input  []Protocol
	Output []Protocol
}

type HWAccel struct {
	Id   string
	Name string
}

// Filter represents a supported filter (e.g. anullsrc, test2)
type Filter struct {
	Id   string
	Name string
}

// Library represents a linked av library
type Library struct {
	Name     string
	Compiled string
	Linked   string
}

type ffmpeg struct {
	Version       string
	Compiler      string
	Configuration string
	Libraries     []Library
}

// Skills are the detected capabilities of a ffmpeg binary
type Skills struct {
	FFmpeg ffmpeg

	Filters  []Filter
	HWAccels []HWAccel

	Codecs    ffCodecs
	Devices   ffDevices
	Formats   ffFormats
	Protocols ffProtocols
}

// New returns all skills that ffmpeg provides
func New(binary string) (Skills, error) {
	c := Skills{}

	ffmpeg, err := version(binary)
	if len(ffmpeg.Version) == 0 || err != nil {
		if err != nil {
			return Skills{}, fmt.Errorf("can't parse ffmpeg version info: %w", err)
		}

		return Skills{}, fmt.Errorf("can't parse ffmpeg version info")
	}

	c.FFmpeg = ffmpeg

	c.Filters = filters(binary)
	c.HWAccels = hwaccels(binary)

	c.Codecs = codecs(binary)
	c.Formats = formats(binary)
	c.Devices = devices(binary)
	c.Protocols = protocols(binary)

	return c, nil
}

func version(binary string) (ffmpeg, error) {
	cmd := exec.Command(binary, "-version")
	cmd.Env = []string{}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return ffmpeg{}, err
	}

	return parseVersion(out), nil
}

func parseVersion(data []byte) ffmpeg {
	f := ffmpeg{}

	reVersion := regexp.MustCompile(`^ffmpeg version ([0-9]+\.[0-9]+(\.[0-9]+)?)`)
	reCompiler := regexp.MustCompile(`(?m)^\s*built with (.*)$`)
	reConfiguration := regexp.MustCompile(`(?m)^\s*configuration: (.*)$`)
	reLibrary := regexp.MustCompile(`(?m)^\s*(lib(?:[a-z]+))\s+([0-9]+\.\s*[0-9]+\.\s*[0-9]+) /\s+([0-9]+\.\s*[0-9]+\.\s*[0-9]+)`)

	if matches := reVersion.FindSubmatch(data); matches != nil {
		f.Version = string(matches[1])
		if len(matches[2]) == 0 {
			f.Version = f.Version + ".0"
		}
	}

	if matches := reCompiler.FindSubmatch(data); matches != nil {
		f.Compiler = string(matches[1])
	}

	if matches := reConfiguration.FindSubmatch(data); matches != nil {
		f.Configuration = string(matches[1])
	}

	for _, matches := range reLibrary.FindAllSubmatch(data, -1) {
		l := Library{
			Name:     string(matches[1]),
			Compiled: string(matches[2]),
			Linked:   string(matches[3]),
		}

		f.Libraries = append(f.Libraries, l)
	}

	return f
}

func filters(binary string) []Filter {
	cmd := exec.Command(binary, "-filters")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseFilters(stdout)
}

func parseFilters(data []byte) []Filter {
	filters := []Filter{}

	re := regexp.MustCompile(`^\s[TSC.]{3} ([0-9A-Za-z_]+)\s+(?:.*?)\s+(.*)?$`)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		filters = append(filters, Filter{
			Id:   matches[1],
			Name: matches[2],
		})
	}

	return filters
}

func codecs(binary string) ffCodecs {
	cmd := exec.Command(binary, "-codecs")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseCodecs(stdout)
}

func parseCodecs(data []byte) ffCodecs {
	codecs := ffCodecs{}

	re := regexp.MustCompile(`^\s([D.])([E.])([VAS]).{3} ([0-9A-Za-z_]+)\s+(.*?)(?:\(decoders:([^\)]+)\))?\s?(?:\(encoders:([^\)]+)\))?$`)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		codec := Codec{
			Id:   matches[4],
			Name: strings.TrimSpace(matches[5]),
		}

		if matches[1] == "D" {
			if len(matches[6]) == 0 {
				codec.Decoders = append(codec.Decoders, matches[4])
			} else {
				codec.Decoders = strings.Split(strings.TrimSpace(matches[6]), " ")
			}
		}

		if matches[2] == "E" {
			if len(matches[7]) == 0 {
				codec.Encoders = append(codec.Encoders, matches[4])
			} else {
				codec.Encoders = strings.Split(strings.TrimSpace(matches[7]), " ")
			}
		}

		switch matches[3] {
		case "V":
			codecs.Video = append(codecs.Video, codec)
		case "A":
			codecs.Audio = append(codecs.Audio, codec)
		case "S":
			codecs.Subtitle = append(codecs.Subtitle, codec)
		}
	}

	return codecs
}

func formats(binary string) ffFormats {
	cmd := exec.Command(binary, "-formats")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseFormats(stdout)
}

func parseFormats(data []byte) ffFormats {
	formats := ffFormats{}

	re := regexp.MustCompile(`^\s([D ])([E ]) ([0-9A-Za-z_,]+)\s+(.*?)$`)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		id := strings.Split(matches[3], ",")[0]

		format := Format{
			Id:   id,
			Name: matches[4],
		}

		if matches[1] == "D" {
			formats.Demuxers = append(formats.Demuxers, format)
		}

		if matches[2] == "E" {
			formats.Muxers = append(formats.Muxers, format)
		}
	}

	return formats
}

func devices(binary string) ffDevices {
	cmd := exec.Command(binary, "-devices")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseDevices(stdout, binary)
}

func parseDevices(data []byte, binary string) ffDevices {
	devices := ffDevices{}

	re := regexp.MustCompile(`^\s([D ])([E ]) ([0-9A-Za-z_,]+)\s+(.*?)$`)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		id := strings.Split(matches[3], ",")[0]

		device := Device{
			Id:   id,
			Name: matches[4],
		}

		switch id {
		case "avfoundation":
			device.Devices, _ = DevicesAvfoundation(binary)
		case "alsa":
			device.Devices, _ = DevicesALSA()
		case "video4linux2":
			device.Devices, _ = DevicesV4L()
		case "fbdev":
			device.Devices, _ = DevicesFramebuffer()
		}

		if matches[1] == "D" {
			devices.Demuxers = append(devices.Demuxers, device)
		}

		if matches[2] == "E" {
			devices.Muxers = append(devices.Muxers, device)
		}
	}

	return devices
}

func protocols(binary string) ffProtocols {
	cmd := exec.Command(binary, "-protocols")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseProtocols(stdout)
}

func parseProtocols(data []byte) ffProtocols {
	protocols := ffProtocols{}

	mode := ""

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "Input:" {
			mode = "input"
			continue
		} else if line == "Output:" {
			mode = "output"
			continue
		}

		if len(mode) == 0 {
			continue
		}

		id := strings.TrimSpace(line)

		protocol := Protocol{
			Id:   id,
			Name: id,
		}

		if mode == "input" {
			protocols.Input = append(protocols.Input, protocol)
		} else if mode == "output" {
			protocols.Output = append(protocols.Output, protocol)
		}
	}

	return protocols
}

func hwaccels(binary string) []HWAccel {
	cmd := exec.Command(binary, "-hwaccels")
	cmd.Env = []string{}

	stdout, _ := cmd.Output()

	return parseHWAccels(stdout)
}

func parseHWAccels(data []byte) []HWAccel {
	hwaccels := []HWAccel{}

	re := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	start := false

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "Hardware acceleration methods:" {
			start = true
			continue
		}

		if !start {
			continue
		}

		if !re.MatchString(line) {
			continue
		}

		id := strings.TrimSpace(line)

		hwaccels = append(hwaccels, HWAccel{
			Id:   id,
			Name: id,
		})
	}

	return hwaccels
}
