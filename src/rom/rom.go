package rom

import (
	"main/src/mapper/enum"
	"os"
)

type Rom interface {
	Load(file *os.File) error
	GetData() []byte
	GetPrgRomSize() uint8
	GetChrRomSize() uint8
	GetMirroringType() enum.MirroringType
	GetMapperId() byte
}
