package main

import (
	"log"
	"os"
	"runtime/pprof"
)

func IsMultiple(x int) bool {
	if x < 1 {
		return false
	}

	if (x%3) == 0 || (x%5) == 0 {
		return true
	}

	return false
}

func SumMultiplesBelow(x int) int {
	total := 0
	for i := 1; i < x; i++ {
		if IsMultiple(i) {
			total += i
		}
	}
	return total
}

func main() {
	// export CPUPROFILE_FREQUENCY=100000

	prof, err := os.Create("p.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer prof.Close()

	pprof.StartCPUProfile(prof)
	defer pprof.StopCPUProfile()

	for i := 0; i < 1000000; i++ {
		out := SumMultiplesBelow(10)
		if out != 23 {
			log.Fatal("the sum of the multiples of 3 or 5 less than 10 should be 23")
		}
	}

	log.Println("finished ok")
}
