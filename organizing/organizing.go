package organizing

import (
	"fmt"
)

// Sqrt returns an approximation to the square root of x.
// It uses Newton's method to iterate toward a more accurate guess
// shamelessly stolen from go tour
func Sqrt(rooted float64) float64 {

	// guess starts with 1 and is updated as iteration progresses
	guess := 1.0

	for iterator := 0; iterator < 1000; iterator++ {
		guess -= (guess*guess - rooted) / (2 * guess)
	}

	// BUG(done) guess should be set as the return value in the function statement instead of using a return statement
	return guess
}
