package ada

import (
	"encoding/binary"
	"math"

	metro "github.com/dgryski/go-metro"
)

// Sketch is a simple count-min-sketch imeplementation with a exp based pre-emphasis and de-emphasis of values
type Sketch struct {
	regs  [][]float64
	w     uint64
	d     uint64
	alpha float64
}

// NewSketch returns a Sketch of width 2^w and depth d for a given alpha
// the alpha is ussed to pre-emphasis and de-emphasis the counters using alpha^timestamp
func NewSketch(w, d uint64, alpha float64) *Sketch {
	regs := make([][]float64, d)
	for i := range regs {
		regs[i] = make([]float64, uint64(math.Pow(2, float64(w))))
	}
	return &Sketch{
		w:     w,
		d:     d,
		regs:  regs,
		alpha: alpha,
	}
}

// Update ...
func (sk *Sketch) Update(value []byte, timestamp, count uint64) {
	for i := range sk.regs {
		j := sk.hash(value, timestamp, uint64(i))
		sk.regs[i][j] += (sk.factor(timestamp) * float64(count))
	}
}

// Count ...
func (sk *Sketch) Count(value []byte, timestamp uint64) uint64 {
	min := math.MaxFloat64
	for i := range sk.regs {
		j := sk.hash(value, timestamp, uint64(i))
		if sk.regs[i][j] < min {
			min = sk.regs[i][j]
		}
	}
	return uint64(min / sk.factor(timestamp))
}

func (sk *Sketch) hash(item []byte, timestamp, hashid uint64) uint64 {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, timestamp)
	hj := metro.Hash64(timeBytes, hashid)
	return metro.Hash64(item, hj) % uint64(len(sk.regs[0]))
}

func (sk *Sketch) factor(timestamp uint64) float64 {
	return math.Pow(sk.alpha, 1/float64(timestamp))
}
