package rasterizer

import (
	"fmt"
	"image"
	"image/color"

	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
)

type fixtureLabel struct {
	fixture              *MVRTypes.Fixture
	fixture_bounding_box boundingBox
}

var fixture_label_padding = padding{5, 10, 5, 10}

const fixture_distance = 15

var fixture_label_background = color.NRGBA{255, 255, 255, 255}

const fixture_label_border_width = 3

var fixture_label_border_color = color.NRGBA{0, 0, 0, 255}

// candidates for label positioning (order in array decides priority)
// if no one fits, fall back to more expensive algorithm
var candidates = []string{
	"bottom",
	"top",
	"right",
	"left",
	"bottom-right",
	"top-right",
	"top-left",
	"bottom-left",
}

//TODO: if no candidate fits, try Force-Directed/Spring, if that puts it too far away, only check for no label overlap?

func fixtureLabelFillZBuf(canvas *Canvas, rect image.Rectangle) {
	for y_index := rect.Min.Y; y_index < rect.Max.Y; y_index++ {
		zBufRowIndex := (y_index) * canvas.width
		for x_index := rect.Min.X; x_index < rect.Max.X; x_index++ {
			canvas.fixture_zbuf[zBufRowIndex+x_index] = true
		}
	}
}

// check if rect space is already occupied, returns true if blocked, false if free
func fixtureLabelCheckZBufOccupied(canvas *Canvas, rect image.Rectangle) bool {
	for y_index := rect.Min.Y; y_index < rect.Max.Y; y_index++ {
		zBufRowIndex := (y_index) * canvas.width
		for x_index := rect.Min.X; x_index < rect.Max.X; x_index++ {
			if canvas.fixture_zbuf[zBufRowIndex+x_index] {
				return true
			}
		}
	}
	return false
}

func fixtureLabelGetCandidatePosition(canvas *Canvas, boundingBox boundingBox, text_width int, text_height int, candidate string, padding padding) (image.Rectangle, error) {
	width := text_width + padding.left + padding.right
	height := text_height + padding.top + padding.bottom
	switch candidate {
	case "bottom":
		x := (boundingBox.right-boundingBox.left-width)/2 + boundingBox.left
		y := boundingBox.bottom + fixture_distance + padding.top
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "right":
		x := boundingBox.right + fixture_distance + padding.left
		y := (boundingBox.bottom-boundingBox.top-height)/2 + boundingBox.top
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "top":
		x := (boundingBox.right-boundingBox.left-width)/2 + boundingBox.left
		y := boundingBox.top - fixture_distance - height - padding.bottom
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "left":
		x := boundingBox.left - fixture_distance - width - padding.right
		y := (boundingBox.bottom-boundingBox.top-height)/2 + boundingBox.top
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "bottom-right":
		x := boundingBox.right + fixture_distance + padding.left
		y := boundingBox.bottom + fixture_distance + padding.top
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "top-right":
		x := boundingBox.right + fixture_distance + padding.left
		y := boundingBox.top - fixture_distance - height - padding.bottom
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "top-left":
		x := boundingBox.left - fixture_distance - width - padding.right
		y := boundingBox.top - fixture_distance - height - padding.bottom
		return getAndCheckLabelRect(canvas, x, y, width, height)
	case "bottom-left":
		x := boundingBox.left - fixture_distance - width - padding.right
		y := boundingBox.bottom + fixture_distance + padding.top
		return getAndCheckLabelRect(canvas, x, y, width, height)
	default:
		return image.Rectangle{}, fmt.Errorf("Unknown Position Candidate: %s", candidate)
	}
}

func fixtureLabelGetPosition(canvas *Canvas, boundingBox boundingBox, text_width int, text_height int, text string) (image.Rectangle, error) {
	for _, candidate := range candidates {
		rect, err := fixtureLabelGetCandidatePosition(canvas, boundingBox, text_width, text_height, candidate, fixture_label_padding)

		if err != nil {
			continue // TODO: implement better error handling to distinguish between unknown candidate and position just being outside canvas
		}

		if !fixtureLabelCheckZBufOccupied(canvas, rect) {
			return rect, nil
		}
	}
	return image.Rectangle{}, fmt.Errorf("Could not find a free position for %s", text)
}

func drawFixtureLabels(canvas *Canvas) error {
	for _, label := range canvas.fixture_labels {
		bb := label.fixture_bounding_box
		label_text := label.fixture.FixtureID

		width, height := calcLabelDimensions(canvas, label_text)
		_ = height

		rect, err := fixtureLabelGetPosition(canvas, bb, width, height, label_text)

		if err != nil {
			// return err
			continue
		}

		drawLabelBackground(canvas, rect, fixture_label_background, fixture_label_border_width, fixture_label_border_color)
		fixtureLabelFillZBuf(canvas, rect)

		drawLabelText(canvas, rect.Min.X+fixture_label_padding.left, rect.Min.Y+fixture_label_padding.top, label_text)
	}
	return nil
}
