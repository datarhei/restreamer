package vp9parser

import (
	"github.com/datarhei/joy4/av"
)

type CodecData struct {
	Record []byte
}

func (codec CodecData) Type() av.CodecType {
	return av.VP9
}

func (codec CodecData) VPDecoderConfRecordBytes() []byte {
	return codec.Record
}

func (codec CodecData) Width() int {
	return 0
}

func (codec CodecData) Height() int {
	return 0
}

func NewCodecDataFromVPDecoderConfRecord(record []byte) (self CodecData, err error) {
	self.Record = record

	return
}
