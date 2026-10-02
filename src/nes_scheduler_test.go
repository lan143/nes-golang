package src

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"main/src/apu"
	"main/src/audio"
	"main/src/audio/null"
	"main/src/bus"
	"main/src/cartridge"
	"main/src/cpu"
	"main/src/display"
	"main/src/mapper"
	"main/src/ppu"
	"main/src/ppu/enum"
	"main/src/ram"
	"main/src/rom"
)

type fakeCPU struct {
	events *[]string
}

func (f *fakeCPU) Init(cartridge *cartridge.Cartridge, ram *ram.Ram) {}
func (f *fakeCPU) Reset()                                            {}
func (f *fakeCPU) RunCycle()                                         { *f.events = append(*f.events, "cpu") }

type fakeAPU struct {
	events *[]string
}

func (f *fakeAPU) Init(sampleRate uint32, audio audio.Audio) {}
func (f *fakeAPU) RunCycle()                                 { *f.events = append(*f.events, "apu") }

type fakePPU struct {
	events *[]string
	ticks  int
	limit  int
	cancel context.CancelFunc
}

func (f *fakePPU) Init(cartridge *cartridge.Cartridge, display display.Display, cpuRam *ram.Ram) {}
func (f *fakePPU) RunCycle() {
	*f.events = append(*f.events, "ppu")
	f.ticks++
	if f.ticks >= f.limit {
		f.cancel()
	}
}

// expectedTickPattern encodes the scheduler contract independently of the
// production implementation: one PPU tick on counts not divisible by 4, a
// CPU+APU pair otherwise, and for non-NTSC systems an extra CPU+APU pair on
// every count divisible by 16.
func expectedTickPattern(videoSystem enum.VideoSystem, maxCounter uint32) []string {
	var events []string
	for counter := uint32(1); counter <= maxCounter; counter++ {
		if counter%4 != 0 {
			events = append(events, "ppu")
		} else {
			events = append(events, "cpu", "apu")
		}

		if videoSystem != enum.VideoSystemNTSC && counter%16 == 0 {
			events = append(events, "cpu", "apu")
		}
	}
	return events
}

// runScheduler drives runInternal synchronously with fakes until the fake PPU
// has produced ppuLimit ticks and cancels the context.
func runScheduler(videoSystem enum.VideoSystem, ppuLimit int) []string {
	var events []string

	ctx, cancel := context.WithCancel(context.Background())
	n := &Nes{
		cpu:         &fakeCPU{events: &events},
		apu:         &fakeAPU{events: &events},
		ppu:         &fakePPU{events: &events, limit: ppuLimit, cancel: cancel},
		videoSystem: videoSystem,
	}

	var wg sync.WaitGroup
	n.runInternal(ctx, &wg)
	wg.Wait()

	return events
}

func countEvents(events []string, name string) int {
	count := 0
	for _, e := range events {
		if e == name {
			count++
		}
	}
	return count
}

func TestSchedulerTickPattern(t *testing.T) {
	// 36 PPU ticks end after counter 47, which is enough to observe the
	// every-16th-count extra CPU pair twice for PAL/Dendy (counts 16 and 32).
	const ppuLimit = 36
	const maxCounter = 47

	tests := []struct {
		name        string
		videoSystem enum.VideoSystem
		cpuTicks    int
	}{
		{name: "NTSC", videoSystem: enum.VideoSystemNTSC, cpuTicks: 11},
		{name: "PAL", videoSystem: enum.VideoSystemPAL, cpuTicks: 13},
		{name: "Dendy", videoSystem: enum.VideoSystemDendy, cpuTicks: 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := runScheduler(tt.videoSystem, ppuLimit)
			expected := expectedTickPattern(tt.videoSystem, maxCounter)

			if got := countEvents(events, "ppu"); got != ppuLimit {
				t.Fatalf("ppu ticks: got %d, want %d", got, ppuLimit)
			}
			if got := countEvents(events, "cpu"); got != tt.cpuTicks {
				t.Fatalf("cpu ticks: got %d, want %d", got, tt.cpuTicks)
			}
			if got := countEvents(events, "apu"); got != tt.cpuTicks {
				t.Fatalf("apu ticks: got %d, want %d (APU must tick once per CPU)", got, tt.cpuTicks)
			}

			if len(events) != len(expected) {
				t.Fatalf("event count: got %d, want %d\ngot:  %v\nwant: %v", len(events), len(expected), events, expected)
			}
			for i := range expected {
				if events[i] != expected[i] {
					t.Fatalf("event %d: got %q, want %q\ngot:  %v\nwant: %v", i, events[i], expected[i], events, expected)
				}
			}

			// Explicit prefix: three PPU ticks then one CPU+APU group.
			wantPrefix := []string{"ppu", "ppu", "ppu", "cpu", "apu", "ppu", "ppu", "ppu", "cpu", "apu"}
			for i, want := range wantPrefix {
				if events[i] != want {
					t.Fatalf("prefix event %d: got %q, want %q", i, events[i], want)
				}
			}
		})
	}
}

func TestSchedulerTickPatternPALDoubleCPUAtSixteen(t *testing.T) {
	const ppuLimit = 36
	events := runScheduler(enum.VideoSystemPAL, ppuLimit)

	// After the first "cpu","apu" pair of count 16 there must be a second pair.
	// Counts 1..12 are three full (ppu,ppu,ppu,cpu,apu) groups (15 events);
	// counts 13..16 form the next block, so it starts at index 15.
	const groupLen = 5
	const start = 3 * groupLen
	wantGroup16 := []string{"ppu", "ppu", "ppu", "cpu", "apu", "cpu", "apu"}
	if got := events[start : start+len(wantGroup16)]; !equalStrings(got, wantGroup16) {
		t.Fatalf("PAL group 16 at event %d: got %v, want %v", start, got, wantGroup16)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type countingDisplay struct {
	frames atomic.Int64
}

func (d *countingDisplay) Init()                   {}
func (d *countingDisplay) Run(ctx context.Context) {}
func (d *countingDisplay) RenderPixel(x, y uint16, color uint32) {
	if x == 0 && y == 0 {
		d.frames.Add(1)
	}
}

func TestSchedulerRealtimeSpeed(t *testing.T) {
	if raceEnabled {
		t.Skip("race instrumentation slows emulation below real-time; pacing bound not meaningful")
	}

	f, err := os.Open("../roms/Super Mario Bros. (Japan, USA).nes")
	if err != nil {
		t.Skipf("test ROM not available: %v", err)
	}
	defer f.Close()

	b := bus.NewBus()
	b.Init()

	r, err := rom.NewFactory().GetRom(f)
	if err != nil {
		t.Fatalf("GetRom: %v", err)
	}

	cart := cartridge.NewCartridge(mapper.NewFactory(b))
	if err := cart.LoadRom(r); err != nil {
		t.Fatalf("LoadRom: %v", err)
	}

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

	const window = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go n.runInternal(ctx, &wg)

	start := time.Now()
	time.Sleep(window)
	cancel()
	wg.Wait()
	elapsed := time.Since(start)

	frames := disp.frames.Load()
	fps := float64(frames) / elapsed.Seconds()
	t.Logf("rendered %d frames in %s -> %.3f fps", frames, elapsed, fps)

	if fps < 55 || fps > 61.5 {
		t.Fatalf("measured fps %.3f outside expected range [55, 61.5]", fps)
	}
}
