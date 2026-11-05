package flv

import (
	"bufio"
	"fmt"
	"io"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/av/avutil"
	"github.com/datarhei/joy4/codec"
	"github.com/datarhei/joy4/codec/aacparser"
	"github.com/datarhei/joy4/codec/av1parser"
	"github.com/datarhei/joy4/codec/fake"
	"github.com/datarhei/joy4/codec/h264parser"
	"github.com/datarhei/joy4/codec/hevcparser"
	"github.com/datarhei/joy4/codec/vp9parser"
	"github.com/datarhei/joy4/format/flv/flvio"
	"github.com/datarhei/joy4/utils/bits/pio"
)

var MaxProbePacketCount = 20

func NewMetadataByStreams(streams []av.CodecData) (metadata flvio.AMFMap, err error) {
	metadata = flvio.AMFMap{}

	for _, _stream := range streams {
		typ := _stream.Type()
		switch {
		case typ.IsVideo():
			stream := _stream.(av.VideoCodecData)
			switch typ {
			case av.H264:
				metadata["videocodecid"] = flvio.VIDEO_H264
			case av.HEVC:
				metadata["videocodecid"] = flvio.FourCCToFloat(flvio.FOURCC_HEVC)
			case av.VP9:
				metadata["videocodecid"] = flvio.FourCCToFloat(flvio.FOURCC_VP9)
			case av.AV1:
				metadata["videocodecid"] = flvio.FourCCToFloat(flvio.FOURCC_AV1)

			default:
				err = fmt.Errorf("flv: metadata: unsupported video codecType=%v", stream.Type())
				return
			}

			width, height := stream.Width(), stream.Height()

			if width != 0 {
				metadata["width"] = width
				metadata["displayWidth"] = width
			}

			if height != 0 {
				metadata["height"] = height
				metadata["displayHeight"] = height
			}

		case typ.IsAudio():
			stream := _stream.(av.AudioCodecData)
			switch typ {
			case av.AAC:
				metadata["audiocodecid"] = flvio.SOUND_AAC

			case av.SPEEX:
				metadata["audiocodecid"] = flvio.SOUND_SPEEX

			default:
				err = fmt.Errorf("flv: metadata: unsupported audio codecType=%v", stream.Type())
				return
			}

			metadata["audiosamplerate"] = stream.SampleRate()
		}
	}

	return
}

type Prober struct {
	HasAudio, HasVideo             bool
	GotAudio, GotVideo             bool
	VideoStreamIdx, AudioStreamIdx int
	PushedCount                    int
	MaxProbePacketCount            int
	Streams                        []av.CodecData
	CachedPkts                     []av.Packet
}

func NewProber(maxProbePacketCount int) *Prober {
	prober := &Prober{
		MaxProbePacketCount: maxProbePacketCount,
	}

	return prober
}

func (prober *Prober) CacheTag(_tag flvio.Tag, timestamp int32) {
	pkt, _ := prober.TagToPacket(_tag, timestamp)
	prober.CachedPkts = append(prober.CachedPkts, pkt)
}

