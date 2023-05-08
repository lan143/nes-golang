package display

import (
	"main/src/bus"
	"main/src/config"
	"main/src/display/pixelgl"
)

type Factory struct {
	bus    *bus.Bus
	config *config.Config
}

func (f *Factory) GetDisplay() Display {
	return pixelgl.NewDisplay(f.config, f.bus)
	//return none.NewDisplay()
}

func NewFactory(bus *bus.Bus, config *config.Config) *Factory {
	return &Factory{
		bus:    bus,
		config: config,
	}
}
