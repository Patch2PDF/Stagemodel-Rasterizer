package rasterizer

type ModelType int

const (
	ModelTypeSceneObject ModelType = iota
	ModelTypeFocusPoint
	ModelTypeFixture
	ModelTypeSupport
	ModelTypeTruss
	ModelTypeVideoScreen
	ModelTypeProjector
)

var modelTypes []ModelType = []ModelType{ // needs to be updated when ModelType enum is extended
	ModelTypeSceneObject,
	ModelTypeFocusPoint,
	ModelTypeFixture,
	ModelTypeSupport,
	ModelTypeTruss,
	ModelTypeVideoScreen,
	ModelTypeProjector,
}
