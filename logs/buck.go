package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
Buck: distributed, partial, bucket sort
* becomes kafkaesq as the query interval approaches zero

AddEvent([]byte):
* read: encoding, replication, priority, time_accuracy, point, crc, length
* pick a random disk and write the message to it.

Query(start_time, end_time):
* for each disk, send the oldest data first
* XXX receiving routine needs to work on the same time period from all channels!!!
   so N disks, there are N incoming channels
   read one event from each channel
     sort the N received events
     send ONLY the oldest event through to the TCP channel
     read another event from THAT channel
     when a channel is empty, close it or stop looping over it

Implement AddEvent and Query as TCP socket listeners
*/

/*
methods required by Query:

QueryDisk(path, start_time, end_time, output_channel)
* find end_time, then recurse directories in time order and ReadFile(path) onto 1 output channel

ReadFile(path)
* read the whole file into [][]byte
* choose start_byte and end_byte based on time digits in "path"
* radix_sort(start_byte, end_byte, [][]byte)
  ram used = 2 * disk_count * file_size
  sort iterations = end_byte - start_byte


methods required by AddEvent:

ReadEventHead

*/

/*
Data Scale Estimation:

event_rate_target:
* average events from all servers in one second?
* currently 6M events/hour/75servers ~ 22 events/second/server
* target rate is 500 events/second/server * 1K servers = 500K events/second
* 500K events/second = 500 events/ms = 0.5 events/us
* 500K * 512 byte events/second = 244MB/second = 20TB/day

bucket_event_count and storage size:
* current 2K events / second * 1024 seconds ~ 2M events/bucket
* target 500K events / second * 1024 seconds ~ 512M events/bucket
* and around 244MB/second * 1024 seconds ~ 244GB/bucket means 2 buckets per ec2 disk
* should be possible to handle the indexing for 512M event id's per bucket... right???
* assuming massive compression is possible...

event_key:
* how many bits of "point" are needed inside directory?
  32^2 = 1024 seconds = 2^10 = the lowest 10 bits
  fractional time = number of 10^-8 second (10ns) intervals since the second began.
  10^8, 10ns intervals = 1 second
  1 second or 10^8 10ns intervals fit in 27 bits
  so... 10 bits for seconds + 27 bits for fractional seconds = 37 bits for "point" inside bucket

* probability of collision using time stamp alone?
  this is less events than the granularity of Event.point = 10ns intervals
  calculate the probability of 2 messages at the same time????
    NTP cloud time accuracy ~ 10^-3 or 1ms
     every 1 ms can contain 10^5 * 10^-8 intervals or ~ 17 bits worth
    given our 37 bit "bucket point":
      20 bits can be considered collision free, but 17 bits must be considered like a random hash

* so how much extra entropy is required to avoid message collision in a "bucket"??
  total number of messages in 1ms = 5*10^5 events/second * 10^-3 seconds/ms = 500 events/ms
  bits of entropy required to avoid collision in 500 events:
    ((500)^3)/2 fits in ~ 26 bits
  we've got 17 bits already and need a total of 26. find 9 bits somewhere or risk it?
  XXXX I'm considering dropping the granularity of the timestamp in favor of checksum
    how many bits of timestamp in bucket for ms granularity?
      1 second or 10^3 ms fits in 2^10 or 10 bits
      so... 10bits for seconds + 10 bits for fractional seconds + 26 bits entropy = 46 bits required
      how to divide up the 6 bytes?  3 time + 3 checksum should be ok

Summary: 3 bytes of time + 3 bytes of checksum is the minimum required to avoid collision

*/

