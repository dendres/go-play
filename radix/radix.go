package radix

import "fmt"

// SortWords does a radix sort by the the "bytes" most significant bytes in each word.
// The input slice is modified in place.
func SortWords(bytes int, input [][]byte) error {
	if bytes < 1 {
		return fmt.Errorf("cannot sort by no words: bytes = %d", bytes)
	}
	inputlen := len(input)
	output := make([][]byte, inputlen)

	var (
		count [256]int
		word  []byte // input[i]
		elem  byte   // input[j][i]
		k, i  int
	)

	// don't try to sort empty or single element arrays
	if inputlen < 2 {
		return nil
	}

	// XXX this should be replaced by [][6]byte input specification if possible... but not possible now.
	for i, word = range input {
		if len(word) < 1 {
			return fmt.Errorf("cannot sort zero length words")
		}
		if bytes > len(word) {
			return fmt.Errorf("bytes = %d is too large for word input[%d] = %v", bytes, i, word)
		}
	}

	// sort one byte at a time starting at the least significant byte
	for k = bytes - 1; k >= 0; k-- {

		// zero count and index
		for i, _ = range count {
			count[i] = 0
		}

		// http://en.wikipedia.org/wiki/Counting_sort
		// count the occurance of element k in each word
		for i, word = range input {
			count[word[k]]++
		}

		// where count was: 01000001000000100000
		// count will be:   01111112222222333333
		for i = 1; i < 256; i++ {
			count[i] += count[i-1]
		}

		fmt.Printf("k = %d, count = %v\n", k, count)

		// count[word[k]] = 01111112222222[3]33333
		// this is one too large. need to move back one spot. see the -1 below:
		//for i, word = range input {
		// fmt.Println("elem =", word[i], ", count[elem] =", count[elem])
		//	output[count[word[k]]-1] = word
		//}

		// Build the output array
		for i = inputlen - 1; i >= 0; i-- {
			word = input[i]
			elem = word[k]
			output[count[elem]-1] = word
			count[elem]--
		}

		// overwrite input with output
		for i, word = range output {
			input[i] = word
		}

		fmt.Printf("k = %d, output = %v\n", k, output)
	}
	return nil
}

// fmt.Println("for k =", k, "and word[k] -1 =", elem, "the count[elem] =", count[elem], "assigning word =", word, "to output index at this c=", c)

/*
Reference:
* http://rosettacode.org/wiki/Sorting_algorithms/Radix_sort#Go
* http://www.cs.princeton.edu/~rs/AlgsDS07/18RadixSort.pdf
* http://en.wikipedia.org/wiki/Counting_sort
* http://pastebin.com/f1222a10
* http://www.geeksforgeeks.org/radix-sort/
*/
