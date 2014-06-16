package q

import (
	"testing"
)

func BenchSetup(b *testing.B) {
	var Bench2 = func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := "hello there"
			a += "again"
		}
	}
}
