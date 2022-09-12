package Ricoh6502

import (
	"main/src/cpu/Ricoh6502/enums"
)

type ANDHandler struct {
}

func (h *ANDHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.A &= value

	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)
	cpu.PC++

	return nil
}

type ASLHandler struct {
}

func (h *ASLHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	twoByteValue := uint16(value)
	twoByteValue <<= 1

	err = cpu.writeWithMemoryAccessType(mode, operand, byte(twoByteValue))
	if err != nil {
		return err
	}

	cpu.PC++
	cpu.P.UpdateN(byte(twoByteValue))
	cpu.P.UpdateZ(byte(twoByteValue))
	cpu.P.UpdateC(twoByteValue)

	return nil
}

type BPLHandler struct {
}

func (h *BPLHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if !cpu.P.IsN() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type CLDHandler struct {
}

func (h *CLDHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.ClearD()
	cpu.PC++

	return nil
}

type DEXHandler struct {
}

func (h *DEXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.X--
	cpu.P.UpdateN(cpu.X)
	cpu.P.UpdateZ(cpu.X)

	cpu.PC++

	return nil
}

type EORHandler struct {
}

func (h *EORHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.A ^= src
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)

	cpu.PC++

	return nil
}

type ISBHandler struct {
}

func (h *ISBHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	value++

	err = cpu.writeWithMemoryAccessType(mode, operand, value)
	if err != nil {
		return err
	}

	cpu.A -= value
	cpu.PC++
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)

	return nil
}

type JMPHandler struct {
}

func (h *JMPHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	if mode == enums.ModeIND {
		addr1 := operand
		addr2 := (operand & 0xff00) | ((operand + 1) & 0xff)

		cpu.PC = uint16(cpu.getByte(addr1)) | (uint16(cpu.getByte(addr2)) << 8)
	} else {
		cpu.PC = operand
	}

	return nil
}

type JSRHandler struct {
}

func (h *JSRHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.S.PushUint16(cpu.PC)
	cpu.PC = operand

	return nil
}

type LDAHandler struct {
}

func (h *LDAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.A = value
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)
	cpu.PC++

	return nil
}

type LSRHandler struct {
}

func (h *LSRHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	leftBit := value & 0x1
	value >>= 1
	err = cpu.writeWithMemoryAccessType(mode, operand, value)
	if err != nil {
		return err
	}

	cpu.PC++

	cpu.P.UpdateN(value)
	cpu.P.UpdateZ(value)

	if leftBit == 0 {
		cpu.P.ClearC()
	} else {
		cpu.P.SetC()
	}

	return nil
}

type RLAHandler struct {
}

func (h *RLAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	var byteC byte
	if cpu.P.IsC() {
		byteC = 1
	} else {
		byteC = 0
	}

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := (uint16(src) << 1) | uint16(byteC)

	err = cpu.writeWithMemoryAccessType(mode, operand, byte(result))
	if err != nil {
		return err
	}

	if result >= 0xFF {
		cpu.P.SetC()
	} else {
		cpu.P.ClearC()
	}

	cpu.A &= byte(result)
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)

	cpu.PC++

	return nil
}

type RORHandler struct {
}

func (h *RORHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	var cByte byte
	if cpu.P.IsC() {
		cByte = 0x80
	} else {
		cByte = 0x0
	}

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := (src >> 1) | cByte

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	if src&1 == 0 {
		cpu.P.ClearC()
	} else {
		cpu.P.SetC()
	}

	cpu.PC++

	return nil
}

type RRAHandler struct {
}

func (h *RRAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.PC++

	return nil
}

type SEIHandler struct {
}

func (h *SEIHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.SetI()
	cpu.PC++

	return nil
}

type STAHandler struct {
}

func (h *STAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	err := cpu.writeWithMemoryAccessType(mode, operand, cpu.A)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type TYAHandler struct {
}

func (h *TYAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.A = cpu.Y
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)
	cpu.PC++

	return nil
}

type LDXHandler struct {
}

func (h *LDXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.X = value
	cpu.P.UpdateZ(cpu.X)
	cpu.P.UpdateN(cpu.X)

	cpu.PC++

	return nil
}

