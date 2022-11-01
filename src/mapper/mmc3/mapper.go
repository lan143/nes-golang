package mmc3

import (
	"main/src/bus"
	"main/src/enum"
)

type Mapper struct {
	bankSelectRegister byte
	mirroringRegister  byte
	irqLatchRegister   byte
	programRegisters   [2]byte
	characterRegisters [6]byte

	irqCounter       byte
	irqCounterReload bool
	irqEnabled       bool

	prgRomSize byte

	bus *bus.Bus
}

func (m *Mapper) Init(prgRomSize byte) error {
	m.prgRomSize = prgRomSize
	m.irqEnabled = true

	m.bus.OnPPUScanline(func() {
		if m.irqCounterReload {
			m.irqCounter = m.irqLatchRegister
			m.irqCounterReload = false
		} else if m.irqEnabled {
			if m.irqCounter > 0 {
				m.irqCounter--

				if m.irqCounter == 0 {
					m.bus.Interrupt(bus.IRQ)
					m.irqCounterReload = true
				}
			}
		}
	})

	return nil
}

func (m *Mapper) GetMirroringType() enum.MirroringType {
	if m.mirroringRegister&0x01 > 0 {
		return enum.Horizontal
	} else {
		return enum.Vertical
	}
}

func (m *Mapper) MapPrgRom(address uint16) uint32 {
	var bank byte
	offset := uint32(address & 0x1FFF)

	if address >= 0x8000 && address < 0xA000 {
		if m.bankSelectRegister&0x40 > 0 {
			bank = m.prgRomSize*2 - 2
		} else {
			bank = m.programRegisters[0]
		}
	} else if address >= 0xA000 && address < 0xC000 {
		bank = m.programRegisters[1]
	} else if address >= 0xC000 && address < 0xE000 {
		if m.bankSelectRegister&0x40 > 0 {
			bank = m.programRegisters[0]
		} else {
			bank = m.prgRomSize*2 - 2
		}
	} else {
		bank = m.prgRomSize*2 - 1
	}

	return uint32(bank)*0x2000 + offset
}

func (m *Mapper) MapChrRom(address uint16) uint32 {
	var bank byte
	offset := uint32(address & 0x03FF)

	if m.bankSelectRegister&0x80 > 0 {
		if address < 0x0400 {
			bank = m.characterRegisters[2]
		} else if address >= 0x0400 && address < 0x0800 {
			bank = m.characterRegisters[3]
		} else if address >= 0x0800 && address < 0x0C00 {
			bank = m.characterRegisters[4]
		} else if address >= 0x0C00 && address < 0x1000 {
			bank = m.characterRegisters[5]
		} else if address >= 0x1000 && address < 0x1800 {
			bank = m.characterRegisters[0]
		} else {
			bank = m.characterRegisters[1]
		}
	} else {
		if address < 0x0800 {
			bank = m.characterRegisters[0]
		} else if address >= 0x0800 && address < 0x1000 {
			bank = m.characterRegisters[1]
		} else if address >= 0x1000 && address < 0x1400 {
			bank = m.characterRegisters[2]
		} else if address >= 0x1400 && address < 0x1800 {
			bank = m.characterRegisters[3]
		} else if address >= 0x1800 && address < 0x1C00 {
			bank = m.characterRegisters[4]
		} else {
			bank = m.characterRegisters[5]
		}
	}

	return uint32(bank)*0x400 + offset
}

func (m *Mapper) PutByte(address uint16, value byte) {
	if address >= 0x8000 && address < 0xA000 {
		if address&0x01 == 0 {
			m.bankSelectRegister = value
		} else {
			registerNumber := m.bankSelectRegister & ((1 << 3) - 1)

			if registerNumber == 0 || registerNumber == 1 {
				value &= 0xFE
			}

			if registerNumber < 6 {
				m.characterRegisters[registerNumber] = value
			} else {
				m.programRegisters[registerNumber-6] = value & 0x3F
			}
		}
	} else if address >= 0xA000 && address < 0xC000 && address&0x01 == 0 {
		m.mirroringRegister = value
	} else if address >= 0xC000 && address < 0xE000 {
		if address&0x01 == 0 {
			m.irqLatchRegister = value
		}

		m.irqCounterReload = true
	} else {
		if address&0x01 == 0 {
			m.irqEnabled = false
		} else {
			m.irqEnabled = true
		}
	}
}

func NewMapper(bus *bus.Bus) *Mapper {
	return &Mapper{bus: bus}
}
