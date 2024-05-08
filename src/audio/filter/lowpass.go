package filter

import "math"

type LowPassFilter struct {
	prevOutput float32
	alpha      float32
}

func NewLowPassFilter(fc float64, sampleRate uint32) *LowPassFilter {
	dt := 1 / float64(sampleRate)
	alpha := float32((2 * math.Pi * fc * dt) / (2*math.Pi*fc*dt + 1))

	return &LowPassFilter{
		prevOutput: 0,
		alpha:      alpha,
	}
}

func (f *LowPassFilter) Filter(input float32) float32 {
	output := f.alpha*input + (1-f.alpha)*f.prevOutput
	f.prevOutput = output

	return output
}
