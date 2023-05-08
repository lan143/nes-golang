package registers

import (
	"unsafe"
)

type Register[T byte | uint16 | uint32] struct {
	value T
	width byte
}

func (r *Register[T]) Init() {
	r.width = byte(unsafe.Sizeof(r.value) * 8)
}

func (r *Register[T]) Get() T {
	return r.value
}

func (r *Register[T]) Set(value T) {
	r.value = value
}

func (r *Register[T]) SetLowerByte(value byte) {
	var andVal T
	andVal = ((1 << r.width) - 1) << 8
	upperBytes := r.value & andVal
	r.value = upperBytes | T(value)
}

func (r *Register[T]) LoadBit(pos byte) T {
	return (r.value >> pos) & 1
}

func (r *Register[T]) Shift(value T) T {
	value &= 0x1
	carry := r.LoadBit(r.width - 1)
	r.value = (r.value << 1) | value

	return carry
}
