package ppu

import "main/src/bus"

type Factory struct {
	b *bus.Bus
}

func (f *Factory) GetPPU() PPU {
	return NewPPU(f.b)
}

func NewFactory(b *bus.Bus) *Factory {
	return &Factory{
		b: b,
	}
}
