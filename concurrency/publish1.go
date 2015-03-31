package main

import (
	"fmt"
	"time"
)

// http://www.nada.kth.se/~snilsson/concurrency/

func Publish(text string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		fmt.Println("BREAKING NEWS:", text)
	}()
}

func main() {
	Publish("A goroutine starts a new thread of execution.", 2*time.Second)
	fmt.Println("hello")
	time.Sleep(3 * time.Second)
	fmt.Println("there")
}
