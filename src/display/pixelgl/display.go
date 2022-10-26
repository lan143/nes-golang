package pixelgl

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"main/src/bus"
)

type Display struct {
	win            *pixelgl.Window
	canvas         *pixelgl.Canvas
	pixelBuffer    []uint8
	updateScreenCh chan any

	bus *bus.Bus
}

func (d *Display) Init() {
	d.pixelBuffer = make([]uint8, 512*480*4)
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
		Title: "NES EMU",
		//Bounds: pixel.R(0, 0, 256, 240),
		Bounds: pixel.R(0, 0, 512, 480),
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
			d.canvas.SetPixels(d.pixelBuffer)
			d.canvas.Draw(d.win, pixel.IM.Moved(d.win.Bounds().Center()))
			d.win.Update()

			if d.win.Pressed(pixelgl.KeyUp) {
				d.bus.KeyEvent(bus.JoyPadButtonUp, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonUp, false)
			}

			if d.win.Pressed(pixelgl.KeyDown) {
				d.bus.KeyEvent(bus.JoyPadButtonDown, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonDown, false)
			}

			if d.win.Pressed(pixelgl.KeyLeft) {
				d.bus.KeyEvent(bus.JoyPadButtonLeft, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonLeft, false)
			}

			if d.win.Pressed(pixelgl.KeyRight) {
				d.bus.KeyEvent(bus.JoyPadButtonRight, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonRight, false)
			}

			if d.win.Pressed(pixelgl.KeyZ) {
				d.bus.KeyEvent(bus.JoyPadButtonA, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonA, false)
			}

			if d.win.Pressed(pixelgl.KeyX) {
				d.bus.KeyEvent(bus.JoyPadButtonB, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonB, false)
			}

			if d.win.Pressed(pixelgl.KeyB) {
				d.bus.KeyEvent(bus.JoyPadButtonSelect, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonSelect, false)
			}

			if d.win.Pressed(pixelgl.KeyEnter) {
				d.bus.KeyEvent(bus.JoyPadButtonStart, true)
			} else {
				d.bus.KeyEvent(bus.JoyPadButtonStart, false)
			}
		}
	}
}

func (d *Display) RenderPixel(x, y uint16, color uint32) {
	y = 239 - y

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			offset := ((int(y)*2)+i)*4*512 + ((int(x)*2)+j)*4

			d.pixelBuffer[offset] = uint8(color)
			d.pixelBuffer[offset+1] = uint8(color >> 8)
			d.pixelBuffer[offset+2] = uint8(color >> 16)
			d.pixelBuffer[offset+3] = uint8(color >> 24)
		}
	}
}

func NewDisplay(bus *bus.Bus) *Display {
	return &Display{
		bus: bus,
	}
}
