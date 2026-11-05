package pktque

import (
	"time"
)

/*
pop                                   push

     seg                 seg        seg
  |--------|         |---------|   |---|
     20ms                40ms       5ms
----------------- time -------------------->
headtm                               tailtm
*/

type tlSeg struct {
	tm, dur time.Duration
}

type Timeline struct {
	segs   []tlSeg
	headtm time.Duration
}

func (self *Timeline) Push(tm time.Duration, dur time.Duration) {
	if len(self.segs) > 0 {
		tail := self.segs[len(self.segs)-1]
		diff := tm - (tail.tm + tail.dur)
		if diff < 0 {
			tm -= diff
		}
	}
	self.segs = append(self.segs, tlSeg{tm, dur})
}

func (self *Timeline) Pop(dur time.Duration) (tm time.Duration) {
	if len(self.segs) == 0 {
		return self.headtm
	}

	tm = self.segs[0].tm
	for dur > 0 && len(self.segs) > 0 {
		seg := &self.segs[0]
		sub := dur
		if seg.dur < sub {
			sub = seg.dur
		}
		seg.dur -= sub
		dur -= sub
		seg.tm += sub
		self.headtm += sub
		if seg.dur == 0 {
			copy(self.segs[0:], self.segs[1:])
			self.segs = self.segs[:len(self.segs)-1]
		}
	}

	return
}
