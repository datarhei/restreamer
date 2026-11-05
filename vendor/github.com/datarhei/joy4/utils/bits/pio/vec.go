package pio

func VecLen(vec [][]byte) (n int) {
	for _, b := range vec {
		n += len(b)
	}
	return
}

func VecSliceTo(in [][]byte, out [][]byte, s int, e int) (n int) {
	if s < 0 {
		s = 0
	}

	if e >= 0 && e < s {
		panic("pio: VecSlice start > end")
	}

	i := 0
	off := 0
	for s > 0 && i < len(in) {
		left := len(in[i])
		read := s
		if left < read {
			read = left
		}
		left -= read
		off += read
		s -= read
		e -= read
		if left == 0 {
			i++
			off = 0
		}
	}
	if s > 0 {
		panic("pio: VecSlice start out of range")
	}

	for e != 0 && i < len(in) {
		left := len(in[i]) - off
		read := left
		if e > 0 && e < read {
			read = e
		}
		out[n] = in[i][off : off+read]
		n++
		left -= read
		e -= read
		off += read
		if left == 0 {
			i++
			off = 0
		}
	}
	if e > 0 {
		panic("pio: VecSlice end out of range")
	}

	return
}

func VecSlice(in [][]byte, s int, e int) (out [][]byte) {
	out = make([][]byte, len(in))
	n := VecSliceTo(in, out, s, e)
	out = out[:n]
	return
}
