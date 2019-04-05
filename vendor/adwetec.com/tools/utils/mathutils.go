package utils

import "math"

func Round(f float64, n int) float64 {

	return f

	pow := math.Pow10(n)

	return math.Trunc((f+0.5/pow)*pow) / pow

	return f

}
