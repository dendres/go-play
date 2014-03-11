package pile

import (
	"github.com/dendres/go-play/event"
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
func TestAll(t *testing.T) {

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

	x := event.EventBytes{[]byte{0x31, 0x32, 0x33}}

	err = p.Append(x)
	if err != nil {
		t.Fatal("error appending:", err)
	}

}
