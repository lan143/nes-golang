package joypad

import (
	"main/src/bus"
)

type JoyPad struct {
	currentButton byte
	buttons       [8]bool
	latch         byte

	bus *bus.Bus
}

func (j *JoyPad) Init() {
	j.bus.OnCPURead(0x4016, func() byte {
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
		switch button {
		case bus.JoyPadButtonA:
			j.buttons[0] = pressed
		case bus.JoyPadButtonB:
			j.buttons[1] = pressed
		case bus.JoyPadButtonSelect:
			j.buttons[2] = pressed
		case bus.JoyPadButtonStart:
			j.buttons[3] = pressed
		case bus.JoyPadButtonUp:
			j.buttons[4] = pressed
		case bus.JoyPadButtonDown:
			j.buttons[5] = pressed
		case bus.JoyPadButtonLeft:
			j.buttons[6] = pressed
		case bus.JoyPadButtonRight:
			j.buttons[7] = pressed
		}
	})
}

func NewJoyPad(bus *bus.Bus) *JoyPad {
	return &JoyPad{
		bus: bus,
	}
}
