// License: Apache License 2.0
// Author: David Anderson
// URL: https://code.google.com/p/curvecp/source/browse/ringbuf/ringbuf.go

// Package ringbuf implements a byte ring buffer. The interface is
// close to that of an io.ReadWriter, but note that the semantics
// differ in significant ways, because this ringbuf is an
// implementation detail of the curvecp package, and it was more
// convenient like this.
package ringbuf

import "fmt"

type Ringbuf struct {
	buf []byte

	// start is the index of the oldest byte in the slice
	// the the index where reading will start
	start int

	// size is the number of data bytes in the slice
	// when the slice is less than full, size is less than len(buf)
	// but after that, size == len(buf)
	size int
}

// New creates a new ring buffer of the given size.
func New(size int) *Ringbuf {
	fmt.Println("making a New Ringbuf of size = ", size)
	return &Ringbuf{make([]byte, size), 0, 0}
}

// Write appends all the bytes in b to the buffer, looping and overwriting
// as needed, while incrementing the start to point to the start of the
// buffer.
func (r *Ringbuf) Write(b []byte) {
	//	fmt.Println("attempting to write ", string(b), "to the ringbuf ", r)
	for len(b) > 0 {

		// the last write ended here
		// r_end must never be greater than 2x len(r.buf)
		//   otherwise, the read would contain duplicate data
		r_end := r.start + r.size

		// so now begin writing after the end of the last write
		start := r_end % len(r.buf)

		fmt.Println("main loop r.start =", r.start, "r.size =", r.size, "start =", start, "r =", r)

		// given the current starting point
		// copy from b into r.buf from start to the end of r.buf
		// n is the number of bytes copied
		n := copy(r.buf[start:], b)

		// now that the beginning of b has been copied into r.buf,
		// reset b to contain only the end of the original input
		b = b[n:]

		// r.size and r.start must now be prepared for the next write

		// change the start only when the buffer is full
		if r.size == len(r.buf) {
			r.start = (start + n) % len(r.buf)
		} else {
			// increase the size until the buffer is full
			r.size += n
		}

		// update start
		// if r.size >= len(r.buf) {
		// 	if n <= len(r.buf) {
		// 		r.start += n
		// 		if r.start >= len(r.buf) {
		// 			r.start = 0
		// 		}
		// 	} else {
		// 		r.start = 0
		// 	}
		// }

		// the size should only change if r.buf is not at capacity already

	}
}

// Read reads as many bytes as possible from the ring buffer into
// b. Returns the number of bytes read.
func (r *Ringbuf) Read(b []byte) int {
	read := 0
	size := r.size
	start := r.start
	for len(b) > 0 && size > 0 {
		end := start + size
		if end > len(r.buf) {
			end = len(r.buf)
		}
		n := copy(b, r.buf[start:end])
		size -= n
		read += n
		b = b[n:]
		start = (start + n) % len(r.buf)
	}
	return read
}

func (r *Ringbuf) Size() int {
	return r.size
}
