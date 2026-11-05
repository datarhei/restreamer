package mp4io

import "github.com/datarhei/joy4/utils/bits/pio"
import "time"

const MOOF = Tag(0x6d6f6f66)

func (self MovieFrag) Tag() Tag {
	return MOOF
}

const HDLR = Tag(0x68646c72)

func (self HandlerRefer) Tag() Tag {
	return HDLR
}

const AVC1 = Tag(0x61766331)

func (self AVC1Desc) Tag() Tag {
	return AVC1
}

const URL = Tag(0x75726c20)

func (self DataReferUrl) Tag() Tag {
	return URL
}

const TREX = Tag(0x74726578)

func (self TrackExtend) Tag() Tag {
	return TREX
}

const ESDS = Tag(0x65736473)

func (self ElemStreamDesc) Tag() Tag {
	return ESDS
}

const MDHD = Tag(0x6d646864)

func (self MediaHeader) Tag() Tag {
	return MDHD
}

const STTS = Tag(0x73747473)

func (self TimeToSample) Tag() Tag {
	return STTS
}

const STSS = Tag(0x73747373)

func (self SyncSample) Tag() Tag {
	return STSS
}

const MFHD = Tag(0x6d666864)

func (self MovieFragHeader) Tag() Tag {
	return MFHD
}

const MVHD = Tag(0x6d766864)

func (self MovieHeader) Tag() Tag {
	return MVHD
}

const MINF = Tag(0x6d696e66)

func (self MediaInfo) Tag() Tag {
	return MINF
}

const MOOV = Tag(0x6d6f6f76)

func (self Movie) Tag() Tag {
	return MOOV
}

const MVEX = Tag(0x6d766578)

func (self MovieExtend) Tag() Tag {
	return MVEX
}

const STSD = Tag(0x73747364)

func (self SampleDesc) Tag() Tag {
	return STSD
}

const MP4A = Tag(0x6d703461)

func (self MP4ADesc) Tag() Tag {
	return MP4A
}

const CTTS = Tag(0x63747473)

func (self CompositionOffset) Tag() Tag {
	return CTTS
}

const STCO = Tag(0x7374636f)

func (self ChunkOffset) Tag() Tag {
	return STCO
}

const TRUN = Tag(0x7472756e)

func (self TrackFragRun) Tag() Tag {
	return TRUN
}

const TRAK = Tag(0x7472616b)

func (self Track) Tag() Tag {
	return TRAK
}

const MDIA = Tag(0x6d646961)

func (self Media) Tag() Tag {
	return MDIA
}

const STSC = Tag(0x73747363)

func (self SampleToChunk) Tag() Tag {
	return STSC
}

const VMHD = Tag(0x766d6864)

func (self VideoMediaInfo) Tag() Tag {
	return VMHD
}

const STBL = Tag(0x7374626c)

func (self SampleTable) Tag() Tag {
	return STBL
}

const AVCC = Tag(0x61766343)

func (self AVC1Conf) Tag() Tag {
	return AVCC
}

const TFDT = Tag(0x74666474)

func (self TrackFragDecodeTime) Tag() Tag {
	return TFDT
}

const DINF = Tag(0x64696e66)

func (self DataInfo) Tag() Tag {
	return DINF
}

const DREF = Tag(0x64726566)

func (self DataRefer) Tag() Tag {
	return DREF
}

const TRAF = Tag(0x74726166)

func (self TrackFrag) Tag() Tag {
	return TRAF
}

const STSZ = Tag(0x7374737a)

func (self SampleSize) Tag() Tag {
	return STSZ
}

const TFHD = Tag(0x74666864)

func (self TrackFragHeader) Tag() Tag {
	return TFHD
}

const TKHD = Tag(0x746b6864)

func (self TrackHeader) Tag() Tag {
	return TKHD
}

const SMHD = Tag(0x736d6864)

func (self SoundMediaInfo) Tag() Tag {
	return SMHD
}

const MDAT = Tag(0x6d646174)

type Movie struct {
	Header      *MovieHeader
	MovieExtend *MovieExtend
	Tracks      []*Track
	Unknowns    []Atom
	AtomPos
}

func (self Movie) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MOOV))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self Movie) marshal(b []byte) (n int) {
	if self.Header != nil {
		n += self.Header.Marshal(b[n:])
	}
	if self.MovieExtend != nil {
		n += self.MovieExtend.Marshal(b[n:])
	}
	for _, atom := range self.Tracks {
		n += atom.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self Movie) Len() (n int) {
	n += 8
	if self.Header != nil {
		n += self.Header.Len()
	}
	if self.MovieExtend != nil {
		n += self.MovieExtend.Len()
	}
	for _, atom := range self.Tracks {
		n += atom.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *Movie) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case MVHD:
			{
				atom := &MovieHeader{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mvhd", n+offset, err)
					return
				}
				self.Header = atom
			}
		case MVEX:
			{
				atom := &MovieExtend{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mvex", n+offset, err)
					return
				}
				self.MovieExtend = atom
			}
		case TRAK:
			{
				atom := &Track{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("trak", n+offset, err)
					return
				}
				self.Tracks = append(self.Tracks, atom)
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self Movie) Children() (r []Atom) {
	if self.Header != nil {
		r = append(r, self.Header)
	}
	if self.MovieExtend != nil {
		r = append(r, self.MovieExtend)
	}
	for _, atom := range self.Tracks {
		r = append(r, atom)
	}
	r = append(r, self.Unknowns...)
	return
}

type MovieHeader struct {
	Version           uint8
	Flags             uint32
	CreateTime        time.Time
	ModifyTime        time.Time
	TimeScale         int32
	Duration          int32
	PreferredRate     float64
	PreferredVolume   float64
	Matrix            [9]int32
	PreviewTime       time.Time
	PreviewDuration   time.Time
	PosterTime        time.Time
	SelectionTime     time.Time
	SelectionDuration time.Time
	CurrentTime       time.Time
	NextTrackId       int32
	AtomPos
}

func (self MovieHeader) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MVHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MovieHeader) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	PutTime32(b[n:], self.CreateTime)
	n += 4
	PutTime32(b[n:], self.ModifyTime)
	n += 4
	pio.PutI32BE(b[n:], self.TimeScale)
	n += 4
	pio.PutI32BE(b[n:], self.Duration)
	n += 4
	PutFixed32(b[n:], self.PreferredRate)
	n += 4
	PutFixed16(b[n:], self.PreferredVolume)
	n += 2
	n += 10
	for _, entry := range self.Matrix {
		pio.PutI32BE(b[n:], entry)
		n += 4
	}
	PutTime32(b[n:], self.PreviewTime)
	n += 4
	PutTime32(b[n:], self.PreviewDuration)
	n += 4
	PutTime32(b[n:], self.PosterTime)
	n += 4
	PutTime32(b[n:], self.SelectionTime)
	n += 4
	PutTime32(b[n:], self.SelectionDuration)
	n += 4
	PutTime32(b[n:], self.CurrentTime)
	n += 4
	pio.PutI32BE(b[n:], self.NextTrackId)
	n += 4
	return
}
func (self MovieHeader) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	n += 2
	n += 10
	n += 4 * len(self.Matrix[:])
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	return
}
func (self *MovieHeader) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("CreateTime", n+offset, err)
		return
	}
	self.CreateTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("ModifyTime", n+offset, err)
		return
	}
	self.ModifyTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("TimeScale", n+offset, err)
		return
	}
	self.TimeScale = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("Duration", n+offset, err)
		return
	}
	self.Duration = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("PreferredRate", n+offset, err)
		return
	}
	self.PreferredRate = GetFixed32(b[n:])
	n += 4
	if len(b) < n+2 {
		err = parseErr("PreferredVolume", n+offset, err)
		return
	}
	self.PreferredVolume = GetFixed16(b[n:])
	n += 2
	n += 10
	if len(b) < n+4*len(self.Matrix) {
		err = parseErr("Matrix", n+offset, err)
		return
	}
	for i := range self.Matrix {
		self.Matrix[i] = pio.I32BE(b[n:])
		n += 4
	}
	if len(b) < n+4 {
		err = parseErr("PreviewTime", n+offset, err)
		return
	}
	self.PreviewTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("PreviewDuration", n+offset, err)
		return
	}
	self.PreviewDuration = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("PosterTime", n+offset, err)
		return
	}
	self.PosterTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("SelectionTime", n+offset, err)
		return
	}
	self.SelectionTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("SelectionDuration", n+offset, err)
		return
	}
	self.SelectionDuration = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("CurrentTime", n+offset, err)
		return
	}
	self.CurrentTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("NextTrackId", n+offset, err)
		return
	}
	self.NextTrackId = pio.I32BE(b[n:])
	n += 4
	return
}
func (self MovieHeader) Children() (r []Atom) {
	return
}

