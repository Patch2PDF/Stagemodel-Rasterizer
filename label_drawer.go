package rasterizer

import (
	_ "embed"
	"fmt"
	"image"
	"strings"

	"image/color"
	"image/draw"
	_ "image/jpeg"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type fontFace struct {
	face        *font.Face
	ascent      int
	descent     int
	font_height int
	drawer      *font.Drawer
}

type padding struct {
	top    int
	right  int
	bottom int
	left   int
}

func initFontFace(size float64, dpi float64, fontBytes []byte) (*fontFace, error) {
	parsed_font, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font bytes: %v", err)
	}

	face, err := opentype.NewFace(parsed_font, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create new face: %v", err)
	}

	metrics := face.Metrics()

	font_face := &fontFace{
		face:    &face,
		ascent:  metrics.Ascent.Ceil(),
		descent: metrics.Descent.Ceil(),
	}

	font_face.font_height = font_face.ascent + font_face.descent

	return font_face, nil
}

func initFontDrawer(font_face *fontFace, img *image.NRGBA, start_x int, start_y int) {
	font_face.drawer = &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.NRGBA{0, 0, 0, 255}),
		Face: *font_face.face,
		Dot:  fixed.Point26_6{X: fixed.I(start_x), Y: fixed.I(start_y + font_face.ascent)},
	}
}

func calcLabelDimensions(canvas *Canvas, text string) (width int, height int) {
	font_face := canvas.label_face

	var line_count int = 0

	// calc text width
	for _, line := range strings.Split(text, "\n") {
		width = font_face.drawer.MeasureString(line).Ceil()
		line_count += 1
	}

	return width, font_face.font_height * line_count
}

func drawLabelBackground(canvas *Canvas, rect image.Rectangle, background_color color.NRGBA, border_width int, border_color color.NRGBA) {
	draw.Draw(
		canvas.canvas,
		image.Rectangle{
			Min: image.Point{rect.Min.X + border_width, rect.Min.Y + border_width},
			Max: image.Point{rect.Max.X - border_width, rect.Max.Y - border_width},
		},
		image.NewUniform(background_color),
		image.Point{},
		draw.Src,
	)

	frame(canvas.canvas, rect, border_width, border_color)
}

func drawLabelText(canvas *Canvas, x int, y int, text string) {
	font_face := canvas.label_face

	font_face.drawer.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y + font_face.ascent)}

	fixed_height := fixed.I(font_face.font_height)

	for _, line := range strings.Split(text, "\n") {
		font_face.drawer.DrawString(line)
		font_face.drawer.Dot = fixed.Point26_6{X: fixed.I(x), Y: font_face.drawer.Dot.Y + fixed_height}
	}
}

func getAndCheckLabelRect(canvas *Canvas, x int, y int, width int, height int) (image.Rectangle, error) {
	img := canvas.canvas
	img_rect := img.Bounds()

	// label border box
	tex_rect := image.Rect(x, y, x+width, y+height)
	if tex_rect.Min.X < img_rect.Min.X || tex_rect.Min.Y < img_rect.Min.Y || tex_rect.Max.X >= img_rect.Max.X || tex_rect.Max.Y >= img_rect.Max.Y {
		return image.Rectangle{}, fmt.Errorf(
			"Label exceeds image bounds (label: (%d|%d),(%d|%d); img: (%d|%d),(%d|%d))",
			tex_rect.Min.X, tex_rect.Min.Y,
			tex_rect.Max.X, tex_rect.Max.Y,
			img_rect.Min.X, img_rect.Min.Y,
			img_rect.Max.X, img_rect.Max.Y,
		)
	}
	return tex_rect, nil
}
