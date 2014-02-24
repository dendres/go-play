package main

/*
replace a given string with an encoded byte slice

how to substitute tokens??????????
http://en.wikipedia.org/wiki/DEFLATE
  Compression is achieved through two steps
    1. The matching and replacement of duplicate strings with pointers.
    2. Replacing symbols with new, weighted symbols based on frequency of use.




attempt1:
  substituted_token = type_byte + value_bytes where each type has a fixed byte length
  encoded_size = token_table_size + (type_byte * token_count) + all_the_token_lengths_added_together????

  XXXXX too expensive to have a full byte of metadata for every rune/token


attempt2:
  1 byte metadata per rune/token
  1 bit type (1 = token), 6bit byte count, bytes
  max number of unicode bytes = 3 ? verify!!!... so allow 2 bits for byte count
  min single bit rune: 0 00 00000 => U+0000
  max single bit rune: 0 11 11111 => U+001F ... just the control characters... so alphanum looks like:
  most alphanum would need 2 bytes!!
  token_index 2 bit byte count.. same as above 0-63 in one byte, up to 247 billion in 5 bytes

  XXXXX too expensive to have a full byte of metadata for every rune/token


attempt3:
  all values are variable byte
  only strings are stored in the token table. no separate type or table for numbers

table_bitset:
  0 = rune, 1 = token
  one bitset entry per variable-length rune/token... NOT one per/byte!!!


writing characters:
  1. append data bytes to byte slice
  2. append type bit to bitset

reading characters:
  1. keep count
  2. bitset[count] determines type token/rune
  3. rune_reader([]byte) (string character, []byte remaining)... or maybe just return offsets or something like scanner!!
     token_reader([]byte) (uint32 token_index, []byte remaining)


http://en.wikipedia.org/wiki/Variable-length_quantity
writing:
*


bitset operations???????????

bitset append:
  ????

bitset output:
  return []byte

bitset read bit:
  return bool



======= listA =======

properties:
* a list of 256 x 32 byte values stored in binary serialization as []byte
* data overhead = only a fixed few bytes overhead to keep the total length.
* single seek to get a value
* single seek to write a value
* length limits established at init.
* id is a variable number of bytes
* does NOT deduplicate
* ordered by insert order ONLY. no sorting.

format:
* type(bytes_in_size):
* bytes_in_size total size of the list + big long unicode string
* id: 5 bit length + varint offset from:
   id 0 of length 3 should fit in 1 byte!

  00000 000  00000000
        0 always
  00000  11 = first 3 tokens in 1 byte

if there was a better way to encode the length, then the first 128 tokens fit in 1 byte????

methods:
* split_id(id []byte) length, offset:
  uint64 = varint(id)
* Read(id): offset, length = split_id(id), seek to offset, read length
* Append(string) = seek to size, write string, size += len(string)

================= listB ====================

must hold the worst case 2K tokens found in a 64KB message

some notes from looking at Event serialization:
* Empty serialized event = 1 8 1 1 4 4 4 = 23 bytes
* 1 unencoded character event: 1 8 1 4 1 12 4 1 4 0 4 2 ~ 42 bytes
* worst_case_uncompressed= 32 header + marks 2^13 + tokens 2^16 +line 2^16 ~ 136K
* normally, as tokens gets larger, marks and lineget smaller. marks getsslower 1/8 speed



from: http://www.zehnet.de/2005/02/12/unicode-utf-8-tutorial/
"The bytes 0xFE (11111110) and 0xFF (11111111) are never used in the UTF-8 encoding."

validated: http://en.wikipedia.org/wiki/UTF-8
"The bytes 0xFE and 0xFF do not appear, so a valid UTF-8 stream never matches the UTF-16 byte order mark and thus cannot be confused with it. The absence of 0xFF (0377) also eliminates the need to escape this byte in Telnet (and FTP control connection)."

properties:
* []byte up to 64KB of up to 256 Byte 0xFF separated unicode strings
* 2 byte length
* id = varint offset from the beginning of the []byte
* data overhead = 2 byte + 1 byte/item

id properties of varint:
0000000 = 0
0111111 = 127
1000001 00000000 -> 000000 10000000 -> 128
1000001 01111111 -> 000000 11111111 -> 255
1111111 01111111 -> 111111 11111111 -> 16383
1000001 10000000 00000000 -> 1000000 00000000 -> 16384
so id's 0 to 127 are 1 byte
and id's 128 to 16383 are 2 byte

methods:
* Init(x,y): length = 0, the empty []byte
  x = fixed number of bytes used to store length
  y = fixed max length of each item = 2^y bytes
* Append(string):
  error if len(string) > 256 Bytes
  read length
  new_length = length + len(string) + 1
  error if new_length >= 4GB
  set id = length, seek to length, write string, write 0xFF
  length = new_length
  return id
* Read(id):
  error if id >= 4GB
  seek to offset = id
  read bytes till 0xFF
  return string


=========== full "line" serialized format =============

* bitset mark encoding for each byte.
* listB of tokens






* the array of Events


XXXX do this and test!!!


================ building the token table incrementally =============

incrementally record the frequency data:

frequency... fixed
word_length
word



??????????????????

file contains: ??????????????

make token_id relative to the timestamp at the start of the bucket??????
 map[token_string] -> []event_id
 map[event_id] -> []token_string




=================== getting the word tokens for compression and search ==================

set limits on word size ????



https://github.com/gyuholee/goling  bunch of language processing. extract numbers, segment

segmenting, stemming etc...
norvig probabilistic word segmenter: https://github.com/llimllib/segment
another segmenter: https://github.com/gyuholee/goling/blob/master/segmentation.go

n-gram based text categorization (guess language): http://godoc.org/github.com/pebbe/textcat


=========== performing search ======
http://swtch.com/~rsc/regexp/regexp4.html
https://code.google.com/p/codesearch/source/browse/index/read.go

given regex(/Google.*Search/), produce query of 3-grams => Goo AND oog AND ogl AND gle AND Sea AND ear AND arc AND rch

build a probabilistic autocomplete based on 3-grams or other??????


*/
func Encode(s string) ([]byte, []byte)

