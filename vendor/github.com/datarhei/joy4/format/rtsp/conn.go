package rtsp

import (
	"net"
	"time"
)

type connWithTimeout struct {
	Timeout time.Duration
	net.Conn
}

func (self connWithTimeout) Read(p []byte) (n int, err error) {
	if self.Timeout > 0 {
		self.Conn.SetReadDeadline(time.Now().Add(self.Timeout))
	}
	return self.Conn.Read(p)
}

func (self connWithTimeout) Write(p []byte) (n int, err error) {
	if self.Timeout > 0 {
		self.Conn.SetWriteDeadline(time.Now().Add(self.Timeout))
	}
	return self.Conn.Write(p)
}
