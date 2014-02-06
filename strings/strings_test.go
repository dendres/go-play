package messedupstrings

import (
	"fmt"
	"testing"
)

/*
find out how strings are supposed to work in go

*/
func Test1(t *testing.T) {

	// http://blog.golang.org/strings
	// * In Go, a string is in effect a read-only slice of bytes.
	// * a string holds arbitrary bytes

	// string literals are ALWAYS utf-8, but strings can contain arbitrary bytes

	// in theory this should turn out to be something readable:
	const sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
	var a, b string

	// %x prints the string as hex digits, 2 per byte
	a = fmt.Sprintf("% x", sample)
	b = "bd b2 3d bc 20 e2 8c 98"
	if a != b {
		t.Fatal(a, "!=", b)
	}

	// %q will escape any non-printable byte sequences in a string so the output is unambiguous
	// double quoted strings can contain escape sequences like  "\u2318"
	// Source code in Go is defined to be UTF-8 text; no other representation is allowed.
	a = fmt.Sprintf("% q", sample)
	b = `\xbd\xb2=\xbc ⌘` // XXX not working... what are a and b???
	if a != b {
		fmt.Println(a, "!=", b)
	}

	// %+q should interpret utf-8 and escape only non-printable byte sequences
	a = fmt.Sprintf("%+q", sample)
	// XXX not sure how to test output?????

	// what is normally thought of as a "character" is called a "rune" in go
	// "rune" is an alias for the type "int32"
	// new_rune := `⌘` // this is a rune with the integer value 0x2318

	// %#U shows a (rune or byte?????)'s unicode value

	// the for range loop does successfully loop over rune's instead of bytes:
	const nihongo = "日本語"
	for index, runeValue := range nihongo {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}

	// http://blog.golang.org/normalization

}
