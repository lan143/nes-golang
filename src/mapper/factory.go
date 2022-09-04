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
	default:
		return nil, fmt.Errorf("unsupported mapper %d", id)
	}
}

func NewFactory() *Factory {
	return &Factory{}
}
