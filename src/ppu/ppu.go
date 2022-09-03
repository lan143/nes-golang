package ppu

import (
	"main/src/bus"
	"main/src/mapper"
)

type PPU interface {
	Init(mapper mapper.Mapper, bus *bus.Bus)
	Run()
}
