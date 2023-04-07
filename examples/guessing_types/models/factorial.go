package models

type FactorialCalculator struct{}

func (c FactorialCalculator) Calculate(from uint64) uint64 {
	if from == 1 || from == 0 {
		return 1
	}
	return from * c.Calculate(from-1)
}
