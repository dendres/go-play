package sbucket

import (
	"fmt"
	"github.com/dendres/go-play/encoding/lp32"
	"time"

	"math"
	"strconv"
	"testing"
)

/*
storing years:
35 bits stores 2^35-1 seconds ~ 1089.54 years
epoch is 1970, so the stamp loops during the year 1970 + 1089 = 3059
there will be no record of this code by Jan 1 3059
*/

/*
storing fractional second:

takes 30 bits to store nanoseconds: 2^29 - 999,999,999 = -463,129,087

but... 35 bits for second + 30 bits for fraction = annoying 65 bits
so take off 5 bits and you get:

25 bits stores the number of ~ 30ns intervals that divide the second

translate 0.999999 into the number of 29ns intervals since the start of the second?
  0.999999 * 2^29 -1

translate 25bit integer back into a fractional second:
  0 / 2^25-1 = 0
  1 / 2^25-1 ~ 29.802323ns
  2 / 2^25-1 ~ 59.604646ns
  2^25-1/2^25-1 = 1 second

translate 20bit integer back into fractional second:
  0 / 2^20-1 = 0
  1 / 2^20-1 = 953.675225ns
  2 / 2^20-1 = 1.90735045us
  2^20-1/2^20-1 = 1 second


so 25bits stores ~ 30 nanosecond granularity
and 20bits stores ~ microsecond granularity

I'm sticking to 25 bits for for now
*/

/*
dividing up directories by date:

what gets expensive as file count grows?
* any use of file globs
* file open and delete times

how many entries in a directory before the cost increase?
* I could not find a cost calculation of any kind... too many variables?
* anywhere around 10k-20k should be reasonable?

time intervals for the base32, 7 character timestamp:

6 = 1073741824 ~ 34.05 years
5 = 33554432 ~ 1.1 years
4 = 1048576 ~ 12.1 days
3 = 32768 ~ 9.1 hours
2 = 1024 ~ 17.1 min
1 = 32 = 32 sec
0 = 1 = 1 sec

/base_dir/s/sssssss stores 34 year interval
...
/base_dir/sssssss/s stores 32 second intervals

time intervals for the base32 5 character fractional timestamp:

4 = 1048576 ~ 31.25ms
3 = 32768 ~ 976.5us
2 = 1024 ~ 30.5us
1 = 32 = 953.7ns
0 = 1 = 29.802323ns
*/

// Stamp1 creates a 7 character, fixed width, base32 timestamp covering 1970 ~ 3059 in seconds.
func Stamp1(ti time.Time) (string, error) {
	sec := uint64(ti.Unix())
	return Enc(sec, 7)
}

// Read1 ???

// Stamp2 creates a 5 character, fixed width, base32 timestamp covering 1 second in 30ns intervals.
func Stamp2(ti time.Time) string {
	nsec := uint64(ti.Nanosecond())
	bs := make([]byte, 5)

	// encode every 5 bits as a character in [0-9a-v]
	for i := 4; i >= 0; i-- {
		bs[i] = base32map[nsec%32]
		nsec = nsec >> 5
	}

	return string(bs)
}

// XXX Read2 ???

// allow a flexible specification that forces "reasonable" base10 buckets
// return , ns, us, ms, s, 10s, and 5m string "buckets" from time.Now()
// negative means 10^-X and truncates digits off nanosecond time
// positive means number of seconds in the bucket
// 0 < x <= 1: the bucket size will be forced to divide evenly into minute
// x > 1: forced to multiple of min
// x > 60: forced to a multiple of hour
// x > 3600: force multiple of a day
func TenStamp(ti time.Time, size int) (string, int) {
	out := "fail"
	sec := ti.Unix() // int64

	switch {
	case size < 0:
		if size < -9 { // force nanosecond minimum
			size = -9
		}

		nsec := int64(ti.Nanosecond())
		pow := int64(math.Pow10(-1 * size)) // like 100 or 1000

		// shift over the second to make room for the fractional second
		// size the fractional second to fit and add them together
		x := (sec * pow) + (nsec / (int64(1000000000) / pow))

		out = strconv.FormatInt(x, 10)
	case size == 0:
		out = strconv.FormatInt(sec, 10)
	case size <= 60:
		for 60%size != 0 { // force even division of 60
			size++
		}
		out = strconv.FormatInt(sec-(sec%int64(size)), 10) // round down to a multiple of size
		if size%10 == 0 {
			out = out[:len(out)-1] // strip digit
		}
	case size <= 3600:
		for size%60 != 0 { // force multiple of minute
			size++
		}
		for 3600%size != 0 { // force divisor of hour
			size++
		}
		out = strconv.FormatInt(sec-(sec%int64(size)), 10)
		if size%10 == 0 {
			out = out[:len(out)-1]
		}
		if size%100 == 0 {
			out = out[:len(out)-1]
		}
	case size <= 86400:
		for size%36000 != 0 {
			size++
		}
		for 86400%size != 0 {
			size++
		}
		out = strconv.FormatInt(sec-(sec%int64(size)), 10)
		if size%10 == 0 {
			out = out[:len(out)-1]
		}
		if size%100 == 0 {
			out = out[:len(out)-1]
		}
		if size%1000 == 0 {
			out = out[:len(out)-1]
		}
	default:
		for size%86400 != 0 { // force multiple of day
			size++
		}
		out = strconv.FormatInt(sec-(sec%int64(size)), 10)
		if size%10 == 0 {
			out = out[:len(out)-1]
		}
	}

	return out, size
}

