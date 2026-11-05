package aacparser

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/utils/bits"
)

// copied from libavcodec/mpeg4audio.h
const (
	AOT_AAC_MAIN        = 1 + iota  ///< Y                       Main
	AOT_AAC_LC                      ///< Y                       Low Complexity
	AOT_AAC_SSR                     ///< N (code in SoC repo)    Scalable Sample Rate
	AOT_AAC_LTP                     ///< Y                       Long Term Prediction
	AOT_SBR                         ///< Y                       Spectral Band Replication
	AOT_AAC_SCALABLE                ///< N                       Scalable
	AOT_TWINVQ                      ///< N                       Twin Vector Quantizer
	AOT_CELP                        ///< N                       Code Excited Linear Prediction
	AOT_HVXC                        ///< N                       Harmonic Vector eXcitation Coding
	AOT_TTSI            = 12 + iota ///< N                       Text-To-Speech Interface
	AOT_MAINSYNTH                   ///< N                       Main Synthesis
	AOT_WAVESYNTH                   ///< N                       Wavetable Synthesis
	AOT_MIDI                        ///< N                       General MIDI
	AOT_SAFX                        ///< N                       Algorithmic Synthesis and Audio Effects
	AOT_ER_AAC_LC                   ///< N                       Error Resilient Low Complexity
	AOT_ER_AAC_LTP      = 19 + iota ///< N                       Error Resilient Long Term Prediction
	AOT_ER_AAC_SCALABLE             ///< N                       Error Resilient Scalable
	AOT_ER_TWINVQ                   ///< N                       Error Resilient Twin Vector Quantizer
	AOT_ER_BSAC                     ///< N                       Error Resilient Bit-Sliced Arithmetic Coding
	AOT_ER_AAC_LD                   ///< N                       Error Resilient Low Delay
	AOT_ER_CELP                     ///< N                       Error Resilient Code Excited Linear Prediction
	AOT_ER_HVXC                     ///< N                       Error Resilient Harmonic Vector eXcitation Coding
	AOT_ER_HILN                     ///< N                       Error Resilient Harmonic and Individual Lines plus Noise
	AOT_ER_PARAM                    ///< N                       Error Resilient Parametric
	AOT_SSC                         ///< N                       SinuSoidal Coding
	AOT_PS                          ///< N                       Parametric Stereo
	AOT_SURROUND                    ///< N                       MPEG Surround
	AOT_ESCAPE                      ///< Y                       Escape Value
	AOT_L1                          ///< Y                       Layer 1
	AOT_L2                          ///< Y                       Layer 2
	AOT_L3                          ///< Y                       Layer 3
	AOT_DST                         ///< N                       Direct Stream Transfer
	AOT_ALS                         ///< Y                       Audio LosslesS
	AOT_SLS                         ///< N                       Scalable LosslesS
	AOT_SLS_NON_CORE                ///< N                       Scalable LosslesS (non core)
	AOT_ER_AAC_ELD                  ///< N                       Error Resilient Enhanced Low Delay
	AOT_SMR_SIMPLE                  ///< N                       Symbolic Music Representation Simple
	AOT_SMR_MAIN                    ///< N                       Symbolic Music Representation Main
	AOT_USAC_NOSBR                  ///< N                       Unified Speech and Audio Coding (no SBR)
	AOT_SAOC                        ///< N                       Spatial Audio Object Coding
	AOT_LD_SURROUND                 ///< N                       Low Delay MPEG Surround
	AOT_USAC                        ///< N                       Unified Speech and Audio Coding
)

type MPEG4AudioConfig struct {
	SampleRate      int
	ChannelLayout   av.ChannelLayout
	ObjectType      uint
	SampleRateIndex uint
	ChannelConfig   uint
}

var sampleRateTable = []int{
	96000, 88200, 64000, 48000, 44100, 32000,
	24000, 22050, 16000, 12000, 11025, 8000, 7350,
}

