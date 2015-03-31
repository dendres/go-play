package main

import (
	"fmt"
	"sync"
)

func Seek(name string, match chan string, wg *sync.WaitGroup) {
	select {
	case peer := <-match:
		fmt.Println(peer, "sent a message to", name)
	case match <- name:
		// wait for someone to receive my message????
	}
	wg.Done()
}

// match up pairs of users on the "match" channel
func main() {
	people := []string{"Alice", "Bob", "Charlie", "David", "Eva"}
	match := make(chan string, 1) // buffered... allows for 1 unmatched send???
	wg := new(sync.WaitGroup)

	wg.Add(len(people))

	for _, name := range people {
		go Seek(name, match, wg)
	}

	wg.Wait()

	select {
	case name := <-match:
		fmt.Println("no one received the message from", name)
	default:
		// no pending send operation??
	}

}

// XXX spend more time looking at send and receive on the same channel!!!!
