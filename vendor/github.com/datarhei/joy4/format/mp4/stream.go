package mp4

import (
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/format/mp4/mp4io"
	"time"
)

type Stream struct {
	av.CodecData

	trackAtom *mp4io.Track
	idx       int

	lastpkt *av.Packet

	timeScale int64
	duration  int64

	muxer   *Muxer
	demuxer *Demuxer

	sample      *mp4io.SampleTable
	sampleIndex int

	sampleOffsetInChunk int64
	syncSampleIndex     int

	dts                    int64
	sttsEntryIndex         int
	sampleIndexInSttsEntry int

	cttsEntryIndex         int
	sampleIndexInCttsEntry int

	chunkGroupIndex    int
	chunkIndex         int
	sampleIndexInChunk int

	sttsEntry *mp4io.TimeToSampleEntry
	cttsEntry *mp4io.CompositionOffsetEntry
}

func timeToTs(tm time.Duration, timeScale int64) int64 {
	return int64(tm * time.Duration(timeScale) / time.Second)
}

func tsToTime(ts int64, timeScale int64) time.Duration {
	return time.Duration(ts) * time.Second / time.Duration(timeScale)
}

func (self *Stream) timeToTs(tm time.Duration) int64 {
	return int64(tm * time.Duration(self.timeScale) / time.Second)
}

func (self *Stream) tsToTime(ts int64) time.Duration {
	return time.Duration(ts) * time.Second / time.Duration(self.timeScale)
}
