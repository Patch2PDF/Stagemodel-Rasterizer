package rasterizer

import (
	"math"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
)

type StageModel = MVRTypes.StageModel

func rotateAndTranslateStageModel(stageModel *StageModel, matrix MeshTypes.Matrix) {
	for scene_object_id := range stageModel.SceneObjectModels {
		scene_object := &stageModel.SceneObjectModels[scene_object_id]
		for part_id := range scene_object.MeshModel {
			(&scene_object.MeshModel[part_id].Mesh).RotateAndTranslate(matrix)
		}
		for geometry_id := range scene_object.Geometries {
			(&scene_object.Geometries[geometry_id]).RotateAndTranslate(matrix)
		}
	}

	for fixture_id := range stageModel.FixtureModels {
		fixture := &stageModel.FixtureModels[fixture_id]
		for part_id := range fixture.MeshModel {
			(&fixture.MeshModel[part_id].Mesh).RotateAndTranslate(matrix)
		}
		for geometry_id := range fixture.Geometries {
			(&fixture.Geometries[geometry_id]).RotateAndTranslate(matrix)
		}
	}
}

func updateMeshMinAndMax(mesh MeshTypes.Mesh, min *MeshTypes.Vector, max *MeshTypes.Vector) {
	for _, triangle := range mesh.Triangles {
		*min = triangle.V0.Position.Min(*min)
		*max = triangle.V0.Position.Max(*max)

		*min = triangle.V1.Position.Min(*min)
		*max = triangle.V1.Position.Max(*max)

		*min = triangle.V2.Position.Min(*min)
		*max = triangle.V2.Position.Max(*max)
	}
}

func calculateStageModelMinAndMax(stageModel *StageModel) (min MeshTypes.Vector, max MeshTypes.Vector) {
	min = MeshTypes.Vector{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)}
	max = MeshTypes.Vector{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)}

	for _, scene_object := range stageModel.SceneObjectModels {
		for _, part := range scene_object.MeshModel {
			updateMeshMinAndMax(part.Mesh, &min, &max)
		}
		for _, geometry := range scene_object.Geometries {
			updateMeshMinAndMax(geometry, &min, &max)
		}
	}

	for _, fixture := range stageModel.FixtureModels {
		for _, part := range fixture.MeshModel {
			updateMeshMinAndMax(part.Mesh, &min, &max)
		}
		for _, geometry := range fixture.Geometries {
			updateMeshMinAndMax(geometry, &min, &max)
		}
	}

	return min, max
}

func normalizeAndRotateStageModel(canvas *Canvas, stageModel *StageModel, rotationMatrix MeshTypes.Matrix) {
	// rotate
	rotateAndTranslateStageModel(stageModel, rotationMatrix)

	// get bounding box
	min, max := calculateStageModelMinAndMax(stageModel)

	// offset calculation for centering
	sum := min.Add(max)
	offset := sum.DivScalar(-2)
	centeringMatrix := MeshTypes.IdentityMatrix()
	centeringMatrix.X03 = centeringMatrix.X03 + offset.X
	centeringMatrix.X13 = centeringMatrix.X13 + offset.Y
	centeringMatrix.X23 = centeringMatrix.X23 + offset.Z

	theoretical_min := centeringMatrix.MulPosition(min)
	theoretical_max := centeringMatrix.MulPosition(max)

	width := float64(canvas.width - 1)
	height := float64(canvas.height - 1)

	// normalizing to canvas dimensions
	diff := theoretical_max.Sub(theoretical_min)
	scale := width / diff.X
	if scale*diff.Y >= float64(canvas.height) {
		scale = height / diff.Y
	}
	centeringMatrix = centeringMatrix.MulScalar(scale)
	centeringMatrix.X03 += width / 2
	centeringMatrix.X13 += height / 2

	rotateAndTranslateStageModel(stageModel, centeringMatrix)
}
