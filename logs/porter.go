package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
porter:

* for each kafka message, form an elasticsearch bulk import
* cluster of 3 to 5 nodes (maybe runs on kafka servers???)
* manage D (maybe 30) days worth of logstash-style (but smaller) "hour" elasticsearch indexes

"tailer" host sends messages to a random topic and changes topic frequently

configure topic = t-<environment_name>

make list of my partitions:
  static partition_count = 15
  configure cluster_size = 3
  configure consumer_number = 0 to (1 - cluster_size)


for each message from a partition (one goroutine per partition)
  decompress the message
  send elasticsearch bulk import
  then ask for the next message

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
