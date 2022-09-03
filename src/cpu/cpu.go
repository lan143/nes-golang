package cpu

import (
	"main/src/mapper"
)

type CPU interface {
	Init(mapper mapper.Mapper)
	Reset()
	Run()
}
