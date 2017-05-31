# adaptive
A probabilistic data structure that serves as a frequency table of events in a stream of data over time range, using a set of slightly modified Count-Min Sketches

# Abstract
[Time Adaptive Sketches (Ada-Sketches) for Summarizing Data Streams](https://www.cs.rice.edu/~as143/Papers/16-ada-sketches.pdf) - Anshumali Shrivastava, Arnd Christian König, Mikhail Bilenko
Obtaining frequency information of data streams, in limited space, is a well-recognized problem in literature. A num- ber of recent practical applications (such as those in com- putational advertising) require temporally-aware solutions: obtaining historical count statistics for both time-points as well as time-ranges. In these scenarios, accuracy of estimates is typically more important for recent instances than for older ones; we call this desirable property Time Adap- tiveness. With this observation, [20] introduced the Hokusai technique based on count-min sketches for estimating the frequency of any given item at any given time. The proposed approach is problematic in practice, as its memory require- ments grow linearly with time, and it produces discontinuities in the estimation accuracy. In this work, we describe a new method, Time-adaptive Sketches, (Ada-sketch), that overcomes these limitations, while extending and providing a strict generalization of several popular sketching algorithms.

# Usage
```go
d := time.Duration(720) // 720 hours range
unit := time.Hour

// Create sketch queryable over 720 hours range, where a unit is an hour
sks := adaptive.NewSketches(d*unit, unit, 9, 7, 1.004)

item := []byte("foo")
exp := uint64(0)
start := time.Now()

for i := uint64(0); i < uint64(d); i++ {
    count := i + 1000
    exp += count
    timestamp := start.Add(time.Duration(i) * time.Hour)

    // Update item for given timestamp
    sks.Update(item, timestamp, count)

    // Estimate for range since start till timestamp
    got, _ := sks.Estimate(item, start, timestamp)
    fmt.Printf("Expected %d, got %d\n", exp, got)

    got, _ = sks.Estimate(item, timestamp.Add(-time.Hour/5), timestamp.Add(time.Hour/5))
    fmt.Printf("Expected %d, got %d\n", count, got)
}
```
