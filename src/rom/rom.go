package rom

import (
	"main/src/enum"
	"os"
)

type Rom interface {
	Load(file *os.File) error
	GetData() []byte
	GetByte(address uint32) byte
	GetPrgRomSize() uint8
	GetChrRomSize() uint8
	GetMirroringType() enum.MirroringType
	GetMapperId() byte
}