func (prober *Prober) PushTag(tag flvio.Tag, timestamp int32) (err error) {
	prober.PushedCount++

	if prober.MaxProbePacketCount <= 0 {
		prober.MaxProbePacketCount = MaxProbePacketCount
	}

	if prober.PushedCount > prober.MaxProbePacketCount {
		err = fmt.Errorf("flv: max probe packet count reached")
		return
	}

	switch tag.Type {
	case flvio.TAG_VIDEO:
		if tag.IsExHeader {
			if tag.FourCC == flvio.FOURCC_HEVC {
				if tag.PacketType == flvio.PKTTYPE_SEQUENCE_START {
					if !prober.GotVideo {
						var stream hevcparser.CodecData
						//fmt.Printf("got HEVC sequence start:\n%s\n", hex.Dump(tag.Data))
						if stream, err = hevcparser.NewCodecDataFromHEVCDecoderConfRecord(tag.Data); err != nil {
							err = fmt.Errorf("flv: hevc seqhdr invalid: %s", err.Error())
							return
						}
						prober.VideoStreamIdx = len(prober.Streams)
						prober.Streams = append(prober.Streams, stream)
						prober.GotVideo = true
					}
				} else if tag.PacketType == flvio.PKTTYPE_CODED_FRAMES || tag.PacketType == flvio.PKTTYPE_CODED_FRAMESX {
					prober.CacheTag(tag, timestamp)
				}
			} else if tag.FourCC == flvio.FOURCC_VP9 {
				if tag.PacketType == flvio.PKTTYPE_SEQUENCE_START {
					if !prober.GotVideo {
						var stream vp9parser.CodecData
						//fmt.Printf("got VP9 sequence start:\n%s\n", hex.Dump(tag.Data))
						if stream, err = vp9parser.NewCodecDataFromVPDecoderConfRecord(tag.Data); err != nil {
							err = fmt.Errorf("flv: vp9 seqhdr invalid: %s", err.Error())
							return
						}
						prober.VideoStreamIdx = len(prober.Streams)
						prober.Streams = append(prober.Streams, stream)
						prober.GotVideo = true
					}
				} else if tag.PacketType == flvio.PKTTYPE_CODED_FRAMES || tag.PacketType == flvio.PKTTYPE_CODED_FRAMESX {
					prober.CacheTag(tag, timestamp)
				}
			} else if tag.FourCC == flvio.FOURCC_AV1 {
				if tag.PacketType == flvio.PKTTYPE_SEQUENCE_START || tag.PacketType == flvio.PKTTYPE_MPEG2TS_SEQUENCE_START {
					if !prober.GotVideo {
						var stream av1parser.CodecData

						if tag.PacketType == flvio.PKTTYPE_SEQUENCE_START {
							//fmt.Printf("got AV1 sequence start:\n%s\n", hex.Dump(tag.Data))
							if stream, err = av1parser.NewCodecDataFromAV1DecoderConfRecord(tag.Data); err != nil {
								err = fmt.Errorf("flv: av1 seqhdr invalid: %s", err.Error())
								return
							}
						} else {
							//fmt.Printf("got AV1 video descriptor:\n%s\n", hex.Dump(tag.Data))
							if stream, err = av1parser.NewCodecDataFromAV1VideoDescriptor(tag.Data); err != nil {
								err = fmt.Errorf("flv: av1 video descriptor invalid: %s", err.Error())
								return
							}
						}
						prober.VideoStreamIdx = len(prober.Streams)
						prober.Streams = append(prober.Streams, stream)
						prober.GotVideo = true
					}
				} else if tag.PacketType == flvio.PKTTYPE_CODED_FRAMES || tag.PacketType == flvio.PKTTYPE_CODED_FRAMESX {
					prober.CacheTag(tag, timestamp)
				}
			}
		} else {
			switch tag.AVCPacketType {
			case flvio.AVC_SEQHDR:
				if !prober.GotVideo {
					var stream h264parser.CodecData
					//fmt.Printf("got H264 sequence start:\n%s\n", hex.Dump(tag.Data))
					if stream, err = h264parser.NewCodecDataFromAVCDecoderConfRecord(tag.Data); err != nil {
						err = fmt.Errorf("flv: h264 seqhdr invalid: %s", err.Error())
						return
					}
					prober.VideoStreamIdx = len(prober.Streams)
					prober.Streams = append(prober.Streams, stream)
					prober.GotVideo = true
				}

			case flvio.AVC_NALU:
				prober.CacheTag(tag, timestamp)
			}
		}

	case flvio.TAG_AUDIO:
		switch tag.SoundFormat {
		case flvio.SOUND_AAC:
			switch tag.AACPacketType {
			case flvio.AAC_SEQHDR:
				if !prober.GotAudio {
					var stream aacparser.CodecData
					if stream, err = aacparser.NewCodecDataFromMPEG4AudioConfigBytes(tag.Data); err != nil {
						err = fmt.Errorf("flv: aac seqhdr invalid")
						return
					}
					prober.AudioStreamIdx = len(prober.Streams)
					prober.Streams = append(prober.Streams, stream)
					prober.GotAudio = true
				}

			case flvio.AAC_RAW:
				prober.CacheTag(tag, timestamp)
			}

		case flvio.SOUND_SPEEX:
			if !prober.GotAudio {
				stream := codec.NewSpeexCodecData(16000, tag.ChannelLayout())
				prober.AudioStreamIdx = len(prober.Streams)
				prober.Streams = append(prober.Streams, stream)
				prober.GotAudio = true
				prober.CacheTag(tag, timestamp)
			}

		case flvio.SOUND_NELLYMOSER:
			if !prober.GotAudio {
				stream := fake.CodecData{
					CodecType_:     av.NELLYMOSER,
					SampleRate_:    16000,
					SampleFormat_:  av.S16,
					ChannelLayout_: tag.ChannelLayout(),
				}
				prober.AudioStreamIdx = len(prober.Streams)
				prober.Streams = append(prober.Streams, stream)
				prober.GotAudio = true
				prober.CacheTag(tag, timestamp)
			}

		}
	}

	return
}

