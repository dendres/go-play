package main

import (
	"fmt"
	"net"
	"time"
)

/*
layer 6 reliability protocol
keep it independent of other event content


12 Byte Header
==============

* head: 1 byte = 2 bit encoding + 2 bit replication + 2 bit priority + 2 bit time accuracy
  head[:1] = event data encoding format
    0 = Binary
    1 = JSON
    2 = a single utf-8 string
  head[2:3] = replication: 0 to 3
    "replication" in this context is the number of other log processing servers
    which must receive the event before the client agent is permitted to delete the local copy.
    0 = the receiving server will not replicate to any other server
    1 = ... will replicate to 1 other server
    2 = ... 2 other servers
    3 = ... 3 other servers
    no higher replication is permitted
  head[4:5] = priority: 0 to 3
   "priority" in this context is used with "replication" to determine message processing order
   0 Action : immediate human intervention is required when this event is received
   1 Error  : non-actionable application error
   2 Info   : noteworthy condition or informational message
   3 Debug  : debug-level events are usually only readable by one engineer
  head[6:7] = crude time accuracy estimate 0 to 3. "point" must always be in ns, but the low bits may be junk.
   0 = 10^-(0*3) = 1s
   1 = 10^-(1*3) = 1ms (default for cloud systems)
   2 = 10^-(2*3) = 1us (default for 10G PTP with GPS)
   3 = 10^-(3*3) = 1ns
     XXX this would be cooler if it was replaced with values taken directly from ntpd describing time accuracy
* point  : 8 byte nanoseconds since unix epoch UTC overflows in year 2554
* crc    : 4 byte crc32 checksum of the event. only checked when reading from disk.
* length : 3 byte length of the event including 0xFF terminator. must come first so the crc can be used!
* EOE    : 1 byte = 0xFF that will never appear in utf-8
  together with length, the EOE provides redundant end of message detection
  to allow recovery from file corruption

XXXXXX what if everything outside the routing fields are never touched?

receive, parse, tokenize, serialize
pass to delivery layer and add the 12 byte fixed header
store message and allow retrieval by time only.
the sorting layer NEVER parses the message
make it the responsibility of the consumer to detect and parse multiple formats??????


don't reinvent the encoding, but still allow the tokenization client side???????

what if the


variable length fields
======================

* words  : 2 byte length + data. each word is 1 byte length + up to 256 bytes of characters
* line   : 2 byte length + up to 64KB of unstructured line text with the words replaced with (0xFE + varint)
* maps   : 1 byte count of the number of key/value pairs in the optional map to follow. 0 - 255
* data   : the key/value data: 0 to 2^8 * (2 + 2^8 + 2^8) ~ 128KB
    * 1 byte key_length
    * 1 byte value_length
    * up to 256 byte key
    * up to 256 byte value


variable length fields are 5 bytes to (5 + 2^16 + 2^16) + (2^8 * (2 + 2^8 + 2^8)) ~ 257KB

empty event size = 12 + 5 + 0 + 1 = 18 bytes
max event size = 12 + (5 + 2^16 + 2^16) + (2^8 * (2 + 2^8 + 2^8)) + 1 ~ 257KB


Why not use syslog's event format?
syslog usess a parsed ASCII format designed for interoperability. see rfc5424 section 6.
the main problems with this format are:
* problem: messages of unlimited size make time and space guarantees much more difficult to achieve
  solution: cap messages at 256KB
* problem: variable width 9 field date specification with ascii separator fields requires X??? operations to parse.
  solution: 8 byte uint64 requires 2 operations: uint64(event[9:17]) to parse.
    before unix epoch = out of scope for log processing
* problem: "priority = facility,severity" is difficult to use
  solution: separate routing from non-routing data:
    routing = 2 bit replication + 2 bit priority
    non-routing = "application name" up to 256 bytes
* problem: structured data is separated by ascii characters, so those characters have to be escaped
  solution: store all utf-8 strings with a length and never parse them


Why not allow typed numbers in "data" ?
  if numbers are distinguished, then they could be excluded from search terms. this would vastly reduce search terms.
  unfortunately, one of the huge search targets is various "ID" type numbers.
  so number values are being treated as characters in utf-8 strings like all other values and strings


Why not do whatever logstash does?
* logstash stores k/v data with lengths and routes by the k/v data
* event consumers may route by k/v data AFTER this system has guaranteed it's durability


Why not do whatever heka does?
* typed protocol buffer... allows efficient parsing of the whole message
  but I am explicitly NOT parsing most of the event = []byte received from the TCPCon.Read()
  full message parsing is done by event consumers
* uses an hmac... not a bad idea, but I'm hoping encryption will cover the intended use of the hmac
* has an int32 severity


How long does a word have to be to be worth picking out for compression and search?
* word gets replaced with 0xFE + varint representing it's spot in the "words" array
* 1 extra byte (word length) is needed to store the word in the words array
* word cost = 3 to 4 bytes more than leaving it in the line
* benefit = ~ 1/3 of the indexing process is now finished


Is there some standard for determining how much reliability is required for event delivery?
* I could not find an existing standard


Input Format:
* client must listen on localhost: udp/tcp 5465 (netops-broker in /etc/services)
* and accept a variety of formats:
  the internally used binary format
  json
  xml
  csv




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



lexing for tokens:

ok... I got a number
* could be an IP4... keep reading.. if it's between 0 and 255 followed by a dot, then keep going to see if it's an IP
* could be hex... keep reading... if it's a series of hex characters followed by anything else, then it's a hex string
* could be an integer keep reading.. if it's a bunch of numbers followed by punctuation or space, then it's a number
* could be the start of some kind of other alphanum sequence... read till punctuation or space and emit as token

ok... I got a letter
* could be a hex... see if it's in the hex range
* could be a natural language word... more letters till punctuation
* could be ???

ok... I got  - . _ / do I care about any of these????  maybe later


go/token ??
Pos and Token types

http://golang.org/src/pkg/go/scanner/scanner.go
isLetter
isDigit
digitVal
scanNumber

bufio.Scanner!!! http://jeremy.marzhillstudios.com/io/



http://golang.org/pkg/text/scanner/

l.accept()




   var txt = `{key1 = "\"value1\"\n" | key2 = { key3 = 10 } | key4 = {key5 = { key6 = value6}}}`
    var s scanner.Scanner
    s.Init(strings.NewReader(txt))
    var b []byte

loop:
    for {
        switch tok := s.Scan(); tok {
        case scanner.EOF:
            break loop
        case '|':
            b = append(b, ',')
        case '=':
            b = append(b, ':')
        case scanner.Ident:
            b = append(b, strconv.Quote(s.TokenText())...)
        default:
            b = append(b, s.TokenText()...)
        }
    }

    var m map[string]interface{}
    err := json.Unmarshal(b, &m)
    if err != nil {
        // handle error
    }

    fmt.Printf("%#v\n",m)













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

	ip_re := rubex.MustCompile(`\d+\.\d+\.\d+\.\d+`)

	for line = range lines {

		// remove the ip from line before looking for other tokens
		//		re.ReplaceAllString(src, repl string) string

		for _, c = range line {

			// XXXX not identifying IP address!!!!
			if ip_character(c) {
				ip += string(c)
			} else {
				if len(ip) >= 2 { // rfc1924: 0 = "::"
					if ip_object := net.ParseIP(ip); ip_object != nil {
						ips <- ip_object
					}

				}
				ip = null_string
				// continue
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
			fmt.Println("sending string =", s)
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
