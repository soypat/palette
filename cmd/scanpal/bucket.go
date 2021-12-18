package main

import (
	"image/color"
	"math"
)

func NewBucket() Bucket {
	return Bucket{
		idx: make([]uint8, 0, 256),
		b:   make(map[uint8]runningAvg),
	}
}

// Bucket gathers similar colors by appearance so as to
// generate palettes of an image.
type Bucket struct {
	b   map[uint8]runningAvg
	idx []uint8 // maps index to uint8 color representation
}

func (bu *Bucket) Add(c color.Color) {
	r, g, b, a := c.RGBA()
	if a != math.MaxUint16 {
		return // ignore non-opaque colors
	}
	cc := compress8(r, g, b)
	avg, found := bu.b[cc]
	if !found {
		bu.idx = append(bu.idx, cc)
	}
	avg.update(float64(r), float64(g), float64(b))
	bu.b[cc] = avg
}

func (bu *Bucket) Filter(frac float64) {
	if frac <= 0 || frac > 1 {
		panic("bad frac")
	}
	var maxOcurrences, meanOcurrences float64
	for _, avg := range bu.b {
		n := avg.n
		if n > maxOcurrences {
			maxOcurrences = n
		}
		meanOcurrences += n
	}
	meanOcurrences /= float64(len(bu.b))

	for cc, avg := range bu.b {
		if avg.n < meanOcurrences*frac {
			delete(bu.b, cc)
		}
	}
}

func (bu Bucket) Palette() color.Palette {
	n := len(bu.b)
	cpal := make(color.Palette, 0, n)
	for _, cc := range bu.idx {
		avg, ok := bu.b[cc]
		if !ok { // may have been filtered out
			continue
		}
		if len(cpal) == n {
			break
		}
		cpal = append(cpal, avg.rgba64())
	}
	return cpal
}

type runningAvg struct {
	r, g, b, n float64
}

func (a *runningAvg) update(r, g, b float64) {
	a.n++ // welfords moving average
	a.r = a.r + (r-a.r)/a.n
	a.g = a.g + (g-a.g)/a.n
	a.b = a.b + (b-a.b)/a.n
}

func (a *runningAvg) rgba64() color.RGBA64 {
	return color.RGBA64{
		R: uint16(a.r),
		G: uint16(a.g),
		B: uint16(a.b),
		A: maxColor,
	}
}

func compress8(r, g, b uint32) (eightbitcolor uint8) {
	rpart := uint8(r>>(imagebits-3)) << 5
	gpart := uint8(g>>(imagebits-3)) << 2
	bpart := uint8(b >> (imagebits - 2))
	return rpart | gpart | bpart
}
