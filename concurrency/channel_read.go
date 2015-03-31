package main

import "fmt"

// http://www.nada.kth.se/~snilsson/concurrency/

/*
chan Sushi      // can be used to send and receive values of type Sushi
chan<- float64  // can only be used to send float64s
<-chan int      // can only be used to receive ints
*/

func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello!"
		close(ch)
	}()

	// unary invocation of '<-' reads from the channel:
	fmt.Println(<-ch)

	// the zero value of the channel type is returned if there is no value to read
	// ... when the channel has been closed
	fmt.Println(<-ch)

	// The value of ok is true if the value received was delivered by a successful send operation to the channel,
	// or false if it is a zero value generated because the channel is closed and empty.
	value, ok := <-ch
	if ok {
		fmt.Println(value)
	} else {
		fmt.Println("not ok")
	}
}
