/*
package drip implements Drip.
*/
package drip

import (
	"fmt"
	"github.com/dendres/go-play/event"
	"github.com/dendres/go-play/pile"
)

// A Drip listens on an EventBytes channel and, writes to a Pile.
type Drip struct {

	// the event pile
	ep *pile.Pile

	// the incoming channel of events
	events chan<- *EventBytes

	// notify Hose that it's time to split
	split <-chan bool

	// stop when the pile is larger than splitsize bytes
	splitsize int64

	// clean up and exit after this inactivity timeout in nanoseconds
	bored int

	// the last operation attempted
	op string

	// the last error encountered
	err error
}

// Finish performs any cleanup required before exiting
func (d *Drip) Finish() {
	close(events)
	d.ep.Finish()
}

// Error cleans up and returns a string representation of the error.
func (d *Drip) Error() error {
	b.Finish()
	return fmt.Errorf("path = %s, op = %s, error: %s", d.ep.Path, d.op, d.err.Error())
}

// XXX rough out Hose to get the methods required

/*
splitting tree on FS:
* ideal time = 35 bit second + 30 bit ns. can reduce fractional second granularity
* 60 bit fits in 12 lp32 characters
* 65 bit fits in 13
* 32 bit of crc ~ 6 characters
* number of files in 3 characters = 32^3 = 32768
* can't list in these directories!
* shield them from mlocate
* any other system processes to wory about???

bucket/time/ss
* 12 days / 1024 x 17min intervals

bucket/time/sss/sss/ssc/ccc/rrr/ccc

fixed file size to balance split rate vs. file scan time. vs. memory pressure:
* 500K * 512 byte events/second = 244MB/second = 20TB/day
* multiple of fs block = 4K
* split rate max of 1/second => 256MB, 2/second => 128MB files, 4/s => 64MB files
* try to allow the file size to not matter to the read/write process so it can be adjusted????
* allow lazy splitting to happen during reduced load???

Message count and size estimates for a fixed 128MB file:
* (8 + 4 + 2, len(E)) * N = max size
* 0 len string produces serialized event size ~ 45 bytes including time crc and len
* 64K max + 14 = 65550
* max number of min sized events in 128MB = 128 * 1024 * 1024 / 45 = 3M events
* max number of average sized events      = 128 * 1024 * 1024 / 1KB = 128K events
* max number of max sized events in 128MB = 128 * 1024 * 1024 / 65550 = 2K events

if there is no fsync, then writes can be re-ordered:
* redundancy inserted over a large time range should cover this

multiple processes appending??????
* one goroutine per open file with a fixed size buffered channel for incoming messages
* pushes back to the tcp listener that blocks
* pushes back to the udp listener that drops
* one receive channel for all incoming events
* this goroutine has an in-memory map of channels to open file processors and it sorts to channel????


write:
[]ts = lp32(ts) and split into a slice of 3 character strings
cd bucket_dir
for i, c3 in []ts:
  if cd c3 works:
    next
  else
    stat c3
    if file
      append_or_split(i, c3, stat, []ts)
    if dir
      next

append_or_split:
  if file > max_file_size
    mv(file, file.p), mkdir(file)
      any other writing process will then write inside the directory and create files as needed
    open file.p and read every message
      call the write process given the current starting offset

read(event_id)
for i, c3 in []ts:
  if cd c3 works:
    next
  else
    stat c3
    if file
      scan_file_for_id(event_id)
    if dir
      race condition between previous cd and now... no worries
      cd c3
      next
    if not there
      return false

make_new_files:
* 2_byte_size,serialized_event(includes time, crc, etc...)

scan_file_for_id:
* 2_byte offset, 8 byte time, 4 byte crc = 14 byte fixed sized buffer
* read buffer. compare time. if match compare crc. if match. make size buffer and read message. prepend time and crc
* if not match, seek forward size bytes and repeat
* return false on EOF


== parts to make =



dispatcher:            writer:
  map:
    writer_interface:
      events->         ->events
      cleanup->        ->cleanup
     switch:
      split<-          <-split
      finished<-       <-finished

event file writer goroutine:
* takes a buffered channel of []byte
* takes a finished channel of bool
* takes a cleanup channel of bool
* takes a split channel of ????
* seek to EOF
* read []byte
* full_write_buffer = len([]byte) + []byte
* write full_write_buffer
* when the channel is empty, sleep for 1 second, fsync, send true on finished channel, and exit
* when cleanup channel is true, fsync and exit
* stat file_size on open, track count of bytes written
* when file_size > max, split:
  mv(file, file.p) && mkdir file

dispatch goroutine:
* takes 1 buffered channel of []byte incoming messages off the wire
* keeps map[try] -> writer_interface_object holds channels
* dispatch new message:
  make a dirs = []try{full/path/to/time/ss, time/ss/sss, time/ss/sss/sss, etc... }
  for try := range dirs {
    channel = map[try]
    if channel
      send to channel
      break
    else
      stat full_path_to_try
      if directory
        continue
      if file or no such file:
        open writer_interface and save entry to map
        send to channel
        break
      if permission denied or disk full etc...
        log/escalate alert
        drop message


Split process:
* writer:
  sees file reach max size
  fsync and close file
  mv(file,file.p) && mkdir file
  send true on split channel
  open file.p and finish writing
  fsync and close file.p
  run split_file_processor
  exit
* dispatch:
  sees true on split channel
  removes entry from map


split_file_processor:
* takes path to file.p static file that will no longer be written to
* takes a channel to send messages back to dispatch.
* read each event. use a small, fixed size lru cache to deduplicate some of the messages
* write intermediate bucket files that match the next level down time/ss/sss/sss etc...
* delete file.p
* read each file. use a small, fixed size lru cache to deduplicate some more of the messages
* send the (slightly more) unique messages back to dispatch
* delete intermediate buckets


https://github.com/golang/groupcache/blob/master/lru/lru.go
*/

/*
store files:
no header.
2 byte message length, message


store.Append(b []byte)
* offset = os.Seek(0, 2) // EOF
* os.Write(len(b))
* os.Write([]byte)
* return offset // new id for the event if needed

store.Read(id int64)
* os.Seek(id, 0) // seek to the offset given by id
* l := make([]byte, 2)
* os.Read(l)
* event_length := int(l)
* e := make([]byte, event_length)
* os.Seek(2,1)
* os.Read(e)
* return e

read the whole thing out into a stream of events????
read l, seek 2, read len(l), emit event
repeat


repair corrupt file?
 events are probably serialized or compressed or both.
 use the attributes of those protocols to seek for the first full valid event, then continue processing

*/
