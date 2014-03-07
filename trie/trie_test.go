package radix

import (
	"bytes"
	"testing"
)

type Case interface {
	Run(t *testing.T)
}

type SortWordsCase struct {
	bytes  int
	input  [][]byte
	output [][]byte
}

func (c SortWordsCase) Run(t *testing.T) {
	t.Logf("bytes = %d, input = %v", c.bytes, c.input)

	err := SortWords(c.bytes, c.input)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	for i, word := range c.input {
		if !bytes.Equal(word, c.output[i]) {
			t.Fatalf("at i = %d: %v != %v", i, word, c.output[i])
		}
	}

	t.Logf("output           = %v", c.output)
}

func TestEncDec(t *testing.T) {

	b0 := []byte{0x00}
	b9 := []byte{0x99}
	bf := []byte{0xFF}

	// hand made specific edge cases
	cases := []Case{
		SortWordsCase{1, [][]byte{b0}, [][]byte{b0}},
		SortWordsCase{1, [][]byte{bf}, [][]byte{bf}},
		SortWordsCase{1, [][]byte{}, [][]byte{}}, // empty byte slice gets a zero

		SortWordsCase{1, [][]byte{b0, b0}, [][]byte{b0, b0}},
		SortWordsCase{1, [][]byte{bf, bf}, [][]byte{bf, bf}},

		SortWordsCase{1, [][]byte{b0, bf}, [][]byte{b0, bf}},
		SortWordsCase{1, [][]byte{bf, b0}, [][]byte{b0, bf}},

		SortWordsCase{1, [][]byte{b0, b0, b0}, [][]byte{b0, b0, b0}},
		SortWordsCase{1, [][]byte{b0, b0, b9}, [][]byte{b0, b0, b9}},
		SortWordsCase{1, [][]byte{b0, b9, b0}, [][]byte{b0, b0, b9}},
		SortWordsCase{1, [][]byte{b0, b9, b9}, [][]byte{b0, b9, b9}},
		SortWordsCase{1, [][]byte{b9, b0, b0}, [][]byte{b0, b0, b9}},
		SortWordsCase{1, [][]byte{b9, b0, b9}, [][]byte{b0, b9, b9}},
		SortWordsCase{1, [][]byte{b9, b9, b0}, [][]byte{b0, b9, b9}},
		SortWordsCase{1, [][]byte{b9, b9, b9}, [][]byte{b9, b9, b9}},

		SortWordsCase{1, [][]byte{b0, b0, bf}, [][]byte{b0, b0, bf}},
		SortWordsCase{1, [][]byte{b0, bf, b0}, [][]byte{b0, b0, bf}},
		SortWordsCase{1, [][]byte{b0, bf, bf}, [][]byte{b0, bf, bf}},
		SortWordsCase{1, [][]byte{bf, b0, b0}, [][]byte{b0, b0, bf}},
		SortWordsCase{1, [][]byte{bf, b0, bf}, [][]byte{b0, bf, bf}},
		SortWordsCase{1, [][]byte{bf, bf, b0}, [][]byte{b0, bf, bf}},
		SortWordsCase{1, [][]byte{bf, bf, bf}, [][]byte{bf, bf, bf}},

		SortWordsCase{1, [][]byte{b9, b9, b9}, [][]byte{b9, b9, b9}},
		SortWordsCase{1, [][]byte{b9, b9, bf}, [][]byte{b9, b9, bf}},
		SortWordsCase{1, [][]byte{b9, bf, b9}, [][]byte{b9, b9, bf}},
		SortWordsCase{1, [][]byte{b9, bf, bf}, [][]byte{b9, bf, bf}},
		SortWordsCase{1, [][]byte{bf, b9, b9}, [][]byte{b9, b9, bf}},
		SortWordsCase{1, [][]byte{bf, b9, bf}, [][]byte{b9, bf, bf}},
		SortWordsCase{1, [][]byte{bf, bf, b9}, [][]byte{b9, bf, bf}},
		SortWordsCase{1, [][]byte{bf, bf, bf}, [][]byte{bf, bf, bf}},

		// SortWordsCase{2, [][]byte{b0, b0}, [][]byte{b0, b0}},
	}

	// generate all possible combinations of 3 bytes and their
	word := []byte{0x00, 0x01, 0x02}
	words := [][]byte{}
	for _, w1 := range word {
		for _, w2 := range word {
			for _, w3 := range word {
				t.Log(w1, w2, w3)
				ww := []byte{w1, w2, w3}
				words = append(words, ww)

				//c := SortWordsCase{1, [][]byte{w1, w2, w3}, [][]byte{bf, bf, bf}},
				//append(cases, c)
			}
		}
	}

	t.Log(words)
	// now scramble words!!!!

	for _, c := range cases {
		c.Run(t)
	}
}
