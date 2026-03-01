package rasterizer

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
	GDTFTypes "github.com/Patch2PDF/GDTF-Parser/pkg/types"
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
)

var colors = map[GDTFTypes.GeometryType]color.NRGBA{
	GDTFTypes.GeometryTypeGeometry:          {25, 25, 25, 255},
	GDTFTypes.GeometryTypeAxis:              {25, 25, 25, 255},
	GDTFTypes.GeometryTypeFilterBeam:        {0, 0, 0, 255},
	GDTFTypes.GeometryTypeFilterColor:       {0, 0, 0, 255},
	GDTFTypes.GeometryTypeFilterGobo:        {0, 0, 0, 255},
	GDTFTypes.GeometryTypeFilterShaper:      {0, 0, 0, 255},
	GDTFTypes.GeometryTypeBeam:              {200, 200, 200, 255},
	GDTFTypes.GeometryTypeMediaServerLayer:  {0, 0, 0, 255},
	GDTFTypes.GeometryTypeMediaServerCamera: {0, 0, 0, 255},
	GDTFTypes.GeometryTypeMediaServerMaster: {0, 0, 0, 255},
	GDTFTypes.GeometryTypeDisplay:           {0, 0, 0, 255},
	GDTFTypes.GeometryTypeGeometryReference: {0, 0, 0, 255},
	GDTFTypes.GeometryTypeLaser:             {0, 0, 0, 255},
	GDTFTypes.GeometryTypeWiringObject:      {0, 0, 0, 255},
	GDTFTypes.GeometryTypeInventory:         {0, 0, 0, 255},
	GDTFTypes.GeometryTypeStructure:         {0, 0, 0, 255},
	GDTFTypes.GeometryTypeSupport:           {0, 0, 0, 255},
	GDTFTypes.GeometryTypeMagnet:            {0, 0, 0, 255},
}

type Rotation struct {
	Alpha float64
	Beta  float64
	Gamma float64
}

type boundingBox struct {
	left   int
	top    int
	right  int
	bottom int
}

func (b *boundingBox) init() {
	b.left = math.MaxInt
	b.top = math.MaxInt
	b.right = math.MinInt
	b.bottom = math.MinInt
}

func drawMesh(mesh MeshTypes.Mesh, canvas *Canvas, color color.NRGBA) {
	for _, triangle := range mesh.Triangles {
		NewTriangleFromMeshTriangle(triangle).boundingTriangle(
			canvas,
			color,
		)
	}
}

func drawMeshUpdateBB(mesh MeshTypes.Mesh, canvas *Canvas, color color.NRGBA, bb boundingBox) (boundingBox, error) {
	triangle_drawn := false

	for _, triangle := range mesh.Triangles {
		bbminx, bbminy, bbmaxx, bbmaxy, err := NewTriangleFromMeshTriangle(triangle).boundingTriangle(
			canvas,
			color,
		)
		if err == nil {
			triangle_drawn = true
			bb.left = min(bb.left, bbminx)
			bb.top = min(bb.top, bbminy)
			bb.right = max(bb.right, bbmaxx)
			bb.bottom = max(bb.bottom, bbmaxy)
		}
	}

	if triangle_drawn {
		return bb, nil
	} else {
		return bb, fmt.Errorf("No triangle was drawn")
	}
}

func drawStageModel(mesh *MVRTypes.StageModel, canvas *Canvas) {
	for _, obj := range mesh.SceneObjectModels {
		for _, part := range obj.MeshModel {
			drawMesh(part.Mesh, canvas, colors[part.GeometryType]) // TODO: obj type specific colors
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, color.NRGBA{100, 100, 100, 255})
		}
	}

	for _, fixture := range mesh.FixtureModels {
		bb := boundingBox{}
		bb.init()

		for _, part := range fixture.MeshModel {
			bb, _ = drawMeshUpdateBB(part.Mesh, canvas, colors[part.GeometryType], bb)
		}

		canvas.fixture_labels = append(
			canvas.fixture_labels,
			fixtureLabel{fixture: fixture.Fixture, fixture_bounding_box: bb},
		)

		for _, geometry := range fixture.Geometries {
			drawMesh(geometry, canvas, color.NRGBA{100, 100, 100, 255})
		}
	}
}

func Draw(mesh *MVRTypes.StageModel, rotation Rotation, filename string) (*Canvas, error) {
	const width = 4000
	const height = 3000

	canvas := &Canvas{}
	err := canvas.Init(width, height)

	if err != nil {
		return nil, err
	}

	normalizeAndRotateStageModel(canvas, mesh, rotation)

	drawStageModel(mesh, canvas)

	err = drawFixtureLabels(canvas)

	if err != nil {
		return nil, err
	}

	return canvas, nil
}

func SaveCanvasAsPNGFile(filename string, canvas *Canvas) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to open render: %v", err)
	}

	err = canvas.SaveAsPNG(f)

	if err != nil {
		log.Fatalf("%s", err)
	}

	f.Close()
}
