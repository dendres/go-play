package event

import (
	"bytes"
	"testing"
)

type Case interface {
	Run(t *testing.T)
}

func TestNewEventHeaderBuffer(t *testing.T) {
	hb := NewEventHeaderBuffer()
	if len(hb) != 16 {
		t.Fatalf("expected event header buffer of length 16, but got length =", len(hb))
	}
}

type RDLCase struct {
	length int
	bytes  []byte
}

func (c RDLCase) Run(t *testing.T) {
	l := int(0)
	ReadDataLength(c.bytes, &l)
	if l != c.length {
		t.Fatalf("expected length =", c.length, "but got length=", l)
	}
}

func TestReadDataLength(t *testing.T) {
	cases := []Case{
		RDLCase{0, make([]byte, 30)},
		RDLCase{0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		RDLCase{1, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}},
		RDLCase{1 << 8, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
		RDLCase{1 << 16, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
		RDLCase{((1 << 24) - 1), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x00}},
	}
	for _, c := range cases {
		c.Run(t)
	}
}

func TestBytes(t *testing.T) {
	// XXXXXXXXXXXX expand this!!!!!
	b := []byte{0x01}
	eb := EventBytes{b}
	if bytes.Compare(eb.Bytes(), b) != 0 {
		t.Fatalf("expected bytes =", b, "but got bytes =", eb.Bytes())
	}
}

type EventCase struct {
	event *Event
}

// func (c EventCase) Run(t *testing.T) {
// 	t.Logf("enc = %x, repl = %x, pri = %x, acc = %x, point = %d, data = %v",
// 		c.event.enc, c.event.repl, c.event.pri, c.event.acc, c.event.point, c.event.data)

// 	e, err := c.event.Encode()
// 	if err != nil {
// 		t.Fatalf("error: %s", err)
// 	}

// 	// make a Decode XXXXXXXXXXXXX

// 	t.Logf("encoded = %v", e)
// }

// func TestEncDec(t *testing.T) {

// 	d1 := []byte{0x00, 0x00, 0x01, 0x02, 0x03, 0x00, 0x00}
// 	t1 := uint64(0x0808080808080808)
// 	cases := []Case{
// 		EventCase{&Event{0, 1, 2, 3, t1, d1}},

// 		// XXXX need many more cases!!!!

// 	}
// 	for _, c := range cases {
// 		c.Run(t)
// 	}
// }
