package mapper

import (
	"fmt"
	"main/src/enum"
	"main/src/mapper/mmc1"
	"main/src/mapper/nrom"
	"main/src/mapper/unrom"
)

type Factory struct {
}

func (f *Factory) GetMapper(id enum.MapperId) (Mapper, error) {
	switch id {
	case enum.MapperNROM:
		return &nrom.Mapper{}, nil
	case enum.MapperMMC1:
		return &mmc1.Mapper{}, nil
	case enum.MapperUnROM:
		return &unrom.Mapper{}, nil
	default:
		return nil, fmt.Errorf("unsupported mapper %d", id)
	}
}

func NewFactory() *Factory {
	return &Factory{}
}
