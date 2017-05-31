package main

import (
	"fmt"
	"time"

	ada "github.com/seiflotfy/ada-sketches"
)

func main() {
	d := time.Duration(720)
	start := time.Now()
	sks := ada.NewSketches(d*time.Hour, time.Hour, 9, 7, 1.004)
	item1 := []byte("foo")
	exp := uint64(0)

	for i := uint64(0); i < uint64(d); i++ {
		count := i + 1000
		exp += count
		end := start.Add(time.Duration(i) * time.Hour)
		sks.Update(item1, end, count)

		got, _ := sks.Estimate(item1, start, end)
		fmt.Printf("Expected %d, got %d\n", exp, got)

		got, _ = sks.Estimate(item1, end.Add(-time.Hour/2), end.Add(time.Hour/2))
		fmt.Printf("Expected %d, got %d\n", count, got)
	}
}
