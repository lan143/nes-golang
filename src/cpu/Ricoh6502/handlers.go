package Ricoh6502

import (
	"main/src/cpu/Ricoh6502/enums"
)

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
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BPL", mode, operand)

	value := int8(operand)

	if cpu.P&N == 0 {
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

	cpu.X--
	cpu.setFlagsByValue(cpu.X)

	cpu.PC++

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
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("JSR", mode, operand)

	stackValue := cpu.PC

	cpu.setByte(uint16(cpu.S)+0x100, byte((stackValue>>8)&0xff))
	cpu.S--
	cpu.setByte(uint16(cpu.S)+0x100, byte(stackValue&0xff))
	cpu.S--

	cpu.PC = operand

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

type BNEHandler struct {
}

func (h *BNEHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BNE", mode, operand)

	value := int8(operand)

	if cpu.P&Z == 0 {
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

func (h *ADCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ADC", mode, operand)

	var cByte byte
	if cpu.P&byte(C) > 0 {
		cByte = 1
	} else {
		cByte = 0
	}

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := uint16(src1) + uint16(src2) + uint16(cByte)
	cpu.A = byte(result)

	cpu.setFlagsByValue(cpu.A)

	if result >= 0xff {
		cpu.setCorrectionBit(1)
	} else {
		cpu.setCorrectionBit(0)
	}

	if !((src1^src2)&0x80 > 0) && ((uint16(src2)^result)&0x80 > 0) {
		cpu.P |= V
	} else {
		cpu.P &= ^byte(V)
	}

	cpu.PC++

	return nil
}

type TXSHandler struct {
}

func (h *TXSHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TXS", mode, 0)

	cpu.S = cpu.X
	cpu.PC++

	return nil
}

type RTSHandler struct {
}

func (h *RTSHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("RTS", mode, 0)

	cpu.S++
	byte1 := cpu.getByte(uint16(cpu.S) + 0x100)
	cpu.S++
	byte2 := cpu.getByte(uint16(cpu.S) + 0x100)

	cpu.PC = (uint16(byte2) << 8) | uint16(byte1)
	cpu.PC++

	return nil
}

type BCCHandler struct {
}

func (h *BCCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BCC", mode, operand)

	value := int8(operand)

	if cpu.P&byte(C) == 0 {
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

func (h *BCSHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BCS", mode, operand)

	value := int8(operand)

	if cpu.P&byte(C) != 0 {
		if cpu.P&N == 0 {
			if value > 0 {
				cpu.PC -= uint16(-value)
			} else {
				cpu.PC += uint16(value)
			}
		}
	}

	cpu.PC++

	return nil
}

type BEQHandler struct {
}

func (h *BEQHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BEQ", mode, operand)

	value := int8(operand)

	if cpu.P&byte(Z) != 0 {
		if cpu.P&N == 0 {
			if value > 0 {
				cpu.PC -= uint16(-value)
			} else {
				cpu.PC += uint16(value)
			}
		}
	}

	cpu.PC++

	return nil
}

type BITHandler struct {
}

func (h *BITHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BIT", mode, operand)

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src1 & src2

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	if (src2 & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if (src2 & 0x40) == 0 {
		cpu.P &= ^byte(V)
	} else {
		cpu.P |= V
	}

	return nil
}

type BMIHandler struct {
}

func (h *BMIHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BMI", mode, operand)

	value := int8(operand)

	if cpu.P&byte(N) != 0 {
		if cpu.P&N == 0 {
			if value > 0 {
				cpu.PC -= uint16(-value)
			} else {
				cpu.PC += uint16(value)
			}
		}
	}

	cpu.PC++

	return nil
}

type BRKHandler struct {
}

func (h *BRKHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("BRK", mode, 0)

	cpu.P |= byte(B)
	cpu.P |= byte(I)
	cpu.PC++

	// @todo: implement run BRK handler

	return nil
}

type BVCHandler struct {
}

func (h *BVCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BVC", mode, operand)

	value := int8(operand)

	if cpu.P&byte(V) == 0 {
		if cpu.P&N == 0 {
			if value > 0 {
				cpu.PC -= uint16(-value)
			} else {
				cpu.PC += uint16(value)
			}
		}
	}

	cpu.PC++

	return nil
}

type BVSHandler struct {
}

func (h *BVSHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("BVS", mode, operand)

	value := int8(operand)

	if cpu.P&byte(V) != 0 {
		if cpu.P&N == 0 {
			if value > 0 {
				cpu.PC -= uint16(-value)
			} else {
				cpu.PC += uint16(value)
			}
		}
	}

	cpu.PC++

	return nil
}

type CLCHandler struct {
}

func (h *CLCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("CLC", mode, 0)

	cpu.P &= ^byte(C)
	cpu.PC++

	return nil
}

type CLIHandler struct {
}

func (h *CLIHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("CLI", mode, 0)

	cpu.P &= ^byte(I)
	cpu.PC++

	return nil
}

type CLVHandler struct {
}

func (h *CLVHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("CLV", mode, 0)

	cpu.P &= ^byte(V)
	cpu.PC++

	return nil
}

type CMPHandler struct {
}

func (h *CMPHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("CMP", mode, operand)

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src1 - src2

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if src1 >= src2 {
		cpu.P |= byte(C)
	} else {
		cpu.P &= ^byte(C)
	}

	cpu.PC++

	return nil
}

type CPXHandler struct {
}

func (h *CPXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("CPX", mode, operand)

	src1 := cpu.X
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src1 - src2

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if src1 >= src2 {
		cpu.P |= byte(C)
	} else {
		cpu.P &= ^byte(C)
	}

	cpu.PC++

	return nil
}

type CPYHandle struct {
}

func (h *CPYHandle) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("CPY", mode, operand)

	src1 := cpu.Y
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src1 - src2

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if src1 >= src2 {
		cpu.P |= byte(C)
	} else {
		cpu.P &= ^byte(C)
	}

	return nil
}

type DECHandler struct {
}

func (h *DECHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("DEC", mode, operand)

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src - 1

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type DEYHandler struct {
}

func (h *DEYHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("DEY", mode, 0)

	cpu.PC++

	cpu.Y--
	cpu.setFlagsByValue(cpu.Y)

	return nil
}

type INCHandler struct {
}

func (h *INCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("INC", mode, operand)

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	result := src + 1

	err = cpu.writeWithMemoryAccessType(mode, operand, result)
	if err != nil {
		return err
	}

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type INXHandler struct {
}

func (h *INXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("INX", mode, 0)

	cpu.X++

	if (cpu.X & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.X & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type INYHandler struct {
}

func (h *INYHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("INY", mode, 0)

	cpu.Y++

	if (cpu.Y & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.Y & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type LDYHandler struct {
}

func (h *LDYHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("LDY", mode, operand)

	value, err := cpu.loadWithMemoryAccessType(mode, operand)

	cpu.Y = value

	if (cpu.Y & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.Y & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type NOPHandler struct {
}

func (h *NOPHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("NOP", mode, 0)
	cpu.PC++

	return nil
}

type ORAHandler struct {
}

func (h *ORAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ORA", mode, operand)

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	result := src1 | src2
	cpu.A = result

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type PHAHandler struct {
}

func (h *PHAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("PHA", mode, 0)

	cpu.setByte(uint16(cpu.S)+0x100, cpu.A)
	cpu.S--

	cpu.PC++

	return nil
}

type PHPHandler struct {
}

func (h *PHPHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("PHP", mode, 0)

	cpu.setByte(uint16(cpu.S)+0x100, cpu.P)
	cpu.S--

	cpu.PC++

	return nil
}

type PLAHandler struct {
}

func (h *PLAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("PLA", mode, 0)

	cpu.S++
	cpu.A = cpu.getByte(uint16(cpu.S) + 0x100)

	if (cpu.A & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.A & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type PLPHandler struct {
}

func (h *PLPHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("PLP", mode, 0)

	cpu.S++
	cpu.P = cpu.getByte(uint16(cpu.S) + 0x100)

	cpu.PC++

	return nil
}

type ROLHandler struct {
}

func (h *ROLHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("ROL", mode, operand)

	src, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	var byteC byte
	if cpu.P&byte(C) > 0 {
		byteC = 1
	} else {
		byteC = 0
	}

	result := (uint16(src) << 1) | uint16(byteC)

	err = cpu.writeWithMemoryAccessType(mode, operand, byte(result))
	if err != nil {
		return err
	}

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if result&0x100 == 0 {
		cpu.P &= ^byte(C)
	} else {
		cpu.P |= byte(C)
	}

	cpu.PC++

	return nil
}

type RTIHandler struct {
}

func (h *RTIHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("RTI", mode, 0)

	cpu.S++
	cpu.S = cpu.getByte(uint16(cpu.S) + 0x100)

	cpu.S++
	byte1 := cpu.getByte(uint16(cpu.S) + 0x100)
	cpu.S++
	byte2 := cpu.getByte(uint16(cpu.S) + 0x100)

	cpu.PC = (uint16(byte2) << 8) | uint16(byte1)
	cpu.PC++

	return nil
}

type SBCHandler struct {
}

func (h *SBCHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("SBC", mode, operand)

	src1 := cpu.A
	src2, err := cpu.loadWithMemoryAccessType(mode, operand)
	if err != nil {
		return err
	}

	var byteC byte
	if cpu.P&byte(C) > 0 {
		byteC = 0
	} else {
		byteC = 1
	}

	result := src1 - src2 - byteC
	cpu.A = result

	if (result & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (result & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	if src1 >= src2+byteC {
		cpu.P |= byte(C)
	} else {
		cpu.P &= ^byte(C)
	}

	if ((src1^result)&0x80 > 0) && ((src1^src2)&0x80 > 0) {
		cpu.P |= byte(V)
	} else {
		cpu.P &= ^byte(V)
	}

	cpu.PC++

	return nil
}

type SECHandler struct {
}

func (h *SECHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("SEC", mode, 0)

	cpu.P |= byte(C)
	cpu.PC++

	return nil
}

type SEDHandler struct {
}

func (h *SEDHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("SED", mode, 0)

	cpu.P |= byte(D)
	cpu.PC++

	return nil
}

type STXHandler struct {
}

func (h *STXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("STX", mode, operand)

	err = cpu.writeWithMemoryAccessType(mode, operand, cpu.X)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type STYHandler struct {
}

func (h *STYHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("STY", mode, operand)

	err = cpu.writeWithMemoryAccessType(mode, operand, cpu.Y)
	if err != nil {
		return err
	}

	cpu.PC++

	return nil
}

type TAXHandler struct {
}

func (h *TAXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TAX", mode, 0)

	cpu.X = cpu.A

	if (cpu.X & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.X & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type TAYHandler struct {
}

func (h *TAYHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TAY", mode, 0)

	cpu.Y = cpu.A

	if (cpu.Y & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.Y & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type TSXHandler struct {
}

func (h *TSXHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TSX", mode, 0)

	cpu.X = cpu.S

	if (cpu.X & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.X & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type TXAHandler struct {
}

func (h *TXAHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	cpu.logExecution("TXA", mode, 0)

	cpu.A = cpu.X

	if (cpu.A & 0x80) == 0 {
		cpu.P &= ^byte(N)
	} else {
		cpu.P |= N
	}

	if (cpu.A & 0xff) == 0 {
		cpu.P |= Z
	} else {
		cpu.P &= ^byte(Z)
	}

	cpu.PC++

	return nil
}

type SLOHandler struct {
}

// Handle @todo: implement
func (h *SLOHandler) Handle(cpu *Cpu, mode enums.Modes) error {
	operand, err := cpu.loadInstructionOperand(mode)
	if err != nil {
		return err
	}

	cpu.logExecution("SLO", mode, operand)
	cpu.PC++

	return nil
}