func (prober *Prober) Probed() (ok bool) {
	if prober.MaxProbePacketCount <= 0 {
		prober.MaxProbePacketCount = MaxProbePacketCount
	}

	if prober.HasAudio || prober.HasVideo {
		if prober.HasAudio == prober.GotAudio && prober.HasVideo == prober.GotVideo {
			return true
		}
	}

	if prober.PushedCount == prober.MaxProbePacketCount {
		return true
	}

	return false
}

func (prober *Prober) TagToPacket(tag flvio.Tag, timestamp int32) (pkt av.Packet, ok bool) {
	switch tag.Type {
	case flvio.TAG_VIDEO:
		pkt.Idx = int8(prober.VideoStreamIdx)
		switch tag.PacketType {
		case flvio.PKTTYPE_CODED_FRAMES, flvio.PKTTYPE_CODED_FRAMESX:
			ok = true
			pkt.Data = tag.Data
			pkt.CompositionTime = flvio.TsToTime(tag.CompositionTime)
			pkt.IsKeyFrame = tag.FrameType == flvio.FRAME_KEY
		}

	case flvio.TAG_AUDIO:
		pkt.Idx = int8(prober.AudioStreamIdx)
		switch tag.SoundFormat {
		case flvio.SOUND_AAC:
			switch tag.AACPacketType {
			case flvio.AAC_RAW:
				ok = true
				pkt.Data = tag.Data
			}

		case flvio.SOUND_SPEEX:
			ok = true
			pkt.Data = tag.Data

		case flvio.SOUND_NELLYMOSER:
			ok = true
			pkt.Data = tag.Data
		}
	}

	pkt.Time = flvio.TsToTime(timestamp)
	return
}

func (prober *Prober) Empty() bool {
	return len(prober.CachedPkts) == 0
}

func (prober *Prober) PopPacket() av.Packet {
	pkt := prober.CachedPkts[0]
	prober.CachedPkts = prober.CachedPkts[1:]
	return pkt
}