type BNEHandler struct {
}

func (h *BNEHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if !cpu.P.IsZ() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type ADCHandler struct {
}

func (h *ADCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	var cByte byte
	if cpu.P.IsC() {
		cByte = 1
	} else {
		cByte = 0
	}

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := uint16(src1) + uint16(src2) + uint16(cByte)
	cpu.A = byte(result)

	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)
	cpu.P.UpdateC(result)

	if !((src1^src2)&0x80 > 0) && ((uint16(src2)^result)&0x80 > 0) {
		cpu.P.SetV()
	} else {
		cpu.P.ClearV()
	}

	cpu.PC++

	return nil
}

type TXSHandler struct {
}

func (h *TXSHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.S.SetValue(cpu.X)
	cpu.PC++

	return nil
}

type RTSHandler struct {
}

func (h *RTSHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := cpu.S.PopUint16()

	cpu.PC = value
	cpu.PC++

	return nil
}

type BCCHandler struct {
}

func (h *BCCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if !cpu.P.IsC() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type BCSHandler struct {
}

func (h *BCSHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if cpu.P.IsC() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type BEQHandler struct {
}

func (h *BEQHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if cpu.P.IsZ() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type BITHandler struct {
}

func (h *BITHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 & src2

	cpu.P.UpdateN(src2)
	cpu.P.UpdateZ(result)

	if (src2 & 0x40) == 0 {
		cpu.P.ClearV()
	} else {
		cpu.P.SetV()
	}

	cpu.PC++

	return nil
}

type BMIHandler struct {
}

func (h *BMIHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if cpu.P.IsN() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type BRKHandler struct {
}

func (h *BRKHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.SetB()
	cpu.PC++
	cpu.interrupt(BRK)

	return nil
}

type BVCHandler struct {
}

func (h *BVCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if !cpu.P.IsV() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type BVSHandler struct {
}

func (h *BVSHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value := int8(operand)

	if cpu.P.IsV() {
		if value > 0 {
			cpu.PC -= uint16(-value)
		} else {
			cpu.PC += uint16(value)
		}
	}

	cpu.PC++

	return nil
}

type CLCHandler struct {
}

func (h *CLCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.ClearC()
	cpu.PC++

	return nil
}

type CLIHandler struct {
}

func (h *CLIHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.ClearI()
	cpu.PC++

	return nil
}

type CLVHandler struct {
}

func (h *CLVHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.ClearV()
	cpu.PC++

	return nil
}

type CMPHandler struct {
}

func (h *CMPHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 - src2
	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	if src1 >= src2 {
		cpu.P.SetC()
	} else {
		cpu.P.ClearC()
	}

	cpu.PC++

	return nil
}

type CPXHandler struct {
}

func (h *CPXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.X
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 - src2
	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	if src1 >= src2 {
		cpu.P.SetC()
	} else {
		cpu.P.ClearC()
	}

	cpu.PC++

	return nil
}

type CPYHandle struct {
}

func (h *CPYHandle) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.Y
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 - src2
	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	if src1 >= src2 {
		cpu.P.SetC()
	} else {
		cpu.P.ClearC()
	}

	cpu.PC++

	return nil
}

type DECHandler struct {
}

func (h *DECHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src - 1

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	cpu.PC++

	return nil
}

type DEYHandler struct {
}

func (h *DEYHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.Y--
	cpu.P.UpdateN(cpu.Y)
	cpu.P.UpdateZ(cpu.Y)

	cpu.PC++

	return nil
}

type INCHandler struct {
}

func (h *INCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src + 1

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	cpu.PC++

	return nil
}

type INXHandler struct {
}

func (h *INXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.X++
	cpu.P.UpdateN(cpu.X)
	cpu.P.UpdateZ(cpu.X)

	cpu.PC++

	return nil
}

type INYHandler struct {
}

func (h *INYHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.Y++
	cpu.P.UpdateN(cpu.Y)
	cpu.P.UpdateZ(cpu.Y)

	cpu.PC++

	return nil
}

type LDYHandler struct {
}

func (h *LDYHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.Y = value
	cpu.P.UpdateN(cpu.Y)
	cpu.P.UpdateZ(cpu.Y)

	cpu.PC++

	return nil
}

type NOPHandler struct {
}

func (h *NOPHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.PC++

	return nil
}

type ORAHandler struct {
}

func (h *ORAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 | src2
	cpu.A = result
	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	cpu.PC++

	return nil
}

type PHAHandler struct {
}

func (h *PHAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.S.PushByte(cpu.A)

	cpu.PC++

	return nil
}

type PHPHandler struct {
}

func (h *PHPHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.S.PushByte(cpu.P.GetValue())

	cpu.PC++

	return nil
}

type PLAHandler struct {
}

func (h *PLAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.A = cpu.S.PopByte()
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)

	cpu.PC++

	return nil
}

type PLPHandler struct {
}

func (h *PLPHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	val := cpu.S.PopByte()
	if cpu.P.IsB() {
		val |= 0x10
	} else {
		val &= ^byte(0x10)
	}

	val |= 0x20

	cpu.P.SetValue(val)
	cpu.PC++

	return nil
}

type ROLHandler struct {
}

func (h *ROLHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	var byteC byte
	if cpu.P.IsC() {
		byteC = 1
	} else {
		byteC = 0
	}

	result := (uint16(src) << 1) | uint16(byteC)

	err = cpu.writeWithMemoryAccessType(mode, operand, byte(result))
	if err != nil {
		return err
	}

	cpu.P.UpdateN(byte(result))
	cpu.P.UpdateZ(byte(result))
	cpu.P.UpdateC(result)

	cpu.PC++

	return nil
}

type RTIHandler struct {
}

func (h *RTIHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.SetValue(cpu.S.PopByte())
	cpu.PC = cpu.S.PopUint16()

	return nil
}

type SBCHandler struct {
}

func (h *SBCHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	var byteC byte
	if cpu.P.IsC() {
		byteC = 0
	} else {
		byteC = 1
	}

	result := src1 - src2 - byteC
	cpu.A = result

	cpu.P.UpdateN(result)
	cpu.P.UpdateZ(result)

	if src1 >= src2+byteC {
		cpu.P.SetC()
	} else {
		cpu.P.ClearC()
	}

	if ((src1^result)&0x80 > 0) && ((src1^src2)&0x80 > 0) {
		cpu.P.SetV()
	} else {
		cpu.P.ClearV()
	}

	cpu.PC++

	return nil
}

type SECHandler struct {
}

func (h *SECHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.SetC()
	cpu.PC++

	return nil
}

type SEDHandler struct {
}

func (h *SEDHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.P.SetD()
	cpu.PC++

	return nil
}

type STXHandler struct {
}

func (h *STXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	err := cpu.writeWithMemoryAccessType(mode, operand, cpu.X)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type STYHandler struct {
}

func (h *STYHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	err := cpu.writeWithMemoryAccessType(mode, operand, cpu.Y)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type TAXHandler struct {
}

func (h *TAXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.X = cpu.A
	cpu.P.UpdateN(cpu.X)
	cpu.P.UpdateZ(cpu.X)

	cpu.PC++

	return nil
}

type TAYHandler struct {
}

func (h *TAYHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.Y = cpu.A
	cpu.P.UpdateN(cpu.Y)
	cpu.P.UpdateZ(cpu.Y)

	cpu.PC++

	return nil
}

type TSXHandler struct {
}

func (h *TSXHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.X = cpu.S.GetValue()
	cpu.P.UpdateN(cpu.X)
	cpu.P.UpdateZ(cpu.X)

	cpu.PC++

	return nil
}

type TXAHandler struct {
}

func (h *TXAHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.A = cpu.X
	cpu.P.UpdateN(cpu.A)
	cpu.P.UpdateZ(cpu.A)

	cpu.PC++

	return nil
}

type SLOHandler struct {
}

// Handle @todo: implement
func (h *SLOHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.PC++

	return nil
}

type SREHandler struct {
}

func (h *SREHandler) Handle(cpu *Cpu, operand uint16, mode enums.Modes) error {
	cpu.PC++

	return nil
}
