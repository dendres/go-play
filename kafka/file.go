package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
file harvester:

watch / tail a file and optionally manage it's rotation

journal/deduplicate lines from the file:
 - message_id ??  increment, timestamp, source, checksum, etc...
 - read a file backlog from a configurable starting point (lines_back, duration_back)
 - once the journal is established, align the journal with the log file to find the correct starting point


kafka topic and partition????


use the existing file as the buffer: and parse timestamps multiple times?????

parse lines and buffer serialized data to disk??? (double disk I/O)






*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
