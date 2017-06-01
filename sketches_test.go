package adaptive

import (
	"fmt"
	"testing"
	"time"

	"github.com/seiflotfy/adaptive"
)

func TestEstimate(t *testing.T) {
	now := time.Now()
	sks := NewSketches(1000*time.Second, time.Second, 4, 7, 1.004)
	item1 := []byte("foo")

	count11 := uint64(1000000)
	timestamp11 := now.Add(1000 * time.Second)
	sks.Insert(item1, timestamp11, count11)

	got, err := sks.Estimate(item1, now, timestamp11)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	} else if got != count11 {
		t.Errorf("expected %d, got %d", count11, got)
	}

	item2 := []byte("bar")
	sks.Insert(item2, timestamp11, count11)
	got, err = sks.Estimate(item2, now, timestamp11)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	} else if got != count11 {
		t.Errorf("expected %d, got %d", count11, got)
	}

	count12 := uint64(1337)
	timestamp12 := now.Add(2000 * time.Second)
	sks.Insert(item1, timestamp12, count12)

	got, err = sks.Estimate(item1, now, timestamp12)
	if err == nil {
		t.Errorf("expected err, got %v", err)
	}

	got, err = sks.Estimate(item1, timestamp11, timestamp12.Add(-time.Second))
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
	timestamp13 := now.Add(10100 * time.Second)
	sks.Insert(item1, timestamp13, count13)
	got, err = sks.Estimate(item1, now, timestamp13)
	if err == nil {
		t.Errorf("expected err, got %v", err)
	} else if err == nil && got != count13 {
		t.Errorf("expected %d, got %d", count13, got)
	}

	got, err = sks.Estimate(item1, timestamp13.Add(-500), timestamp13.Add(500))
	if err != nil {
		t.Errorf("expected err, got %v", err)
	} else if got != count13 {
		t.Errorf("expected %d, got %d", count13, got)
	}
}

func TestEstimateReal(t *testing.T) {
	d := time.Duration(720)
	start := time.Now()
	sks := NewSketches(d*time.Hour, time.Hour, 7, 7, 1.004)
	item1 := []byte("foo")
	exp := uint64(0)

	for i := uint64(0); i < uint64(d*2); i++ {
		count := i + 1000
		exp += count
		pexp := exp * 5 / 100
		expRange := [2]uint64{exp - pexp, exp + pexp}
		end := start.Add(time.Duration(i) * time.Hour)
		sks.Insert(item1, end, count)
		got := sks.MultiEstimate(item1, start, end)
		if got < expRange[0] || got > expRange[1] {
			t.Errorf("expected %d, got %d", exp, got)
		}
	}
}

func Example() {
	duration := time.Duration(720 * time.Hour) // 720 hours range
	unit := time.Hour

	// Create sketch queryable with
	// duation = 720 hours range
	// unit = 1 hour
	// width per sketch = 2^9
	// depth per sketch = 8
	// alpha = 1.004 (used for emphasizing and de-emphasizing)
	sks := adaptive.NewSketches(duration, unit, 9, 7, 1.004)

	item := []byte("foo")
	exp := uint64(0)
	start := time.Now()

	count := uint64(1337)
	exp += count
	timestamp := start.Add(time.Hour)

	// Update item for given timestamp
	sks.Insert(item, timestamp, count)

	// Estimate count of item within time range [start, timestamp]
	got, _ := sks.Estimate(item, start, timestamp)
	fmt.Printf("Expected count for \"item\" %s in timerange [%v, %v] to be %d, got %d \n",
		string(item), start.Format(time.Kitchen), timestamp.Format(time.Kitchen), exp, got)

	// Move one hour
	timestamp = timestamp.Add(time.Hour)
	count = 100000
	exp += count

	// Update item for given timestamp
	sks.Insert(item, timestamp, count)

	// Estimate count of item within time range [start, timestamp]
	got, _ = sks.Estimate(item, start, timestamp)

	fmt.Printf("Expected count for \"item\" %s in timerange [%v, %v] to be %d, got %d \n",
		string(item), start.Format(time.Kitchen), timestamp.Format(time.Kitchen), got, exp)

	// Estimate count of item within time range [timestamp-12m, timestamp+12m]
	got, _ = sks.Estimate(item, timestamp.Add(-time.Hour/5), timestamp.Add(time.Hour/5))

	fmt.Printf("Expected count for \"item\" %s in timerange [%v, %v] to be %d, got %d \n",
		string(item), timestamp.Add(-time.Hour/5).Format(time.Kitchen), timestamp.Add(time.Hour/5).Format(time.Kitchen), got, count)

	// Output:
	// Expected count for "item" foo in timerange [2:56PM, 3:56PM] to be 1337, got 1337
	// Expected count for "item" foo in timerange [2:56PM, 4:56PM] to be 101337, got 101337
	// Expected count for "item" foo in timerange [4:44PM, 5:08PM] to be 100000, got 100000
}
