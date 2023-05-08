package ppu

import (
	"main/src/cartridge"
	"main/src/display"
	"main/src/ram"
)

type PPU interface {
	Init(cartridge *cartridge.Cartridge, display display.Display, cpuRam *ram.Ram)
	RunCycle()
}
