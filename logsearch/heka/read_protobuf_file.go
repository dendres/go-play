package main

import (
	"fmt"
	"github.com/mozilla-services/heka/message"
)

func main() {
	// read file
	// figure out how heka delimits the protobuf messages?????
	// undelimit them

	// make new messages from each []byte

	m := message.Message{}
	// err = proto.Unmarshal(filebuf, m)

	fmt.Println(m)
}
