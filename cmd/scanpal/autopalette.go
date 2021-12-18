package main

import (
	"image"
	"image/color"
	"sort"
)

func autoPalette(im image.Image) color.Palette {
	rect := im.Bounds()
	buckets := NewBucket()
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			buckets.Add(im.At(x, y))
		}
	}
	buckets.Filter(*palfilter)
	cpal := buckets.Palette()
	sort.Sort(ByDirection{
		cpal: cpal,
	})
	return cpal
}

type ByDirection struct {
	midpoint, vdir Vec
	cpal           color.Palette
}

// Len is the number of elements in the collection.
func (h ByDirection) Len() int {
	n := len(h.cpal)
	// First find color center of gravity or midpoint.
	var midpoint Vec
	for i := range h.cpal {
		v := h.Vec(i)
		midpoint = Add(v, midpoint)
	}
	h.midpoint = Scale(1/float64(n), midpoint)

	// Next find principal direction.
	var maxDist float64
	for i := range h.cpal {
		v := h.Vec(i)
		vdir := Sub(v, midpoint)
		vnorm := Norm(vdir)
		if vnorm > maxDist {
			maxDist = vnorm
			h.vdir = vdir
		}
	}
	return n
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (h ByDirection) Less(i, j int) bool {
	v1, v2 := h.Vec(i), h.Vec(j)
	a := Dot(Sub(v1, h.midpoint), h.vdir)
	b := Dot(Sub(v2, h.midpoint), h.vdir)
	return a < b
}

// Swap swaps the elements with indexes i and j.
func (h ByDirection) Swap(i, j int) {
	h.cpal[i], h.cpal[j] = h.cpal[j], h.cpal[i]
}

func (h ByDirection) Vec(i int) Vec {
	r, g, b, _ := h.cpal[i].RGBA()
	return Vec{X: float64(r), Y: float64(g), Z: float64(b)}
}
