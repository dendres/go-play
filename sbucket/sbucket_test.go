package sbucket

import (
	"testing"
)

// EncDecCase represents a test case involving the Enc and Dec functions.
// when the package name stops sucking, "EncDecCase" can become "Case" or "Test" or something like that
type EncDecCase struct {
	// the number to be encoded
	Number uint64

	// how many characters to use to encode N
	Count int

	// the string that should result from encoding
	String string

	// the number that can be retrieved from the encoded N... if N is truncated because C is too small
	N2 uint64
}

// Run runs a single EncDec test case.
func (c EncDecCase) Run(t *testing.T) {
	encoded, err := Enc(c.Number, c.Count)
	if err != nil {
		t.Fatal("Error encoding", c.Number, "in", c.Count, "characters:", err)
	}

	if encoded != c.String {
		t.Fatal("Enc should have encoded", c.Number, "in", c.Count, "characters as", c.String, "not", encoded)
	}

	decoded, err := Dec(encoded)
	if err != nil {
		t.Fatal("Error decoding", encoded, " err =", err)
	}

	if decoded != c.N2 {
		t.Fatal("Dec should have decoded", encoded, "as", c.N2, "not", decoded)
	}

	t.Log(c.Number, "encoded in", c.Count, "characters =", encoded, ", then decoded =", decoded)
}

// XXX run test cases that should fail during decode??????

func TestEnc(t *testing.T) {

	uint64max := uint64(18446744073709551615)

	cases := []EncDecCase{
		// min N
		EncDecCase{0, 1, "0", 0},
		EncDecCase{0, 13, "0000000000000", 0},

		// max N
		EncDecCase{uint64max, 1, "v", uint64(31)},
		EncDecCase{uint64max, 13, "fvvvvvvvvvvvv", uint64max},

		// easy to calculate
		EncDecCase{9, 2, "09", 9},
		EncDecCase{10, 2, "0a", 10},
		EncDecCase{11, 2, "0b", 11},

		EncDecCase{31, 2, "0v", 31},
		EncDecCase{32, 2, "10", 32},
		EncDecCase{33, 2, "11", 33},

		EncDecCase{999, 2, "v7", 999},

		EncDecCase{1023, 3, "0vv", 1023},
		EncDecCase{1024, 3, "100", 1024},
		EncDecCase{1025, 3, "101", 1025},

		// truncating 1
		EncDecCase{1, 1, "1", 1},
		EncDecCase{1, 2, "01", 1},
		EncDecCase{1, 3, "001", 1},
		EncDecCase{1, 4, "0001", 1},
		EncDecCase{1, 5, "00001", 1},
		EncDecCase{1, 6, "000001", 1},
		EncDecCase{1, 7, "0000001", 1},
		EncDecCase{1, 8, "00000001", 1},
		EncDecCase{1, 9, "000000001", 1},
		EncDecCase{1, 10, "0000000001", 1},
		EncDecCase{1, 11, "00000000001", 1},
		EncDecCase{1, 12, "000000000001", 1},
		EncDecCase{1, 13, "0000000000001", 1},

		// truncating uint64max
		EncDecCase{uint64max, 1, "v", uint64(31)}, // 32^1-1
		EncDecCase{uint64max, 2, "vv", uint64(1023)},
		EncDecCase{uint64max, 3, "vvv", uint64(32767)},
		EncDecCase{uint64max, 4, "vvvv", uint64(1048575)},
		EncDecCase{uint64max, 5, "vvvvv", uint64(33554431)},
		EncDecCase{uint64max, 6, "vvvvvv", uint64(1073741823)},
		EncDecCase{uint64max, 7, "vvvvvvv", uint64(34359738367)},
		EncDecCase{uint64max, 8, "vvvvvvvv", uint64(1099511627775)},
		EncDecCase{uint64max, 9, "vvvvvvvvv", uint64(35184372088831)},
		EncDecCase{uint64max, 10, "vvvvvvvvvv", uint64(1125899906842623)},
		EncDecCase{uint64max, 11, "vvvvvvvvvvv", uint64(36028797018963967)},
		EncDecCase{uint64max, 12, "vvvvvvvvvvvv", uint64(1152921504606846975)}, // 32^2-1

		// truncating non-zero and non-max numbers????

	}

	for _, c := range cases {
		c.Run(t)
	}

	// make sure these cases are caught as errors
	// EncDecCase{0, 15, "000000000000000", 0},

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
