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
bloom filter files:
* deduplicate incoming events with a fixed 0.001% false positive rate
* easily rebuilt from accumulated events if lost or corrupted
* discarded after "freeze"
* key size is fixed at bucket creation time

adjust key size for each new bucket:
* 24 hours = 84.4 time buckets
* event_estimate = max(last 200 time bucket event_count)
* bk := bloomkeybits(event_estimate)

if bk < 24
 24 bit time + 24 bit crc32 minimum!
if bk < 32
 bk bit time + bk bit crc
else
  37 bit time + 32 bit crc +
  start adding bits from shn, app, and level (the only remaining visible fields)

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


leveldb:
* http://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/
* https://code.google.com/p/leveldb-go/

bitcask:
* https://code.google.com/p/gocask/
* http://downloads.basho.com/papers/bitcask-intro.pdf

http://godoc.org/github.com/cznic/kv

http://godoc.org/github.com/cznic/exp/dbm

more about the vfs cache: https://www.kernel.org/doc/Documentation/filesystems/vfs.txt

*/

/*

buckets/tim/es:
* 32768, 12.1 day directories / 1024, 17.1 minute buckets

incoming.filter
* 128MB - 32GB bloom filter deduplicates messages
* discard on freeze

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
Buck:

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
