package registers

type Uint16Register struct {
	value uint16
}

func (r *Uint16Register) GetValue() uint16 {
	return r.value
}

func (r *Uint16Register) Shift(value uint16) uint16 {
	value &= 0x1
	carry := value & (1 << 15)
	r.value = (r.value << 1) | value

	return carry
}

func (r *Uint16Register) StoreLowerByte(value byte) {
	r.value = uint16(value)
}
