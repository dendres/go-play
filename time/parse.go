package main

import (
	"fmt"
	"time"
)

func main() {

	// reference time: Mon Jan 2 15:04:05 -0700 MST 2006
	// reference time: 2006-01-02T15:04:05.000000-07:00
	format := "2006-01-02T15:04:05.000000-07:00"
	s := "2014-01-30T00:30:01.246899+00:00"
	time.Parse(format, s)
	fmt.Println(s)
}
