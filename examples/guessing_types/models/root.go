package models

import "math"

type RootCalculator struct {
	Radical   int32
	Precision uint16
}

func (c RootCalculator) executeCalculation(
	firstNum, increment, actualResult float64, totalCalls uint16,
) float64 {
	var temp float64 = 0
	for temp == 0 {
		if math.Pow(actualResult, float64(c.Radical)) <= firstNum &&
			math.Pow(actualResult+increment, float64(c.Radical)) > firstNum {
			temp = actualResult
		}
		actualResult = actualResult + increment
	}
	actualResult = temp
	if math.Pow(actualResult, float64(c.Radical)) == firstNum || totalCalls >= c.Precision {
		return actualResult
	}

	return c.executeCalculation(firstNum, increment/10, actualResult, totalCalls+1)
}

func (c RootCalculator) Calculate(from float64) (result float64) {
	shouldDivide := c.Radical < 0
	if shouldDivide {
		c.Radical = -1 * c.Radical
	}

	if c.Radical == 0 {
		return 1
	}

	result = c.executeCalculation(from, 1, 1, 0)

	if shouldDivide {
		result = 1 / result
	}

	return
}
