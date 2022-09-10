package registers

import (
	"math"
	"unsafe"
)

type Register[T byte | uint16] struct {
	value T
}

func (r *Register[T]) Get() T {
	return r.value
}

func (r *Register[T]) Set(value T) {
	r.value = value
}

func (r *Register[T]) SetLowerByte(value byte) {
	var andVal T
	andVal = T(math.Pow(2, float64(r.GetWidth())) - 1)
	andVal <<= 8
	upperBytes := r.value & andVal
	r.value = upperBytes | T(value)
}

func (r *Register[T]) GetWidth() byte {
	return byte(unsafe.Sizeof(r.value) * 8)
}

func (r *Register[T]) LoadBit(pos byte) T {
	return (r.value >> pos) & 1
}

func (r *Register[T]) Shift(value T) T {
	value &= 0x1
	carry := r.LoadBit(r.GetWidth() - 1)
	r.value = (r.value << 1) | value

	return carry
}
