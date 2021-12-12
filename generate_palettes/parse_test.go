package main

import (
	"image/color"
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	rd := strings.NewReader(testInput)
	palettes := GenPalettes(rd)
	t.Error(palettes)
}

func TestParseGenericColor(t *testing.T) {
	testCases := []struct {
		got    string
		expect color.Color
	}{
		{
			got:    "[.2,.2,.1]",
			expect: rgbf{r: .2, g: .2, b: .1},
		},
	}
	for _, tC := range testCases {
		got, err := parseGenericColor(tC.got)
		if err != nil {
			t.Error(err)
		} else if !colorEqual(got, tC.expect) {
			t.Errorf("color %s got unexpected value during parse", tC.got)
			t.Error(got.RGBA())
			t.Error(tC.expect.RGBA())
		}
	}

}

func colorEqual(a, b color.Color) bool {
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := a.RGBA()
	if ra != rb || ga != gb || ba != bb || aa != ab {
		return false
	}
	return true
}

const testInput = `Matlab taken from http://math.loyola.edu/~loberbro/matlab/html/colorsInMatlab.html
[0, 0.4470, 0.7410]
[0.8500, 0.3250, 0.0980]
[0.9290, 0.6940, 0.1250]
[0.4940, 0.1840, 0.5560]
[0.4660, 0.6740, 0.1880]
[0.3010, 0.7450, 0.9330]
[0.6350, 0.0780, 0.1840]]

MatlabOld taken from http://math.loyola.edu/~loberbro/matlab/html/colorsInMatlab.html
[0, 0, 1]
[0, 0.5, 0]
[1, 0, 0]
[0, 0.75, 0.75]
[0.75, 0, 0.75]
[0.75, 0.75, 0]
[0.25, 0.25, 0.25]`
