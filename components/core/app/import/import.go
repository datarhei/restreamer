package main

// TODO: Sonderzeichen im userpass von importierten URLs
// TODO: import von internal RTMP (external.stream), wenn jemand z.B. von OBS reinschiesst

import (
	gojson "encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/datarhei/core/v16/encoding/json"
	"github.com/datarhei/core/v16/ffmpeg"
	"github.com/datarhei/core/v16/ffmpeg/skills"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/restream"
	"github.com/datarhei/core/v16/restream/app"
	"github.com/datarhei/core/v16/restream/store"

	"github.com/google/uuid"
)

type storeDataV1 struct {
	Addresses struct {
		Output  string `json:"optionalOutputAddress"`
		Input   string `json:"srcAddress"`
		Streams struct {
			Audio *storeDataV1Stream `json:"audio"`
			Video *storeDataV1Stream `json:"video"`
		} `json:"srcStreams"`
	} `json:"addresses"`

	Options struct {
		Audio struct {
			Bitrate  string `json:"bitrate"`
			Channels string `json:"channels"`
			Codec    string `json:"codec"`
			Preset   string `json:"preset"`
			Sampling string `json:"sampling"`
		} `json:"audio"`

		Output struct {
			HLS struct {
				ListSize string `json:"listSize"`
				Method   string `json:"method"`
				Time     string `json:"time"`
				Timeout  string `json:"timeout"`
			} `json:"hls"`
			Type string `json:"type"`
		} `json:"output"`

		Player struct {
			Autoplay   bool   `json:"autoplay"`
			Color      string `json:"color"`
			Mute       bool   `json:"mute"`
			Statistics bool   `json:"statistics"`
			Logo       struct {
				Image    string `json:"image"`
				Link     string `json:"link"`
				Position string `json:"position"`
			} `json:"logo"`
		} `json:"player"`

		RTSPTCP bool `json:"rtspTcp"`

		Video struct {
			Bitrate string `json:"bitrate"`
			Codec   string `json:"codec"`
			FPS     string `json:"fps"`
			GOP     string `json:"gop"`
			Preset  string `json:"preset"`
			Profile string `json:"profile"`
			Tune    string `json:"tune"`
		}
	} `json:"options"`

	States struct {
		Ingest struct {
			Type string `json:"type"`
		} `json:"repeatToLocalNginx"`
		Egress struct {
			Type string `json:"type"`
		} `json:"repeatToOptionalOutput"`
	} `json:"states"`

	Actions struct {
		Ingest string `json:"repeatToLocalNginx"`
		Egress string `json:"repeatToOptionalOutput"`
	} `json:"userActions"`
}

type storeDataV1Stream struct {
	Index int    `json:"index"`
	Type  string `json:"type"`
	Codec string `json:"codec"`

	// Video
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Format string `json:"format,omitempty"`

	// Audio
	Layout   string `json:"layout,omitempty"`
	Channels int    `json:"channels,omitempty"`
	Sampling string `json:"sampling,omitempty"`
}

type restreamerUISourceSettingsNetwork struct {
	Address string `json:"address"`
	General struct {
		Flags           []string `json:"fflags"`
		ThreadQueueSize int      `json:"thread_queue_size"`
	} `json:"general"`
	HTTP struct {
		ForceFramerate bool `json:"forceFramerate"`
		Framerate      int  `json:"framerate"`
		ReadNative     bool `json:"readNative"`
	} `json:"http"`
	Mode string `json:"mode"`
	RTSP struct {
		STimeout int  `json:"stimeout"`
		UDP      bool `json:"udp"`
	} `json:"rtsp"`
}

type restreamerUISourceSettingsV4L struct {
	Device    string `json:"device"`
	Format    string `json:"format"`
	Framerate int    `json:"framerate"`
	Size      string `json:"size"`
}

type restreamerUISourceSettingsVirtualAudio struct {
	Amplitude  int    `json:"amplitude"`
	BeepFactor int    `json:"beepfactor"`
	Color      string `json:"color"`
	Frequency  int    `json:"frequency"`
	Layout     string `json:"layout"`
	Sampling   int    `json:"sampling"`
	Source     string `json:"source"`
}

type restreamerUISourceSettingsALSA struct {
	Address  string `json:"address"`
	Device   string `json:"device"`
	Sampling int    `json:"sampling"`
	Channels int    `json:"channels"`
	Delay    int    `json:"delay"`
}

