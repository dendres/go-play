package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
interface kafka for the logio system??????

listen for incoming messages

send them to the local log.io system?????

*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
