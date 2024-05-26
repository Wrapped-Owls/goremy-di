package models

import "math"

type RootCalculator struct {
	Radical   int32
	Precision uint16
}

func (c RootCalculator) executeCalculation(
	number, increment, actualResult float64, totalCalls uint16,
) float64 {
	var (
		temp          float64 = 0
		nextIncrement         = increment / 10
	)
	if increment > number {
		temp = nextIncrement
	}
	for temp == 0 {
		if math.Pow(actualResult, float64(c.Radical)) <= number &&
			math.Pow(actualResult+increment, float64(c.Radical)) > number {
			temp = actualResult
		}
		actualResult = actualResult + increment
	}
	actualResult = temp
	if math.Pow(actualResult, float64(c.Radical)) == number || totalCalls >= c.Precision {
		return actualResult
	}

	return c.executeCalculation(number, nextIncrement, actualResult, totalCalls+1)
}

func (c RootCalculator) Calculate(from float64) (result float64) {
	shouldDivide := c.Radical < 0
	if shouldDivide {
		c.Radical = -1 * c.Radical
	}

	if c.Radical == 0 {
		return 1
	}

	if from < 0 {
		from *= -1
	}
	result = c.executeCalculation(from, 1, 1, 0)

	if shouldDivide {
		result = 1 / result
	}

	return
}
