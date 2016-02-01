package fuzzy

type Defuzzificator interface {
	Defuzzify([]float64) float64
}

type weightedAverageDefuzzificator struct {
	funcs []TriangularFunc
}

func NewWeightedAverageDefuzzificator(funcs ...TriangularFunc) Defuzzificator {
	return &weightedAverageDefuzzificator{
		funcs: funcs,
	}
}

func (d *weightedAverageDefuzzificator) Defuzzify(weightVector []float64) float64 {
	topSum, downSum := 0.0, 0.0
	for i, w := range weightVector {
		topSum += w * d.funcs[i].Mid()
		downSum += w
	}
	return topSum / downSum
}
