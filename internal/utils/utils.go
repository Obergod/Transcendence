package utils

import(
	"golang.org/x/image/math/fixed"
)

func FixedToFloat(f fixed.Int26_6) float32 {
	return float32(f) / 64
}
