// package event makes an event []byte from a data []byte + routing paramaters
package event

import (
	"fmt"
	"hash/crc32"
)

/*
Event Serialization for files and network
* event[0:3] crc of event[4:] including the EOE=0xFF
* event[4] 2 bit encoding + 2 bit replication + 2 bit priority + 2 bit time accuracy
* event[5:12] time ns since unix epoch UTC overflows in year 2554
* event[13:15] length of data NOT including EOE=0xFF
* event[16:len(event)-1] data
* event[len(event)-1] EOE=0xFF 1 byte "end of event" that will never appear in utf-8

crc + data_length + EOE in serialization provides a high probability of detecting a partial or corrupted events
and allows for several variations and optimizations depending on the data in hand.

Rough Wire and Disk Read process:
hb := NewEventHeaderBuffer()
thing.Read(hb)
dl := DataLengthFromHeader(hb)
... XXX
eb := NewEventBuffer(hb)
thing.Read(eb)

XXX Don't implement this here because retry and error handling process may not be the same for all Readers!?!?!


*/

const EOE = byte(0xFF)
const HeaderSize = 16
const HeaderFooterSize = 17
const MaxDataLength = int(1<<24 - 1)

// NewEventHeaderBuffer returns a byte slice sized for reading the event header.
func NewEventHeaderBuffer() []byte {
	return make([]byte, HeaderSize)
}

// DataLengthFromHeader returns the int length of the data as read from the event header.
func DataLengthFromHeader(header []byte) (int, error) {
	data_length := int(0)

	if len(header) < HeaderSize {
		return data_length, fmt.Errorf("header is not long enough to read the length from: %v", header)
	}

	data_length |= int(header[13]) << 16
	data_length |= int(header[14]) << 8
	data_length |= int(header[15])

	return data_length, nil
}

// An EventBytes is a byte slice of the encoded event including header and footer.
// Getter and Setter methods are provided for each encoded object.
// bytes is protected to ensure it meets the minimum length requirement.
type EventBytes struct {
	bytes []byte
}

// NewEventBytes creates a new EventBytes from the given []byte if it is large enough.
// the event data is guaranteed to be at least 1 byte long.
func NewEventBytes(b []byte) (*EventBytes, error) {
	eb := &EventBytes{b}

	if len(b) < HeaderFooterSize {
		return eb, fmt.Errorf("supplied byte slice is too small to be a valid event")
	}

	return eb, nil
}

// GetBytes returns the encoded byte slice.
func (e *EventBytes) GetBytes() []byte {
	return e.bytes
}

// GetData returns the data byte slice without header or footer.
// using len(e.bytes) here instead of e.GetDataLength() because e.bytes is known to be the correct length following initialization
func (e *EventBytes) GetData() ([]byte, error) {
	last_data_index := len(e.bytes) - 1
	if last_data_index < HeaderSize {
		return []byte{0}, fmt.Errorf("invalid EventBytes is too small to read data from: %v", e.bytes)
	}
	return e.bytes[HeaderSize:last_data_index], nil
}

// SetData writes the data bytes to the data location.
// The []Byte provided must be EXACTLY the correct size.
func (e *EventBytes) SetData(b []byte) error {
	if len(b) != len(e.bytes)-HeaderFooterSize {
		return fmt.Errorf("SetData required length = %d, but got length = %d", len(e.bytes)-HeaderFooterSize, len(b))
	}

	for i := 0; i < len(b); i++ {
		e.bytes[i+HeaderSize] = b[i]
	}

	return nil
}

// GetDataLength returns the 3 byte length of the data as an int.
// It does NOT include the header or footer when counting the length.
func (e *EventBytes) GetDataLength() (dl int) {
	dl |= int(e.bytes[13]) << 16
	dl |= int(e.bytes[14]) << 8
	dl |= int(e.bytes[15])
	return dl
}

// SetDataLength writes the length of the event data to the data length location.
func (e *EventBytes) SetDataLength(dl int) error {
	if dl < 0 || dl > MaxDataLength {
		return fmt.Errorf("data length out of range 0 to 2^24-1: %d", dl)
	}

	e.bytes[13] = byte(dl >> 16)
	e.bytes[14] = byte(dl >> 8)
	e.bytes[15] = byte(dl)

	return nil
}

