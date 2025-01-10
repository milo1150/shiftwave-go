package utils

import "math"

func RoundFloat64(value float64) float64 {
	fraction := value - math.Floor(value)
	if fraction >= 0.5 {
		return math.Ceil(value)
	} else {
		return math.Floor(value)
	}
}

func RoundToTwoDecimals(value float64) float64 {
	// Multiply by 100, truncate, and divide back by 100
	return math.Floor(value*100) / 100
}
