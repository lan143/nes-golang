package RP2A03

import (
	"testing"

	"main/src/audio/null"
	"main/src/bus"
)

// irqRig wires a bare APU onto a fresh bus and counts IRQ pulses.
type irqRig struct {
	b     *bus.Bus
	a     *APU
	count int
}

func newIRQRig(t *testing.T) *irqRig {
	t.Helper()

	b := bus.NewBus()
	b.Init()

	a := NewApu(b)
	a.Init(44100, null.NewAudio())

	r := &irqRig{b: b, a: a}
	b.OnInterrupt(bus.IRQ, func() { r.count++ })
	// The DMC fetch path calls CPUSkipCycles; no CPU is wired in this rig.
	b.OnCPUSkipCycles(func(cycles uint16) {})

	return r
}

func (r *irqRig) run(cycles int) {
	for i := 0; i < cycles; i++ {
		r.a.RunCycle()
	}
}

// TestNoDMCFlagOnStatusWrite verifies that writing $4015 never raises the DMC
// IRQ flag. Hardware only starts/stops channels on write; the flag is set by
// the DMC generator when a sample ends and is cleared by a read.
func TestNoDMCFlagOnStatusWrite(t *testing.T) {
	r := newIRQRig(t)

	// Disable frame-sequencer IRQs so only the DMC flag behavior is observed.
	r.b.WriteByCPU(0x4017, 0x40)

	r.b.WriteByCPU(0x4015, 0x10) // enable DMC
	r.b.WriteByCPU(0x4015, 0x00) // disable DMC

	r.run(10 * 7457)

	if r.count != 0 {
		t.Fatalf("spurious IRQ pulses after $4015 writes: got %d, want 0", r.count)
	}
}

// TestStatusReadAcknowledgesDMC verifies that a read of $4015 reports the DMC
// IRQ flag and then clears it, so it does not re-fire on later frame windows.
func TestStatusReadAcknowledgesDMC(t *testing.T) {
	r := newIRQRig(t)

	// Disable frame-sequencer IRQs so the only source under test is the DMC.
	r.b.WriteByCPU(0x4017, 0x40)

	r.b.WriteByCPU(0x4010, 0x8F) // IRQ enabled, non-loop, timer index 0xF (fastest)
	r.b.WriteByCPU(0x4012, 0x00)
	r.b.WriteByCPU(0x4013, 0x01) // length 17 bytes
	r.b.WriteByCPU(0x4015, 0x10) // start DMC

	r.run(8000) // sample ends and one 7457-cycle IRQ window is crossed

	if r.count < 1 {
		t.Fatalf("DMC sample did not raise an IRQ: got %d, want >= 1", r.count)
	}

	before := r.count
	status := r.b.ReadByCPU(0x4015)
	if status&0x80 == 0 {
		t.Fatalf("status read did not report DMC IRQ: got %02X, want bit 7 set", status)
	}

	r.run(4 * 7457)

	if r.count != before {
		t.Fatalf("DMC IRQ re-fired after acknowledge: got %d pulses, want %d", r.count, before)
	}
}

// TestStatusReadReturnsFlagsAndClearsFrameIRQ verifies the frame-sequencer IRQ
// flag is reported and acknowledged by a $4015 read, while the sequencer keeps
// firing on its own schedule.
func TestStatusReadReturnsFlagsAndClearsFrameIRQ(t *testing.T) {
	r := newIRQRig(t)

	r.b.WriteByCPU(0x4017, 0x00) // 4-step mode, IRQ enabled

	r.run(4*7457 + 10)
	if r.count < 1 {
		t.Fatalf("frame sequencer did not raise an IRQ: got %d, want >= 1", r.count)
	}

	status := r.b.ReadByCPU(0x4015)
	if status&0x40 == 0 {
		t.Fatalf("status read did not report frame IRQ: got %02X, want bit 6 set", status)
	}

	before := r.count
	r.run(4 * 7457)

	if delta := r.count - before; delta != 1 {
		t.Fatalf("frame sequencer pulses after acknowledge: got delta %d, want exactly 1", delta)
	}
}
