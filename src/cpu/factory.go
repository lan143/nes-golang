package cpu

import (
	"main/src/bus"
	"main/src/cpu/Ricoh6502"
)

type Factory struct {
	b *bus.Bus
}

func (f *Factory) GetCPU() CPU {
	return Ricoh6502.NewCPU(f.b)
}

func NewFactory(b *bus.Bus) *Factory {
	return &Factory{
		b: b,
	}
}
