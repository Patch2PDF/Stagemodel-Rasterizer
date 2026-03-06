package rasterizer

import "github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"

type RasterizerConfig struct {
	CanvasConfig   CanvasConfig
	RenderLabels   bool
	Rotation       MeshTypes.Matrix
	OverrideColors OverrideColorMap
}

type CanvasConfig struct {
	Width         int
	Height        int
	LabelFont     []byte
	LabelDPI      float64
	LabelFontSize float64
}
