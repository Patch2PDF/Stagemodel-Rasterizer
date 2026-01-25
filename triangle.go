package rasterizer

import (
	"image/color"
	"image/draw"
	"math"
)

// func scanTriangle(ax int, ay int, bx int, by int, cx int, triangle.c.y int, canvas draw.Image, color color.Color) {
func scanTriangle(triangle Triangle, canvas draw.Image, color color.Color) {
	// sort the vertices, a,b,c in ascending y order (bubblesort yay!)
	if triangle.a.y > triangle.b.y {
		triangle.a.x, triangle.b.x = swap(triangle.a.x, triangle.b.x)
		triangle.a.y, triangle.b.y = swap(triangle.a.y, triangle.b.y)
	}
	if triangle.a.y > triangle.c.y {
		triangle.a.x, triangle.c.x = swap(triangle.a.x, triangle.c.x)
		triangle.a.y, triangle.c.y = swap(triangle.a.y, triangle.c.y)
	}
	if triangle.b.y > triangle.c.y {
		triangle.b.x, triangle.c.x = swap(triangle.b.x, triangle.c.x)
		triangle.b.y, triangle.c.y = swap(triangle.b.y, triangle.c.y)
	}
	total_height := triangle.c.y - triangle.a.y

	if triangle.a.y != triangle.b.y { // if the bottom half is not degenerate
		segment_height := triangle.b.y - triangle.a.y
		for y := triangle.a.y; y <= triangle.b.y; y++ { // sweep the horizontal line from ay to by
			x1 := triangle.a.x + ((triangle.c.x-triangle.a.x)*(y-triangle.a.y))/total_height
			x2 := triangle.a.x + ((triangle.b.x-triangle.a.x)*(y-triangle.a.y))/segment_height
			for x := math.Min(float64(x1), float64(x2)); x < math.Max(float64(x1), float64(x2)); x++ { // draw a horizontal line
				canvas.Set(int(x), y, color)
			}
		}
	}
	if triangle.b.y != triangle.c.y { // if the upper half is not degenerate
		segment_height := triangle.c.y - triangle.b.y
		for y := triangle.b.y; y <= triangle.c.y; y++ { // sweep the horizontal line from by to cy
			x1 := triangle.a.x + ((triangle.c.x-triangle.a.x)*(y-triangle.a.y))/total_height
			x2 := triangle.b.x + ((triangle.c.x-triangle.b.x)*(y-triangle.b.y))/segment_height
			for x := math.Min(float64(x1), float64(x2)); x < math.Max(float64(x1), float64(x2)); x++ { // draw a horizontal line
				canvas.Set(int(x), y, color)
			}
		}
	}
}
