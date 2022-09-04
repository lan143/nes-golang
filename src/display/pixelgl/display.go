package pixelgl

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	color2 "image/color"
)

type Display struct {
	win            *pixelgl.Window
	canvas         *pixelgl.Canvas
	buffer         *image.RGBA
	updateScreenCh chan any
}

func (d *Display) Init() {
	d.buffer = image.NewRGBA(image.Rect(0, 0, 256, 240))
	d.updateScreenCh = make(chan any)
}

func (d *Display) Run() {
	pixelgl.Run(d.runInternal)
}

func (d *Display) UpdateScreen() {
	d.updateScreenCh <- struct{}{}
}

func (d *Display) runInternal() {
	cfg := pixelgl.WindowConfig{
		Title:  "NES EMU",
		Bounds: pixel.R(0, 0, 256, 240),
	}

	var err error
	d.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	d.canvas = pixelgl.NewCanvas(d.win.Bounds())

	for {
		<-d.updateScreenCh

		if !d.win.Closed() {
			d.canvas.SetPixels(d.buffer.Pix)
			d.canvas.Draw(d.win, pixel.IM.Moved(d.win.Bounds().Center()))
			d.win.Update()
		}
	}
}

func (d *Display) RenderPixel(x, y uint16, color uint32) {
	d.buffer.Set(int(x), int(y), color2.RGBA{
		R: uint8(color >> 16),
		G: uint8(color >> 8),
		B: uint8(color),
		A: uint8(color >> 24),
	})
}

func NewDisplay() *Display {
	return &Display{}
}
