package bits

import (
	"io"
)

type GolombBitReader struct {
	R    io.Reader
	buf  [1]byte
	left byte
}

func (self *GolombBitReader) ReadBit() (res uint, err error) {
	if self.left == 0 {
		if _, err = self.R.Read(self.buf[:]); err != nil {
			return
		}
		self.left = 8
	}
	self.left--
	res = uint(self.buf[0]>>self.left) & 1
	return
}

func (self *GolombBitReader) ReadBits(n int) (res uint, err error) {
	for i := 0; i < n; i++ {
		var bit uint
		if bit, err = self.ReadBit(); err != nil {
			return
		}
		res |= bit << uint(n-i-1)
	}
	return
}

func (self *GolombBitReader) ReadBits32(n uint) (r uint32, err error) {
	var t uint
	for i := uint(0); i < n; i++ {
		t, err = self.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | uint32(t)
	}
	return
}

func (self *GolombBitReader) ReadBits64(n uint) (r uint64, err error) {
	var t uint
	for i := uint(0); i < n; i++ {
		t, err = self.ReadBit()
		if err != nil {
			return
		}
		r = (r << 1) | uint64(t)
	}
	return
}

func (self *GolombBitReader) ReadExponentialGolombCode() (res uint, err error) {
	i := 0
	for {
		var bit uint
		if bit, err = self.ReadBit(); err != nil {
			return
		}
		if !(bit == 0 && i < 32) {
			break
		}
		i++
	}
	if res, err = self.ReadBits(i); err != nil {
		return
	}
	res += (1 << uint(i)) - 1
	return
}

func (self *GolombBitReader) ReadSE() (res uint, err error) {
	if res, err = self.ReadExponentialGolombCode(); err != nil {
		return
	}
	if res&0x01 != 0 {
		res = (res + 1) / 2
	} else {
		res = -res / 2
	}
	return
}
