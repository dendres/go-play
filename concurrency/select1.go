package main

import "fmt"

func RandomBits() <-chan int {
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()

	return ch
}

func main() {
	r := RandomBits()

	// prints forever
	for i := range r {
		fmt.Println(i)
	}
}
