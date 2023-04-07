package models

type FibonacciCalculator struct{}

func (c FibonacciCalculator) Calculate(from uint64) (result uint64) {
	var (
		cursors = [2]uint64{0, 1}
		values  = append(make([]uint64, 0, from), cursors[0], cursors[1])
	)

	for iterations := uint64(1); iterations < from; iterations++ {
		result = cursors[0] + cursors[1]
		cursors[0], cursors[1] = cursors[1], result
		values = append(values, result)
	}

	return
}
