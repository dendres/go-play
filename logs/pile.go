/*
package pile implements an append only file that almost never calls fsync(2).
each file entry is a 2 byte uint16 length followed by the entry data.
file entries are not guaranteed to be written in the order received.
the OS chooses when to flush pages to disk
the page flush rate increases with memory pressure and IO load.

2 types are implemented. each has one key method:
Writer.Write()
Reader.Next()

*/

package pile

import (
	"fmt"
	"io"
	"os"
)

type Writer struct {
	file   *file.File
	err    error
	offset int64

	length int
	write_bytes int // number of bytes that should be written
	bytes_written int // number of bytes that were written

	op     string // last operation attempted
	// buf must be a different size for each message, so not including it in the Writer struct.
}

func (p *Writer) Finish() {
	p.file.Sync()
	p.file.Close()
}

func (p *Writer) Error() error {
	p.Finish()
	return fmt.Errorf("file = %s, offset = %d, op = %s, error: %s",
		p.file.Name(), p.offset, p.op, p.err.Error())
}

// NewWriter creates the file at path if it is missing, then opens it and returns a new Write.
// it does NOT create parent directories
func NewWriter(path string) (*Writer, error) {
	w := &Writer{path}
	w.file, w.err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if w.err != nil {
		return w, w.Error()
	}
	return w, nil
}

// Writes appends len(event) and event []byte to the end of the file
// it does not fsync or close the file
func (p *Writer) Write(event []byte) (int, error) {
	w.op = "Write Validate"
	length = uint16(len(event))
	write_bytes = length + 2

	if length < 1 {
		w.err = fmt.Errorf("not appending a zero length event")
		return 0, w.Error()
	}
	if length > 65536 {
		w.err = fmt.Errorf("not appending an event larger than 65535 bytes")
		return 0, w.Error()
	}

	// 2 calls to write(2) are not guaranteed to be executed in order
	// combine length + event into a single buffer to pass to Write
	// fixed size byte slice for now.. not sure if there's a write method that takes [X]byte ????
	buf := make([]byte, write_bytes, write_bytes)
	buf[0] = byte(length >> 8)
	buf[1] = byte(length)
	buf[2:] = event


	w.written, w.err = w.file.Write(buf)
	if w.err != nil {
		w.op = "Write Write"
		return 0, p.Error()
	}
	if w.written != length+2


	w.offset, w.err = w.file.Seek(0, 2)
	if w.err != nil {
		w.op = "Write Seek"
		return 0, p.Error()
	}


}

type Reader struct {
}

func (*Reader) Read(p []byte) (n int, err error)

/*
2 byte size. 64K max element.

NewWriter
Writer
 Write

NewPile

Next


methods:
* append([]byte)
* read(offset)
* scan()
* reverse()
* recover()
* recovergzip()
*/

/*
s := NewStore
event := s.Next()



*/

/*
store.Append(b []byte)
* offset = os.Seek(0, 2) // EOF
* os.Write(len(b))
* os.Write([]byte)
* return offset // new id for the event if needed
*/

/*
store.Read(id int64)
* os.Seek(id, 0) // seek to the offset given by id
* l := make([]byte, 2)
* os.Read(l)
* event_length := int(l)
* e := make([]byte, event_length)
* os.Seek(2,1)
* os.Read(e)
* return e

*/

/*
read the whole thing out into a stream of events????
read l, seek 2, read len(l), emit event
repeat
*/

/*
repair corrupt file?
 raw events... use crc32
 gzip... use the crc and length in the gzip

read and discard bytes until gzip magic byte. backup 2 and read length.?????

*/

/*
http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file

*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
