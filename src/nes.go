package src

import (
	"context"
	"main/src/bus"
	"main/src/cpu"
	"main/src/display"
	"main/src/mapper"
	"main/src/mapper/enum"
	"main/src/ppu"
	"main/src/ram"
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

	args := os.Args[1:]
	f, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}

	r, err := n.romFactory.GetRom(f)
	if err != nil {
		return err
	}

	// @todo: transfer mapper id from rom
	m, err := n.mapperFactory.GetMapper(enum.NROM)
	if err != nil {
		return err
	}

	m.LoadRom(r)

	n.display = n.displayFactory.GetDisplay()
	n.display.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	n.cpu = n.cpuFactory.GetCPU()
	n.cpu.Init(m, cpuRam)

	n.ppu = n.ppuFactory.GetPPU()
	n.ppu.Init(m, n.display, cpuRam)

	return nil
}

func (n *Nes) Run(ctx context.Context) {
	// @todo: use wait group, run all in goroutines, process signals from OS...
	go func() {
		//timer := time.NewTimer(time.Second / 600000)
		for {
			//<-timer.C
			//timer.Reset(time.Second / 600000)
			n.cpu.Run()
			n.ppu.Run()
		}
	}()

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
