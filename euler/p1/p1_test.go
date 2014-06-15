package p0

import (
	"testing"
)

// If we list all the natural numbers below 10 that are multiples of 3 or 5, we get 3, 5, 6 and 9.
// The sum of these multiples is 23.
// Find the sum of all the multiples of 3 or 5 below 1000.

// mod 3 , mod 5
func IsMultiple(x int) bool {
	if x < 1 {
		return false
	}

	if (x%3) == 0 || (x%5) == 0 {
		return true
	}

	return false
}

// BenchmarkSumMultiplesBelow 200000 6976 ns/op
func SumMultiplesBelow(x int) int {
	total := 0
	for i := 1; i < x; i++ {
		if IsMultiple(i) {
			total += i
		}
	}
	return total
}

func TestIsMultiple(t *testing.T) {
	should := []int{3, 5, 6, 9}
	should_not := []int{-1, 0, 1, 2, 4, 7, 8}

	for _, s := range should {
		if !IsMultiple(s) {
			t.Fatal("IsMultiple should not think", s, "is a multiple of 3 or 5")
		}
	}

	for _, n := range should_not {
		if IsMultiple(n) {
			t.Fatal("IsMultiple should think", n, "is a multiple of 3 or 5")
		}
	}
}

func TestSumMultiplesBelow(t *testing.T) {
	out := SumMultiplesBelow(10)
	if out != 23 {
		t.Fatal("the sum of the multiples of 3 or 5 less than 10 should be 23")
	}
}

func BenchmarkSumMultiplesBelow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumMultiplesBelow(1000)
	}
}

// see if you can do it with only addition operations
// XXX multiples of 15 are counted twice... have to subtract to count them only once
// BenchmarkSB2 5000000 456 ns/op
func SB2(max int) int {
	total := 0
	for t := 3; t < max; t = t + 3 {
		total += t
	}
	for f := 5; f < max; f = f + 5 {
		total += f
	}
	for f := 15; f < max; f = f + 15 {
		total -= f
	}
	return total
}

type Case struct {
	Max      int
	Expected int
}

func TestSB2(t *testing.T) {
	cases := []Case{
		Case{10, 23},
		Case{1000, 233168},
	}

	for _, c := range cases {
		if out := SB2(c.Max); out != c.Expected {
			t.Fatal(c.Max, "should be", c.Expected, ", but got", out)
		}
	}
}

func BenchmarkSB2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SB2(1000)
	}
}