// parse_rune takes a byte slice and returns a character and a byte slice
// gets called only when the start of the byte slice is a known unicode character
func parse_rune([]byte) {
}

// parse_token takes a byte slice and returns a uint64 and a byte slice
func parse_token([]byte) {
}

// XXXXXXXXXXXX i find it very difficult to believe that this could work.

// look at the unicode parser and see how it figures out weather or not to read the next byte!!!!

/*
http://www.zehnet.de/2005/02/12/unicode-utf-8-tutorial/

 utf-8 is one byte: U+0000 (00000000) to U+007F (01111111)
 All Unicode characters > U+007F are encoded as a sequence of several bytes, each of which has the most significant bit set. Therefore, no ASCII byte (0×00-0×7F) can appear as part of any other character.

The first byte of a multibyte sequence that represents a non-ASCII character is always in the range 0xC0 (11000000) to 0xFD (11111101) and it indicates how many bytes follow for this character. All further bytes in a multibyte sequence are in the range 0×80 (10000000) to 0xBF (10111111). This allows easy resynchronization and makes the encoding stateless and robust against missing bytes.

The bytes 0xFE (11111110) and 0xFF (11111111) are never used in the UTF-8 encoding

Unicode character number (decimal)bit sequence
U-0        – U-127       :0xxxxxxx (ASCII characters)
U-128      – U-2047      :110xxxxx 10xxxxxx
U-2048     – U-65535     :1110xxxx 10xxxxxx 10xxxxxx
U-65536    – U-2097151   :11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
U-2097152  – U-67108863  :111110xx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx
U-67108864 – U-2147483647:1111110x 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx


ok... I'm convinced. the first "continuation" bit followed for bits representing number of additional bytes.








*/

func main() {

	// iterate over all possible values of 4 bytes
	for i := 0; i < 4294967296; i++ {
		a := i
		a = a - i + i
	}

}
