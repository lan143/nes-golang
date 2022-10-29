package rom

import (
	"log"
	"main/src/enum"
	"os"
)

type INes struct {
	prgRomSize uint8 // Size of PRG ROM in 16 KB units
	chrRomSize uint8 // Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
	flags6     uint8 // Mapper, mirroring, battery, trainer
	flags7     uint8 // Mapper, VS/Playchoice, NES 2.0
	flags8     uint8 // PRG-RAM size (rarely used extension)
	flags9     uint8 // TV system (rarely used extension)
	flags10    uint8 // TV system, PRG-RAM presence (unofficial, rarely used extension)

	data []byte
}

func (n *INes) GetPrgRomSize() uint8 {
	return n.prgRomSize
}

func (n *INes) GetChrRomSize() uint8 {
	return n.chrRomSize
}

func (n *INes) GetMirroringType() enum.MirroringType {
	if n.flags6&0x1 > 0 {
		return enum.Vertical
	} else {
		return enum.Horizontal
	}
}

func (n *INes) Load(file *os.File) error {
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	log.Printf("ROM size: %d", fi.Size())

	n.data = make([]byte, fi.Size()-0x10)

	err = n.loadHeader(file)
	if err != nil {
		return err
	}

	_, err = file.Seek(0x10, 0)

	size, err := file.Read(n.data)
	if err != nil {
		return err
	}

	log.Printf("Read %d bytes of data", size)

	return nil
}

func (n *INes) GetData() []byte {
	return n.data
}

func (n *INes) GetByte(address uint32) byte {
	return n.data[address]
}

func (n *INes) GetMapperId() byte {
	return (n.flags6 >> 4) | ((n.flags7 >> 4) << 4)
}

func (n *INes) loadHeader(file *os.File) error {
	header := make([]byte, 0x10)

	file.Seek(0, 0)
	_, err := file.Read(header)
	if err != nil {
		return err
	}

	n.prgRomSize = header[4]
	n.chrRomSize = header[5]
	n.flags6 = header[6] // control byte 1
	n.flags7 = header[7] // control byte 2
	n.flags8 = header[8]
	n.flags9 = header[9]
	n.flags10 = header[10]

	return nil
}
