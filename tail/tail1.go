package main

import (
	"bufio"
	"log"
	"os/exec"
	"time"
)

// tail -n0 -F path
// sleep one second and restart if the tail command exits for any reason
// send lines to the supplied channel
func tailer(path string, out chan string) {
	for {
		log.Println("starting tailer")

		cmd := exec.Command("tail", "-n0", "-F", path)

		stdoutpipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		stdout := bufio.NewScanner(stdoutpipe)

		log.Println("running tail cmd")
		cmd.Start()

		for stdout.Scan() {
			out <- stdout.Text()
		}

		cmd.Wait()
		time.Sleep(time.Second)
	}
}

// listen to the supplied channel and log each line
func logger(out chan string) {
	log.Println("starting logger")
	for line := range out {
		log.Println("got:", line)
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("start main")

	out := make(chan string, 1000)

	go tailer("/tmp/tailing.txt", out)
	logger(out)
}
