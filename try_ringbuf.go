package main

import (
	//	"fmt"
	"github.com/dendres/go-play/ringbuf"
)

func main() {
	//	out := make([]byte, 5)
	r := ringbuf.New(5)
	for c := 1; c < 10; c++ {
		r.Write([]byte{byte(c)})
	}

	//	r.Read(out)
	//	fmt.Println("out =", string(out))

}
