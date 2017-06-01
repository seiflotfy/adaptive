# adaptive

## Adaptive Sketches
A probabilistic datastructure to estimate **How many times did i see item "x" within the timerange "T"**, using a set Adaptive Count-Min Sketches.

The Adaptive Count-Min Sketch algorithm (Ada-CMS), is just CMS but with the update and query mechanisms adapted to use the pre-emphasis and de-emphasis mechanism.

For more information read the post by Adrian Colyer [Time Adaptive Sketches (Ada-Sketches) for Summarizing Data Streams](https://blog.acolyer.org/2016/07/21/time-adaptive-sketches-ada-sketches-for-summarizing-data-streams/) or [the official paper with the same title](https://www.cs.rice.edu/~as143/Papers/16-ada-sketches.pdf) by Anshumali Shrivastava, Arnd Christian König, Mikhail Bilenko) 

## Usage
```go
duration := time.Duration(720 * time.Houw) // 720 hours range
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

for i := uint64(0); i < uint64(d); i++ {
    count := i + 1000
    exp += count
    timestamp := start.Add(time.Duration(i) * time.Hour)

    // Insert item for given timestamp
    sks.Insert(item, timestamp, count)

    // Estimate for range since start till timestamp
    got, _ := sks.Estimate(item, start, timestamp)
    fmt.Printf("Expected %d, got %d\n", exp, got)

    got, _ = sks.Estimate(item, timestamp.Add(-time.Hour/5), timestamp.Add(time.Hour/5))
    fmt.Printf("Expected %d, got %d\n", count, got)
}
```
