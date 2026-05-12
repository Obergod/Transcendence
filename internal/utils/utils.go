package utils

import(
	"golang.org/x/image/math/fixed"
)

func FixedToFloat32(f fixed.Int26_6) float32 {
	return float32(f) / 64
}

func FixedToFloat64(f fixed.Int26_6) float64 {
	return float64(f) / 64
}

func FixedDiv(a, b fixed.Int26_6) fixed.Int26_6 {
    if b == 0 {
        return 0
    }
    return fixed.Int26_6((int64(a) * 64) / int64(b))
}

func FixedSqrt(a fixed.Int26_6) fixed.Int26_6 {
    if a <= 0 {
        return 0
    }
    x := a
    for i := 0; i < 10; i++ {
        if x == 0 {
            break
        }
        x = (x + FixedDiv(a, x)) / 2
    }
    return x
}