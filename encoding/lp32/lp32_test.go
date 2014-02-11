package lp32

import (
	"fmt"
	"testing"
)

// A Case is any test case that can be run
type Case interface {
	Run(t *testing.T)
}

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
	encoded, err := Encode(c.Number, c.Count)
	if err != nil {
		t.Fatal("Error encoding", c.Number, "in", c.Count, "characters:", err)
	}

	if encoded != c.String {
		t.Fatal("Enc should have encoded", c.Number, "in", c.Count, "characters as", c.String, "not", encoded)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal("Error decoding", encoded, " err =", err)
	}

	if decoded != c.N2 {
		t.Fatal("Dec should have decoded", encoded, "as", c.N2, "not", decoded)
	}

	t.Log(c.Number, "encoded in", c.Count, "characters =", encoded, ", then decoded =", decoded)
}

// EncErrorCase is a test case that passes when there is an encoding error.
type EncErrorCase struct {
	// the number to be encoded
	Number uint64

	// how many characters to use to encode N
	Count int

	// the output expected during the error case
	String string

	// the error expected
	Error error
}

// Run passes when Enc returns the expected output and error.
func (c EncErrorCase) Run(t *testing.T) {
	t.Log("encoding", c.Number, "in", c.Count, ". expecting output =", c.String, "and error =", c.Error)

	s, e := Encode(c.Number, c.Count)
	if s != c.String {
		t.Fatal("Expected output =", c.String, ", but got output =", s)
	}

	if e.Error() != c.Error.Error() {
		t.Fatal("Expected error =", c.Error.Error(), ", but got error =", e.Error())
	}
}

// DecErrorCase is a test case that passes when there is a decoding error.
type DecErrorCase struct {
	// the input string to decode
	String string

	// the output expected during the error case
	Number uint64

	// the error expected
	Error error
}

// RunF passes when Enc or Dec returns an error.
// the output of Enc and Dec are ignored
func (c DecErrorCase) Run(t *testing.T) {
	t.Log("decoding", c.String, ". expecting output =", c.Number, "and error =", c.Error)

	n, e := Decode(c.String)
	if n != c.Number {
		t.Fatal("Expected output =", c.Number, ", but got output =", n)
	}

	if e.Error() != c.Error.Error() {
		t.Fatal("Expected error =", c.Error.Error(), ", but got error =", e.Error())
	}
}

func TestEncDec(t *testing.T) {

	uint64max := uint64(18446744073709551615)

	cases := []Case{
		// Enc bad input
		EncErrorCase{0, 0, "", fmt.Errorf("error encoding n = 0. count = 0 is too small to encode")},
		// EncErrorCase{0, 0, "0", fmt.Errorf("error encoding n = 0. count = 0 is too small to encode")},
		// EncErrorCase{0, 0, "", fmt.Errorf("fake error message")},
		EncErrorCase{0, 14, "", fmt.Errorf("error encoding n = 0. count = 14 would cause uint64 overflow during decoding")},
		EncErrorCase{0, 9223372036854775807, "", fmt.Errorf("error encoding n = 0. count = 9223372036854775807 would cause uint64 overflow during decoding")},
		// EncErrorCase{0, 9223372036854775808 2^63 overflows "int"??? looks like int64 to me!!!
		EncErrorCase{0, -1, "", fmt.Errorf("error encoding n = 0. count = -1 is too small to encode")},
		EncErrorCase{0, -9223372036854775808, "", fmt.Errorf("error encoding n = 0. count = -9223372036854775808 is too small to encode")},
		// EncErrorCase{0, -9223372036854775809  (2^63)+1 overflows "int"???? looks like int64 to me!!!

		// Dec bad input
		DecErrorCase{"", 0, fmt.Errorf("the null string has no corresponding integer")},
		DecErrorCase{"00000000000000", 0, fmt.Errorf("string = 00000000000000 is 1 characters too long to decode")},
		DecErrorCase{"w", 0, fmt.Errorf("invalid rune 77 at offset 0 while decoding string w")},
		DecErrorCase{"⌘", 0, fmt.Errorf("invalid rune 2318 at offset 0 while decoding string ⌘")},

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
		EncDecCase{999, 4, "00v7", 999},
		EncDecCase{999, 1, "7", 7},

		EncDecCase{1023, 3, "0vv", 1023},
		EncDecCase{1024, 3, "100", 1024},
		EncDecCase{1025, 3, "101", 1025},

		// padding 1
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

	}

	for _, c := range cases {
		c.Run(t)
	}
}

// around 300 ns/op on m1.large
func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(uint64(i), 13)
	}
}

// around 1000 ns/op on m1.large
func BenchmarkDecode(b *testing.B) {
	test_string := "vvvvvvvvvvvv"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(test_string)
	}
}
