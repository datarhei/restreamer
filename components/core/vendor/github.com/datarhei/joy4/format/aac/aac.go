package aac

import (
	"bufio"
	"fmt"
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/av/avutil"
	"github.com/datarhei/joy4/codec/aacparser"
	"io"
	"time"
)

type Muxer struct {
	w       io.Writer
	config  aacparser.MPEG4AudioConfig
	adtshdr []byte
}

func NewMuxer(w io.Writer) *Muxer {
	return &Muxer{
		adtshdr: make([]byte, aacparser.ADTSHeaderLength),
		w:       w,
	}
}

func (self *Muxer) WriteHeader(streams []av.CodecData) (err error) {
	if len(streams) > 1 || streams[0].Type() != av.AAC {
		err = fmt.Errorf("aac: must be only one aac stream")
		return
	}
	self.config = streams[0].(aacparser.CodecData).Config
	if self.config.ObjectType > aacparser.AOT_AAC_LTP {
		err = fmt.Errorf("aac: AOT %d is not allowed in ADTS", self.config.ObjectType)
	}
	return
}

func (self *Muxer) WritePacket(pkt av.Packet) (err error) {
	aacparser.FillADTSHeader(self.adtshdr, self.config, 1024, len(pkt.Data))
	if _, err = self.w.Write(self.adtshdr); err != nil {
		return
	}
	if _, err = self.w.Write(pkt.Data); err != nil {
		return
	}
	return
}

func (self *Muxer) WriteTrailer() (err error) {
	return
}

type Demuxer struct {
	r         *bufio.Reader
	config    aacparser.MPEG4AudioConfig
	codecdata av.CodecData
	ts        time.Duration
}

func NewDemuxer(r io.Reader) *Demuxer {
	return &Demuxer{
		r: bufio.NewReader(r),
	}
}

func (self *Demuxer) Streams() (streams []av.CodecData, err error) {
	if self.codecdata == nil {
		var adtshdr []byte
		var config aacparser.MPEG4AudioConfig
		if adtshdr, err = self.r.Peek(9); err != nil {
			return
		}
		if config, _, _, _, err = aacparser.ParseADTSHeader(adtshdr); err != nil {
			return
		}
		if self.codecdata, err = aacparser.NewCodecDataFromMPEG4AudioConfig(config); err != nil {
			return
		}
	}
	streams = []av.CodecData{self.codecdata}
	return
}

func (self *Demuxer) ReadPacket() (pkt av.Packet, err error) {
	var adtshdr []byte
	var config aacparser.MPEG4AudioConfig
	var hdrlen, framelen, samples int
	if adtshdr, err = self.r.Peek(9); err != nil {
		return
	}
	if config, hdrlen, framelen, samples, err = aacparser.ParseADTSHeader(adtshdr); err != nil {
		return
	}

	pkt.Data = make([]byte, framelen)
	if _, err = io.ReadFull(self.r, pkt.Data); err != nil {
		return
	}
	pkt.Data = pkt.Data[hdrlen:]

	pkt.Time = self.ts
	self.ts += time.Duration(samples) * time.Second / time.Duration(config.SampleRate)
	return
}

func Handler(h *avutil.RegisterHandler) {
	h.Ext = ".aac"

	h.ReaderDemuxer = func(r io.Reader) av.Demuxer {
		return NewDemuxer(r)
	}

	h.WriterMuxer = func(w io.Writer) av.Muxer {
		return NewMuxer(w)
	}

	h.Probe = func(b []byte) bool {
		_, _, _, _, err := aacparser.ParseADTSHeader(b)
		return err == nil
	}

	h.CodecTypes = []av.CodecType{av.AAC}
}
