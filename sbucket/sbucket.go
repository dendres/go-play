package sbucket

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"
)

// allow a flexible specification that forces "reasonable" base10 buckets
// return , ns, us, ms, s, 10s, and 5m string "buckets" from time.Now()
// negative means 10^-X and truncates digits off nanosecond time
// positive means number of seconds in the bucket
// 0 < x <= 1: the bucket size will be forced to divide evenly into minute
// x > 1: forced to multiple of min
// x > 1: forced to a multiple of hour
func Stamp(ti time.Time, size int) (string, int) {
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
		out = strconv.FormatInt(sec-(sec%int64(size)), 10) // round down to a multiple of size

		if size%10 == 0 {
			out = out[:len(out)-1] // strip digit
		}
		if size%100 == 0 {
			out = out[:len(out)-1] // strip digit
		}
		if size%1000 == 0 {
			out = out[:len(out)-1] // strip digit
		}
	default:
		for size%3600 != 0 { // force multiple of hour
			size++
		}
		out = strconv.FormatInt(sec-(sec%int64(size)), 10) // round down to a multiple of size
		// find what power of 10 size is.
		// strip that many digits
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
// ((2^64 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000000 = 292 years +- unix epoch in ns
// ((2^64 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000 = 292K years +- unix epoch in us
// ((2^64 / 2) - 1) / 365 / 24 / 60 / 60 / 1000 = 292M years +- unix epoch in ms
// convert time.Now() into each of these formats

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