/*
These are the channel configurations:
0: Defined in AOT Specifc Config
1: 1 channel: front-center
2: 2 channels: front-left, front-right
3: 3 channels: front-center, front-left, front-right
4: 4 channels: front-center, front-left, front-right, back-center
5: 5 channels: front-center, front-left, front-right, back-left, back-right
6: 6 channels: front-center, front-left, front-right, back-left, back-right, LFE-channel
7: 8 channels: front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel
8-15: Reserved
*/
var chanConfigTable = []av.ChannelLayout{
	0,
	av.CH_FRONT_CENTER,
	av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT,
	av.CH_FRONT_CENTER | av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT,
	av.CH_FRONT_CENTER | av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT | av.CH_BACK_CENTER,
	av.CH_FRONT_CENTER | av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT | av.CH_BACK_LEFT | av.CH_BACK_RIGHT,
	av.CH_FRONT_CENTER | av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT | av.CH_BACK_LEFT | av.CH_BACK_RIGHT | av.CH_LOW_FREQ,
	av.CH_FRONT_CENTER | av.CH_FRONT_LEFT | av.CH_FRONT_RIGHT | av.CH_SIDE_LEFT | av.CH_SIDE_RIGHT | av.CH_BACK_LEFT | av.CH_BACK_RIGHT | av.CH_LOW_FREQ,
}

func ParseADTSHeader(frame []byte) (config MPEG4AudioConfig, hdrlen int, framelen int, samples int, err error) {
	if frame[0] != 0xff || frame[1]&0xf6 != 0xf0 {
		err = fmt.Errorf("aacparser: not adts header")
		return
	}
	config.ObjectType = uint(frame[2]>>6) + 1
	config.SampleRateIndex = uint(frame[2] >> 2 & 0xf)
	config.ChannelConfig = uint(frame[2]<<2&0x4 | frame[3]>>6&0x3)
	if config.ChannelConfig == uint(0) {
		err = fmt.Errorf("aacparser: adts channel count invalid")
		return
	}
	(&config).Complete()
	framelen = int(frame[3]&0x3)<<11 | int(frame[4])<<3 | int(frame[5]>>5)
	samples = (int(frame[6]&0x3) + 1) * 1024
	hdrlen = 7
	if frame[1]&0x1 == 0 {
		hdrlen = 9
	}
	if framelen < hdrlen {
		err = fmt.Errorf("aacparser: adts framelen < hdrlen")
		return
	}
	return
}

const ADTSHeaderLength = 7

func FillADTSHeader(header []byte, config MPEG4AudioConfig, samples int, payloadLength int) {
	payloadLength += 7
	//AAAAAAAA AAAABCCD EEFFFFGH HHIJKLMM MMMMMMMM MMMOOOOO OOOOOOPP (QQQQQQQQ QQQQQQQQ)
	header[0] = 0xff
	header[1] = 0xf1
	header[2] = 0x50
	header[3] = 0x80
	header[4] = 0x43
	header[5] = 0xff
	header[6] = 0xcd
	//config.ObjectType = uint(frames[2]>>6)+1
	//config.SampleRateIndex = uint(frames[2]>>2&0xf)
	//config.ChannelConfig = uint(frames[2]<<2&0x4|frames[3]>>6&0x3)
	header[2] = (byte(config.ObjectType-1)&0x3)<<6 | (byte(config.SampleRateIndex)&0xf)<<2 | byte(config.ChannelConfig>>2)&0x1
	header[3] = header[3]&0x3f | byte(config.ChannelConfig&0x3)<<6
	header[3] = header[3]&0xfc | byte(payloadLength>>11)&0x3
	header[4] = byte(payloadLength >> 3)
	header[5] = header[5]&0x1f | (byte(payloadLength)&0x7)<<5
	header[6] = header[6]&0xfc | byte(samples/1024-1)
	return
}

func readObjectType(r *bits.Reader) (objectType uint, err error) {
	if objectType, err = r.ReadBits(5); err != nil {
		return
	}
	if objectType == AOT_ESCAPE {
		var i uint
		if i, err = r.ReadBits(6); err != nil {
			return
		}
		objectType = 32 + i
	}
	return
}

func writeObjectType(w *bits.Writer, objectType uint) (err error) {
	if objectType >= 32 {
		if err = w.WriteBits(AOT_ESCAPE, 5); err != nil {
			return
		}
		if err = w.WriteBits(objectType-32, 6); err != nil {
			return
		}
	} else {
		if err = w.WriteBits(objectType, 5); err != nil {
			return
		}
	}
	return
}

