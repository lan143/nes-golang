package Ricoh6502

import (
	"fmt"
	"log"
	"main/src/bus"
	"main/src/mapper"
	"main/src/ram"
)

type InterruptHandler uint16

const (
	NMI   InterruptHandler = 0xFFFA
	Reset                  = 0xFFFC
	IRQ                    = 0xFFFE
	BRK                    = 0xFFFE
)

type Cpu struct {
	A    byte      // аккумулятор, 8 бит;
	X, Y byte      // индексные регистры, 8 бит;
	PC   uint16    // счетчик команд, 16 бит;
	S    SRegister // указатель стека, 8 бит;
	P    PRegister // регистр состояния;
	ram  *ram.Ram  // 2KB internal RAM

	mapper           mapper.Mapper
	hasInterrupt     bool
	interruptHandler InterruptHandler

	b *bus.Bus

	skipCycles uint16
}

func (c *Cpu) Init(mapper mapper.Mapper, ram *ram.Ram) {
	c.mapper = mapper
	c.ram = ram

	c.P.Init()
	c.S.Init(c.ram)

	c.Reset()

	c.b.OnInterrupt(bus.NMI, func() {
		c.hasInterrupt = true
		c.interruptHandler = NMI
	})

	c.b.OnInterrupt(bus.IRQ, func() {
		c.hasInterrupt = true
		c.interruptHandler = IRQ
	})

	c.b.OnReadFromCPU(func(address uint16) byte {
		return c.getByte(address)
	})

	c.b.OnCPUSkipCycles(func(cycles uint16) {
		c.skipCycles += cycles
	})
}

func (c *Cpu) Reset() {
	c.skipCycles = 0

	c.hasInterrupt = true
	c.interruptHandler = Reset
}

func (c *Cpu) Run() {
	if c.skipCycles > 0 {
		c.skipCycles--
		return
	}

	if c.hasInterrupt {
		c.hasInterrupt = false
		c.interrupt(c.interruptHandler)
	}

	err := c.processCommand()
	if err != nil {
		log.Println(err)
	}
}

func (c *Cpu) processCommand() error {
	command := c.getByte(c.PC)
	d := commandHandlers[command]
	if d == nil {
		c.PC++

		return fmt.Errorf("handler for command 0x%X not found", command)
	}

	operand, err := c.loadInstructionOperand(d.Mode)
	if err != nil {
		return err
	}

	err = d.Handler.Handle(c, operand, d.Mode)

	if err != nil {
		return err
	}

	c.skipCycles = d.SkipCycles

	return nil
}

func (c *Cpu) getNextByte() byte {
	c.PC++

	return c.getByte(c.PC)
}

func (c *Cpu) getNextUint16() uint16 {
	byteOne := c.getNextByte()
	byteTwo := c.getNextByte()

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) getByte(address uint16) byte {
	if address < 0x2000 {
		return c.ram.GetByte(address)
	}

	if address >= 0x2000 && address < 0x4000 {
		address &= 0x2007

		if address == 0x2002 || address == 0x2004 || address == 0x2007 {
			return c.b.ReadByCPU(address)
		}
	}

	if address == 0x4016 {
		return c.b.ReadByCPU(address)
	}

	return c.mapper.GetByte(address)
}

func (c *Cpu) getUint16(address uint16) uint16 {
	byteOne := c.getByte(address)
	byteTwo := c.getByte(address + 1)

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) getUint16FromZeroPage(address uint16) uint16 {
	byteOne := c.getByte(address & 0xff)
	byteTwo := c.getByte((address + 1) & 0xff)

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) setByte(address uint16, value byte) {
	if address < 0x2000 {
		c.ram.SetByte(address, value)
		return
	}

	if address >= 0x2000 && address < 0x4000 {
		address &= 0x2007

		switch address {
		case 0x2000:
			c.b.WriteByCPU(address, value)
			break
		case 0x2001:
			c.b.WriteByCPU(address, value)
			break
		case 0x2003:
			c.b.WriteByCPU(address, value)
			break
		case 0x2004:
			c.b.WriteByCPU(address, value)
			break
		case 0x2005:
			c.b.WriteByCPU(address, value)
			break
		case 0x2006:
			c.b.WriteByCPU(address, value)
			break
		case 0x2007:
			c.b.WriteByCPU(address, value)
			break
		}
		return
	}

	if address >= 0x4000 && address <= 0x4017 {
		c.b.WriteByCPU(address, value)

		if address == 0x4014 {
			c.skipCycles += 514
		}
		return
	}

	if address >= 0x8000 {
		c.mapper.PutByte(address, value)
	}
}

func (c *Cpu) interrupt(handler InterruptHandler) {
	if handler == IRQ && c.P.IsI() {
		return
	}

	if handler != Reset {
		if handler != BRK {
			c.P.ClearB()
		}

		stackValue := c.PC

		c.S.PushUint16(stackValue)
		c.S.PushByte(c.P.GetValue())

		c.P.SetI()
	}

	address1 := c.getByte(uint16(handler))
	address2 := c.getByte(uint16(handler) + 1)
	c.PC = uint16(address1) | (uint16(address2))<<8
}

func NewCPU(b *bus.Bus) *Cpu {
	return &Cpu{
		b: b,
	}
}
