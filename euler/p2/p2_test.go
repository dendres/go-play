package p2

import (
	"fmt"
	"testing"
)

// sum the even-valued fibonacci terms <= 4 million

// memoize fib
// sum += f if f % 2 == 0

// step1: fib
func Fib1(max_fib_val uint64) uint64 {
	var sum, l, zero uint64
	two := uint64(2)

	for f := uint64(1); sum <= max_fib_val; l, f = f, f+l {
		fmt.Println(f, sum)
		if f%two == zero {
			sum += f
		}
	}
	return sum
}

type Case1 struct {
	In  uint64
	Sum uint64
}

func TestFib(t *testing.T) {

	cases := []Case1{
		Case1{1, 0},
		Case1{2, 2},
		Case1{3, 2},
		Case1{4, 2},
		Case1{5, 2},
		Case1{6, 2},
		Case1{7, 2},
		Case1{8, 10},
		Case1{9, 10},
		Case1{40, 44},
		Case1{4000000, 5},
	}

	for _, c := range cases {
		sum := Fib1(c.In)
		if sum != c.Sum {
			t.Fatal("for", c.In, ", expected", c.Sum, ", but got", sum)
		}
	}
}

// step2: memoize fib

// step3: sum
