package main

import (
	"fmt"
	"sync"
)

/*
http://www.nada.kth.se/~snilsson/concurrency/

sometimes find race conditions:
  go run --race race1.go

*/

// race condition the counter runs through to the end
// before the goroutines start
func race1() {
	var wg sync.WaitGroup
	routines := 9

	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}
	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("finished")
}

// pass the counter as an argument to each goroutine
func fix1() {
	var wg sync.WaitGroup
	routines := 9

	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func(n int) {
			fmt.Println(n)
			wg.Done()
		}(i)
	}
	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("finished")
}

// make a separate local variable for each goroutine
func fix2() {
	var wg sync.WaitGroup
	routines := 9

	wg.Add(routines)
	for i := 0; i < routines; i++ {
		n := i
		go func() {
			fmt.Println(n)
			wg.Done()
		}()
	}
	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("finished")
}

func main() {
	// race1()
	fix1()
	fix2()
}
