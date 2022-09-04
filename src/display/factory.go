package display

import "main/src/display/pixelgl"

type Factory struct {
}

func (f *Factory) GetDisplay() Display {
	return pixelgl.NewDisplay()
}

func NewFactory() *Factory {
	return &Factory{}
}
