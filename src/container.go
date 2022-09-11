package src

import (
	"go.uber.org/dig"
	"main/src/bus"
	"main/src/cpu"
	"main/src/display"
	"main/src/joypad"
	"main/src/mapper"
	"main/src/ppu"
	"main/src/rom"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	processError(container.Provide(bus.NewBus))
	processError(container.Provide(NewNes))
	processError(container.Provide(rom.NewFactory))
	processError(container.Provide(mapper.NewFactory))
	processError(container.Provide(cpu.NewFactory))
	processError(container.Provide(ppu.NewFactory))
	processError(container.Provide(display.NewFactory))
	processError(container.Provide(joypad.NewJoyPad))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
