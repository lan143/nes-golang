package Ricoh6502

import (
	"fmt"
	"log"
	"main/src/bus"
	"main/src/mapper"
	"sync"
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

	mapper           mapper.Mapper
	decoder          Decoder
	hasInterrupt     bool
	interruptHandler InterruptHandler
	interruptMX      sync.RWMutex

	b *bus.Bus
}

func (c *Cpu) Init(mapper mapper.Mapper) {
	c.mapper = mapper
	c.decoder.InitCommands()

	c.P.Init()
	c.S.Init(mapper)

	c.b.Subscribe(bus.NMIInterrupt, func() {
		c.interruptMX.Lock()
		c.hasInterrupt = true
		c.interruptHandler = NMI
		c.interruptMX.Unlock()
	})
}

func (c *Cpu) Reset() {
	c.interruptMX.Lock()
	c.hasInterrupt = true
	c.interruptHandler = Reset
	c.interruptMX.Unlock()
}

func (c *Cpu) Run() {
	if c.PC == 0 {
		c.Reset()
	}

	for {
		c.interruptMX.RLock()
		if c.hasInterrupt {
			c.hasInterrupt = false
			c.interrupt(c.interruptHandler)
		}
		c.interruptMX.RUnlock()

		err := c.processCommand()
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func (c *Cpu) processCommand() error {
	position := c.PC
	command := c.mapper.GetByte(c.PC)

	found := false

	for _, d := range c.decoder.Commands {
		if d.Command == command {
			found = true
			operand, err := c.loadInstructionOperand(d.Mode)

			if err != nil {
				return err
			}

			c.logExecution(position, d.OpcodeName, d.Mode, operand)

			err = d.Handler.Handle(c, operand, d.Mode)

			if err != nil {
				return err
			}
		}
	}

	if !found {
		return fmt.Errorf("handler for command 0x%X not found", command)
	}

	return nil
}

func (c *Cpu) getNextByte() byte {
	c.PC++

	return c.mapper.GetByte(c.PC)
}

func (c *Cpu) getByte(address uint16) byte {
	return c.mapper.GetByte(address)
}

func (c *Cpu) getUin16(address uint16) uint16 {
	byteOne := c.mapper.GetByte(address)
	byteTwo := c.mapper.GetByte(address + 1)

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) getNextUint16() uint16 {
	byteOne := c.getNextByte()
	byteTwo := c.getNextByte()

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) setByte(address uint16, value byte) {
	c.mapper.PutByte(address, value)

	switch address {
	case 0x2000:
		c.b.PushEvent(bus.Write2000)
		break
	case 0x2001:
		c.b.PushEvent(bus.Write2001)
		break
	case 0x2003:
		c.b.PushEvent(bus.Write2003)
		break
	case 0x2005:
		c.b.PushEvent(bus.Write2005)
		break
	case 0x2006:
		c.b.PushEvent(bus.Write2006)
		break
	case 0x2007:
		c.b.PushEvent(bus.Write2007)
		break
	case 0x4014:
		c.b.PushEvent(bus.Write4014)
		break
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

		c.P.SetI()

		stackValue := c.PC

		c.S.PushUint16(stackValue)
		c.S.PushByte(c.P.GetValue())
	}

	address1 := c.mapper.GetByte(uint16(handler))
	address2 := c.mapper.GetByte(uint16(handler) + 1)
	c.PC = uint16(address1) | (uint16(address2))<<8

	log.Printf("Interrupt: 0x%X. Address: 0x%X", handler, c.PC)
}

func NewCPU(b *bus.Bus) *Cpu {
	return &Cpu{
		b: b,
	}
}
