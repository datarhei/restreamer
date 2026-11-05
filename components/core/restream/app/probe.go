package app

type ProbeIO struct {
	Address string

	// General
	Index    uint64
	Stream   uint64
	Language string
	Format   string
	Type     string
	Codec    string
	Coder    string
	Bitrate  float64 // kbit/s
	Duration float64

	// Video
	Pixfmt string
	Width  uint64
	Height uint64
	FPS    float64

	// Audio
	Sampling uint64
	Layout   string
	Channels uint64
}

type Probe struct {
	Streams []ProbeIO
	Log     []string
}
