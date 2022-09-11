package Ricoh6502

import (
	"fmt"
	"log"
	"main/src/cpu/Ricoh6502/enums"
)

func (c *Cpu) logExecution(position uint16, opcode string, mode enums.Modes, operand uint16) {
	switch mode {
	case enums.ModeREL:
		log.Printf("(0x%04X) (REL) %s %04X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, int8(operand), c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeIND:
		log.Printf("(0x%04X) (IND) %s $%04X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeIMM:
		log.Printf("(0x%04X) (IMM) %s #%04X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeAcc:
		log.Printf("(0x%04X) %s A	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeZP:
		log.Printf("(0x%04X) (ZP) %s $%04X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeZPX:
		log.Printf("(0x%04X) (ZPX) %s $%04X,X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeZPY:
		log.Printf("(0x%04X) (ZPY) %s $%04X,X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeINDX:
		log.Printf("(0x%04X) (INDX) %s $%02X,X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeINDY:
		log.Printf("(0x%04X) (INDY) %s $%04X,Y	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeABS:
		log.Printf("(0x%04X) (ABS) %s $%04X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeABSX:
		log.Printf("(0x%04X) (ABSX) %s $%04X,X	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeABSY:
		log.Printf("(0x%04X) (ABSY) %s $%04X,Y	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, operand, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	case enums.ModeIMP:
		log.Printf("(0x%04X) (IMP) %s	A:%02X X:%02X Y:%02X P:%02X SP:%02X", position, opcode, c.A, c.X, c.Y, c.P.GetValue(), c.S.value)
		break
	default:
		log.Printf("(0x%04X) %s %04X", position, opcode, operand)
	}
}

func (c *Cpu) loadInstructionOperand(mode enums.Modes) (uint16, error) {
	switch mode {
	case enums.ModeIND:
		return c.getNextUint16(), nil
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
		return byte(operand), nil
	case enums.ModeIND:
		return c.getByte(operand), nil
	case enums.ModeAcc:
		return c.A, nil
	case enums.ModeZP:
		return c.getByte(operand), nil
	case enums.ModeZPX:
		return c.getByte((operand + uint16(c.X)) & 0xFF), nil
	case enums.ModeZPY:
		return c.getByte((operand + uint16(c.Y)) & 0xFF), nil
	case enums.ModeABS:
		return c.getByte(operand), nil
	case enums.ModeABSX:
		return c.getByte(operand + uint16(c.X)), nil
	case enums.ModeABSY:
		return c.getByte(operand + uint16(c.Y)), nil
	case enums.ModeINDX:
		address := (operand + uint16(c.X)) & 0xff

		return c.getByte(c.getUin16(address)), nil
	case enums.ModeINDY:
		address := c.getUin16(operand)
		address += uint16(c.Y)

		return c.getByte(address), nil
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
		c.setByte((operand+uint16(c.X))&0xFF, value)
		break
	case enums.ModeZPY:
		c.setByte((operand+uint16(c.Y))&0xFF, value)
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
		address := (operand + uint16(c.X)) & 0xff
		c.setByte(c.getUin16(address), value)
		break
	case enums.ModeINDY:
		address := c.getUin16(operand)
		address += uint16(c.Y)

		c.setByte(address, value)
		break
	default:
		return fmt.Errorf("writeWithMemoryAccessType: unsupported %d memory access type", mode)
	}

	return nil
}
