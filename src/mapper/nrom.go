package mapper

import (
	"main/src/mapper/enum"
	"main/src/rom"
)

type NROMMapper struct {
	memory [0xffff]byte
	chrRom [0x2000]byte
	rom    rom.Rom
}

func (m *NROMMapper) HasChrRom() bool {
	return true
}

func (m *NROMMapper) GetMirroringType() enum.MirroringType {
	return m.rom.GetMirroringType()
}

func (m *NROMMapper) LoadRom(rom rom.Rom) {
	m.rom = rom

	var i, j uint16
	data := rom.GetData()
	j = 0

	// Prg ROM
	if rom.GetPrgRomSize() == 2 {
		for i = 0x8000; i <= 0xBFFF; i++ {
			m.memory[i] = data[j]
			j++
		}
	}

	// Prg ROM
	for i = 0xC000; i <= 0xFFFE; i++ {
		m.memory[i] = data[j]
		j++
	}

	j++

	// Chr ROM
	for i = 0x0000; i <= 0x1FFF; i++ {
		m.chrRom[i] = data[j]
		j++
	}
}

func (m *NROMMapper) GetByte(address uint16) byte {
	return m.memory[address]
}

func (m *NROMMapper) PutByte(address uint16, value byte) {
	m.memory[address] = value
}

func (m *NROMMapper) GetChrByte(address uint16) byte {
	return m.chrRom[address]
}

func (m *NROMMapper) PutChrByte(address uint16, value byte) {
	m.chrRom[address] = value
}
