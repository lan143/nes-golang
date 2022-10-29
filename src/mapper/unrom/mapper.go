package unrom

import (
	"main/src/enum"
)

type Mapper struct {
	bank       byte
	prgRomSize byte
}

func (m *Mapper) Init(prgRomSize byte) error {
	m.prgRomSize = prgRomSize

	return nil
}

func (m *Mapper) GetMirroringType() enum.MirroringType {
	return 0
}

func (m *Mapper) MapPrgRom(address uint16) uint32 {
	var bank uint32

	if address < 0xC000 {
		bank = uint32(m.bank)
	} else {
		bank = uint32(m.prgRomSize - 1)
	}

	offset := uint32(address & 0x3FFF)

	return 0x4000*bank + offset
}

func (m *Mapper) PutByte(address uint16, value byte) {
	m.bank = value & 0xf
}

func (m *Mapper) MapChrRom(address uint16) uint32 {
	return uint32(address)
}
