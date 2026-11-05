package pktque

import (
	"github.com/datarhei/joy4/av"
)

type Buf struct {
	Head, Tail BufPos
	pkts       []av.Packet
	Size       int
	Count      int
}

func NewBuf() *Buf {
	return &Buf{
		pkts: make([]av.Packet, 64),
	}
}

func (self *Buf) Pop() av.Packet {
	if self.Count == 0 {
		panic("pktque.Buf: Pop() when count == 0")
	}

	i := int(self.Head) & (len(self.pkts) - 1)
	pkt := self.pkts[i]
	self.pkts[i] = av.Packet{}
	self.Size -= len(pkt.Data)
	self.Head++
	self.Count--

	return pkt
}

func (self *Buf) grow() {
	newpkts := make([]av.Packet, len(self.pkts)*2)
	for i := self.Head; i.LT(self.Tail); i++ {
		newpkts[int(i)&(len(newpkts)-1)] = self.pkts[int(i)&(len(self.pkts)-1)]
	}
	self.pkts = newpkts
}

func (self *Buf) Push(pkt av.Packet) {
	if self.Count == len(self.pkts) {
		self.grow()
	}
	self.pkts[int(self.Tail)&(len(self.pkts)-1)] = pkt
	self.Tail++
	self.Count++
	self.Size += len(pkt.Data)
}

func (self *Buf) Get(pos BufPos) av.Packet {
	return self.pkts[int(pos)&(len(self.pkts)-1)]
}

func (self *Buf) IsValidPos(pos BufPos) bool {
	return pos.GE(self.Head) && pos.LT(self.Tail)
}

type BufPos int

func (self BufPos) LT(pos BufPos) bool {
	return self-pos < 0
}

func (self BufPos) GE(pos BufPos) bool {
	return self-pos >= 0
}

func (self BufPos) GT(pos BufPos) bool {
	return self-pos > 0
}
