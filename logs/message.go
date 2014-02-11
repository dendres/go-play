package main

import (
	"log"
)

/*
message format as found written to kafka
https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol#AGuideToTheKafkaProtocol-CommonRequestandResponseStructure


The produce API is used to send message sets to the server. For efficiency it allows sending message sets intended for many topic partitions in a single request. The produce API uses the generic message set format, but since no offset has been assigned to the messages at the time of the send the producer is free to fill in that field in any way it likes.


ProduceRequest => RequiredAcks Timeout [TopicName [Partition MessageSetSize MessageSet]]
  RequiredAcks => int16: number of servers to wait for
    0 = no ack, 1 = server's local log, -1 = all in sync replicas
  Timeout => int32: milliseconds to wait for RequiredAcks (NOT STRICT)
  Partition => int32
  MessageSetSize => int32
  Topicname => ?????
  Partition => ?????
  MessageSet => ?????

MESSAGE SET:
  MessageSize: uint32: make([]byte, 4+msgLen):
    The MessageSize field gives the size of the subsequent request or response message in bytes.
    The client can read requests by first reading this 4 byte size as an integer N,
    and then reading and parsing the subsequent N bytes of the request.
    max message length = 2Gb

  MAGIC: 1 byte: msg[4] = m.magic
  COMPRESSION: 1 byte: msg[5] = m.compression
  CHECKSUM: uint32: msg[6:] ... crc32.ChecksumIEEE(message.payload)
  MESSAGE PAYLOAD: msg[10:] ... bytes: codec.Encode(payload)


XXX can a "partition" be a time slice?????
 - kafka does not enforce what data goes in what partition under a topic. clients do.
 - partition is an int32
 - can a single partition be deleted????   "Data is deleted one log segment at a time. "


# from scala: A message. The format of an N byte message is the following:
* 4 byte CRC32 of the message
* 1 byte "magic" identifier to allow format changes, value is 2 currently
* 1 byte "attributes" identifier to allow annotations on the message independent of the version
* 4 byte key length, containing length K
* K byte key ??????
* 4 byte payload length, containing length V
* V byte payload






XXX message set
https://github.com/apache/kafka/blob/0.8/core/src/main/scala/kafka/message/MessageSet.scala
A set of messages with offsets. A message set has a fixed serialized form, though the container for the bytes could be either in-memory or on disk. The format of each message is as follows:
* 8 byte message offset number
* 4 byte size containing an integer N
* N message bytes as described in the Message class



message.magic = byte(MAGIC_DEFAULT) # = 1
message.compression = codec.Id()
message.payload = codec.Encode(payload)





func (m *Message) Encode() []byte {
  msgLen := NO_LEN_HEADER_SIZE + len(m.payload)
  msg := make([]byte, 4+msgLen)
  binary.BigEndian.PutUint32(msg[0:], uint32(msgLen))
  msg[4] = m.magic
  msg[5] = m.compression

  copy(msg[6:], m.checksum[0:])
  copy(msg[10:], m.payload)

  return msg
}




4 byte big-endian int: length of message in bytes (including the rest of the header, but excluding the length field itself)

1 byte: "magic" identifier (format version number)

If the magic byte == 0, there is one more header field:
  4 byte big-endian int: CRC32 checksum of the payload

If the magic byte == 1, there are two more header fields:
  1 byte: "attributes" (flags for compression, codec etc)
  4 byte big-endian int: CRC32 checksum of the payload




topic???
A message. The format of an N byte message is the following:
If magic byte is 0
1. 1 byte "magic" identifier to allow format changes
2. 4 byte CRC32 of the payload
3. N - 5 byte payload
If magic byte is 1
1. 1 byte "magic" identifier to allow format changes
2. 1 byte "attributes" identifier to allow annotations on the message independent of the version (e.g. compression enabled, type of codec used)
3. 4 byte CRC32 of the payload
4. N - 6 byte payload







*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
