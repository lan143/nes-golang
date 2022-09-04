package src

import (
	"context"
	"main/src/bus"
	"main/src/cpu"
	"main/src/display"
	"main/src/mapper"
	"main/src/ppu"
	"main/src/rom"
	"os"
)

type Nes struct {
	romFactory     *rom.Factory
	mapperFactory  *mapper.Factory
	cpuFactory     *cpu.Factory
	ppuFactory     *ppu.Factory
	displayFactory *display.Factory

	bus *bus.Bus

	cpu     cpu.CPU
	ppu     ppu.PPU
	display display.Display
}

func (n *Nes) Init() error {
	n.bus.Init()

	// @todo: implement load file name from cmd
	//f, err := os.Open("battle_city.nes")
	f, err := os.Open("helloworld.nes")
	if err != nil {
		panic(err)
	}

	r, err := n.romFactory.GetRom(f)
	if err != nil {
		return err
	}

	// @todo: transfer mapper id from rom
	m, err := n.mapperFactory.GetMapper(mapper.NROM)
	if err != nil {
		return err
	}

	m.LoadRom(r.GetData())

	n.display = n.displayFactory.GetDisplay()
	n.display.Init()

	n.cpu = n.cpuFactory.GetCPU()
	n.cpu.Init(m)

	n.ppu = n.ppuFactory.GetPPU()
	n.ppu.Init(m, n.display)

	return nil
}

func (n *Nes) Run(ctx context.Context) {
	// @todo: use wait group, run all in goroutines, process signals from OS...
	go n.ppu.Run()
	go n.cpu.Run()

	n.display.Run()
}

func NewNes(
	bus *bus.Bus,
	romFactory *rom.Factory,
	mapperFactory *mapper.Factory,
	cpuFactory *cpu.Factory,
	ppuFactory *ppu.Factory,
	displayFactory *display.Factory,
) *Nes {
	return &Nes{
		bus:            bus,
		romFactory:     romFactory,
		mapperFactory:  mapperFactory,
		cpuFactory:     cpuFactory,
		ppuFactory:     ppuFactory,
		displayFactory: displayFactory,
	}
}
