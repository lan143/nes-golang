package mapper

import (
	"log"
	"main/src/mapper/enum"
	"main/src/rom"
)

type UnROMMapper struct {
	memory []byte
	rom    rom.Rom
	bank   byte

	_bank uint32
}

func (m *UnROMMapper) HasChrRom() bool {
	return false
}

func (m *UnROMMapper) GetMirroringType() enum.MirroringType {
	return m.rom.GetMirroringType()
}

func (m *UnROMMapper) LoadRom(rom rom.Rom) {
	m.rom = rom

	log.Printf("Mapper: UnROM")
	log.Printf("PRG ROM Size: %d", m.rom.GetPrgRomSize())

	// Prg ROM
	m.memory = rom.GetData()
}

func (m *UnROMMapper) GetByte(address uint16) byte {
	if address < 0xC000 {
		m._bank = uint32(m.bank)
	} else {
		m._bank = uint32(m.rom.GetPrgRomSize() - 1)
	}

	offset := uint32(address & 0x3FFF)

	return m.memory[0x4000*m._bank+offset]
}

func (m *UnROMMapper) PutByte(address uint16, value byte) {
	m.bank = value & 0xf
}

// GetChrByte Unsupported
func (m *UnROMMapper) GetChrByte(address uint16) byte {
	return 0
}

// PutChrByte Unsupported
func (m *UnROMMapper) PutChrByte(address uint16, value byte) {
}
