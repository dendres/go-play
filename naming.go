package main

import "fmt"

/* https://blog.golang.org/package-names

If you cannot come up with a package name
that's a meaningful prefix for the package's contents,
the package abstraction boundary may be wrong.
Write code that uses your package as a client would,
and restructure your packages if the result seems poor.
This approach will yield packages that are easier for clients
to understand and for the package developers to maintain.


package names:
 - are short and clear
 - lower case
 - NO under_scores or mixedCaps
 - often simple nouns like (time, list, http)


Abbreviate judiciously.
Package names may be abbreviated when the abbreviation is familiar to the programmer.
Widely-used packages often have compressed names:
 - strconv (string conversion)
 - syscall (system call)
 - fmt (formatted I/O)

Directory names:
 runtime/pproff
 net/http/pproff


Don't steal good names from the user.
Avoid giving a package a name that is commonly used in client code.
For example, the buffered I/O package is called bufio, not buf, since buf is a good variable name for a buffer.

*/

/*
client code uses the package name together with exported functions, so...
 - avoid stutter when possible: http.Server NOT http.HttpServer

simplify function names:

  New: When do I use New vs. when do I include the Type in the function name:
    list.New()                returns a *list.List
    time.ParseDuration("10s") returns a time.Duration
    time.Since(start)         returns a time.Duration
    time.NewTicker(d)         returns a *time.Ticker
    time.NewTimer(d)          returns a *time.Timer

The same Type name can be used in different packages:
jpeg.Reader
bufio.Reader
csv.Reader


*/

/* Identifying bad package names:

 - is it meaningless?  util, common, misc
 - are all the api's in a single package?  api, types, interfaces

if so, break it up

*/

func main() {
	fmt.Println("hi")
}
