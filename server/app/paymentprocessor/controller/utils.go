package controller

import "math"

func round(num float64) int {
	// https://stackoverflow.com/a/29786394
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	// https://stackoverflow.com/a/29786394
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func fromStripeFormat(num int64) float64 {
	return toFixed(float64(num)/100, 2)
}
