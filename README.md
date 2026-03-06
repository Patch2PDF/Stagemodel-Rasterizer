# Patch2PDF Rasterizer

This rasterizer is used for 2D renderings of stage models.

Triangle Rasterization is inspired by [this course](https://haqr.eu/tinyrenderer/) and has been enhanced with custom performance optimisations.

In addition to only drawing the stage model, labels with the corresponding fixture id can also be added.

## Configuration

Each Canvas output can be customized via the `RasterizerConfig` struct.
See below for examples

## Examples

### Without Labels

```go
RasterizerConfig{
  RenderLabels: false,
  Rotation: rasterizer.Rotation{Alpha: 80, Beta: 0, Gamma: 200},
  OverrideColors: OverrideColorMap{},
  CanvasConfig: CanvasConfig{
		Width:         4000,
		Height:        3000,
		LabelFont:     fontBytes, // raw bytes of an embedded ttf font
		LabelDPI:      300,
		LabelFontSize: 10,
	}
}
```

![Rendering without Labels](images/side.png)

### With Labels

```go
RasterizerConfig{
  RenderLabels: true,
  Rotation: rasterizer.Rotation{Alpha: 90, Beta: 0, Gamma: 180},
  OverrideColors: OverrideColorMap{},
  CanvasConfig: CanvasConfig{
		Width:         4000,
		Height:        3000,
		LabelFont:     fontBytes, // raw bytes of an embedded ttf font
		LabelDPI:      300,
		LabelFontSize: 10,
	}
}
```

![Rendering with Labels](images/front.png)
