/*
Package coverage is a test coverage example from http://blog.golang.org/cover

References:

    go help test
    go help testfunc
    go help testflag
    go help packages

setup go test coverage as root:

    yum install golang-googlecode-tools-devel
    export GOPATH=/home/done/go
    export PATH=~$PATH:$GOPATH/bin
    go get code.google.com/p/go.tools/cmd/cover

run all examples with "// Output: " blocks and all tests

    go test

use test coverage from http://blog.golang.org/cover

    go test --cover

some test coverage data that looks a lot like profiling data could also be available:

   go test --coverprofile=coverage.out
   go test --covermode=count --coverprofile=count.out

run benchmarks:

   go test -bench=.

Why does Go not have assertions? http://golang.org/doc/faq#assertions

Go doesn't provide assertions. They are undeniably convenient, but our experience has been that programmers use them as a crutch to avoid thinking about proper error handling and reporting. Proper error handling means that servers continue operation after non-fatal errors instead of crashing. Proper error reporting means that errors are direct and to the point, saving the programmer from interpreting a large crash trace. Precise errors are particularly important when the programmer seeing the errors is not familiar with the code.

We understand that this is a point of contention. There are many things in the Go language and libraries that differ from modern practices, simply because we feel it's sometimes worth trying a different approach.

*/
package coverage
