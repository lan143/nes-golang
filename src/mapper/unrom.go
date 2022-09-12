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
	var bank uint32

	if address < 0xC000 {
		bank = uint32(m.bank)
	} else {
		bank = uint32(m.rom.GetPrgRomSize() - 1)
	}

	offset := uint32(address & 0x3FFF)

	if 0x4000*bank+offset > uint32(len(m.memory)) {
		log.Printf("Address: 0x%X bank: %d Offset: 0x%X", address, bank, 0x4000*bank+offset)
	}

	return m.memory[0x4000*bank+offset]
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
