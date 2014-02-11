package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

/*
event serialization format must be:
* a very compact binary format for long term storage of 5min chunks of log data from a single source/process
* json-compatable nested hash and arrays of string, float, int, bool, etc...
* schema free
* does not need to be streamable because the routing/indexing fields will be part of the kafka protocol
* based on existing standards
* clean utf-8 strings would be nice
* 8-bit clean binary fields would also be nice
* relatively easy to port the decoder to a new language


input format:
* syslog lines, json, xml and various other formats become some kind of go structure??
* how does heka parse plaintext with regex????





ideally it should have a table of field names like:
key1,key2,key3,key4,0:value1,1:value2



existing research:
* https://github.com/eishay/jvm-serializers/wiki
* http://leopard.in.ua/2013/10/13/binary-serialization-formats/
* http://en.wikipedia.org/wiki/Comparison_of_data_serialization_formats
* http://web-resource-optimization.blogspot.com/2011/06/json-compression-algorithms.html

candidates:
* compress(capn(event_object))
  https://github.com/jmckaskill/go-capnproto

* kryo????
* msgpack


* thrift??
* compress(hpack(json(event_object)))


http://msgpack.org/??? not tabular???





*/
func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")
}
