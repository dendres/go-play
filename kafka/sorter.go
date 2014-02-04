package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/* sorter:
* a single-instance non-redundant service that can be offline for 2 days with minimal impact
* consume ALL partitions
* if-e processed/<topic>/<5mstamp>, write the message to recovered/<topic>/etc...
* bucket messages on ephemeral disk: buckets/<topic>/<5mstamp>/<us_stamp>/<small_checksum>
* after 48 hours process the oldest buckets
* sort, bz2, send to s3
* mark processed by writing some metadata to processed/<topic>/<5mstamp>
* delete bucket


Definitions:
* 5mstamp: "find" and "ls" must sort time correctly for this sortable base36 timestamp representing a 5min boundary



sorter.Listen(topic, partition_number)
* kafka client for a single partition
* decompress and separate kafkamessages into events
* for each event, parseonly the timestamp!
* copy serialized event to disk.
XXXXX recover and replay mechanism???????
  find out what offset to start at??????
XXXX disk not mounted recovery????
consumer group: each message published to a topic is delivered to a single consumer instance in the consumer group

multiple instances in a consumer group = queue balanced across instances
one instance per consumer group = broadcast to all instances




sorter.Process(5mbucket_path, processed_path)
* find buckets older than 48 hours
* touch the processed file (empty processed file means work in progress.. no more new messages... you had 48 hours)
* initialize bz2_writer(bz2_file_path)
* recurse bucket directory in sort order without having to read all filenames into memory?????
* for each file, read the file to bz2_writer
* close bz2_writer
* s3_sender(path_to_bz_file)
* collect a bunch of metadata and write it to the processed file.

s3 sender????
How to write to S3????

* s3 object can tell you:
   content_type, content_length, etag??, exists, expiration_date???, last_modified, metadata??, versions???




*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
