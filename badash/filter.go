/*
package badhash implements an oversized, sparse file based table
given a uint64 key < 40 bytes long
seek to key and read/write the single bit there.

This is massively oversized, but requires a consistent single disk seek per operation

*/
package filter

import (
	"fmt"
	"os/file"
)


type Filter struct {
	// handle to the filter file
	file *os.File

	// the number of bits in the filter key. 41
	size uint8
}


/*

* makes file
* writes fixed width hader field with key_length
* other header fields????
*/
func Open(size int64, path\)


/*
filter.Set(key int64) error
* 0 < key < (1<<40) // never make a 1TB bit array.
* offset = len(header) + key
* file.Seek(offset, 0)
* b := []byte{0}
* file.Read(b) // read a whole byte, but only modify 1 bit.
* b = b[0] | 128
* file.Write(b)
* return error on disk error
1 seek, 1 read, 1 write
*/

/*
filter.Get(key int64) bool
* validate key
* offset = len(header) + key
* b := []byte{0}
* ReadAt(b, offset)
* if b[0] & 128 == 1 return true
* return false
1 seek 1 read
*/

/*
func (f *filter) GetSet(key int64) (found bool)
* validate key
* offset = len(header) + key
* b := []byte{0}
* file.Seek(offset, 0)
* file.Read(b) // read a whole byte, but only modify 1 bit.
* if b[0] & 128 == 1 return true
* b = b[0] | 128 // set the 1 bit
* file.Write(b)


* ReadAt(b, offset)
* if b[0] & 128 == 1 return true
* return false




*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
