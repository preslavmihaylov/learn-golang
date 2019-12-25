package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	newImageWithHistograms(800, 600, []int{70, 30, 60, 80})
}

func newImageWithHistograms(width, height int, barHeights []int) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	barWidth := width / (len(barHeights) + 1)
	offsetWidth := barWidth / (len(barHeights) + 1)

	bars := []image.Rectangle{}
	for i := range barHeights {
		offset := (i + 1) * offsetWidth
		topLeft := image.Point{offset + i*barWidth, height - height*barHeights[i]/100}
		botRight := image.Point{offset + (i+1)*barWidth, height}

		bars = append(bars, image.Rectangle{topLeft, botRight})
	}

	cyan := color.RGBA{100, 200, 200, 0xff}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
			for _, bar := range bars {
				if image.Pt(x, y).In(bar) {
					img.Set(x, y, cyan)
				}
			}
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, img)
}
