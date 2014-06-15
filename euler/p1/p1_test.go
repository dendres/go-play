package p0

import (
	"fmt"
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

func BenchmarkAdd(b *testing.B) {
	out := SumMultiplesBelow(1000)
	fmt.Println(out)
}
