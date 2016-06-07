package ada

import (
	"errors"
	"fmt"
	"math"

	"github.com/dgryski/go-farm"
)

// CountMinSketch ...
type CountMinSketch struct {
	t uint64
	w uint64
	d uint64
	m [][]uint64
}

// NewAdaCountMinSketch ...
func NewAdaCountMinSketch(t uint64, delta, epsilon float64) (*CountMinSketch, error) {
	if epsilon <= 0 || epsilon >= 1 {
		return nil, errors.New("countminsketch: value of epsilon should be in range of (0, 1)")
	}
	if delta <= 0 || delta >= 1 {
		return nil, errors.New("countminsketch: value of delta should be in range of (0, 1)")
	}

	w := uint64(math.Ceil(math.E / epsilon))
	d := uint64(math.Ceil(math.Log(1 / delta)))
	m := make([][]uint64, d)
	for i := range m {
		m[i] = make([]uint64, w)
	}

	fmt.Println(w, d)

	return &CountMinSketch{
		t: t,
		w: w,
		d: d,
		m: m,
	}, nil
}

// Update ...
func (cms *CountMinSketch) Update(value []byte, timestamp uint64) {
	ft := functTime(timestamp, cms.t)
	hsum := farm.Hash64(value)
	h1 := uint32(hsum & 0xffffffff)
	h2 := uint32((hsum >> 32) & 0xffffffff)

	for j := range cms.m {
		h := (uint64(h1+uint32(j)*h2) * ft) % cms.w
		// cjt is always 1 so just add ft
		cms.m[j][h] += ft
	}
}

// Query ...
func (cms *CountMinSketch) Query(value []byte, timestamp uint64) uint64 {
	ft := functTime(timestamp, cms.t)
	hsum := farm.Hash64(value)
	h1 := uint32(hsum & 0xffffffff)
	h2 := uint32((hsum >> 32) & 0xffffffff)

	var res uint64 = math.MaxUint64

	for j := range cms.m {
		h := (uint64(h1+uint32(j)*h2) * ft) % cms.w
		temp := cms.m[j][h] / ft
		if temp < res {
			res = temp
		}
	}
	return res
}
