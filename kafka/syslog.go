package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
syslog harvester:

listen on syslog port
parse and multiline messages and determine message identity fields
serialize and bucket by timestamp
compress and send to kafka

allow configurable message durability for: disk io vs. memory footprint needs ???
message journal and deduplication???

crypto message authenticity???

map level and facility to kafka topics???

map messages to kafka partitions???
 - client's responsibility to determine which partition to send which message to??
 - default: hash(key) % numPartitions

rough disk based queue process for a single "queue":
* /queue/<source>/<interval>.txt
* wait "delay", then read file and send as a single kafka message
* delete file after ack


*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
