package mmc1

import (
	"log"
	"main/src/mapper/enum"
	"main/src/rom"
)

type Mapper struct {
	data []byte
	rom  rom.Rom

	controlRegister  byte
	chrBank0Register byte
	chrBank1Register byte
	prgBankRegister  byte
	latch            byte

	registerWriteCount byte
}

func (m *Mapper) HasChrRom() bool {
	return true
}

func (m *Mapper) GetMirroringType() enum.MirroringType {
	mirroringType := m.controlRegister & 0x03

	// 0: one-screen, lower bank; 1: one-screen, upper bank
	if mirroringType == 0 || mirroringType == 1 {
		return enum.SingleScreen
	} else if mirroringType == 2 {
		return enum.Vertical
	} else {
		return enum.Horizontal
	}
}

func (m *Mapper) LoadRom(rom rom.Rom) {
	m.controlRegister = 0x0C
	m.rom = rom
	m.data = rom.GetData()

	log.Printf("Mapper: MMC1")
	log.Printf("PRG ROM Size: %d", m.rom.GetPrgRomSize())
	log.Printf("CHR ROM Size: %d", m.rom.GetChrRomSize())
}

func (m *Mapper) GetByte(address uint16) byte {
	var bank uint32 = 0
	offset := uint32(address & 0x3FFF)
	bankNum := uint32(m.prgBankRegister & 0x0F)
	prgRomBankMode := (m.controlRegister >> 2) & 0x03

	if prgRomBankMode == 0 || prgRomBankMode == 1 { // switch 32 KB at $8000, ignoring low bit of bank number
		offset |= uint32(address & 0x4000)
		bank = bankNum & 0x0E
	} else if prgRomBankMode == 2 { // fix first bank at $8000 and switch 16 KB bank at $C000
		if address < 0xC000 {
			bank = 0
		} else {
			bank = bankNum
		}
	} else if prgRomBankMode == 3 { // fix last bank at $C000 and switch 16 KB bank at $8000
		if address >= 0xC000 {
			bank = uint32(m.rom.GetPrgRomSize() - 1)
		} else {
			bank = bankNum
		}
	}

	return m.data[bank*0x4000+offset]
}

func (m *Mapper) PutByte(address uint16, value byte) {
	if value&0x80 > 0 {
		m.registerWriteCount = 0
		m.latch = 0

		if address&0x6000 == 0 {
			m.controlRegister |= 0x03 << 2
		}
	} else {
		m.latch = ((value & 0x01) << 4) | (m.latch >> 1)
		m.registerWriteCount++

		if m.registerWriteCount >= 5 {
			switch address & 0x6000 {
			case 0x0000:
				m.controlRegister = m.latch
				break
			case 0x2000:
				m.chrBank0Register = m.latch
				break
			case 0x4000:
				m.chrBank1Register = m.latch
				break
			case 0x6000:
				m.prgBankRegister = m.latch
				break
			}

			m.registerWriteCount = 0
			m.latch = 0
		}
	}
}

func (m *Mapper) GetChrByte(address uint16) byte {
	var bank uint32 = 0
	offset := uint32(address & 0x0FFF)

	// CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
	if m.controlRegister&0x10 == 0 {
		bank = uint32(m.chrBank0Register & 0x1E)
		offset |= uint32(address & 0x1000)
	} else {
		if address < 0x1000 {
			bank = uint32(m.chrBank0Register)
		} else {
			bank = uint32(m.chrBank1Register & 0x1F)
		}
	}

	return m.data[uint32(m.rom.GetPrgRomSize())*0x4000+bank*0x1000+offset]
}

func (m *Mapper) PutChrByte(address uint16, value byte) {
}
