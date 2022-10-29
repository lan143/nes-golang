package nrom

import (
	"main/src/enum"
)

type Mapper struct {
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
	if m.prgRomSize == 1 && address >= 0xC000 {
		address -= 0x4000
	}

	return uint32(address - 0x8000)
}

func (m *Mapper) MapChrRom(address uint16) uint32 {
	return uint32(address)
}

func (m *Mapper) PutByte(address uint16, value byte) {
}
