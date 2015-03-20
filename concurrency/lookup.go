package main

import (
	"fmt"
	"time"
)

// concurrent dns lookup. first response wins
// XXX later add caching


var servers = []string{
	"4.2.2.2",
	"8.8.8.8"
}


func lookup(query string) string {
    ch := make(chan Result, len(servers))  // buffered
    for _, conn := range conns {
        go func(c Conn) {
            ch <- c.DoQuery(query):
        }(conn)
    }
    return <-ch
}

// teardown of late finishers????

func main() {
	fmt.Println("hello")
}
