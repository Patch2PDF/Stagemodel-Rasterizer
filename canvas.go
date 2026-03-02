package rasterizer

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
)

type Canvas struct {
	width          int
	height         int
	canvas         *image.NRGBA
	zbuffer        []float64
	label_face     *fontFace
	fixture_labels []fixtureLabel
	fixture_zbuf   []bool // used for label placement to search for "free" pixels
}

func (cv *Canvas) Init(width int, height int) error {
	cv.width = width
	cv.height = height

	cv.canvas = image.NewNRGBA(image.Rect(0, 0, width, height))

	cv.fixture_zbuf = make([]bool, height*width)

	cv.zbuffer = make([]float64, height*width)
	for i := range cv.zbuffer {
		cv.zbuffer[i] = math.Inf(-1)
	}

	var err error

	cv.label_face, err = initFontFace(10, 300)
	if err != nil {
		return err
	}

	initFontDrawer(cv.label_face, cv.canvas, 0, 0)

	return nil
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
