# adaptive

A probabilistic datastructure to estimate **How many times did i see item "x" within the timerange "T"**, using a set Adaptive Count-Min Sketches.

The Adaptive Count-Min Sketch algorithm (Ada-CMS), is just CMS but with the update and query mechanisms adapted to use the pre-emphasis and de-emphasis mechanism.

For more information read the post by Adrian Colyer [Time Adaptive Sketches (Ada-Sketches) for Summarizing Data Streams](https://blog.acolyer.org/2016/07/21/time-adaptive-sketches-ada-sketches-for-summarizing-data-streams/) or [the official paper with the same title](https://www.cs.rice.edu/~as143/Papers/16-ada-sketches.pdf) by Anshumali Shrivastava, Arnd Christian König, Mikhail Bilenko) 

## Usage
```go
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
	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
		string(item), start.Format(time.Kitchen), timestamp.Format(time.Kitchen), exp, got)

	// Move one hour and set new count
	timestamp = timestamp.Add(time.Hour)
	count = 100000
	exp += count

	// Update item for given timestamp
	sks.Insert(item, timestamp, count)

	// Estimate count of item within time range [start, timestamp]
	got, _ = sks.Estimate(item, start, timestamp)

	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
		string(item), start.Format(time.Kitchen), timestamp.Format(time.Kitchen), got, exp)

	// Estimate count of item within time range [timestamp-12m, timestamp+12m]
	t1 := timestamp.Add(-time.Hour / 5)
	t2 := timestamp.Add(time.Hour / 5)
	got, _ = sks.Estimate(item, t1, t2)

	fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d \n",
		string(item), t1.Format(time.Kitchen), t2.Format(time.Kitchen), got, count)

	// Output:
	// Expected count for "foo" in timerange [2:56PM, 3:56PM] to be 1337, got 1337
	// Expected count for "foo" in timerange [2:56PM, 4:56PM] to be 101337, got 101337
	// Expected count for "foo" in timerange [4:44PM, 5:08PM] to be 100000, got 100000
```
