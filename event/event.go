// package event makes an event []byte from a data []byte + routing paramaters
package event

import (
	"fmt"
	"hash/crc32"
)

/*
Event Serialization:
* event[0:3] crc of event[4:] including the EOE=0xFF
* event[4] 2 bit encoding + 2 bit replication + 2 bit priority + 2 bit time accuracy
* event[5:12] time ns since unix epoch UTC overflows in year 2554
* event[13:15] length of data NOT including EOE=0xFF
* event[16:len(event)-1] data
* event[len(event)-1:] EOE=0xFF 1 byte "end of event" that will never appear in utf-8

crc + data_length + EOE in serialization provides a high probability of detecting a partial or corrupted events
and allows for several variations and optimizations depending on the data in hand.
*/

const headersize = 16
const headerfootersize = 17

// NewEventHeaderBuffer returns a byte slice sized for reading the event header.
func NewEventHeaderBuffer() []byte {
	return make([]byte, headersize, headersize)
}

// ReadDataLength reads the length bytes from the provided byte slice and updates the provided *int.
func ReadDataLength(b []byte, l *int) {
	*l = 0
	*l |= int(b[13]) << 16
	*l |= int(b[14]) << 8
	*l |= int(b[15]) << 0
}

// An EventBytes is a byte slice of the encoded event including header and footer
type EventBytes struct {
	b []byte
}

// Bytes returns the whole event byte slice.
func (e *EventBytes) Bytes() []byte {
	return e.b
}

// Crc returns the 4 byte IEEE crc32 as an int.
// The CRC counts every byte in the event except the 4 bytes of CRC/
func (e *EventBytes) CRC() (crc uint32) {
	crc |= uint32(e.b[0]) << 24
	crc |= uint32(e.b[1]) << 16
	crc |= uint32(e.b[2]) << 8
	crc |= uint32(e.b[3])
	return crc
}

// DL returns the 3 byte data length as int.
// this does not count the EOE=0xFF
func (e *EventBytes) DL() (dl int) {
	dl |= int(e.b[13]) << 16
	dl |= int(e.b[14]) << 8
	dl |= int(e.b[15])
	return dl
}

// Data returns the data without header or footer.
// XXX better to use len(e) or e.DL() ??????
func (e *EventBytes) Data() []byte {
	return e.b[16 : 16+e.DL()]
}

// EOE returns the last byte of the event.
// it is expected to equal 0xFF but is not checked here.
// XXX better to use len(e) or e.DL() ??????
func (e *EventBytes) EOE() int {
	return int(e.b[len(e.b)-1])
}

// PointBytes returns a slice of the 8 bytes of the uint64 event point.
// This is used for fast access routing and writing.
func (e *EventBytes) PointBytes() []byte {
	return e.b[5:12]
}

// Point returns the event time as a uint64 nanoseconds since unix epoch.
func (e *EventBytes) Point() (p uint64) {
	p |= uint64(e.b[5]) << 56
	p |= uint64(e.b[6]) << 48
	p |= uint64(e.b[7]) << 40
	p |= uint64(e.b[8]) << 32
	p |= uint64(e.b[9]) << 24
	p |= uint64(e.b[10]) << 16
	p |= uint64(e.b[11]) << 8
	p |= uint64(e.b[12])
	return p
}

// CheckCRC calculates and compares the CRC and returns true if the message is valid.
func (e *EventBytes) CheckCRC() bool {
	check := crc32.ChecksumIEEE(e.b[4:])
	if e.CRC() == check {
		return true
	}
	return false
}

// CheckEOE reads the last byte of the message and returns true if it is 0xFF
func (e *EventBytes) CheckEOE() bool {
	if e.EOE() == 0xFF {
		return true
	}
	return false
}

// Encoding returns the event encoding format 0-3 as an int.
func (e *EventBytes) Encoding() (h int) {
	return int(e.b[4]) >> 6 & 0x03
}

// Replication returns the event Replication count as an int.
func (e *EventBytes) Replication() (h int) {
	return int(e.b[4]) >> 4 & 0x03
}

// Priority returns the event Priority as an int.
func (e *EventBytes) Priority() (h int) {
	return int(e.b[4]) >> 2 & 0x03
}

// TimeAccuracy returns the time accuracy estimate multiplier as an int.
func (e *EventBytes) TimeAccuracy() (h int) {
	return int(e.b[4]) & 0x03
}

