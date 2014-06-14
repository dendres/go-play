package p0

import (
	"testing"
)

func Add(a int, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	out := Add(2, 2)
	if out != 4 {
		t.Fail()
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out := Add(i, i)
		check := i * 2
		if out != check {
			b.Fail()
		}
	}
}
