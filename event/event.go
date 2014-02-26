// package event makes an event []byte from a data []byte + routing paramaters
package event

import (
	"hash/crc32"
)

/*
Event Serialization:
* event[0:3] crc of event[4:]
* event[4] 2 bit encoding + 2 bit replication + 2 bit priority + 2 bit time accuracy
* event[5:12] time ns since unix epoch UTC overflows in year 2554
* event[13:15] length of data NOT including EOE=0xFF
* event[16:len(event)-1] data
* event[len(event)-1:] EOE=0xFF 1 byte "end of event" that will never appear in utf-8

crc + data_length + EOE in serialization provides a high probability of detecting a partial or corrupted events
and allows for several variations and optimizations depending on the data in hand.
*/

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

// Encode takes data + routing attributes and returns a byte slice with header and footer
//  data, replication factor 0-3, priority 0-3, nanoseconds since epoch, and time accuracy 0-3
// and returns a byte slice of the string with header and footer added
func Encode(data []byte, enc int, repl int, pri int, point uint64, acc int) ([]byte, error) {
	head := int(0)
	head |= enc << 6
	head |= repl << 4
	head |= pri << 2
	head |= acc

	event := make([]byte, 16, 256)
	event[4] = byte(head)
	event[5] = byte(point >> 56)
	event[6] = byte(point >> 48)
	event[7] = byte(point >> 40)
	event[8] = byte(point >> 32)
	event[9] = byte(point >> 24)
	event[10] = byte(point >> 16)
	event[11] = byte(point >> 8)
	event[12] = byte(point)

	data_length := len(data)
	event[13] = byte(data_length >> 16)
	event[14] = byte(data_length >> 8)
	event[15] = byte(data_length)

	event = append(event, data...) // XXXX might be expensive?????
	event = append(event, byte(0xFF))

	crc := crc32.ChecksumIEEE(event[4:])
	event[0] = byte(crc >> 24)
	event[1] = byte(crc >> 16)
	event[2] = byte(crc >> 8)
	event[3] = byte(crc)

	return event, nil
}

func Decode([]byte) ([]byte, error) {
}
