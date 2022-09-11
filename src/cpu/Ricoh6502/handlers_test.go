package Ricoh6502

import (
	bus2 "main/src/bus"
	"main/src/cpu/Ricoh6502/enums"
	"main/src/ram"
	"testing"
)

func TestROLHandlerWithC(t *testing.T) {
	handler := ROLHandler{}
	cpu := &Cpu{}
	cpu.A = 0x80
	cpu.P.SetC()

	err := handler.Handle(cpu, 0, enums.ModeAcc)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0x01 {
		t.Errorf("A should be 0x01. Actual: 0x%X", cpu.A)
	}

	if !cpu.P.IsC() {
		t.Error("C should be set")
	}

	if cpu.P.IsZ() {
		t.Error("C bit should be equals 0")
	}

	if cpu.P.IsN() {
		t.Errorf("N bit should be equals 0")
	}
}

func TestROLHandlerWithoutC(t *testing.T) {
	handler := ROLHandler{}
	cpu := &Cpu{}
	cpu.A = 0x80
	cpu.P.ClearC()

	err := handler.Handle(cpu, 0, enums.ModeAcc)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0x00 {
		t.Errorf("A should be 0x00. Actual: 0x%X", cpu.A)
	}

	if !cpu.P.IsC() {
		t.Error("C should be set")
	}

	if !cpu.P.IsZ() {
		t.Error("C bit should be equals 1")
	}

	if cpu.P.IsN() {
		t.Errorf("N bit should be equals 0")
	}
}

func TestROLHandlerWithoutOverflow(t *testing.T) {
	handler := ROLHandler{}
	cpu := &Cpu{}
	cpu.A = 0x40
	cpu.P.SetC()

	err := handler.Handle(cpu, 0, enums.ModeAcc)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0x81 {
		t.Errorf("A should be 0x81. Actual: 0x%X", cpu.A)
	}

	if cpu.P.IsC() {
		t.Error("C bit should be equals 0")
	}

	if !cpu.P.IsN() {
		t.Errorf("N bit should be equals 1")
	}

	if cpu.P.IsZ() {
		t.Errorf("Z bit should be equals 0")
	}
}

func TestSBCHandler(t *testing.T) {
	handler := SBCHandler{}

	bus := bus2.NewBus()
	bus.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	cpu := NewCPU(bus)
	cpu.Init(nil, cpuRam)
	cpu.A = 0x40
	cpu.P.SetC()

	err := handler.Handle(cpu, 0x40, enums.ModeIMM)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0x00 {
		t.Errorf("A should be 0x00. Actual: 0x%X", cpu.A)
	}

	if !cpu.P.IsC() {
		t.Error("C bit should be equals 1")
	}

	if cpu.P.IsN() {
		t.Errorf("N bit should be equals 0")
	}

	if !cpu.P.IsZ() {
		t.Errorf("Z bit should be equals 1")
	}

	if cpu.P.IsV() {
		t.Errorf("V bit should be equals 0")
	}
}

func TestSBCHandlerTwo(t *testing.T) {
	handler := SBCHandler{}

	bus := bus2.NewBus()
	bus.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	cpu := NewCPU(bus)
	cpu.Init(nil, cpuRam)
	cpu.A = 0x40
	cpu.P.SetC()

	err := handler.Handle(cpu, 0x3F, enums.ModeIMM)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0x01 {
		t.Errorf("A should be 0x01. Actual: 0x%X", cpu.A)
	}

	if !cpu.P.IsC() {
		t.Error("C bit should be equals 1")
	}

	if cpu.P.IsN() {
		t.Errorf("N bit should be equals 0")
	}

	if cpu.P.IsZ() {
		t.Errorf("Z bit should be equals 0")
	}

	if cpu.P.IsV() {
		t.Errorf("V bit should be equals 0")
	}
}

func TestSBCHandlerThree(t *testing.T) {
	handler := SBCHandler{}

	bus := bus2.NewBus()
	bus.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	cpu := NewCPU(bus)
	cpu.Init(nil, cpuRam)
	cpu.A = 0x40
	cpu.P.SetC()

	err := handler.Handle(cpu, 0x41, enums.ModeIMM)
	if err != nil {
		t.Error(err)
		return
	}

	if cpu.A != 0xFF {
		t.Errorf("A should be 0xFF. Actual: 0x%X", cpu.A)
	}

	if cpu.P.IsC() {
		t.Error("C bit should be equals 0")
	}

	if !cpu.P.IsN() {
		t.Errorf("N bit should be equals 1")
	}

	if cpu.P.IsZ() {
		t.Errorf("Z bit should be equals 0")
	}

	if cpu.P.IsV() {
		t.Errorf("V bit should be equals 0")
	}
}
