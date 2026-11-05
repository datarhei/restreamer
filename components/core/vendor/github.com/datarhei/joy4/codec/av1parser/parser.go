package av1parser

import (
	"github.com/datarhei/joy4/av"
)

type CodecData struct {
	Record    []byte
	IsMpeg2TS bool
}

func (codec CodecData) Type() av.CodecType {
	return av.AV1
}

func (codec CodecData) AV1DecoderConfRecordBytes() []byte {
	return codec.Record
}

func (codec CodecData) AV1VideoDescriptorBytes() []byte {
	return codec.Record
}

func (codec CodecData) Width() int {
	return 0
}

func (codec CodecData) Height() int {
	return 0
}

func NewCodecDataFromAV1DecoderConfRecord(record []byte) (data CodecData, err error) {
	data.Record = record
	data.IsMpeg2TS = false

	return
}

func NewCodecDataFromAV1VideoDescriptor(record []byte) (data CodecData, err error) {
	data.Record = record
	data.IsMpeg2TS = true

	return
}
