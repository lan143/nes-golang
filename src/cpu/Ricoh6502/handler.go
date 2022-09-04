package Ricoh6502

import (
	"main/src/cpu/Ricoh6502/enums"
)

type Handler interface {
	Handle(cpu *Cpu, operand uint16, mode enums.Modes) error
}
