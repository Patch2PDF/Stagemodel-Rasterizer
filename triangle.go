package rasterizer

import (
	"fmt"
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

func get_signed_triangle_area_delta(bx float64, by float64, cx float64, cy float64, inv_total_area float64) (dx float64, dy float64) {
	factor := inv_total_area * .5
	return (factor * (by - cy)), (factor * (cx - bx))
}

func (triangle Triangle) boundingTriangle(canvas *Canvas, color color.NRGBA) (bbminx int, bbminy int, bbmaxx int, bbmaxy int, err error) {
	if triangle.a.x == triangle.b.x && triangle.b.x == triangle.c.x {
		return 0, 0, 0, 0, fmt.Errorf("Triangle has 0 width")
	}
	if triangle.a.y == triangle.b.y && triangle.b.y == triangle.c.y {
		return 0, 0, 0, 0, fmt.Errorf("Triangle has 0 height")
	}

	total_area := triangle.signed_triangle_area()
	if total_area < 1 {
		return 0, 0, 0, 0, fmt.Errorf("Triangle area is less than 1 pixel") // backface culling + discarding triangles that cover less than a pixel // TODO: test / improve
	}

	// TODO: add canvas bounds check (via min/max with canvas bounds?)
	bbminx = int(math.Min(math.Min(triangle.a.x, triangle.b.x), triangle.c.x))
	bbminy = int(math.Min(math.Min(triangle.a.y, triangle.b.y), triangle.c.y))
	bbmaxx = int(math.Max(math.Max(triangle.a.x, triangle.b.x), triangle.c.x))
	bbmaxy = int(math.Max(math.Max(triangle.a.y, triangle.b.y), triangle.c.y))

	inv_total_area := 1 / total_area

	xf, yf := float64(bbminx), float64(bbminy) // casting

	// calculate triangle area change per x / y step
	alpha_dx, alpha_dy := get_signed_triangle_area_delta(triangle.b.x, triangle.b.y, triangle.c.x, triangle.c.y, inv_total_area)
	beta_dx, beta_dy := get_signed_triangle_area_delta(triangle.c.x, triangle.c.y, triangle.a.x, triangle.a.y, inv_total_area)
	gamma_dx, gamma_dy := get_signed_triangle_area_delta(triangle.a.x, triangle.a.y, triangle.b.x, triangle.b.y, inv_total_area)

	// base triangle area in upper left corner
	row_alpha := signed_triangle_area(xf, yf, triangle.b.x, triangle.b.y, triangle.c.x, triangle.c.y) * inv_total_area
	row_beta := signed_triangle_area(xf, yf, triangle.c.x, triangle.c.y, triangle.a.x, triangle.a.y) * inv_total_area
	row_gamma := signed_triangle_area(xf, yf, triangle.a.x, triangle.a.y, triangle.b.x, triangle.b.y) * inv_total_area

	for y := bbminy; y <= bbmaxy; y++ {
		drawing_canvas := canvas.canvas
		zBufRowIndex := y * canvas.width

		alpha := row_alpha
		beta := row_beta
		gamma := row_gamma

		var z float64
		var zBufIndex int

		for x := bbminx; x <= bbmaxx; x++ {
			if alpha < 0 || beta < 0 || gamma < 0 {
				goto inc_area_calc // negative barycentric coordinate => the pixel is outside the triangle
			}

			z = (alpha*triangle.a.z + beta*triangle.b.z + gamma*triangle.c.z)
			zBufIndex = zBufRowIndex + x
			if z <= (canvas.zbuffer)[zBufIndex] {
				goto inc_area_calc
			}
			canvas.zbuffer[zBufIndex] = z

			drawing_canvas.SetNRGBA(x, y, color)

		inc_area_calc:
			alpha += alpha_dx
			beta += beta_dx
			gamma += gamma_dx
		}
		row_alpha += alpha_dy
		row_beta += beta_dy
		row_gamma += gamma_dy
	}

	return bbminx, bbminy, bbmaxx, bbmaxy, nil
}
