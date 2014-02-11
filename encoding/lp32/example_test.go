package lp32_test

import (
	"fmt"
	"github.com/dendres/go-play/encoding/lp32"
)

func ExampleEncode() {
	n := uint64(999)
	s, err := lp32.Encode(n, 4)
	if err != nil {
		fmt.Printf("error encoding %d:%v", n, err)
	} else {
		fmt.Println(s)
	}
	// Output:
	// 00v7
}

func ExampleDecode() {
	s := "00v7"
	n, err := lp32.Decode(s)
	if err != nil {
		fmt.Printf("error decoding %s:%v", s, err)
	} else {
		fmt.Println(n)
	}
	// Output:
	// 999
}
