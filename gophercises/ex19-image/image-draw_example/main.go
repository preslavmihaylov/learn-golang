package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	newImageWithHistograms(800, 600, []int{70, 30, 60, 80})
}

func newImageWithHistograms(width, height int, barHeights []int) {
	destImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	barWidth := width / (len(barHeights) + 1)
	offsetWidth := barWidth / (len(barHeights) + 1)

	draw.Draw(destImg, destImg.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	cyan := color.RGBA{100, 200, 200, 0xff}
	for i := range barHeights {
		offset := (i + 1) * offsetWidth
		topLeft := image.Point{offset + i*barWidth, height - height*barHeights[i]/100}
		botRight := image.Point{offset + (i+1)*barWidth, height}

		src := image.NewRGBA(image.Rectangle{topLeft, botRight})
		draw.Draw(destImg, src.Bounds(), &image.Uniform{cyan}, image.ZP, draw.Src)
	}

	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, destImg)
}
