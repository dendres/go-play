package main

import (
	"bytes"
	"log"
	"log/syslog"
)

func main() {

	// logger, err := NewLogger(syslog.LOG_WARNING | LOG_DAEMON, logFlag int)

	l2, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DAEMON, "go_syslog_output_test")
	defer l2.Close()
	if err != nil {
		log.Fatal("error writing syslog!")
	}

	// l2.Notice("this is a notice test")
	// l2.Debug("this is a debug test")
	// 2014-02-06T01:05:03.135904+00:00,notice,daemon,go_test_tag, this is a notice test
	// 2014-02-06T01:05:03.135942+00:00,debug,daemon,go_test_tag, this is a debug test

	// try to send a huge string!
	var b bytes.Buffer

	for i := 0; i < 10000; i++ {
		b.WriteString("0123456789")
	}

	l2.Debug(b.String())
	// Success!  got a huge string in syslog
	// with $MaxMessageSize 64k, got a message with 65556 characters

}