func CodecDataToTag(stream av.CodecData) (_tag flvio.Tag, ok bool, err error) {
	switch stream.Type() {
	case av.H264:
		h264 := stream.(h264parser.CodecData)
		tag := flvio.Tag{
			Type:          flvio.TAG_VIDEO,
			AVCPacketType: flvio.AVC_SEQHDR,
			CodecID:       flvio.VIDEO_H264,
			Data:          h264.AVCDecoderConfRecordBytes(),
			FrameType:     flvio.FRAME_KEY,
		}
		//fmt.Printf("set H264 sequence start:\n%v\n", hex.Dump(tag.Data))
		ok = true
		_tag = tag

	case av.HEVC:
		hevc := stream.(hevcparser.CodecData)
		tag := flvio.Tag{
			Type:       flvio.TAG_VIDEO,
			IsExHeader: true,
			PacketType: flvio.PKTTYPE_SEQUENCE_START,
			FourCC:     flvio.FOURCC_HEVC,
			Data:       hevc.HEVCDecoderConfRecordBytes(),
			FrameType:  flvio.FRAME_KEY,
		}
		//fmt.Printf("set HEVC sequence start:\n%v\n", hex.Dump(tag.Data))
		ok = true
		_tag = tag

	case av.VP9:
		vp9 := stream.(vp9parser.CodecData)
		tag := flvio.Tag{
			Type:       flvio.TAG_VIDEO,
			IsExHeader: true,
			PacketType: flvio.PKTTYPE_SEQUENCE_START,
			FourCC:     flvio.FOURCC_VP9,
			Data:       vp9.VPDecoderConfRecordBytes(),
			FrameType:  flvio.FRAME_KEY,
		}
		//fmt.Printf("set VP9 sequence start:\n%v\n", hex.Dump(tag.Data))
		ok = true
		_tag = tag

	case av.AV1:
		av1 := stream.(av1parser.CodecData)
		tag := flvio.Tag{
			Type:       flvio.TAG_VIDEO,
			IsExHeader: true,
			PacketType: flvio.PKTTYPE_SEQUENCE_START,
			FourCC:     flvio.FOURCC_AV1,
			Data:       av1.AV1DecoderConfRecordBytes(),
			FrameType:  flvio.FRAME_KEY,
		}

		if av1.IsMpeg2TS {
			tag.PacketType = flvio.PKTTYPE_MPEG2TS_SEQUENCE_START
			tag.Data = av1.AV1VideoDescriptorBytes()
		}

		//fmt.Printf("set AV1 sequence start:\n%v\n", hex.Dump(tag.Data))
		ok = true
		_tag = tag

	case av.NELLYMOSER:
	case av.SPEEX:

	case av.AAC:
		aac := stream.(aacparser.CodecData)
		tag := flvio.Tag{
			Type:          flvio.TAG_AUDIO,
			SoundFormat:   flvio.SOUND_AAC,
			SoundRate:     flvio.SOUND_44Khz,
			AACPacketType: flvio.AAC_SEQHDR,
			Data:          aac.MPEG4AudioConfigBytes(),
		}
		switch aac.SampleFormat().BytesPerSample() {
		case 1:
			tag.SoundSize = flvio.SOUND_8BIT
		default:
			tag.SoundSize = flvio.SOUND_16BIT
		}
		switch aac.ChannelLayout().Count() {
		case 1:
			tag.SoundType = flvio.SOUND_MONO
		case 2:
			tag.SoundType = flvio.SOUND_STEREO
		}
		ok = true
		_tag = tag

	default:
		err = fmt.Errorf("flv: unspported codecType=%v", stream.Type())
		return
	}
	return
}

