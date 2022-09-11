package Ricoh6502

import "testing"

func TestUpdateN(t *testing.T) {
	register := PRegister{}

	register.UpdateN(34)
	if register.IsN() {
		t.Error("Dont should be N set")
	}

	register.UpdateN(47)
	if register.IsN() {
		t.Error("Dont should be N set")
	}

	register.UpdateN(123)
	if register.IsN() {
		t.Error("Dont should be N set")
	}

	register.UpdateN(0x81)
	if !register.IsN() {
		t.Error("Should be N set")
	}

	register.UpdateN(0x83)
	if !register.IsN() {
		t.Error("Should be N set")
	}

	register.UpdateN(0x99)
	if !register.IsN() {
		t.Error("Should be N set")
	}
}

func TestUpdateZ(t *testing.T) {
	register := PRegister{}

	register.UpdateZ(0)
	if !register.IsZ() {
		t.Error("Should be Z set")
	}

	register.UpdateZ(1)
	if register.IsZ() {
		t.Error("Dont should be Z set")
	}

	register.UpdateZ(15)
	if register.IsZ() {
		t.Error("Dont should be Z set")
	}

	register.UpdateZ(147)
	if register.IsZ() {
		t.Error("Dont should be Z set")
	}

	register.UpdateZ(223)
	if register.IsZ() {
		t.Error("Dont should be Z set")
	}
}

func TestUpdateC(t *testing.T) {
	register := PRegister{}

	register.UpdateC(0x80 << 1)
	if !register.IsC() {
		t.Error("Should be C set")
	}

	register.UpdateC(0x40 << 1)
	if register.IsC() {
		t.Error("Dont should be C set")
	}

	register.UpdateC(0xFF << 1)
	if !register.IsC() {
		t.Error("Should be C set")
	}
}
