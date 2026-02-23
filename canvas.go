package rasterizer

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math"
)

type Canvas struct {
	width   int
	height  int
	canvas  draw.Image
	zbuffer [][]float64
}

func (cv *Canvas) Init(width int, height int) {
	cv.width = width
	cv.height = height

	cv.canvas = image.NewRGBA(image.Rect(0, 0, width, height))

	cv.zbuffer = make([][]float64, height)
	for i := range cv.zbuffer {
		cv.zbuffer[i] = make([]float64, width)
		for j := range cv.zbuffer[i] {
			cv.zbuffer[i][j] = math.Inf(-1)
		}
	}
}

func (cv *Canvas) SaveAsPNG(w io.Writer) error {
	err := png.Encode(w, cv.canvas)
	if err != nil {
		return fmt.Errorf("Error Saving Canvas as PNG: %s", err)
	}
	return nil
}
