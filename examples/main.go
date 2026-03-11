package main

import (
	"archive/zip"
	_ "embed"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	GDTFMeshReader "github.com/Patch2PDF/GDTF-Mesh-Reader/v2"
	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
	MVRParser "github.com/Patch2PDF/MVR-Parser"
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
	rasterizer "github.com/Patch2PDF/Stagemodel-Rasterizer"
)

var config = MVRTypes.MVRParserConfig{
	MeshHandling:      MVRTypes.BuildStageModel,
	ReadThumbnail:     true,
	GDTFParserWorkers: 4,
	StageMeshWorkers:  4,
}

var (
	//go:embed fonts/DejaVuSansMono.ttf
	fontBytes []byte
)

func main() {
	mem_f, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer mem_f.Close()

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("%s", err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatalf("%s", err)
	}
	GDTFMeshReader.LoadPrimitives()

	mvr, err := zip.OpenReader("test5.mvr")
	if err != nil {
		log.Fatal(err)
	}
	defer mvr.Close()

	GDTFMeshReader.LoadPrimitives()

	mvrData, err := MVRParser.ParseMVRZipReader(&mvr.Reader, config)

	if err != nil {
		log.Fatalf("%s", err)
	}
	// rasterizer.Draw(mesh)
	// jobs := make(chan *rasterizer.Canvas, 2)
	// jobs <- rasterizer.Draw(mvrData.StageModel, rasterizer.Rotation{Alpha: 80, Beta: 0, Gamma: 200}, "render.png")
	// jobs <- rasterizer.Draw(mvrData.StageModel, rasterizer.Rotation{Alpha: 90, Beta: 0, Gamma: 0}, "front.png")
	// close(jobs)

	model_config := MVRTypes.ModelConfig{
		Global: MVRTypes.GlobalModelConfig{
			RenderOnlyAddressedFixture: true,
		},
		Individual: map[string]MVRTypes.ModelNodeConfig{
			"FA992217-CB18-D844-9D42-5B791B2BF05E": {
				Exclude:                    MVRTypes.GetBoolPtr(false),
				RenderOnlyAddressedFixture: nil,
			},
			"FA992217-0AE0-E31C-4C97-F45431626CD8": {
				Exclude: MVRTypes.GetBoolPtr(true),
			},
		},
	}

	stage_model := mvrData.GetStageModel(model_config)

	overrideColors := rasterizer.OverrideColorMap{
		// rasterizer.ModelTypeFixture: map[rasterizer.GeometryType]*color.NRGBA{
		// 	rasterizer.GeometryTypeAxis: {R: 255, G: 0, B: 0, A: 255},
		// 	rasterizer.GeometryTypeBeam: {R: 0, G: 255, B: 0, A: 255},
		// },
	}

	canvasConfig := rasterizer.CanvasConfig{
		Width:         4000,
		Height:        3000,
		LabelFont:     fontBytes,
		LabelDPI:      300,
		LabelFontSize: 10,
	}

	rotation1 := MeshTypes.GenerateRotationMatrix(10, 0, -20)
	canvas1, err := rasterizer.Draw(&stage_model, rasterizer.RasterizerConfig{RenderLabels: false, Rotation: rotation1, OverrideColors: overrideColors, CanvasConfig: canvasConfig})
	if err != nil {
		log.Fatal(err)
	}

	rotation2 := MeshTypes.GenerateRotationMatrix(0, 0, 0)
	rotation2 = rotation2.ReverseTransformation(rotation1)
	rasterizer.SaveCanvasAsPNGFile("side.png", canvas1)
	canvas2, err := rasterizer.Draw(&stage_model, rasterizer.RasterizerConfig{RenderLabels: true, Rotation: rotation2, OverrideColors: overrideColors, CanvasConfig: canvasConfig})
	if err != nil {
		log.Fatal(err)
	}
	rasterizer.SaveCanvasAsPNGFile("front.png", canvas2)

	pprof.StopCPUProfile()

	runtime.GC() // Get up-to-date statistics
	if err := pprof.WriteHeapProfile(mem_f); err != nil {
		panic(err)
	}

	f.Close()
}