type restreamerUISourceStream struct {
	Bitrate     int    `json:"bitrate_kbps"`
	Channels    int    `json:"channels"`
	Codec       string `json:"codec"`
	Coder       string `json:"coder"`
	Duration    int    `json:"duration_sec"`
	Format      string `json:"format"`
	FPS         int    `json:"fps"`
	Height      int    `json:"height"`
	Index       int    `json:"index"`
	Language    string `json:"language"`
	Layout      string `json:"layout"`
	PixelFormat string `json:"pix_fmt"`
	Sampling    int    `json:"sampling_hz"`
	Stream      int    `json:"stream"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Width       int    `json:"width"`
}

type restreamerUISourceInput struct {
	Address string   `json:"address"`
	Options []string `json:"options"`
}

type restreamerUISource struct {
	Inputs   []restreamerUISourceInput  `json:"inputs"`
	Settings interface{}                `json:"settings"`
	Streams  []restreamerUISourceStream `json:"streams"`
	Type     string                     `json:"type"`
}

type restreamerUIPlayer struct {
	Autoplay bool `json:"autoplay"`
	Color    struct {
		Buttons string `json:"buttons"`
		Seekbar string `json:"seekbar"`
	} `json:"color"`
	GA struct {
		Account string `json:"account"`
		Name    string `json:"name"`
	} `json:"ga"`
	Logo struct {
		Image    string `json:"image"`
		Link     string `json:"link"`
		Position string `json:"position"`
	} `json:"logo"`
	Mute       bool `json:"mute"`
	Statistics bool `json:"statistics"`
}

type restreamerUIProfileCoderSettingsAAC struct {
	Bitrate  string `json:"bitrate"`
	Channels string `json:"channels"`
	Layout   string `json:"layout"`
	Sampling string `json:"sampling"`
}

type restreamerUIProfileCoderSettingsX264 struct {
	Bitrate string `json:"bitrate"`
	FPS     string `json:"fps"`
	Preset  string `json:"preset"`
	Profile string `json:"profile"`
	Tune    string `json:"tune"`
}

type restreamerUIProfileCoderSettingsCopy struct{}

type restreamerUIProfileAV struct {
	Coder    string      `json:"coder"`
	Codec    string      `json:"codec"`
	Mapping  []string    `json:"mapping"`
	Settings interface{} `json:"settings"`
}

type restreamerUIProfile struct {
	Audio struct {
		Encoder restreamerUIProfileAV `json:"encoder"`
		Decoder restreamerUIProfileAV `json:"decoder"`
		Source  int                   `json:"source"`
		Stream  int                   `json:"stream"`
	} `json:"audio"`
	Video struct {
		Encoder restreamerUIProfileAV `json:"encoder"`
		Decoder restreamerUIProfileAV `json:"decoder"`
		Source  int                   `json:"source"`
		Stream  int                   `json:"stream"`
	} `json:"video"`
}

type restreamerUIIngest struct {
	Control struct {
		HLS struct {
			LHLS            bool `json:"lhls"`
			SegmentDuration int  `json:"segmentDuration"`
			ListSize        int  `json:"listSize"`
		} `json:"hls"`
		Process struct {
			Autostart    bool `json:"autostart"`
			Reconnect    bool `json:"reconnect"`
			Delay        int  `json:"delay"`
			StaleTimeout int  `json:"staleTimeout"`
		} `json:"process"`
		Snapshot struct {
			Enable   bool `json:"enable"`
			Interval int  `json:"interval"`
		} `json:"snapshot"`
	} `json:"control"`
	License string `json:"license"`
	Meta    struct {
		Author struct {
			Description string `json:"description"`
			Name        string `json:"name"`
		} `json:"author"`
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"meta"`
	Player   restreamerUIPlayer    `json:"player"`
	Profiles []restreamerUIProfile `json:"profiles"`
	Sources  []restreamerUISource  `json:"sources"`
	Version  int                   `json:"version"`
	Imported bool                  `json:"imported"`
}

type restreamerUIEgressSettingsRTMP struct {
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
}

type restreamerUIEgressSettingsHLS struct {
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Options  struct {
		DeleteThreshold   string   `json:"hls_delete_threshold"`
		Flags             []string `json:"hls_flags"`
		InitTime          string   `json:"hls_init_time"`
		ListSize          string   `json:"hls_list_size"`
		SegmentType       string   `json:"hls_segment_type"`
		Time              string   `json:"hls_time"`
		Method            string   `json:"method"`
		StartNumber       string   `json:"start_number"`
		Timeout           string   `json:"timeout"`
		StartNumberSource string   `json:"hls_start_number_source"`
		AllowCache        string   `json:"hls_allow_cache"`
		Enc               string   `json:"hls_enc"`
		IgnoreIOErrors    string   `json:"ignore_io_errors"`
		HTTPersistent     string   `json:"http_persistent"`
	} `json:"options"`
}

type restreamerUIEgress struct {
	Control struct {
		Process struct {
			Autostart    bool `json:"autostart"`
			Reconnect    bool `json:"reconnect"`
			Delay        int  `json:"delay"`
			StaleTimeout int  `json:"staleTimeout"`
		} `json:"process"`
	} `json:"control"`
	Name   string `json:"name"`
	Output struct {
		Address string   `json:"address"`
		Options []string `json:"options"`
	} `json:"output"`
	Settings interface{} `json:"settings"`
	Version  int         `json:"version"`
}

func importSnapshotInterval(value string, defval int) int {
	interval, err := strconv.Atoi(value)
	if err == nil {
		return interval / 1000
	}

	duration, err := time.ParseDuration(value)
	if err == nil {
		return int(duration.Seconds())
	}

	return defval
}

var v1Environment = map[string]struct {
	value  string
	defval string
}{
	"RS_SNAPSHOT_INTERVAL": {defval: "60s"},
	"RS_MODE":              {defval: ""},

	"RS_INPUTSTREAM":  {defval: ""},
	"RS_OUTPUTSTREAM": {defval: ""},

	"RS_USBCAM_VIDEODEVICE":   {defval: "/dev/video"},
	"RS_USBCAM_FPS":           {defval: "25"},
	"RS_USBCAM_GOP":           {defval: "50"},
	"RS_USBCAM_BITRATE":       {defval: "5000000"},
	"RS_USBCAM_H264PRESET":    {defval: "ultrafast"},
	"RS_USBCAM_H264PROFILE":   {defval: "baseline"},
	"RS_USBCAM_WIDTH":         {defval: "1280"},
	"RS_USBCAM_HEIGHT":        {defval: "720"},
	"RS_USBCAM_AUDIO":         {defval: "false"},
	"RS_USBCAM_AUDIODEVICE":   {defval: "0"},
	"RS_USBCAM_AUDIOBITRATE":  {defval: "64000"},
	"RS_USBCAM_AUDIOCHANNELS": {defval: "1"},
	"RS_USBCAM_AUDIOLAYOUT":   {defval: "mono"},
	"RS_USBCAM_AUDIOSAMPLING": {defval: "44100"},

	"RS_RASPICAM_FORMAT":        {defval: "h264"},
	"RS_RASPICAM_FPS":           {defval: "25"},
	"RS_RASPICAM_WIDTH":         {defval: "1920"},
	"RS_RASPICAM_HEIGHT":        {defval: "1080"},
	"RS_RASPICAM_AUDIO":         {defval: "false"},
	"RS_RASPICAM_AUDIODEVICE":   {defval: "0"},
	"RS_RASPICAM_AUDIOBITRATE":  {defval: "64000"},
	"RS_RASPICAM_AUDIOCHANNELS": {defval: "1"},
	"RS_RASPICAM_AUDIOLAYOUT":   {defval: "mono"},
	"RS_RASPICAM_AUDIOSAMPLING": {defval: "44100"},
}

