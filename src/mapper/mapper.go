package mapper

import (
	"main/src/enum"
)

type Mapper interface {
	Init(prgRomSize byte) error
	GetMirroringType() enum.MirroringType
	MapPrgRom(address uint16) uint32
	MapChrRom(address uint16) uint32
	PutByte(address uint16, value byte)
}
