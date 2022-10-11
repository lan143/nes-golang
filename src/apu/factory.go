package apu

import (
	"main/src/apu/RP2A03"
	"main/src/bus"
)

type Factory struct {
	b *bus.Bus
}

func (f *Factory) GetAPU() APU {
	return RP2A03.NewApu(f.b)
}

func NewFactory(b *bus.Bus) *Factory {
	return &Factory{
		b: b,
	}
}
