package prelude

import (
	"regexp"
	"strconv"
	"strings"
)

type IO struct {
	// common
	Address  string
	Format   string
	Index    uint64
	Stream   uint64
	Language string
	Type     string
	Codec    string
	Coder    string
	Bitrate  float64 // kbps
	Duration float64 // sec

	// video
	FPS    float64
	Pixfmt string
	Width  uint64
	Height uint64

	// audio
	Sampling uint64 // Hz
	Layout   string
	Channels uint64
}

// Parse parses the inputs and outputs from the default FFmpeg output. It returns a list of
// detected inputs and outputs as well as the number of outputs according to the stream mapping.
func Parse(lines []string) (inputs, outputs []IO, noutputs int) {
	// Input #0, lavfi, from 'testsrc=size=1280x720:rate=25':
	// Input #1, lavfi, from 'anullsrc=r=44100:cl=stereo':
	// Output #0, hls, to './data/testsrc.m3u8':
	reFormat := regexp.MustCompile(`^(Input|Output) #([0-9]+), (.*?), (from|to) '([^']+)`)

	// Duration: 00:01:02.28, start: 0.000000, bitrate: 5895 kb/s
	// Duration: N/A, start: 0.000000, bitrate: 5895 kb/s
	reDuration := regexp.MustCompile(`Duration: ([0-9]+):([0-9]+):([0-9]+)\.([0-9]+)`)

	// Stream #0:0: Video: rawvideo (RGB[24] / 0x18424752), rgb24, 1280x720 [SAR 1:1 DAR 16:9], 25 tbr, 25 tbn, 25 tbc
	// Stream #1:0: Audio: pcm_u8, 44100 Hz, stereo, u8, 705 kb/s
	// Stream #0:0: Video: h264 (libx264), yuv420p(progressive), 1280x720 [SAR 1:1 DAR 16:9], q=-1--1, 25 fps, 90k tbn, 25 tbc
	// Stream #0:1: Audio: aac (LC), 44100 Hz, stereo, fltp, 64 kb/s
	// Stream #4:0[0x100]: Video: h264 (Main) ([27][0][0][0] / 0x001B), yuv420p(tv, smpte170m/bt709/bt709, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 25 tbr, 90k tbn
	// Stream #4:1[0x101]: Audio: aac (LC) ([15][0][0][0] / 0x000F), 48000 Hz, stereo, fltp, 162 kb/s
	reStream := regexp.MustCompile(`Stream #([0-9]+):([0-9]+)(?:\[[0-9a-fx]+\])?(?:\(([a-z]+)\))?: (Video|Audio|Subtitle): (.*)`)
	reStreamCodec := regexp.MustCompile(`^([^\s,]+)`)
	reStreamVideoPixfmtSize := regexp.MustCompile(`, ([0-9A-Za-z]+)(\([^\)]+\))?, ([0-9]+)x([0-9]+)`)
	reStreamVideoFPS := regexp.MustCompile(`, ([0-9]+(\.[0-9]+)?) fps`)
	reStreamAudio := regexp.MustCompile(`, ([0-9]+) Hz, ([^,]+)`)
	reStreamBitrate := regexp.MustCompile(`, ([0-9]+) kb/s`)

	reStreamMapping := regexp.MustCompile(`^Stream mapping:`)
	reStreamMap := regexp.MustCompile(`^[\s]+Stream #[0-9]+:[0-9]+`)

	iotype := ""
	format := ""
	address := ""
	var duration float64 = 0.0

	streamMapping := false

	for _, line := range lines {
		if reStreamMapping.MatchString(line) {
			streamMapping = true
			continue
		}

		if streamMapping {
			if reStreamMap.MatchString(line) {
				noutputs++
			} else {
				streamMapping = false
			}

			continue
		}

		if matches := reFormat.FindStringSubmatch(line); matches != nil {
			iotype = matches[1]
			format = matches[3]
			address = matches[5]
			duration = 0

			continue
		}

		if matches := reDuration.FindStringSubmatch(line); matches != nil {
			duration = 0.0

			// hours
			if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
				duration += x * 60 * 60
			}

			// minutes
			if x, err := strconv.ParseFloat(matches[2], 64); err == nil {
				duration += x * 60
			}

			// seconds
			if x, err := strconv.ParseFloat(matches[3], 64); err == nil {
				duration += x
			}

			// fractions
			if x, err := strconv.ParseFloat(matches[4], 64); err == nil {
				duration += x / 100
			}

			continue
		}

		if matches := reStream.FindStringSubmatch(line); matches != nil {
			io := IO{}

			io.Address = address
			io.Format = format
			io.Duration = duration

			if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				io.Index = x
			}

			if x, err := strconv.ParseUint(matches[2], 10, 64); err == nil {
				io.Stream = x
			}

			io.Language = "und"
			if len(matches[3]) == 3 {
				io.Language = matches[3]
			}

			io.Type = strings.ToLower(matches[4])

			streamDetail := matches[5]

			if matches = reStreamCodec.FindStringSubmatch(streamDetail); matches != nil {
				io.Codec = matches[1]
			}

			if matches = reStreamBitrate.FindStringSubmatch(streamDetail); matches != nil {
				if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
					io.Bitrate = x
				}
			}

			if io.Type == "video" {
				if matches = reStreamVideoPixfmtSize.FindStringSubmatch(streamDetail); matches != nil {
					io.Pixfmt = matches[1]

					if x, err := strconv.ParseUint(matches[3], 10, 64); err == nil {
						io.Width = x
					}
					if x, err := strconv.ParseUint(matches[4], 10, 64); err == nil {
						io.Height = x
					}
				}

				if matches = reStreamVideoFPS.FindStringSubmatch(streamDetail); matches != nil {
					if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
						io.FPS = x
					}
				}

			} else if io.Type == "audio" {
				if matches = reStreamAudio.FindStringSubmatch(streamDetail); matches != nil {
					if x, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
						io.Sampling = x
					}

					io.Layout = matches[2]
				}
			}

			if iotype == "Input" {
				inputs = append(inputs, io)
			} else if iotype == "Output" {
				outputs = append(outputs, io)
			}
		}
	}

	return
}
