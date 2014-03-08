/*
package pile implements an append-only file that rarely calls fsync(2).
"events" with a very specific layout are written to the file.
"events" are not guaranteed to be written in the order received.
the OS chooses when to flush pages to disk
the page flush rate increases with memory pressure and IO load.

# XXX use the Event and EventByte Types!!!!!

*/
package pile

import (
	"container/list"
	"fmt"
	"github.com/dendres/go-play/event"
	"io"
	"os"
)

// A Pile is an interface to an append-only file with weak guarantees about ordering, duplication, and syncing.
type Pile struct {
	path   string
	writer *os.File
	op     string // last operation attempted
	err    error  // last error received
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
	p := new(Pile)
	p.path = path
	p.writer, p.err = os.OpenFile(p.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if p.err != nil {
		p.op = "New"
		return p, p.Error()
	}
	return p, nil
}

// Append adds event []byte to the end of the file and advances the fd offset.
// It does not validate any part of the event.
// It does not fsync or close the file.
// File.Write loops forever to prevent short write, so no need to check for short write here.
func (p *Pile) Append(eb *event.EventBytes) error {
	p.op = "Append"
	_, p.err = p.writer.Write(eb.Bytes())
	if p.err != nil {
		return p.Error()
	}
	return nil
}

// Read reads all events in a file onto the given List.
func Read(p *Pile, events *list.List) error {
	var reader *os.File
	reader, p.err = os.Open(p.path)
	if p.err != nil {
		p.op = "Read Open"
		return p.Error()
	}

	event_header := event.NewEventHeaderBuffer()
	data_length := int(0)

	for {
		_, p.err = reader.Read(event_header)
		if p.err == io.EOF {
			break
		}
		if p.err != nil {
			p.op = "Read Header"
			return p.Error()
		}

		event.ReadDataLength(event_header, &data_length)

		// event_remainder is a different size each event so the slice cannot be reused
		event_remainder := make([]byte, data_length+1)
		_, p.err = reader.Read(event_remainder)
		if p.err != nil {
			p.op = "Read Event"
			return p.Error()
		}

		event := append(event_header, event_remainder...) // XXX might be a lot of copy and resize operations???
		events.PushBack(event)
	}
	return nil
}

/*
	n.fileinfo, n.err := f.Stat()
	if err != nil {
		n.op = "Stat"
		return n, n.Error()
	}
	n.filesize = n.fileinfo.Size()


filesystem optimization for Append:
http://kernelnewbies.org/Ext4
https://ext4.wiki.kernel.org/index.php/Ext4_Disk_Layout
http://digital-forensics.sans.org/blog/2010/12/20/digital-forensics-understanding-ext4-part-1-extents
http://www.youtube.com/watch?v=WRAKFG1xqxA  class about the vfs cache etc...
http://www.thomas-krenn.com/en/wiki/Linux_Page_Cache_Basics
http://www.westnet.com/~gsmith/content/linux-pdflush.htm
http://superuser.com/questions/479379/how-long-can-file-system-writes-be-cached-with-ext4


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
