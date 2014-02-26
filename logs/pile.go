/*
package pile implements an append-only file that rarely calls fsync(2).
"events" with a very specific layout are written to the file.
"events" are not guaranteed to be written in the order received.
the OS chooses when to flush pages to disk
the page flush rate increases with memory pressure and IO load.

the "event" layout is:
* event[0:3] crc of event[4:]
* event[4] header of 4x 2 bit routing values
* event[5:12] time ns since unix epoch
* event[13:15] length of data not including EOE=0xFF
* event[16:len(event)-1] data (DO NOT READ)
* event[len(event)-1:] EOE=0xFF

pile.Append() Appends to the file
pile.Read() Returns a [][]byte of sorted events

XXXXXXXXXXXXXXXX need refactor for single Pile type!!!!!

*/

package pile

import (
	"fmt"
	"io"
	"os"
)

const headersize = 16
const headerfootersize = 17

// A Pile is an interface to an append-only file with weak guarantees about ordering, duplication, and syncing.
type Pile struct {
	path string
	writer *os.File
	op   string // last operation attempted
	err  error  // last error received
}

// Finish fsyncs and closes the file.
// Finish is called from Error, so it can't call Error on error!!
// if the Sync or Close error out, ignore the errors and continue XXXXXX
func (p *Pile) Finish() {
	p.writer.Sync()
	p.writer.Close()
}

// Error calls Finish() then returns an error object that shows state info
func (p *Pile) Error() error {
	p.Finish()
	return fmt.Errorf("file = %s, op = %s, error: %s", p.writer.Name(), p.op, p.err.Error())
}

// NewPile creates the file at path if it is missing, then opens it and returns a Pile.
// it does NOT create parent directories
func NewPile(path string) (*Pile, error) {
	p := &Pile{path}
	p.file, p.err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if p.err != nil {
		p.op = "New"
		return p, p.Error()
	}
	return p, nil
}

// Write appends event []byte to the end of the file and advances the fd offset.
// It does not validate any part of the event.
// It does not fsync or close the file.
// File.Write loops forever to prevent short write, so no need to check for short write here.
func (w *Writer) Write(event []byte) (int, error) {
	_, w.err = w.file.Write(event)
	if w.err != nil {
		w.op = "Write Write"
		return 0, w.Error()
	}
}

// XXX remove file related info... only needed in Sort()
type Nexter struct {
	path     string

	op       string
	err      error
	offset   int
	events   [][]byte // ~ 64MB of sorted events
	eheader  []byte   // buffer containing the last event header read
	elength  int      // the length value from the header of the last event to be read
}

// Sort reads all events in the file, sorts them by time, and returns a Nexter.
// path is parsed to allow some of the high order time bytes to be ignored.
func Sort(path string) (*Nexter, error) {
	n := &Nexter{path}
	eheader = make([]byte, headersize, headersize)


	file     *os.File
	fileinfo *os.FileInfo
	filesize int


	n.fileinfo, n.err := f.Stat()
	if err != nil {
		n.op = "Stat"
		return n, n.Error()
	}
	n.filesize = n.fileinfo.Size()


	n.file, n.err = os.Open(path)
	if n.err != nil {
		n.op = "New"
		return n, n.Error()
	}

	for n.offset < n.filesize {
		// read header
		// read whole event including 0xFF
		// increment offset by 

	return n, nil
}

// Next returns 1 event at a time and advances the offset.
func (n *Nexter) Next() error {

	_, n.err = n.file.Read(n.header)
	if n.err != nil {
		n.op = "Read Header"
		return n.Error()
	}

	n.elength = int(n.header[13:15])

	buf := make([]byte, n.elength)
	_, n.err = n.file.Read()
	if n.err != nil {
		n.op = "Read Header"
		return n.Error()
	}
}

/*
filesystem optimization for Append:
http://kernelnewbies.org/Ext4
https://ext4.wiki.kernel.org/index.php/Ext4_Disk_Layout
http://digital-forensics.sans.org/blog/2010/12/20/digital-forensics-understanding-ext4-part-1-extents
http://www.youtube.com/watch?v=WRAKFG1xqxA  class about the vfs cache etc...
http://www.thomas-krenn.com/en/wiki/Linux_Page_Cache_Basics
http://www.westnet.com/~gsmith/content/linux-pdflush.htm


filesystem optimization for Read:
mount -t debugfs nodev /sys/kernel/debug
https://access.redhat.com/site/documentation/en-US/Red_Hat_Enterprise_MRG/2/html/Realtime_Tuning_Guide/sect-Realtime_Tuning_Guide-Realtime_Specific_Tuning-Using_the_ftrace_Utility_for_Tracing_Latencies.html


	// 2 calls to write(2) are not guaranteed to be executed in order
	// combine length + event into a single buffer to pass to Write
	// fixed size byte slice for now.. not sure if there's a write method that takes [X]byte ????
	buf := make([]byte, w.write_bytes, w.write_bytes)
	buf[0] = byte(length >> 8)
	buf[1] = byte(length)
	buf[2:] = event



*/

/*
Recover closes the file, reopens the file, verifies the end of the file is good, and seeks to EOF for the next write.
errors likely to be returned by syscall.write == write(2):
ENOSPC: The device containing the file referred to by fd has no room for the data.
EFBIG: An attempt was made to write a file that exceeds the implementation-defined maximum file size
  or the process’s file size limit, or to write at a position past the maximum allowed offset.
EINTR: The call was interrupted by a signal before any data was written; see signal(7).
EBADF: fd is not a valid file descriptor or is not open for writing.
EINVAL: fd is attached to an object which is unsuitable for writing; or the file was opened with the O_DIRECT flag,
  and either the address specified in buf, the value specified in count, or the current file offset is not suitably aligned.
EIO: A low-level I/O error occurred while modifying the inode.
*/
// func (w *Writer) Recover() error {
// 	w.Finish()

// 	w.file, w.err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
// 	if w.err != nil {
// 		w.op = "Recover Open"
// 		return w.Error()
// 	}

// 	lastbyte := []byte{0}

// 	_, w.err = w.file.Seek(-1, 2)
// 	if w.err != nil {
// 		w.op = "Recover Seek"
// 		return w.Error()
// 	}

// 	_, w.err = w.file.Read(lastbyte)
// 	if w.err != nil {
// 		w.op = "Recover Read"
// 		return w.Error()
// 	}

// 	if lastbyte[0] != 0xFF {
// 		fmt.Println("Found an incomplete Write!!! XXXX Recovery process not yet implemented!!!!!")
// 		// find the last 0xFF in the file or the beginning of the file
// 		// for i := 0; i < 100; i++ {
// 		//   Seek(1 - i) and read 100 bytes
// 		//   loop over the 100 bytes. if 0xFF is found, calculate the location of 0xFF from the loop counters
// 	}
// 	return nil
// }

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
