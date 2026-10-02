package mmc3

import (
	"testing"

	"main/src/bus"
)

// irqHarness wires an MMC3 mapper to a bus and counts IRQ pulses, so the
// tests can drive the IRQ unit through the same bus surface production uses.
type irqHarness struct {
	bus    *bus.Bus
	mapper *Mapper
	irqs   int
}

func newIRQHarness(t *testing.T) *irqHarness {
	t.Helper()

	b := bus.NewBus()
	b.Init()

	h := &irqHarness{bus: b}
	b.OnInterrupt(bus.IRQ, func() { h.irqs++ })

	h.mapper = NewMapper(b)
	if err := h.mapper.Init(8); err != nil {
		t.Fatalf("Init: %v", err)
	}

	return h
}

func (h *irqHarness) clock(n int) {
	for i := 0; i < n; i++ {
		h.bus.DrivePPUScanline()
	}
}

func (h *irqHarness) write(address uint16, value byte) {
	h.mapper.PutByte(address, value)
}

func (h *irqHarness) assertIRQs(t *testing.T, want int) {
	t.Helper()

	if h.irqs != want {
		t.Fatalf("IRQ pulses: got %d, want %d", h.irqs, want)
	}
}

func TestMMC3IRQDisabledAtPowerOn(t *testing.T) {
	h := newIRQHarness(t)

	h.clock(100)

	h.assertIRQs(t, 0)
}

func TestMMC3IRQArmsAndReloadsPeriodically(t *testing.T) {
	h := newIRQHarness(t)

	h.write(0xC000, 3) // latch = 3 (must not arm by itself)
	h.write(0xC001, 0) // arm reload
	h.write(0xE001, 0) // enable

	// Clock 1 loads the latch (3), clocks 2-4 decrement 3->0, pulse on clock 4.
	h.clock(3)
	h.assertIRQs(t, 0)

	h.clock(1)
	h.assertIRQs(t, 1)

	// Auto-reload: next zero crossing is 4 clocks later, so 2 pulses total.
	h.clock(4)
	h.assertIRQs(t, 2)
}

func TestMMC3CounterRunsWhileDisabled(t *testing.T) {
	h := newIRQHarness(t)

	h.write(0xC000, 2)
	h.write(0xC001, 0)
	h.write(0xE001, 0)

	h.clock(3) // 2 -> 1 -> 0, pulse on clock 3
	h.assertIRQs(t, 1)

	h.write(0xE000, 0) // disable (and acknowledge)
	h.clock(10)
	h.assertIRQs(t, 1) // disabled: free-running counter never pulses

	// Re-enabling must not pulse immediately: the counter was already
	// reloaded to a non-zero value by the free-running auto-reload cycle.
	h.write(0xE001, 0)
	h.clock(1)
	h.assertIRQs(t, 1)

	// It pulses again at the next zero crossing, keeping the 3-clock period.
	h.clock(1)
	h.assertIRQs(t, 2)
}

func TestMMC3LatchWriteDoesNotArmReload(t *testing.T) {
	h := newIRQHarness(t)

	h.write(0xC000, 5)
	h.write(0xC001, 0)
	h.write(0xE001, 0)

	h.clock(3) // reload 5, then 4, then 3
	h.assertIRQs(t, 0)

	// $C000 changes the latch only. If it armed a reload the counter would
	// jump to 1 and pulse one clock early.
	h.write(0xC000, 1)

	h.clock(2) // 3 -> 2 -> 1
	h.assertIRQs(t, 0)

	h.clock(1) // 1 -> 0
	h.assertIRQs(t, 1)

	// The next zero crossing uses the new latch value 1 (not the old 5).
	h.clock(1) // reload 1
	h.assertIRQs(t, 1)

	h.clock(1) // 1 -> 0
	h.assertIRQs(t, 2)
}

func TestMMC3NonIRQWritesDoNotTouchIRQState(t *testing.T) {
	// $A001 must not enable an IRQ that is disabled at power-on.
	disabled := newIRQHarness(t)
	disabled.write(0xA001, 0xFF)
	disabled.clock(20)
	disabled.assertIRQs(t, 0)

	// $A001 and sub-$8000 writes must not disable an enabled IRQ either.
	enabled := newIRQHarness(t)
	enabled.write(0xC000, 1)
	enabled.write(0xC001, 0)
	enabled.write(0xE001, 0)
	enabled.write(0xA001, 0xFF)
	enabled.write(0x0000, 0x00)
	enabled.write(0x1234, 0x00)

	enabled.clock(2) // reload 1, then 1 -> 0
	enabled.assertIRQs(t, 1)
}

func TestMMC3LatchZeroSharpBehavior(t *testing.T) {
	h := newIRQHarness(t)

	h.write(0xC000, 0) // latch = 0
	h.write(0xC001, 0) // arm reload
	h.write(0xE001, 0) // enable (clears pending-ack)

	// First scanline pulses; the pending-ack latch then suppresses repeats
	// until the handler acknowledges through $E000/$E001.
	h.clock(1)
	h.assertIRQs(t, 1)

	h.clock(3)
	h.assertIRQs(t, 1)

	// Acknowledging each scanline re-arms the assertion, which is the Sharp
	// "one IRQ every scanline" behaviour real MMC3 boards exhibit.
	for i := 0; i < 4; i++ {
		h.write(0xE001, 0)
		h.clock(1)
	}
	h.assertIRQs(t, 5)

	// Disabling stops the pulses entirely.
	h.write(0xE000, 0)
	h.clock(5)
	h.assertIRQs(t, 5)

	// Re-enabling acknowledges and allows the next scanline to pulse again.
	h.write(0xE001, 0)
	h.clock(1)
	h.assertIRQs(t, 6)
}
