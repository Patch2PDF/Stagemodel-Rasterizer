package rasterizer

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

type Point struct {
	x int
	y int
}

type Triangle struct {
	a Point
	b Point
	c Point
}

func Draw() {
	var canvas draw.Image = image.NewRGBA(image.Rect(0, 0, 127, 127))

	scanTriangle(Triangle{a: Point{7, 45}, b: Point{35, 100}, c: Point{45, 60}}, canvas, color.RGBA{255, 0, 0, 255})
	scanTriangle(Triangle{a: Point{120, 35}, b: Point{90, 5}, c: Point{45, 110}}, canvas, color.RGBA{0, 255, 0, 255})
	scanTriangle(Triangle{a: Point{115, 83}, b: Point{80, 90}, c: Point{85, 120}}, canvas, color.RGBA{0, 0, 255, 255})

	f, err := os.Create("render.png")
	if err != nil {
		log.Fatalf("failed to open render: %v", err)
	}

	png.Encode(f, canvas)

	f.Close()
}
