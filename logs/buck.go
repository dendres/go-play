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

message_count_max:
* how many unique events total from all servers in 17.1 minutes?
* currently 6M events/hour/75servers ~ 22 events/second/server
* target rate is 1K events/second/server * 10K servers = 10M events/second
* 10M events / second * 1024 seconds (17.1 minutes) ~ 10B events
* all indexing and hashing methods MUST allow 10B events per bucket
* 10M events/second = 10^7 events every 10^9 ns = 1 event every 100ns

event_key:
* how many bits of "point" are needed inside directory?
  32^2 = 1024 seconds = the lowest 10 bits
  fractional time = number of 10^-8 second (10ns) intervals since the second began.
  10^8, 10ns intervals = 1 second
  1 second or 10^8 10ns intervals fit in 27 bits
  need 10 bits for 1024 seconds
  so 10 bits for seconds + 27 bits for fractional seconds = 37 bits for "point"
  which gets rounded up to 40 bits or 5 bytes.

* probability of collision using time stamp alone?
  this is less messages than the granularity of Event.point = 10ns intervals
  calculate the probability of 2 messages at the same time????
    NTP cloud time accuracy ~ 10^-3 or 1ms
     every 1 ms can contain 10^5 * 10^-8 intervals or ~ 17 bits worth
    given our 37 bit "bucket point":
      20 bits can be considered collision free, but 17 bits must be considered like a random hash

* so how much extra entropy is required to avoid message collision in a "bucket"??
  total number of messages in 1ms = 10^7 events/second * 10^-3 seconds/ms = 10^4 events/ms
  bits of entropy required to avoid collision in 10^4 events => max number of ((10^4)^3)/2 fits in ~ 39 bits

* we've got 17 bits already and need a total of 39, so a 22 bit hash function would be perfect
  a crc32 is overkill by 10 bits, so 3 or 4 bytes of crc should be ok
  rounding up to the full crc32 for now for the added benefit of having it available

* the bucket event key will be 5 + 4 = 9 bytes long


token_count_max:
* 10B events... how many words in 17min from 1k servers ???
* worst case every message has a new sha256 and we save the whole thing
* 6 byte count allows 2^48 tokens / 256 bytes/token ~ 1T events

buckets/tim/es:
* 32768, 12.1 day directories / 1024, 17.1 minute buckets


Data updated for every incoming event:

events.kv:
* deduplicates incoming events
* key = event_key
* value = event bytes directly off wire

need a better on-disk deduplicating data structure in a single file??????
* file operations:
  ReadAt(b []byte, off int64) = pread(2)
  WriteAt(b []byte, off int64) = pwrite(2)
  Truncate(size int64) = syscall.Ftruncate ?????
* more about the vfs cache: https://www.kernel.org/doc/Documentation/filesystems/vfs.txt

* not sure how sparse files get treated in vfs cache... don't want a bunch of zeros in memory!!!!
  probably avoid sparse files for now?????

* O(1) key lookup????
* append the key and value to the end of the file if it does not exist??????

* could do sparse file where key = offset????




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
