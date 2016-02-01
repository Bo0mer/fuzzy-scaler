package fuzzy

import "math"

type TriangularFunc struct {
	F     func(float64) float64
	start float64
	end   float64
}

func NewTriangularFunc(start, end float64) TriangularFunc {
	return TriangularFunc{
		start: start,
		end:   end,
		F: func(x float64) float64 {
			mid := (start + end) / 2.0
			dist := (end - start) / 2.0
			return math.Max(dist-math.Abs(x-mid), 0.0) / dist
		},
	}
}

func (t TriangularFunc) Mid() float64 {
	return (t.start + t.end) / 2.0
}
