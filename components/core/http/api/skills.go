package api

import (
	"github.com/datarhei/core/v16/ffmpeg/skills"
)

// SkillsFilter represents an ffmpeg filter
type SkillsFilter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Unmarshal converts a skills filter to its API representation
func (s *SkillsFilter) Unmarshal(filter skills.Filter) {
	s.ID = filter.Id
	s.Name = filter.Name
}

// SkillsHWAccel represents an ffmpeg HW accelerator
type SkillsHWAccel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Unmarshal converts a skills HWAccel to its API representation
func (s *SkillsHWAccel) Unmarshal(hwaccel skills.HWAccel) {
	s.ID = hwaccel.Id
	s.Name = hwaccel.Name
}

// SkillsCodec represents an ffmpeg codec
type SkillsCodec struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Encoders []string `json:"encoders"`
	Decoders []string `json:"decoders"`
}

// Unmarshal converts a skills codec to its API representation
func (s *SkillsCodec) Unmarshal(codec skills.Codec) {
	s.ID = codec.Id
	s.Name = codec.Name
	s.Encoders = codec.Encoders
	s.Decoders = codec.Decoders
}

// SkillsHWDevice represents an ffmpeg hardware device
type SkillsHWDevice struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Extra string `json:"extra"`
	Media string `json:"media"`
}

// Unmarshal converts a skills HW device to its API representation
func (s *SkillsHWDevice) Unmarshal(hwdevice skills.HWDevice) {
	s.ID = hwdevice.Id
	s.Name = hwdevice.Name
	s.Extra = hwdevice.Extra
	s.Media = hwdevice.Media
}

// SkillsDevice represents a group of ffmpeg hardware devices
type SkillsDevice struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Devices []SkillsHWDevice `json:"devices"`
}

// Unmarshal converts a skills device to its API representation
func (s *SkillsDevice) Unmarshal(device skills.Device) {
	s.ID = device.Id
	s.Name = device.Name
	s.Devices = make([]SkillsHWDevice, len(device.Devices))

	for i, dev := range device.Devices {
		s.Devices[i].Unmarshal(dev)
	}
}

// SkillsFormat represents an ffmpeg format
type SkillsFormat struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Unmarshal converts a skills format to its API representation
func (s *SkillsFormat) Unmarshal(format skills.Format) {
	s.ID = format.Id
	s.Name = format.Name
}

// SkillsProtocol represents an ffmpeg protocol
type SkillsProtocol struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Unmarshal converts a skills protocol to its API representation
func (s *SkillsProtocol) Unmarshal(proto skills.Protocol) {
	s.ID = proto.Id
	s.Name = proto.Name
}

// SkillsLibrary represents an avlib ffmpeg is compiled and linked with
type SkillsLibrary struct {
	Name     string `json:"name"`
	Compiled string `json:"compiled"`
	Linked   string `json:"linked"`
}

// Unmarshal converts a skills library to its API representation
func (s *SkillsLibrary) Unmarshal(lib skills.Library) {
	s.Name = lib.Name
	s.Compiled = lib.Compiled
	s.Linked = lib.Linked
}

// Skills represents a set of ffmpeg capabilities
type Skills struct {
	FFmpeg struct {
		Version       string          `json:"version"`
		Compiler      string          `json:"compiler"`
		Configuration string          `json:"configuration"`
		Libraries     []SkillsLibrary `json:"libraries"`
	} `json:"ffmpeg"`

	Filters  []SkillsFilter  `json:"filter"`
	HWAccels []SkillsHWAccel `json:"hwaccels"`

	Codecs struct {
		Audio    []SkillsCodec `json:"audio"`
		Video    []SkillsCodec `json:"video"`
		Subtitle []SkillsCodec `json:"subtitle"`
	} `json:"codecs"`

	Devices struct {
		Demuxers []SkillsDevice `json:"demuxers"`
		Muxers   []SkillsDevice `json:"muxers"`
	} `json:"devices"`

	Formats struct {
		Demuxers []SkillsFormat `json:"demuxers"`
		Muxers   []SkillsFormat `json:"muxers"`
	} `json:"formats"`

	Protocols struct {
		Input  []SkillsProtocol `json:"input"`
		Output []SkillsProtocol `json:"output"`
	} `json:"protocols"`
}

// Unmarshal converts skills to their API representation
func (s *Skills) Unmarshal(skills skills.Skills) {
	s.FFmpeg.Version = skills.FFmpeg.Version
	s.FFmpeg.Compiler = skills.FFmpeg.Compiler
	s.FFmpeg.Configuration = skills.FFmpeg.Configuration
	s.FFmpeg.Libraries = make([]SkillsLibrary, len(skills.FFmpeg.Libraries))
	for id, x := range skills.FFmpeg.Libraries {
		s.FFmpeg.Libraries[id].Unmarshal(x)
	}

	s.Filters = make([]SkillsFilter, len(skills.Filters))
	for id, x := range skills.Filters {
		s.Filters[id].Unmarshal(x)
	}

	s.HWAccels = make([]SkillsHWAccel, len(skills.HWAccels))
	for id, x := range skills.HWAccels {
		s.HWAccels[id].Unmarshal(x)
	}

	s.Codecs.Audio = make([]SkillsCodec, len(skills.Codecs.Audio))
	for id, x := range skills.Codecs.Audio {
		s.Codecs.Audio[id].Unmarshal(x)
	}

	s.Codecs.Video = make([]SkillsCodec, len(skills.Codecs.Video))
	for id, x := range skills.Codecs.Video {
		s.Codecs.Video[id].Unmarshal(x)
	}

	s.Codecs.Subtitle = make([]SkillsCodec, len(skills.Codecs.Subtitle))
	for id, x := range skills.Codecs.Subtitle {
		s.Codecs.Subtitle[id].Unmarshal(x)
	}

	s.Devices.Demuxers = make([]SkillsDevice, len(skills.Devices.Demuxers))
	for id, x := range skills.Devices.Demuxers {
		s.Devices.Demuxers[id].Unmarshal(x)
	}

	s.Devices.Muxers = make([]SkillsDevice, len(skills.Devices.Muxers))
	for id, x := range skills.Devices.Muxers {
		s.Devices.Muxers[id].Unmarshal(x)
	}

	s.Formats.Demuxers = make([]SkillsFormat, len(skills.Formats.Demuxers))
	for id, x := range skills.Formats.Demuxers {
		s.Formats.Demuxers[id].Unmarshal(x)
	}

	s.Formats.Muxers = make([]SkillsFormat, len(skills.Formats.Muxers))
	for id, x := range skills.Formats.Muxers {
		s.Formats.Muxers[id].Unmarshal(x)
	}

	s.Protocols.Input = make([]SkillsProtocol, len(skills.Protocols.Input))
	for id, x := range skills.Protocols.Input {
		s.Protocols.Input[id].Unmarshal(x)
	}

	s.Protocols.Output = make([]SkillsProtocol, len(skills.Protocols.Output))
	for id, x := range skills.Protocols.Output {
		s.Protocols.Output[id].Unmarshal(x)
	}
}
