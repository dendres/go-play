package coverage

import (
	"fmt"
	"testing"
)

type Test struct {
	in  int
	out string
}

var tests = []Test{
	{-1, "negative"},
	{0, "zero"},
	{5, "small"},
}

// TestSize demonstrates marginal coverage of the branching in Size.
func TestSize(t *testing.T) {
	for i, test := range tests {
		size := Size(test.in)
		if size != test.out {
			t.Errorf("#%d: Size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}

// ExampleSize gives an example invocation of Size.
func ExampleSize() {
	int := 999
	fmt.Println(int, "is", Size(999))
	// Output: 999 is huge
}

// BenchmarkSize tests repeated calls to Size with a fixed input.
func BenchmarkSize(b *testing.B) {
	zero := 0
	for i := 0; i < b.N; i++ {
		Size(zero)
	}
}
