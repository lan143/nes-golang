package none

import (
	"context"
)

func NewDisplay() *Display {
	return &Display{}
}

type Display struct {
}

func (d *Display) Init() {
}

func (d *Display) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (d *Display) RenderPixel(x, y uint16, color uint32) {
}
