package main

import (
	"fmt"
)

func main() {
	input := []int{5, 6, 2, 7, 8, 1, 4, 9, 0}
	fmt.Println(input)

	var i, j, a, b, N int

	n := len(input)

	for {
		N = 0
		for j = 1; j < n; j++ {
			i = j - 1
			a = input[i]
			b = input[j]

			if b < a {
				input[i], input[j] = input[j], input[i]
				N = j
			}
		}
		if N == 1 {
			break
		}
		n = N
		fmt.Println(input)
	}
	fmt.Println(input)
}

// I don't really get the optimization about changing n

// I don't see a way to make this concurrent
