/*
package bita implements a sparse file based bit array.
Given a uint64, seek to that offset in the file and read/write the single bit there.
Pass through os.File errors unaltered.
The bit array file has no header.
The caller must track the meaning of the offset/index externally.
NOT concurrent!
*/
package bita

import (
	"fmt"
	"io"
	"os"
)

// type Bita holds the handle to the bit array sparse file.
type Bita struct {
	// the file to treat like a bit array
	file *os.File

	// the last offset used by ReadAt or WriteAt
	offset int64

	// the bit mask used to read/write a single bit from the byte
	mask uint8

	// holds the byte read or written
	buf []byte

	// the last operation attempted
	op string

	// the last error encountered
	err error
}

func (b *Bita) Finish() {
	_ = b.file.Sync()
	_ = b.file.Close()
}

// Error cleans up and returns a string representation of the error.
func (b *Bita) Error() error {
	b.Finish()
	return fmt.Errorf("file = %s, offset = %d, mask = %d, buf = %v, op = %s, error: %s",
		b.file.Name(), b.offset, b.mask, b.buf, b.op, b.err.Error())
}

// Open opens an existing bit array file by path or creates a new empty one if needed and returns the new struct.
// There is no fixed size specification.
// The caller must track the desired size externally!
func Open(path string) (*Bita, error) {
	b := new(Bita)
	b.op = "OpenFile"
	b.buf = make([]byte, 1, 1)
	if b.file, b.err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666); b.err != nil {
		return b, b.Error()
	}
	return b, nil
}

// Set writes a single byte to the array
// seek, read, seek, write = 4 iops
func (b *Bita) Set(index int64) error {
	b.offset = index / 8
	b.mask = uint8(128) >> uint8(index%8)
	b.buf[0] = byte(0)

	_, b.err = b.file.ReadAt(b.buf, b.offset)
	if b.err == io.EOF {
		b.buf[0] = byte(0)
	} else if b.err != nil {
		b.op = "ReadAt"
		return b.Error()
	}

	fmt.Println("file", b.file.Name(), "old buffer", b.buf)
	b.buf[0] = b.buf[0] | b.mask
	fmt.Println("file", b.file.Name(), "new buffer", b.buf)

	if _, b.err = b.file.WriteAt(b.buf, b.offset); b.err != nil {
		b.op = "WriteAt"
		return b.Error()
	}

	return nil
}

// Get reads the bit at the given index and returns a bool.
// seek read = 2 iops
func (b *Bita) Get(index int64) (bool, error) {
	b.offset = index / 8
	b.mask = uint8(128) >> uint8(index%8)

	_, b.err = b.file.ReadAt(b.buf, b.offset)
	if b.err == io.EOF {
		b.buf[0] = byte(0)
	} else if b.err != nil {
		b.op = "ReadAt"
		return false, b.Error()
	}

	if b.buf[0]&b.mask == b.mask {
		return true, nil
	}
	return false, nil
}
