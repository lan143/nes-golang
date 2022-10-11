package src

import (
	"context"
	"main/src/apu"
	"main/src/audio"
	"main/src/bus"
	"main/src/cpu"
	"main/src/display"
	"main/src/joypad"
	"main/src/mapper"
	"main/src/mapper/enum"
	"main/src/ppu"
	"main/src/ram"
	"main/src/rom"
	"os"
	"time"
)

type Nes struct {
	romFactory     *rom.Factory
	mapperFactory  *mapper.Factory
	cpuFactory     *cpu.Factory
	ppuFactory     *ppu.Factory
	displayFactory *display.Factory
	apuFactory     *apu.Factory
	audioFactory   *audio.Factory

	bus *bus.Bus

	cpu     cpu.CPU
	ppu     ppu.PPU
	display display.Display
	apu     apu.APU
	joypad  *joypad.JoyPad
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

	m, err := n.mapperFactory.GetMapper(enum.Id(r.GetMapperId()))
	if err != nil {
		return err
	}

	m.LoadRom(r)

	n.joypad.Init()

	n.display = n.displayFactory.GetDisplay()
	n.display.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	n.cpu = n.cpuFactory.GetCPU()
	n.cpu.Init(m, cpuRam)

	n.ppu = n.ppuFactory.GetPPU()
	n.ppu.Init(m, n.display, cpuRam)

	aud := n.audioFactory.GetAudio()
	err = aud.Init()
	if err != nil {
		return err
	}

	n.apu = n.apuFactory.GetAPU()
	n.apu.Init(aud.GetSampleRate(), aud)

	return nil
}

func (n *Nes) Run(ctx context.Context) {
	// @todo: use wait group, run all in goroutines, process signals from OS...

	go func() {
		for {
			oldTime := time.Now()

			n.cpu.Run()

			// 1 CPU cycle = 3 PPU cycles
			for i := 0; i < 3; i++ {
				n.ppu.Run()
			}

			n.apu.Run()

			for time.Since(oldTime) < 500*time.Nanosecond {
			}

			// 0,000000558730074
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
	apuFactory *apu.Factory,
	audioFactory *audio.Factory,
	joypad *joypad.JoyPad,
) *Nes {
	return &Nes{
		bus:            bus,
		romFactory:     romFactory,
		mapperFactory:  mapperFactory,
		cpuFactory:     cpuFactory,
		ppuFactory:     ppuFactory,
		displayFactory: displayFactory,
		apuFactory:     apuFactory,
		audioFactory:   audioFactory,
		joypad:         joypad,
	}
}
