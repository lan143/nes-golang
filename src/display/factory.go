package display

import (
	"main/src/bus"
	"main/src/display/pixelgl"
)

type Factory struct {
	bus *bus.Bus
}

func (f *Factory) GetDisplay() Display {
	return pixelgl.NewDisplay(f.bus)
}

func NewFactory(bus *bus.Bus) *Factory {
	return &Factory{
		bus: bus,
	}
}
