package cpu

import (
	"main/src/cartridge"
	"main/src/ram"
)

type CPU interface {
	Init(cartridge *cartridge.Cartridge, ram *ram.Ram)
	Reset()
	Run()
}
