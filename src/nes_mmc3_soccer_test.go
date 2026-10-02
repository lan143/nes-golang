package src

import (
	"os"
	"testing"

	"main/src/apu"
	"main/src/audio/null"
	"main/src/bus"
	"main/src/cartridge"
	"main/src/cpu"
	"main/src/joypad"
	"main/src/mapper"
	"main/src/ppu"
	"main/src/ppu/enum"
	"main/src/ram"
	"main/src/rom"
)

const soccerROMPath = "../roms/Kunio Kun no Nekketsu Soccer League (J) [T-Eng1.1_PentarouZero].nes"

// soccerTestRig holds a fully wired NTSC NES running the Kunio soccer ROM.
type soccerTestRig struct {
	b *bus.Bus
	n *Nes
	f *os.File
}

// newSoccerTestRig builds the emulator wiring shared by the regression phases.
func newSoccerTestRig(t *testing.T) *soccerTestRig {
	t.Helper()

	f, err := os.Open(soccerROMPath)
	if err != nil {
		t.Skipf("test ROM not available: %v", err)
	}

	b := bus.NewBus()
	b.Init()

	r, err := rom.NewFactory().GetRom(f)
	if err != nil {
		f.Close()
		t.Fatalf("GetRom: %v", err)
	}

	cart := cartridge.NewCartridge(mapper.NewFactory(b))
	if err := cart.LoadRom(r); err != nil {
		f.Close()
		t.Fatalf("LoadRom: %v", err)
	}

	jp := joypad.NewJoyPad(b)
	jp.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	disp := &countingDisplay{}

	n := &Nes{videoSystem: enum.VideoSystemNTSC}
	n.cpu = cpu.NewFactory(b).GetCPU()
	n.cpu.Init(cart, cpuRam)
	n.ppu = ppu.NewFactory(b).GetPPU()
	n.ppu.Init(cart, disp, cpuRam)
	n.apu = apu.NewFactory(b).GetAPU()
	n.apu.Init(44100, null.NewAudio())

	return &soccerTestRig{b: b, n: n, f: f}
}

// runFrames advances the emulator for the given number of PPU frames using the
// same 3-PPU/1-CPU tick pattern as the scheduler. onFrame runs before each
// frame's ticks (nil for no callback).
func (r *soccerTestRig) runFrames(frames int, onFrame func(frame int)) {
	const ticksPerFrame = 89342

	tick := 0
	for frame := 0; frame < frames; frame++ {
		if onFrame != nil {
			onFrame(frame)
		}

		for i := 0; i < ticksPerFrame; i++ {
			tick++
			if tick%4 != 0 {
				r.n.ppu.RunCycle()
			} else {
				r.n.cpu.RunCycle()
				r.n.apu.RunCycle()
			}
		}
	}
}

// TestMMC3SoccerNoPanic is a regression test for Kunio Kun no Nekketsu Soccer
// League (J), an MMC3 game whose palette RAM can hold values >= 0x40. Before
// the palette fix the renderer indexed the 64-entry palette directly and
// panicked at around frame 497.
//
// It has two deterministic phases. The first exercises the MMC3 IRQ and CPU
// interrupt paths under scripted input (must not panic). The second runs with
// no input and asserts the APU does not flood the CPU with spurious DMC IRQ
// pulses: with the $4015 write-sets-flag bug the MMC3 band counter in zero-page
// $55 grows past 0x20 within a few hundred frames.
func TestMMC3SoccerNoPanic(t *testing.T) {
	type keyEvent struct {
		frame   int
		button  bus.JoyPadButton
		pressed bool
	}
	script := []keyEvent{
		{90, bus.JoyPadButtonStart, true}, {100, bus.JoyPadButtonStart, false},
		{150, bus.JoyPadButtonStart, true}, {160, bus.JoyPadButtonStart, false},
		{220, bus.JoyPadButtonA, true}, {230, bus.JoyPadButtonA, false},
		{300, bus.JoyPadButtonStart, true}, {310, bus.JoyPadButtonStart, false},
		{380, bus.JoyPadButtonDown, true}, {390, bus.JoyPadButtonDown, false},
		{420, bus.JoyPadButtonA, true}, {430, bus.JoyPadButtonA, false},
		{500, bus.JoyPadButtonStart, true}, {510, bus.JoyPadButtonStart, false},
		{560, bus.JoyPadButtonA, true}, {570, bus.JoyPadButtonA, false},
		{620, bus.JoyPadButtonB, true}, {630, bus.JoyPadButtonB, false},
		{680, bus.JoyPadButtonA, true}, {690, bus.JoyPadButtonA, false},
	}

	// A panic anywhere in the synchronous run must fail the test rather than
	// crash the test binary.
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("emulator panic during run: %v", rec)
		}
	}()

	t.Run("scripted input", func(t *testing.T) {
		rig := newSoccerTestRig(t)
		defer rig.f.Close()

		event := 0
		rig.runFrames(1000, func(frame int) {
			for event < len(script) && script[event].frame == frame {
				rig.b.KeyEvent(script[event].button, script[event].pressed)
				event++
			}
		})
	})

	t.Run("no input stays free of DMC IRQ storm", func(t *testing.T) {
		rig := newSoccerTestRig(t)
		defer rig.f.Close()

		rig.runFrames(600, nil)

		// MMC3 band counter: with fixed APU semantics it stays 0; the
		// write-sets-DMC-flag bug drives it past 0x40.
		if got := rig.b.ReadFromCPU(0x55); got >= 0x20 {
			t.Fatalf("spurious DMC IRQ storm: band counter $55 = %02X after 600 no-input frames (want < 0x20)", got)
		}
	})
}
