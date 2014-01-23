package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	file, err := os.Open("fake_file") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(file)
}
