package cpu

import (
	"main/src/bus"
	"main/src/mapper"
)

type CPU interface {
	Init(mapper mapper.Mapper)
	Reset()
	Run()
	OnEvent(event bus.Event)
}
