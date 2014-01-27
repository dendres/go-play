package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
log messasge consumers:
 - get(topic, start_date, end_date) returns multiple compressed archives of log messages
   optionally inflate, parse, encode as desired, etc...





*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
