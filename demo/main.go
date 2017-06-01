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

	for i := uint64(0); i < uint64(duration); i++ {
		count := i + 1000
		exp += count
		timestamp := start.Add(time.Duration(i) * time.Hour)

		// Update item for given timestamp
		sks.Update(item, timestamp, count)

		// Estimate count of item within time range [start, timestamp]
		got, _ := sks.Estimate(item, start, timestamp)
		fmt.Printf("Expected %d, got %d\n", exp, got)

		// Estimate count of item within time range [timestamp-12m, timestamp+12m]
		got, _ = sks.Estimate(item, timestamp.Add(-time.Hour/5), timestamp.Add(time.Hour/5))
		fmt.Printf("Expected %d, got %d\n", count, got)
	}
}
