package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
single kafka partition reading goroutine.
single consumer_group for all partitions (queue balance messages across consumers)
buck.Listen(topic, partition_number)
* determine if data is missing and the offset needs a rewind ??????
* decompress and separate kafkamessages into events
* for each event, parseonly the timestamp!
* if-e processed/<topic>/<5mstamp>, write the message to recovered/<topic>/etc...
* write compressed event to ephemeral disk: buckets/<topic>/<5mstamp>/<us_stamp>/<small_checksum>
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
buck:

* bucket messages onto local ephemeral filesystem then compress and archive to s3
* a SPOF server that can be offline for up to 2 days with minimal impact
* after 48 hours process the oldest buckets
* sort, bz2, send to s3

configure topic = t-<environment_name>
configure buckets_folder = /data/buck/buckets
configure recovered_folder = /data/buck/recovered
configure upload_goroutine_count = 3

consume ALL 15 partitions. one goroutine each

5mstamp: "find" and "ls" must sort time correctly for this sortable base36 timestamp representing a 5min boundary
*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
