package rasterizer

import (
	"math"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
)

func swap[variable any](a variable, b variable) (variable, variable) {
	return b, a
}

func generateRotationMatrix(alpha float64, beta float64, gamma float64) MeshTypes.Matrix {
	alphaSin := math.Sin(alpha / 180 * math.Pi)
	alphaCos := math.Cos(alpha / 180 * math.Pi)
	betaSin := math.Sin(beta / 180 * math.Pi)
	betaCos := math.Cos(beta / 180 * math.Pi)
	gammaSin := math.Sin(gamma / 180 * math.Pi)
	gammaCos := math.Cos(gamma / 180 * math.Pi)

	return MeshTypes.Matrix{
		X00: betaCos * gammaCos, X01: -betaCos * gammaSin, X02: betaSin, X03: 0,
		X10: alphaCos*gammaSin + alphaSin*betaSin*gammaCos, X11: alphaCos*gammaCos - alphaSin*betaSin*gammaSin, X12: -alphaSin * betaCos, X13: 0,
		X20: alphaSin*gammaSin - alphaCos*betaSin*gammaCos, X21: alphaSin*gammaCos + alphaCos*betaSin*gammaSin, X22: alphaCos * betaCos, X23: 0,
		X30: 0, X31: 0, X32: 0, X33: 1,
	}
}
