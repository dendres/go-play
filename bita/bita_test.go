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
	stamp := time.Now().Unix()
	path := base_path + "/" + strconv.FormatInt(stamp, 10)
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
	size    int64
	blksize int64
	blocks  int64
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
	t.Log("Blksize =", stat.Blksize)
	t.Log("Blocks =", stat.Blocks)

	expected := c.blksize * c.blocks
	actual := stat.Blksize * stat.Blocks
	if actual > expected {
		t.Fatal("actual =", actual, "> expected =", expected)
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
	tb15 := int64(140500000000000) // 15T not taking the time to find exact largest number without error.. in theory 16T
	tb7 := int64(70250000000000)
	cases := []Case{
		SetGetCase{[]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		SetGetCase{[]int64{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
		FileSizeCase{[]int64{tb15, 0}, 17179869185, 4096, 16},
		FileSizeCase{[]int64{0, tb7, tb15}, 17179869185, 4096, 24},
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

// empty m1.large: 200000, 10300 ns/op
func BenchmarkRandomSet(b *testing.B) {
	ba := newBenchBita(b)
	b32 := int64(1 << 32)
	rand.Seed(time.Now().UTC().UnixNano())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := rand.Int63n(b32)
		if err := ba.Set(value); err != nil {
			b.Fatal("Error Setting index =", i, ", err =", err)
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
