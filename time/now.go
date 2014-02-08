package main

import (
	// "bytes"
	// "encoding/base32"
	// "encoding/binary"
	"fmt"
	// "hash/crc32"
	"strconv"
	"time"
)

// print the current time in a number of serializations
func main() {
	// now := time.Now()
	now := time.Unix(9223372036854775807, 999999999)
	sec := now.Unix() // int64
	nsec := int64(now.Nanosecond())

	sec10 := strconv.FormatInt(sec, 10)
	nsec10 := strconv.FormatInt(nsec, 10)

	fmt.Println("standard:", now)
	fmt.Println("2x int64 as base10:", sec10, nsec10)

	sec16 := strconv.FormatInt(sec, 16)
	nsec16 := strconv.FormatInt(nsec, 16)
	fmt.Println("2x int64 as base16:", sec16, nsec16)

	sec32 := strconv.FormatInt(sec, 32)
	nsec32 := strconv.FormatInt(nsec, 32)
	fmt.Println("2x int64 as base32:", sec32, nsec32)

	sec36 := strconv.FormatInt(sec, 36)
	nsec36 := strconv.FormatInt(nsec, 36)
	fmt.Println("2x int64 as base36:", sec36, nsec36)

	// Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, mpeg-2, ...
	// ChecksumIEEE(data []byte) uint32

	// n := 26
	// testBytes := make([]byte, n)
	// for i := 0; i < n; i++ {
	// 	testBytes[i] = 'a' + byte(i%26)
	// }
	// fmt.Println("testBytes =", string(testBytes))
	// cs := int64(crc32.ChecksumIEEE(testBytes))
	// cs16 := strconv.FormatInt(cs, 16)
	// cs32 := strconv.FormatInt(cs, 32)
	// cs36 := strconv.FormatInt(cs, 36)
	// fmt.Println("ChecksumIEEE(testBytes) 10 =", cs)
	// fmt.Println("ChecksumIEEE(testBytes) 16 =", cs16)
	// fmt.Println("ChecksumIEEE(testBytes) 32 =", cs32)
	// fmt.Println("ChecksumIEEE(testBytes) 36 =", cs36)

	// serialize after combining sec, nsec, and crc32 onto a single byte slice???

	// out := new(bytes.Buffer)
	// lowercase := base32.NewEncoding("0123456789abcdefghijklmnopqrstuv")
	// encoder := base32.NewEncoder(lowercase, out)

	// input := []byte("foo\x00bar")
	// encoder.Write(input)
	// encoder.Close() // required to write the last character!

	// fmt.Println(out)

	// var out string
	// fmt.Println("appended =", strconv.AppendInt(out, int64(5), 10))

	// combined := new(bytes.Buffer)
	// c2 := strconv.AppendInt(combined, testint5, 10)

	// combined := make([]byte, 2)
	// binary.LittleEndian.PutUint16(combined, uint16(5))
	// binary.LittleEndian.PutInt64(combined, int64(5))

	//	binary.Write(combined, binary.LittleEndian, int64(5))
	//	binary.Write(combined, binary.LittleEndian, int32(777777))
	//	binary.Write(combined, binary.LittleEndian, int32(88))

	//	fmt.Println("str =", str)

	// binary.Read(combined, binary.LittleEndian, &c1)

	// RFC 4648 base32 ??

	// print whole buffer????
	// AppendInt(dst []byte, i int64, base int) []byte

}

// myfirstint, err := binary.ReadVarint(buf)
// anotherint, err := binary.ReadVarint(buf)
