package adaptive

import (
	"fmt"
	"testing"
	"time"
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

func TestExample(t *testing.T) {
	duration := time.Duration(720 * time.Hour) // 720 hours range
	unit := time.Hour

	// Create sketch queryable with
	// duation = 720 hours range
	// unit = 1 hour
	// width per sketch = 2^9
	// depth per sketch = 8
	// alpha = 1.004 (used for emphasizing and de-emphasizing)
	sks := NewSketches(duration, unit, 9, 7, 1.004)

	item := []byte("foo")
	t1 := time.Date(2017, 06, 03, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	t4 := t3.Add(time.Hour)
	count1 := uint64(1337)
	count2 := uint64(100000)

	// Update item for given timestamps
	sks.Insert(item, t1, count1)
	sks.Insert(item, t3, count2)

	// Estimate count of item within time range [t1, t2]
	got, _ := sks.Estimate(item, t1, t2)
	if count1 != got {
		t.Errorf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
			string(item), t1.Format(time.Kitchen), t2.Format(time.Kitchen), count1, got)
	}

	// Estimate count of item within time range [t1, t3]
	got, _ = sks.Estimate(item, t1, t3)
	if count1+count2 != got {
		t.Errorf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
			string(item), t1.Format(time.Kitchen), t3.Format(time.Kitchen), count1+count2, got)
	}

	// Estimate count of item within time range [t1, t3]
	got, _ = sks.Estimate(item, t3, t4)
	if count2 != got {
		t.Errorf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
			string(item), t3.Format(time.Kitchen), t4.Format(time.Kitchen), count2, got)
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
	sks := NewSketches(duration, unit, 9, 7, 1.004)

	item := []byte("foo")
	t1 := time.Date(2017, 06, 03, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	t4 := t3.Add(time.Hour)
	count1 := uint64(1337)
	count2 := uint64(100000)

	// Update item for given timestamps
	sks.Insert(item, t1, count1)
	sks.Insert(item, t3, count2)

	// Estimate count of item within time range [t1, t2]
	got, _ := sks.Estimate(item, t1, t2)
	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
		string(item), t1.Format(time.Kitchen), t2.Format(time.Kitchen), count1, got)

	// Estimate count of item within time range [t1, t3]
	got, _ = sks.Estimate(item, t1, t3)
	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
		string(item), t1.Format(time.Kitchen), t3.Format(time.Kitchen), count1+count2, got)

	// Estimate count of item within time range [t1, t3]
	got, _ = sks.Estimate(item, t3, t4)
	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
		string(item), t3.Format(time.Kitchen), t4.Format(time.Kitchen), count2, got)

	// Output:
	// Expected count for "foo" in timerange [12:00AM, 1:00AM] to be 1337, got 1337
	// Expected count for "foo" in timerange [12:00AM, 2:00AM] to be 101337, got 101337
	// Expected count for "foo" in timerange [2:00AM, 3:00AM] to be 100000, got 100000
}
