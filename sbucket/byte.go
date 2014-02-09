package main

import (
	"fmt"
)

const base32map string = "0123456789abcdefghijklmnopqrstuvwxyz"

func main() {
	sec := uint64(1555) // 1, 16, 19

	x := make([]byte, 8)

	for i := 7; i >= 0; i-- {
		x[i] = base32map[sec%32]
		sec = sec >> 5
	}

	fmt.Println("byte array =", x)
	fmt.Println("string =", string(x))
}
