package chat

import (
	"log"
	"net"
	"testing"
)

func isTimeout(err error) bool {
	e, ok := err.(Error)
	return ok && e.Timeout()
}

// BUG(done) find a tcp connection library... Dial???
// start the chat go routine
// dial the port
// send text
// check the response text
// close the dial connection
// stop the chat go routine????
func TestXXXX(t *testing.T) {
	for i, test := range tests {
		size := Size(test.in)
		if size != test.out {
			t.Errorf("#%d: Size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}
