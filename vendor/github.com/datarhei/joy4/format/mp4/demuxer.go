package mp4

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/codec/aacparser"
	"github.com/datarhei/joy4/codec/h264parser"
	"github.com/datarhei/joy4/format/mp4/mp4io"
)

type Demuxer struct {
	r         io.ReadSeeker
	streams   []*Stream
	movieAtom *mp4io.Movie
}

func NewDemuxer(r io.ReadSeeker) *Demuxer {
	return &Demuxer{
		r: r,
	}
}

func (self *Demuxer) Streams() (streams []av.CodecData, err error) {
	if err = self.probe(); err != nil {
		return
	}
	for _, stream := range self.streams {
		streams = append(streams, stream.CodecData)
	}
	return
}

func (self *Demuxer) readat(pos int64, b []byte) (err error) {
	if _, err = self.r.Seek(pos, 0); err != nil {
		return
	}
	if _, err = io.ReadFull(self.r, b); err != nil {
		return
	}
	return
}

func (self *Demuxer) probe() (err error) {
	if self.movieAtom != nil {
		return
	}

	var moov *mp4io.Movie
	var atoms []mp4io.Atom

	if atoms, err = mp4io.ReadFileAtoms(self.r); err != nil {
		return
	}
	if _, err = self.r.Seek(0, 0); err != nil {
		return
	}

	for _, atom := range atoms {
		if atom.Tag() == mp4io.MOOV {
			moov = atom.(*mp4io.Movie)
		}
	}

	if moov == nil {
		err = fmt.Errorf("mp4: 'moov' atom not found")
		return
	}

	self.streams = []*Stream{}
	for i, atrack := range moov.Tracks {
		stream := &Stream{
			trackAtom: atrack,
			demuxer:   self,
			idx:       i,
		}
		if atrack.Media != nil && atrack.Media.Info != nil && atrack.Media.Info.Sample != nil {
			stream.sample = atrack.Media.Info.Sample
			stream.timeScale = int64(atrack.Media.Header.TimeScale)
		} else {
			err = fmt.Errorf("mp4: sample table not found")
			return
		}

		if avc1 := atrack.GetAVC1Conf(); avc1 != nil {
			if stream.CodecData, err = h264parser.NewCodecDataFromAVCDecoderConfRecord(avc1.Data); err != nil {
				return
			}
			self.streams = append(self.streams, stream)
		} else if esds := atrack.GetElemStreamDesc(); esds != nil {
			if stream.CodecData, err = aacparser.NewCodecDataFromMPEG4AudioConfigBytes(esds.DecConfig); err != nil {
				return
			}
			self.streams = append(self.streams, stream)
		}
	}

	self.movieAtom = moov
	return
}

func (self *Stream) setSampleIndex(index int) (err error) {
	found := false
	start := 0
	self.chunkGroupIndex = 0

	for self.chunkIndex = range self.sample.ChunkOffset.Entries {
		if self.chunkGroupIndex+1 < len(self.sample.SampleToChunk.Entries) &&
			uint32(self.chunkIndex+1) == self.sample.SampleToChunk.Entries[self.chunkGroupIndex+1].FirstChunk {
			self.chunkGroupIndex++
		}
		n := int(self.sample.SampleToChunk.Entries[self.chunkGroupIndex].SamplesPerChunk)
		if index >= start && index < start+n {
			found = true
			self.sampleIndexInChunk = index - start
			break
		}
		start += n
	}
	if !found {
		err = fmt.Errorf("mp4: stream[%d]: cannot locate sample index in chunk", self.idx)
		return
	}

	if self.sample.SampleSize.SampleSize != 0 {
		self.sampleOffsetInChunk = int64(self.sampleIndexInChunk) * int64(self.sample.SampleSize.SampleSize)
	} else {
		if index >= len(self.sample.SampleSize.Entries) {
			err = fmt.Errorf("mp4: stream[%d]: sample index out of range", self.idx)
			return
		}
		self.sampleOffsetInChunk = int64(0)
		for i := index - self.sampleIndexInChunk; i < index; i++ {
			self.sampleOffsetInChunk += int64(self.sample.SampleSize.Entries[i])
		}
	}

	self.dts = int64(0)
	start = 0
	found = false
	self.sttsEntryIndex = 0
	for self.sttsEntryIndex < len(self.sample.TimeToSample.Entries) {
		entry := self.sample.TimeToSample.Entries[self.sttsEntryIndex]
		n := int(entry.Count)
		if index >= start && index < start+n {
			self.sampleIndexInSttsEntry = index - start
			self.dts += int64(index-start) * int64(entry.Duration)
			found = true
			break
		}
		start += n
		self.dts += int64(n) * int64(entry.Duration)
		self.sttsEntryIndex++
	}
	if !found {
		err = fmt.Errorf("mp4: stream[%d]: cannot locate sample index in stts entry", self.idx)
		return
	}

	if self.sample.CompositionOffset != nil && len(self.sample.CompositionOffset.Entries) > 0 {
		start = 0
		found = false
		self.cttsEntryIndex = 0
		for self.cttsEntryIndex < len(self.sample.CompositionOffset.Entries) {
			n := int(self.sample.CompositionOffset.Entries[self.cttsEntryIndex].Count)
			if index >= start && index < start+n {
				self.sampleIndexInCttsEntry = index - start
				found = true
				break
			}
			start += n
			self.cttsEntryIndex++
		}
		if !found {
			err = fmt.Errorf("mp4: stream[%d]: cannot locate sample index in ctts entry", self.idx)
			return
		}
	}

	if self.sample.SyncSample != nil {
		self.syncSampleIndex = 0
		for self.syncSampleIndex < len(self.sample.SyncSample.Entries)-1 {
			if self.sample.SyncSample.Entries[self.syncSampleIndex+1]-1 > uint32(index) {
				break
			}
			self.syncSampleIndex++
		}
	}

	if false {
		fmt.Printf("mp4: stream[%d]: setSampleIndex chunkGroupIndex=%d chunkIndex=%d sampleOffsetInChunk=%d\n",
			self.idx, self.chunkGroupIndex, self.chunkIndex, self.sampleOffsetInChunk)
	}

	self.sampleIndex = index
	return
}

