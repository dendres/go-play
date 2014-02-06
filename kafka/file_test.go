package fileassumptions

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

/*
test assumptions about file operations... open, close, read, seek, etc...
*/
func TestOpen(t *testing.T) {

	now := time.Now()
	fmt.Println("now =", now)

	s := now.Nanosecond()
	fmt.Println("ns =", ns)

	ns := now.Nanosecond()
	fmt.Println("ns =", ns)

	nss := strconv.Itoa(ns)
	fmt.Println("nss =", nss)

	//t.Fatal("ssssshhhhhhiiiiitttt")

}
