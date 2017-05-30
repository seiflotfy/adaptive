package ada

import (
	"testing"
	"time"
)

func TestEstimate(t *testing.T) {
	sks := NewSketches(1000, 4, 7, 1.004)
	item1 := []byte("foo")

	count11 := uint64(1000000)
	timestamp11 := uint64(1000)
	sks.Update(item1, timestamp11, count11)
	got, err := sks.Estimate(item1, 0, timestamp11)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	} else if got != count11 {
		t.Errorf("expected %d, got %d", count11, got)
	}

	item2 := []byte("bar")
	sks.Update(item2, timestamp11, count11)
	got, err = sks.Estimate(item2, 0, timestamp11)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	} else if got != count11 {
		t.Errorf("expected %d, got %d", count11, got)
	}

	count12 := uint64(1337)
	timestamp12 := uint64(2000)
	sks.Update(item1, timestamp12, count12)
	got, err = sks.Estimate(item1, 0, timestamp12)
	if err == nil {
		t.Errorf("expected err, got %v", err)
	} else if err == nil && got != count12 {
		t.Errorf("expected %d, got %d", count12, got)
	}

	got, err = sks.Estimate(item1, timestamp11, timestamp12-1)
	if err != nil {
		t.Errorf("expected err, got %v", err)
	} else if err == nil && got != count11 {
		t.Errorf("expected %d, got %d", count11, got)
	}

	got, err = sks.Estimate(item1, timestamp11, timestamp12)
	if err != nil {
		t.Errorf("expected err, got %v", err)
	} else if err == nil && got != count11+count12 {
		t.Errorf("expected %d, got %d", count11+count12, got)
	}

	count13 := uint64(900000)
	timestamp13 := uint64(10100)
	sks.Update(item1, timestamp13, count13)
	got, err = sks.Estimate(item1, 0, timestamp13)
	if err == nil {
		t.Errorf("expected err, got %v", err)
	} else if err == nil && got != count13 {
		t.Errorf("expected %d, got %d", count13, got)
	}

	got, err = sks.Estimate(item1, timestamp13-500, timestamp13+500)
	if err != nil {
		t.Errorf("expected err, got %v", err)
	} else if got != count13 {
		t.Errorf("expected %d, got %d", count13, got)
	}
}

func TestEstimateReal(t *testing.T) {
	d := 720
	now := uint64(time.Now().UnixNano() / int64(time.Hour))
	sks := NewSketches(uint64(d), 7, 7, 1.004)
	item1 := []byte("foo")
	exp := uint64(0)

	for i := uint64(0); i < uint64(d); i++ {
		count := i + 1000
		exp += count
		pexp := exp * 5 / 100
		expRange := [2]uint64{exp - pexp, exp + pexp}
		sks.Update(item1, now+i, count)
		got, err := sks.Estimate(item1, now, now+uint64(i))
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		} else if got < expRange[0] || got > expRange[1] {
			t.Errorf("expected %d, got %d", exp, got)
		}
	}
}
