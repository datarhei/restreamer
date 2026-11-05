package fake

import (
	"github.com/datarhei/joy4/av"
)

type CodecData struct {
	CodecType_     av.CodecType
	SampleRate_    int
	SampleFormat_  av.SampleFormat
	ChannelLayout_ av.ChannelLayout
}

func (self CodecData) Type() av.CodecType {
	return self.CodecType_
}

func (self CodecData) SampleFormat() av.SampleFormat {
	return self.SampleFormat_
}

func (self CodecData) ChannelLayout() av.ChannelLayout {
	return self.ChannelLayout_
}

func (self CodecData) SampleRate() int {
	return self.SampleRate_
}
