// go help test
// go help testfunc
// go help testflag
// go help packages

// setup go test coverage as root:
//   yum install golang-googlecode-tools-devel
//   export GOPATH=/home/done/go
//   export PATH=~$PATH:$GOPATH/bin
//   go get code.google.com/p/go.tools/cmd/cover

// use test coverage from http://blog.golang.org/cover
//  go test --cover

// some test coverage data that looks a lot like profiling data could also be available:
//  go test --coverprofile=coverage.out
//  go test --covermode=count --coverprofile=count.out

// http://blog.golang.org/organizing-go-code

// coverage is a test coverage example from http://blog.golang.org/cover
package coverage

func Size(a int) string {
	switch {
	case a < 0:
		return "negative"
	case a == 0:
		return "zero"
	case a < 10:
		return "small"
	case a < 100:
		return "big"
	case a < 1000:
		return "huge"
	}
	return "enormous"
}