func (self *Stream) isSampleValid() bool {
	if self.chunkIndex >= len(self.sample.ChunkOffset.Entries) {
		return false
	}
	if self.chunkGroupIndex >= len(self.sample.SampleToChunk.Entries) {
		return false
	}
	if self.sttsEntryIndex >= len(self.sample.TimeToSample.Entries) {
		return false
	}
	if self.sample.CompositionOffset != nil && len(self.sample.CompositionOffset.Entries) > 0 {
		if self.cttsEntryIndex >= len(self.sample.CompositionOffset.Entries) {
			return false
		}
	}
	if self.sample.SyncSample != nil {
		if self.syncSampleIndex >= len(self.sample.SyncSample.Entries) {
			return false
		}
	}
	if self.sample.SampleSize.SampleSize != 0 {
		if self.sampleIndex >= len(self.sample.SampleSize.Entries) {
			return false
		}
	}
	return true
}

func (self *Stream) incSampleIndex() (duration int64) {
	if false {
		fmt.Printf("incSampleIndex sampleIndex=%d sampleOffsetInChunk=%d sampleIndexInChunk=%d chunkGroupIndex=%d chunkIndex=%d\n",
			self.sampleIndex, self.sampleOffsetInChunk, self.sampleIndexInChunk, self.chunkGroupIndex, self.chunkIndex)
	}

	self.sampleIndexInChunk++
	if uint32(self.sampleIndexInChunk) == self.sample.SampleToChunk.Entries[self.chunkGroupIndex].SamplesPerChunk {
		self.chunkIndex++
		self.sampleIndexInChunk = 0
		self.sampleOffsetInChunk = int64(0)
	} else {
		if self.sample.SampleSize.SampleSize != 0 {
			self.sampleOffsetInChunk += int64(self.sample.SampleSize.SampleSize)
		} else {
			self.sampleOffsetInChunk += int64(self.sample.SampleSize.Entries[self.sampleIndex])
		}
	}

	if self.chunkGroupIndex+1 < len(self.sample.SampleToChunk.Entries) &&
		uint32(self.chunkIndex+1) == self.sample.SampleToChunk.Entries[self.chunkGroupIndex+1].FirstChunk {
		self.chunkGroupIndex++
	}

	sttsEntry := self.sample.TimeToSample.Entries[self.sttsEntryIndex]
	duration = int64(sttsEntry.Duration)
	self.sampleIndexInSttsEntry++
	self.dts += duration
	if uint32(self.sampleIndexInSttsEntry) == sttsEntry.Count {
		self.sampleIndexInSttsEntry = 0
		self.sttsEntryIndex++
	}

	if self.sample.CompositionOffset != nil && len(self.sample.CompositionOffset.Entries) > 0 {
		self.sampleIndexInCttsEntry++
		if uint32(self.sampleIndexInCttsEntry) == self.sample.CompositionOffset.Entries[self.cttsEntryIndex].Count {
			self.sampleIndexInCttsEntry = 0
			self.cttsEntryIndex++
		}
	}

	if self.sample.SyncSample != nil {
		entries := self.sample.SyncSample.Entries
		if self.syncSampleIndex+1 < len(entries) && entries[self.syncSampleIndex+1]-1 == uint32(self.sampleIndex+1) {
			self.syncSampleIndex++
		}
	}

	self.sampleIndex++
	return
}

