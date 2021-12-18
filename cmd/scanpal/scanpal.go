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
	target    = flag.String("target", "", "file target expect png.")
	dir       = flag.String("spectrum", "x", "For color scales. Set to color scale direction. Expects 'x' or 'y'. If none provided will try naively generating pallete with whole image.")
	palfilter = flag.Float64("cthreshold", 0.2, "filter out less dense colors in image")
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
	buckets := NewBucket()
	rect := im.Bounds()
	y := (rect.Max.Y - rect.Min.Y) / 2
	for x := rect.Min.X; x < rect.Max.X; x++ {
		buckets.Add(im.At(x, y))
	}
	buckets.Filter(*palfilter)
	return buckets.Palette()
}

func yPalette(im image.Image) color.Palette {
	buckets := NewBucket()
	rect := im.Bounds()
	x := (rect.Max.X - rect.Min.X) / 2
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		buckets.Add(im.At(x, y))
	}
	buckets.Filter(*palfilter)
	return buckets.Palette()
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
