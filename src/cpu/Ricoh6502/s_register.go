package Ricoh6502

import "main/src/mapper"

const (
	StackOffset uint16 = 0x100
)

type SRegister struct {
	value  byte
	mapper mapper.Mapper
}

func (r *SRegister) Init(mapper mapper.Mapper) {
	r.mapper = mapper
	r.value = 0xFD
}

func (r *SRegister) GetValue() byte {
	return r.value
}

func (r *SRegister) SetValue(value byte) {
	r.value = value
}

func (r *SRegister) PushUint16(value uint16) {
	r.PushByte(byte((value >> 8) & 0xff))
	r.PushByte(byte(value & 0xff))
}

func (r *SRegister) PushByte(value byte) {
	r.mapper.PutByte(uint16(r.value)+StackOffset, value)
	r.value--
}

func (r *SRegister) PopUint16() uint16 {
	byte1 := r.PopByte()
	byte2 := r.PopByte()

	return (uint16(byte2) << 8) | uint16(byte1)
}

func (r *SRegister) PopByte() byte {
	r.value++

	return r.mapper.GetByte(uint16(r.value) + StackOffset)
}
