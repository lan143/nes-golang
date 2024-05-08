package filter

import (
	"math"
)

type HighPassFilter struct {
	prevOutput float32
	alpha      float64
}

func NewHighPassFilter(fc float64, sampleRate uint32) *HighPassFilter {
	dt := 1 / float64(sampleRate)
	alpha := 2 * math.Pi * fc * dt

	return &HighPassFilter{
		prevOutput: 0,
		alpha:      alpha,
	}
}

func (f *HighPassFilter) Filter(input float32) float32 {
	output := float32(f.alpha) * (f.prevOutput + input)
	f.prevOutput = output

	return input - output
}
