package rawv2_test

import "math"

func ae(f1, f2 float32) bool {
	return math.Abs(float64(f1-f2)) < 1e-4
}
