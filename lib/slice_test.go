package lib

import (
	"testing"
)

func TestArraySliceInt64(t *testing.T) {
	sids := ArraySliceInt64([]int64{1}, 101)
	if len(sids) != 1 && len(sids[0]) != 1 {
		t.Fatal("empty")
	}

	ids := []int64{}

	sids = ArraySliceInt64(ids, 101)
	if len(sids) != 0 {
		t.Fatal("0")
	}

	for i := 0; i < 1000; i++ {
		ids = append(ids, int64(i))
	}
	sids = ArraySliceInt64(ids, 100)
	for _, vs := range sids {
		if len(vs) != 100 {
			t.Fatal("1")
		}
	}

	sids = ArraySliceInt64(ids, 101)
	for k, vs := range sids {
		if len(vs) != 101 && k != len(sids)-1 {
			t.Fatal("2")
		}
	}
}
