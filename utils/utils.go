package utils

import (
	"math"

	"github.com/thoas/go-funk"
)

func CalcReducedQuadraticScores(scoresTotal float64, percentages []float64) []float64 {
	return funk.Map(percentages, func(p float64) float64 {
		return p * scoresTotal
	}).([]float64)
}

func CalcPercentageOfSum(choice float64, choices []float64) float64 {
	if choice == 0.0 {
		return 0.0
	}

	whole := funk.Reduce(choices, func(acc float64, c float64) float64 {
		return acc + c
	}, 0).(float64)

	if whole == 0.0 {
		return 0.0
	}

	return choice / whole
}

func FloatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}
