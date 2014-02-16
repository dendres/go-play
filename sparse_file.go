package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("/tmp/sparse_file_test_1") // For read access.
	if err != nil {
		fmt.Println("error creating test file:", err)
	}

	bw, err := f.Write([]byte("hello"))
	if err != nil {
		fmt.Println("error writing to test file:", err)
	}
	fmt.Println("wrote this many bytes:", bw)

	x := []byte("second word to be written")
	o := int64(10000)
	bw, err = f.WriteAt(x, o)
	if err != nil {
		fmt.Println("error writing to test file:", err)
	}
	fmt.Println("wrote this many bytes:", bw, "at offset =", o)

}
