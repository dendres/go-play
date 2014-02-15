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

event_key: 16 bytes
* how many bits of "point" are needed inside directory?
  32^2 = 1024 seconds = the lowest 10 bits
  fractional time = number of 10^-8 intervals since the second began.
  10^8, 10ns intervals = 1 second
  1 second or 10^8 10ns intervals fit in 27 bits
  need 10 bits for 1024 seconds
  so 10 second bits + 27 fractional bits = 37 bits
  which is contained by... 5 bytes.
* 4 bytes of crc32
* 3 bytes of shn[1:3]
* 4 bytes of app[0:3]
* 5 + 4 + 3 + 4 = 16 byte

message_count_max:
* how many unique words from all servers in 17.1 minutes?
* currently 6M events/hour/75servers ~ 22 events/second/server
* target rate is 1K events/second/server * 10K servers = 10M events/second
* indexes should not prevent ~ 10B events per bucket

token_count_max:
* 10B events... how many words in 17min from 1k servers ???
* worst case every message has a new sha256 and we save the whole thing
* 6 byte count allows 2^48 tokens / 256 bytes/token ~ 1T events

buckets/tim/es:
* 32768, 12.1 day directories / 1024, 17.1 minute buckets


Data updated for every incoming message:

events.kv:
* key = event_key
* value = event bytes directly off wire

tokens.kv:
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
