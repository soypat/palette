package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"image/color"
	"io"
	"log"
	"os"
	"strings"
)

//go:embed _template.go
var paletteTemplate string

var (
	inputName  = flag.String("input", "palettes.txt", "Name of file with palettes")
	outputName = flag.String("o", "palettes.go", "Resulting file")
)

func main() {
	flag.Parse()

	ifp, err := os.Open(*inputName)
	if err != nil {
		log.Fatal(err)
	}
	defer ifp.Close()
	ofp, err := os.Create(*outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer ofp.Close()

	p := Parameters{
		Template: paletteTemplate,
		Input:    ifp,
		Output:   ofp,
	}
	err = Execute(p)
	if err != nil {
		log.Fatal(err)
	}

}

type Parameters struct {
	Template string
	Input    io.Reader
	Output   io.Writer
}

type Pal struct {
	Name   string
	Doc    string
	Colors []color.Color
}

func (p Pal) Format(idx int) string {
	col := p.Colors[idx]
	r, g, b, a := col.RGBA()
	return fmt.Sprintf("color.RGBA64{R: %#x, G: %#x, B: %#x, A: %#x},\n", r, g, b, a)
}

func Execute(p Parameters) error {
	palettes := GenPalettes(p.Input)
	tpl, err := template.New("").Parse(p.Template)
	if err != nil {
		return err
	}
	return tpl.Execute(p.Output, palettes)
}

func GenPalettes(rd io.Reader) []Pal {
	sc := bufio.NewScanner(rd)
	var palettes []Pal
	var currentPal Pal
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0][0] <= 'Z' && fields[0][0] >= 'A' {
			currentPal.Name = fields[0]
			currentPal.Doc = line
			continue
		}
		color, err := parseGenericColor(line)
		if err != nil && currentPal.Name != "" {
			palettes = append(palettes, currentPal)
			currentPal = Pal{}
		} else {
			currentPal.Colors = append(currentPal.Colors, color)
		}
	}
	if currentPal.Name != "" {
		// If last parsed not appended finish job.
		palettes = append(palettes, currentPal)
		currentPal = Pal{}
	}
	return palettes
}
