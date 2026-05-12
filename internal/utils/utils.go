package utils

import(
	"golang.org/x/image/math/fixed"
)

func FixedToFloat(f fixed.Int26_6) float64 {
	return float64(f) / 64
}
