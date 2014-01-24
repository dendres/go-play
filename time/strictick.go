package main

import (
	"log"
	"time"
)

// actually... this is probably not needed, since the message timestamps will be used

// strictick returns a channel that emits the current time on the provided interval
// sleeps till an interval boundary
func strictick(interval time.Duration, out chan time.Time) {
	// BUG(done) ensure interval divides evenly into Hour or is a multiple of Hour

	for {
		now := time.Now()
		next := now.Round(interval)
		if next.Before(now) {
			next = next.Add(interval)
		}

		time.Sleep(next.Sub(now))
		out <- next
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	c := make(chan time.Time)

	i := time.Millisecond * 500

	go strictick(i, c)

	for t := range c {
		log.Println("tick =", t)
	}
}
