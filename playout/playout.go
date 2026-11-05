package playout

type StatusIO struct {
	State  string `json:"state"`
	Packet uint64 `json:"packet"`
	Time   uint64 `json:"time"`
	Size   uint64 `json:"size_kb"`
}

type StatusSwap struct {
	Address     string `json:"url"`
	Status      string `json:"status"`
	LastAddress string `json:"lasturl"`
	LastError   string `json:"lasterror"`
}

type Status struct {
	ID          string      `json:"id"`
	Address     string      `json:"url"`
	Stream      uint64      `json:"stream"`
	Queue       uint64      `json:"queue"`
	AQueue      uint64      `json:"aqueue"`
	Dup         uint64      `json:"dup"`
	Drop        uint64      `json:"drop"`
	Enc         uint64      `json:"enc"`
	Looping     bool        `json:"looping"`
	Duplicating bool        `json:"duplicating"`
	GOP         string      `json:"gop"`
	Debug       interface{} `json:"debug"`
	Input       StatusIO    `json:"input"`
	Output      StatusIO    `json:"output"`
	Swap        StatusSwap  `json:"swap"`
}
