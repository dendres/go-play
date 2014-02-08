package sbucket

import (
	"testing"
	"time"
)

func TestTenStamp(t *testing.T) {
	now := time.Now()
	t.Log(now.Unix(), now.Nanosecond())

	for i := -11; i < 65; i++ {
		stamp, size := TenStamp(now, i)
		t.Log(i, "\t", size, "\t", stamp)
	}

	// see if 12 and 15 are ok to drop a digit ??? NO!
	m := []int{65, 129, 333, 7777, 88888, 9999999, 3600, 36000, 360000, 3600000, 36000000, 360000000, 3600000000, 36000000000}

	for _, i := range m {
		stamp, size := TenStamp(now, i)
		t.Log(i, "\t", size, "\t", stamp)
	}

}
