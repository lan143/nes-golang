package Ricoh6502

import (
	"main/src/cpu/Ricoh6502/enums"
)

type CommandDescription struct {
	Command    byte
	OpcodeName string
	Mode       enums.Modes
	Handler    Handler
}

type Decoder struct {
	Commands []CommandDescription
}

func (c *Decoder) InitCommands() {
	c.Commands = []CommandDescription{
		{Command: 0x98, OpcodeName: "TYA", Mode: enums.ModeIMP, Handler: &TYAHandler{}},
		{Command: 0x20, OpcodeName: "JSR", Mode: enums.ModeABS, Handler: &JSRHandler{}},
		{Command: 0x4C, OpcodeName: "JMP", Mode: enums.ModeABS, Handler: &JMPHandler{}},
		{Command: 0x6C, OpcodeName: "JMP", Mode: enums.ModeIND, Handler: &JMPHandler{}},
		{Command: 0x4C, OpcodeName: "JMP", Mode: enums.ModeABS, Handler: &JMPHandler{}},
		{Command: 0x10, OpcodeName: "BPL", Mode: enums.ModeREL, Handler: &BPLHandler{}},
		{Command: 0x78, OpcodeName: "SEI", Mode: enums.ModeIMP, Handler: &SEIHandler{}},
		{Command: 0xD8, OpcodeName: "CLD", Mode: enums.ModeIMP, Handler: &CLDHandler{}},
		{Command: 0x49, OpcodeName: "EOR", Mode: enums.ModeIMM, Handler: &EORHandler{}},
		{Command: 0x45, OpcodeName: "EOR", Mode: enums.ModeZP, Handler: &EORHandler{}},
		{Command: 0x55, OpcodeName: "EOR", Mode: enums.ModeZPX, Handler: &EORHandler{}},
		{Command: 0x4D, OpcodeName: "EOR", Mode: enums.ModeABS, Handler: &EORHandler{}},
		{Command: 0x5D, OpcodeName: "EOR", Mode: enums.ModeABSX, Handler: &EORHandler{}},
		{Command: 0x59, OpcodeName: "EOR", Mode: enums.ModeABSY, Handler: &EORHandler{}},
		{Command: 0x41, OpcodeName: "EOR", Mode: enums.ModeINDX, Handler: &EORHandler{}},
		{Command: 0x51, OpcodeName: "EOR", Mode: enums.ModeINDY, Handler: &EORHandler{}},
		{Command: 0x29, OpcodeName: "AND", Mode: enums.ModeIMM, Handler: &ANDHandler{}},
		{Command: 0x25, OpcodeName: "AND", Mode: enums.ModeZP, Handler: &ANDHandler{}},
		{Command: 0x35, OpcodeName: "AND", Mode: enums.ModeZPX, Handler: &ANDHandler{}},
		{Command: 0x2D, OpcodeName: "AND", Mode: enums.ModeABS, Handler: &ANDHandler{}},
		{Command: 0x3D, OpcodeName: "AND", Mode: enums.ModeABSX, Handler: &ANDHandler{}},
		{Command: 0x39, OpcodeName: "AND", Mode: enums.ModeABSY, Handler: &ANDHandler{}},
		{Command: 0x21, OpcodeName: "AND", Mode: enums.ModeINDX, Handler: &ANDHandler{}},
		{Command: 0x31, OpcodeName: "AND", Mode: enums.ModeINDY, Handler: &ANDHandler{}},
		{Command: 0x85, OpcodeName: "STA", Mode: enums.ModeZP, Handler: &STAHandler{}},
		{Command: 0x95, OpcodeName: "STA", Mode: enums.ModeZPX, Handler: &STAHandler{}},
		{Command: 0x8D, OpcodeName: "STA", Mode: enums.ModeABS, Handler: &STAHandler{}},
		{Command: 0x9D, OpcodeName: "STA", Mode: enums.ModeABSX, Handler: &STAHandler{}},
		{Command: 0x99, OpcodeName: "STA", Mode: enums.ModeABSY, Handler: &STAHandler{}},
		{Command: 0x81, OpcodeName: "STA", Mode: enums.ModeINDX, Handler: &STAHandler{}},
		{Command: 0x91, OpcodeName: "STA", Mode: enums.ModeINDY, Handler: &STAHandler{}},
		{Command: 0xA9, OpcodeName: "LDA", Mode: enums.ModeIMM, Handler: &LDAHandler{}},
		{Command: 0xA5, OpcodeName: "LDA", Mode: enums.ModeZP, Handler: &LDAHandler{}},
		{Command: 0xB5, OpcodeName: "LDA", Mode: enums.ModeZPX, Handler: &LDAHandler{}},
		{Command: 0xAD, OpcodeName: "LDA", Mode: enums.ModeABS, Handler: &LDAHandler{}},
		{Command: 0xBD, OpcodeName: "LDA", Mode: enums.ModeABSX, Handler: &LDAHandler{}},
		{Command: 0xB9, OpcodeName: "LDA", Mode: enums.ModeABSY, Handler: &LDAHandler{}},
		{Command: 0xA1, OpcodeName: "LDA", Mode: enums.ModeINDX, Handler: &LDAHandler{}},
		{Command: 0xB1, OpcodeName: "LDA", Mode: enums.ModeINDY, Handler: &LDAHandler{}},
		{Command: 0x8A, OpcodeName: "XTA", Mode: enums.ModeIMP, Handler: &XTAHandler{}},
		{Command: 0x6A, OpcodeName: "ROR", Mode: enums.ModeAcc, Handler: &RORHandler{}},
		{Command: 0x66, OpcodeName: "ROR", Mode: enums.ModeZP, Handler: &RORHandler{}},
		{Command: 0x76, OpcodeName: "ROR", Mode: enums.ModeZPX, Handler: &RORHandler{}},
		{Command: 0x6E, OpcodeName: "ROR", Mode: enums.ModeABS, Handler: &RORHandler{}},
		{Command: 0x7E, OpcodeName: "ROR", Mode: enums.ModeABSX, Handler: &RORHandler{}},
		{Command: 0xCA, OpcodeName: "DEX", Mode: enums.ModeIMP, Handler: &DEXHandler{}},
		{Command: 0x0A, OpcodeName: "ASL", Mode: enums.ModeAcc, Handler: &ASLHandler{}},
		{Command: 0x06, OpcodeName: "ASL", Mode: enums.ModeZP, Handler: &ASLHandler{}},
		{Command: 0x16, OpcodeName: "ASL", Mode: enums.ModeZPX, Handler: &ASLHandler{}},
		{Command: 0x0E, OpcodeName: "ASL", Mode: enums.ModeABS, Handler: &ASLHandler{}},
		{Command: 0x1E, OpcodeName: "ASL", Mode: enums.ModeABSX, Handler: &ASLHandler{}},
		{Command: 0x4A, OpcodeName: "LSR", Mode: enums.ModeAcc, Handler: &LSRHandler{}},
		{Command: 0x46, OpcodeName: "LSR", Mode: enums.ModeZP, Handler: &LSRHandler{}},
		{Command: 0x56, OpcodeName: "LSR", Mode: enums.ModeZPX, Handler: &LSRHandler{}},
		{Command: 0x4E, OpcodeName: "LSR", Mode: enums.ModeABS, Handler: &LSRHandler{}},
		{Command: 0x5E, OpcodeName: "LSR", Mode: enums.ModeABSX, Handler: &LSRHandler{}},
		{Command: 0x27, OpcodeName: "RLA", Mode: enums.ModeZP, Handler: &RLAHandler{}},
		{Command: 0x37, OpcodeName: "RLA", Mode: enums.ModeZPX, Handler: &RLAHandler{}},
		{Command: 0x2F, OpcodeName: "RLA", Mode: enums.ModeABS, Handler: &RLAHandler{}},
		{Command: 0x3F, OpcodeName: "RLA", Mode: enums.ModeABSX, Handler: &RLAHandler{}},
		{Command: 0x3B, OpcodeName: "RLA", Mode: enums.ModeABSY, Handler: &RLAHandler{}},
		{Command: 0x23, OpcodeName: "RLA", Mode: enums.ModeINDX, Handler: &RLAHandler{}},
		{Command: 0x33, OpcodeName: "RLA", Mode: enums.ModeINDY, Handler: &RLAHandler{}},
		{Command: 0x67, OpcodeName: "RRA", Mode: enums.ModeZP, Handler: &RRAHandler{}},
		{Command: 0x77, OpcodeName: "RRA", Mode: enums.ModeZPX, Handler: &RRAHandler{}},
		{Command: 0x6F, OpcodeName: "RRA", Mode: enums.ModeABS, Handler: &RRAHandler{}},
		{Command: 0x7F, OpcodeName: "RRA", Mode: enums.ModeABSX, Handler: &RRAHandler{}},
		{Command: 0x7B, OpcodeName: "RRA", Mode: enums.ModeABSY, Handler: &RRAHandler{}},
		{Command: 0x63, OpcodeName: "RRA", Mode: enums.ModeINDX, Handler: &RRAHandler{}},
		{Command: 0x73, OpcodeName: "RRA", Mode: enums.ModeINDY, Handler: &RRAHandler{}},
		{Command: 0xE7, OpcodeName: "ISB", Mode: enums.ModeZP, Handler: &ISBHandler{}},
		{Command: 0xF7, OpcodeName: "ISB", Mode: enums.ModeZPX, Handler: &ISBHandler{}},
		{Command: 0xEF, OpcodeName: "ISB", Mode: enums.ModeABS, Handler: &ISBHandler{}},
		{Command: 0xFF, OpcodeName: "ISB", Mode: enums.ModeABSX, Handler: &ISBHandler{}},
		{Command: 0xFB, OpcodeName: "ISB", Mode: enums.ModeABSY, Handler: &ISBHandler{}},
		{Command: 0xE3, OpcodeName: "ISB", Mode: enums.ModeINDX, Handler: &ISBHandler{}},
		{Command: 0xF3, OpcodeName: "ISB", Mode: enums.ModeINDY, Handler: &ISBHandler{}},
		{Command: 0xA2, OpcodeName: "LDX", Mode: enums.ModeIMM, Handler: &LDXHandler{}},
		{Command: 0xA6, OpcodeName: "LDX", Mode: enums.ModeZP, Handler: &LDXHandler{}},
		{Command: 0xB6, OpcodeName: "LDX", Mode: enums.ModeZPY, Handler: &LDXHandler{}},
		{Command: 0xAE, OpcodeName: "LDX", Mode: enums.ModeABS, Handler: &LDXHandler{}},
		{Command: 0xBE, OpcodeName: "LDX", Mode: enums.ModeABSY, Handler: &LDXHandler{}},
	}
}