// GetCRC returns the 4 byte IEEE crc32 as an int.
func (e *EventBytes) GetCRC() (crc uint32) {
	crc |= uint32(e.bytes[0]) << 24
	crc |= uint32(e.bytes[1]) << 16
	crc |= uint32(e.bytes[2]) << 8
	crc |= uint32(e.bytes[3])
	return crc
}

// SetCRC calculates the CRC of bytes[4:] and writes it to the CRC location.
// All event attributes MUST already be written to e.
func (e *EventBytes) SetCRC() {
	crc := crc32.ChecksumIEEE(e.bytes[4:])
	e.bytes[0] = byte(crc >> 24)
	e.bytes[1] = byte(crc >> 16)
	e.bytes[2] = byte(crc >> 8)
	e.bytes[3] = byte(crc)
}

// InvalidCRC calculates and compares the CRC and returns true if the message is invalid.
func (e *EventBytes) InvalidCRC() bool {
	if e.GetCRC() == crc32.ChecksumIEEE(e.bytes[4:]) {
		return false
	}
	return true
}

// GetEOE returns the last byte of the event.
// it is expected to equal 0xFF but is not checked here.
// XXX better to use len(e) or e.GetDataLength() ??????
func (e *EventBytes) GetEOE() byte {
	return e.bytes[len(e.bytes)-1]
}

// SetEOE writes the constant EOE byte to the last byte of the event.
func (e *EventBytes) SetEOE() {
	e.bytes[len(e.bytes)-1] = EOE
}

// InvalidEOE returns true if the last byte of the message is not 0xFF
func (e *EventBytes) InvalidEOE() bool {
	if e.GetEOE() == EOE {
		return false
	}
	return true
}

// GetPointBytes returns a slice of the 8 bytes of the uint64 event point.
// This is used for fast access routing and writing.
func (e *EventBytes) GetPointBytes() []byte {
	return e.bytes[5:12]
}

// SetPointBytes writes the first 8 bytes in the slice provided to the point location.
// added for consistency, but SetPoint should probably be used instead.
func (e *EventBytes) SetPointBytes(b []byte) {
	for i := 0; i < 8; i++ {
		e.bytes[i+5] = b[i]
	}
}

// GetPoint returns the event time as a uint64 nanoseconds since unix epoch.
func (e *EventBytes) GetPoint() (p uint64) {
	p |= uint64(e.bytes[5]) << 56
	p |= uint64(e.bytes[6]) << 48
	p |= uint64(e.bytes[7]) << 40
	p |= uint64(e.bytes[8]) << 32
	p |= uint64(e.bytes[9]) << 24
	p |= uint64(e.bytes[10]) << 16
	p |= uint64(e.bytes[11]) << 8
	p |= uint64(e.bytes[12])
	return p
}

// SetPoint divides a uint64 into bytes and writes them to the Point location.
func (e *EventBytes) SetPoint(p uint64) {
	e.bytes[5] = byte(p >> 56)
	e.bytes[6] = byte(p >> 48)
	e.bytes[7] = byte(p >> 40)
	e.bytes[8] = byte(p >> 32)
	e.bytes[9] = byte(p >> 24)
	e.bytes[10] = byte(p >> 16)
	e.bytes[11] = byte(p >> 8)
	e.bytes[12] = byte(p)
}

// CheckEOE reads the last byte of the message and returns true if it is 0xFF
func (e *EventBytes) CheckEOE() bool {
	if e.GetEOE() == EOE {
		return true
	}
	return false
}

// GetEncoding returns the event encoding format 0-3 as an int.
func (e *EventBytes) GetEncoding() int {
	return int(e.bytes[4]) >> 6 & 0x03
}

// SetEncoding writes the 2 bits of encoding in the event head byte.
func (e *EventBytes) SetEncoding(encoding int) error {
	if encoding < 0 || encoding > 3 {
		return fmt.Errorf("encoding must be 0 to 3, not %d", encoding)
	}
	e.bytes[4] |= byte(encoding << 6)
	return nil
}

// GetReplication returns the event Replication count as an int.
func (e *EventBytes) GetReplication() int {
	return int(e.bytes[4]) >> 4 & 0x03
}

