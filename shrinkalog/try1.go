package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

/*
256 item list of strings/tokens that represent the most repetition outside the compression block size

list of (offset, token_number) that describe how to re-inflate to the original file

crc32 of the original file
*/

// sort a map's keys in descending order of its values.

// https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func main() {

	// need to be able to sort by the total number of bytes represented by each token,count
	// XXX remember to test the case of a GB log file of hex sha512's!!!
	tokens := make(map[string]int)

	input, err := os.Open("test1.log")
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		token := scanner.Text()
		if len(token) > 5 {
			// add the length of the token each time to get the total number of bytes taken by this token
			tokens[token] += len(token)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

	// iterate the tokens by most bytes replaced
	s := sortedKeys(tokens)
	for i := 0; i < 10; i++ {
		k := s[i]
		v := tokens[k]
		fmt.Println(k, v)
	}

	// now do the replacement???????
	// for each token to be replaced
	//   scan the whole thing, find each, delete, save the offset in a slice

}
