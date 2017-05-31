package adaptive

import (
	"encoding/binary"
	"math"

	metro "github.com/dgryski/go-metro"
)

// CMS is a simple count-min-sketch imeplementation with a exp based pre-emphasis and de-emphasis of values
type CMS struct {
	regs  [][]float64
	w     uint64
	d     uint64
	alpha float64
}

// NewCMS returns a Sketch of width 2^w and depth d for a given alpha
// the alpha is ussed to pre-emphasis and de-emphasis the counters using alpha^timestamp
func NewCMS(w, d uint64, alpha float64) *CMS {
	regs := make([][]float64, d)
	for i := range regs {
		regs[i] = make([]float64, uint64(math.Pow(2, float64(w))))
	}
	return &CMS{
		w:     w,
		d:     d,
		regs:  regs,
		alpha: alpha,
	}
}

// Update ...
func (cms *CMS) Update(value []byte, timestamp, count uint64) {
	for i := range cms.regs {
		j := cms.hash(value, timestamp, uint64(i))
		cms.regs[i][j] += (cms.factor(timestamp) * float64(count))
	}
}

// Count ...
func (cms *CMS) Count(value []byte, timestamp uint64) uint64 {
	min := math.MaxFloat64
	for i := range cms.regs {
		j := cms.hash(value, timestamp, uint64(i))
		if cms.regs[i][j] < min {
			min = cms.regs[i][j]
		}
	}
	return uint64(min / cms.factor(timestamp))
}

func (cms *CMS) hash(item []byte, timestamp, hashid uint64) uint64 {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, timestamp)
	hj := metro.Hash64(timeBytes, hashid)
	return metro.Hash64(item, hj) % uint64(len(cms.regs[0]))
}

func (cms *CMS) factor(timestamp uint64) float64 {
	return math.Pow(cms.alpha, 1/float64(timestamp))
}
