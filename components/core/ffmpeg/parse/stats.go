package parse

type statsData struct {
	frame  uint64 // counter
	packet uint64 // counter
	size   uint64 // bytes
	dup    uint64 // counter
	drop   uint64 // counter
}

type stats struct {
	last statsData
	diff statsData
}

func (s *stats) updateFromProgress(p *ffmpegProgress) {
	s.diff.frame = p.Frame - s.last.frame
	s.diff.packet = p.Packet - s.last.packet
	s.diff.size = p.Size - s.last.size
	s.diff.drop = p.Drop - s.last.drop
	s.diff.dup = p.Dup - s.last.dup

	s.last.frame = p.Frame
	s.last.packet = p.Packet
	s.last.size = p.Size
	s.last.dup = p.Dup
	s.last.drop = p.Drop
}

func (s *stats) updateFromProgressIO(p *ffmpegProgressIO) {
	s.diff.frame = p.Frame - s.last.frame
	s.diff.packet = p.Packet - s.last.packet
	s.diff.size = p.Size - s.last.size

	s.last.frame = p.Frame
	s.last.packet = p.Packet
	s.last.size = p.Size
}