type Track struct {
	Header   *TrackHeader
	Media    *Media
	Unknowns []Atom
	AtomPos
}

func (self Track) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TRAK))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self Track) marshal(b []byte) (n int) {
	if self.Header != nil {
		n += self.Header.Marshal(b[n:])
	}
	if self.Media != nil {
		n += self.Media.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self Track) Len() (n int) {
	n += 8
	if self.Header != nil {
		n += self.Header.Len()
	}
	if self.Media != nil {
		n += self.Media.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *Track) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case TKHD:
			{
				atom := &TrackHeader{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("tkhd", n+offset, err)
					return
				}
				self.Header = atom
			}
		case MDIA:
			{
				atom := &Media{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mdia", n+offset, err)
					return
				}
				self.Media = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self Track) Children() (r []Atom) {
	if self.Header != nil {
		r = append(r, self.Header)
	}
	if self.Media != nil {
		r = append(r, self.Media)
	}
	r = append(r, self.Unknowns...)
	return
}

type TrackHeader struct {
	Version        uint8
	Flags          uint32
	CreateTime     time.Time
	ModifyTime     time.Time
	TrackId        int32
	Duration       int32
	Layer          int16
	AlternateGroup int16
	Volume         float64
	Matrix         [9]int32
	TrackWidth     float64
	TrackHeight    float64
	AtomPos
}

func (self TrackHeader) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TKHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackHeader) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	PutTime32(b[n:], self.CreateTime)
	n += 4
	PutTime32(b[n:], self.ModifyTime)
	n += 4
	pio.PutI32BE(b[n:], self.TrackId)
	n += 4
	n += 4
	pio.PutI32BE(b[n:], self.Duration)
	n += 4
	n += 8
	pio.PutI16BE(b[n:], self.Layer)
	n += 2
	pio.PutI16BE(b[n:], self.AlternateGroup)
	n += 2
	PutFixed16(b[n:], self.Volume)
	n += 2
	n += 2
	for _, entry := range self.Matrix {
		pio.PutI32BE(b[n:], entry)
		n += 4
	}
	PutFixed32(b[n:], self.TrackWidth)
	n += 4
	PutFixed32(b[n:], self.TrackHeight)
	n += 4
	return
}
func (self TrackHeader) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	n += 8
	n += 2
	n += 2
	n += 2
	n += 2
	n += 4 * len(self.Matrix[:])
	n += 4
	n += 4
	return
}
func (self *TrackHeader) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("CreateTime", n+offset, err)
		return
	}
	self.CreateTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("ModifyTime", n+offset, err)
		return
	}
	self.ModifyTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("TrackId", n+offset, err)
		return
	}
	self.TrackId = pio.I32BE(b[n:])
	n += 4
	n += 4
	if len(b) < n+4 {
		err = parseErr("Duration", n+offset, err)
		return
	}
	self.Duration = pio.I32BE(b[n:])
	n += 4
	n += 8
	if len(b) < n+2 {
		err = parseErr("Layer", n+offset, err)
		return
	}
	self.Layer = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("AlternateGroup", n+offset, err)
		return
	}
	self.AlternateGroup = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Volume", n+offset, err)
		return
	}
	self.Volume = GetFixed16(b[n:])
	n += 2
	n += 2
	if len(b) < n+4*len(self.Matrix) {
		err = parseErr("Matrix", n+offset, err)
		return
	}
	for i := range self.Matrix {
		self.Matrix[i] = pio.I32BE(b[n:])
		n += 4
	}
	if len(b) < n+4 {
		err = parseErr("TrackWidth", n+offset, err)
		return
	}
	self.TrackWidth = GetFixed32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("TrackHeight", n+offset, err)
		return
	}
	self.TrackHeight = GetFixed32(b[n:])
	n += 4
	return
}
func (self TrackHeader) Children() (r []Atom) {
	return
}

type HandlerRefer struct {
	Version uint8
	Flags   uint32
	Type    [4]byte
	SubType [4]byte
	Name    []byte
	AtomPos
}

func (self HandlerRefer) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(HDLR))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self HandlerRefer) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	copy(b[n:], self.Type[:])
	n += len(self.Type[:])
	copy(b[n:], self.SubType[:])
	n += len(self.SubType[:])
	copy(b[n:], self.Name[:])
	n += len(self.Name[:])
	return
}
func (self HandlerRefer) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += len(self.Type[:])
	n += len(self.SubType[:])
	n += len(self.Name[:])
	return
}
func (self *HandlerRefer) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+len(self.Type) {
		err = parseErr("Type", n+offset, err)
		return
	}
	copy(self.Type[:], b[n:])
	n += len(self.Type)
	if len(b) < n+len(self.SubType) {
		err = parseErr("SubType", n+offset, err)
		return
	}
	copy(self.SubType[:], b[n:])
	n += len(self.SubType)
	self.Name = b[n:]
	n += len(b[n:])
	return
}
func (self HandlerRefer) Children() (r []Atom) {
	return
}

type Media struct {
	Header   *MediaHeader
	Handler  *HandlerRefer
	Info     *MediaInfo
	Unknowns []Atom
	AtomPos
}

