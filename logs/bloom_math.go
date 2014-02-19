package main

import (
	"fmt"
	"math"
)

/*
bloomkeybits determines how many bits wide the key should be given event count
m = bit count of the bloom filter array
k = 2 hash functions
n = events
p = false positive probability
p = (1- ( 1 - (1/m)  )^(kn) )^k

Sticking to a fixed number of hash functions k=2 (time, crc)

m = 2^key_length
key_length = math.Log2(m)

The ideal false positive probability is 1/n... meaning one message false positive total.
A fixed rate of 0.001% is a bit more realistic


*/
func bloomkeybits(item_count uint64) int64 {
	n := float64(item_count)
	p := 0.00001
	a := 1 - math.Pow(p, 0.5)
	b := 1 - math.Pow(a, 1/(2*n))
	m := 1 / b
	bits := int64(math.Log2(m))
	return bits
}

func main() {
	for i := 0; i < 32; i++ {
		n := uint64(1) << uint64(i)
		bits := bloomkeybits(n)
		fmt.Println("n =", n, "key bits =", bits)
	}
}

/*

rates:
* current: 2M events, 30 bits, 128MB, 20 false positives
* target: 500M events, 38 bits, 32GB, 5K false positives


n = 32768 m = 20691518 key bits = 24
n = 65536 m = 41383035 key bits = 25
n = 131072 m = 82766070 key bits = 26
n = 262144 m = 165532141 key bits = 27
n = 524288 m = 331064277 key bits = 28
n = 1048576 m = 662128579 key bits = 29
n = 2097152 m = 1324257061 key bits = 30
n = 4194304 m = 2648514122 key bits = 31
n = 8388608 m = 5297028245 key bits = 32
n = 16777216 m = 10594062721 key bits = 33
n = 33554432 m = 21188125443 key bits = 34
n = 67108864 m = 42376250886 key bits = 35
n = 134217728 m = 84752103039 key bits = 36
n = 268435456 m = 169505801022 key bits = 37
n = 536870912 m = 339011602045 key bits = 38
n = 1073741824 m = 677997685716 key bits = 39
n = 2147483648 m = 1356097448771 key bits = 40
*/
