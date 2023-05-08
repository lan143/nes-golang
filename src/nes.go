package src

import (
	"context"
	"fmt"
	"main/src/apu"
	"main/src/audio"
	"main/src/bus"
	"main/src/cartridge"
	"main/src/config"
	"main/src/cpu"
	"main/src/display"
	"main/src/joypad"
	"main/src/ppu"
	"main/src/ppu/enum"
	"main/src/ram"
	"main/src/rom"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
)

type Nes struct {
	sigs chan os.Signal
	wg   sync.WaitGroup

	config *config.Config

	romFactory     *rom.Factory
	cpuFactory     *cpu.Factory
	ppuFactory     *ppu.Factory
	displayFactory *display.Factory
	apuFactory     *apu.Factory
	audioFactory   *audio.Factory

	bus *bus.Bus

	cpu         cpu.CPU
	ppu         ppu.PPU
	display     display.Display
	apu         apu.APU
	joypad      *joypad.JoyPad
	cartridge   *cartridge.Cartridge
	videoSystem enum.VideoSystem
}

func (n *Nes) Init(romName string) error {
	n.sigs = make(chan os.Signal, 1)
	signal.Notify(n.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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

	switch n.config.VideoSystem {
	case "pal":
		fallthrough
	case "PAL":
		fallthrough
	case "Pal":
		n.videoSystem = enum.VideoSystemPAL
	case "ntsc":
		fallthrough
	case "NTSC":
		fallthrough
	case "Ntsc":
		n.videoSystem = enum.VideoSystemNTSC
	case "dendy":
		fallthrough
	case "DENDY":
		fallthrough
	case "Dendy":
		n.videoSystem = enum.VideoSystemDendy
	default:
		return fmt.Errorf("unsupported video system %s", n.videoSystem)
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
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go n.processSignals(cancelFunc)

	if n.config.PprofEnabled {
		go n.runProfiler(cancelCtx, &n.wg)
	}

	go n.runInternal(cancelCtx, &n.wg)

	n.display.Run(cancelCtx)
	n.wg.Wait()
}

func (n *Nes) runInternal(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var counter uint32

	cpuRunCycleChan := make(chan bool, 1)
	ppuRunCycleChan := make(chan bool, 1)

	go n.runCPU(cpuRunCycleChan)
	go n.runPPU(ppuRunCycleChan)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			counter++
		}

		if n.videoSystem == enum.VideoSystemNTSC {
			if counter%4 != 0 {
				ppuRunCycleChan <- true
			} else {
				cpuRunCycleChan <- true
			}
		} else {
			if counter%4 != 0 {
				ppuRunCycleChan <- true
			} else {
				cpuRunCycleChan <- true
			}

			if counter%16 == 0 {
				cpuRunCycleChan <- true
			}
		}
	}
}

func (n *Nes) runCPU(cpuRunCycleChan chan bool) {
	for {
		<-cpuRunCycleChan
		n.cpu.RunCycle()
		n.apu.RunCycle()
	}
}

func (n *Nes) runPPU(ppuRunCycleChan chan bool) {
	for {
		<-ppuRunCycleChan
		n.ppu.RunCycle()
	}
}

func (n *Nes) runProfiler(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	<-ctx.Done()
	pprof.StopCPUProfile()
	wg.Done()
}

func (n *Nes) processSignals(cancelFunc context.CancelFunc) {
	select {
	case <-n.sigs:
		cancelFunc()
		break
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
