package ppu

import (
	"main/src/bus"
	"main/src/mapper"
	"time"
)

type P2C02 struct {
	mapper mapper.Mapper
	bus    *bus.Bus
}

func (p *P2C02) Init(mapper mapper.Mapper) {
	p.mapper = mapper
}

func (p *P2C02) Run() {
	timer := time.NewTicker(time.Second / 25)

	for {
		select {
		case <-timer.C:
			p.mapper.PutByte(0x2002, 0x80)
			p.bus.PushEvent(bus.VBlink)
			break
		}
	}
}

func NewPPU(b *bus.Bus) PPU {
	return &P2C02{
		bus: b,
	}
}
