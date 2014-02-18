package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
Data Scale Estimation:

event_rate_target:
* average events from all servers in one second?
* currently 6M events/hour/75servers ~ 22 events/second/server
* target rate is 500 events/second/server * 1K servers = 500K events/second
* 500K events/second = 500 events/ms = 0.5 events/us
* 500K * 512 byte events/second = 244MB/second = 20TB/day

bucket_event_count and storage size:
* 500K events / second * 1024 seconds ~ 512M events/bucket
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
filter files:
* fixed sized bit arrays that answer "Seen this event?" with "Exactly Never" or "Maybe"
* basically the opposite of a bloom filter giving a probabilistic result
* efficiently deduplicate incoming events
* initially created as sparse files, but should fill up fast.
* easily rebuilt from accumulated events if lost or corrupted
* discarded after "freeze"

filter.New(key_length int, path)
* makes file
* writes fixed width hader field with key_length
* other header fields????

filter.Set(key int64) error
* 0 < key < (1<<40) // never make a 1TB bit array.
* offset = len(header) + key
* file.Seek(offset, 0)
* b := []byte{0}
* file.Read(b) // read a whole byte, but only modify 1 bit.
* b = b[0] | 128
* file.Write(b)
* return error on disk error
1 seek, 1 read, 1 write

filter.Get(key int64) bool
* validate key
* offset = len(header) + key
* b := []byte{0}
* ReadAt(b, offset)
* if b[0] & 128 == 1 return true
* return false
1 seek 1 read

func (f *filter) GetSet(key int64) (found bool)
* validate key
* offset = len(header) + key
* b := []byte{0}
* file.Seek(offset, 0)
* file.Read(b) // read a whole byte, but only modify 1 bit.
* if b[0] & 128 == 1 return true
* b = b[0] | 128 // set the 1 bit
* file.Write(b)


* ReadAt(b, offset)
* if b[0] & 128 == 1 return true
* return false


ok... try again:
* 2MB 3,3 bloom filter needs 48 bits of input
* give it: 16 bits of time + 32 bits of crc
* and only cover 500K events
* ( 1 - 2.71828^(-2 * (500000) / (2^24)) )^2 = 0.33% false positives
* 2MB per second = 2GB of filter per bucket

512MB 4,4 bloom filter needs 64 bits of input
* give it 32 time + 32 crc
* and cover 5*10^8 events/bucket
* ( 1 - 2.71828^(-2 * (5*10^8) / (2^32)) )^2 = 4% false positive rate
* and it drops off fast 0.002% at 10^7


ok... try again:

how about ditching the event key and combining deduplication with the token store in a pre-fab kv store?
is the event a duplicate?
* look up all the unique tokens in the store. return all the docs they point to. find the intersection.
* that sounds really slow and doesn't account for the same message being permitted over time.

tokens[token]->event_key







*/

/*
store files:
no header.
2 byte message length, message


store.Write(b []byte)
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
Bucket Data Structures



2^16 * 4 = 256K


events.store: 244GB

accumulate kv of tokens -> []something???  all 6 bytes of event_id???? checksum only???
* non-lossy normal on disk hash table

6 bytes * 5*10^8 = 2.7GB token database per 17min











* index should test existence in ~ 2 disk seeks
* events should be expected to arrive out of order
* key = event_key
* value = event bytes directly off wire
* file operations:
  ReadAt(b []byte, off int64) = pread(2)
  WriteAt(b []byte, off int64) = pwrite(2)
  Truncate(size int64) = syscall.Ftruncate ?????
  more about the vfs cache: https://www.kernel.org/doc/Documentation/filesystems/vfs.txt
* sparse event_id list would be 2^48 bits or 32 TB. larger than FS file size limit.


thought experiment for rough tree structure:
first 16 blocks are first layer index:
256x256 matrix of pointers to all combinations of the first 2 bytes.

read:
* seek to the first 2 bytes i1. read 3 bytes i2.
* seek from beginning to i2. read 6 bytes i3.
* if i3 == event_id, then it's alreay here

????????????????????

http://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/

SSTable and LevelDB and Log Structured Merge Trees

Leveldb would be redoing compression and checksum work that we've already done!
 - re-implement using these pre-computed values as input







extendible hashing:
 - when the data gets over filled, double the index size and add more buckets

linear hashing:
 - split policy or "load factor"
 - state: i = current round of splitting
 - state: p = the next bucket to split



4GB bloom filter for every 244GB ?

4GB bloom filter as 4,4
m = 2^32 bits in the bloom filter array
k = 2 hash functions
n = 5*10^8 events

( 1 - 2.71828^(-2 * (5*10^8) / (2^32)) )^2 = 4% chance of false positive is still too high


http://stackoverflow.com/questions/635728/opposite-of-bloom-filter
http://www.somethingsimilar.com/2012/05/21/the-opposite-of-a-bloom-filter/
http://www.i-programmer.info/programming/theory/4641-the-invertible-bloom-filter.html
* lossy hash table or LRU cache ???
* https://github.com/jmhodges/opposite_of_a_bloom_filter

how big is a prefix tree of all possible 9 byte id's ???

break messages up into even sized blocks for fixed offset file storage

sparse bit array: 2^(9*8) bits long. seek to the 9byte address and set 1 or 0

* bisection or binary search of a linear list. fibonaccian search: same, but no division

binary tree of fixed width keys:
key,left,right    where left and right are other keys

reading:
* read root key and decide next_key = left or right
* seek to and read next_key and repeat till ????

prefix-tree is also a good fit for autocompleting dictionary

bloom_filter_false_positive_rate = ( 1 - 2.71828^(-kn/m) )^k


buckets/tim/es:
* 32768, 12.1 day directories / 1024, 17.1 minute buckets



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




* bucket messages onto local ephemeral filesystem then compress and archive to s3

*/

/*
Buck:

receive Tailer "Events"

build a multi-level index:
3gram(h_e) -> []token
token(hello) -> []event_key
event(event_key) -> event

and respond to requests for bulk event movement:
events(start_time, end_time) -> []bucket_path
*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
