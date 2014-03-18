package pile

import (
	"github.com/dendres/go-play/event"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

const base_path = "/opt/kafka/disks/0/test/pile"

// A Case is any test case that can be run
type Case interface {
	Run(t *testing.T)
}

func newTestPile(t *testing.T) *Pile {
	stamp := time.Now().Nanosecond()
	path := base_path + "/" + strconv.Itoa(stamp)
	t.Log("new test path =", path)

	p, err := NewPile(path)
	if err != nil {
		t.Fatal("Error opening path =", path, ", error =", err)
	}

	t.Log("opened pile p =", p.writer.Name())
	return p
}

// run all the test cases
func TestSingleReadWrite(t *testing.T) {

	err := os.RemoveAll(base_path)
	if err != nil {
		t.Fatal("error removing base_path =", base_path, ", err =", err)
	}

	err = os.Mkdir(base_path, 0777)
	if err != nil {
		t.Fatal("error creating base_path =", base_path, ", err =", err)
	}

	// make an event []byte
	// write it to file
	p := newTestPile(t)
	t.Log("p = ", p)

	e1 := event.Event{
		Enc:   0,
		Repl:  1,
		Pri:   2,
		Acc:   3,
		Point: event.Point(time.Now()),
		Data:  []byte("hello there"),
	}
	eb1, err := e1.Encode()
	if err != nil {
		t.Fatal("Encoding error:", err)
	}

	err = p.Append(eb1)
	if err != nil {
		t.Fatal("error appending:", err)
	}

	// read the newly created file and parse it into an event
	events := make([]*event.EventBytes, 1, 99)
	p.Read(events)

	for _, e := range events {
		t.Log("got event bytes =", e)
		event, err := e.Decode()
		if err != nil {
			t.Fatal("event decode error:", err)
		}

		t.Log("got event =", event)
	}
}

// RandEvent chooses a random value for all attributes and returns the Event
func RandEvent() *event.Event {
	words := []string{"bit", "manipulation", "is", "the", "act", "of",
		"algorithmically", "manipulating", "bits", "or", "other",
		"pieces", "of", "data", "shorter", "than", "a", "word"}

	sentence := ""
	for i := 0; i < 9; i++ {
		sentence += words[rand.Intn(len(words))] + " "
	}

	e := event.Event{
		Enc:   rand.Intn(3),
		Repl:  rand.Intn(3),
		Pri:   rand.Intn(3),
		Acc:   rand.Intn(3),
		Point: uint64(rand.Int63()),
		Data:  []byte(sentence),
	}
	return &e
}

func TestWriteRead(t *testing.T) {
	events := make([]*event.Event, 20, 40)

	for i := 0; i < 20; i++ {
		events[i] = RandEvent()
	}

	p := newTestPile(t)

	// write the events
	for _, e := range events {
		eb, err := e.Encode()
		if err != nil {
			t.Fatal("Encoding error:", err)
		}

		err = p.Append(eb)
		if err != nil {
			t.Fatal("error appending:", err)
		}
	}

	// read the events
	events2 := make([]*event.EventBytes, len(events), 40)
	p.Read(events2)

	for i, e := range events2 {
		event, err := e.Decode()
		if err != nil {
			t.Fatal("event decode error:", err)
		}

		if event.Point != events[i].Point {
			t.Fatalf("expected, vs got: %v, %v", event, events[i])
		}
	}

}

// XXX test 1 write and 1 read goroutine

// XXX test with 1 write and many read goroutines
