// Copyright (c) <year> <copyright holder>.
// etc...

/*
Package organizing demonstrates go conventions related to code organization and documentation.
Special thanks to http://blog.golang.org/organizing-go-code
and http://blog.golang.org/godoc-documenting-go-code

To setup a new development environment:

    export GOPATH=$HOME/go
    export PATH=~$PATH:$GOPATH/bin
    [ -e $GOPATH ] || mkdir -p $GOPATH
    [ -e $GOPATH/bin ] || mkdir -p $GOPATH/bin
    [ -e $GOPATH/pkg ] || mkdir -p $GOPATH/pkg
    [ -e $GOPATH/src ] || mkdir -p $GOPATH/src

package name:
 name carefully.
 provides context for everything in the package.
 A well-chosen name is therefore the starting point for good documentation
 short, concise, evocative
 By convention, packages are given lower case, single-word names; there should be no need for underscores or mixedCaps.

package import path:
 make your code "go get"-able
 globally unique
 The last element of the import path should be the same as the package name.
  code.google.com/p/go.net/websocket
  github.com/you/your-package
  github.com/you/your-project/your-package

minimize the exported interface: The larger the interface you provide, the more you must support.

package main is often larger than other packages. Complex commands contain a lot of code that is of little use outside the context of the executable, and often it's simpler to just keep it all in the one place.

If a package will change in a backward incompatible way:
1. keep the old "package import path" backward compatible
2. make a new "package import path" with the new incompatible code
2. Write a gofix to automate the migration to the new package ?????

Indentation:
    We use tabs for indentation and gofmt emits them by default. Use spaces only if you must.

Line length:
    Go has no line length limit. Don't worry about overflowing a punched card.
        If a line feels too long, wrap it and indent with an extra tab.

Parentheses:
    Go needs fewer parentheses than C and Java:
    control structures (if, for, switch) do not have parentheses in their syntax.
    Also, the operator precedence hierarchy is shorter and clearer, so
        x<<8 + y<<16
    means what the spacing implies, unlike in the other languages.

Comments:
    to document a (type, variable, constant, function, package)
    write a regular comment directly preceding its declaration, with no intervening blank line
    Doc comments work best as complete sentences, which allow a wide variety of automated presentations.
    no need for comment formatting like banners of stars etc...
    plain text only. avoid HTML or other annotations
    if a package required a long introduction, then make a separate file with just documentation doc.go
    gdoc -> html formatting rules:
    * leave a blank line to separate paragraphs
    * indent preformatted text
    * url's become html links

Name conventions in go:
    type, variable, constant, function:
    http://golang.org/doc/effective_go.html#names
    Long names don't automatically make things more readable.
    A helpful doc comment can often be more valuable than an extra long name.
    cammelCase seems to appear frequently in twoWord function and variable names???
*/
package organizing
