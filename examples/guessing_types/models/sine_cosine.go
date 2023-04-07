package models

import "math"

type (
	SineCalculator struct {
		factorial FactorialCalculator
	}
	CosineCalculator struct {
		factorial FactorialCalculator
	}
)

func (c SineCalculator) Calculate(angle float64) (result float64) {
	if angle > 360 {
		temp := uint32(angle)
		angle = float64(temp % 360)
	}

	radian := (angle * math.Pi) / 180

	for n := 0; n < 11; n = n + 1 {
		factorialResult := c.factorial.Calculate(uint64((2 * n) + 1))

		powResults := [...]float64{
			math.Pow(-1.0, float64(n)),
			math.Pow(radian, float64((2*n)+1)),
		}
		result = result + powResults[0]*powResults[1]/float64(factorialResult)
	}

	return result
}

func (c CosineCalculator) Calculate(angle float64) (result float64) {
	if angle > 360 {
		temp := uint32(angle)
		angle = float64(temp % 360)
	}

	radian := (angle * math.Pi) / 180

	for n := 0; n < 11; n = n + 1 {
		factorialResult := c.factorial.Calculate(uint64(2 * n))
		powResults := [...]float64{
			math.Pow(-1.0, float64(n)), math.Pow(radian, float64(2*n)),
		}
		result = result + powResults[0]*powResults[1]/float64(factorialResult)
	}

	return result
}
