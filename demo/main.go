package main

import (
	"fmt"
	"time"

	"github.com/seiflotfy/adaptive"
)

func main() {
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
