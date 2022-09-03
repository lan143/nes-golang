package mapper

import "fmt"

type Factory struct {
}

func (f *Factory) GetMapper(id Id) (Mapper, error) {
	switch id {
	case NROM:
		return &NROMMapper{}, nil
	default:
		return nil, fmt.Errorf("unsupported mapper %d", id)
	}
}

func NewFactory() *Factory {
	return &Factory{}
}
