package flvio

import (
	"fmt"
	"io"
	"time"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/utils/bits/pio"
)

func TsToTime(ts int32) time.Duration {
	return time.Millisecond * time.Duration(ts)
}

func TimeToTs(tm time.Duration) int32 {
	return int32(tm / time.Millisecond)
}

const MaxTagSubHeaderLength = 16

const (
	TAG_AUDIO      = 8
	TAG_VIDEO      = 9
	TAG_SCRIPTDATA = 18
)

const (
	SOUND_MP3                   = 2
	SOUND_NELLYMOSER_16KHZ_MONO = 4
	SOUND_NELLYMOSER_8KHZ_MONO  = 5
	SOUND_NELLYMOSER            = 6
	SOUND_ALAW                  = 7
	SOUND_MULAW                 = 8
	SOUND_AAC                   = 10
	SOUND_SPEEX                 = 11

	SOUND_5_5Khz = 0
	SOUND_11Khz  = 1
	SOUND_22Khz  = 2
	SOUND_44Khz  = 3

	SOUND_8BIT  = 0
	SOUND_16BIT = 1

	SOUND_MONO   = 0
	SOUND_STEREO = 1

	AAC_SEQHDR = 0
	AAC_RAW    = 1
)

const (
	AVC_SEQHDR = 0
	AVC_NALU   = 1
	AVC_EOS    = 2

	FRAME_KEY   = 1
	FRAME_INTER = 2

	VIDEO_H264 = 7
)

const (
	PKTTYPE_SEQUENCE_START         = 0
	PKTTYPE_CODED_FRAMES           = 1
	PKTTYPE_SEQUENCE_END           = 2
	PKTTYPE_CODED_FRAMESX          = 3
	PKTTYPE_METADATA               = 4
	PKTTYPE_MPEG2TS_SEQUENCE_START = 5
)

var (
	FOURCC_AV1  = [4]byte{'a', 'v', '0', '1'}
	FOURCC_VP9  = [4]byte{'v', 'p', '0', '9'}
	FOURCC_HEVC = [4]byte{'h', 'v', 'c', '1'}
)

func FourCCToFloat(fourcc [4]byte) float64 {
	i := int(fourcc[0])<<24 | int(fourcc[1])<<16 | int(fourcc[2])<<8 | int(fourcc[3])

	return float64(i)
}

type Tag struct {
	/*
		8 = Audio
		9 = Video
		18 = Script data
	*/
	Type uint8

	/*
		SoundFormat: UB[4]
		0 = Linear PCM, platform endian
		1 = ADPCM
		2 = MP3
		3 = Linear PCM, little endian
		4 = Nellymoser 16-kHz mono
		5 = Nellymoser 8-kHz mono
		6 = Nellymoser
		7 = G.711 A-law logarithmic PCM
		8 = G.711 mu-law logarithmic PCM
		9 = reserved
		10 = AAC
		11 = Speex
		14 = MP3 8-Khz
		15 = Device-specific sound
		Formats 7, 8, 14, and 15 are reserved for internal use
		AAC is supported in Flash Player 9,0,115,0 and higher.
		Speex is supported in Flash Player 10 and higher.
	*/
	SoundFormat uint8

	/*
		SoundRate: UB[2]
		Sampling rate
		0 = 5.5-kHz For AAC: always 3
		1 = 11-kHz
		2 = 22-kHz
		3 = 44-kHz
	*/
	SoundRate uint8

	/*
		SoundSize: UB[1]
		0 = snd8Bit
		1 = snd16Bit
		Size of each sample.
		This parameter only pertains to uncompressed formats.
		Compressed formats always decode to 16 bits internally
	*/
	SoundSize uint8

	/*
		SoundType: UB[1]
		0 = sndMono
		1 = sndStereo
		Mono or stereo sound For Nellymoser: always 0
		For AAC: always 1
	*/
	SoundType uint8

	/*
		0: AAC sequence header
		1: AAC raw
	*/
	AACPacketType uint8

	/*
		0: reserved
		1: keyframe (for AVC, a seekable frame)
		2: inter frame (for AVC, a non- seekable frame)
		3: disposable inter frame (H.263 only)
		4: generated keyframe (reserved for server use only)
		5: video info/command frame
		6: reserved
		7: reserved
	*/
	FrameType uint8

	/*
		FrameType & 0b1000 != 0
	*/
	IsExHeader bool

	/*
		1: JPEG (currently unused)
		2: Sorenson H.263
		3: Screen video
		4: On2 VP6
		5: On2 VP6 with alpha channel
		6: Screen video version 2
		7: AVC
	*/
	CodecID uint8

	/*
		0: PacketTypeSequenceStart
		1: PacketTypeCodedFrames
		2: PacketTypeSequenceEnd
		3: PacketTypeCodedFramesX
		4: PacketTypeMetadata
		5: PacketTypeMPEG2TSSequenceStart
	*/
	PacketType uint8

	/*
		0: AVC sequence header
		1: AVC NALU
		2: AVC end of sequence (lower level NALU sequence ender is not required or supported)
	*/
	AVCPacketType uint8

	CompositionTime int32

	FourCC [4]byte

	Data []byte
}

