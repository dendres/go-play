package main

import (
	"log"

	"github.com/ActiveState/tail"
)

// check on the test cases https://github.com/ActiveState/tail/blob/master/tail_test.go
// and maybe write some expicit tests to compare to tail -F behavior!!!

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	file := "/tmp/tailing.txt"

	t, err := tail.TailFile(file, tail.Config{ReOpen: true, Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		log.Println(line.Text)
		// parse the lines.. see parsing.go
	}

	// log.Println("waiting 90 seconds")
	// <-time.After(90 * time.Second)

}
