package main

import (
	"fmt"
	"github.com/davecheney/profile"
)

// http://golang.org/doc/faq#Pointers

// the method's argument should be a value when:
// - For types such as basic types, slices, and small structs, a value receiver is very cheap so unless the semantics of the method requires a pointer, a value receiver is efficient and clear.

// the method's argument should be a pointer when:
// - the method needs to modify the argument
//     slices and maps contain references
//     if the length of the slice needs to change, then pass a pointer
// - if the argument is very large

// try to be consistent... method sets????

type Thing struct {
	a, b int
}

func (t Thing) p() {
	fmt.Println("thing = ", t.a, t.b)
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	x := Thing{a: 11, b: 12}
	x.p()

	y := &Thing{a: 22, b: 33}
	y.p()
}