func (self Media) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MDIA))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self Media) marshal(b []byte) (n int) {
	if self.Header != nil {
		n += self.Header.Marshal(b[n:])
	}
	if self.Handler != nil {
		n += self.Handler.Marshal(b[n:])
	}
	if self.Info != nil {
		n += self.Info.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self Media) Len() (n int) {
	n += 8
	if self.Header != nil {
		n += self.Header.Len()
	}
	if self.Handler != nil {
		n += self.Handler.Len()
	}
	if self.Info != nil {
		n += self.Info.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *Media) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case MDHD:
			{
				atom := &MediaHeader{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mdhd", n+offset, err)
					return
				}
				self.Header = atom
			}
		case HDLR:
			{
				atom := &HandlerRefer{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("hdlr", n+offset, err)
					return
				}
				self.Handler = atom
			}
		case MINF:
			{
				atom := &MediaInfo{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("minf", n+offset, err)
					return
				}
				self.Info = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self Media) Children() (r []Atom) {
	if self.Header != nil {
		r = append(r, self.Header)
	}
	if self.Handler != nil {
		r = append(r, self.Handler)
	}
	if self.Info != nil {
		r = append(r, self.Info)
	}
	r = append(r, self.Unknowns...)
	return
}

type MediaHeader struct {
	Version    uint8
	Flags      uint32
	CreateTime time.Time
	ModifyTime time.Time
	TimeScale  int32
	Duration   int32
	Language   int16
	Quality    int16
	AtomPos
}

func (self MediaHeader) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MDHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MediaHeader) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	PutTime32(b[n:], self.CreateTime)
	n += 4
	PutTime32(b[n:], self.ModifyTime)
	n += 4
	pio.PutI32BE(b[n:], self.TimeScale)
	n += 4
	pio.PutI32BE(b[n:], self.Duration)
	n += 4
	pio.PutI16BE(b[n:], self.Language)
	n += 2
	pio.PutI16BE(b[n:], self.Quality)
	n += 2
	return
}
func (self MediaHeader) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4
	n += 4
	n += 4
	n += 2
	n += 2
	return
}
func (self *MediaHeader) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("CreateTime", n+offset, err)
		return
	}
	self.CreateTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("ModifyTime", n+offset, err)
		return
	}
	self.ModifyTime = GetTime32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("TimeScale", n+offset, err)
		return
	}
	self.TimeScale = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("Duration", n+offset, err)
		return
	}
	self.Duration = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+2 {
		err = parseErr("Language", n+offset, err)
		return
	}
	self.Language = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Quality", n+offset, err)
		return
	}
	self.Quality = pio.I16BE(b[n:])
	n += 2
	return
}
func (self MediaHeader) Children() (r []Atom) {
	return
}

type MediaInfo struct {
	Sound    *SoundMediaInfo
	Video    *VideoMediaInfo
	Data     *DataInfo
	Sample   *SampleTable
	Unknowns []Atom
	AtomPos
}

func (self MediaInfo) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MINF))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MediaInfo) marshal(b []byte) (n int) {
	if self.Sound != nil {
		n += self.Sound.Marshal(b[n:])
	}
	if self.Video != nil {
		n += self.Video.Marshal(b[n:])
	}
	if self.Data != nil {
		n += self.Data.Marshal(b[n:])
	}
	if self.Sample != nil {
		n += self.Sample.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self MediaInfo) Len() (n int) {
	n += 8
	if self.Sound != nil {
		n += self.Sound.Len()
	}
	if self.Video != nil {
		n += self.Video.Len()
	}
	if self.Data != nil {
		n += self.Data.Len()
	}
	if self.Sample != nil {
		n += self.Sample.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *MediaInfo) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case SMHD:
			{
				atom := &SoundMediaInfo{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("smhd", n+offset, err)
					return
				}
				self.Sound = atom
			}
		case VMHD:
			{
				atom := &VideoMediaInfo{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("vmhd", n+offset, err)
					return
				}
				self.Video = atom
			}
		case DINF:
			{
				atom := &DataInfo{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("dinf", n+offset, err)
					return
				}
				self.Data = atom
			}
		case STBL:
			{
				atom := &SampleTable{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stbl", n+offset, err)
					return
				}
				self.Sample = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self MediaInfo) Children() (r []Atom) {
	if self.Sound != nil {
		r = append(r, self.Sound)
	}
	if self.Video != nil {
		r = append(r, self.Video)
	}
	if self.Data != nil {
		r = append(r, self.Data)
	}
	if self.Sample != nil {
		r = append(r, self.Sample)
	}
	r = append(r, self.Unknowns...)
	return
}

type DataInfo struct {
	Refer    *DataRefer
	Unknowns []Atom
	AtomPos
}

func (self DataInfo) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(DINF))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self DataInfo) marshal(b []byte) (n int) {
	if self.Refer != nil {
		n += self.Refer.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self DataInfo) Len() (n int) {
	n += 8
	if self.Refer != nil {
		n += self.Refer.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *DataInfo) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case DREF:
			{
				atom := &DataRefer{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("dref", n+offset, err)
					return
				}
				self.Refer = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self DataInfo) Children() (r []Atom) {
	if self.Refer != nil {
		r = append(r, self.Refer)
	}
	r = append(r, self.Unknowns...)
	return
}

type DataRefer struct {
	Version uint8
	Flags   uint32
	Url     *DataReferUrl
	AtomPos
}

func (self DataRefer) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(DREF))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self DataRefer) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	_childrenNR := 0
	if self.Url != nil {
		_childrenNR++
	}
	pio.PutI32BE(b[n:], int32(_childrenNR))
	n += 4
	if self.Url != nil {
		n += self.Url.Marshal(b[n:])
	}
	return
}
func (self DataRefer) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	if self.Url != nil {
		n += self.Url.Len()
	}
	return
}
func (self *DataRefer) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	n += 4
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case URL:
			{
				atom := &DataReferUrl{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("url ", n+offset, err)
					return
				}
				self.Url = atom
			}
		}
		n += size
	}
	return
}
func (self DataRefer) Children() (r []Atom) {
	if self.Url != nil {
		r = append(r, self.Url)
	}
	return
}

type DataReferUrl struct {
	Version uint8
	Flags   uint32
	AtomPos
}

func (self DataReferUrl) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(URL))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self DataReferUrl) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	return
}
func (self DataReferUrl) Len() (n int) {
	n += 8
	n += 1
	n += 3
	return
}
func (self *DataReferUrl) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	return
}
func (self DataReferUrl) Children() (r []Atom) {
	return
}

type SoundMediaInfo struct {
	Version uint8
	Flags   uint32
	Balance int16
	AtomPos
}

func (self SoundMediaInfo) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(SMHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SoundMediaInfo) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutI16BE(b[n:], self.Balance)
	n += 2
	n += 2
	return
}
func (self SoundMediaInfo) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 2
	n += 2
	return
}
func (self *SoundMediaInfo) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+2 {
		err = parseErr("Balance", n+offset, err)
		return
	}
	self.Balance = pio.I16BE(b[n:])
	n += 2
	n += 2
	return
}
func (self SoundMediaInfo) Children() (r []Atom) {
	return
}

type VideoMediaInfo struct {
	Version      uint8
	Flags        uint32
	GraphicsMode int16
	Opcolor      [3]int16
	AtomPos
}

