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

// newExhaustiveTestCPU builds a Cpu using the same rig pattern as the existing
// SBC tests so that exhaustive flag oracles exercise the real handler path.
func newExhaustiveTestCPU() *Cpu {
	bus := bus2.NewBus()
	bus.Init()

	cpuRam := &ram.Ram{}
	cpuRam.Init(0x0800)

	cpu := NewCPU(bus)
	cpu.Init(nil, cpuRam)

	return cpu
}

func TestSBCHandlerExhaustiveFlags(t *testing.T) {
	handler := SBCHandler{}

	for ai := 0; ai < 256; ai++ {
		A := byte(ai)
		for mi := 0; mi < 256; mi++ {
			M := byte(mi)
			for _, cin := range []byte{0, 1} {
				cpu := newExhaustiveTestCPU()
				cpu.A = A
				if cin == 1 {
					cpu.P.SetC()
				} else {
					cpu.P.ClearC()
				}

				if err := handler.Handle(cpu, uint16(M), enums.ModeIMM); err != nil {
					t.Fatalf("SBC A=0x%02X M=0x%02X Cin=%d: unexpected error: %v", A, M, cin, err)
				}

				borrow := uint16(1 - cin)
				res := uint16(A) - uint16(M) - borrow
				byteRes := byte(res)

				wantA := byteRes
				wantC := uint16(A) >= uint16(M)+borrow
				wantZ := byteRes == 0
				wantN := byteRes&0x80 != 0
				wantV := ((A^byteRes)&(A^M))&0x80 != 0

				if cpu.A != wantA {
					t.Errorf("SBC A=0x%02X M=0x%02X Cin=%d: A got 0x%02X want 0x%02X", A, M, cin, cpu.A, wantA)
				}
				if cpu.P.IsC() != wantC {
					t.Errorf("SBC A=0x%02X M=0x%02X Cin=%d: C got %v want %v", A, M, cin, cpu.P.IsC(), wantC)
				}
				if cpu.P.IsZ() != wantZ {
					t.Errorf("SBC A=0x%02X M=0x%02X Cin=%d: Z got %v want %v", A, M, cin, cpu.P.IsZ(), wantZ)
				}
				if cpu.P.IsN() != wantN {
					t.Errorf("SBC A=0x%02X M=0x%02X Cin=%d: N got %v want %v", A, M, cin, cpu.P.IsN(), wantN)
				}
				if cpu.P.IsV() != wantV {
					t.Errorf("SBC A=0x%02X M=0x%02X Cin=%d: V got %v want %v", A, M, cin, cpu.P.IsV(), wantV)
				}
			}
		}
	}
}

func TestADCHandlerExhaustiveFlags(t *testing.T) {
	handler := ADCHandler{}

	for ai := 0; ai < 256; ai++ {
		A := byte(ai)
		for mi := 0; mi < 256; mi++ {
			M := byte(mi)
			for _, cin := range []byte{0, 1} {
				cpu := newExhaustiveTestCPU()
				cpu.A = A
				if cin == 1 {
					cpu.P.SetC()
				} else {
					cpu.P.ClearC()
				}

				if err := handler.Handle(cpu, uint16(M), enums.ModeIMM); err != nil {
					t.Fatalf("ADC A=0x%02X M=0x%02X Cin=%d: unexpected error: %v", A, M, cin, err)
				}

				res := uint16(A) + uint16(M) + uint16(cin)
				byteRes := byte(res)

				wantA := byteRes
				wantC := res > 0xFF
				wantZ := byteRes == 0
				wantN := byteRes&0x80 != 0
				wantV := ((^(A ^ M))&(A^byteRes))&0x80 != 0

				if cpu.A != wantA {
					t.Errorf("ADC A=0x%02X M=0x%02X Cin=%d: A got 0x%02X want 0x%02X", A, M, cin, cpu.A, wantA)
				}
				if cpu.P.IsC() != wantC {
					t.Errorf("ADC A=0x%02X M=0x%02X Cin=%d: C got %v want %v", A, M, cin, cpu.P.IsC(), wantC)
				}
				if cpu.P.IsZ() != wantZ {
					t.Errorf("ADC A=0x%02X M=0x%02X Cin=%d: Z got %v want %v", A, M, cin, cpu.P.IsZ(), wantZ)
				}
				if cpu.P.IsN() != wantN {
					t.Errorf("ADC A=0x%02X M=0x%02X Cin=%d: N got %v want %v", A, M, cin, cpu.P.IsN(), wantN)
				}
				if cpu.P.IsV() != wantV {
					t.Errorf("ADC A=0x%02X M=0x%02X Cin=%d: V got %v want %v", A, M, cin, cpu.P.IsV(), wantV)
				}
			}
		}
	}
}
