package main

import (
	"fmt"
	"time"
)

func s(v ...interface{}) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println(v)
}

// There are some subtle points about the Receive Operator
// http://golang.org/ref/spec#Receive_operator
// http://golang.org/ref/mem#tmp_6

// A receive from an unbuffered channel happens before the send on that channel completes.

var c = make(chan int)
var a string

func f() {
	s("f1")
	a = "hello, world"

	s("f2")

	x := <-c
	s("x =", x)

	s("f3")
}

func main() {
	s("main1")
	a = "initial value"

	s("main2")
	go f()
	s("main3")

	c <- 0

	s("a =", a)
}
