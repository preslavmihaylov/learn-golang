package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	newImageWithHistograms(1200, 600,
		[]int{95, 100, 45, 40, 30, 20, 15, 25, 0, 0, 0, 0},
		[]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	)
}

func newImageWithHistograms(width, height int, barHeights []int, labels []string) {
	if len(barHeights) != len(labels) {
		panic("bar heights count different from labels count")
	}

	f, err := os.Create("image.svg")
	if err != nil {
		panic(err)
	}

	canvas := svg.New(f)
	canvas.Start(width, height)

	barWidth := width / (len(barHeights) + 1)
	offsetWidth := barWidth / (len(barHeights) + 1)
	offsetHeight := 50
	maxBarHeight := height - offsetHeight

	// draw line between bars & months
	canvas.Line(0, height-offsetHeight, width, height-offsetHeight, "stroke:gray;stroke-width:2")

	for i := range barHeights {
		offsetW := (i + 1) * offsetWidth
		offsetH := offsetHeight

		topLeftX, topLeftY := offsetW+i*barWidth, height-maxBarHeight*barHeights[i]/100-offsetH
		botRightX, botRightY := offsetW+(i+1)*barWidth, height-offsetH

		// draw bars & labels
		canvas.Rect(topLeftX, topLeftY, botRightX-topLeftX, botRightY-topLeftY, "fill:deepskyblue")
		canvas.Text(topLeftX+(botRightX-topLeftX)/2, botRightY+(offsetH/2), labels[i], "fill:gray;text-anchor:middle")
	}

	// draw warning line threshold
	warningPercent := 70
	topLeftX, topLeftY := 0, 0
	botRightX, botRightY := width, height-maxBarHeight*warningPercent/100-offsetHeight

	canvas.Rect(topLeftX, topLeftY, botRightX, botRightY, "fill:red;fill-opacity:0.2")
	canvas.Line(topLeftX, botRightY, botRightX, botRightY, "stroke:red;stroke-width:3")

	canvas.End()
}
