package event

import (
	"testing"
)

type Case interface {
	Run(t *testing.T)
}

type EventCase struct {
	event *Event
}

func (c EventCase) Run(t *testing.T) {
	t.Logf("enc = %x, repl = %x, pri = %x, acc = %x, point = %d, data = %v",
		c.event.enc, c.event.repl, c.event.pri, c.event.acc, c.event.point, c.event.data)

	e, err := c.event.Encode()
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
		EventCase{&Event{0, 1, 2, 3, t1, d1}},

		// XXXX need many more cases!!!!

	}
	for _, c := range cases {
		c.Run(t)
	}
}
