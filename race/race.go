package main

import (
	"fmt"
)

// generate and test some race conditions ???
// http://blog.golang.org/race-detector
// https://code.google.com/p/thread-sanitizer/wiki/Algorithm

// go [test|run|build|install] --race ????
// requires realistic workloads

// shared demonstrates a race condition where multiple goroutines write to the same shared map.
func shared() {
	done := make(chan bool)

	// a map that will be written to from a goroutine that closes over it
	shared := make(map[string]string)
	shared["name"] = "world"

	// write concurrently to the shared map
	go func() {
		shared["name"] = "data race"
		done <- true
	}()

	// read from the shared map
	fmt.Println("Hello,", shared["name"])
	<-done
}

func main() {
	shared()
}