// var x Sb  // returns Sb
// x := new(Sb) // returns *Sb
// x := Sb(x: 5, y: 4, z: 3)  // not pointer

// x := Sb(5,4,3)
// can pass &x to a function if a pointer is required

//
// func Start(t time, power int) time {

// func End(t time, power) {

// return , ns, us, ms, s, 10s, and 5m string "buckets" from time.Now()

func Open(t *testing.T) {

	now := time.Now()
	sec := now.Unix()
	nsec := now.Nanosecond()

	fmt.Println("now =", now)
	fmt.Println("sec =", sec)
	fmt.Println("nsec =", nsec)

	// 1 second string: convert int to base10 string
	dir1 := "/path/" + strconv.FormatInt(sec, 10) + "/more/path"
	fmt.Println("second   =  ", dir1)

	// 10 sec buckets: truncate 1 digit from the second string
	ss := strconv.FormatInt(sec, 10)
	dir2 := "/path/" + ss[:len(ss)-1] + "/more/path"
	fmt.Println("10 second = ", dir2)

	// 5 min buckets:
	// s - s % 300
	// bucket sizes that are multiples of 10^X can truncate X digits
	for i := 0; i < 10; i++ {
		s := sec + int64(i*33)
		fm := strconv.FormatInt(s-s%300, 10)
		dir3 := "/path/" + fm[:len(fm)-2] + "/more/path"
		fmt.Println("5 min = ", dir3)
	}

	// ms buckets: (sec * 1000) + (nsec / 1000000)
	for i := 0; i < 1; i++ {
		now1 := time.Now()
		sec1 := now1.Unix()
		nsec1 := now1.Nanosecond()
		ms := (sec1 * 1000) + (int64(nsec1) / 1000000)
		mss := strconv.FormatInt(ms, 10)
		dir3 := "/path/" + mss + "/more/path"
		fmt.Println("sec =", sec1, "ns =", nsec1, "ms =", dir3)
	}

	// // us buckets: (sec * 1000000) + (nsec / 1000)
	// // ns buckets = strconv.FormatInt(time.Now().UnixNano(), 10)

	// // now... store the buckets as FormatInt(ms, 36) and figure out how many digits to keep???
	// b36 := strconv.FormatInt(sec, 36)
	// back, _ := strconv.ParseInt(b36, 36, 64)
	// fmt.Println("started with sec = ", sec, "b36 =", b36, "back =", back)
	// if back != sec {
	// 	t.Fail()
	// }

	// // base 36... strip digits on intervals of 36^n... not really any convenient intervals
	// // are any of 36^n convenient for segmenting time?
	// // 36^2 = 21m, 36s
	// // 36^3 = 12.96 hours
	// // 36^4 = 19.44 days
	// // so no... this is not really a useful encoding... so stick to base10 for readability!
	// start := time.Date(2014, time.February, 6, 0, 0, 0, 0, time.UTC)
	// for i := 0; i < 1; i++ {
	// 	m := time.Duration(i*1679616) * time.Second
	// 	n := start.Add(m)
	// 	fmt.Println("bucket =", n, "sec =", n.Unix(), "b36 =", strconv.FormatInt(n.Unix(), 36))
	// }

	// base30 might be useful!!!

}

// determine a reasonable precision time that can be stored in a int64
// ((2^63 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000000 = 146 years +- unix epoch in ns
// ((2^63 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000 = 146K years +- unix epoch in us
// ((2^63 / 2) - 1) / 365 / 24 / 60 / 60 / 1000 = 146M years +- unix epoch in ms
// convert time.Now() into each of these formats

// uint64 is ((2^64 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000000 = 292 years

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64. Note that this
// means the result of calling UnixNano on the zero Time is undefined.

// secondsPerMinute = 60
// secondsPerHour   = 60 * 60
// secondsPerDay    = 24 * secondsPerHour
// secondsPerWeek   = 7 * secondsPerDay

// unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
// internalToUnix int64 = -unixToInternal

// Unix() int64 = t.sec + internalToUnix

// UnixNano() int64 = (t.sec+internalToUnix)*1e9 + int64(t.nsec)

//    .Bytes() ???

// buf := bytes.NewBuffer(b) // []byte
// myfirstint, err := binary.ReadVarint(buf)
// anotherint, err := binary.ReadVarint(buf)
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
