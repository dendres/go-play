package main

import (
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(1000000000)
	log.Println("started ticker")

	for r := range ticker.C {
		log.Println("tick", r)
	}
}
