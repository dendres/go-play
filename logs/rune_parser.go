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


=========== full serialized format =============

Message:
* the array of token strings
* the array of Events

Event:
* point uint64: 10^-8 since unix epoch loops in year 7819
* shn string:   short host name
* app string:   overloaded = app_name || process_name || process_id
* line []byte:  encoded line with app, shn, and point decoded


XXXX do this and test!!!


================ building the token table incrementally =============

incrementally record the frequency data:

frequency... fixed
word_length
word



??????????????????

file contains: ??????????????




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
