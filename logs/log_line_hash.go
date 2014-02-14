package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
we currently produce 3Mm/h or 833m/s
target a max of 36Mm/h or 10Km/s

assume log line strings from 80 to 100k bytes
string content contains microsecond timestamp and hostname
how much entropy is required to prevent collission based on message content?

http://preshing.com/20110504/hash-collision-probabilities
a 32 bit hash function should be good up to about 1K messages, but pretty bad after 2K
a 64 bit hash function should be good up to 1M messages, but pretty bad after 5M

How to divide up the directory tree so that there are less than 1k entries in each directory?

5_min_stamp:
* 3M messages in a 5min interval should be covered by a 64 bit hash
* this model breaks down at around 5Mm/5min or 16Km/s

1_min_stamp:
* 64 bit hash breaks down at 5Mm/1min = 83Km/s

1_sec_stamp:
* 64 bit hash breaks down as 5Mm/s

1_us_stamp:
* 10Km/s = 0.01m/us
* 1Mm/s = 1m/us
* a 32bit hash should be good up to 1Km/us

(2**32 - 1).to_s(10) = "4294967295"
(2**32 - 1).to_s(16) = "ffffffff"
(2**32 - 1).to_s(36) = "1z141z3"


best hash function for this purpose available in go?
http://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
https://github.com/spaolacci/murmur3
https://github.com/reusee/mmh3
http://golang.org/pkg/hash/crc32/

https://github.com/dgryski/dgohash


*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