func (self VideoMediaInfo) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(VMHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self VideoMediaInfo) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutI16BE(b[n:], self.GraphicsMode)
	n += 2
	for _, entry := range self.Opcolor {
		pio.PutI16BE(b[n:], entry)
		n += 2
	}
	return
}
func (self VideoMediaInfo) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 2
	n += 2 * len(self.Opcolor[:])
	return
}
func (self *VideoMediaInfo) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+2 {
		err = parseErr("GraphicsMode", n+offset, err)
		return
	}
	self.GraphicsMode = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2*len(self.Opcolor) {
		err = parseErr("Opcolor", n+offset, err)
		return
	}
	for i := range self.Opcolor {
		self.Opcolor[i] = pio.I16BE(b[n:])
		n += 2
	}
	return
}
func (self VideoMediaInfo) Children() (r []Atom) {
	return
}

type SampleTable struct {
	SampleDesc        *SampleDesc
	TimeToSample      *TimeToSample
	CompositionOffset *CompositionOffset
	SampleToChunk     *SampleToChunk
	SyncSample        *SyncSample
	ChunkOffset       *ChunkOffset
	SampleSize        *SampleSize
	AtomPos
}

func (self SampleTable) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STBL))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SampleTable) marshal(b []byte) (n int) {
	if self.SampleDesc != nil {
		n += self.SampleDesc.Marshal(b[n:])
	}
	if self.TimeToSample != nil {
		n += self.TimeToSample.Marshal(b[n:])
	}
	if self.CompositionOffset != nil {
		n += self.CompositionOffset.Marshal(b[n:])
	}
	if self.SampleToChunk != nil {
		n += self.SampleToChunk.Marshal(b[n:])
	}
	if self.SyncSample != nil {
		n += self.SyncSample.Marshal(b[n:])
	}
	if self.ChunkOffset != nil {
		n += self.ChunkOffset.Marshal(b[n:])
	}
	if self.SampleSize != nil {
		n += self.SampleSize.Marshal(b[n:])
	}
	return
}
func (self SampleTable) Len() (n int) {
	n += 8
	if self.SampleDesc != nil {
		n += self.SampleDesc.Len()
	}
	if self.TimeToSample != nil {
		n += self.TimeToSample.Len()
	}
	if self.CompositionOffset != nil {
		n += self.CompositionOffset.Len()
	}
	if self.SampleToChunk != nil {
		n += self.SampleToChunk.Len()
	}
	if self.SyncSample != nil {
		n += self.SyncSample.Len()
	}
	if self.ChunkOffset != nil {
		n += self.ChunkOffset.Len()
	}
	if self.SampleSize != nil {
		n += self.SampleSize.Len()
	}
	return
}
func (self *SampleTable) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case STSD:
			{
				atom := &SampleDesc{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stsd", n+offset, err)
					return
				}
				self.SampleDesc = atom
			}
		case STTS:
			{
				atom := &TimeToSample{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stts", n+offset, err)
					return
				}
				self.TimeToSample = atom
			}
		case CTTS:
			{
				atom := &CompositionOffset{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("ctts", n+offset, err)
					return
				}
				self.CompositionOffset = atom
			}
		case STSC:
			{
				atom := &SampleToChunk{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stsc", n+offset, err)
					return
				}
				self.SampleToChunk = atom
			}
		case STSS:
			{
				atom := &SyncSample{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stss", n+offset, err)
					return
				}
				self.SyncSample = atom
			}
		case STCO:
			{
				atom := &ChunkOffset{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stco", n+offset, err)
					return
				}
				self.ChunkOffset = atom
			}
		case STSZ:
			{
				atom := &SampleSize{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("stsz", n+offset, err)
					return
				}
				self.SampleSize = atom
			}
		}
		n += size
	}
	return
}
func (self SampleTable) Children() (r []Atom) {
	if self.SampleDesc != nil {
		r = append(r, self.SampleDesc)
	}
	if self.TimeToSample != nil {
		r = append(r, self.TimeToSample)
	}
	if self.CompositionOffset != nil {
		r = append(r, self.CompositionOffset)
	}
	if self.SampleToChunk != nil {
		r = append(r, self.SampleToChunk)
	}
	if self.SyncSample != nil {
		r = append(r, self.SyncSample)
	}
	if self.ChunkOffset != nil {
		r = append(r, self.ChunkOffset)
	}
	if self.SampleSize != nil {
		r = append(r, self.SampleSize)
	}
	return
}

type SampleDesc struct {
	Version  uint8
	AVC1Desc *AVC1Desc
	MP4ADesc *MP4ADesc
	Unknowns []Atom
	AtomPos
}

func (self SampleDesc) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STSD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SampleDesc) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	n += 3
	_childrenNR := 0
	if self.AVC1Desc != nil {
		_childrenNR++
	}
	if self.MP4ADesc != nil {
		_childrenNR++
	}
	_childrenNR += len(self.Unknowns)
	pio.PutI32BE(b[n:], int32(_childrenNR))
	n += 4
	if self.AVC1Desc != nil {
		n += self.AVC1Desc.Marshal(b[n:])
	}
	if self.MP4ADesc != nil {
		n += self.MP4ADesc.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self SampleDesc) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	if self.AVC1Desc != nil {
		n += self.AVC1Desc.Len()
	}
	if self.MP4ADesc != nil {
		n += self.MP4ADesc.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *SampleDesc) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	n += 3
	n += 4
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case AVC1:
			{
				atom := &AVC1Desc{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("avc1", n+offset, err)
					return
				}
				self.AVC1Desc = atom
			}
		case MP4A:
			{
				atom := &MP4ADesc{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mp4a", n+offset, err)
					return
				}
				self.MP4ADesc = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self SampleDesc) Children() (r []Atom) {
	if self.AVC1Desc != nil {
		r = append(r, self.AVC1Desc)
	}
	if self.MP4ADesc != nil {
		r = append(r, self.MP4ADesc)
	}
	r = append(r, self.Unknowns...)
	return
}

type MP4ADesc struct {
	DataRefIdx       int16
	Version          int16
	RevisionLevel    int16
	Vendor           int32
	NumberOfChannels int16
	SampleSize       int16
	CompressionId    int16
	SampleRate       float64
	Conf             *ElemStreamDesc
	Unknowns         []Atom
	AtomPos
}

