package codec

import (
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/codec/fake"
	"time"
)

type PCMUCodecData struct {
	typ av.CodecType
}

func (self PCMUCodecData) Type() av.CodecType {
	return self.typ
}

func (self PCMUCodecData) SampleRate() int {
	return 8000
}

func (self PCMUCodecData) ChannelLayout() av.ChannelLayout {
	return av.CH_MONO
}

func (self PCMUCodecData) SampleFormat() av.SampleFormat {
	return av.S16
}

func (self PCMUCodecData) PacketDuration(data []byte) (time.Duration, error) {
	return time.Duration(len(data)) * time.Second / time.Duration(8000), nil
}

func NewPCMMulawCodecData() av.AudioCodecData {
	return PCMUCodecData{
		typ: av.PCM_MULAW,
	}
}

func NewPCMAlawCodecData() av.AudioCodecData {
	return PCMUCodecData{
		typ: av.PCM_ALAW,
	}
}

type SpeexCodecData struct {
	fake.CodecData
}

func (self SpeexCodecData) PacketDuration(data []byte) (time.Duration, error) {
	// libavcodec/libspeexdec.c
	// samples = samplerate/50
	// duration = 0.02s
	return time.Millisecond * 20, nil
}

func NewSpeexCodecData(sr int, cl av.ChannelLayout) SpeexCodecData {
	codec := SpeexCodecData{}
	codec.CodecType_ = av.SPEEX
	codec.SampleFormat_ = av.S16
	codec.SampleRate_ = sr
	codec.ChannelLayout_ = cl
	return codec
}
