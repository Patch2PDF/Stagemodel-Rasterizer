package rasterizer

import (
	_ "embed"
	"fmt"
	"image"
	"os"
	"strings"

	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	//go:embed fonts/DejaVuSansMono.ttf
	fontBytes []byte
)

var face font.Face

const dpi = 72
const size = 18

var font_height int
var ascent int
var descent int

func InitFace() error {
	parsed_font, err := opentype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse font bytes: %v", err)
	}
	face, err = opentype.NewFace(parsed_font, &opentype.FaceOptions{
		Size:    size, // float64(img.Bounds().Dx() / 5),
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return fmt.Errorf("failed to create new face: %v", err)
	}

	metrics := face.Metrics()
	ascent = metrics.Ascent.Ceil()
	descent = metrics.Descent.Ceil()
	font_height = ascent + descent

	return nil
}

func DrawLabelBox(img *image.NRGBA, x int, y int, text string) (width int, height int, err error) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.NRGBA{0, 0, 0, 255}),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y + ascent)},
	}

	var line_count int = 0

	// calc text width
	for _, line := range strings.Split(text, "\n") {
		width = d.MeasureString(line).Ceil()
		line_count += 1
	}

	// label border box
	tex_rect := image.Rect(x, y, x+width, y+font_height*line_count)
	if tex_rect.Min.X < img.Rect.Min.X || tex_rect.Min.Y < img.Rect.Min.Y || tex_rect.Max.X > img.Rect.Max.X || tex_rect.Max.Y > img.Rect.Max.Y {
		return 0, 0, fmt.Errorf(
			"Label exceeds image bounds (label: (%d|%d),(%d|%d); img: (%d|%d),(%d|%d))",
			tex_rect.Min.X, tex_rect.Min.Y,
			tex_rect.Max.X, tex_rect.Max.Y,
			img.Rect.Min.X, img.Rect.Min.Y,
			img.Rect.Max.X, img.Rect.Max.Y,
		)
	}

	// set background
	draw.Draw(img, tex_rect, image.NewUniform(color.White), image.Point{}, draw.Src)

	fixed_height := fixed.I(font_height)

	for _, line := range strings.Split(text, "\n") {
		d.DrawString(line)
		d.Dot = fixed.Point26_6{X: fixed.I(x), Y: d.Dot.Y + fixed_height}
	}

	size := tex_rect.Size()

	return size.X, size.Y, nil
}

func DrawLabel(x int, y int) error {
	InitFace()

	dst := image.NewNRGBA(image.Rect(0, -25, 50, 25))

	dst.Set(x, y, color.NRGBA{0, 255, 0, 255})

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.NRGBA{255, 255, 255, 255}),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y + ascent)},
	}
	fmt.Printf("%s", d.Dot)
	width := d.MeasureString("LgTM")
	d.DrawString("LgTM")
	fmt.Printf("%s", d.Dot)

	// metrics := d.Face.Metrics()
	// height2 := metrics.Ascent.Ceil()

	dst.Set(x, y+font_height, color.NRGBA{255, 0, 0, 255})
	// dst.Set(x, y+height2, color.NRGBA{255, 0, 255, 255})
	dst.Set(x+width.Ceil(), y+font_height, color.NRGBA{0, 0, 255, 255})
	// dst.Set(x+width.Ceil(), y+height2, color.NRGBA{255, 255, 0, 255})

	f, err := os.Create("out.png")
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, dst); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	return nil
}
