package mapper

import (
	"main/src/mapper/enum"
	"main/src/rom"
)

type Mapper interface {
	LoadRom(rom rom.Rom)
	GetByte(address uint16) byte
	GetChrByte(address uint16) byte
	PutByte(address uint16, value byte)
	PutChrByte(address uint16, value byte)
	HasChrRom() bool
	GetMirroringType() enum.MirroringType
}
