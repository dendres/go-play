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

// google hint to use geometric series... avoiding the repetative mod operation
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

/*
folks on the forum claim you can calculate the series without iteration????

take +3 series...
0       3       6       9       12
3 * 0 + 3 * 1 + 3 * 2 + 3 * 3 + 3 * 4
3 (0 + 1 + 2.... + n)

http://en.wikipedia.org/wiki/1_%2B_2_%2B_3_%2B_4_%2B_%E2%8B%AF#Partial_sums

0 + 1 + 2 + 3 + 4 + 5 = 15
triangular number because stack points in equilateral triangle
it's like (5*5)/2... but if you do that you skip half of the bottom row of the triangle
so you have to add another row: (5*5 + 5)/2
which can also be written n(n+1)/2, n=5

so now the sum of the nth number of the 3 series is 3( n(n+1)/2  )

also... using uint to avoid having to check for negative number input
*/
func SumSeries(a, n uint) uint {
	return a * (n * (n + 1)) / 2
}

type SSCase struct {
	Series uint
	N      uint
	Sum    uint
}

func TestSumSeries(t *testing.T) {
	cases := []SSCase{
		SSCase{3, 0, 0},
		SSCase{3, 1, 3},
		SSCase{3, 2, 9},
		SSCase{3, 3, 18},
		SSCase{3, 4, 30},
		SSCase{5, 0, 0},
		SSCase{5, 1, 5},
		SSCase{5, 2, 15},
		SSCase{5, 3, 30},
		SSCase{5, 4, 50},
	}

	for _, c := range cases {
		if s := SumSeries(c.Series, c.N); s != c.Sum {
			t.Fatal("SumSeries of", c.N, "should be", c.Sum, ", but got", s)
		}
	}
}

func SB3(max uint) uint {
	max = max - 1 // "less than, not equal 2, max"
	return SumSeries(3, max/3) + SumSeries(5, max/5) - SumSeries(15, max/15)
}

type Case2 struct {
	Max, Expected uint
}

func TestSB3(t *testing.T) {
	cases := []Case2{
		Case2{10, 23},
		Case2{1000, 233168},
	}

	for _, c := range cases {
		if out := SB3(c.Max); out != c.Expected {
			t.Fatal(c.Max, "should be", c.Expected, ", but got", out)
		}
	}
}

// BenchmarkSB3100000000 11.6 ns/op
func BenchmarkSB3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SB3(1000)
	}
}
