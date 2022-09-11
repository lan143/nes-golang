package joypad

import (
	"main/src/bus"
	"sync"
)

type JoyPad struct {
	currentButton byte
	buttons       [8]bool
	buttonsMx     sync.RWMutex
	latch         byte

	bus *bus.Bus
}

func (j *JoyPad) Init() {
	j.bus.OnCPURead(0x4016, func() byte {
		j.buttonsMx.RLock()
		defer j.buttonsMx.RUnlock()

		var button byte

		if j.latch == 1 {
			button = 0
		} else {
			button = j.currentButton
			j.currentButton++
		}

		if button >= 8 || j.buttons[button] {
			return 1
		} else {
			return 0
		}
	})
	j.bus.OnCPUWrite(0x4016, func(value byte) {
		value &= 1

		if value == 1 {
			j.currentButton = 0
		}

		j.latch = value
	})
	j.bus.OnKeyEvent(func(button bus.JoyPadButton, pressed bool) {
		j.buttonsMx.Lock()
		defer j.buttonsMx.Unlock()

		switch button {
		case bus.JoyPadButtonA:
			j.buttons[0] = pressed
			break
		case bus.JoyPadButtonB:
			j.buttons[1] = pressed
			break
		case bus.JoyPadButtonSelect:
			j.buttons[2] = pressed
			break
		case bus.JoyPadButtonStart:
			j.buttons[3] = pressed
			break
		case bus.JoyPadButtonUp:
			j.buttons[4] = pressed
			break
		case bus.JoyPadButtonDown:
			j.buttons[5] = pressed
			break
		case bus.JoyPadButtonLeft:
			j.buttons[6] = pressed
			break
		case bus.JoyPadButtonRight:
			j.buttons[7] = pressed
			break
		}
	})
}

func NewJoyPad(bus *bus.Bus) *JoyPad {
	return &JoyPad{
		bus: bus,
	}
}