func (self MP4ADesc) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MP4A))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MP4ADesc) marshal(b []byte) (n int) {
	n += 6
	pio.PutI16BE(b[n:], self.DataRefIdx)
	n += 2
	pio.PutI16BE(b[n:], self.Version)
	n += 2
	pio.PutI16BE(b[n:], self.RevisionLevel)
	n += 2
	pio.PutI32BE(b[n:], self.Vendor)
	n += 4
	pio.PutI16BE(b[n:], self.NumberOfChannels)
	n += 2
	pio.PutI16BE(b[n:], self.SampleSize)
	n += 2
	pio.PutI16BE(b[n:], self.CompressionId)
	n += 2
	n += 2
	PutFixed32(b[n:], self.SampleRate)
	n += 4
	if self.Conf != nil {
		n += self.Conf.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self MP4ADesc) Len() (n int) {
	n += 8
	n += 6
	n += 2
	n += 2
	n += 2
	n += 4
	n += 2
	n += 2
	n += 2
	n += 2
	n += 4
	if self.Conf != nil {
		n += self.Conf.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *MP4ADesc) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	n += 6
	if len(b) < n+2 {
		err = parseErr("DataRefIdx", n+offset, err)
		return
	}
	self.DataRefIdx = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("RevisionLevel", n+offset, err)
		return
	}
	self.RevisionLevel = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+4 {
		err = parseErr("Vendor", n+offset, err)
		return
	}
	self.Vendor = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+2 {
		err = parseErr("NumberOfChannels", n+offset, err)
		return
	}
	self.NumberOfChannels = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("SampleSize", n+offset, err)
		return
	}
	self.SampleSize = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("CompressionId", n+offset, err)
		return
	}
	self.CompressionId = pio.I16BE(b[n:])
	n += 2
	n += 2
	if len(b) < n+4 {
		err = parseErr("SampleRate", n+offset, err)
		return
	}
	self.SampleRate = GetFixed32(b[n:])
	n += 4
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case ESDS:
			{
				atom := &ElemStreamDesc{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("esds", n+offset, err)
					return
				}
				self.Conf = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self MP4ADesc) Children() (r []Atom) {
	if self.Conf != nil {
		r = append(r, self.Conf)
	}
	r = append(r, self.Unknowns...)
	return
}

type AVC1Desc struct {
	DataRefIdx           int16
	Version              int16
	Revision             int16
	Vendor               int32
	TemporalQuality      int32
	SpatialQuality       int32
	Width                int16
	Height               int16
	HorizontalResolution float64
	VorizontalResolution float64
	FrameCount           int16
	CompressorName       [32]byte
	Depth                int16
	ColorTableId         int16
	Conf                 *AVC1Conf
	Unknowns             []Atom
	AtomPos
}

func (self AVC1Desc) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(AVC1))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self AVC1Desc) marshal(b []byte) (n int) {
	n += 6
	pio.PutI16BE(b[n:], self.DataRefIdx)
	n += 2
	pio.PutI16BE(b[n:], self.Version)
	n += 2
	pio.PutI16BE(b[n:], self.Revision)
	n += 2
	pio.PutI32BE(b[n:], self.Vendor)
	n += 4
	pio.PutI32BE(b[n:], self.TemporalQuality)
	n += 4
	pio.PutI32BE(b[n:], self.SpatialQuality)
	n += 4
	pio.PutI16BE(b[n:], self.Width)
	n += 2
	pio.PutI16BE(b[n:], self.Height)
	n += 2
	PutFixed32(b[n:], self.HorizontalResolution)
	n += 4
	PutFixed32(b[n:], self.VorizontalResolution)
	n += 4
	n += 4
	pio.PutI16BE(b[n:], self.FrameCount)
	n += 2
	copy(b[n:], self.CompressorName[:])
	n += len(self.CompressorName[:])
	pio.PutI16BE(b[n:], self.Depth)
	n += 2
	pio.PutI16BE(b[n:], self.ColorTableId)
	n += 2
	if self.Conf != nil {
		n += self.Conf.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self AVC1Desc) Len() (n int) {
	n += 8
	n += 6
	n += 2
	n += 2
	n += 2
	n += 4
	n += 4
	n += 4
	n += 2
	n += 2
	n += 4
	n += 4
	n += 4
	n += 2
	n += len(self.CompressorName[:])
	n += 2
	n += 2
	if self.Conf != nil {
		n += self.Conf.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *AVC1Desc) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	n += 6
	if len(b) < n+2 {
		err = parseErr("DataRefIdx", n+offset, err)
		return
	}
	self.DataRefIdx = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Revision", n+offset, err)
		return
	}
	self.Revision = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+4 {
		err = parseErr("Vendor", n+offset, err)
		return
	}
	self.Vendor = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("TemporalQuality", n+offset, err)
		return
	}
	self.TemporalQuality = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("SpatialQuality", n+offset, err)
		return
	}
	self.SpatialQuality = pio.I32BE(b[n:])
	n += 4
	if len(b) < n+2 {
		err = parseErr("Width", n+offset, err)
		return
	}
	self.Width = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("Height", n+offset, err)
		return
	}
	self.Height = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+4 {
		err = parseErr("HorizontalResolution", n+offset, err)
		return
	}
	self.HorizontalResolution = GetFixed32(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("VorizontalResolution", n+offset, err)
		return
	}
	self.VorizontalResolution = GetFixed32(b[n:])
	n += 4
	n += 4
	if len(b) < n+2 {
		err = parseErr("FrameCount", n+offset, err)
		return
	}
	self.FrameCount = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+len(self.CompressorName) {
		err = parseErr("CompressorName", n+offset, err)
		return
	}
	copy(self.CompressorName[:], b[n:])
	n += len(self.CompressorName)
	if len(b) < n+2 {
		err = parseErr("Depth", n+offset, err)
		return
	}
	self.Depth = pio.I16BE(b[n:])
	n += 2
	if len(b) < n+2 {
		err = parseErr("ColorTableId", n+offset, err)
		return
	}
	self.ColorTableId = pio.I16BE(b[n:])
	n += 2
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case AVCC:
			{
				atom := &AVC1Conf{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("avcC", n+offset, err)
					return
				}
				self.Conf = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self AVC1Desc) Children() (r []Atom) {
	if self.Conf != nil {
		r = append(r, self.Conf)
	}
	r = append(r, self.Unknowns...)
	return
}

type AVC1Conf struct {
	Data []byte
	AtomPos
}

func (self AVC1Conf) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(AVCC))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self AVC1Conf) marshal(b []byte) (n int) {
	copy(b[n:], self.Data[:])
	n += len(self.Data[:])
	return
}
func (self AVC1Conf) Len() (n int) {
	n += 8
	n += len(self.Data[:])
	return
}
func (self *AVC1Conf) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	self.Data = b[n:]
	n += len(b[n:])
	return
}
func (self AVC1Conf) Children() (r []Atom) {
	return
}

type TimeToSample struct {
	Version uint8
	Flags   uint32
	Entries []TimeToSampleEntry
	AtomPos
}

func (self TimeToSample) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STTS))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TimeToSample) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		PutTimeToSampleEntry(b[n:], entry)
		n += LenTimeToSampleEntry
	}
	return
}
func (self TimeToSample) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += LenTimeToSampleEntry * len(self.Entries)
	return
}
func (self *TimeToSample) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]TimeToSampleEntry, _len_Entries)
	if len(b) < n+LenTimeToSampleEntry*len(self.Entries) {
		err = parseErr("TimeToSampleEntry", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = GetTimeToSampleEntry(b[n:])
		n += LenTimeToSampleEntry
	}
	return
}
func (self TimeToSample) Children() (r []Atom) {
	return
}

type TimeToSampleEntry struct {
	Count    uint32
	Duration uint32
}

