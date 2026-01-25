package rasterizer

import (
	"image/color"
	"image/draw"
	"math"
)

func line(ax int, ay int, bx int, by int, canvas draw.Image, color color.Color) {
	steep := math.Abs(float64(ax-bx)) < math.Abs(float64(ay-by))
	if steep { // if the line is steep, we transpose the image
		ax, ay = swap(ax, ay)
		bx, by = swap(bx, by)
	}
	if ax > bx { // make it left−to−right
		ax, bx = swap(ax, bx)
		ay, by = swap(ay, by)
	}

	var y float64 = float64(ay)
	for x := int(ax); x <= bx; x++ {
		if steep { // if transposed, de−transpose
			canvas.Set(int(math.Round(y)), x, color)
		} else {
			canvas.Set(x, int(math.Round(y)), color)
		}
		y += float64(by-ay) / float64(bx-ax)
	}
}
