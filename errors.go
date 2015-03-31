package main

import (
	"fmt"
	"io"
)

/*
https://blog.golang.org/errors-are-values

Errors are Values!

if err != nil


If you want to repeat some action
and not worry about errors till after, then:

stuff = pkg.Something()
more_stuff = pkg.Something()
err := pkg.Err();
if err != nil {
   print or return the error
}

or... see the example below for another implementation of this pattern:
*/

type errWriter struct {
	err Error
	w   io.Writer
}

func (ew *errWriter) write(buf []byte) {
	// don't bother writing more if there was an error
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

func main() {
	ew := &errWriter{w: fd}
	ew.write(buf1)
	ew.write(buf2)
	ew.write(buf3)

	if ew.err() != nil {
		fmt.Println("error writing")
	}
	fmt.Println("wrote ok")
}
