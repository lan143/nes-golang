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
	irqPendingAck    bool

	prgRomSize byte

	bus *bus.Bus
}

func (m *Mapper) Init(prgRomSize byte) error {
	m.prgRomSize = prgRomSize
	// MMC3 IRQs are disabled at power-on.
	m.irqEnabled = false

	m.bus.OnPPUScanline(func() {
		// The counter runs regardless of the enable bit. A12 rising edges are
		// approximated by the existing per-scanline hook.
		if m.irqCounterReload || m.irqCounter == 0 {
			m.irqCounter = m.irqLatchRegister
			m.irqCounterReload = false
		} else {
			m.irqCounter--
		}

		// A non-zero counter clears the pending-ack latch so the next
		// zero-crossing asserts the IRQ again.
		if m.irqCounter != 0 {
			m.irqPendingAck = false
		}

		if m.irqCounter == 0 && m.irqEnabled && !m.irqPendingAck {
			m.bus.Interrupt(bus.IRQ)
			m.irqPendingAck = true
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
		} else if address >= 0x1000 && address < 0x1400 {
			bank = m.characterRegisters[0] & 0xFE
		} else if address >= 0x1400 && address < 0x1800 {
			bank = m.characterRegisters[0] | 0x01
		} else if address >= 0x1800 && address < 0x1C00 {
			bank = m.characterRegisters[1] & 0xFE
		} else {
			bank = m.characterRegisters[1] | 0x01
		}
	} else {
		if address < 0x0400 {
			bank = m.characterRegisters[0] & 0xFE
		} else if address >= 0x0400 && address < 0x0800 {
			bank = m.characterRegisters[0] | 0x01
		} else if address >= 0x0800 && address < 0x0C00 {
			bank = m.characterRegisters[1] & 0xFE
		} else if address >= 0x0C00 && address < 0x1000 {
			bank = m.characterRegisters[1] | 0x01
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
			registerNumber := m.bankSelectRegister & 0x7

			if registerNumber == 0 || registerNumber == 1 {
				value &= 0xFE
			}

			if registerNumber < 6 {
				m.characterRegisters[registerNumber] = value
			} else {
				m.programRegisters[registerNumber-6] = value & 0x3F
			}
		}
	} else if address >= 0xA000 && address < 0xC000 {
		if address&0x01 == 0 {
			m.mirroringRegister = value
		}
		// $A001 (PRG RAM protect) is not modeled and must not touch IRQ state.
	} else if address >= 0xC000 && address < 0xE000 {
		if address&0x01 == 0 {
			// $C000-$DFFE even: IRQ latch write only, never arms a reload.
			m.irqLatchRegister = value
		} else {
			// $C001-$DFFF odd: arm the reload.
			m.irqCounterReload = true
		}
	} else if address >= 0xE000 {
		if address&0x01 == 0 {
			m.irqEnabled = false
			m.irqPendingAck = false
		} else {
			m.irqEnabled = true
			m.irqPendingAck = false
		}
	}
	// Addresses below $8000 are ignored and must not affect IRQ state.
}

func NewMapper(bus *bus.Bus) *Mapper {
	return &Mapper{bus: bus}
}
