package display

import "context"

type Display interface {
	Init()
	Run(ctx context.Context)
	RenderPixel(x, y uint16, color uint32)
}
