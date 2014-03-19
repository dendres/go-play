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

func newTestPath() string {
	stamp := time.Now().Nanosecond()
	return base_path + "/" + strconv.Itoa(stamp)
}

func newTestPile(t *testing.T) *Pile {
	path := newTestPath()
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
	for i := 0; i < len(words); i++ {
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

// RandomEvents returns a slice of events containing pseudorandom values
func RandomEvents(count int) []*event.Event {
	events := make([]*event.Event, count)

	for i := 0; i < count; i++ {
		events[i] = RandEvent()
	}
	return events
}

// WriteEvents serializes and appends the events from the slice provided to the Pile provided.
// on a quiet server, I have not been able to detect out of order or even disjointed writes.
// vfs and ext4 docs suggest it is possible, so I will continue to guard against it.
func WriteEvents(events []*event.Event, p *Pile, t *testing.T, s int) {
	for _, e := range events {

		// XXX a sleep of 1ns causes the sequential write pattern to be measurable
		// XXX remove the sleep, and the writes seem to appear as a single write ???
		if s > 0 {
			time.Sleep(time.Duration(s))
		}
		eb, err := e.Encode()
		if err != nil {
			t.Fatal("Encoding error:", err)
		}

		err = p.Append(eb)
		if err != nil {
			t.Fatal("error appending:", err)
		}
	}
}

// TestWriteRead Writes random events to a file, reads them from file, then decodes and verifies them.
func TestWriteRead(t *testing.T) {
	test_event_count := 50
	p := newTestPile(t)
	events := RandomEvents(test_event_count)
	WriteEvents(events, p, t, 0)

	events2 := make([]*event.EventBytes, len(events), test_event_count)
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

// try to demonstrate reading from a pile while it is being written
func TestSequentialWriteWithSleep(t *testing.T) {
	pile_path := newTestPath()
	t.Log("pile_path =", pile_path)
	test_event_count := 5

	go func() {
		p1, err := NewPile(pile_path)
		if err != nil {
			t.Fatal("Error opening path =", pile_path, ", error =", err)
		}

		events := RandomEvents(test_event_count)
		WriteEvents(events, p1, t, 1) // sleep 1ns between writes
	}()

	p2, err := NewPile(pile_path)
	if err != nil {
		t.Fatal("Error opening path =", pile_path, ", error =", err)
	}

	events2 := make([]*event.EventBytes, test_event_count)

	t.Log("len(events2) and test_event_count are ", len(events2), test_event_count)

	for i := 0; i < 10; i++ {
		time.Sleep(1) // sleep 1ns between reads
		p2.Read(events2)

		// print before all the events get written
		if events2[0] != nil && events2[test_event_count-1] == nil {
			t.Log("events2 =", events2)
		}
	}
}

func Read2(path string, count int, t *testing.T) {
	p2, err := NewPile(path)
	if err != nil {
		t.Fatal("Error opening path =", path, ", error =", err)
	}

	events2 := make([]*event.EventBytes, count)

	t.Log("len(events2) and test_event_count are ", len(events2), count)

	for i := 0; i < 10; i++ {
		time.Sleep(1) // sleep 1ns between reads
		p2.Read(events2)

		if events2[0] != nil && events2[count-1] == nil {
			t.Log("events2 =", events2)
		}
	}
}

// demonstrate one write goroutine and many read goroutines
func TestMultipleReadGoroutines(t *testing.T) {
	pile_path := newTestPath()
	t.Log("pile_path =", pile_path)
	test_event_count := 5

	go func() {
		p1, err := NewPile(pile_path)
		if err != nil {
			t.Fatal("Error opening path =", pile_path, ", error =", err)
		}

		events := RandomEvents(test_event_count)
		WriteEvents(events, p1, t, 1) // sleep 1ns between writes
	}()

	go Read2(pile_path, test_event_count, t)
	go Read2(pile_path, test_event_count, t)
	go Read2(pile_path, test_event_count, t)
	go Read2(pile_path, test_event_count, t)
	go Read2(pile_path, test_event_count, t)
	go Read2(pile_path, test_event_count, t)

	time.Sleep(1 * time.Second)
}
