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
