package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

/*
find ways to make go a useful replacement for string searching on the command line

takes a regex argument
takes stdin
prints lines matching regex
exit 0 on success, non-zero on error
*/

func main() {
	regex_string := ""
	args := os.Args
	if len(args) > 1 {
		regex_string = os.Args[1]
	}

	r, err := regexp.Compile(regex_string)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			fmt.Println(line)
		}
	}

}