func readSampleRateIndex(r *bits.Reader) (index uint, err error) {
	if index, err = r.ReadBits(4); err != nil {
		return
	}
	if index == 0xf {
		if index, err = r.ReadBits(24); err != nil {
			return
		}
	}
	return
}

func writeSampleRateIndex(w *bits.Writer, index uint) (err error) {
	if index >= 0xf {
		if err = w.WriteBits(0xf, 4); err != nil {
			return
		}
		if err = w.WriteBits(index, 24); err != nil {
			return
		}
	} else {
		if err = w.WriteBits(index, 4); err != nil {
			return
		}
	}
	return
}

func (self MPEG4AudioConfig) IsValid() bool {
	return self.ObjectType > 0
}

func (self *MPEG4AudioConfig) Complete() {
	if int(self.SampleRateIndex) < len(sampleRateTable) {
		self.SampleRate = sampleRateTable[self.SampleRateIndex]
	}
	if int(self.ChannelConfig) < len(chanConfigTable) {
		self.ChannelLayout = chanConfigTable[self.ChannelConfig]
	}
	return
}

func ParseMPEG4AudioConfigBytes(data []byte) (config MPEG4AudioConfig, err error) {
	// copied from libavcodec/mpeg4audio.c avpriv_mpeg4audio_get_config()
	r := bytes.NewReader(data)
	br := &bits.Reader{R: r}
	if config.ObjectType, err = readObjectType(br); err != nil {
		return
	}
	if config.SampleRateIndex, err = readSampleRateIndex(br); err != nil {
		return
	}
	if config.ChannelConfig, err = br.ReadBits(4); err != nil {
		return
	}
	(&config).Complete()
	return
}

func WriteMPEG4AudioConfig(w io.Writer, config MPEG4AudioConfig) (err error) {
	bw := &bits.Writer{W: w}
	if err = writeObjectType(bw, config.ObjectType); err != nil {
		return
	}

	if config.SampleRateIndex == 0 {
		for i, rate := range sampleRateTable {
			if rate == config.SampleRate {
				config.SampleRateIndex = uint(i)
			}
		}
	}
	if err = writeSampleRateIndex(bw, config.SampleRateIndex); err != nil {
		return
	}

	if config.ChannelConfig == 0 {
		for i, layout := range chanConfigTable {
			if layout == config.ChannelLayout {
				config.ChannelConfig = uint(i)
			}
		}
	}
	if err = bw.WriteBits(config.ChannelConfig, 4); err != nil {
		return
	}

	if err = bw.FlushBits(); err != nil {
		return
	}
	return
}

type CodecData struct {
	ConfigBytes []byte
	Config      MPEG4AudioConfig
}

func (self CodecData) Type() av.CodecType {
	return av.AAC
}

func (self CodecData) MPEG4AudioConfigBytes() []byte {
	return self.ConfigBytes
}

func (self CodecData) ChannelLayout() av.ChannelLayout {
	return self.Config.ChannelLayout
}

func (self CodecData) SampleRate() int {
	return self.Config.SampleRate
}

func (self CodecData) SampleFormat() av.SampleFormat {
	return av.FLTP
}

func (self CodecData) PacketDuration(data []byte) (dur time.Duration, err error) {
	dur = time.Duration(1024) * time.Second / time.Duration(self.Config.SampleRate)
	return
}

func NewCodecDataFromMPEG4AudioConfig(config MPEG4AudioConfig) (self CodecData, err error) {
	b := &bytes.Buffer{}
	WriteMPEG4AudioConfig(b, config)
	return NewCodecDataFromMPEG4AudioConfigBytes(b.Bytes())
}

func NewCodecDataFromMPEG4AudioConfigBytes(config []byte) (self CodecData, err error) {
	self.ConfigBytes = config
	if self.Config, err = ParseMPEG4AudioConfigBytes(config); err != nil {
		err = fmt.Errorf("aacparser: parse MPEG4AudioConfig failed(%s)", err)
		return
	}
	return
}
