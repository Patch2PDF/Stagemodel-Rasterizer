package rasterizer

import (
	"image/color"
	"math"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
)

type Point struct {
	x float64
	y float64
	z float64
}

type Triangle struct {
	a Point
	b Point
	c Point
}

func NewTriangleFromMeshTriangle(triangle MeshTypes.Triangle) Triangle {
	return Triangle{
		Point{x: math.Round(triangle.V0.Position.X), y: math.Round(triangle.V0.Position.Y), z: triangle.V0.Position.Z},
		Point{x: math.Round(triangle.V1.Position.X), y: math.Round(triangle.V1.Position.Y), z: triangle.V1.Position.Z},
		Point{x: math.Round(triangle.V2.Position.X), y: math.Round(triangle.V2.Position.Y), z: triangle.V2.Position.Z},
	}
}

func (triangle Triangle) signed_triangle_area() float64 {
	return signed_triangle_area(triangle.a.x, triangle.a.y, triangle.b.x, triangle.b.y, triangle.c.x, triangle.c.y)
}

func signed_triangle_area(ax float64, ay float64, bx float64, by float64, cx float64, cy float64) float64 {
	return .5 * ((by-ay)*(bx+ax) + (cy-by)*(cx+bx) + (ay-cy)*(ax+cx))
}

func (triangle Triangle) boundingTriangle(canvas *Canvas, color color.RGBA) {
	if triangle.a.x == triangle.b.x && triangle.b.x == triangle.c.x {
		return
	}
	if triangle.a.y == triangle.b.y && triangle.b.y == triangle.c.y {
		return
	}

	total_area := triangle.signed_triangle_area()
	if total_area < 1 {
		return // backface culling + discarding triangles that cover less than a pixel // TODO: test / improve
	}

	// TODO: add canvas bounds check
	bbminx := int(math.Min(math.Min(triangle.a.x, triangle.b.x), triangle.c.x))
	bbminy := int(math.Min(math.Min(triangle.a.y, triangle.b.y), triangle.c.y))
	bbmaxx := int(math.Max(math.Max(triangle.a.x, triangle.b.x), triangle.c.x))
	bbmaxy := int(math.Max(math.Max(triangle.a.y, triangle.b.y), triangle.c.y))

	inv_total_area := 1 / total_area
	factor := .5 * inv_total_area

	// precalculating constant edge for each barycentric weight (not related to (x|y))
	edge_bc := (triangle.c.y - triangle.b.y) * (triangle.c.x + triangle.b.x)
	edge_ca := (triangle.a.y - triangle.c.y) * (triangle.a.x + triangle.c.x)
	edge_ab := (triangle.b.y - triangle.a.y) * (triangle.b.x + triangle.a.x)

	for y := bbminy; y <= bbmaxy; y++ {
		drawing_canvas := canvas.canvas
		zBufRowIndex := y * canvas.width

		for x := bbminx; x <= bbmaxx; x++ {
			xf, yf := float64(x), float64(y) // casting

			alpha := factor * ((triangle.b.y-yf)*(triangle.b.x+xf) + edge_bc + (yf-triangle.c.y)*(xf+triangle.c.x))
			beta := factor * ((triangle.c.y-yf)*(triangle.c.x+xf) + edge_ca + (yf-triangle.a.y)*(xf+triangle.a.x))
			gamma := factor * ((triangle.a.y-yf)*(triangle.a.x+xf) + edge_ab + (yf-triangle.b.y)*(xf+triangle.b.x))

			if alpha < 0 || beta < 0 || gamma < 0 {
				continue // negative barycentric coordinate => the pixel is outside the triangle
			}

			z := (alpha*triangle.a.z + beta*triangle.b.z + gamma*triangle.c.z)
			zBufIndex := zBufRowIndex + x
			if z <= (canvas.zbuffer)[zBufIndex] {
				continue
			}
			canvas.zbuffer[zBufIndex] = z

			drawing_canvas.SetRGBA(x, y, color)
		}
	}
}
