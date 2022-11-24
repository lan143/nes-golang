package pixelgl

import (
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

var buttonNames = map[pixelgl.Button]string{
	pixelgl.MouseButton4:      "MouseButton4",
	pixelgl.MouseButton5:      "MouseButton5",
	pixelgl.MouseButton6:      "MouseButton6",
	pixelgl.MouseButton7:      "MouseButton7",
	pixelgl.MouseButton8:      "MouseButton8",
	pixelgl.MouseButtonLeft:   "MouseButtonLeft",
	pixelgl.MouseButtonRight:  "MouseButtonRight",
	pixelgl.MouseButtonMiddle: "MouseButtonMiddle",
	pixelgl.KeyUnknown:        "Unknown",
	pixelgl.KeySpace:          "Space",
	pixelgl.KeyApostrophe:     "Apostrophe",
	pixelgl.KeyComma:          "Comma",
	pixelgl.KeyMinus:          "Minus",
	pixelgl.KeyPeriod:         "Period",
	pixelgl.KeySlash:          "Slash",
	pixelgl.Key0:              "0",
	pixelgl.Key1:              "1",
	pixelgl.Key2:              "2",
	pixelgl.Key3:              "3",
	pixelgl.Key4:              "4",
	pixelgl.Key5:              "5",
	pixelgl.Key6:              "6",
	pixelgl.Key7:              "7",
	pixelgl.Key8:              "8",
	pixelgl.Key9:              "9",
	pixelgl.KeySemicolon:      "Semicolon",
	pixelgl.KeyEqual:          "Equal",
	pixelgl.KeyA:              "A",
	pixelgl.KeyB:              "B",
	pixelgl.KeyC:              "C",
	pixelgl.KeyD:              "D",
	pixelgl.KeyE:              "E",
	pixelgl.KeyF:              "F",
	pixelgl.KeyG:              "G",
	pixelgl.KeyH:              "H",
	pixelgl.KeyI:              "I",
	pixelgl.KeyJ:              "J",
	pixelgl.KeyK:              "K",
	pixelgl.KeyL:              "L",
	pixelgl.KeyM:              "M",
	pixelgl.KeyN:              "N",
	pixelgl.KeyO:              "O",
	pixelgl.KeyP:              "P",
	pixelgl.KeyQ:              "Q",
	pixelgl.KeyR:              "R",
	pixelgl.KeyS:              "S",
	pixelgl.KeyT:              "T",
	pixelgl.KeyU:              "U",
	pixelgl.KeyV:              "V",
	pixelgl.KeyW:              "W",
	pixelgl.KeyX:              "X",
	pixelgl.KeyY:              "Y",
	pixelgl.KeyZ:              "Z",
	pixelgl.KeyLeftBracket:    "LeftBracket",
	pixelgl.KeyBackslash:      "Backslash",
	pixelgl.KeyRightBracket:   "RightBracket",
	pixelgl.KeyGraveAccent:    "GraveAccent",
	pixelgl.KeyWorld1:         "World1",
	pixelgl.KeyWorld2:         "World2",
	pixelgl.KeyEscape:         "Escape",
	pixelgl.KeyEnter:          "Enter",
	pixelgl.KeyTab:            "Tab",
	pixelgl.KeyBackspace:      "Backspace",
	pixelgl.KeyInsert:         "Insert",
	pixelgl.KeyDelete:         "Delete",
	pixelgl.KeyRight:          "Right",
	pixelgl.KeyLeft:           "Left",
	pixelgl.KeyDown:           "Down",
	pixelgl.KeyUp:             "Up",
	pixelgl.KeyPageUp:         "PageUp",
	pixelgl.KeyPageDown:       "PageDown",
	pixelgl.KeyHome:           "Home",
	pixelgl.KeyEnd:            "End",
	pixelgl.KeyCapsLock:       "CapsLock",
	pixelgl.KeyScrollLock:     "ScrollLock",
	pixelgl.KeyNumLock:        "NumLock",
	pixelgl.KeyPrintScreen:    "PrintScreen",
	pixelgl.KeyPause:          "Pause",
	pixelgl.KeyF1:             "F1",
	pixelgl.KeyF2:             "F2",
	pixelgl.KeyF3:             "F3",
	pixelgl.KeyF4:             "F4",
	pixelgl.KeyF5:             "F5",
	pixelgl.KeyF6:             "F6",
	pixelgl.KeyF7:             "F7",
	pixelgl.KeyF8:             "F8",
	pixelgl.KeyF9:             "F9",
	pixelgl.KeyF10:            "F10",
	pixelgl.KeyF11:            "F11",
	pixelgl.KeyF12:            "F12",
	pixelgl.KeyF13:            "F13",
	pixelgl.KeyF14:            "F14",
	pixelgl.KeyF15:            "F15",
	pixelgl.KeyF16:            "F16",
	pixelgl.KeyF17:            "F17",
	pixelgl.KeyF18:            "F18",
	pixelgl.KeyF19:            "F19",
	pixelgl.KeyF20:            "F20",
	pixelgl.KeyF21:            "F21",
	pixelgl.KeyF22:            "F22",
	pixelgl.KeyF23:            "F23",
	pixelgl.KeyF24:            "F24",
	pixelgl.KeyF25:            "F25",
	pixelgl.KeyKP0:            "KP0",
	pixelgl.KeyKP1:            "KP1",
	pixelgl.KeyKP2:            "KP2",
	pixelgl.KeyKP3:            "KP3",
	pixelgl.KeyKP4:            "KP4",
	pixelgl.KeyKP5:            "KP5",
	pixelgl.KeyKP6:            "KP6",
	pixelgl.KeyKP7:            "KP7",
	pixelgl.KeyKP8:            "KP8",
	pixelgl.KeyKP9:            "KP9",
	pixelgl.KeyKPDecimal:      "KPDecimal",
	pixelgl.KeyKPDivide:       "KPDivide",
	pixelgl.KeyKPMultiply:     "KPMultiply",
	pixelgl.KeyKPSubtract:     "KPSubtract",
	pixelgl.KeyKPAdd:          "KPAdd",
	pixelgl.KeyKPEnter:        "KPEnter",
	pixelgl.KeyKPEqual:        "KPEqual",
	pixelgl.KeyLeftShift:      "LeftShift",
	pixelgl.KeyLeftControl:    "LeftControl",
	pixelgl.KeyLeftAlt:        "LeftAlt",
	pixelgl.KeyLeftSuper:      "LeftSuper",
	pixelgl.KeyRightShift:     "RightShift",
	pixelgl.KeyRightControl:   "RightControl",
	pixelgl.KeyRightAlt:       "RightAlt",
	pixelgl.KeyRightSuper:     "RightSuper",
	pixelgl.KeyMenu:           "Menu",
}

type Display struct {
	config *config.Config

	win            *pixelgl.Window
	canvas         *pixelgl.Canvas
	updateScreenCh chan any

	bus *bus.Bus

	coeff           int
	resizeLock      sync.Mutex
	pixelBuffer     []uint8
	currentWinSizeX float64
	currentWinSizeY float64

	buttonsJoyPad1 map[bus.JoyPadButton]pixelgl.Button
}

func (d *Display) Init() {
	d.coeff = 1
	d.currentWinSizeX = NesImageSizeX
	d.currentWinSizeY = NesImageSizeY
	d.pixelBuffer = make([]uint8, (d.coeff*NesImageSizeX)*(d.coeff*NesImageSizeY)*4)
	d.updateScreenCh = make(chan any)

	d.buttonsJoyPad1 = make(map[bus.JoyPadButton]pixelgl.Button)

	for button, name := range buttonNames {
		switch name {
		case d.config.JoyPad1.A:
			d.buttonsJoyPad1[bus.JoyPadButtonA] = button
			break
		case d.config.JoyPad1.B:
			d.buttonsJoyPad1[bus.JoyPadButtonB] = button
			break
		case d.config.JoyPad1.Start:
			d.buttonsJoyPad1[bus.JoyPadButtonStart] = button
			break
		case d.config.JoyPad1.Select:
			d.buttonsJoyPad1[bus.JoyPadButtonSelect] = button
			break
		case d.config.JoyPad1.Up:
			d.buttonsJoyPad1[bus.JoyPadButtonUp] = button
			break
		case d.config.JoyPad1.Down:
			d.buttonsJoyPad1[bus.JoyPadButtonDown] = button
			break
		case d.config.JoyPad1.Left:
			d.buttonsJoyPad1[bus.JoyPadButtonLeft] = button
			break
		case d.config.JoyPad1.Right:
			d.buttonsJoyPad1[bus.JoyPadButtonRight] = button
			break
		}
	}
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
	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonUp]) {
		d.bus.KeyEvent(bus.JoyPadButtonUp, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonUp, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonDown]) {
		d.bus.KeyEvent(bus.JoyPadButtonDown, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonDown, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonLeft]) {
		d.bus.KeyEvent(bus.JoyPadButtonLeft, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonLeft, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonRight]) {
		d.bus.KeyEvent(bus.JoyPadButtonRight, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonRight, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonA]) {
		d.bus.KeyEvent(bus.JoyPadButtonA, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonA, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonB]) {
		d.bus.KeyEvent(bus.JoyPadButtonB, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonB, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonSelect]) {
		d.bus.KeyEvent(bus.JoyPadButtonSelect, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonSelect, false)
	}

	if d.win.Pressed(d.buttonsJoyPad1[bus.JoyPadButtonStart]) {
		d.bus.KeyEvent(bus.JoyPadButtonStart, true)
	} else {
		d.bus.KeyEvent(bus.JoyPadButtonStart, false)
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
