package rasterizer

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
)

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
			defaultPerPixelPreCallback,
		)
	}
}

func perPixelPreCallbackUpdateFixtureZBuf(canvas *Canvas, x int, y int) bool {
	canvas.fixture_zbuf[y*canvas.width+x] = true
	return false
}

func drawMeshUpdateBB(mesh MeshTypes.Mesh, canvas *Canvas, color color.NRGBA, bb boundingBox) (boundingBox, error) {
	triangle_drawn := false

	for _, triangle := range mesh.Triangles {
		bbminx, bbminy, bbmaxx, bbmaxy, err := NewTriangleFromMeshTriangle(triangle).boundingTriangle(
			canvas,
			color,
			perPixelPreCallbackUpdateFixtureZBuf,
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

func drawStageModel(mesh *MVRTypes.StageModel, canvas *Canvas, config RasterizerConfig) error {
	colors := getOverrideColors(config.OverrideColors)
	for _, obj := range mesh.SceneObjectModels {
		for _, part := range obj.MeshModel {
			color, err := colors.getColor(ModelTypeSceneObject, part.GeometryType)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			drawMesh(part.Mesh, canvas, color)
		}
		geometry_color, err := colors.getGeometriesColor(ModelTypeSceneObject)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	for _, obj := range mesh.FocusPointModels {
		geometry_color, err := colors.getGeometriesColor(ModelTypeFocusPoint)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	if config.RenderLabels {
		for _, fixture := range mesh.FixtureModels {
			bb := boundingBox{}
			bb.init()

			for _, part := range fixture.MeshModel {
				color, err := colors.getColor(ModelTypeFixture, part.GeometryType)
				if err != nil {
					return fmt.Errorf("Could not draw stage model: %s", err)
				}
				bb, _ = drawMeshUpdateBB(part.Mesh, canvas, color, bb)
			}

			canvas.fixture_labels = append(
				canvas.fixture_labels,
				fixtureLabel{fixture: fixture.Fixture, fixture_bounding_box: bb},
			)

			geometry_color, err := colors.getGeometriesColor(ModelTypeFixture)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			for _, geometry := range fixture.Geometries {
				drawMesh(geometry, canvas, geometry_color)
			}
		}
	} else {
		for _, fixture := range mesh.FixtureModels {
			for _, part := range fixture.MeshModel {
				color, err := colors.getColor(ModelTypeFixture, part.GeometryType)
				if err != nil {
					return fmt.Errorf("Could not draw stage model: %s", err)
				}
				drawMesh(part.Mesh, canvas, color)
			}

			geometry_color, err := colors.getGeometriesColor(ModelTypeFixture)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			for _, geometry := range fixture.Geometries {
				drawMesh(geometry, canvas, geometry_color)
			}
		}
	}

	for _, obj := range mesh.SupportModels {
		for _, part := range obj.MeshModel {
			color, err := colors.getColor(ModelTypeSupport, part.GeometryType)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			drawMesh(part.Mesh, canvas, color)
		}
		geometry_color, err := colors.getGeometriesColor(ModelTypeSupport)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	for _, obj := range mesh.TrussModels {
		for _, part := range obj.MeshModel {
			color, err := colors.getColor(ModelTypeTruss, part.GeometryType)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			drawMesh(part.Mesh, canvas, color)
		}
		geometry_color, err := colors.getGeometriesColor(ModelTypeTruss)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	for _, obj := range mesh.VideoScreenModels {
		for _, part := range obj.MeshModel {
			color, err := colors.getColor(ModelTypeVideoScreen, part.GeometryType)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			drawMesh(part.Mesh, canvas, color)
		}
		geometry_color, err := colors.getGeometriesColor(ModelTypeVideoScreen)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	for _, obj := range mesh.ProjectorModels {
		for _, part := range obj.MeshModel {
			color, err := colors.getColor(ModelTypeProjector, part.GeometryType)
			if err != nil {
				return fmt.Errorf("Could not draw stage model: %s", err)
			}
			drawMesh(part.Mesh, canvas, color)
		}
		geometry_color, err := colors.getGeometriesColor(ModelTypeProjector)
		if err != nil {
			return fmt.Errorf("Could not draw stage model: %s", err)
		}
		for _, geometry := range obj.Geometries {
			drawMesh(geometry, canvas, geometry_color)
		}
	}

	return nil
}

func Draw(mesh *MVRTypes.StageModel, config RasterizerConfig) (*Canvas, error) {
	canvas := &Canvas{}
	err := canvas.Init(config.CanvasConfig)

	if err != nil {
		return nil, err
	}

	normalizeAndRotateStageModel(canvas, mesh, config.Rotation)

	err = drawStageModel(mesh, canvas, config)

	if err != nil {
		return nil, err
	}

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
