package main

import (
	"fmt"
	"time"
)

/*
how many bytes are required to identify the UTC day in a json string?
what format should be used to allow easy consumption of the time by standard javascript time functions?

var js_date = new Date(milliseconds_since_utc_unix_epoch)
var js_utc_date_string = js_date.toUTCString()

would be great to append a fixed number of zeros to the millisecond time, or bit-shift
  XXX can't bit-shift because of javascript's crappy numbers

best bet seems to be to append a fixed number of zeros?

so print out some day stamps to see the pattern:
  it appears as though the last 2 digits of uint64(start_of_a_day.Unix()) are always zero
  can't prove it
  it's only 2 bytes

might as well send uint64(t.Unix()) as json to the browser

determine limits ????

*/

func main() {
	format := "2006-01-02T15:04:05.000000-07:00"

	times := []string{
		"1970-01-01T00:00:00.000000+00:00",
		//		"2014-01-30T00:30:01.246899+00:00",
	}

	for m := 1; m < 13; m++ {
		for d := 0; d < 31; d++ {
			x := fmt.Sprintf("2014-%02d-%02dT00:00:00.000000+00:00", m, d)
			times = append(times, x)
		}
	}

	for i, time_string := range times {
		t, err := time.Parse(format, time_string)
		if err != nil {
			fmt.Println("error parsing timestamp", err)
		}
		time_int := uint64(t.Unix())
		fmt.Println(i, time_string, t, time_int)
	}
}
