package main

import (
	"fmt"
	"math/rand"
	"time"
)

// continuing to work through http://blog.golang.org/race-detector

// randomDuration returns a random time.Duration between 0 and 1 second
func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}

// racyTimer demonstrates variable access from the main function body and from a goroutine called from the function.
// time.AfterFunc runs the function argument in a goroutine
func racyTimer() {
	start := time.Now()

	// XXX have to initialize timer here or it shows up as undefined in the function below
	var timer *time.Timer

	// assignment to timer in the function body
	timer = time.AfterFunc(randomDuration(), func() {
		diff := time.Now().Sub(start)

		fmt.Println(diff)

		// read from timer from the goroutine
		timer.Reset(randomDuration())
	})

	time.Sleep(5 * time.Second)
}

// lessRacyTimer resolves the racyTimer problem by not accessing timer from the goroutine
func lessRacyTimer() {
	start := time.Now()
	reset := make(chan bool)
	timer := time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		reset <- true
	})

	for time.Since(start) < 5*time.Second {
		<-reset
		timer.Reset(randomDuration())
	}
}

// r3 uses separate timers instead of resetting the timer
// this is simpler, but apparently less efficient????
func r3() {
	start := time.Now()

	time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		time.AfterFunc(randomDuration(), f)
	})

	time.Sleep(5 * time.Second)
}

func main() {
	fmt.Println("running racyTimer")
	racyTimer()

	fmt.Println("running lessRacyTimer")
	lessRacyTimer()

	fmt.Println("running r3")
	r3()
}
