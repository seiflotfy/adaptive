package ada

import (
	"encoding/binary"
	"fmt"
	"math"

	metro "github.com/dgryski/go-metro"
)

// Sketches ...
type Sketches struct {
	sketches []*sketch
	maxTime  uint64
}

// NewSketches ...
func NewSketches(maxTime, w, d uint64, alpha float64) *Sketches {
	num := uint64(math.Log2(float64(maxTime))) + 1
	sketches := make([]*sketch, num)
	for i := range sketches {
		sketches[i] = newSketch(w, d, alpha)
	}
	return &Sketches{
		sketches: sketches,
		maxTime:  maxTime,
	}
}

// Update ...
func (sks *Sketches) Update(item []byte, timestamp, count uint64) {
	if timestamp > sks.maxTime {
		fmt.Println("Error: More than Max Time")
		return
	}
	for i := range sks.sketches {
		t := 1 + timestamp/uint64(math.Pow(2, float64(i)))
		sks.sketches[i].Update(item, t, count)
	}
}

// Estimate ...
func (sks *Sketches) Estimate(item []byte, start, end uint64) uint64 {
	estimate := uint64(0)
	for start <= end {
		pow2 := float64(start & (^start + 1))
		logpow2 := math.Log2(pow2)
		for i := logpow2; i >= 0; i-- {
			if start+uint64(math.Pow(2, i))-1 <= end {
				timestamp := 1 + start/uint64(math.Pow(2, i))
				estimate += sks.sketches[uint64(i)].Count(item, timestamp)
				start += uint64(math.Pow(2, i))
			}
		}
	}
	return estimate
}

// sketch ...
type sketch struct {
	regs  [][]float64
	w     uint64
	d     uint64
	alpha float64
}

func newSketch(w, d uint64, alpha float64) *sketch {
	regs := make([][]float64, d)
	for i := range regs {
		regs[i] = make([]float64, uint64(math.Pow(2, float64(w))))
	}

	return &sketch{
		w:     w,
		d:     d,
		regs:  regs,
		alpha: alpha,
	}
}

// Update ...
func (sk *sketch) Update(value []byte, timestamp, count uint64) {
	for i := range sk.regs {
		j := sk.hash(value, timestamp, uint64(i))
		sk.regs[i][j] += (sk.factor(timestamp) * float64(count))
	}
}

// Count ...
func (sk *sketch) Count(value []byte, timestamp uint64) uint64 {
	min := math.MaxFloat64
	for i := range sk.regs {
		j := sk.hash(value, timestamp, uint64(i))
		if sk.regs[i][j] < min {
			min = sk.regs[i][j]
		}
	}
	return uint64((min / sk.factor(timestamp)) + 1)
}

func (sk *sketch) hash(item []byte, timestamp, hashid uint64) uint64 {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, timestamp)
	hj := metro.Hash64(timeBytes, hashid)
	return metro.Hash64(item, hj) % uint64(len(sk.regs[0]))
}

func (sk *sketch) factor(timestamp uint64) float64 {
	return math.Pow(sk.alpha, float64(timestamp))
}
