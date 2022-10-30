package pixelgl

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"main/src/bus"
	"sync"
)

const (
	NesImageSizeX = 256
	NesImageSizeY = 240
)

type Display struct {
	win            *pixelgl.Window
	canvas         *pixelgl.Canvas
	updateScreenCh chan any

	bus *bus.Bus

	coeff           int
	resizeLock      sync.Mutex
	pixelBuffer     []uint8
	currentWinSizeX float64
	currentWinSizeY float64
}

func (d *Display) Init() {
	d.coeff = 1
	d.currentWinSizeX = NesImageSizeX
	d.currentWinSizeY = NesImageSizeY
	d.pixelBuffer = make([]uint8, (d.coeff*NesImageSizeX)*(d.coeff*NesImageSizeY)*4)

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
		Title:     "NES EMU",
		Bounds:    pixel.R(0, 0, d.currentWinSizeX, d.currentWinSizeY),
		VSync:     true,
		Resizable: true,
	}

	var err error
	d.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	d.canvas = pixelgl.NewCanvas(d.win.Bounds())

	for {
		<-d.updateScreenCh
		d.checkButtons()
		d.checkWindowResize()

		if !d.win.Closed() {
			d.canvas.SetPixels(d.pixelBuffer)
			d.canvas.Draw(d.win, pixel.IM.Moved(d.win.Bounds().Center()))
			d.win.Update()
		}
	}
}

func (d *Display) RenderPixel(x, y int, color uint32) {
	y = NesImageSizeY - 1 - y

	d.resizeLock.Lock()

	for i := 0; i < d.coeff; i++ {
		for j := 0; j < d.coeff; j++ {
			offset := ((y*d.coeff)+i)*4*(NesImageSizeX*d.coeff) + ((x*d.coeff)+j)*4

			d.pixelBuffer[offset] = uint8(color)
			d.pixelBuffer[offset+1] = uint8(color >> 8)
			d.pixelBuffer[offset+2] = uint8(color >> 16)
			d.pixelBuffer[offset+3] = uint8(color >> 24)
		}
	}

	d.resizeLock.Unlock()
}

func (d *Display) checkWindowResize() {
	bounds := d.win.Bounds()

	currentWinSizeX := bounds.Max.X - bounds.Min.X
	currentWinSizeY := bounds.Max.Y - bounds.Min.Y

	if currentWinSizeX != d.currentWinSizeX || currentWinSizeY != d.currentWinSizeY {
		d.currentWinSizeX = currentWinSizeX
		d.currentWinSizeY = currentWinSizeY

		coeffX := currentWinSizeX / NesImageSizeX
		coeffY := currentWinSizeX / NesImageSizeY
		coeff := 0

		newSizeX := NesImageSizeX * coeffX
		newSizeY := NesImageSizeY * coeffY

		if newSizeX > d.currentWinSizeX {
			coeff = int(coeffY)
		}

		if newSizeY > d.currentWinSizeY {
			coeff = int(coeffX)
		}

		if coeff != 1 && coeff%2 != 0 {
			coeff = coeff - 1
		}

		if coeff != d.coeff {
			d.resizeLock.Lock()
			d.coeff = coeff
			d.canvas = pixelgl.NewCanvas(pixel.Rect{Min: pixel.Vec{
				X: 0,
				Y: 0,
			}, Max: pixel.Vec{
				X: float64(d.coeff * NesImageSizeX),
				Y: float64(d.coeff * NesImageSizeY),
			}})
			d.win.Clear(color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			})
			d.pixelBuffer = make([]uint8, (d.coeff*NesImageSizeX)*(d.coeff*NesImageSizeY)*4)
			d.resizeLock.Unlock()
		}
	}
}

func (d *Display) checkButtons() {
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

func NewDisplay(bus *bus.Bus) *Display {
	return &Display{
		bus: bus,
	}
}
