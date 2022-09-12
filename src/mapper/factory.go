package mapper

import (
	"fmt"
	"main/src/mapper/enum"
)

type Factory struct {
}

func (f *Factory) GetMapper(id enum.Id) (Mapper, error) {
	switch id {
	case enum.NROM:
		return &NROMMapper{}, nil
	case enum.UnROM:
		return &UnROMMapper{}, nil
	default:
		return nil, fmt.Errorf("unsupported mapper %d", id)
	}
}

func NewFactory() *Factory {
	return &Factory{}
}
