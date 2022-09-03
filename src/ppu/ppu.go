package ppu

import (
	"main/src/mapper"
)

type PPU interface {
	Init(mapper mapper.Mapper)
	Run()
}