func PacketToTag(pkt av.Packet, stream av.CodecData) (tag flvio.Tag, timestamp int32) {
	switch stream.Type() {
	case av.H264:
		tag = flvio.Tag{
			Type:            flvio.TAG_VIDEO,
			AVCPacketType:   flvio.AVC_NALU,
			CodecID:         flvio.VIDEO_H264,
			Data:            pkt.Data,
			CompositionTime: flvio.TimeToTs(pkt.CompositionTime),
		}
		if pkt.IsKeyFrame {
			tag.FrameType = flvio.FRAME_KEY
		} else {
			tag.FrameType = flvio.FRAME_INTER
		}

	case av.HEVC:
		tag = flvio.Tag{
			Type:            flvio.TAG_VIDEO,
			IsExHeader:      true,
			PacketType:      flvio.PKTTYPE_CODED_FRAMES,
			CompositionTime: flvio.TimeToTs(pkt.CompositionTime),
			FourCC:          flvio.FOURCC_HEVC,
			Data:            pkt.Data,
		}

		if pkt.CompositionTime == 0 {
			tag.PacketType = flvio.PKTTYPE_CODED_FRAMESX
		}

		if pkt.IsKeyFrame {
			tag.FrameType = flvio.FRAME_KEY
		} else {
			tag.FrameType = flvio.FRAME_INTER
		}

	case av.VP9:
		tag = flvio.Tag{
			Type:            flvio.TAG_VIDEO,
			IsExHeader:      true,
			PacketType:      flvio.PKTTYPE_CODED_FRAMES,
			CompositionTime: flvio.TimeToTs(pkt.CompositionTime),
			FourCC:          flvio.FOURCC_VP9,
			Data:            pkt.Data,
		}

		if pkt.IsKeyFrame {
			tag.FrameType = flvio.FRAME_KEY
		} else {
			tag.FrameType = flvio.FRAME_INTER
		}

	case av.AV1:
		tag = flvio.Tag{
			Type:            flvio.TAG_VIDEO,
			IsExHeader:      true,
			PacketType:      flvio.PKTTYPE_CODED_FRAMES,
			CompositionTime: flvio.TimeToTs(pkt.CompositionTime),
			FourCC:          flvio.FOURCC_AV1,
			Data:            pkt.Data,
		}

		if pkt.IsKeyFrame {
			tag.FrameType = flvio.FRAME_KEY
		} else {
			tag.FrameType = flvio.FRAME_INTER
		}

	case av.AAC:
		tag = flvio.Tag{
			Type:          flvio.TAG_AUDIO,
			SoundFormat:   flvio.SOUND_AAC,
			SoundRate:     flvio.SOUND_44Khz,
			AACPacketType: flvio.AAC_RAW,
			Data:          pkt.Data,
		}
		astream := stream.(av.AudioCodecData)
		switch astream.SampleFormat().BytesPerSample() {
		case 1:
			tag.SoundSize = flvio.SOUND_8BIT
		default:
			tag.SoundSize = flvio.SOUND_16BIT
		}
		switch astream.ChannelLayout().Count() {
		case 1:
			tag.SoundType = flvio.SOUND_MONO
		case 2:
			tag.SoundType = flvio.SOUND_STEREO
		}

	case av.SPEEX:
		tag = flvio.Tag{
			Type:        flvio.TAG_AUDIO,
			SoundFormat: flvio.SOUND_SPEEX,
			Data:        pkt.Data,
		}

	case av.NELLYMOSER:
		tag = flvio.Tag{
			Type:        flvio.TAG_AUDIO,
			SoundFormat: flvio.SOUND_NELLYMOSER,
			Data:        pkt.Data,
		}
	}

	timestamp = flvio.TimeToTs(pkt.Time)
	return
}

type Muxer struct {
	bufw    writeFlusher
	b       []byte
	streams []av.CodecData
}

type writeFlusher interface {
	io.Writer
	Flush() error
}

func NewMuxerWriteFlusher(w writeFlusher) *Muxer {
	return &Muxer{
		bufw: w,
		b:    make([]byte, 256),
	}
}

func NewMuxer(w io.Writer) *Muxer {
	return NewMuxerWriteFlusher(bufio.NewWriterSize(w, pio.RecommendBufioSize))
}

var CodecTypes = []av.CodecType{av.H264, av.HEVC, av.VP9, av.AV1, av.AAC, av.SPEEX}

