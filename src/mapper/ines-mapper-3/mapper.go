package ines_mapper_3

import "main/src/enum"

type Mapper struct {
	bankSelectRegister byte
}

func (m *Mapper) Init(prgRomSize byte) error {
	return nil
}

func (m *Mapper) GetMirroringType() enum.MirroringType {
	return 0
}

func (m *Mapper) MapPrgRom(address uint16) uint32 {
	return uint32(address) - 0x8000
}

func (m *Mapper) MapChrRom(address uint16) uint32 {
	return uint32(m.bankSelectRegister)*0x2000 + uint32(address&0x1FFF)
}

func (m *Mapper) PutByte(address uint16, value byte) {
	m.bankSelectRegister = value & 0xf
}
