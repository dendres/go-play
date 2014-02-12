package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
Case 1: client -> server -> message_time_buckets
* client sends list of time,encoded_msg
* server converts time to file path and writes encoded_msg

Case 2: client -> server -> log.io
* client sends a list of encoded messages
* server parses msg, encodes in log.io format and sends to waiting socket.io connections

server has to decode and re-encode each message for log.io anyway...

Case 1 + 2: client -> server -> message_time_buckets + log.io
* client sends a list of encoded_msgs
* server buffers each msg in memory
* server buffers each parsed message in memory
* disk writer reads the timestamp from a message, forms the file name, opens file, writes buffered unparsed message to disk, close file
* log.io handler processes and encodes the message for log.io


Therefore, the Event and Message format should look something like:

Event:
* time object or time.Unix(), time.Nanoseconds()
* hostname string
* app string ... overloaded = app_name || process_name || process_id
* line_checksum uint64 murmur3, or siphash: this should never be checksum("") or checksum("  ")
* line string (contains the log line after removing time, hostname, and app) don't send empty lines or lines of all whitespace

Message:
* some kind of header essentially compressing the rest of the messages?????
* dedup "line" strings by checksum
* tokenize away hostname if possible
* tokenize away the high order bits of time is possible.
* list of tokenized events

Message2:
* hash table of lines... map[line_checksum]string of lines or something?????
* []string of unique hostnames
* []string of unique app names
* []Event{time, line_checksum, hostname_index, app_index}
* []DuplicateEvent[]
this handles large string deduplication

can't have a format that requires reading all the messages into memory during encoding or decoding!

could implement huffman encoding with space separated tokens???  where space is also a token


* json-compatable nested hash and arrays of string, float, int, bool, etc...
* schema free
* does not need to be streamable because the routing/indexing fields will be part of the kafka protocol
* based on existing standards
* clean utf-8 strings would be nice
* 8-bit clean binary fields would also be nice
* relatively easy to port the decoder to a new language


input format:
* syslog lines, json, xml and various other formats become some kind of go structure??
* how does heka parse plaintext with regex????





ideally it should have a table of field names like:
key1,key2,key3,key4,0:value1,1:value2



existing research:
* https://github.com/eishay/jvm-serializers/wiki
* http://leopard.in.ua/2013/10/13/binary-serialization-formats/
* http://en.wikipedia.org/wiki/Comparison_of_data_serialization_formats
* http://web-resource-optimization.blogspot.com/2011/06/json-compression-algorithms.html
* http://mainroach.blogspot.com/2013/08/boosting-text-compression-with-dense.html


can an inverted index of space separated terms be used to deduplicate words in a group of thousands of log lines?

can a frequency tree be made (like the hoffman encoding tree) to assign replacement tokens to words by their frequency?

if the tokens are too big, then the tree will be really big....
  would the index + compressed document be bigger than the original document?


can a tree structure of words be made to create a "term dictionary" use for:
* compression and decompression
* quickly tell you if the word exists in the file

log_line_storage_format:
* start with a series of space separated words
* return a set of int32 or int16

prefix code???
dictionary coder??
"Huffword"
concordance = list of words that appear in a book

maybe need a max word length to prevent encoding all those sha256 hashes!  XXX pull those out and make a separate index!!!!!!

token replacement: pull the
fixed length encoding.  before writing the files, we'll know how many words need encoded and can choose 2,3,4, or 5 byte encoding as needed.

the position of the word in the index is the code for the word!


might want to pull out and convert to binary when tokenizing:
* ip addresses convert to uint32
* hex strings 32 characters or longer... md5, sha1 etc..., convert pairs of characters into a byte... right?


then pull out "words" for the rest of the index:
* [a-zA-Z0-9] 3 to ? characters in length... then downcase them! lossy, but OK!



lets get some stats about tokens to be encoded....









https://github.com/jmckaskill/go-capnproto  Notes:
https://github.com/jmckaskill/go-capnproto/blob/master/doc.go
* Most getters/setters don't return error, instead they return zero value????? not very idiomatic I would say
* errors are provided on the lower level object methods!!!!

candidates:
* capn(event_object) optional protocol specific packing
  https://github.com/jmckaskill/go-capnproto

* kryo????
* msgpack

* thrift??
* compress(hpack(json(event_object)))

http://msgpack.org/??? not tabular???

how about creating a static hoffman encoding ahead of time and sharing it.


http://blog.golang.org/gobs-of-data
* gob streams are self-describing
* GobEncoder, GobDecoder interfaces

http://stackoverflow.com/questions/11202058/unable-to-send-gob-data-over-tcp-in-go-programming
* example







*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
