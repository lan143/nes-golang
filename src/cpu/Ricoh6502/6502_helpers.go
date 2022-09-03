package Ricoh6502

import (
	"fmt"
	"log"
	"main/src/cpu/Ricoh6502/enums"
)

func (c *Cpu) logExecution(opcode string, mode enums.Modes, operand uint16) {
	switch mode {
	case enums.ModeREL:
		log.Printf("%s %X", opcode, int8(operand))
		break
	case enums.ModeIND:
		log.Printf("%s $%X", opcode, operand)
		break
	case enums.ModeIMM:
		log.Printf("%s #%X", opcode, operand)
		break
	case enums.ModeAcc:
		log.Printf("%s A", opcode)
		break
	case enums.ModeZP:
		log.Printf("%s $%X", opcode, operand)
		break
	case enums.ModeZPX:
		log.Printf("%s $%X,X", opcode, operand)
		break
	case enums.ModeZPY:
		log.Printf("%s $%X,X", opcode, operand)
		break
	case enums.ModeINDX:
		log.Printf("%s $%X,X", opcode, operand)
		break
	case enums.ModeINDY:
		log.Printf("%s $%X,Y", opcode, operand)
		break
	case enums.ModeABS:
		log.Printf("%s $%X", opcode, operand)
		break
	case enums.ModeABSX:
		log.Printf("%s $%X,X", opcode, operand)
		break
	case enums.ModeABSY:
		log.Printf("%s $%X,Y", opcode, operand)
		break
	case enums.ModeIMP:
		log.Printf("%s", opcode)
		break
	default:
		log.Printf("%s %X", opcode, operand)
	}
}

func (c *Cpu) loadInstructionOperand(mode enums.Modes) (uint16, error) {
	switch mode {
	case enums.ModeIND:
		return uint16(c.getByte(c.getNextUint16())), nil
	case enums.ModeAcc:
		return 0, nil
	case enums.ModeIMM:
		return uint16(c.getNextByte()), nil
	case enums.ModeZP:
		return uint16(c.getNextByte()), nil
	case enums.ModeZPX:
		return uint16(c.getNextByte()), nil
	case enums.ModeZPY:
		return uint16(c.getNextByte()), nil
	case enums.ModeINDX:
		return uint16(c.getNextByte()), nil
	case enums.ModeINDY:
		return uint16(c.getNextByte()), nil
	case enums.ModeABS:
		return c.getNextUint16(), nil
	case enums.ModeABSX:
		return c.getNextUint16(), nil
	case enums.ModeABSY:
		return c.getNextUint16(), nil
	case enums.ModeIMP:
		return 0, nil
	case enums.ModeREL:
		return uint16(c.getNextByte()), nil
	default:
		return 0, fmt.Errorf("loadInstructionOperand: unsupported %d memory access type", mode)
	}
}

func (c *Cpu) loadWithMemoryAccessType(mode enums.Modes, operand uint16) (byte, error) {
	switch mode {
	case enums.ModeIMM:
		return c.getByte(c.PC), nil
	case enums.ModeAcc:
		return c.A, nil
	case enums.ModeZP:
		return c.getByte(operand), nil
	case enums.ModeZPX:
		return c.getByte(operand + uint16(c.X)), nil
	case enums.ModeZPY:
		return c.getByte(operand + uint16(c.Y)), nil
	case enums.ModeABS:
		return c.getByte(operand), nil
	case enums.ModeABSX:
		return c.getByte(operand + uint16(c.X)), nil
	case enums.ModeABSY:
		return c.getByte(operand + uint16(c.Y)), nil
	case enums.ModeINDX:
		return c.getByte(operand + uint16(c.X)), nil
	case enums.ModeINDY:
		address := c.getByte(operand)

		return c.getByte(uint16(address) + uint16(c.Y)), nil
	case enums.ModeIMP:
		return 0, nil
	default:
		return 0, fmt.Errorf("loadWithMemoryAccessType: unsupported %d memory access type", mode)
	}
}

func (c *Cpu) writeWithMemoryAccessType(mode enums.Modes, operand uint16, value byte) error {
	switch mode {
	case enums.ModeAcc:
		c.A = value
		break
	case enums.ModeZP:
		c.setByte(operand, value)
		break
	case enums.ModeZPX:
		c.setByte(operand+uint16(c.X), value)
		break
	case enums.ModeZPY:
		c.setByte(operand+uint16(c.Y), value)
		break
	case enums.ModeABS:
		c.setByte(operand, value)
		break
	case enums.ModeABSX:
		c.setByte(operand+uint16(c.X), value)
		break
	case enums.ModeABSY:
		c.setByte(operand+uint16(c.Y), value)
		break
	case enums.ModeINDX:
		c.setByte(operand+uint16(c.X), value)
		break
	case enums.ModeINDY:
		address := c.getByte(operand)
		c.setByte(uint16(address)+uint16(c.Y), value)
		break
	default:
		return fmt.Errorf("writeWithMemoryAccessType: unsupported %d memory access type", mode)
	}

	return nil
}