// An Event contains data and routing attributes
type Event struct {
	// encoding format:
	// 0 = utf-8 string including json, xml, and any other text based serializations
	// 1 = gzip(utf-8 string)
	// 2 = Internal Binary Format
	// 3 = Reserved Internal Binary Format
	enc int

	// replication count:
	// the number of additional log processing servers which must receive the event
	// before the client will discard the local copy
	// 0 = the receiving server will not replicate to any other server
	// 1 = " " replicate to 1 other server
	// 2 = " " 2 other servers
	// 3 = " " 3 others
	// no higher replication count is permitted
	repl int

	// routing priority:
	// Priority queues are used throughout the event processing implementation.
	// Higher priority events will have:
	// * a higher probability of delivery
	// * a possibly lower latency
	// * a higher cost of delivery
	//  3 Action : immediate human intervention is required when this event is received
	//  2 Error  : non-actionable application error
	//  1 Info   : noteworthy condition or informational message
	//  0 Debug  : debug-level events are usually only readable by one engineer
	pri int

	// a crude time time accuracy estimate
	// "point" must always be in ns, but the low bits may be junk.
	// XXX this should be replaced with a specification directly from ntp or ptp
	// 0 = 10^-(0*3) = 1s
	// 1 = 10^-(1*3) = 1ms (default for cloud systems)
	// 2 = 10^-(2*3) = 1us (default for 10G PTP with GPS)
	// 3 = 10^-(3*3) = 1ns
	acc int // time accuracy

	// unsigned nanoseconds since Jan 1 1970 UTC overflows in year 2554
	point uint64

	// any binary data XXX size limitation?????????
	data []byte
}

// Encode takes an Event and returns a []byte with header and footer added.
func (e *Event) Encode() ([]byte, error) {

	if e.enc < 0 || e.enc > 3 {
		fmt.Errorf("enc must be 0 to 3")
	}

	if e.repl < 0 || e.repl > 3 {
		fmt.Errorf("repl must be 0 to 3")
	}

	if e.pri < 0 || e.pri > 3 {
		fmt.Errorf("pri must be 0 to 3")
	}

	if e.acc < 0 || e.acc > 3 {
		fmt.Errorf("acc must be 0 to 3")
	}

	head := int(0)
	head |= e.enc << 6
	head |= e.repl << 4
	head |= e.pri << 2
	head |= e.acc

	event := make([]byte, 16, 256)
	event[4] = byte(head)
	event[5] = byte(e.point >> 56)
	event[6] = byte(e.point >> 48)
	event[7] = byte(e.point >> 40)
	event[8] = byte(e.point >> 32)
	event[9] = byte(e.point >> 24)
	event[10] = byte(e.point >> 16)
	event[11] = byte(e.point >> 8)
	event[12] = byte(e.point)

	data_length := len(e.data)
	if data_length < 0 || data_length > int(1<<24-1) {
		fmt.Errorf("data must be 0 to 2^24-1 bytes long")
	}

	event[13] = byte(data_length >> 16)
	event[14] = byte(data_length >> 8)
	event[15] = byte(data_length)

	event = append(event, e.data...) // XXXX might be expensive?????
	event = append(event, byte(0xFF))

	crc := crc32.ChecksumIEEE(event[4:])
	event[0] = byte(crc >> 24)
	event[1] = byte(crc >> 16)
	event[2] = byte(crc >> 8)
	event[3] = byte(crc)

	return event, nil
}

// Decode takes a byte slice, extracts the fields and returns a new Event.
func Decode(event []byte) (*Event, error) {
	e := Event{}

	crc := uint32(0)
	crc |= uint32(event[0]) << 24
	crc |= uint32(event[1]) << 16
	crc |= uint32(event[2]) << 8
	crc |= uint32(event[3])

	// validate all event bytes other than the crc itself
	check := crc32.ChecksumIEEE(event[4:])
	if check != crc {
		return &e, fmt.Errorf("encoded crc = %x does not match calculated crc = %x", crc, check)
	}

	data_length := int(0)
	ReadDataLength(event, &data_length)

	// read the data
	e.data = event[16 : 16+data_length]

	// validate the EOE=0xFF
	eoe := int(event[16+data_length+1])
	if eoe != int(0xFF) {
		return &e, fmt.Errorf("last byte of event is %x, but needs to be %x", eoe, 0xFF)
	}

	// test with 0xaa = 10101010 and 0x1B = 00011011
	head := int(event[4])
	e.enc = head >> 6 & 0x03
	e.repl = head >> 4 & 0x03
	e.pri = head >> 2 & 0x03
	e.acc = head & 0x03

	e.point = 0
	e.point |= uint64(event[5]) << 56
	e.point |= uint64(event[6]) << 48
	e.point |= uint64(event[7]) << 40
	e.point |= uint64(event[8]) << 32
	e.point |= uint64(event[9]) << 24
	e.point |= uint64(event[10]) << 16
	e.point |= uint64(event[11]) << 8
	e.point |= uint64(event[12])

	return &e, nil
}

// Decode extracts the values from EventBytes and returns an Event.
func (eb *EventBytes) Decode() (*Event, error) {
	event := Event{}

	if !eb.CheckCRC() {
		return &event, fmt.Errorf("checksum fail")
	}
	if !eb.CheckEOE() {
		return &event, fmt.Errorf("truncated message does not end in 0xFF")
	}
	event.data = eb.Data()
	event.point = eb.Point()
	event.enc = eb.Encoding()
	event.repl = eb.Replication()
	event.pri = eb.Priority()
	event.acc = eb.TimeAccuracy()
	return &event, nil
}
