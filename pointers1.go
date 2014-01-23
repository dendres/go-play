// demonstrate passing values and pointers to and from functions
// XXX not there yet

package main

import (
	"fmt"
	"github.com/davecheney/profile"
)

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
