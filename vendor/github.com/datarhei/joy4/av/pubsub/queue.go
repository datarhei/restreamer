// Packege pubsub implements publisher-subscribers model used in multi-channel streaming.
package pubsub

import (
	"io"
	"sync"
	"time"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/av/pktque"
)

//        time
// ----------------->
//
// V-A-V-V-A-V-V-A-V-V
// |                 |
// 0        5        10
// head             tail
// oldest          latest
//

// One publisher and multiple subscribers thread-safe packet buffer queue.
type Queue struct {
	buf                      *pktque.Buf
	head, tail               int
	lock                     *sync.RWMutex
	cond                     *sync.Cond
	curgopcount, maxgopcount int
	streams                  []av.CodecData
	videoidx                 int
	closed                   bool
}

func NewQueue() *Queue {
	q := &Queue{}
	q.buf = pktque.NewBuf()
	q.maxgopcount = 2
	q.lock = &sync.RWMutex{}
	q.cond = sync.NewCond(q.lock.RLocker())
	q.videoidx = -1
	return q
}

func (self *Queue) SetMaxGopCount(n int) {
	self.lock.Lock()
	self.maxgopcount = n
	self.lock.Unlock()
	return
}

func (self *Queue) WriteHeader(streams []av.CodecData) error {
	self.lock.Lock()

	self.streams = streams
	for i, stream := range streams {
		if stream.Type().IsVideo() {
			self.videoidx = i
		}
	}
	self.cond.Broadcast()

	self.lock.Unlock()

	return nil
}

func (self *Queue) WriteTrailer() error {
	return nil
}

// After Close() called, all QueueCursor's ReadPacket will return io.EOF.
func (self *Queue) Close() (err error) {
	self.lock.Lock()

	self.closed = true
	self.cond.Broadcast()

	self.lock.Unlock()
	return
}

// Put packet into buffer, old packets will be discared.
func (self *Queue) WritePacket(pkt av.Packet) (err error) {
	self.lock.Lock()

	self.buf.Push(pkt)
	if pkt.Idx == int8(self.videoidx) && pkt.IsKeyFrame {
		self.curgopcount++
	}

	for self.curgopcount >= self.maxgopcount && self.buf.Count > 1 {
		pkt := self.buf.Pop()
		if pkt.Idx == int8(self.videoidx) && pkt.IsKeyFrame {
			self.curgopcount--
		}
		if self.curgopcount < self.maxgopcount {
			break
		}
	}

	self.cond.Broadcast()

	self.lock.Unlock()
	return
}

type QueueCursor struct {
	que    *Queue
	pos    pktque.BufPos
	gotpos bool
	init   func(buf *pktque.Buf, videoidx int) pktque.BufPos
}

func (self *Queue) newCursor() *QueueCursor {
	return &QueueCursor{
		que: self,
	}
}

// Create cursor position at latest packet.
func (self *Queue) Latest() *QueueCursor {
	cursor := self.newCursor()
	cursor.init = func(buf *pktque.Buf, videoidx int) pktque.BufPos {
		return buf.Tail
	}
	return cursor
}

// Create cursor position at oldest buffered packet.
func (self *Queue) Oldest() *QueueCursor {
	cursor := self.newCursor()
	cursor.init = func(buf *pktque.Buf, videoidx int) pktque.BufPos {
		return buf.Head
	}
	return cursor
}

// Create cursor position at specific time in buffered packets.
func (self *Queue) DelayedTime(dur time.Duration) *QueueCursor {
	cursor := self.newCursor()
	cursor.init = func(buf *pktque.Buf, videoidx int) pktque.BufPos {
		i := buf.Tail - 1
		if buf.IsValidPos(i) {
			end := buf.Get(i)
			for buf.IsValidPos(i) {
				if end.Time-buf.Get(i).Time > dur {
					break
				}
				i--
			}
		}
		return i
	}
	return cursor
}

// Create cursor position at specific delayed GOP count in buffered packets.
func (self *Queue) DelayedGopCount(n int) *QueueCursor {
	cursor := self.newCursor()
	cursor.init = func(buf *pktque.Buf, videoidx int) pktque.BufPos {
		i := buf.Tail - 1
		if videoidx != -1 {
			for gop := 0; buf.IsValidPos(i) && gop < n; i-- {
				pkt := buf.Get(i)
				if pkt.Idx == int8(self.videoidx) && pkt.IsKeyFrame {
					gop++
				}
			}
		}
		return i
	}
	return cursor
}

func (self *QueueCursor) Streams() (streams []av.CodecData, err error) {
	self.que.cond.L.Lock()
	for self.que.streams == nil && !self.que.closed {
		self.que.cond.Wait()
	}
	if self.que.streams != nil {
		streams = self.que.streams
	} else {
		err = io.EOF
	}
	self.que.cond.L.Unlock()
	return
}

// ReadPacket will not consume packets in Queue, it's just a cursor.
func (self *QueueCursor) ReadPacket() (pkt av.Packet, err error) {
	self.que.cond.L.Lock()
	buf := self.que.buf
	if !self.gotpos {
		self.pos = self.init(buf, self.que.videoidx)
		self.gotpos = true
	}
	for {
		if self.pos.LT(buf.Head) {
			self.pos = buf.Head
		} else if self.pos.GT(buf.Tail) {
			self.pos = buf.Tail
		}
		if buf.IsValidPos(self.pos) {
			pkt = buf.Get(self.pos)
			self.pos++
			break
		}
		if self.que.closed {
			err = io.EOF
			break
		}
		self.que.cond.Wait()
	}
	self.que.cond.L.Unlock()
	return
}
