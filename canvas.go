package rasterizer

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
)

type Canvas struct {
	width   int
	height  int
	canvas  *image.NRGBA
	zbuffer []float64
}

func (cv *Canvas) Init(width int, height int) {
	cv.width = width
	cv.height = height

	cv.canvas = image.NewNRGBA(image.Rect(0, 0, width, height))

	cv.zbuffer = make([]float64, height*width)
	for i := range cv.zbuffer {
		cv.zbuffer[i] = math.Inf(-1)
	}
}

func (cv *Canvas) SaveAsPNG(w io.Writer) error {
	encoder := png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	err := encoder.Encode(w, cv.canvas)
	if err != nil {
		return fmt.Errorf("Error Saving Canvas as PNG: %s", err)
	}
	return nil
}
