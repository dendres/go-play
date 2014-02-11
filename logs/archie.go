package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
goroutine initated by REST/JSON api
archie.Request(environment_name, start_time, end_time)
* generates a tmp elasticsearch index name based on the input paramaters
* makes a list of 5min time blocks to retrieve
* adds to a disk based queue of the 5min blocks
* returns to caller: elasticsearch index name, number of blocks that need retrieved, time estimate???
*/
func Request() {
}

/*
goroutine to read disk based queue, download from s3, and bulk upload to elasticsearch
configurable pool size
archie.Index()
* find next job on queue to work on
* lock job / mark as in progress
* check for downloaded file
   if missing, download
   if there, get checksum from s3
   if checksum is wrong, delete and download again
* extract messages, convert to json / bulk upload format and write all bulk uploads to disk
* form bulk upload / uploads ... not sure how big they should be ?????
* send all data to elasticsearch
* clean up tmp files




XXXX figure out how to do this as a stream???

get some blocks off the s3 download
uncompress the blocks
when a full message is available, decode and write to current es bulk upload
when es bulk upload is ready, send to es



*/
func Index() {
}

/*
archie:

* generate elasticsearch indexes on demand from archive
* takes (environment_name, start_time, end_time)
* generates a tmp index name based on the input paramaters and return it to the caller
* makes a list of 5min time blocks to retrieve
* downloads blocks from s3
* form elasticsearch bulk import messages

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
