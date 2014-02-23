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
	path   string
	file   *file.File
	err    error
	offset int64

	length        int
	write_bytes   int // number of bytes that should be written
	bytes_written int // number of bytes that were written

	op string // last operation attempted
	// buf must be a different size for each message, so not including it in the Writer struct.
}

func (p *Writer) Finish() error {
	w.err = w.file.Sync()
	if w.err != nil {
		w.op = "Finish Sync"
		return p.Error()
	}

	w.err = w.file.Close()
	if w.err != nil {
		w.op = "Finish Close"
		return p.Error()
	}
	return nil
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
func (w *Writer) Write(event []byte) (int, error) {
	w.op = "Write Validate"
	w.length = uint16(len(event))
	w.write_bytes = w.length + 2

	if w.length < 1 {
		w.err = fmt.Errorf("not appending a zero length event")
		return 0, w.Error()
	}
	if w.length > 65536 {
		w.err = fmt.Errorf("not appending an event larger than 65535 bytes")
		return 0, w.Error()
	}

	// 2 calls to write(2) are not guaranteed to be executed in order
	// combine length + event into a single buffer to pass to Write
	// fixed size byte slice for now.. not sure if there's a write method that takes [X]byte ????
	buf := make([]byte, w.write_bytes, w.write_bytes)
	buf[0] = byte(length >> 8)
	buf[1] = byte(length)
	buf[2:] = event

	w.bytes_written, w.err = w.file.Write(buf)
	if w.err != nil {
		w.op = "Write Write"
		return 0, p.Error()
	}

	if w.write_bytes != w.bytes_written {
		// recover. sync, close, open, find the last valid message, and then continue.
		// do this here so that the fast-path of successful writes is maintained
		w.err = w.Finish()
		if w.err != nil {
			return 0, p.Error()
		}

		w.file, w.err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if w.err != nil {
			w.op = "Write Recover Open"
			return 0, w.Error()
		}

		// find the correct place to sync to to begin writting again??????????
		// start at the beginning and read in chunks till EOF and event length don't add up.
		// then seek to that position and finish the current write!

		// ???????? maybe this should be done on NewWriter every time??????????   probably a good idea... and should not be that expensive.
		// since the fast-track case should be open empty file and write till full

	}

	// http://www.youtube.com/watch?v=WRAKFG1xqxA  class about the vfs cache etc...
	// http://www.thomas-krenn.com/en/wiki/Linux_Page_Cache_Basics
	// http://www.westnet.com/~gsmith/content/linux-pdflush.htm

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
filesystem optimization for Append:
http://kernelnewbies.org/Ext4
https://ext4.wiki.kernel.org/index.php/Ext4_Disk_Layout
http://digital-forensics.sans.org/blog/2010/12/20/digital-forensics-understanding-ext4-part-1-extents




filesystem optimization for Read:




mount -t debugfs nodev /sys/kernel/debug
https://access.redhat.com/site/documentation/en-US/Red_Hat_Enterprise_MRG/2/html/Realtime_Tuning_Guide/sect-Realtime_Tuning_Guide-Realtime_Specific_Tuning-Using_the_ftrace_Utility_for_Tracing_Latencies.html


*/

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
