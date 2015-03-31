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
struct{}        // empty data means all channel data will/should be ignored

*/

func Publish(text string, delay time.Duration) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		time.Sleep(delay)
		fmt.Println("BREAKING NEWS:", text)
		close(ch)
	}()
	return ch
}

func main() {
	wait := Publish("something", 1*time.Second)
	fmt.Println("waiting...")
	<-wait // blocks execution until the breaking news has been printed
	fmt.Println("got something")
}
