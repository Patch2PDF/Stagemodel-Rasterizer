package rasterizer

type RasterizerConfig struct {
	CanvasConfig   CanvasConfig
	RenderLabels   bool
	Rotation       Rotation
	OverrideColors OverrideColorMap
}

type CanvasConfig struct {
	Width         int
	Height        int
	LabelFont     []byte
	LabelDPI      float64
	LabelFontSize float64
}
