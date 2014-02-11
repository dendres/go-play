// lp32 encodes a uint64 as a padded, sortable string using up to 13 characters in [0-9a-z].
package lp32

import "fmt"

// The encoded rune is looked up in base32map using the string slice operator.
const base32map string = "0123456789abcdefghijklmnopqrstuv"

// Each utf-8 code point or "rune" in the encoded string can be converted to uint64 via a lookup in map32base.
// the runes are defined here using their http://en.wikipedia.org/wiki/UTF-8 hex integer
var map32base = map[rune]uint64{
	0x30: 0, 0x31: 1, 0x32: 2, 0x33: 3, 0x34: 4, 0x35: 5, 0x36: 6, 0x37: 7, 0x38: 8, 0x39: 9,
	0x61: 10, 0x62: 11, 0x63: 12, 0x64: 13, 0x65: 14, 0x66: 15, 0x67: 16, 0x68: 17, 0x69: 18, 0x6A: 19, 0x6B: 20,
	0x6C: 21, 0x6D: 22, 0x6E: 23, 0x6F: 24, 0x70: 25, 0x71: 26, 0x72: 27, 0x73: 28, 0x74: 29, 0x75: 30, 0x76: 31,
}

// slice of the first 13 powers of 32
// this allows "uint64(math.Pow(float64(32), float64(n)))"
// to be replaced with a single lookup power32[n]
// 32^13 overflows uint64
var power32 = map[int]uint64{
	0: 1, 1: 32, 2: 1024, 3: 32768, 4: 1048576, 5: 33554432,
	6: 1073741824, 7: 34359738368, 8: 1099511627776,
	9: 35184372088832, 10: 1125899906842624, 11: 36028797018963968, 12: 1152921504606846976,
}

// Encode the given integer (n) as a base32 string of length count.
func Encode(n uint64, count int) (string, error) {
	if count < 1 {
		return string(""), fmt.Errorf("error encoding n = %d. count = %d is too small to encode", n, count)
	}

	if count > 13 {
		return string(""), fmt.Errorf("error encoding n = %d. count = %d would cause uint64 overflow during decoding", n, count)
		count = 13
	}

	// byte slice containing the encoded characters
	// this will only work for single byte characters!
	bs := make([]byte, count)

	// encode every 5 bits as a character in [0-9a-v]
	for i := (count - 1); i >= 0; i-- {
		bs[i] = base32map[n%32]
		n = n >> 5
	}

	return string(bs), nil
}

// Decode restores the given string to a uint64.
func Decode(txt string) (n uint64, err error) {
	txt_length := len(txt)
	if txt_length < 1 {
		return n, fmt.Errorf("the null string has no corresponding integer")
	}

	if txt_length > 13 {
		return n, fmt.Errorf("string = %s is %d characters too long to decode", txt, txt_length-13)
	}

	var power int
	var i, p, tmp_n uint64
	var ok bool

	for txt_index, c := range txt {
		power = txt_length - txt_index - 1

		i, ok = map32base[c]
		if ok == false {
			return n, fmt.Errorf("invalid rune %x at offset %d while decoding string %s", c, txt_index, txt)
		}

		p, ok = power32[power]
		if ok == false {
			return n, fmt.Errorf("power out of range %d at offset %d while decoding string %s", power, txt_index, txt)
		}

		// base32 conversion: 5*32^2 + 6*32^1 + 7*32^0
		tmp_n = n + (i * p)

		// uint64 overflow loops back to 0. catch this!
		if tmp_n < n {
			return n, fmt.Errorf("uint64 overflow at offset %d while decoding string %s", txt_index, txt)
		}

		n = tmp_n
	}
	return n, nil
}
