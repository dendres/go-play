package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	start := time.Now()

	sf1 := "/opt/kafka/disks/0/test/sparse_file_test_1"

	f, err := os.Create(sf1) // For read access.
	if err != nil {
		fmt.Println("error creating test file:", err)
	}

	// write to the beginning of the file as usual:
	bw, err := f.Write([]byte("hello"))
	if err != nil {
		fmt.Println("error writing to test file:", err)
	}
	fmt.Println("wrote this many bytes:", bw)

	// then seek to max ext4 file size on 4k blocks = 16T ?
	// seek to 17592186040319 was the largest successful write.
	// this is somewhere between 2^40 and 2^41 bytes
	// the file takes up 8192 bytes.. exactly 2 blocks
	x := []byte("w")
	o := int64(17592186040319)
	bw, err = f.WriteAt(x, o)
	if err != nil {
		fmt.Println("error writing to test file:", err)
	}
	fmt.Println("wrote this many bytes:", bw, "at offset =", o)

	d := time.Now().Sub(start)
	fmt.Println("create and sparse write took", d)

	// reaad size and apparent size????
	// best I can do for now is count the blocks from stat
	stat := syscall.Stat_t{}
	err = syscall.Stat(sf1, &stat)
	if err != nil {
		fmt.Println("error from syscall.Stat:", err)
	}
	fmt.Println("Size =", stat.Size)
	fmt.Println("Blksize =", stat.Blksize)
	fmt.Println("Blocks =", stat.Blocks)

	// conclusions thus far:
	// sparse files can be efficient, but if you ever have to read the whole file for any reason (copy or gzip etc..)
	// then all the bytes must be read which is very slow.

	// next test: does data density impact the number of blocks allocated???
	// seek to 4x block size and write one character.
	// how many blocks are allocated?
	for i := 1; i < 5; i++ {
		x = []byte("B")
		o = int64(i * 4 * 4096)
		bw, err = f.WriteAt(x, o)
		if err != nil {
			fmt.Println("error writing to test file:", err)
		}
		fmt.Println("wrote this many bytes:", bw, "at offset =", o)
	}

	stat = syscall.Stat_t{}
	err = syscall.Stat(sf1, &stat)
	if err != nil {
		fmt.Println("error from syscall.Stat:", err)
	}
	fmt.Println("After Size =", stat.Size)
	fmt.Println("Blksize =", stat.Blksize)
	fmt.Println("Blocks =", stat.Blocks)

	// Blocks = 48
	// du -B=1 gave 28672. which is 17 blocks.

}