func (muxer *Muxer) WriteHeader(streams []av.CodecData) (err error) {
	var flags uint8
	for _, stream := range streams {
		if stream.Type().IsVideo() {
			flags |= flvio.FILE_HAS_VIDEO
		} else if stream.Type().IsAudio() {
			flags |= flvio.FILE_HAS_AUDIO
		}
	}

	n := flvio.FillFileHeader(muxer.b, flags)
	if _, err = muxer.bufw.Write(muxer.b[:n]); err != nil {
		return
	}

	for _, stream := range streams {
		var tag flvio.Tag
		var ok bool
		if tag, ok, err = CodecDataToTag(stream); err != nil {
			return
		}
		if ok {
			if err = flvio.WriteTag(muxer.bufw, tag, 0, muxer.b); err != nil {
				return
			}
		}
	}

	muxer.streams = streams
	return
}

func (muxer *Muxer) WritePacket(pkt av.Packet) (err error) {
	stream := muxer.streams[pkt.Idx]
	tag, timestamp := PacketToTag(pkt, stream)

	if err = flvio.WriteTag(muxer.bufw, tag, timestamp, muxer.b); err != nil {
		return
	}
	return
}

func (muxer *Muxer) WriteTrailer() (err error) {
	if err = muxer.bufw.Flush(); err != nil {
		return
	}
	return
}

type Demuxer struct {
	prober *Prober
	bufr   *bufio.Reader
	b      []byte
	stage  int
}

func NewDemuxer(r io.Reader) *Demuxer {
	return &Demuxer{
		bufr:   bufio.NewReaderSize(r, pio.RecommendBufioSize),
		prober: &Prober{},
		b:      make([]byte, 256),
	}
}

func (demuxer *Demuxer) prepare() (err error) {
	for demuxer.stage < 2 {
		switch demuxer.stage {
		case 0:
			if _, err = io.ReadFull(demuxer.bufr, demuxer.b[:flvio.FileHeaderLength]); err != nil {
				return
			}
			var flags uint8
			var skip int
			if flags, skip, err = flvio.ParseFileHeader(demuxer.b); err != nil {
				return
			}
			if _, err = demuxer.bufr.Discard(skip); err != nil {
				return
			}
			if flags&flvio.FILE_HAS_AUDIO != 0 {
				demuxer.prober.HasAudio = true
			}
			if flags&flvio.FILE_HAS_VIDEO != 0 {
				demuxer.prober.HasVideo = true
			}
			demuxer.stage++

		case 1:
			for !demuxer.prober.Probed() {
				var tag flvio.Tag
				var timestamp int32
				if tag, timestamp, err = flvio.ReadTag(demuxer.bufr, demuxer.b); err != nil {
					return
				}
				if err = demuxer.prober.PushTag(tag, timestamp); err != nil {
					return
				}
			}
			demuxer.stage++
		}
	}
	return
}

func (demuxer *Demuxer) Streams() (streams []av.CodecData, err error) {
	if err = demuxer.prepare(); err != nil {
		return
	}
	streams = demuxer.prober.Streams
	return
}

func (demuxer *Demuxer) ReadPacket() (pkt av.Packet, err error) {
	if err = demuxer.prepare(); err != nil {
		return
	}

	if !demuxer.prober.Empty() {
		pkt = demuxer.prober.PopPacket()
		return
	}

	for {
		var tag flvio.Tag
		var timestamp int32
		if tag, timestamp, err = flvio.ReadTag(demuxer.bufr, demuxer.b); err != nil {
			return
		}

		var ok bool
		if pkt, ok = demuxer.prober.TagToPacket(tag, timestamp); ok {
			return
		}
	}
}

func Handler(h *avutil.RegisterHandler) {
	h.Probe = func(b []byte) bool {
		return b[0] == 'F' && b[1] == 'L' && b[2] == 'V'
	}

	h.Ext = ".flv"

	h.ReaderDemuxer = func(r io.Reader) av.Demuxer {
		return NewDemuxer(r)
	}

	h.WriterMuxer = func(w io.Writer) av.Muxer {
		return NewMuxer(w)
	}

	h.CodecTypes = CodecTypes
}
