package mapper

import (
	"fmt"
	"main/src/enum"
	ines_mapper_3 "main/src/mapper/ines-mapper-3"
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
	case enum.MapperINES003:
		return &ines_mapper_3.Mapper{}, nil
	default:
		return nil, fmt.Errorf("unsupported mapper (%d) %s", id, id.String())
	}
}

func NewFactory() *Factory {
	return &Factory{}
}
