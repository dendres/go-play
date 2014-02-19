package bita

import (
	"math/rand"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"
)

const base_path = "/opt/kafka/disks/0/test/bita"

// A Case is any test case that can be run
type Case interface {
	Run(t *testing.T)
}

type SetGetCase struct {
	// offsets to set and get
	data []int64
}

func newTestBita(t *testing.T) *Bita {
	stamp := time.Now().Nanosecond()
	path := base_path + "/" + strconv.Itoa(stamp)
	t.Log("new test path =", path)

	b, err := Open(path)
	if err != nil {
		t.Fatal("Error opening path =", path, ", error =", err)
	}

	t.Log("opened b=", b.file.Name())
	return b
}

// cleans out the test directory
// makes 1 file
// sets and gets a series of indexes
func (c SetGetCase) Run(t *testing.T) {

	b := newTestBita(t)

	for _, i := range c.data {
		err := b.Set(i)
		if err != nil {
			t.Fatal("Error Setting index =", i, ", err =", err)
		}
		t.Log("set i=", i)

		var is_set bool
		is_set, err = b.Get(i)
		if err != nil {
			t.Fatal("Error Getting index =", i, ", err =", err)
		}
		t.Log("got", is_set, "from i=", i)

		if is_set == false {
			t.Fatal("did not set index", i)
		}
	}
}

// do some bit array operations and then report on the size of the file
type FileSizeCase struct {
	// data to set
	data []int64

	// expected output of Stat following the inserts
	size   int64 // total file size in bytes
	blocks int64 // number of 512 byte blocks
}

func (c FileSizeCase) Run(t *testing.T) {

	b := newTestBita(t)

	for _, i := range c.data {
		err := b.Set(i)
		if err != nil {
			t.Fatal("Error Setting index =", i, ", err =", err)
		}
		t.Log("set i=", i)
	}

	// check the size of the file!
	stat := syscall.Stat_t{}
	err := syscall.Stat(b.file.Name(), &stat)
	if err != nil {
		t.Fatal("error from syscall.Stat:", err)
	}
	t.Log("After Size =", stat.Size)
	t.Log("Blocks =", stat.Blocks)

	if stat.Size > c.size {
		t.Fatal("file is too large:", stat.Size, ">", c.size)
	}
	if stat.Blocks > c.blocks {
		t.Fatal("too many 512 byte blocks:", stat.Blocks, ">", c.blocks)
	}
}

// run all the test cases
func TestAll(t *testing.T) {

	err := os.RemoveAll(base_path)
	if err != nil {
		t.Fatal("error removing base_path =", base_path, ", err =", err)
	}

	err = os.Mkdir(base_path, 0777)
	if err != nil {
		t.Fatal("error creating base_path =", base_path, ", err =", err)
	}

	// int64max := int64(9223372036854775807)
	//gb128 := int64(99999999999999)

	// tb15 := int64(140500000000000)
	// tb7 := int64(70250000000000)

	cases := []Case{
		SetGetCase{[]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		SetGetCase{[]int64{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},

		FileSizeCase{[]int64{0}, 1, 8}, // 1 byte long file, 8 512 byte blocks = 1 4096 byte filesystem block
		FileSizeCase{[]int64{7}, 1, 8},
		FileSizeCase{[]int64{0, 7}, 1, 8},
		FileSizeCase{[]int64{0, 1, 2, 3, 4, 5, 6, 7}, 1, 8},

		FileSizeCase{[]int64{8}, 2, 8},
		FileSizeCase{[]int64{9}, 2, 8},

		// write a 1KB sparse file
		FileSizeCase{[]int64{8183}, 1023, 8}, // 8*(2^10 -1) -1
		FileSizeCase{[]int64{8184}, 1024, 8}, // 8*(2^10 -1)
		FileSizeCase{[]int64{8191}, 1024, 8}, // 8*(2^10) -1
		FileSizeCase{[]int64{8192}, 1025, 8}, // 8*(2^10)

		// 1MB sparse file
		FileSizeCase{[]int64{8388607}, 1048576, 8}, // 8*(2^17) -1 gives a file sized 2^17

		// 1GB sparse file
		FileSizeCase{[]int64{8589934591}, 1073741824, 8}, // 8*(2^30) -1 gives a file sized 2^30

		// 8TB sparse file
		FileSizeCase{[]int64{70368744177663}, 8796093022208, 8}, // 8*(2^43) -1 gives a file sized 2^43

		// largest bit in 16T should be 8*(2^44)-1, but this errored out
	}

	for _, c := range cases {
		c.Run(t)
	}
}

func newBenchBita(b *testing.B) *Bita {
	stamp := time.Now().Unix()
	path := base_path + "/" + strconv.FormatInt(stamp, 10)
	b.Log("new test path =", path)

	ba, err := Open(path)
	if err != nil {
		b.Fatal("Error opening path =", path, ", error =", err)
	}

	b.Log("opened b=", ba.file.Name())
	return ba
}

// empty m1.large: 500000, 6542 ns/op
func BenchmarkLinearSet(b *testing.B) {
	ba := newBenchBita(b)
	b.ResetTimer()
	for i := int64(0); i < int64(b.N); i++ {
		if err := ba.Set(i); err != nil {
			b.Fatal("Error Setting index =", i, ", err =", err)
		}
	}
}

// empty m1.large 500000, 6875 ns/op
func BenchmarkRandomSmallSet(b *testing.B) {
	ba := newBenchBita(b)
	b24 := int64(1 << 24)
	rand.Seed(time.Now().UTC().UnixNano())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := rand.Int63n(b24)
		if err := ba.Set(value); err != nil {
			b.Fatal("Error Setting index =", i, ", err =", err)
		}
	}
}

// empty m1.large: 100000     15148 ns/op
func BenchmarkRandomSetGet(b *testing.B) {
	ba := newBenchBita(b)
	b32 := int64(1 << 32)
	rand.Seed(time.Now().UTC().UnixNano())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value1 := rand.Int63n(b32)
		value2 := rand.Int63n(b32)

		if err := ba.Set(value1); err != nil {
			b.Fatal("Error Setting index =", value1, ", err =", err)
		}
		if _, err := ba.Get(value2); err != nil {
			b.Fatal("Error Getting index =", value2, ", err =", err)
		}
	}
}

// // empty m1.large: 200000, 305492 ns/op
// func BenchmarkRandomBigSet(b *testing.B) {
// 	ba := newBenchBita(b)
// 	b40 := int64(1 << 40)
// 	rand.Seed(time.Now().UTC().UnixNano())

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		value := rand.Int63n(b40)
// 		if err := ba.Set(value); err != nil {
// 			b.Fatal("Error Setting index =", i, ", err =", err)
// 		}
// 	}
// }
