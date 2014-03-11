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

func TestNewEventBytes(t *testing.T) {
	b1 := []byte{0x01}
	_, err := NewEventBytes(b1)
	if err == nil {
		t.Fatalf("expected an error passing a small byte slice to NewEventBytes")
	}

	b2 := make([]byte, 16)
	_, err = NewEventBytes(b2)
	if err == nil {
		t.Fatalf("NewEventBytes must NOT allow zero length data")
	}

	b3 := make([]byte, 17)
	_, err = NewEventBytes(b3)
	if err != nil {
		t.Fatalf("NewEventBytes must allow a length with 1 byte of data, but instead it gave error = %v", err)
	}
}

type DataCase struct {
	Data []byte
}

func (c *DataCase) Run(t *testing.T) {

	b := make([]byte, len(c.Data)+HeaderFooterSize)
	eb, err := NewEventBytes(b)
	if err != nil {
		t.Fatalf("error creating new event bytes: %v", err)
	}

	err = eb.SetData(c.Data)
	if err != nil {
		t.Fatalf("error setting data: %v", err)
	}

	t.Logf("eb.Bytes() = %v", eb.GetBytes())

	data_out, err := eb.GetData()
	if err != nil {
		t.Fatalf("error getting data: %v", err)
	}

	if bytes.Compare(c.Data, data_out) != 0 {
		t.Fatalf("expected data = %v but got data = %v", c.Data, data_out)
	}
}

func TestGetSetData(t *testing.T) {
	cases := []DataCase{
		{[]byte{0}},
		{[]byte{1}},
		{[]byte{0xFF}},
		{[]byte{1, 2}},
		{[]byte{1, 2, 3}},
	}

	for _, c := range cases {
		c.Run(t)
	}
}

type CRCase struct {
	Expected uint32
	Data     []byte
}

func (c *CRCase) Run(t *testing.T) {
	t.Logf("testing that crc(%v) = 0x%x", c.Data, c.Expected)

	b := make([]byte, len(c.Data)+HeaderFooterSize)
	eb, err := NewEventBytes(b)
	if err != nil {
		t.Fatalf("error creating new event bytes: %v", err)
	}

	err = eb.SetData(c.Data)
	if err != nil {
		t.Fatalf("error setting data: %v", err)
	}
	eb.SetCRC()
	t.Logf("bytes with crc = %v", eb.GetBytes())
	crc := eb.GetCRC()
	if crc != c.Expected {
		t.Fatalf("expected crc = %x, but got crc = %x", c.Expected, crc)
	}

	if eb.InvalidCRC() {
		t.Fatalf("expected a valid crc, but told crc is invalid")
	}
}

func TestGetSetCRC(t *testing.T) {
	// crc32.ChecksumIEEE([]byte{0,0,0,0,0,0,0,0,0,0,0,0,<data>,0})
	cases := []CRCase{
		{0xd1bb79c7, []byte{0}},
		{0xe4270b52, []byte{1, 2}},
		{0xad5d2321, []byte("a")},
		{0xd1022e7, []byte("abcdefghij")},
		{0xae39243e, []byte("C is as portable as Stonehedge!!")},
	}

	for _, c := range cases {
		c.Run(t)
	}
}

// type RDLCase struct {
// 	length int
// 	bytes  []byte
// }

// func (c RDLCase) Run(t *testing.T) {
// 	l := int(0)
// 	ReadDataLength(c.bytes, &l)
// 	if l != c.length {
// 		t.Fatalf("expected length =", c.length, "but got length=", l)
// 	}
// }

// func TestReadDataLength(t *testing.T) {
// 	cases := []Case{
// 		RDLCase{0, make([]byte, 30)},
// 		RDLCase{0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
// 		RDLCase{1, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}},
// 		RDLCase{1 << 8, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
// 		RDLCase{1 << 16, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
// 		RDLCase{((1 << 24) - 1), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x00}},
// 	}
// 	for _, c := range cases {
// 		c.Run(t)
// 	}
// }

// func TestBytes(t *testing.T) {
// 	b := []byte{0x01, 0x02, 0x03, 0xFF}
// 	eb := EventBytes{b}
// 	if bytes.Compare(eb.Bytes(), b) != 0 {
// 		t.Fatalf("expected bytes =", b, "but got bytes =", eb.Bytes())
// 	}
// }

// func TestCRC(t *testing.T) {
// 	crc := []byte{0xFF, 0xFF, 0xFF, 0xFF}
// 	eb := EventBytes{crc}
// 	ebc := eb.CRC()
// 	exp := uint32(4294967295)
// 	if ebc != exp {
// 		t.Fatalf("expected crc =", exp, "but got crc =", ebc)
// 	}
// }

// type EventCase struct {
// 	event *Event
// }

// // test event header byte with 0xaa = 10101010 and 0x1B = 00011011

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
