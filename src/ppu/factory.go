package ppu

import (
	"main/src/bus"
	"main/src/ppu/RP2A03"
)

type Factory struct {
	b *bus.Bus
}

func (f *Factory) GetPPU() PPU {
	return RP2A03.NewPPU(f.b)
}

func NewFactory(b *bus.Bus) *Factory {
	return &Factory{
		b: b,
	}
}
