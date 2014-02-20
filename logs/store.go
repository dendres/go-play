package store

import (
	"fmt"
	"os"
)

/*
interface files that store events
* events already have timestamp and crc32
* 2 byte size. 64K max element

methods:
* append([]byte)
* read(offset)
* scan()
* reverse()
* recover()
* recovergzip()
*/

/*
s := NewStore
event := s.Next()



*/

/*
store.Append(b []byte)
* offset = os.Seek(0, 2) // EOF
* os.Write(len(b))
* os.Write([]byte)
* return offset // new id for the event if needed
*/

/*
store.Read(id int64)
* os.Seek(id, 0) // seek to the offset given by id
* l := make([]byte, 2)
* os.Read(l)
* event_length := int(l)
* e := make([]byte, event_length)
* os.Seek(2,1)
* os.Read(e)
* return e

*/

/*
read the whole thing out into a stream of events????
read l, seek 2, read len(l), emit event
repeat
*/

/*
repair corrupt file?
 raw events... use crc32
 gzip... use the crc and length in the gzip

read and discard bytes until gzip magic byte. backup 2 and read length.?????

*/

/*
http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file

*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
