package main

import (
	"fmt"
	"time"
)

/*
http://dave.cheney.net/2013/04/30/curious-channels

chan Sushi      // can be used to send and receive values of type Sushi
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
*/

/*
a nil chan always blocks!!
  not initialized... or explicitly set to nil

var ch chan bool // note: not the same as make(chan bool)
ch <- true       // deadlock
<-ch             // also deadlock

*/

// wait till both channels are closed
func WaitMany(a, b chan bool) {
	for a != nil || b != nil {
		fmt.Println("a =", a, "b =", b)
		select {
		case <-a:
			fmt.Println("setting a = nil")
			a = nil
		case <-b:
			fmt.Println("setting b = nil")
			b = nil
		}
	}
}

func main() {
	a, b := make(chan bool), make(chan bool)
	t0 := time.Now()
	go func() {
		close(a)
		close(b)
	}()

	WaitMany(a, b)
	fmt.Printf("waited %v for WaitMany\n", time.Since(t0))
}
