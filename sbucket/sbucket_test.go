package sbucket

import (
	"math"
	"testing"
)

func TestEnc(t *testing.T) {

	// Enc should handle edge cases 0 and 2^32 -1
	// maxuint := uint64(math.Pow(2, 32) - 1)

	integers := map[int]uint64{
		0: uint64(0),
		1: uint64(math.Pow(2, 32) - 1),
		2: uint64(0),
	}

	encoded_integers := map[int]string{
		0: "0",
		1: "bignumber",
	}

	for i := 0; i < 1; i++ {
		n := integers[i]
		s := encoded_integers[i]
		e := Enc(n, i)
		t.Log(i, n, s, e)
		if e != s {
			t.Fatal(e, "!=", s)
		}

	}

}

// func TestTenStamp(t *testing.T) {
// 	now := time.Now()
// 	t.Log(now.Unix(), now.Nanosecond())

// 	for i := -11; i < 65; i++ {
// 		stamp, size := TenStamp(now, i)
// 		t.Log(i, "\t", size, "\t", stamp)
// 	}

// 	// see if 12 and 15 are ok to drop a digit ??? NO!
// 	m := []int{65, 129, 333, 7777, 88888, 9999999, 3600, 36000, 360000, 3600000, 36000000, 360000000, 3600000000, 36000000000}

// 	for _, i := range m {
// 		stamp, size := TenStamp(now, i)
// 		t.Log(i, "\t", size, "\t", stamp)
// 	}

// }
