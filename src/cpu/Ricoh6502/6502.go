package Ricoh6502

import (
	"fmt"
	"log"
	"main/src/bus"
	"main/src/mapper"
	"sync"
)

type PFlag byte

const (
	C PFlag = 0x1  // перенос
	Z       = 0x2  // ноль
	I       = 0x4  // запрет внешних прерываний — IRQ (I=0 — прерывания разрешены)
	D       = 0x8  // режим BCD для инструкций сложения и вычитания с переносом;
	B       = 0x10 // обработка прерывания (B=1 после выполнения команды BRK);
	// V 0x20 не используется, равен 1;
	V = 0x40 // переполнение;
	N = 0x80 // знак. Равен старшему биту значения, загруженного в A, X или Y в результате выполнения операции (кроме TXS).
)

type InterruptHandler uint16

const (
	NMI   InterruptHandler = 0xFFFA
	Reset                  = 0xFFFC
	IRQ                    = 0xFFFE
	BRK                    = 0xFFFE
)

type Cpu struct {
	A    byte   // аккумулятор, 8 бит;
	X, Y byte   // индексные регистры, 8 бит;
	PC   uint16 // счетчик команд, 16 бит;
	S    byte   // указатель стека, 8 бит;
	P    byte   // регистр состояния;

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

	c.b.Subscribe(bus.VBlink, func() {
		c.interruptMX.Lock()
		c.hasInterrupt = true
		c.interruptHandler = NMI
		c.interruptMX.Unlock()
	})

	c.P = 0
	c.P |= 0x20
	c.S = 0xFD
}

func (c *Cpu) Reset() {
	c.interruptMX.Lock()
	c.hasInterrupt = true
	c.interruptHandler = Reset
	c.interruptMX.Unlock()
}

func (c *Cpu) Run() {
	breakLoop := false

	if c.PC == 0 {
		c.Reset()
	}

	for {
		if breakLoop {
			break
		}

		c.interruptMX.RLock()
		if c.hasInterrupt {
			c.hasInterrupt = false
			c.interrupt(c.interruptHandler)
		}
		c.interruptMX.RUnlock()

		fmt.Printf("PC: 0x%X ", c.PC)
		command := c.mapper.GetByte(c.PC)

		found := false

		for _, d := range c.decoder.Commands {
			if d.Command == command {
				found = true
				err := d.Handler.Handle(c, d.Mode)

				if err != nil {
					breakLoop = true
					log.Printf("Error: %s", err)
				}
				break
			}
		}

		if !found {
			breakLoop = true
			log.Printf("Handler for command 0x%X not found", command)
		}
	}
}

func (c *Cpu) getNextByte() byte {
	c.PC++

	return c.mapper.GetByte(c.PC)
}

func (c *Cpu) getByte(address uint16) byte {
	return c.mapper.GetByte(address)
}

func (c *Cpu) getNextUint16() uint16 {
	byteOne := c.getNextByte()
	byteTwo := c.getNextByte()

	return uint16(byteOne) | (uint16(byteTwo))<<8
}

func (c *Cpu) setByte(address uint16, value byte) {
	c.mapper.PutByte(address, value)
}

func (c *Cpu) setFlagsByValue(value byte) {
	if value&0x80 > 0 {
		c.P |= N
	} else {
		c.P &= ^byte(N)
	}

	if value == 0 {
		c.P |= Z
	} else {
		c.P &= ^byte(Z)
	}
}

func (c *Cpu) setCorrectionBit(value byte) {
	if value > 0 {
		c.P |= byte(C)
	} else {
		c.P &= ^byte(C)
	}
}

func (c *Cpu) interrupt(handler InterruptHandler) {
	if handler == IRQ && c.P&I > 0 {
		return
	}

	if handler != Reset {
		if handler != BRK {
			c.P &= ^byte(B)
		}

		c.P |= I

		stackValue := c.PC

		c.setByte(uint16(c.S)+0x100, byte((stackValue>>8)&0xff))
		c.S--
		c.setByte(uint16(c.S)+0x100, byte(stackValue&0xff))
		c.S--
		c.setByte(uint16(c.S)+0x100, c.P)
		c.S--
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
