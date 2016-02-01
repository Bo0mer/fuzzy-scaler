package fuzzy

type Fuzzificator interface {
	Fuzzify(Metric) []float64
}

type fuzzificator struct {
	selector func(Metric) float64
	funcs    []TriangularFunc
}

func NewFuzzificator(selector func(Metric) float64, fn ...TriangularFunc) *fuzzificator {
	return &fuzzificator{
		selector: selector,
		funcs:    fn,
	}
}

func (f *fuzzificator) Fuzzify(m Metric) []float64 {
	r := make([]float64, len(f.funcs))
	for i, fn := range f.funcs {
		val := f.selector(m)
		r[i] = fn.F(val)
	}
	return r
}
