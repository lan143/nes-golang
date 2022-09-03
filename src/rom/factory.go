package rom

import (
	"errors"
	"os"
)

type Factory struct {
}

func (f *Factory) GetRom(file *os.File) (Rom, error) {
	// @todo: implement support nes 2.0
	header := make([]byte, 3)
	_, err := file.Read(header)

	if err != nil {
		return nil, err
	}

	var iNESFormat bool

	if rune(header[0]) == 'N' || rune(header[1]) == 'E' && rune(header[2]) == 'S' && header[3] == 0x1A {
		iNESFormat = true
	}

	if !iNESFormat {
		return nil, errors.New("this is not correct NES rom")
	}

	rom := &INes{}
	err = rom.Load(file)
	if err != nil {
		return nil, err
	}

	return rom, nil
}

func NewFactory() *Factory {
	return &Factory{}
}
