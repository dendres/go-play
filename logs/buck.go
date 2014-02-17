package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
goroutine to handle incoming messages from tailer
*/
func Listen() {
}

/*
upload goroutine:
buck.Upload(5mbucket_path, processed_path)
* list buckets older than 48 hours
* pick the oldest bucket
* start over if processed_file exists and is younger than 2 hours
* XXX possibly use contents of processed_file to avoid looping on problems??? like mark error in processed_file???
* initialize bz2_writer(bz2_file_path)
* recurse bucket directory in sort order without having to read all filenames into memory?????
* for each file, read the file to bz2_writer
* close bz2_writer
* s3_sender(path_to_bz_file)
* collect a bunch of metadata and write it to the processed file.
* delete bucket

*/
func Upload() {
}

/*

s3 sender????
How to write to S3????

* s3 object can tell you:
   content_type, content_length, etag??, exists, expiration_date???, last_modified, metadata??, versions???
s3 sender????
How to write to S3????

* s3 object can tell you:
   content_type, content_length, etag??, exists, expiration_date???, last_modified, metadata??, versions???

*/
func s3_send() {
}

/*
Buck:

receive Tailer "Events"

build a multi-level index:
3gram(h_e) -> []token
token(hello) -> []event_key
event(event_key) -> event

and respond to requests for bulk event movement:
events(start_time, end_time) -> []bucket_path

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

* roughing out indexing with a 6 byte event key.


Data updated for every incoming event:

events.kv:
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









4GB bloom filter for every 244GB

4GB bloom filter as 4,4
m = 2^32 bits in the bloom filter array
k = 2 hash functions
n = 5*10^8 events

( 1 - 2.71828^(-2 * (5*10^8) / (2^32)) )^2 = 4% chance of false positive is still too high



maybe take the 9 bytes, base32 encode them, abc/def/
15 characters







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

tokens.kv:
* deduplicates incoming events
* records event frequency
* key = token_string
* value = []event_key

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
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
