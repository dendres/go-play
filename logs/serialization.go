package main

import (
	"fmt"
	"net"
	"time"
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

* remove big binary or number-ish things
* remove all punctuation
* anything left that is longer than 3 characters can be added to the index

if an efficient "characters that make up multilingual words" character set can be constructed, then this might ok.


lets get some stats about tokens to be encoded....



XXX need a stemmed + downcased index for search + exact match index for compression
* cool because it can quickly give the exact full term responses based on a stemmed search

so... the data set then becomes:
* document with terms replaced by term index
* term index list
* downcased, stemmed index with pointers to the term index
  should be pretty small
  made by parsing the index





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

const (
	null_string     = ``
	space_character = 0x20
	dot_character   = 0x2E
	colon_character = 0x3A

	number_min = 0x30
	number_max = 0x39

	uppercase_min = 0x41
	uppercase_max = 0x5a

	uppercase_hex_max = 0x46

	lowercase_min = 0x61
	lowercase_max = 0x7a

	lowercase_hex_max = 0x66
)

// token_character returns true if the rune is in [0-9a-zA-Z] and false otherwise.
// should expand the character set to cover most words in most languages
// guessing character frequencies: lowercase_letter_count > number_count > uppercase_letter_count
func token_character(c rune) bool {
	if c >= lowercase_min && c <= lowercase_max {
		return true
	} else if c >= number_min && c <= number_max {
		return true
	} else if c >= uppercase_min && c <= uppercase_max {
		return true
	}
	return false
}

// hex_character returns true if the rune is one of [0-9a-fA-F] and false otherwise.
func hex_character(c rune) bool {
	if c >= number_min && c <= number_max {
		return true
	} else if c >= uppercase_min && c <= uppercase_hex_max {
		return true
	} else if c >= lowercase_min && c <= lowercase_hex_max {
		return true
	}
	return false
}

// ip_character returns true if rune is one of [\.\:0-9a-fA-F] and false otherwise.
func ip_character(c rune) bool {
	if c == dot_character {
		return true
	} else if c == colon_character {
		return true
	}
	return hex_character(c)
}

// Tokenize takes a channel of log lines and writes to a channel of terms
// it never terminates ?????
// XXX add a channel of all of the different types of numbers???? or one channel for all int like things????
func Tokenize(lines <-chan string, words chan<- string, ips chan<- net.IP) {
	var line, ip, word string
	var c rune

	for line = range lines {
		for _, c = range line {

			// XXXX not identifying IP address!!!!
			if ip_character(c) {
				ip += string(c)
			} else {
				if len(word) >= 2 { // rfc1924: 0 = "::"
					if ip_object := net.ParseIP(ip); ip_object != nil {
						ips <- ip_object
					}

				}
				ip = null_string
				continue
			}

			// check if it's a long hex string like 0xafffff or fffaaaa0000000045666544fffaaa ?????

			if token_character(c) {
				word += string(c)
			} else {
				if len(word) > 0 {
					words <- word
				}
				word = null_string
			}
		}
	}
}

func ShowStuff(words <-chan string, ips <-chan net.IP) {

	fmt.Println("reading words")

	for {
		select {
		case ip := <-ips:
			fmt.Println("found an ip =", ip)
		case word := <-words:
			fmt.Println("found a word =", word)
		}
	}

	fmt.Println("finished reading words")
}

func main() {
	test_strings := []string{
		"Received disconnect from 60.199.196.144: 11: Bye ⌘ Bye",
		"DHCPREQUEST on eth0 to 172.31.0.1 port 67 (xid=0x21aeff47)",
		"Anacron started on 2014-01-31",
		"pam_unix(su-l:session): session opened for user root by done(uid=0)",
	}

	fmt.Println("making channels")
	lines := make(chan string, 10)
	words := make(chan string, 10)
	ips := make(chan net.IP, 10)

	fmt.Println("starting ShowWords")
	go ShowStuff(words, ips)

	fmt.Println("starting Tokenize")
	go Tokenize(lines, words, ips)

	for {
		for _, s := range test_strings {
			time.Sleep(1 * time.Second)
			lines <- s
		}
	}

	fmt.Println("finished main")
}

/*
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)
}

*/
