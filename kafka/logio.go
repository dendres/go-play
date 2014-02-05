package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
port the log.io server
replace tcp connection with kafka consumer

https://github.com/NarrativeScience/Log.io/blob/master/src/server.coffee

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
