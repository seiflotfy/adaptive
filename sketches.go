package adaptive

import (
	"fmt"
	"math"
	"time"
)

// Sketches represents an set of Adaptive Count-Min Sketch algorithm (Ada-CMS),
// which is just CMS but with the update and query mechanisms adapted to use the pre-emphasis and de-emphasis mechanism.
type Sketches struct {
	sketches    []*CMS
	maxDuration time.Duration
	timeUnit    time.Duration
	w           uint64
	d           uint64
	alpha       float64
}

// NewSketches returns an adaptive Sketch (representing multiple Count-Min Sketches) with:
// maxDuration: specifies the max queriable range
// timeUnit: degines the unit of measurment, ns, ms, s, etc...
// w: width of CMS is 2^w
// d: depth of CMS
// alpha: ussed to pre-emphasis and de-emphasis the counters using alpha^timestamp
func NewSketches(maxDuration, timeUnit time.Duration, w, d uint64, alpha float64) *Sketches {
	convDuration := maxDuration / timeUnit
	numSketches := uint64(math.Log2(float64(convDuration))) + 1
	sketches := make([]*CMS, numSketches)
	for i := range sketches {
		sketches[i] = NewCMS(w, d, alpha)
	}
	return &Sketches{
		sketches:    sketches,
		maxDuration: maxDuration,
		timeUnit:    timeUnit,
		w:           w,
		d:           d,
		alpha:       alpha,
	}
}

func (sks *Sketches) generateTimestamp(timestamp uint64, i int) uint64 {
	timestamp = timestamp % uint64(sks.maxDuration)
	return 1 + timestamp/(uint64(math.Pow(2, float64(i))))
}

// Update an item with a new count with at a given timestamp
func (sks *Sketches) Update(item []byte, timestamp time.Time, count uint64) error {
	tmpTimestamp := uint64(timestamp.UnixNano() / int64(sks.timeUnit))
	for i := range sks.sketches {
		t := sks.generateTimestamp(tmpTimestamp, i)
		sks.sketches[i].Update(item, t, count)
	}
	return nil
}

func (sks *Sketches) estimate(item []byte, start, end uint64) uint64 {
	estimate := uint64(0)
	for start <= end {
		pow2 := float64(start & (^start + 1))
		logpow2 := math.Log2(pow2)
		for i := logpow2; i >= 0; i-- {
			if float64(start)+math.Pow(2, i)-1 <= float64(end) {
				t := sks.generateTimestamp(start, int(i))
				estimate += sks.sketches[uint64(i)].Count(item, t)
				start += uint64(math.Pow(2, i))
				break
			}
		}
	}
	return estimate
}

// Estimate returns the count of an item within a given timerange [start, end], if timerange > maxDuration return error
func (sks *Sketches) Estimate(item []byte, start, end time.Time) (uint64, error) {
	tmpStart := uint64(start.UnixNano() / int64(sks.timeUnit))
	tmpEnd := uint64(end.UnixNano() / int64(sks.timeUnit))
	maxDuration := uint64(sks.maxDuration) / uint64(sks.timeUnit)

	if tmpEnd-tmpStart > maxDuration {
		return 0, fmt.Errorf("window to big [start, end] %d > %d", end.Sub(start)/sks.timeUnit, sks.maxDuration/sks.timeUnit)
	}

	return sks.estimate(item, tmpStart, tmpEnd), nil
}

// MultiEstimate  ...
func (sks *Sketches) MultiEstimate(item []byte, start, end time.Time) uint64 {
	estimate := uint64(0)
	orgEnd := uint64(end.UnixNano() / int64(sks.timeUnit))
	tmpStart := uint64(start.UnixNano() / int64(sks.timeUnit))
	maxDuration := uint64(sks.maxDuration) / uint64(sks.timeUnit)

	for i := tmpStart; i <= orgEnd; i += maxDuration {
		tmpEnd := i + maxDuration - 1
		if tmpEnd > orgEnd {
			tmpEnd = orgEnd
		}
		res := sks.estimate(item, i, tmpEnd)
		estimate += res
	}
	return estimate
}
