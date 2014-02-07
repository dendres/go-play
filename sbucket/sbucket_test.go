package sbucket

import (
	"testing"
	"time"
)

func TestStamp(t *testing.T) {
	now := time.Now()
	t.Log(now.Unix(), now.Nanosecond())

	for i := -11; i < 65; i++ {
		stamp, size := Stamp(now, i)
		t.Log(i, "\t", size, "\t", stamp)
	}

	// see if 12 and 15 are ok to drop a digit ??? NO!
	m := []int{65, 129, 333, 7777, 88888, 9999999}

	for _, i := range m {
		stamp, size := Stamp(now, i)
		t.Log(i, "\t", size, "\t", stamp)
	}

}
