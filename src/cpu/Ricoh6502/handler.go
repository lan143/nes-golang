package Ricoh6502

import (
	"main/src/cpu/Ricoh6502/enums"
)

type CommandHandler struct {
	Command    byte
	OpcodeName string
	Mode       enums.Modes
	Handler    Handler
	SkipCycles uint16
}

type Handler interface {
	Handle(cpu *Cpu, operand uint16, mode enums.Modes) error
}

var commandHandlers [256]*CommandHandler

func init() {
	commandHandlers[0xB0] = &CommandHandler{OpcodeName: "BCS", Mode: 13, Handler: &BCSHandler{}, SkipCycles: 2}
	commandHandlers[0xE4] = &CommandHandler{OpcodeName: "CPX", Mode: 2, Handler: &CPXHandler{}, SkipCycles: 3}
	commandHandlers[0xF7] = &CommandHandler{OpcodeName: "ISB", Mode: 3, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x0A] = &CommandHandler{OpcodeName: "ASL", Mode: 1, Handler: &ASLHandler{}, SkipCycles: 2}
	commandHandlers[0x35] = &CommandHandler{OpcodeName: "AND", Mode: 3, Handler: &ANDHandler{}, SkipCycles: 4}
	commandHandlers[0x49] = &CommandHandler{OpcodeName: "EOR", Mode: 8, Handler: &EORHandler{}, SkipCycles: 2}
	commandHandlers[0xAD] = &CommandHandler{OpcodeName: "LDA", Mode: 5, Handler: &LDAHandler{}, SkipCycles: 4}
	commandHandlers[0xFD] = &CommandHandler{OpcodeName: "SBC", Mode: 6, Handler: &SBCHandler{}, SkipCycles: 4}
	commandHandlers[0x68] = &CommandHandler{OpcodeName: "PLA", Mode: 11, Handler: &PLAHandler{}, SkipCycles: 4}
	commandHandlers[0x75] = &CommandHandler{OpcodeName: "ADC", Mode: 3, Handler: &ADCHandler{}, SkipCycles: 4}
	commandHandlers[0xD1] = &CommandHandler{OpcodeName: "CMP", Mode: 10, Handler: &CMPHandler{}, SkipCycles: 5}
	commandHandlers[0xDD] = &CommandHandler{OpcodeName: "CMP", Mode: 6, Handler: &CMPHandler{}, SkipCycles: 4}
	commandHandlers[0x73] = &CommandHandler{OpcodeName: "RRA", Mode: 10, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0xB5] = &CommandHandler{OpcodeName: "LDA", Mode: 3, Handler: &LDAHandler{}, SkipCycles: 4}
	commandHandlers[0x18] = &CommandHandler{OpcodeName: "CLC", Mode: 11, Handler: &CLCHandler{}, SkipCycles: 2}
	commandHandlers[0x2D] = &CommandHandler{OpcodeName: "AND", Mode: 5, Handler: &ANDHandler{}, SkipCycles: 4}
	commandHandlers[0x59] = &CommandHandler{OpcodeName: "EOR", Mode: 7, Handler: &EORHandler{}, SkipCycles: 4}
	commandHandlers[0x6A] = &CommandHandler{OpcodeName: "ROR", Mode: 1, Handler: &RORHandler{}, SkipCycles: 2}
	commandHandlers[0xED] = &CommandHandler{OpcodeName: "SBC", Mode: 5, Handler: &SBCHandler{}, SkipCycles: 4}
	commandHandlers[0xF0] = &CommandHandler{OpcodeName: "BEQ", Mode: 13, Handler: &BEQHandler{}, SkipCycles: 2}
	commandHandlers[0x53] = &CommandHandler{OpcodeName: "SRE", Mode: 10, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x7E] = &CommandHandler{OpcodeName: "ROR", Mode: 6, Handler: &RORHandler{}, SkipCycles: 7}
	commandHandlers[0x94] = &CommandHandler{OpcodeName: "STY", Mode: 3, Handler: &STYHandler{}, SkipCycles: 4}
	commandHandlers[0xD5] = &CommandHandler{OpcodeName: "CMP", Mode: 3, Handler: &CMPHandler{}, SkipCycles: 4}
	commandHandlers[0x08] = &CommandHandler{OpcodeName: "PHP", Mode: 11, Handler: &PHPHandler{}, SkipCycles: 3}
	commandHandlers[0x1E] = &CommandHandler{OpcodeName: "ASL", Mode: 6, Handler: &ASLHandler{}, SkipCycles: 7}
	commandHandlers[0x2A] = &CommandHandler{OpcodeName: "ROL", Mode: 1, Handler: &ROLHandler{}, SkipCycles: 2}
	commandHandlers[0x4D] = &CommandHandler{OpcodeName: "EOR", Mode: 5, Handler: &EORHandler{}, SkipCycles: 4}
	commandHandlers[0x25] = &CommandHandler{OpcodeName: "AND", Mode: 2, Handler: &ANDHandler{}, SkipCycles: 3}
	commandHandlers[0x79] = &CommandHandler{OpcodeName: "ADC", Mode: 7, Handler: &ADCHandler{}, SkipCycles: 4}
	commandHandlers[0x96] = &CommandHandler{OpcodeName: "STX", Mode: 4, Handler: &STXHandler{}, SkipCycles: 4}
	commandHandlers[0x36] = &CommandHandler{OpcodeName: "ROL", Mode: 3, Handler: &ROLHandler{}, SkipCycles: 6}
	commandHandlers[0x40] = &CommandHandler{OpcodeName: "RTI", Mode: 11, Handler: &RTIHandler{}, SkipCycles: 6}
	commandHandlers[0x4C] = &CommandHandler{OpcodeName: "JMP", Mode: 5, Handler: &JMPHandler{}, SkipCycles: 0}
	commandHandlers[0xEE] = &CommandHandler{OpcodeName: "INC", Mode: 5, Handler: &INCHandler{}, SkipCycles: 6}
	commandHandlers[0x8E] = &CommandHandler{OpcodeName: "STX", Mode: 5, Handler: &STXHandler{}, SkipCycles: 4}
	commandHandlers[0xB1] = &CommandHandler{OpcodeName: "LDA", Mode: 10, Handler: &LDAHandler{}, SkipCycles: 5}
	commandHandlers[0xD0] = &CommandHandler{OpcodeName: "BNE", Mode: 13, Handler: &BNEHandler{}, SkipCycles: 2}
	commandHandlers[0x05] = &CommandHandler{OpcodeName: "ORA", Mode: 2, Handler: &ORAHandler{}, SkipCycles: 3}
	commandHandlers[0x0E] = &CommandHandler{OpcodeName: "ASL", Mode: 5, Handler: &ASLHandler{}, SkipCycles: 6}
	commandHandlers[0x21] = &CommandHandler{OpcodeName: "AND", Mode: 9, Handler: &ANDHandler{}, SkipCycles: 6}
	commandHandlers[0x69] = &CommandHandler{OpcodeName: "ADC", Mode: 8, Handler: &ADCHandler{}, SkipCycles: 2}
	commandHandlers[0xC0] = &CommandHandler{OpcodeName: "CPY", Mode: 8, Handler: &CPYHandle{}, SkipCycles: 2}
	commandHandlers[0xEA] = &CommandHandler{OpcodeName: "NOP", Mode: 11, Handler: &NOPHandler{}, SkipCycles: 2}
	commandHandlers[0x15] = &CommandHandler{OpcodeName: "ORA", Mode: 3, Handler: &ORAHandler{}, SkipCycles: 4}
	commandHandlers[0x3F] = &CommandHandler{OpcodeName: "RLA", Mode: 6, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0x58] = &CommandHandler{OpcodeName: "CLI", Mode: 11, Handler: &CLIHandler{}, SkipCycles: 2}
	commandHandlers[0x77] = &CommandHandler{OpcodeName: "RRA", Mode: 3, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0x2C] = &CommandHandler{OpcodeName: "BIT", Mode: 5, Handler: &BITHandler{}, SkipCycles: 4}
	commandHandlers[0x38] = &CommandHandler{OpcodeName: "SEC", Mode: 11, Handler: &SECHandler{}, SkipCycles: 2}
	commandHandlers[0x8D] = &CommandHandler{OpcodeName: "STA", Mode: 5, Handler: &STAHandler{}, SkipCycles: 4}
	commandHandlers[0xE1] = &CommandHandler{OpcodeName: "SBC", Mode: 9, Handler: &SBCHandler{}, SkipCycles: 6}
	commandHandlers[0x16] = &CommandHandler{OpcodeName: "ASL", Mode: 3, Handler: &ASLHandler{}, SkipCycles: 6}
	commandHandlers[0xAC] = &CommandHandler{OpcodeName: "LDY", Mode: 5, Handler: &LDYHandler{}, SkipCycles: 4}
	commandHandlers[0xB9] = &CommandHandler{OpcodeName: "LDA", Mode: 7, Handler: &LDAHandler{}, SkipCycles: 4}
	commandHandlers[0xBE] = &CommandHandler{OpcodeName: "LDX", Mode: 7, Handler: &LDXHandler{}, SkipCycles: 4}
	commandHandlers[0xDE] = &CommandHandler{OpcodeName: "DEC", Mode: 6, Handler: &DECHandler{}, SkipCycles: 7}
	commandHandlers[0x6E] = &CommandHandler{OpcodeName: "ROR", Mode: 5, Handler: &RORHandler{}, SkipCycles: 6}
	commandHandlers[0x6F] = &CommandHandler{OpcodeName: "RRA", Mode: 5, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0xC8] = &CommandHandler{OpcodeName: "INY", Mode: 11, Handler: &INYHandler{}, SkipCycles: 2}
	commandHandlers[0xCA] = &CommandHandler{OpcodeName: "DEX", Mode: 11, Handler: &DEXHandler{}, SkipCycles: 2}
	commandHandlers[0xD9] = &CommandHandler{OpcodeName: "CMP", Mode: 7, Handler: &CMPHandler{}, SkipCycles: 4}
	commandHandlers[0xF8] = &CommandHandler{OpcodeName: "SED", Mode: 11, Handler: &SEDHandler{}, SkipCycles: 2}
	commandHandlers[0x13] = &CommandHandler{OpcodeName: "SLO", Mode: 10, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x1D] = &CommandHandler{OpcodeName: "ORA", Mode: 6, Handler: &ORAHandler{}, SkipCycles: 4}
	commandHandlers[0x2E] = &CommandHandler{OpcodeName: "ROL", Mode: 5, Handler: &ROLHandler{}, SkipCycles: 6}
	commandHandlers[0x60] = &CommandHandler{OpcodeName: "RTS", Mode: 11, Handler: &RTSHandler{}, SkipCycles: 6}
	commandHandlers[0xFB] = &CommandHandler{OpcodeName: "ISB", Mode: 7, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x29] = &CommandHandler{OpcodeName: "AND", Mode: 8, Handler: &ANDHandler{}, SkipCycles: 2}
	commandHandlers[0x8A] = &CommandHandler{OpcodeName: "TXA", Mode: 11, Handler: &TXAHandler{}, SkipCycles: 2}
	commandHandlers[0x91] = &CommandHandler{OpcodeName: "STA", Mode: 10, Handler: &STAHandler{}, SkipCycles: 6}
	commandHandlers[0xF5] = &CommandHandler{OpcodeName: "SBC", Mode: 3, Handler: &SBCHandler{}, SkipCycles: 4}
	commandHandlers[0x24] = &CommandHandler{OpcodeName: "BIT", Mode: 2, Handler: &BITHandler{}, SkipCycles: 3}
	commandHandlers[0x33] = &CommandHandler{OpcodeName: "RLA", Mode: 10, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0xB4] = &CommandHandler{OpcodeName: "LDY", Mode: 3, Handler: &LDYHandler{}, SkipCycles: 4}
	commandHandlers[0xB6] = &CommandHandler{OpcodeName: "LDX", Mode: 4, Handler: &LDXHandler{}, SkipCycles: 4}
	commandHandlers[0x63] = &CommandHandler{OpcodeName: "RRA", Mode: 9, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0x67] = &CommandHandler{OpcodeName: "RRA", Mode: 2, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0xBC] = &CommandHandler{OpcodeName: "LDY", Mode: 6, Handler: &LDYHandler{}, SkipCycles: 4}
	commandHandlers[0xE7] = &CommandHandler{OpcodeName: "ISB", Mode: 2, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x28] = &CommandHandler{OpcodeName: "PLP", Mode: 11, Handler: &PLPHandler{}, SkipCycles: 4}
	commandHandlers[0x41] = &CommandHandler{OpcodeName: "EOR", Mode: 9, Handler: &EORHandler{}, SkipCycles: 6}
	commandHandlers[0x5B] = &CommandHandler{OpcodeName: "SRE", Mode: 7, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x61] = &CommandHandler{OpcodeName: "ADC", Mode: 9, Handler: &ADCHandler{}, SkipCycles: 6}
	commandHandlers[0xEF] = &CommandHandler{OpcodeName: "ISB", Mode: 5, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x6D] = &CommandHandler{OpcodeName: "ADC", Mode: 5, Handler: &ADCHandler{}, SkipCycles: 4}
	commandHandlers[0x84] = &CommandHandler{OpcodeName: "STY", Mode: 2, Handler: &STYHandler{}, SkipCycles: 3}
	commandHandlers[0xA8] = &CommandHandler{OpcodeName: "TAY", Mode: 11, Handler: &TAYHandler{}, SkipCycles: 2}
	commandHandlers[0xE0] = &CommandHandler{OpcodeName: "CPX", Mode: 8, Handler: &CPXHandler{}, SkipCycles: 2}
	commandHandlers[0x78] = &CommandHandler{OpcodeName: "SEI", Mode: 11, Handler: &SEIHandler{}, SkipCycles: 2}
	commandHandlers[0xA9] = &CommandHandler{OpcodeName: "LDA", Mode: 8, Handler: &LDAHandler{}, SkipCycles: 2}
	commandHandlers[0xEC] = &CommandHandler{OpcodeName: "CPX", Mode: 5, Handler: &CPXHandler{}, SkipCycles: 4}
	commandHandlers[0xF1] = &CommandHandler{OpcodeName: "SBC", Mode: 10, Handler: &SBCHandler{}, SkipCycles: 5}
	commandHandlers[0x0D] = &CommandHandler{OpcodeName: "ORA", Mode: 5, Handler: &ORAHandler{}, SkipCycles: 4}
	commandHandlers[0x4F] = &CommandHandler{OpcodeName: "SRE", Mode: 5, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x5E] = &CommandHandler{OpcodeName: "LSR", Mode: 6, Handler: &LSRHandler{}, SkipCycles: 7}
	commandHandlers[0x76] = &CommandHandler{OpcodeName: "ROR", Mode: 3, Handler: &RORHandler{}, SkipCycles: 6}
	commandHandlers[0xAA] = &CommandHandler{OpcodeName: "TAX", Mode: 11, Handler: &TAXHandler{}, SkipCycles: 2}
	commandHandlers[0x10] = &CommandHandler{OpcodeName: "BPL", Mode: 13, Handler: &BPLHandler{}, SkipCycles: 2}
	commandHandlers[0x2F] = &CommandHandler{OpcodeName: "RLA", Mode: 5, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0x3E] = &CommandHandler{OpcodeName: "ROL", Mode: 6, Handler: &ROLHandler{}, SkipCycles: 7}
	commandHandlers[0x46] = &CommandHandler{OpcodeName: "LSR", Mode: 2, Handler: &LSRHandler{}, SkipCycles: 5}
	commandHandlers[0xFF] = &CommandHandler{OpcodeName: "ISB", Mode: 6, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x03] = &CommandHandler{OpcodeName: "SLO", Mode: 9, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x39] = &CommandHandler{OpcodeName: "AND", Mode: 7, Handler: &ANDHandler{}, SkipCycles: 4}
	commandHandlers[0x7F] = &CommandHandler{OpcodeName: "RRA", Mode: 6, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0xCC] = &CommandHandler{OpcodeName: "CPY", Mode: 5, Handler: &CPYHandle{}, SkipCycles: 4}
	commandHandlers[0xE3] = &CommandHandler{OpcodeName: "ISB", Mode: 9, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0xE8] = &CommandHandler{OpcodeName: "INX", Mode: 11, Handler: &INXHandler{}, SkipCycles: 2}
	commandHandlers[0x5F] = &CommandHandler{OpcodeName: "SRE", Mode: 6, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x81] = &CommandHandler{OpcodeName: "STA", Mode: 9, Handler: &STAHandler{}, SkipCycles: 6}
	commandHandlers[0x9A] = &CommandHandler{OpcodeName: "TXS", Mode: 11, Handler: &TXSHandler{}, SkipCycles: 2}
	commandHandlers[0xC6] = &CommandHandler{OpcodeName: "DEC", Mode: 2, Handler: &DECHandler{}, SkipCycles: 5}
	commandHandlers[0x86] = &CommandHandler{OpcodeName: "STX", Mode: 2, Handler: &STXHandler{}, SkipCycles: 3}
	commandHandlers[0xC5] = &CommandHandler{OpcodeName: "CMP", Mode: 2, Handler: &CMPHandler{}, SkipCycles: 3}
	commandHandlers[0x0F] = &CommandHandler{OpcodeName: "SLO", Mode: 5, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x1B] = &CommandHandler{OpcodeName: "SLO", Mode: 7, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x4E] = &CommandHandler{OpcodeName: "LSR", Mode: 5, Handler: &LSRHandler{}, SkipCycles: 6}
	commandHandlers[0x85] = &CommandHandler{OpcodeName: "STA", Mode: 2, Handler: &STAHandler{}, SkipCycles: 3}
	commandHandlers[0x71] = &CommandHandler{OpcodeName: "ADC", Mode: 10, Handler: &ADCHandler{}, SkipCycles: 5}
	commandHandlers[0xA0] = &CommandHandler{OpcodeName: "LDY", Mode: 8, Handler: &LDYHandler{}, SkipCycles: 2}
	commandHandlers[0xA6] = &CommandHandler{OpcodeName: "LDX", Mode: 2, Handler: &LDXHandler{}, SkipCycles: 3}
	commandHandlers[0xF9] = &CommandHandler{OpcodeName: "SBC", Mode: 7, Handler: &SBCHandler{}, SkipCycles: 4}
	commandHandlers[0x01] = &CommandHandler{OpcodeName: "ORA", Mode: 9, Handler: &ORAHandler{}, SkipCycles: 6}
	commandHandlers[0x45] = &CommandHandler{OpcodeName: "EOR", Mode: 2, Handler: &EORHandler{}, SkipCycles: 3}
	commandHandlers[0x51] = &CommandHandler{OpcodeName: "EOR", Mode: 10, Handler: &EORHandler{}, SkipCycles: 5}
	commandHandlers[0x65] = &CommandHandler{OpcodeName: "ADC", Mode: 2, Handler: &ADCHandler{}, SkipCycles: 3}
	commandHandlers[0xA4] = &CommandHandler{OpcodeName: "LDY", Mode: 2, Handler: &LDYHandler{}, SkipCycles: 3}
	commandHandlers[0x00] = &CommandHandler{OpcodeName: "BRK", Mode: 11, Handler: &BRKHandler{}, SkipCycles: 7}
	commandHandlers[0x23] = &CommandHandler{OpcodeName: "RLA", Mode: 9, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0x7D] = &CommandHandler{OpcodeName: "ADC", Mode: 6, Handler: &ADCHandler{}, SkipCycles: 4}
	commandHandlers[0x88] = &CommandHandler{OpcodeName: "DEY", Mode: 11, Handler: &DEYHandler{}, SkipCycles: 2}
	commandHandlers[0xC1] = &CommandHandler{OpcodeName: "CMP", Mode: 9, Handler: &CMPHandler{}, SkipCycles: 6}
	commandHandlers[0xC9] = &CommandHandler{OpcodeName: "CMP", Mode: 8, Handler: &CMPHandler{}, SkipCycles: 2}
	commandHandlers[0xD8] = &CommandHandler{OpcodeName: "CLD", Mode: 11, Handler: &CLDHandler{}, SkipCycles: 2}
	commandHandlers[0xF6] = &CommandHandler{OpcodeName: "INC", Mode: 3, Handler: &INCHandler{}, SkipCycles: 6}
	commandHandlers[0x06] = &CommandHandler{OpcodeName: "ASL", Mode: 2, Handler: &ASLHandler{}, SkipCycles: 5}
	commandHandlers[0x3B] = &CommandHandler{OpcodeName: "RLA", Mode: 7, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0x5D] = &CommandHandler{OpcodeName: "EOR", Mode: 6, Handler: &EORHandler{}, SkipCycles: 4}
	commandHandlers[0x98] = &CommandHandler{OpcodeName: "TYA", Mode: 11, Handler: &TYAHandler{}, SkipCycles: 2}
	commandHandlers[0x7B] = &CommandHandler{OpcodeName: "RRA", Mode: 7, Handler: &RRAHandler{}, SkipCycles: 0}
	commandHandlers[0xC4] = &CommandHandler{OpcodeName: "CPY", Mode: 2, Handler: &CPYHandle{}, SkipCycles: 3}
	commandHandlers[0xF3] = &CommandHandler{OpcodeName: "ISB", Mode: 10, Handler: &ISBHandler{}, SkipCycles: 0}
	commandHandlers[0x09] = &CommandHandler{OpcodeName: "ORA", Mode: 8, Handler: &ORAHandler{}, SkipCycles: 2}
	commandHandlers[0x55] = &CommandHandler{OpcodeName: "EOR", Mode: 3, Handler: &EORHandler{}, SkipCycles: 4}
	commandHandlers[0x57] = &CommandHandler{OpcodeName: "SRE", Mode: 3, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x6C] = &CommandHandler{OpcodeName: "JMP", Mode: 12, Handler: &JMPHandler{}, SkipCycles: 0}
	commandHandlers[0x95] = &CommandHandler{OpcodeName: "STA", Mode: 3, Handler: &STAHandler{}, SkipCycles: 4}
	commandHandlers[0xA1] = &CommandHandler{OpcodeName: "LDA", Mode: 9, Handler: &LDAHandler{}, SkipCycles: 6}
	commandHandlers[0xB8] = &CommandHandler{OpcodeName: "CLV", Mode: 11, Handler: &CLVHandler{}, SkipCycles: 2}
	commandHandlers[0xCD] = &CommandHandler{OpcodeName: "CMP", Mode: 5, Handler: &CMPHandler{}, SkipCycles: 4}
	commandHandlers[0x1F] = &CommandHandler{OpcodeName: "SLO", Mode: 6, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x27] = &CommandHandler{OpcodeName: "RLA", Mode: 2, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0x30] = &CommandHandler{OpcodeName: "BMI", Mode: 13, Handler: &BMIHandler{}, SkipCycles: 2}
	commandHandlers[0x47] = &CommandHandler{OpcodeName: "SRE", Mode: 2, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x19] = &CommandHandler{OpcodeName: "ORA", Mode: 7, Handler: &ORAHandler{}, SkipCycles: 4}
	commandHandlers[0x31] = &CommandHandler{OpcodeName: "AND", Mode: 10, Handler: &ANDHandler{}, SkipCycles: 5}
	commandHandlers[0x37] = &CommandHandler{OpcodeName: "RLA", Mode: 3, Handler: &RLAHandler{}, SkipCycles: 0}
	commandHandlers[0xD6] = &CommandHandler{OpcodeName: "DEC", Mode: 3, Handler: &DECHandler{}, SkipCycles: 6}
	commandHandlers[0x9D] = &CommandHandler{OpcodeName: "STA", Mode: 6, Handler: &STAHandler{}, SkipCycles: 5}
	commandHandlers[0xE5] = &CommandHandler{OpcodeName: "SBC", Mode: 2, Handler: &SBCHandler{}, SkipCycles: 3}
	commandHandlers[0x11] = &CommandHandler{OpcodeName: "ORA", Mode: 10, Handler: &ORAHandler{}, SkipCycles: 5}
	commandHandlers[0x17] = &CommandHandler{OpcodeName: "SLO", Mode: 3, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x50] = &CommandHandler{OpcodeName: "BVC", Mode: 13, Handler: &BVCHandler{}, SkipCycles: 2}
	commandHandlers[0x99] = &CommandHandler{OpcodeName: "STA", Mode: 7, Handler: &STAHandler{}, SkipCycles: 5}
	commandHandlers[0x70] = &CommandHandler{OpcodeName: "BVS", Mode: 13, Handler: &BVSHandler{}, SkipCycles: 2}
	commandHandlers[0x8C] = &CommandHandler{OpcodeName: "STY", Mode: 5, Handler: &STYHandler{}, SkipCycles: 4}
	commandHandlers[0xBA] = &CommandHandler{OpcodeName: "TSX", Mode: 11, Handler: &TSXHandler{}, SkipCycles: 2}
	commandHandlers[0xBD] = &CommandHandler{OpcodeName: "LDA", Mode: 6, Handler: &LDAHandler{}, SkipCycles: 4}
	commandHandlers[0x07] = &CommandHandler{OpcodeName: "SLO", Mode: 2, Handler: &SLOHandler{}, SkipCycles: 0}
	commandHandlers[0x20] = &CommandHandler{OpcodeName: "JSR", Mode: 5, Handler: &JSRHandler{}, SkipCycles: 0}
	commandHandlers[0x48] = &CommandHandler{OpcodeName: "PHA", Mode: 11, Handler: &PHAHandler{}, SkipCycles: 3}
	commandHandlers[0x56] = &CommandHandler{OpcodeName: "LSR", Mode: 3, Handler: &LSRHandler{}, SkipCycles: 6}
	commandHandlers[0xCE] = &CommandHandler{OpcodeName: "DEC", Mode: 5, Handler: &DECHandler{}, SkipCycles: 6}
	commandHandlers[0xE9] = &CommandHandler{OpcodeName: "SBC", Mode: 8, Handler: &SBCHandler{}, SkipCycles: 2}
	commandHandlers[0xFE] = &CommandHandler{OpcodeName: "INC", Mode: 6, Handler: &INCHandler{}, SkipCycles: 7}
	commandHandlers[0x26] = &CommandHandler{OpcodeName: "ROL", Mode: 2, Handler: &ROLHandler{}, SkipCycles: 5}
	commandHandlers[0x3D] = &CommandHandler{OpcodeName: "AND", Mode: 6, Handler: &ANDHandler{}, SkipCycles: 4}
	commandHandlers[0x4A] = &CommandHandler{OpcodeName: "LSR", Mode: 1, Handler: &LSRHandler{}, SkipCycles: 2}
	commandHandlers[0xAE] = &CommandHandler{OpcodeName: "LDX", Mode: 5, Handler: &LDXHandler{}, SkipCycles: 4}
	commandHandlers[0xA5] = &CommandHandler{OpcodeName: "LDA", Mode: 2, Handler: &LDAHandler{}, SkipCycles: 3}
	commandHandlers[0xE6] = &CommandHandler{OpcodeName: "INC", Mode: 2, Handler: &INCHandler{}, SkipCycles: 5}
	commandHandlers[0x43] = &CommandHandler{OpcodeName: "SRE", Mode: 9, Handler: &SREHandler{}, SkipCycles: 0}
	commandHandlers[0x66] = &CommandHandler{OpcodeName: "ROR", Mode: 2, Handler: &RORHandler{}, SkipCycles: 5}
	commandHandlers[0x90] = &CommandHandler{OpcodeName: "BCC", Mode: 13, Handler: &BCCHandler{}, SkipCycles: 2}
	commandHandlers[0xA2] = &CommandHandler{OpcodeName: "LDX", Mode: 8, Handler: &LDXHandler{}, SkipCycles: 2}
}
