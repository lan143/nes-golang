package ppu

import (
	"main/src/mapper"
	"time"
)

type PPU struct {
	mapper mapper.Mapper
}

func (p *PPU) Init(mapper mapper.Mapper) {
	p.mapper = mapper
}

func (p *PPU) Run() {
	timer := time.NewTicker(time.Second / 25)

	for {
		select {
		case <-timer.C:
			p.mapper.PutByte(0x2002, 0x80)
			break
		}
	}
}
