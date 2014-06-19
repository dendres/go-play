package main

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
)

/*
for compression only... ignore search for now
find the longest repeating []byte in the file

XXX later, events will be stored with delimiter of some kind. newline or offset delimited
    at that time, each event's size will be available. and a buffer can have the whole event
    so won't have to worry about things getting cut off by the buffering
  so just read in the whole file for now as one iteration

http://en.wikipedia.org/wiki/Longest_repeated_substring_problem
make suffix tree, then find the highest node with at least 2 descendants

http://www.allisons.org/ll/AlgDS/Tree/Suffix/
http://en.wikipedia.org/wiki/Suffix_tree
http://en.wikipedia.org/wiki/Suffix_array

http://stackoverflow.com/questions/9452701/ukkonens-suffix-tree-algorithm-in-plain-english

https://www.cs.helsinki.fi/u/ukkonen/SuffixT1withFigs.pdf

http://golang.org/pkg/index/suffixarray/

*/

func main() {

	buf, err := ioutil.ReadFile("test1.log")
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	index := suffixarray.New(buf)

	fmt.Println(index)

}
