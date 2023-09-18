package consistenthash_test

import (
	"strconv"
	"testing"

	"github.com/454270186/GoCache/gocache/consistenthash"
)

func TestHashRing(t *testing.T) {
	hr := consistenthash.New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	hr.Add("2", "4", "6")

	testCases := map[string]string{
		"2": "2",
		"13": "4",
		"25": "6",
		"99": "2",
	}

	for k, v := range testCases {
		if hr.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hr.Add("5")
	if hr.Get("25") != "5" {
		t.Errorf("Asking for %s, should have yielded %s", "25", "5")
	}
}