func GetTimeToSampleEntry(b []byte) (self TimeToSampleEntry) {
	self.Count = pio.U32BE(b[0:])
	self.Duration = pio.U32BE(b[4:])
	return
}
func PutTimeToSampleEntry(b []byte, self TimeToSampleEntry) {
	pio.PutU32BE(b[0:], self.Count)
	pio.PutU32BE(b[4:], self.Duration)
}

const LenTimeToSampleEntry = 8

type SampleToChunk struct {
	Version uint8
	Flags   uint32
	Entries []SampleToChunkEntry
	AtomPos
}

func (self SampleToChunk) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STSC))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SampleToChunk) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		PutSampleToChunkEntry(b[n:], entry)
		n += LenSampleToChunkEntry
	}
	return
}
func (self SampleToChunk) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += LenSampleToChunkEntry * len(self.Entries)
	return
}
func (self *SampleToChunk) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]SampleToChunkEntry, _len_Entries)
	if len(b) < n+LenSampleToChunkEntry*len(self.Entries) {
		err = parseErr("SampleToChunkEntry", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = GetSampleToChunkEntry(b[n:])
		n += LenSampleToChunkEntry
	}
	return
}
func (self SampleToChunk) Children() (r []Atom) {
	return
}

type SampleToChunkEntry struct {
	FirstChunk      uint32
	SamplesPerChunk uint32
	SampleDescId    uint32
}

func GetSampleToChunkEntry(b []byte) (self SampleToChunkEntry) {
	self.FirstChunk = pio.U32BE(b[0:])
	self.SamplesPerChunk = pio.U32BE(b[4:])
	self.SampleDescId = pio.U32BE(b[8:])
	return
}
func PutSampleToChunkEntry(b []byte, self SampleToChunkEntry) {
	pio.PutU32BE(b[0:], self.FirstChunk)
	pio.PutU32BE(b[4:], self.SamplesPerChunk)
	pio.PutU32BE(b[8:], self.SampleDescId)
}

const LenSampleToChunkEntry = 12

type CompositionOffset struct {
	Version uint8
	Flags   uint32
	Entries []CompositionOffsetEntry
	AtomPos
}

func (self CompositionOffset) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(CTTS))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self CompositionOffset) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		PutCompositionOffsetEntry(b[n:], entry)
		n += LenCompositionOffsetEntry
	}
	return
}
func (self CompositionOffset) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += LenCompositionOffsetEntry * len(self.Entries)
	return
}
func (self *CompositionOffset) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]CompositionOffsetEntry, _len_Entries)
	if len(b) < n+LenCompositionOffsetEntry*len(self.Entries) {
		err = parseErr("CompositionOffsetEntry", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = GetCompositionOffsetEntry(b[n:])
		n += LenCompositionOffsetEntry
	}
	return
}
func (self CompositionOffset) Children() (r []Atom) {
	return
}

type CompositionOffsetEntry struct {
	Count  uint32
	Offset uint32
}

func GetCompositionOffsetEntry(b []byte) (self CompositionOffsetEntry) {
	self.Count = pio.U32BE(b[0:])
	self.Offset = pio.U32BE(b[4:])
	return
}
func PutCompositionOffsetEntry(b []byte, self CompositionOffsetEntry) {
	pio.PutU32BE(b[0:], self.Count)
	pio.PutU32BE(b[4:], self.Offset)
}

const LenCompositionOffsetEntry = 8

type SyncSample struct {
	Version uint8
	Flags   uint32
	Entries []uint32
	AtomPos
}

func (self SyncSample) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STSS))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SyncSample) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		pio.PutU32BE(b[n:], entry)
		n += 4
	}
	return
}
func (self SyncSample) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4 * len(self.Entries)
	return
}
func (self *SyncSample) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]uint32, _len_Entries)
	if len(b) < n+4*len(self.Entries) {
		err = parseErr("uint32", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = pio.U32BE(b[n:])
		n += 4
	}
	return
}
func (self SyncSample) Children() (r []Atom) {
	return
}

type ChunkOffset struct {
	Version uint8
	Flags   uint32
	Entries []uint32
	AtomPos
}

func (self ChunkOffset) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STCO))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self ChunkOffset) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		pio.PutU32BE(b[n:], entry)
		n += 4
	}
	return
}
func (self ChunkOffset) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4 * len(self.Entries)
	return
}
func (self *ChunkOffset) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]uint32, _len_Entries)
	if len(b) < n+4*len(self.Entries) {
		err = parseErr("uint32", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = pio.U32BE(b[n:])
		n += 4
	}
	return
}
func (self ChunkOffset) Children() (r []Atom) {
	return
}

type MovieFrag struct {
	Header   *MovieFragHeader
	Tracks   []*TrackFrag
	Unknowns []Atom
	AtomPos
}

func (self MovieFrag) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MOOF))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MovieFrag) marshal(b []byte) (n int) {
	if self.Header != nil {
		n += self.Header.Marshal(b[n:])
	}
	for _, atom := range self.Tracks {
		n += atom.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self MovieFrag) Len() (n int) {
	n += 8
	if self.Header != nil {
		n += self.Header.Len()
	}
	for _, atom := range self.Tracks {
		n += atom.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *MovieFrag) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case MFHD:
			{
				atom := &MovieFragHeader{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("mfhd", n+offset, err)
					return
				}
				self.Header = atom
			}
		case TRAF:
			{
				atom := &TrackFrag{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("traf", n+offset, err)
					return
				}
				self.Tracks = append(self.Tracks, atom)
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self MovieFrag) Children() (r []Atom) {
	if self.Header != nil {
		r = append(r, self.Header)
	}
	for _, atom := range self.Tracks {
		r = append(r, atom)
	}
	r = append(r, self.Unknowns...)
	return
}

type MovieFragHeader struct {
	Version uint8
	Flags   uint32
	Seqnum  uint32
	AtomPos
}

func (self MovieFragHeader) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MFHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MovieFragHeader) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], self.Seqnum)
	n += 4
	return
}
func (self MovieFragHeader) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	return
}
func (self *MovieFragHeader) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("Seqnum", n+offset, err)
		return
	}
	self.Seqnum = pio.U32BE(b[n:])
	n += 4
	return
}
func (self MovieFragHeader) Children() (r []Atom) {
	return
}

type TrackFrag struct {
	Header     *TrackFragHeader
	DecodeTime *TrackFragDecodeTime
	Run        *TrackFragRun
	Unknowns   []Atom
	AtomPos
}

