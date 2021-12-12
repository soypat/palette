package main

import (
	"errors"
	"image/color"
	"math"
	"strconv"
)

// parseGenericColor is the main color parsing function.
func parseGenericColor(s string) (color.Color, error) {
	floats := scanFloats(s)
	if len(floats) == 3 && floats[0] <= 1 {
		return rgbf{r: floats[0], g: floats[1], b: floats[2]}, nil
	}

	return nil, errors.New("unable to parse color " + s)
}

// Floats contained in [0,1]
type rgbf struct {
	r, g, b float64
}

func (cf rgbf) RGBA() (r, g, b, a uint32) {
	col := color.RGBA64{
		R: uint16(cf.r * math.MaxUint16),
		G: uint16(cf.g * math.MaxUint16),
		B: uint16(cf.b * math.MaxUint16),
		A: math.MaxUint16,
	}
	return col.RGBA()
}

func scanFloats(s string) (nums []float64) {
	start := -1
	for i, c := range s {
		isNum := isNumTok(c)
		if (isNum || c == '.') && start < 0 {
			start = i
			if i > 0 && s[i-1] == '-' {
				start--
			}
		}
		isTok := isFloatTok(c)
		if (start >= 0 && !isTok) || (i == len(s)-1 && start >= 0) {
			if isNum {
				// Include number if at end of string.
				i++
			}
			num, err := strconv.ParseFloat(s[start:i], 64)
			if err == nil {
				nums = append(nums, num)
			}
			start = -1
		}
	}
	return nums
}

func isNumTok(r rune) bool {
	return r^'0' < 10
}

func isFloatTok(r rune) bool {
	return isNumTok(r) || r == '.' || r == 'E' || r == '+' || r == '-' || r == 'e'
}
