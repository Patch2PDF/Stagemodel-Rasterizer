package main

import (
	"archive/zip"
	_ "embed"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	GDTFMeshReader "github.com/Patch2PDF/GDTF-Mesh-Reader/v2"
	MVRParser "github.com/Patch2PDF/MVR-Parser"
	MVRTypes "github.com/Patch2PDF/MVR-Parser/pkg/types"
	rasterizer "github.com/Patch2PDF/Rasterizer"
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

	// content, err := gdtf.ParseGDTFByFilename("test.gdtf", true, false)
	// // content, err := gdtf.ParseGDTFByFilename("test3.gdtf", true, false)
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }
	// mesh, err := content.BuildMesh("32Ch")
	// // mesh, err := content.BuildMesh("36 channel")

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

	stageModelCopy1 := stage_model.Copy()
	stageModelCopy2 := stage_model.Copy()
	// TODO: instead of copy, rotate back? (multiply new matrix with inverse of previous) --> need to add function to Mesh-Reader module

	// buf1 := &bytes.Buffer{}
	// buf2 := &bytes.Buffer{}

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

	canvas1, err := rasterizer.Draw(&stageModelCopy1, rasterizer.RasterizerConfig{RenderLabels: true, Rotation: rasterizer.Rotation{Alpha: 80, Beta: 0, Gamma: 200}, OverrideColors: overrideColors, CanvasConfig: canvasConfig})
	if err != nil {
		log.Fatal(err)
	}
	// canvas1.SaveAsPNG(buf1)
	rasterizer.SaveCanvasAsPNGFile("side.png", canvas1)
	canvas2, err := rasterizer.Draw(&stageModelCopy2, rasterizer.RasterizerConfig{RenderLabels: true, Rotation: rasterizer.Rotation{Alpha: 90, Beta: 0, Gamma: 180}, OverrideColors: overrideColors, CanvasConfig: canvasConfig})
	if err != nil {
		log.Fatal(err)
	}
	// canvas2.SaveAsPNG(buf2)
	rasterizer.SaveCanvasAsPNGFile("front.png", canvas2)

	// eg := errgroup.Group{}

	// for range 2 {
	// 	eg.Go(func() error {
	// 		rasterizer.SaveCanvasAsPNGFile()
	// 	})
	// }

	// wg.Wait()

	pprof.StopCPUProfile()

	runtime.GC() // Get up-to-date statistics
	if err := pprof.WriteHeapProfile(mem_f); err != nil {
		panic(err)
	}

	f.Close()
}

// func main() {
// 	// fmt.Printf("%s", rasterizer.DrawLabel(0, 0))
// 	// rasterizer.InitFace()

// 	// dst := image.NewNRGBA(image.Rect(0, 0, 110, 50))

// 	// fmt.Println(rasterizer.DrawLabelBox(dst, 0, 0, "Lgtm rtfm\nexcept you"))

// 	// f, err := os.Create("out.png")
// 	// if err != nil {
// 	// 	log.Fatalf("failed to create file: %v", err)
// 	// }
// 	// defer f.Close()

// 	// if err := png.Encode(f, dst); err != nil {
// 	// 	log.Fatalf("failed to encode image: %v", err)
// 	// }
// 	err := rasterizer.DrawLabel(0, 0, "Lgtm rtfm\nexcept you")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
