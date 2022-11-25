package src

import (
	"context"
	"main/src/apu"
	"main/src/audio"
	"main/src/bus"
	"main/src/cartridge"
	"main/src/config"
	"main/src/cpu"
	"main/src/display"
	"main/src/joypad"
	"main/src/ppu"
	"main/src/ram"
	"main/src/rom"
	"os"
	"runtime/pprof"
)

type Nes struct {
	config *config.Config

	romFactory     *rom.Factory
	cpuFactory     *cpu.Factory
	ppuFactory     *ppu.Factory
	displayFactory *display.Factory
	apuFactory     *apu.Factory
	audioFactory   *audio.Factory

	bus *bus.Bus

	cpu       cpu.CPU
	ppu       ppu.PPU
	display   display.Display
	apu       apu.APU
	joypad    *joypad.JoyPad
	cartridge *cartridge.Cartridge
}

func (n *Nes) Init(romName string) error {
	err := n.config.Init()
	if err != nil {
		return err
	}

	if n.config.PprofEnabled {
		cpuProfile, err := os.Create("cpuprofile")
		if err != nil {
			return err
		}

		err = pprof.StartCPUProfile(cpuProfile)
		if err != nil {
			return err
		}
	}

	n.bus.Init()

	f, err := os.Open(romName)
	if err != nil {
		return err
	}

	r, err := n.romFactory.GetRom(f)
	if err != nil {
		return err
	}

	err = n.cartridge.LoadRom(r)
	if err != nil {
		return err
	}

	n.joypad.Init()

	n.display = n.displayFactory.GetDisplay()
	n.display.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	n.cpu = n.cpuFactory.GetCPU()
	n.cpu.Init(n.cartridge, cpuRam)

	n.ppu = n.ppuFactory.GetPPU()
	n.ppu.Init(n.cartridge, n.display, cpuRam)

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
			//oldTime := time.Now()

			n.cpu.Run()

			// 1 CPU cycle = 3 PPU cycles
			for i := 0; i < 3; i++ {
				n.ppu.Run()
			}

			n.apu.Run()

			/*for time.Since(oldTime) < 500*time.Nanosecond {
			}*/

			// 0,000000558730074
		}
	}()

	n.display.Run()

	if n.config.PprofEnabled {
		pprof.StopCPUProfile()
	}
}

func NewNes(
	config *config.Config,
	bus *bus.Bus,
	romFactory *rom.Factory,
	cpuFactory *cpu.Factory,
	ppuFactory *ppu.Factory,
	displayFactory *display.Factory,
	apuFactory *apu.Factory,
	audioFactory *audio.Factory,
	joypad *joypad.JoyPad,
	cartridge *cartridge.Cartridge,
) *Nes {
	return &Nes{
		config:         config,
		bus:            bus,
		romFactory:     romFactory,
		cpuFactory:     cpuFactory,
		ppuFactory:     ppuFactory,
		displayFactory: displayFactory,
		apuFactory:     apuFactory,
		audioFactory:   audioFactory,
		joypad:         joypad,
		cartridge:      cartridge,
	}
}
