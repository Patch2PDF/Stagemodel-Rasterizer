module github.com/Patch2PDF/Stagemodel-Rasterizer

go 1.25.6

require (
	github.com/Patch2PDF/GDTF-Mesh-Reader/v2 v2.2.0
	github.com/Patch2PDF/GDTF-Parser v0.4.1
	github.com/Patch2PDF/MVR-Parser v0.4.1
)

require golang.org/x/text v0.35.0 // indirect

require (
	github.com/qmuntal/gltf v0.28.0 // indirect
	golang.org/x/image v0.37.0
	golang.org/x/sync v0.20.0 // indirect
)

// replace github.com/Patch2PDF/GDTF-Mesh-Reader/v2 => ../GDTF-Mesh-Reader

// replace github.com/Patch2PDF/GDTF-Parser => ../GDTF-Parser

// replace github.com/Patch2PDF/MVR-Parser => ../MVR-Parser
