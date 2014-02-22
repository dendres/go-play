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

assume log line strings from 80 to 64k bytes
string content contains microsecond timestamp and hostname
how much entropy is required to prevent collission based on message content?

http://preshing.com/20110504/hash-collision-probabilities
a 32 bit hash function should be good up to about 1K messages, but pretty bad after 2K
a 64 bit hash function should be good up to 1M messages, but pretty bad after 5M


pick the hash function with the lowest collision rate on 80 byte - 64KB punctuated english with occasional hex strings
if there is a tie, pick the faster one.





http://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
  inputs were very short strings...

start with crc32, collect a data set to use for benchmarks.
then mesure speed and collision rate on that data set


https://github.com/spaolacci/murmur3
https://github.com/reusee/mmh3
http://golang.org/pkg/hash/crc32/
https://github.com/dgryski/dgohash


*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
