package event

import (
	"testing"
)

type Case interface {
	Run(t *testing.T)
}

type EventCase struct {
	data  []byte
	enc   int
	repl  int
	pri   int
	point uint64
	acc   int
}

func (c EventCase) Run(t *testing.T) {
	t.Logf("enc = %x, repl = %x, pri = %x, point = %d, acc = %x, data = %v",
		c.enc, c.repl, c.pri, c.point, c.acc, c.data)

	e, err := Encode(c.data, c.enc, c.repl, c.pri, c.point, c.acc)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	// make a Decode XXXXXXXXXXXXX

	t.Logf("encoded = %v", e)
}

func TestEncDec(t *testing.T) {

	d1 := []byte{0x00, 0x00, 0x01, 0x02, 0x03, 0x00, 0x00}
	t1 := uint64(0x0808080808080808)
	cases := []Case{
		EventCase{d1, 0, 1, 2, t1, 1},

		// XXXX need many more cases!!!!

	}
	for _, c := range cases {
		c.Run(t)
	}
}