// SetReplication writes the 2 bits of replication to the event head byte.
func (e *EventBytes) SetReplication(replication int) error {
	if replication < 0 || replication > 3 {
		return fmt.Errorf("replication must be 0 to 3, not %d", replication)
	}
	e.bytes[4] |= byte(replication << 4)
	return nil
}

// GetPriority returns the event Priority as an int.
func (e *EventBytes) GetPriority() int {
	return int(e.bytes[4]) >> 2 & 0x03
}

// SetPriority writes the 2 bit priority to the event head byte.
func (e *EventBytes) SetPriority(priority int) error {
	if priority < 0 || priority > 3 {
		return fmt.Errorf("priority must be 0 to 3, not %d", priority)
	}
	e.bytes[4] |= byte(priority << 2)
	return nil
}

// GetTimeAccuracy returns the time accuracy estimate multiplier as an int.
func (e *EventBytes) GetTimeAccuracy() int {
	return int(e.bytes[4]) & 0x03
}

func (e *EventBytes) SetTimeAccuracy(ta int) error {
	if ta < 0 || ta > 3 {
		return fmt.Errorf("time accuracy must be 0 to 3, not %d", ta)
	}
	e.bytes[4] |= byte(ta)
	return nil
}

// An Event contains data and routing attributes
type Event struct {
	// encoding format:
	// 0 = utf-8 string including json, xml, and any other text based serializations
	// 1 = gzip(utf-8 string)
	// 2 = Internal Binary Format
	// 3 = Reserved Internal Binary Format
	Enc int

	// replication count:
	// the number of additional log processing servers which must receive the event
	// before the client will discard the local copy
	// 0 = the receiving server will not replicate to any other server
	// 1 = " " replicate to 1 other server
	// 2 = " " 2 other servers
	// 3 = " " 3 others
	// no higher replication count is permitted
	Repl int

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
	Pri int

	// a crude time time accuracy estimate
	// "point" must always be in ns, but the low bits may be junk.
	// XXX this should be replaced with a specification directly from ntp or ptp
	// 0 = 10^-(0*3) = 1s
	// 1 = 10^-(1*3) = 1ms (default for cloud systems)
	// 2 = 10^-(2*3) = 1us (default for 10G PTP with GPS)
	// 3 = 10^-(3*3) = 1ns
	Acc int // time accuracy

	// unsigned nanoseconds since Jan 1 1970 UTC overflows in year 2554
	Point uint64

	// any binary data XXX size limitation?????????
	Data []byte
}

// Encode the Event and return the new EventBytes.
func (e *Event) Encode() (*EventBytes, error) {
	eb := &EventBytes{make([]byte, HeaderFooterSize+len(e.Data))}

	// write the data first, so any errors returned have the data that caused the error.
	err := eb.SetData(e.Data)
	if err != nil {
		return eb, err
	}

	err = eb.SetDataLength(len(e.Data))
	if err != nil {
		return eb, err
	}

	err = eb.SetEncoding(e.Enc)
	if err != nil {
		return eb, err
	}

	err = eb.SetReplication(e.Enc)
	if err != nil {
		return eb, err
	}

	err = eb.SetPriority(e.Enc)
	if err != nil {
		return eb, err
	}

	err = eb.SetTimeAccuracy(e.Enc)
	if err != nil {
		return eb, err
	}

	eb.SetPoint(e.Point)
	eb.SetEOE()
	eb.SetCRC()
	return eb, nil
}

// Decode extracts the fields from EventBytes and returns the new Event.
func (eb *EventBytes) Decode() (*Event, error) {
	e := &Event{}

	if eb.InvalidEOE() {
		return e, fmt.Errorf("Invalid EOE")
	}

	if eb.InvalidCRC() {
		return e, fmt.Errorf("Invalid CRC")
	}

	var err error
	e.Data, err = eb.GetData()
	if err != nil {
		return e, err
	}

	e.Point = eb.GetPoint()

	e.Enc = eb.GetEncoding()
	e.Repl = eb.GetReplication()
	e.Pri = eb.GetPriority()
	e.Acc = eb.GetTimeAccuracy()

	return e, nil
}