/*
Sorting Incoming Events by "point" using all available disks equally
* do a partial sort across N disks in ./disks/N/
* bucket interval is NOT determined at write time.
  consumers may request any interval.
  only certain intervals will be available.
  a completeness probability can be assigned to each interval as of the query time.

 and do the sorts separately on each disk.
 - keep the full time directory tree on each disk, but put less files in it!
   this makes less splits on each disk! ... and probably warrants more pre-allocation of directories?

need a way to map ./time/ss... onto ./disk1/timess...
  or some other way to spread disk access across more than one disk


  XXXX 90% of incoming goes to the same interval_id...
   right????? this is the basis of the VFS cache optimization
   so why try to spread it across disks???
    because it speeds up the disk io when writes are required!!!


calculate file size based on radix_sort memory requirements
Assume there are multiple incoming sorts and multiple outgoing sorts
* m1.medium = 3.7, 1x410: 2 * 1 * 64 = 128: 15 simultaenous sorts in 1/2 of available ram
* m1.large  = 7.5, 2x420: 2 * 2 * 64 = 256: 15 " "
* m1.xlarge = 15,  4x420: 2 * 4 * 64 = 512: 15 " "



Deduplicate Events on READ, NOT WRITE:
* accept the 3x event count in intermediate store files
* stip high order bits from key on write (high order bits are stored in the file name)
* split on fixed file size
* sort and remove duplicates in memory on split when reading the keys into memory
* ensure that in-memory sort structure takes up less than N x 4096 blocks in the fs cache


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

the token frequency analysis purposes:
* identify bucket/time/ss containing the token for retrieval from long term storage
* sort tokens for decompression
* I don't think either of these is a huge problem if it's 1/3 duplication!!!!!!!!!!!!!
* so only try to deduplicate if it's cheap or necessary (like at freeze time when events can be statically sorted.)



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
back to the same dedup problem as before, but on a fixed size data set!!!
* append-only on-disk tree structure
* fast append
* fast key lookup
* ideally less than 1% of raw event storage 64MB to 292MB => 655k to 3MB of index

http://guide.couchdb.org/draft/btree.html
* append-only b+tree
* the root node must be re-written every time the file is updated
* old parts are never overwritten... so every old root is a consistent "snapshot" of the db???
* append new leaf node
* read parent???
* append new parent to reference new leaf node
* continue appending new parents all the way through to the root node
* "commit" is when the root gets overwritten
* read: root node is the end of the file. traverse tree pointers.
* is this a huge amout of wasted space or no???
* how to locate the last valid root node????
  write root node with length and some magic ID byte ????
  if the last node in the file a corrupted root node or not a root node,
    then seek backward one node at a time till the last valid root node is found.

could have fixed width nodes with fixed width id's and overwrite in place... won't handle corruption well


if there is no fsync, then writes can be re-ordered!
write barriers????

http://sphia.org/architecture.html
* LSM is a collection of sorted files that are periodically merged
* region in-memory index with a in-memory key index
* no internal page to page links???



leveldb:
* http://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/
* https://code.google.com/p/leveldb-go/

bitcask:
* https://code.google.com/p/gocask/
* http://downloads.basho.com/papers/bitcask-intro.pdf
* crc,ts,key_size,value_size,key,value
* key -> file_id,value_size,value_offset,tstamp
http://godoc.org/github.com/cznic/kv

http://godoc.org/github.com/cznic/exp/dbm

more about the vfs cache: https://www.kernel.org/doc/Documentation/filesystems/vfs.txt

http://en.wikipedia.org/wiki/B-tree

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

/*
choose specialized data structures for the following purposes:
* 3 - 64 character token -> frequency int32
* 3gram -> []token
* global token -> []bucket
* global 3gram -> []bucket




*/

/*

buckets/tim/es:
* 32768, 12.1 day directories / 1024, 17.1 minute buckets

token frequency data????
* compact, searchable, balanced, prefix tree or normal hash table
* token -> count

3gram.kv:
* key = 3gram
* value = []token

Data updated during "freeze" D days after tim/es

tokens.b:
* listB(6,8)
* tokens sorted by frequency
* created during freeze

static.kv:
* key = event_key
* value = events with "tokens" removed and "line" tokenized using static_token.list

global/all_3gram.kv
* key = 3gram
* value = []tim/es

global/all_tokens.kv
* key = token_string
* value = []tim/es


*/

/*
Buck Processes:


*/

/*
Buck Data:

receive "Events" from clients and peers


build a multi-level index per bucket:
* bloom filter to deduplicate incoming events
* token frequecy token[hello] -> count
* 3gram[h_e] -> []token

build a global index:
* token[hello] -> []bucket
* 3gram[l_o] -> []token

freeze the bucket for long term storage:
* frequency sorted token list
* gzip events in key order with tokens replaced


and respond to requests for bulk event movement:
* bucket(time_range) -> []bucket_path
* bucket([]3gram, time_range) -> []bucket_path
* bucket([]token, time_range) -> []bucket_path
* events(bucket_path) -> []event

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
