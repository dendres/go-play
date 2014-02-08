package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func main() {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte("hello, world\n"))
	w.Close() // You must close this first to flush the bytes to the buffer.
	fmt.Printf("%x", b.Bytes())

}

//https://github.com/antirez/smaz/tree/master
