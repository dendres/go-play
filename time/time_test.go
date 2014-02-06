package fileassumptions

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

/*
timestamps

efficiently produce, ns, us, ms, s, 10s, and 5m timestamps from time.Now()

time.Round(duration)





*/
func TestOpen(t *testing.T) {

	// determine a reasonable precision time that can be stored in a int64
	// ((2^64 / 2) - 1) / 365 / 24 / 60 / 60 / 1000000 = 292471.208 years +- unix epoch in us
	// convert time.Now() into this format????

	now := time.Now()
	fmt.Println("now =", now)

	// get Second
	sec := now.Unix()
	fmt.Println("sec =", sec)

	// get Nanosecond
	nsec := now.Nanosecond()
	fmt.Println("nsec =", nsec)

	// bit shift Second
	x := sec << 20
	fmt.Println("sec shifted up 20 =", x)

	// get 20 bits worth of microseconds added to the shifted Second ????

	nss := strconv.Itoa(nsec)
	fmt.Println("nss =", nss)

	//t.Fatal("ssssshhhhhhiiiiitttt")

}
