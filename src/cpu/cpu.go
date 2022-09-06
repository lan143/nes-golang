package cpu

import (
	"main/src/mapper"
	"main/src/ram"
)

type CPU interface {
	Init(mapper mapper.Mapper, ram *ram.Ram)
	Reset()
	Run()
}