func (self TrackFrag) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TRAF))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackFrag) marshal(b []byte) (n int) {
	if self.Header != nil {
		n += self.Header.Marshal(b[n:])
	}
	if self.DecodeTime != nil {
		n += self.DecodeTime.Marshal(b[n:])
	}
	if self.Run != nil {
		n += self.Run.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self TrackFrag) Len() (n int) {
	n += 8
	if self.Header != nil {
		n += self.Header.Len()
	}
	if self.DecodeTime != nil {
		n += self.DecodeTime.Len()
	}
	if self.Run != nil {
		n += self.Run.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *TrackFrag) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case TFHD:
			{
				atom := &TrackFragHeader{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("tfhd", n+offset, err)
					return
				}
				self.Header = atom
			}
		case TFDT:
			{
				atom := &TrackFragDecodeTime{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("tfdt", n+offset, err)
					return
				}
				self.DecodeTime = atom
			}
		case TRUN:
			{
				atom := &TrackFragRun{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("trun", n+offset, err)
					return
				}
				self.Run = atom
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self TrackFrag) Children() (r []Atom) {
	if self.Header != nil {
		r = append(r, self.Header)
	}
	if self.DecodeTime != nil {
		r = append(r, self.DecodeTime)
	}
	if self.Run != nil {
		r = append(r, self.Run)
	}
	r = append(r, self.Unknowns...)
	return
}

type MovieExtend struct {
	Tracks   []*TrackExtend
	Unknowns []Atom
	AtomPos
}

func (self MovieExtend) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(MVEX))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self MovieExtend) marshal(b []byte) (n int) {
	for _, atom := range self.Tracks {
		n += atom.Marshal(b[n:])
	}
	for _, atom := range self.Unknowns {
		n += atom.Marshal(b[n:])
	}
	return
}
func (self MovieExtend) Len() (n int) {
	n += 8
	for _, atom := range self.Tracks {
		n += atom.Len()
	}
	for _, atom := range self.Unknowns {
		n += atom.Len()
	}
	return
}
func (self *MovieExtend) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	for n+8 < len(b) {
		tag := Tag(pio.U32BE(b[n+4:]))
		size := int(pio.U32BE(b[n:]))
		if len(b) < n+size {
			err = parseErr("TagSizeInvalid", n+offset, err)
			return
		}
		switch tag {
		case TREX:
			{
				atom := &TrackExtend{}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("trex", n+offset, err)
					return
				}
				self.Tracks = append(self.Tracks, atom)
			}
		default:
			{
				atom := &Dummy{Tag_: tag, Data: b[n : n+size]}
				if _, err = atom.Unmarshal(b[n:n+size], offset+n); err != nil {
					err = parseErr("", n+offset, err)
					return
				}
				self.Unknowns = append(self.Unknowns, atom)
			}
		}
		n += size
	}
	return
}
func (self MovieExtend) Children() (r []Atom) {
	for _, atom := range self.Tracks {
		r = append(r, atom)
	}
	r = append(r, self.Unknowns...)
	return
}

type TrackExtend struct {
	Version               uint8
	Flags                 uint32
	TrackId               uint32
	DefaultSampleDescIdx  uint32
	DefaultSampleDuration uint32
	DefaultSampleSize     uint32
	DefaultSampleFlags    uint32
	AtomPos
}

func (self TrackExtend) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TREX))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackExtend) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], self.TrackId)
	n += 4
	pio.PutU32BE(b[n:], self.DefaultSampleDescIdx)
	n += 4
	pio.PutU32BE(b[n:], self.DefaultSampleDuration)
	n += 4
	pio.PutU32BE(b[n:], self.DefaultSampleSize)
	n += 4
	pio.PutU32BE(b[n:], self.DefaultSampleFlags)
	n += 4
	return
}
func (self TrackExtend) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	n += 4
	n += 4
	n += 4
	n += 4
	return
}
func (self *TrackExtend) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("TrackId", n+offset, err)
		return
	}
	self.TrackId = pio.U32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("DefaultSampleDescIdx", n+offset, err)
		return
	}
	self.DefaultSampleDescIdx = pio.U32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("DefaultSampleDuration", n+offset, err)
		return
	}
	self.DefaultSampleDuration = pio.U32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("DefaultSampleSize", n+offset, err)
		return
	}
	self.DefaultSampleSize = pio.U32BE(b[n:])
	n += 4
	if len(b) < n+4 {
		err = parseErr("DefaultSampleFlags", n+offset, err)
		return
	}
	self.DefaultSampleFlags = pio.U32BE(b[n:])
	n += 4
	return
}
func (self TrackExtend) Children() (r []Atom) {
	return
}

type SampleSize struct {
	Version    uint8
	Flags      uint32
	SampleSize uint32
	Entries    []uint32
	AtomPos
}

func (self SampleSize) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(STSZ))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self SampleSize) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], self.SampleSize)
	n += 4
	if self.SampleSize != 0 {
		return
	}
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	for _, entry := range self.Entries {
		pio.PutU32BE(b[n:], entry)
		n += 4
	}
	return
}
func (self SampleSize) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	if self.SampleSize != 0 {
		return
	}
	n += 4
	n += 4 * len(self.Entries)
	return
}
func (self *SampleSize) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if len(b) < n+4 {
		err = parseErr("SampleSize", n+offset, err)
		return
	}
	self.SampleSize = pio.U32BE(b[n:])
	n += 4
	if self.SampleSize != 0 {
		return
	}
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]uint32, _len_Entries)
	if len(b) < n+4*len(self.Entries) {
		err = parseErr("uint32", n+offset, err)
		return
	}
	for i := range self.Entries {
		self.Entries[i] = pio.U32BE(b[n:])
		n += 4
	}
	return
}
func (self SampleSize) Children() (r []Atom) {
	return
}

type TrackFragRun struct {
	Version          uint8
	Flags            uint32
	DataOffset       uint32
	FirstSampleFlags uint32
	Entries          []TrackFragRunEntry
	AtomPos
}