func initV1Environment() {
	for key, env := range v1Environment {
		value, ok := os.LookupEnv(key)
		if !ok {
			value = v1Environment[key].defval
		}

		env.value = value

		v1Environment[key] = env
	}
}

func getV1Environment(key string) string {
	value := v1Environment[key]

	return value.value
}

func importConfigFromEnvironment() importConfig {
	initV1Environment()

	c := importConfig{
		id:               uuid.New().String(),
		snapshotInterval: importSnapshotInterval(getV1Environment("RS_SNAPSHOT_INTERVAL"), 0),

		inputstream:  getV1Environment("RS_INPUTSTREAM"),
		outputstream: getV1Environment("RS_OUTPUTSTREAM"),
	}

	mode := getV1Environment("RS_MODE")

	if mode == "USBCAM" {
		c.usbcam.enable = true
		c.usbcam.device = getV1Environment("RS_USBCAM_VIDEODEVICE")
		c.usbcam.fps = getV1Environment("RS_USBCAM_FPS")
		c.usbcam.gop = getV1Environment("RS_USBCAM_GOP")
		c.usbcam.bitrate = getV1Environment("RS_USBCAM_BITRATE")
		c.usbcam.preset = getV1Environment("RS_USBCAM_H264PRESET")
		c.usbcam.profile = getV1Environment("RS_USBCAM_H264PROFILE")
		c.usbcam.width = getV1Environment("RS_USBCAM_WIDTH")
		c.usbcam.height = getV1Environment("RS_USBCAM_HEIGHT")

		if getV1Environment("RS_USBCAM_AUDIO") == "true" {
			c.audio.enable = true
			c.audio.device = getV1Environment("RS_USBCAM_AUDIODEVICE")
			c.audio.bitrate = getV1Environment("RS_USBCAM_AUDIOBITRATE")
			c.audio.channels = getV1Environment("RS_USBCAM_AUDIOCHANNELS")
			c.audio.layout = getV1Environment("RS_USBCAM_AUDIOLAYOUT")
			c.audio.sampling = getV1Environment("RS_USBCAM_AUDIOSAMPLING")
		}
	} else if mode == "RASPICAM" {
		c.raspicam.enable = true

		devices, _ := skills.DevicesV4L()
		for _, device := range devices {
			if strings.Contains(device.Extra, "bcm2835-v4l2") {
				c.raspicam.device = device.Id
				break
			}
		}

		c.raspicam.format = getV1Environment("RS_RASPICAM_FORMAT")
		c.raspicam.fps = getV1Environment("RS_RASPICAM_FPS")
		c.raspicam.width = getV1Environment("RS_RASPICAM_WIDTH")
		c.raspicam.height = getV1Environment("RS_RASPICAM_HEIGHT")

		if getV1Environment("RS_RASPICAM_AUDIO") == "true" {
			c.audio.enable = true
			c.audio.device = getV1Environment("RS_RASPICAM_AUDIODEVICE")
			c.audio.bitrate = getV1Environment("RS_RASPICAM_AUDIOBITRATE")
			c.audio.channels = getV1Environment("RS_RASPICAM_AUDIOCHANNELS")
			c.audio.layout = getV1Environment("RS_RASPICAM_AUDIOLAYOUT")
			c.audio.sampling = getV1Environment("RS_RASPICAM_AUDIOSAMPLING")
		}
	}

	return c
}

type importConfig struct {
	id               string
	snapshotInterval int

	binary string

	inputstream  string
	outputstream string

	usbcam   importConfigUSBCam
	raspicam importConfigRaspiCam
	audio    importConfigAudio
}

type importConfigUSBCam struct {
	enable  bool
	device  string
	fps     string
	gop     string
	bitrate string
	preset  string
	profile string
	width   string
	height  string
}

type importConfigRaspiCam struct {
	enable bool
	device string
	format string
	fps    string
	width  string
	height string
}

type importConfigAudio struct {
	enable   bool
	device   string
	bitrate  string
	channels string
	layout   string
	sampling string
}

