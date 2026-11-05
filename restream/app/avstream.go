package app

type AVstreamIO struct {
	State  string
	Packet uint64 // counter
	Time   uint64 // sec
	Size   uint64 // bytes
}

type AVstream struct {
	Input          AVstreamIO
	Output         AVstreamIO
	Aqueue         uint64 // gauge
	Queue          uint64 // gauge
	Dup            uint64 // counter
	Drop           uint64 // counter
	Enc            uint64 // counter
	Looping        bool
	LoopingRuntime uint64 // sec
	Duplicating    bool
	GOP            string
}