func (self TrackFragRun) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TRUN))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackFragRun) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	pio.PutU32BE(b[n:], uint32(len(self.Entries)))
	n += 4
	if self.Flags&TRUN_DATA_OFFSET != 0 {
		{
			pio.PutU32BE(b[n:], self.DataOffset)
			n += 4
		}
	}
	if self.Flags&TRUN_FIRST_SAMPLE_FLAGS != 0 {
		{
			pio.PutU32BE(b[n:], self.FirstSampleFlags)
			n += 4
		}
	}

	for i, entry := range self.Entries {
		var flags uint32
		if i > 0 {
			flags = self.Flags
		} else {
			flags = self.FirstSampleFlags
		}
		if flags&TRUN_SAMPLE_DURATION != 0 {
			pio.PutU32BE(b[n:], entry.Duration)
			n += 4
		}
		if flags&TRUN_SAMPLE_SIZE != 0 {
			pio.PutU32BE(b[n:], entry.Size)
			n += 4
		}
		if flags&TRUN_SAMPLE_FLAGS != 0 {
			pio.PutU32BE(b[n:], entry.Flags)
			n += 4
		}
		if flags&TRUN_SAMPLE_CTS != 0 {
			pio.PutU32BE(b[n:], entry.Cts)
			n += 4
		}
	}
	return
}
func (self TrackFragRun) Len() (n int) {
	n += 8
	n += 1
	n += 3
	n += 4
	if self.Flags&TRUN_DATA_OFFSET != 0 {
		{
			n += 4
		}
	}
	if self.Flags&TRUN_FIRST_SAMPLE_FLAGS != 0 {
		{
			n += 4
		}
	}

	for i := range self.Entries {
		var flags uint32
		if i > 0 {
			flags = self.Flags
		} else {
			flags = self.FirstSampleFlags
		}
		if flags&TRUN_SAMPLE_DURATION != 0 {
			n += 4
		}
		if flags&TRUN_SAMPLE_SIZE != 0 {
			n += 4
		}
		if flags&TRUN_SAMPLE_FLAGS != 0 {
			n += 4
		}
		if flags&TRUN_SAMPLE_CTS != 0 {
			n += 4
		}
	}
	return
}
func (self *TrackFragRun) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	var _len_Entries uint32
	_len_Entries = pio.U32BE(b[n:])
	n += 4
	self.Entries = make([]TrackFragRunEntry, _len_Entries)
	if self.Flags&TRUN_DATA_OFFSET != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("DataOffset", n+offset, err)
				return
			}
			self.DataOffset = pio.U32BE(b[n:])
			n += 4
		}
	}
	if self.Flags&TRUN_FIRST_SAMPLE_FLAGS != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("FirstSampleFlags", n+offset, err)
				return
			}
			self.FirstSampleFlags = pio.U32BE(b[n:])
			n += 4
		}
	}

	for i := 0; i < int(_len_Entries); i++ {
		var flags uint32
		if i > 0 {
			flags = self.Flags
		} else {
			flags = self.FirstSampleFlags
		}
		entry := &self.Entries[i]
		if flags&TRUN_SAMPLE_DURATION != 0 {
			entry.Duration = pio.U32BE(b[n:])
			n += 4
		}
		if flags&TRUN_SAMPLE_SIZE != 0 {
			entry.Size = pio.U32BE(b[n:])
			n += 4
		}
		if flags&TRUN_SAMPLE_FLAGS != 0 {
			entry.Flags = pio.U32BE(b[n:])
			n += 4
		}
		if flags&TRUN_SAMPLE_CTS != 0 {
			entry.Cts = pio.U32BE(b[n:])
			n += 4
		}
	}
	return
}
func (self TrackFragRun) Children() (r []Atom) {
	return
}

type TrackFragRunEntry struct {
	Duration uint32
	Size     uint32
	Flags    uint32
	Cts      uint32
}

func GetTrackFragRunEntry(b []byte) (self TrackFragRunEntry) {
	self.Duration = pio.U32BE(b[0:])
	self.Size = pio.U32BE(b[4:])
	self.Flags = pio.U32BE(b[8:])
	self.Cts = pio.U32BE(b[12:])
	return
}
func PutTrackFragRunEntry(b []byte, self TrackFragRunEntry) {
	pio.PutU32BE(b[0:], self.Duration)
	pio.PutU32BE(b[4:], self.Size)
	pio.PutU32BE(b[8:], self.Flags)
	pio.PutU32BE(b[12:], self.Cts)
}

const LenTrackFragRunEntry = 16

type TrackFragHeader struct {
	Version         uint8
	Flags           uint32
	BaseDataOffset  uint64
	StsdId          uint32
	DefaultDuration uint32
	DefaultSize     uint32
	DefaultFlags    uint32
	AtomPos
}

func (self TrackFragHeader) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TFHD))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackFragHeader) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	if self.Flags&TFHD_BASE_DATA_OFFSET != 0 {
		{
			pio.PutU64BE(b[n:], self.BaseDataOffset)
			n += 8
		}
	}
	if self.Flags&TFHD_STSD_ID != 0 {
		{
			pio.PutU32BE(b[n:], self.StsdId)
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_DURATION != 0 {
		{
			pio.PutU32BE(b[n:], self.DefaultDuration)
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_SIZE != 0 {
		{
			pio.PutU32BE(b[n:], self.DefaultSize)
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_FLAGS != 0 {
		{
			pio.PutU32BE(b[n:], self.DefaultFlags)
			n += 4
		}
	}
	return
}
func (self TrackFragHeader) Len() (n int) {
	n += 8
	n += 1
	n += 3
	if self.Flags&TFHD_BASE_DATA_OFFSET != 0 {
		{
			n += 8
		}
	}
	if self.Flags&TFHD_STSD_ID != 0 {
		{
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_DURATION != 0 {
		{
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_SIZE != 0 {
		{
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_FLAGS != 0 {
		{
			n += 4
		}
	}
	return
}
func (self *TrackFragHeader) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if self.Flags&TFHD_BASE_DATA_OFFSET != 0 {
		{
			if len(b) < n+8 {
				err = parseErr("BaseDataOffset", n+offset, err)
				return
			}
			self.BaseDataOffset = pio.U64BE(b[n:])
			n += 8
		}
	}
	if self.Flags&TFHD_STSD_ID != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("StsdId", n+offset, err)
				return
			}
			self.StsdId = pio.U32BE(b[n:])
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_DURATION != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("DefaultDuration", n+offset, err)
				return
			}
			self.DefaultDuration = pio.U32BE(b[n:])
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_SIZE != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("DefaultSize", n+offset, err)
				return
			}
			self.DefaultSize = pio.U32BE(b[n:])
			n += 4
		}
	}
	if self.Flags&TFHD_DEFAULT_FLAGS != 0 {
		{
			if len(b) < n+4 {
				err = parseErr("DefaultFlags", n+offset, err)
				return
			}
			self.DefaultFlags = pio.U32BE(b[n:])
			n += 4
		}
	}
	return
}
func (self TrackFragHeader) Children() (r []Atom) {
	return
}

type TrackFragDecodeTime struct {
	Version uint8
	Flags   uint32
	Time    time.Time
	AtomPos
}

func (self TrackFragDecodeTime) Marshal(b []byte) (n int) {
	pio.PutU32BE(b[4:], uint32(TFDT))
	n += self.marshal(b[8:]) + 8
	pio.PutU32BE(b[0:], uint32(n))
	return
}
func (self TrackFragDecodeTime) marshal(b []byte) (n int) {
	pio.PutU8(b[n:], self.Version)
	n += 1
	pio.PutU24BE(b[n:], self.Flags)
	n += 3
	if self.Version != 0 {
		PutTime64(b[n:], self.Time)
		n += 8
	} else {

		PutTime32(b[n:], self.Time)
		n += 4
	}
	return
}
func (self TrackFragDecodeTime) Len() (n int) {
	n += 8
	n += 1
	n += 3
	if self.Version != 0 {
		n += 8
	} else {

		n += 4
	}
	return
}
func (self *TrackFragDecodeTime) Unmarshal(b []byte, offset int) (n int, err error) {
	(&self.AtomPos).setPos(offset, len(b))
	n += 8
	if len(b) < n+1 {
		err = parseErr("Version", n+offset, err)
		return
	}
	self.Version = pio.U8(b[n:])
	n += 1
	if len(b) < n+3 {
		err = parseErr("Flags", n+offset, err)
		return
	}
	self.Flags = pio.U24BE(b[n:])
	n += 3
	if self.Version != 0 {
		self.Time = GetTime64(b[n:])
		n += 8
	} else {

		self.Time = GetTime32(b[n:])
		n += 4
	}
	return
}
func (self TrackFragDecodeTime) Children() (r []Atom) {
	return
}