func (t Tag) ChannelLayout() av.ChannelLayout {
	if t.SoundType == SOUND_MONO {
		return av.CH_MONO
	} else {
		return av.CH_STEREO
	}
}

func (t *Tag) audioParseHeader(b []byte) (n int, err error) {
	if len(b) < n+1 {
		err = fmt.Errorf("audiodata: parse invalid")
		return
	}

	flags := b[n]
	n++
	t.SoundFormat = flags >> 4
	t.SoundRate = (flags >> 2) & 0x3
	t.SoundSize = (flags >> 1) & 0x1
	t.SoundType = flags & 0x1

	switch t.SoundFormat {
	case SOUND_AAC:
		if len(b) < n+1 {
			err = fmt.Errorf("audiodata: parse invalid")
			return
		}
		t.AACPacketType = b[n]
		n++
	}

	return
}

func (t Tag) audioFillHeader(b []byte) (n int) {
	var flags uint8
	flags |= t.SoundFormat << 4
	flags |= t.SoundRate << 2
	flags |= t.SoundSize << 1
	flags |= t.SoundType
	b[n] = flags
	n++

	switch t.SoundFormat {
	case SOUND_AAC:
		b[n] = t.AACPacketType
		n++
	}

	return
}

func (t *Tag) videoParseHeader(b []byte) (n int, err error) {
	if len(b) < n+1 {
		err = fmt.Errorf("videodata: parse invalid")
		return
	}
	flags := b[n]
	t.FrameType = flags >> 4
	t.CodecID = flags & 0b1111

	//fmt.Printf("%#8b\n", flags)
	n++

	if (t.FrameType & 0b1000) != 0 {
		t.IsExHeader = true
		t.PacketType = t.CodecID
		t.CodecID = 0

		if t.PacketType != PKTTYPE_METADATA {
			t.FrameType = t.FrameType & 0b0111
		}
	}

	if !t.IsExHeader {
		if t.FrameType == FRAME_INTER || t.FrameType == FRAME_KEY {
			if len(b) < n+4 {
				err = fmt.Errorf("videodata: parse invalid: neither interframe nor keyframe")
				return
			}
			t.AVCPacketType = b[n]
			switch t.AVCPacketType {
			case AVC_SEQHDR:
				t.PacketType = PKTTYPE_SEQUENCE_START
			case AVC_NALU:
				t.PacketType = PKTTYPE_CODED_FRAMES
			case AVC_EOS:
				t.PacketType = PKTTYPE_SEQUENCE_END
			}
			n++

			t.CompositionTime = pio.I24BE(b[n:])
			n += 3
		}
	} else {
		if len(b) < n+4 {
			err = fmt.Errorf("videodata: parse invalid: not enough bytes for the fourCC value")
			return
		}

		t.FourCC[0] = b[n]
		t.FourCC[1] = b[n+1]
		t.FourCC[2] = b[n+2]
		t.FourCC[3] = b[n+3]

		n += 4

		t.CompositionTime = 0

		if t.FourCC == FOURCC_HEVC {
			if t.PacketType == PKTTYPE_CODED_FRAMES {
				t.CompositionTime = pio.I24BE(b[n:])
				n += 3
			}
		}
	}

	//fmt.Printf("parseVideoHeader: PacketType: %d\n", t.PacketType)

	return
}

func (t Tag) videoFillHeader(b []byte) (n int) {
	if t.IsExHeader {
		flags := t.FrameType<<4 | t.PacketType | 0b10000000
		b[n] = flags
		n++
		b[n] = t.FourCC[0]
		b[n+1] = t.FourCC[1]
		b[n+2] = t.FourCC[2]
		b[n+3] = t.FourCC[3]
		n += 4

		if t.FourCC == FOURCC_HEVC {
			if t.PacketType == PKTTYPE_CODED_FRAMES {
				pio.PutI24BE(b[n:], t.CompositionTime)
				n += 3
			}
		}
	} else {
		flags := t.FrameType<<4 | t.CodecID
		b[n] = flags
		n++
		b[n] = t.AVCPacketType
		n++
		pio.PutI24BE(b[n:], t.CompositionTime)
		n += 3
	}

	//fmt.Printf("videoFillHeader: PacketType: %d\n%s\n", t.PacketType, hex.Dump(b[:n]))

	return
}