func importV1(fs fs.Filesystem, path string, cfg importConfig) (store.StoreData, error) {
	if len(cfg.id) == 0 {
		cfg.id = uuid.New().String()
	}

	r := store.NewStoreData()

	jsondata, err := fs.ReadFile(path)
	if err != nil {
		return r, fmt.Errorf("failed to read data from %s: %w", path, err)
	}

	var v1data storeDataV1

	if err := gojson.Unmarshal(jsondata, &v1data); err != nil {
		return r, json.FormatError(jsondata, err)
	}

	if len(v1data.Addresses.Input) == 0 {
		if len(cfg.inputstream) != 0 {
			v1data.Addresses.Input = cfg.inputstream
			v1data.Addresses.Streams.Video = nil
			v1data.Addresses.Streams.Audio = nil
			v1data.States.Ingest.Type = "connected"
			v1data.Actions.Ingest = "start"

			if len(v1data.Addresses.Output) == 0 {
				v1data.Addresses.Output = cfg.outputstream
				v1data.States.Egress.Type = "connected"
				v1data.Actions.Egress = "stop"
			}
		}
	}

	if cfg.usbcam.enable || cfg.raspicam.enable {
		if cfg.usbcam.enable {
			if v1data.Options.Video.Codec == "copy" {
				v1data.Options.Video.Codec = "h264"
				v1data.Options.Video.Preset = cfg.usbcam.preset

				if bitrate, err := strconv.Atoi(cfg.usbcam.bitrate); err == nil {
					v1data.Options.Video.Bitrate = strconv.Itoa(bitrate / 1000)
				} else {
					v1data.Options.Video.Bitrate = "2048"
				}

				v1data.Options.Video.FPS = cfg.usbcam.fps
				v1data.Options.Video.Profile = cfg.usbcam.profile
				v1data.Options.Video.Tune = "zerolatency"
			}
		}

		if v1data.Options.Audio.Codec == "silence" || v1data.Options.Audio.Codec == "none" {
			cfg.audio.enable = false
		}

		if cfg.audio.enable {
			v1data.Options.Audio.Codec = "aac"
			v1data.Options.Audio.Preset = "device"

			if bitrate, err := strconv.Atoi(cfg.audio.bitrate); err == nil {
				v1data.Options.Audio.Bitrate = strconv.Itoa(bitrate / 1000)
			} else {
				v1data.Options.Audio.Bitrate = "64"
			}

			v1data.Options.Audio.Channels = "mono"
			if cfg.audio.channels == "2" {
				v1data.Options.Audio.Channels = "stereo"
			}

			v1data.Options.Audio.Sampling = cfg.audio.sampling
		}
	}

	process := &app.Process{
		ID:        "restreamer-ui:ingest:" + cfg.id,
		Reference: cfg.id,
		CreatedAt: time.Now().Unix(),
		Order:     "stop",
	}

	if v1data.Actions.Ingest == "start" {
		process.Order = "start"
	}

	config := &app.Config{
		ID:             "restreamer-ui:ingest:" + cfg.id,
		Reference:      cfg.id,
		Reconnect:      true,
		ReconnectDelay: 15,
		Autostart:      true,
		StaleTimeout:   30,
		Options: []string{
			"-err_detect",
			"ignore_err",
		},
	}

	// UI Settings

	ui := restreamerUIIngest{
		License:  "CC BY 4.0",
		Player:   restreamerUIPlayer{},
		Profiles: []restreamerUIProfile{},
		Sources:  []restreamerUISource{},
		Version:  1,
		Imported: true,
	}

	ui.Control.HLS.LHLS = false
	ui.Control.HLS.SegmentDuration = 2
	ui.Control.HLS.ListSize = 6

	ui.Control.Process.Autostart = true
	ui.Control.Process.Reconnect = true
	ui.Control.Process.Delay = 15
	ui.Control.Process.StaleTimeout = 30

	ui.Control.Snapshot.Enable = false
	ui.Control.Snapshot.Interval = 60
	if cfg.snapshotInterval != 0 {
		ui.Control.Snapshot.Enable = true
		ui.Control.Snapshot.Interval = cfg.snapshotInterval
	}

	ui.Meta.Author.Description = ""
	ui.Meta.Author.Name = ""

	ui.Meta.Description = "Live from earth. Powered by datarhei Restreamer."
	ui.Meta.Name = "Livestream"

	ui.Player.Autoplay = v1data.Options.Player.Autoplay
	ui.Player.Color.Buttons = v1data.Options.Player.Color
	ui.Player.Color.Seekbar = v1data.Options.Player.Color
	ui.Player.Logo.Image = v1data.Options.Player.Logo.Image
	ui.Player.Logo.Link = v1data.Options.Player.Logo.Link
	ui.Player.Logo.Position = v1data.Options.Player.Logo.Position
	ui.Player.Mute = v1data.Options.Player.Mute
	ui.Player.Statistics = v1data.Options.Player.Statistics

	if cfg.usbcam.enable || cfg.raspicam.enable {
		source := restreamerUISource{
			Type: "video4linux2",
		}

		settings := restreamerUISourceSettingsV4L{}

		if cfg.usbcam.enable {
			settings.Device = cfg.usbcam.device
			settings.Format = "nv12"
		} else {
			settings.Device = cfg.raspicam.device
			settings.Format = cfg.raspicam.format
		}

		if x, err := strconv.Atoi(cfg.usbcam.fps); err == nil {
			settings.Framerate = x
		}

		var width int = 0
		var height int = 0

		if w, err := strconv.Atoi(cfg.usbcam.width); err == nil {
			width = w
			if h, err := strconv.Atoi(cfg.usbcam.height); err == nil {
				height = h
				settings.Size = strconv.Itoa(w) + "x" + strconv.Itoa(h)
			}
		}

		if len(settings.Size) == 0 {
			settings.Size = "1280x720"
		}

		source.Settings = settings

		input := restreamerUISourceInput{
			Address: settings.Device,
			Options: []string{
				"-thread_queue_size", "512",
				"-f", "v4l2",
				"-framerate", strconv.Itoa(settings.Framerate),
				"-video_size", settings.Size,
				"-input_format", settings.Format,
			},
		}

		source.Inputs = append(source.Inputs, input)

		source.Streams = append(source.Streams, restreamerUISourceStream{
			URL:         settings.Device,
			Type:        "video",
			Index:       0,
			Stream:      0,
			Codec:       "rawvideo",
			Width:       width,
			Height:      height,
			PixelFormat: settings.Format,
		})

		ui.Sources = append(ui.Sources, source)

		if cfg.audio.enable {
			source := restreamerUISource{
				Type: "alsa",
			}

			settings := restreamerUISourceSettingsALSA{
				Address: "hw:" + cfg.audio.device,
				Device:  cfg.audio.device,
			}

			re := regexp.MustCompile("([0-9]+),([0-9]+)")
			if !re.MatchString(cfg.audio.device) {
				settings.Address = "hw:1,0"
				settings.Device = "1,0"
			}

			if x, err := strconv.Atoi(cfg.audio.sampling); err == nil {
				settings.Sampling = x
			}

			if x, err := strconv.Atoi(cfg.audio.channels); err == nil {
				settings.Channels = x
			}

			source.Settings = settings

			input := restreamerUISourceInput{
				Address: settings.Address,
				Options: []string{
					"-thread_queue_size", "512",
					"-f", "alsa",
					"-ac", cfg.audio.channels,
					"-ar", cfg.audio.sampling,
				},
			}

			source.Inputs = append(source.Inputs, input)

			source.Streams = append(source.Streams, restreamerUISourceStream{
				URL:      settings.Address,
				Type:     "audio",
				Index:    0,
				Stream:   0,
				Codec:    "pcm_u8",
				Layout:   cfg.audio.layout,
				Channels: settings.Channels,
				Sampling: settings.Sampling,
				Format:   "alsa",
			})

			ui.Sources = append(ui.Sources, source)
		}
	} else {
		var inputURL *url.URL = nil

		if len(v1data.Addresses.Input) != 0 {
			var err error
			inputURL, err = url.Parse(v1data.Addresses.Input)
			if err != nil {
				inputURL = nil
			}
		}

		if inputURL == nil {
			return r, nil
		}

		source := restreamerUISource{
			Type: "network",
		}

		settings := restreamerUISourceSettingsNetwork{
			Address: v1data.Addresses.Input,
			Mode:    "pull",
		}

		settings.General.Flags = []string{
			"genpts",
		}
		settings.General.ThreadQueueSize = 512

		settings.HTTP.ForceFramerate = false
		settings.HTTP.Framerate = 25
		settings.HTTP.ReadNative = true

		settings.RTSP.STimeout = 5000000
		settings.RTSP.UDP = !v1data.Options.RTSPTCP

		source.Settings = settings

		input := restreamerUISourceInput{
			Address: v1data.Addresses.Input,
			Options: []string{
				"-fflags",
				"+genpts",
				"-thread_queue_size",
				"512",
			},
		}

		if strings.HasPrefix(v1data.Addresses.Input, "rtsp://") {
			input.Options = append(input.Options, "-stimeout", "5000000")
			if v1data.Options.RTSPTCP {
				input.Options = append(input.Options, "-rtsp_transport", "tcp")
			}
		} else if strings.HasPrefix(v1data.Addresses.Input, "http://") || strings.HasPrefix(v1data.Addresses.Input, "https://") {
			input.Options = append(input.Options, "-analyzeduration", "20000000")
			input.Options = append(input.Options, "-re")
		} else if strings.HasPrefix(v1data.Addresses.Input, "rtmp://") {
			input.Options = append(input.Options, "-analyzeduration", "3000000")
		}

		source.Inputs = append(source.Inputs, input)

		if v1data.Addresses.Streams.Video == nil {
			// In case there's no video stream info (pre 0.6.7 database, RS_INPUTSTREAM) we have to probe the stream.
			// If this doesn't work, we have to make some assumptions.
			v1data.Addresses.Streams.Video = &storeDataV1Stream{
				Index:  0,
				Type:   "video",
				Codec:  "h264",
				Width:  1280,
				Height: 720,
				Format: "yuv420p",
			}

			config := app.Config{
				ID: "process",
				Input: []app.ConfigIO{
					{
						ID:      "in",
						Address: input.Address,
						Options: input.Options,
					},
				},
				Output: []app.ConfigIO{
					{
						ID:      "out",
						Address: "-",
						Options: []string{
							"-codec",
							"copy",
							"-f",
							"null",
						},
					},
				},
				Reconnect:      true,
				ReconnectDelay: 10,
				Autostart:      false,
				StaleTimeout:   0,
			}

			probe := probeInput(cfg.binary, config)

			for _, stream := range probe.Streams {
				if stream.Type == "video" && v1data.Addresses.Streams.Video == nil {
					v1data.Addresses.Streams.Video = &storeDataV1Stream{
						Index:  int(stream.Index),
						Type:   "video",
						Codec:  stream.Codec,
						Width:  int(stream.Width),
						Height: int(stream.Height),
						Format: stream.Format,
					}
				} else if stream.Type == "audio" && v1data.Addresses.Streams.Audio == nil {
					v1data.Addresses.Streams.Video = &storeDataV1Stream{
						Index:    int(stream.Index),
						Type:     "audio",
						Codec:    stream.Codec,
						Layout:   stream.Layout,
						Channels: int(stream.Channels),
						Sampling: strconv.FormatUint(stream.Sampling, 10),
					}
				}
			}
		}

		source.Streams = append(source.Streams, restreamerUISourceStream{
			URL:         v1data.Addresses.Input,
			Type:        "video",
			Index:       0,
			Stream:      v1data.Addresses.Streams.Video.Index,
			Codec:       v1data.Addresses.Streams.Video.Codec,
			Width:       v1data.Addresses.Streams.Video.Width,
			Height:      v1data.Addresses.Streams.Video.Height,
			PixelFormat: v1data.Addresses.Streams.Video.Format,
		})

		if v1data.Addresses.Streams.Audio != nil {
			sampling, _ := strconv.Atoi(v1data.Addresses.Streams.Audio.Sampling)

			source.Streams = append(source.Streams, restreamerUISourceStream{
				URL:      v1data.Addresses.Input,
				Type:     "audio",
				Index:    0,
				Stream:   v1data.Addresses.Streams.Audio.Index,
				Codec:    v1data.Addresses.Streams.Audio.Codec,
				Layout:   v1data.Addresses.Streams.Audio.Layout,
				Channels: v1data.Addresses.Streams.Audio.Channels,
				Sampling: sampling,
			})
		}

		ui.Sources = append(ui.Sources, source)
	}

	// Audio

	if v1data.Options.Audio.Codec == "auto" {
		if v1data.Addresses.Streams.Audio == nil {
			v1data.Options.Audio.Codec = "aac"
			v1data.Options.Audio.Preset = "silence"
			v1data.Options.Audio.Bitrate = "32"
			v1data.Options.Audio.Channels = "stereo"
			v1data.Options.Audio.Sampling = "44100"
		} else {
			if v1data.Addresses.Streams.Audio.Codec == "aac" {
				v1data.Options.Audio.Codec = "copy"
			} else {
				v1data.Options.Audio.Codec = "aac"
				v1data.Options.Audio.Preset = "encode"
			}

			if v1data.Options.Audio.Sampling != "inherit" {
				v1data.Options.Audio.Sampling = v1data.Addresses.Streams.Audio.Sampling
			}

			if v1data.Options.Audio.Channels != "inherit" {
				v1data.Options.Audio.Channels = v1data.Addresses.Streams.Audio.Layout
			}
		}
	}

	if v1data.Options.Audio.Preset == "silence" {
		if v1data.Options.Audio.Codec == "aac" || v1data.Options.Audio.Codec == "mp3" {
			source := restreamerUISource{
				Settings: restreamerUISourceSettingsVirtualAudio{
					Amplitude:  1,
					BeepFactor: 4,
					Color:      "white",
					Frequency:  440,
					Layout:     "stereo",
					Sampling:   44100,
					Source:     "silence",
				},
				Streams: []restreamerUISourceStream{},
				Type:    "virtualaudio",
			}

			source.Inputs = append(source.Inputs, restreamerUISourceInput{
				Address: "anullsrc=r=44100:cl=stereo",
				Options: []string{
					"-f",
					"lavfi",
				},
			})

			source.Streams = append(source.Streams, restreamerUISourceStream{
				Bitrate:     0,
				Channels:    2,
				Codec:       "pcm_u8",
				Coder:       "",
				Duration:    0,
				Format:      "lavfi",
				FPS:         0,
				Height:      0,
				Index:       0,
				Language:    "",
				Layout:      "stereo",
				PixelFormat: "",
				Sampling:    44100,
				Stream:      0,
				Type:        "audio",
				URL:         "anullsrc=r=44100:cl=stereo",
				Width:       0,
			})

			ui.Sources = append(ui.Sources, source)
		}
	}

	profile := restreamerUIProfile{}

	profile.Video.Source = 0
	profile.Video.Stream = v1data.Addresses.Streams.Video.Index
	profile.Video.Decoder = restreamerUIProfileAV{
		Coder:    "default",
		Mapping:  []string{},
		Settings: restreamerUIProfileCoderSettingsCopy{},
	}
	profile.Video.Encoder = restreamerUIProfileAV{
		Coder: "copy",
		Codec: "h264",
		Mapping: []string{
			"-codec:v",
			"copy",
			"-vsync",
			"0",
			"-copyts",
			"-start_at_zero",
		},
		Settings: restreamerUIProfileCoderSettingsCopy{},
	}

	if v1data.Options.Video.Codec != "copy" {
		fps, err := strconv.ParseFloat(v1data.Options.Video.FPS, 64)
		if err != nil {
			fps = 25
		}

		r := strconv.FormatFloat(fps, 'f', 0, 64)
		if math.Floor(fps) != fps {
			r = strconv.FormatFloat(fps, 'f', 2, 64)
		}

		profile.Video.Encoder.Coder = "libx264"
		profile.Video.Encoder.Codec = "h264"
		profile.Video.Encoder.Mapping = []string{
			"-codec:v",
			"libx264",
			"-preset:v", v1data.Options.Video.Preset,
			"-b:v", v1data.Options.Video.Bitrate + "k",
			"-maxrate:v", v1data.Options.Video.Bitrate + "k",
			"-bufsize:v", v1data.Options.Video.Bitrate + "k",
			"-r", r,
			"-g", strconv.FormatFloat(fps*2, 'f', 0, 64),
			"-pix_fmt", "yuv420p",
			"-vsync", "1",
		}

		if v1data.Options.Video.Profile != "auto" {
			profile.Video.Encoder.Mapping = append(profile.Video.Encoder.Mapping, "-profile:v", v1data.Options.Video.Profile)
		}

		if v1data.Options.Video.Tune != "none" {
			profile.Video.Encoder.Mapping = append(profile.Video.Encoder.Mapping, "-tune:v", v1data.Options.Video.Tune)
		}

		profile.Video.Encoder.Settings = restreamerUIProfileCoderSettingsX264{
			Bitrate: v1data.Options.Video.Bitrate,
			FPS:     v1data.Options.Video.FPS,
			Preset:  v1data.Options.Video.Preset,
			Profile: v1data.Options.Video.Profile,
			Tune:    v1data.Options.Video.Tune,
		}
	}

	profile.Audio.Source = -1
	profile.Audio.Stream = -1
	profile.Audio.Decoder = restreamerUIProfileAV{
		Coder:    "default",
		Mapping:  []string{},
		Settings: restreamerUIProfileCoderSettingsCopy{},
	}
	profile.Audio.Encoder = restreamerUIProfileAV{
		Coder:    "none",
		Codec:    "none",
		Mapping:  []string{},
		Settings: restreamerUIProfileCoderSettingsCopy{},
	}

	if v1data.Options.Audio.Codec != "none" {
		if v1data.Options.Audio.Codec == "copy" {
			profile.Audio.Encoder.Coder = "copy"
			profile.Audio.Encoder.Codec = "copy"
			profile.Audio.Encoder.Mapping = []string{
				"-codec:a", "copy",
			}
			profile.Audio.Source = 0
			profile.Audio.Stream = v1data.Addresses.Streams.Audio.Index
		} else if v1data.Options.Audio.Codec == "aac" {
			profile.Audio.Encoder.Coder = "aac"
			profile.Audio.Encoder.Codec = "aac"
			profile.Audio.Encoder.Mapping = []string{
				"-codec:a", profile.Audio.Encoder.Coder,
				"-b:a", v1data.Options.Audio.Bitrate + "k",
				"-shortest",
				"-af", "aresample=osr=" + v1data.Options.Audio.Sampling + ":ocl=" + v1data.Options.Audio.Channels,
			}

			if v1data.Options.Audio.Preset == "encode" {
				if v1data.Addresses.Streams.Audio.Codec == "aac" {
					profile.Audio.Encoder.Mapping = append(profile.Audio.Encoder.Mapping, "-bsf:a", "aac_adtstoasc")
				}
				profile.Audio.Source = 0
				profile.Audio.Stream = v1data.Addresses.Streams.Audio.Index
			} else { // silence or device
				profile.Audio.Source = 1
				profile.Audio.Stream = 0
			}

			channels := "1"
			if v1data.Options.Audio.Channels == "stereo" {
				channels = "2"
			}

			profile.Audio.Encoder.Settings = restreamerUIProfileCoderSettingsAAC{
				Bitrate:  v1data.Options.Audio.Bitrate,
				Channels: channels,
				Layout:   v1data.Options.Audio.Channels,
				Sampling: v1data.Options.Audio.Sampling,
			}
		} else if v1data.Options.Audio.Codec == "mp3" {
			profile.Audio.Encoder.Coder = "libmp3lame"
			profile.Audio.Encoder.Codec = "mp3"
			profile.Audio.Encoder.Mapping = []string{
				"-codec:a", profile.Audio.Encoder.Coder,
				"-b:a", v1data.Options.Audio.Bitrate + "k",
				"-shortest",
				"-af", "aresample=osr=" + v1data.Options.Audio.Sampling + ":ocl=" + v1data.Options.Audio.Channels,
			}

			if v1data.Options.Audio.Preset == "encode" {
				profile.Audio.Source = 0
				profile.Audio.Stream = v1data.Addresses.Streams.Audio.Index
			} else { // silence or device
				profile.Audio.Source = 1
				profile.Audio.Stream = 0
			}

			channels := "1"
			if v1data.Options.Audio.Channels == "stereo" {
				channels = "2"
			}

			profile.Audio.Encoder.Settings = restreamerUIProfileCoderSettingsAAC{
				Bitrate:  v1data.Options.Audio.Bitrate,
				Channels: channels,
				Layout:   v1data.Options.Audio.Channels,
				Sampling: v1data.Options.Audio.Sampling,
			}
		}
	}

	ui.Profiles = append(ui.Profiles, profile)

	// Input

	config.Input = append(config.Input, app.ConfigIO{
		ID:      "input_0",
		Address: ui.Sources[0].Inputs[0].Address,
		Options: ui.Sources[0].Inputs[0].Options,
	})

	if ui.Profiles[0].Audio.Source == 1 {
		config.Input = append(config.Input, app.ConfigIO{
			ID:      "input_1",
			Address: ui.Sources[1].Inputs[0].Address,
			Options: ui.Sources[1].Inputs[0].Options,
		})
	}

	//  Output

	output := app.ConfigIO{
		ID:      "output_0",
		Address: "{memfs}/" + cfg.id + ".m3u8",
		Options: []string{
			"-dn",
			"-sn",
			"-map", strconv.Itoa(ui.Profiles[0].Video.Source) + ":" + strconv.Itoa(ui.Profiles[0].Video.Stream),
		},
	}

	output.Options = append(output.Options, ui.Profiles[0].Video.Encoder.Mapping...)

	if ui.Profiles[0].Audio.Source != -1 {
		output.Options = append(output.Options, "-map", strconv.Itoa(ui.Profiles[0].Audio.Source)+":"+strconv.Itoa(ui.Profiles[0].Audio.Stream))
		output.Options = append(output.Options, ui.Profiles[0].Audio.Encoder.Mapping...)
	} else {
		output.Options = append(output.Options, "-an")
	}

	output.Options = append(output.Options,
		"-f", "hls",
		"-start_number", "0",
		"-hls_time", "2",
		"-hls_list_size", "6",
		"-hls_flags", "append_list+delete_segments",
		"-hls_segment_filename", "{memfs}/"+cfg.id+"_%04d.ts",
		"-y",
		"-method", "PUT",
	)

	config.Output = append(config.Output, output)

	process.Config = config
	r.Process[process.ID] = process

	r.Metadata.Process["restreamer-ui:ingest:"+cfg.id] = make(map[string]interface{})
	r.Metadata.Process["restreamer-ui:ingest:"+cfg.id]["restreamer-ui"] = ui

	// Snapshot

	if ui.Control.Snapshot.Enable {
		snapshotProcess := &app.Process{
			ID:        "restreamer-ui:ingest:" + cfg.id + "_snapshot",
			Reference: cfg.id,
			CreatedAt: time.Now().Unix(),
			Order:     "stop",
		}

		if v1data.Actions.Ingest == "start" {
			process.Order = "start"
		}

		snapshotConfig := &app.Config{
			ID:             "restreamer-ui:ingest:" + cfg.id + "_snapshot",
			Reference:      cfg.id,
			Reconnect:      true,
			ReconnectDelay: uint64(ui.Control.Snapshot.Interval),
			Autostart:      true,
			StaleTimeout:   30,
			Options: []string{
				"-err_detect",
				"ignore_err",
			},
		}

		snapshotInput := app.ConfigIO{
			ID:      "input_0",
			Address: "#restreamer-ui:ingest:" + cfg.id + ":output=output_0",
			Options: []string{},
		}

		snapshotConfig.Input = append(snapshotConfig.Input, snapshotInput)

		snapshotOutput := app.ConfigIO{
			ID:      "output_0",
			Address: "{memfs}/" + cfg.id + ".jpg",
			Options: []string{
				"-vframes", "1",
				"-f", "image2",
				"-update", "1",
			},
		}

		snapshotConfig.Output = append(snapshotConfig.Output, snapshotOutput)

		snapshotProcess.Config = snapshotConfig
		r.Process[snapshotProcess.ID] = snapshotProcess

		r.Metadata.Process["restreamer-ui:ingest:"+cfg.id+"_snapshot"] = nil
	}

	// Optional publication

	var outputURL *url.URL = nil

	if v1data.Options.Output.Type != "rtmp" && v1data.Options.Output.Type != "hls" {
		v1data.Addresses.Output = ""
	}

	if len(v1data.Addresses.Output) != 0 {
		var err error
		outputURL, err = url.Parse(v1data.Addresses.Output)
		if err == nil {
			if !strings.HasPrefix(outputURL.Scheme, "http") && !strings.HasPrefix(outputURL.Scheme, "https") && !strings.HasPrefix(outputURL.Scheme, "rtmp") && !strings.HasPrefix(outputURL.Scheme, "rtmps") {
				outputURL = nil
			}
		} else {
			outputURL = nil
		}
	}

	if outputURL != nil {
		egressId := "restreamer-ui:egress:rtmp:" + cfg.id
		if v1data.Options.Output.Type == "hls" {
			egressId = "restreamer-ui:egress:hls:" + cfg.id
		}

		process := &app.Process{
			ID:        egressId,
			Reference: cfg.id,
			CreatedAt: time.Now().Unix(),
			Order:     "stop",
		}

		if v1data.Actions.Egress == "start" {
			process.Order = "start"
		}

		egress := restreamerUIEgress{
			Version: 1,
		}

		if v1data.Options.Output.Type == "hls" {
			egress.Name = "HLS"

			settings := restreamerUIEgressSettingsHLS{
				Address:  strings.TrimPrefix(v1data.Addresses.Output, outputURL.Scheme+"://"),
				Protocol: outputURL.Scheme + "://",
			}

			settings.Options.DeleteThreshold = "1"
			settings.Options.Flags = []string{"delete_segments", "append_list"}
			settings.Options.InitTime = "0"
			settings.Options.ListSize = v1data.Options.Output.HLS.ListSize
			settings.Options.SegmentType = "mpegts"
			settings.Options.Time = v1data.Options.Output.HLS.Time
			settings.Options.Method = v1data.Options.Output.HLS.Method
			settings.Options.StartNumber = "0"
			settings.Options.Timeout = v1data.Options.Output.HLS.Timeout
			settings.Options.StartNumberSource = "generic"
			settings.Options.AllowCache = "0"
			settings.Options.Enc = "0"
			settings.Options.IgnoreIOErrors = "0"
			settings.Options.HTTPersistent = "0"

			egress.Settings = settings

			egress.Output.Address = v1data.Addresses.Output
			egress.Output.Options = []string{
				"-codec",
				"copy",
				"-f",
				"hls",
				"-hls_init_time",
				settings.Options.InitTime,
				"-hls_time",
				settings.Options.Time,
				"-hls_list_size",
				settings.Options.ListSize,
				"-hls_delete_threshold",
				settings.Options.DeleteThreshold,
				"-hls_start_number_source",
				settings.Options.StartNumberSource,
				"-start_number",
				settings.Options.StartNumber,
				"-hls_allow_cache",
				settings.Options.AllowCache,
				"-hls_enc",
				settings.Options.Enc,
				"-hls_segment_type",
				settings.Options.SegmentType,
				"-hls_flags",
				strings.Join(settings.Options.Flags, ","),
				"-method",
				settings.Options.Method,
				"-http_persistent",
				settings.Options.HTTPersistent,
				"-timeout",
				settings.Options.Timeout,
				"-ignore_io_errors",
				settings.Options.IgnoreIOErrors,
			}

		} else if v1data.Options.Output.Type == "rtmp" {
			egress.Name = "RTMP"

			settings := restreamerUIEgressSettingsRTMP{
				Address:  strings.TrimPrefix(v1data.Addresses.Output, outputURL.Scheme+"://"),
				Protocol: outputURL.Scheme + "://",
			}

			egress.Settings = settings

			egress.Output.Address = v1data.Addresses.Output
			egress.Output.Options = []string{
				"-codec",
				"copy",
				"-f",
				"flv",
				"-rtmp_flashver",
				"FMLE/3.0",
			}
		}

		egress.Control.Process.Autostart = true
		egress.Control.Process.Reconnect = true
		egress.Control.Process.Delay = 30
		egress.Control.Process.StaleTimeout = 30

		config := &app.Config{
			ID:             egressId,
			Reference:      cfg.id,
			Reconnect:      egress.Control.Process.Reconnect,
			ReconnectDelay: uint64(egress.Control.Process.Delay),
			Autostart:      egress.Control.Process.Autostart,
			StaleTimeout:   uint64(egress.Control.Process.StaleTimeout),
			Options: []string{
				"-err_detect",
				"ignore_err",
			},
		}

		config.Input = append(config.Input, app.ConfigIO{
			ID:      "input_0",
			Address: "#restreamer-ui:ingest:" + cfg.id + ":output=output_0",
			Options: []string{"-re"},
		})

		output := app.ConfigIO{
			ID:      "output_0",
			Address: v1data.Addresses.Output,
			Options: egress.Output.Options,
		}

		config.Output = append(config.Output, output)

		process.Config = config
		r.Process[process.ID] = process

		r.Metadata.Process[egressId] = make(map[string]interface{})
		r.Metadata.Process[egressId]["restreamer-ui"] = egress
	}

	return r, nil
}

func probeInput(binary string, config app.Config) app.Probe {
	ffmpeg, err := ffmpeg.New(ffmpeg.Config{
		Binary: binary,
	})
	if err != nil {
		return app.Probe{}
	}

	dummyfs, _ := fs.NewMemFilesystem(fs.MemConfig{})
	store, err := store.NewJSON(store.JSONConfig{
		Filesystem: dummyfs,
		Filepath:   "/",
		Logger:     nil,
	})
	if err != nil {
		return app.Probe{}
	}

	rs, err := restream.New(restream.Config{
		FFmpeg: ffmpeg,
		Store:  store,
	})
	if err != nil {
		return app.Probe{}
	}

	rs.AddProcess(&config)
	probe := rs.Probe(config.ID)
	rs.DeleteProcess(config.ID)

	return probe
}
