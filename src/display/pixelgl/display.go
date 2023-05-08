package pixelgl

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"main/src/bus"
	"main/src/config"
	"sync"
)

const (
	NesImageSizeX = 256
	NesImageSizeY = 240
)

type Display struct {
	config *config.Config

	win    *pixelgl.Window
	canvas *pixelgl.Canvas

	bus *bus.Bus

	coeff           int
	resizeLock      sync.Mutex
	resizeBuffer    []uint8
	pixelBuffer     []uint32
	currentWinSizeX float64
	currentWinSizeY float64

	buttonsJoyPad1 map[bus.JoyPadButton]pixelgl.Button

	ctx context.Context
}

func (d *Display) Init() {
	d.coeff = 1
	d.currentWinSizeX = NesImageSizeX
	d.currentWinSizeY = NesImageSizeY
	d.resizeBuffer = make([]uint8, (d.coeff*NesImageSizeX)*(d.coeff*NesImageSizeY)*4)
	d.pixelBuffer = make([]uint32, NesImageSizeX*NesImageSizeY)

	d.initInput()
}

func (d *Display) Run(ctx context.Context) {
	d.ctx = ctx
	pixelgl.Run(d.runInternal)
}

func (d *Display) RenderPixel(x, y uint16, color uint32) {
	y = NesImageSizeY - 1 - y
	d.pixelBuffer[y*NesImageSizeX+x] = color
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

	d.canvas = d.win.Canvas()

	for !d.win.Closed() {
		select {
		case <-d.ctx.Done():
			return
		default:
		}

		d.checkButtons()
		d.checkWindowResize()
		d.resizeFrame()

		d.canvas.SetPixels(d.resizeBuffer)
		d.canvas.Draw(d.win, pixel.IM.Moved(d.win.Bounds().Center()))

		d.win.Update()
	}
}

func (d *Display) resizeFrame() {
	var c uint32
	var offset int

	d.resizeLock.Lock()
	for y := 0; y < NesImageSizeY; y++ {
		for x := 0; x < NesImageSizeX; x++ {
			c = d.pixelBuffer[y*NesImageSizeX+x]

			for i := 0; i < d.coeff; i++ {
				for j := 0; j < d.coeff; j++ {
					offset = ((y*d.coeff)+i)*4*(NesImageSizeX*d.coeff) + ((x*d.coeff)+j)*4

					d.resizeBuffer[offset] = uint8(c)
					d.resizeBuffer[offset+1] = uint8(c >> 8)
					d.resizeBuffer[offset+2] = uint8(c >> 16)
					d.resizeBuffer[offset+3] = uint8(c >> 24)
				}
			}
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
			d.resizeBuffer = make([]uint8, (d.coeff*NesImageSizeX)*(d.coeff*NesImageSizeY)*4)
			d.resizeLock.Unlock()
		}
	}
}

func NewDisplay(
	config *config.Config,
	bus *bus.Bus,
) *Display {
	return &Display{
		config: config,
		bus:    bus,
	}
}
