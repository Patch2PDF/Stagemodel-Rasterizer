package rasterizer

import (
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
)

type fixtureLabel struct {
	fixture              *MVRTypes.Fixture
	fixture_bounding_box boundingBox
}

var fixture_label_padding = padding{5, 10, 5, 10}

func drawFixtureLabels(canvas *Canvas) error {
	for _, label := range canvas.fixture_labels {
		bb := label.fixture_bounding_box
		label_text := label.fixture.FixtureID

		width, height := calcLabelDimensions(canvas, label_text)
		_ = height

		// TODO: using zbuffer implement label positioning logic so that there are no overlaps
		x := (bb.right-bb.left-width)/2 + bb.left
		// y := (bb.bottom-bb.top-height)/2 + bb.top
		y := bb.bottom + 20

		rect, err := getAndCheckLabelRect(canvas, x, y, width, height, fixture_label_padding)
		if err != nil {
			// TODO: once logic for auto positioning is there, reactivate error forwarding
			// return err
			continue
		}

		drawLabelBackground(canvas, rect)
		labelFillZBuf(canvas, rect)

		drawLabelText(canvas, x, y, label_text)
	}
	return nil
}
