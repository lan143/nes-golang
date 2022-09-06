package ppu

import (
	"main/src/display"
	"main/src/mapper"
	"main/src/ram"
)

type PPU interface {
	Init(mapper mapper.Mapper, display display.Display, cpuRam *ram.Ram)
	Run()
}
