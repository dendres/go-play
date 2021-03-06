/*
package pile appends and reads EventBytes while rarely calling fsync(2).
the OS chooses when to flush pages to diks
the page flush rate increases with memory pressure and IO load.
"events" are not guaranteed to be written in the order received.
*/
package pile

import (
	"fmt"
	"github.com/dendres/go-play/event"
	"io"
	"os"
)

// A Pile is an interface to an append-only file with weak guarantees about ordering, duplication, and syncing.
type Pile struct {
	// the path to the file containing the offset delimited events
	Path string

	// the handle to the file
	writer *os.File

	// the approximate number of bytes in the file
	written int64

	// the last operation attempted
	op string

	// the last error received
	err error
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
	return fmt.Errorf("path = %s, op = %s, error: %s", p.Path, p.op, p.err.Error())
}

// Size performs a stat(2) and returns the current file size in bytes as int64
func (p *Pile) Size() (int64, error) {
	fi, p.err = p.writer.Stat()
	if p.err != nil {
		p.op = "Stat"
		return 0, p.Error()
	}

	return fileinfo.Size(), nil
}

// NewPile creates the file at path if it is missing, then opens it and returns a Pile.
// it does NOT create parent directories
func NewPile(path string) (*Pile, error) {
	p := new(Pile)
	p.Path = path
	p.writer, p.err = os.OpenFile(p.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if p.err != nil {
		p.op = "New"
		return p, p.Error()
	}

	p.written, p.err = p.Size()
	if p.err != nil {
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
	b := eb.GetBytes()
	_, p.err = p.writer.Write(b)
	if p.err != nil {
		return p.Error()
	}

	p.written += len(b)

	return nil
}

// Read up to len(events) events from the file.
func (p *Pile) Read(events []*event.EventBytes) error {
	var reader *os.File
	reader, p.err = os.Open(p.Path)
	if p.err != nil {
		p.op = "Read Open"
		return p.Error()
	}

	// event_header is fixed length, so it can be reused
	event_header := event.NewEventHeaderBuffer()

	for i := 0; i < len(events); i++ {
		_, p.err = reader.Read(event_header)
		if p.err == io.EOF {
			break
		}
		if p.err != nil {
			p.op = "Read Header"
			return p.Error()
		}

		var event_remainder []byte
		event_remainder, p.err = event.NewEventRemainderBuffer(event_header)
		if p.err != nil {
			p.op = "New Event Remainder buffer"
			return p.Error()
		}

		_, p.err = reader.Read(event_remainder)
		if p.err != nil {
			// if EOF or any read error appears while trying to read an event, the caller is expected to retry??
			p.op = "Read Event"
			return p.Error()
		}

		events[i] = event.NewEventFromBuffers(event_header, event_remainder)
	}
	return nil
}

/*


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
