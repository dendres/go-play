package main

import (
	"fmt"
	"sync"
	"time"
)

/*
http://dave.cheney.net/2013/04/30/curious-channels

chan Sushi      // can be used to send and receive values of type Sushi
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
*/

func main() {
	const n = 25
	finish := make(chan bool) // make(chan struct{}) says "this channel will only be used for its closed property"
	var done sync.WaitGroup

	for i := 0; i < n; i++ {
		done.Add(1)
		go func() {
			select {
			case <-time.After(1 * time.Hour):
				fmt.Println("waited 1 hour")
			case <-finish:
				// when the channel is closed, this case will trigger and the goroutine will finish
				fmt.Println("received from the finish channel")
			}
			done.Done()
		}()
	}

	t0 := time.Now()

	fmt.Println("closing the finish channel")
	close(finish)

	fmt.Println("waiting on the wait group")
	done.Wait()

	fmt.Printf("Waited %v for %d goroutines to stop\n", time.Since(t0), n)
}