func (t Tag) FillHeader(b []byte) (n int) {
	switch t.Type {
	case TAG_AUDIO:
		return t.audioFillHeader(b)

	case TAG_VIDEO:
		return t.videoFillHeader(b)
	}

	return
}

func (t *Tag) ParseHeader(b []byte) (n int, err error) {
	switch t.Type {
	case TAG_AUDIO:
		return t.audioParseHeader(b)

	case TAG_VIDEO:
		return t.videoParseHeader(b)
	}

	return
}

const (
	// TypeFlagsReserved UB[5]
	// TypeFlagsAudio    UB[1] Audio tags are present
	// TypeFlagsReserved UB[1] Must be 0
	// TypeFlagsVideo    UB[1] Video tags are present
	FILE_HAS_AUDIO = 0x4
	FILE_HAS_VIDEO = 0x1
)

const TagHeaderLength = 11
const TagTrailerLength = 4

func ParseTagHeader(b []byte) (tag Tag, ts int32, datalen int, err error) {
	tagtype := b[0]

	switch tagtype {
	case TAG_AUDIO, TAG_VIDEO, TAG_SCRIPTDATA:
		tag = Tag{Type: tagtype}

	default:
		err = fmt.Errorf("flvio: ReadTag tagtype=%d invalid", tagtype)
		return
	}

	datalen = int(pio.U24BE(b[1:4]))

	var tslo uint32
	var tshi uint8
	tslo = pio.U24BE(b[4:7])
	tshi = b[7]
	ts = int32(tslo | uint32(tshi)<<24)

	return
}

func ReadTag(r io.Reader, b []byte) (tag Tag, ts int32, err error) {
	if _, err = io.ReadFull(r, b[:TagHeaderLength]); err != nil {
		return
	}
	var datalen int
	if tag, ts, datalen, err = ParseTagHeader(b); err != nil {
		return
	}

	data := make([]byte, datalen)
	if _, err = io.ReadFull(r, data); err != nil {
		return
	}

	var n int
	if n, err = (&tag).ParseHeader(data); err != nil {
		return
	}
	tag.Data = data[n:]

	if _, err = io.ReadFull(r, b[:4]); err != nil {
		return
	}
	return
}

func FillTagHeader(b []byte, tagtype uint8, datalen int, ts int32) (n int) {
	b[n] = tagtype
	n++
	pio.PutU24BE(b[n:], uint32(datalen))
	n += 3
	pio.PutU24BE(b[n:], uint32(ts&0xffffff))
	n += 3
	b[n] = uint8(ts >> 24)
	n++
	pio.PutI24BE(b[n:], 0)
	n += 3
	return
}

func FillTagTrailer(b []byte, datalen int) (n int) {
	pio.PutU32BE(b[n:], uint32(datalen+TagHeaderLength))
	n += 4
	return
}

func WriteTag(w io.Writer, tag Tag, ts int32, b []byte) (err error) {
	data := tag.Data

	n := tag.FillHeader(b[TagHeaderLength:])
	datalen := len(data) + n

	n += FillTagHeader(b, tag.Type, datalen, ts)

	if _, err = w.Write(b[:n]); err != nil {
		return
	}

	if _, err = w.Write(data); err != nil {
		return
	}

	n = FillTagTrailer(b, datalen)
	if _, err = w.Write(b[:n]); err != nil {
		return
	}

	return
}

const FileHeaderLength = 9

func FillFileHeader(b []byte, flags uint8) (n int) {
	// 'FLV', version 1
	pio.PutU32BE(b[n:], 0x464c5601)
	n += 4

	b[n] = flags
	n++

	// DataOffset: UI32 Offset in bytes from start of file to start of body (that is, size of header)
	// The DataOffset field usually has a value of 9 for FLV version 1.
	pio.PutU32BE(b[n:], 9)
	n += 4

	// PreviousTagSize0: UI32 Always 0
	pio.PutU32BE(b[n:], 0)
	n += 4

	return
}

func ParseFileHeader(b []byte) (flags uint8, skip int, err error) {
	flv := pio.U24BE(b[0:3])
	if flv != 0x464c56 { // 'FLV'
		err = fmt.Errorf("flvio: file header cc3 invalid")
		return
	}

	flags = b[4]

	skip = int(pio.U32BE(b[5:9])) - 9 + 4
	if skip < 0 {
		err = fmt.Errorf("flvio: file header datasize invalid")
		return
	}

	return
}
