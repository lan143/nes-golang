package ppu

import (
	"main/src/display"
	"main/src/mapper"
)

type PPU interface {
	Init(mapper mapper.Mapper, display display.Display)
	Run()
}
