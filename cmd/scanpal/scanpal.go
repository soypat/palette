package main

import (
	"flag"
	"fmt"
	"image"

	"image/color"
	"image/png"
	"log"
	"os"

	"gonum.org/v1/gonum/spatial/r3"
)

var (
	target = flag.String("target", "", "file target expect png.")
	dir    = flag.String("spectrum", "x", "For color scales. Set to color scale direction. Expects 'x' or 'y'. If none provided will try naively generating pallete.")
)

const (
	imagebits = 16
	maxColor  = (1 << imagebits) - 1
)

func main() {
	flag.Parse()
	if *target == "" {
		flag.CommandLine.Usage()
		log.Fatal("Please pass a target png file.\n")
	}
	fp, err := os.Open(*target)
	if err != nil {
		log.Fatal(err)
	}
	im, err := png.Decode(fp)
	if err != nil {
		log.Fatal(err)
	}
	var cpal color.Palette
	switch *dir {
	case "x":
		cpal = xPalette(im)
	case "y":
		cpal = yPalette(im)
	case "":
		cpal = autoPalette(im)
	default:
		log.Fatal("unknown spectrum argument ", *dir)
	}

	savePalette(cpal, 100, "palette.png")
	for i := range cpal {
		fmt.Printf("%.6f\n", toFloat(cpal[i]))
	}
}

func xPalette(im image.Image) color.Palette {
	rect := im.Bounds()
	y := (rect.Max.Y - rect.Min.Y) / 2
	buckets := make(map[uint8]struct{})
	var cpal color.Palette
	for x := rect.Min.X; x < rect.Max.X; x++ {
		rbig, gbig, bbig, abig := im.At(x, y).RGBA()
		if abig != maxColor {
			// skip transparent pixels
			continue
		}
		col := compress8(rbig, gbig, bbig)
		if _, ok := buckets[col]; !ok {
			buckets[col] = struct{}{}
			cpal = append(cpal, color.RGBA64{R: uint16(rbig), G: uint16(gbig), B: uint16(bbig), A: maxColor})
		}
	}
	return cpal
}

func yPalette(im image.Image) color.Palette {
	rect := im.Bounds()
	x := (rect.Max.X - rect.Min.X) / 2
	buckets := make(map[uint8]struct{})
	var cpal color.Palette
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		rbig, gbig, bbig, abig := im.At(x, y).RGBA()
		if abig != maxColor {
			// skip transparent pixels
			continue
		}
		col := compress8(rbig, gbig, bbig)
		if _, ok := buckets[col]; !ok {
			buckets[col] = struct{}{}
			cpal = append(cpal, color.RGBA64{R: uint16(rbig), G: uint16(gbig), B: uint16(bbig), A: maxColor})
		}
	}
	return cpal
}

func savePalette(cpal color.Palette, size int, name string) {
	// n := len(cpal)
	fp, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	s := palimg{
		size: size,
		cpal: cpal,
	}

	err = png.Encode(fp, s)
	if err != nil {
		panic(err)
	}
}

func compress8(r, g, b uint32) (eightbitcolor uint8) {
	rpart := uint8(r>>(imagebits-3)) << 5
	gpart := uint8(g>>(imagebits-3)) << 2
	bpart := uint8(b >> (imagebits - 2))
	return rpart | gpart | bpart
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

type palimg struct {
	size int
	cpal color.Palette
}

// ColorModel returns the Image's color model.
func (pi palimg) ColorModel() color.Model {
	return color.RGBA64Model
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (pi palimg) Bounds() image.Rectangle {
	return image.Rect(0, 0, pi.size*len(pi.cpal), pi.size)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (pi palimg) At(x, y int) color.Color {
	return pi.cpal[x/pi.size]
}

func toFloat(c color.Color) r3.Vec {
	r, g, b, _ := c.RGBA()
	return r3.Vec{X: float64(r) / float64(maxColor), Y: float64(g) / float64(maxColor), Z: float64(b) / float64(maxColor)}
}