func (self *Stream) sampleCount() int {
	if self.sample.SampleSize.SampleSize == 0 {
		chunkGroupIndex := 0
		count := 0
		for chunkIndex := range self.sample.ChunkOffset.Entries {
			n := int(self.sample.SampleToChunk.Entries[chunkGroupIndex].SamplesPerChunk)
			count += n
			if chunkGroupIndex+1 < len(self.sample.SampleToChunk.Entries) &&
				uint32(chunkIndex+1) == self.sample.SampleToChunk.Entries[chunkGroupIndex+1].FirstChunk {
				chunkGroupIndex++
			}
		}
		return count
	} else {
		return len(self.sample.SampleSize.Entries)
	}
}

func (self *Demuxer) ReadPacket() (pkt av.Packet, err error) {
	if err = self.probe(); err != nil {
		return
	}
	if len(self.streams) == 0 {
		err = errors.New("mp4: no streams available while trying to read a packet")
		return
	}

	var chosen *Stream
	var chosenidx int
	for i, stream := range self.streams {
		if chosen == nil || stream.tsToTime(stream.dts) < chosen.tsToTime(chosen.dts) {
			chosen = stream
			chosenidx = i
		}
	}
	if false {
		fmt.Printf("ReadPacket: chosen index=%v time=%v\n", chosen.idx, chosen.tsToTime(chosen.dts))
	}
	tm := chosen.tsToTime(chosen.dts)
	if pkt, err = chosen.readPacket(); err != nil {
		return
	}
	pkt.Time = tm
	pkt.Idx = int8(chosenidx)
	return
}

func (self *Demuxer) CurrentTime() (tm time.Duration) {
	if len(self.streams) > 0 {
		stream := self.streams[0]
		tm = stream.tsToTime(stream.dts)
	}
	return
}

func (self *Demuxer) SeekToTime(tm time.Duration) (err error) {
	for _, stream := range self.streams {
		if stream.Type().IsVideo() {
			if err = stream.seekToTime(tm); err != nil {
				return
			}
			tm = stream.tsToTime(stream.dts)
			break
		}
	}

	for _, stream := range self.streams {
		if !stream.Type().IsVideo() {
			if err = stream.seekToTime(tm); err != nil {
				return
			}
		}
	}

	return
}

func (self *Stream) readPacket() (pkt av.Packet, err error) {
	if !self.isSampleValid() {
		err = io.EOF
		return
	}
	//fmt.Println("readPacket", self.sampleIndex)

	chunkOffset := self.sample.ChunkOffset.Entries[self.chunkIndex]
	sampleSize := uint32(0)
	if self.sample.SampleSize.SampleSize != 0 {
		sampleSize = self.sample.SampleSize.SampleSize
	} else {
		sampleSize = self.sample.SampleSize.Entries[self.sampleIndex]
	}

	sampleOffset := int64(chunkOffset) + self.sampleOffsetInChunk
	pkt.Data = make([]byte, sampleSize)
	if err = self.demuxer.readat(sampleOffset, pkt.Data); err != nil {
		return
	}

	if self.sample.SyncSample != nil {
		if self.sample.SyncSample.Entries[self.syncSampleIndex]-1 == uint32(self.sampleIndex) {
			pkt.IsKeyFrame = true
		}
	}

	//println("pts/dts", self.ptsEntryIndex, self.dtsEntryIndex)
	if self.sample.CompositionOffset != nil && len(self.sample.CompositionOffset.Entries) > 0 {
		cts := int64(self.sample.CompositionOffset.Entries[self.cttsEntryIndex].Offset)
		pkt.CompositionTime = self.tsToTime(cts)
	}

	self.incSampleIndex()

	return
}

func (self *Stream) seekToTime(tm time.Duration) (err error) {
	index := self.timeToSampleIndex(tm)
	if err = self.setSampleIndex(index); err != nil {
		return
	}
	if false {
		fmt.Printf("stream[%d]: seekToTime index=%v time=%v cur=%v\n", self.idx, index, tm, self.tsToTime(self.dts))
	}
	return
}

func (self *Stream) timeToSampleIndex(tm time.Duration) int {
	targetTs := self.timeToTs(tm)
	targetIndex := 0

	startTs := int64(0)
	endTs := int64(0)
	startIndex := 0
	endIndex := 0
	found := false
	for _, entry := range self.sample.TimeToSample.Entries {
		endTs = startTs + int64(entry.Count*entry.Duration)
		endIndex = startIndex + int(entry.Count)
		if targetTs >= startTs && targetTs < endTs {
			targetIndex = startIndex + int((targetTs-startTs)/int64(entry.Duration))
			found = true
		}
		startTs = endTs
		startIndex = endIndex
	}
	if !found {
		if targetTs < 0 {
			targetIndex = 0
		} else {
			targetIndex = endIndex - 1
		}
	}

	if self.sample.SyncSample != nil {
		entries := self.sample.SyncSample.Entries
		for i := len(entries) - 1; i >= 0; i-- {
			if entries[i]-1 < uint32(targetIndex) {
				targetIndex = int(entries[i] - 1)
				break
			}
		}
	}

	return targetIndex
}
