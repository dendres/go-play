package main

import (
	"fmt"
	"time"
)

/*
http://www.nada.kth.se/~snilsson/concurrency/

chan Sushi      // can be used to send and receive values of type Sushi
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
*/

type Sushi string

func Producer() <-chan Sushi {
	ch := make(chan Sushi)
	go func() {
		ch <- Sushi("hi")
		time.Sleep(1 * time.Second)
		ch <- Sushi("there")
		close(ch)
	}()
	return ch
}

func main() {
	var ch <-chan Sushi = Producer()

	// read values from the channel
	for s := range ch {
		fmt.Println("Consumed", s)
	}
}
