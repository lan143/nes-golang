package pixelgl

import (
	"github.com/faiface/pixel/pixelgl"
	"main/src/bus"
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

func (d *Display) initInput() {
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
