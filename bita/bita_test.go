package bita

import (
	"os"
	"strconv"
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

// GetEmptyCase!!!!

// cleans out the test directory
// makes 1 file
// sets and gets a series of indexes
func (c SetGetCase) Run(t *testing.T) {

	stamp := time.Now().Unix()
	path := base_path + "/" + strconv.FormatInt(stamp, 10)
	t.Log("new test path =", path)

	b, err := Open(path)
	if err != nil {
		t.Fatal("Error opening path =", path, ", error =", err)
	}
	t.Log("opened b=", b)

	for _, i := range c.data {
		err = b.Set(i)
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

func TestSetGet(t *testing.T) {

	err := os.RemoveAll(base_path)
	if err != nil {
		t.Fatal("error removing base_path =", base_path, ", err =", err)
	}

	err = os.Mkdir(base_path, 0777)
	if err != nil {
		t.Fatal("error creating base_path =", base_path, ", err =", err)
	}

	// int64max := int64(9223372036854775807)
	gb128 := int64(137438953472)
	cases := []Case{
		SetGetCase{[]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		SetGetCase{[]int64{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
		SetGetCase{[]int64{gb128}},
	}

	for _, c := range cases {
		c.Run(t)
	}
}

// t.Log(c.Number, "encoded in", c.Count, "characters =", encoded, ", then decoded =", decoded)

// around 300 ns/op on m1.large
// func BenchmarkEncode(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Encode(uint64(i), 13)
// 	}
// }
