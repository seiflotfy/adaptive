package adaptive

import (
	"encoding/binary"
	"math"

	metro "github.com/dgryski/go-metro"
)

// ACMS is a simple count-min-sketch imeplementation with a exp based pre-emphasis and de-emphasis of values
type ACMS struct {
	regs  [][]float64
	w     uint64
	d     uint64
	alpha float64
}

// NewACMS returns a Sketch of width 2^w and depth d for a given alpha
// the alpha is ussed to pre-emphasis and de-emphasis the counters using alpha^timestamp
func NewACMS(w, d uint64, alpha float64) *ACMS {
	regs := make([][]float64, d)
	for i := range regs {
		regs[i] = make([]float64, uint64(math.Pow(2, float64(w))))
	}
	return &ACMS{
		w:     w,
		d:     d,
		regs:  regs,
		alpha: alpha,
	}
}

// Insert item at given timestamp n times (n = count)
func (acms *ACMS) Insert(item []byte, timestamp, count uint64) {
	for i := range acms.regs {
		j := acms.hash(item, timestamp, uint64(i))
		acms.regs[i][j] += (acms.factor(timestamp) * float64(count))
	}
}

// Estimate how many times item appeared at given timestamp
func (acms *ACMS) Estimate(item []byte, timestamp uint64) uint64 {
	min := math.MaxFloat64
	for i := range acms.regs {
		j := acms.hash(item, timestamp, uint64(i))
		if acms.regs[i][j] < min {
			min = acms.regs[i][j]
		}
	}
	return uint64(min / acms.factor(timestamp))
}

func (acms *ACMS) hash(item []byte, timestamp, hashid uint64) uint64 {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, timestamp)
	hj := metro.Hash64(timeBytes, hashid)
	return metro.Hash64(item, hj) % uint64(len(acms.regs[0]))
}

func (acms *ACMS) factor(timestamp uint64) float64 {
	return math.Pow(acms.alpha, 1/float64(timestamp))
}
