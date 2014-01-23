// style notes: when in rome!!!!!!
//
// Indentation:
//     We use tabs for indentation and gofmt emits them by default. Use spaces only if you must.
// Line length:
//     Go has no line length limit. Don't worry about overflowing a punched card.
//         If a line feels too long, wrap it and indent with an extra tab.
// Parentheses:
//     Go needs fewer parentheses than C and Java:
//     control structures (if, for, switch) do not have parentheses in their syntax.
//     Also, the operator precedence hierarchy is shorter and clearer, so
//         x<<8 + y<<16
//     means what the spacing implies, unlike in the other languages.
// Comments:
//     Doc comments work best as complete sentences, which allow a wide variety of automated presentations. The first sentence should be a one-sentence summary that starts with the name being declared.
//     no need for comment formatting like banners of stars etc...
//     plain text only. avoid HTML or other annotations
//
// Name conventions in go: long

package main

import (
	"fmt"
)

// Sqrt returns an approximation to the square root of x.
// It uses Newton's method to iterate toward a more accurate guess
// shamelessly stolen from go tour
func Sqrt(rooted float64) float64 {
	guess := 1.0
	for iterator := 0; iterator < 1000; iterator++ {
		guess -= (guess*guess - rooted) / (2 * guess)
	}
	return guess
}

func main() {
	y := Sqrt(5)
	fmt.Println("sqrt of 5 =", y)
}
