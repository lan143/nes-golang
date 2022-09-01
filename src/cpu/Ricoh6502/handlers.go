package Ricoh6502

import "main/src/cpu/Ricoh6502/enums"

type ANDHandler struct {
}

func (h *ANDHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("AND", mode, operand)

	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.A &= value

	cpu.setFlagsByValue(cpu.A)
	cpu.PC++

	return nil
}

type ASLHandler struct {
}

func (h *ASLHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ASL", mode, operand)

	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	leftBit := value & 0x80
	value <<= 1
	err = cpu.writeWithMemoryAccessType(mode, operand, value)
	if err != nil {
		return err
	}

	cpu.PC++
	cpu.setFlagsByValue(value)
	cpu.setCorrectionBit(leftBit)

	return nil
}

type BPLHandler struct {
}

func (h *BPLHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	value := int8(cpu.getNextByte())
	cpu.logExecution("BPL", mode, uint16(value))

	if cpu.P&N == 0 {
		if value > 0 {
			cpu.PC += uint16(value)
		} else {
			cpu.PC -= uint16(-value)
		}
	}

	cpu.PC++

	return nil
}

type CLDHandler struct {
}

func (h *CLDHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("CLD", mode, 0)

	cpu.P &= ^uint8(D)
	cpu.PC++

	return nil
}

type DEXHandler struct {
}

func (h *DEXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("DEX", mode, 0)

	cpu.PC++

	cpu.X--
	cpu.setFlagsByValue(cpu.X)

	return nil
}

type EORHandler struct {
}

func (h *EORHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("EOR", mode, operand)

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	cpu.A ^= src
	cpu.setFlagsByValue(cpu.A)

	cpu.PC++

	return nil
}

type ISBHandler struct {
}

func (h *ISBHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ISB", mode, operand)

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
	cpu.setFlagsByValue(cpu.A)

	return nil
}

type JMPHandler struct {
}

func (h *JMPHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	address, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("JMP", mode, address)

	cpu.PC = address

	return nil
}

type JSRHandler struct {
}

func (h *JSRHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand := cpu.getNextUint16()

	cpu.logExecution("JSR", mode, operand)

	value := cpu.PC + 1
	cpu.PC = operand

	cpu.setByte(uint16(cpu.S), byte((value>>8)&0xff))
	cpu.S--
	cpu.setByte(uint16(cpu.S), byte(value&0xff))
	cpu.S--

	return nil
}

type LDAHandler struct {
}

func (h *LDAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("LDA", mode, operand)

	value, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	cpu.A = value
	cpu.setFlagsByValue(cpu.A)
	cpu.PC++

	return nil
}

type LSRHandler struct {
}

func (h *LSRHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("LSR", mode, operand)

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

	cpu.setFlagsByValue(value)
	cpu.setCorrectionBit(leftBit)

	return nil
}

type RLAHandler struct {
}

func (h *RLAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("RLA", mode, operand)

	var byteC byte
	if cpu.P&byte(C) > 0 {
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
		cpu.P |= byte(C)
	} else {
		cpu.P &= ^byte(C)
	}

	cpu.A &= byte(result)
	cpu.setFlagsByValue(cpu.A)
	cpu.PC++

	return nil
}

type RORHandler struct {
}

func (h *RORHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ROR", mode, operand)

	var cByte byte
	if cpu.P&byte(C) > 0 {
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

	if src&1 == 0 {
		cpu.P &= ^byte(C)
		cByte = 0
	} else {
		cpu.P |= byte(C)
		cByte = 1
	}

	cpu.PC++

	return nil
}

type RRAHandler struct {
}

func (h *RRAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("RRA", mode, operand)

	var cByte byte
	if cpu.P&byte(C) > 0 {
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

	if src&1 == 0 {
		cpu.P &= ^byte(C)
		cByte = 0
	} else {
		cpu.P |= byte(C)
		cByte = 1
	}

	src1 := cpu.A
	src2 := result
	result1 := uint16(src1) + uint16(src2) + uint16(cByte)
	cpu.A = byte(result1)

	cpu.setFlagsByValue(cpu.A)

	if result1 >= 0xff {
		cpu.setCorrectionBit(1)
	} else {
		cpu.setCorrectionBit(0)
	}

	if !((src1^src2)&0x80 > 0) && ((src2^result)&0x80 > 0) {
		cpu.P |= V
	} else {
		cpu.P &= ^byte(V)
	}

	cpu.PC++

	return nil
}

type SEIHandler struct {
}

func (h *SEIHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("SEI", mode, 0)

	cpu.P |= 0x20
	cpu.PC++

	return nil
}

type STAHandler struct {
}

func (h *STAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("STA", mode, operand)

	err = cpu.writeWithMemoryAccessType(mode, operand, cpu.A)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type TYAHandler struct {
}

func (h *TYAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TYA", mode, 0)

	cpu.A = cpu.Y
	cpu.PC++

	return nil
}

type XTAHandler struct {
}

func (h *XTAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("XTA", mode, 0)

	cpu.A = cpu.X
	cpu.PC++

	return nil
}

type LDXHandler struct {
}

func (h *LDXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("LDX", mode, operand)

	value, err := cpu.loadWithMemoryAccessType(mode, operand)

	cpu.X = value
	cpu.PC++

	return nil
}
