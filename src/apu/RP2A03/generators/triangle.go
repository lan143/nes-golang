package generators

type Triangle struct {
	enabled bool

	timerCounter  uint16
	timerSequence byte

	lengthCounter byte

	linearReloadFlag bool
	linearCounter    byte

	registers [4]byte

	sequenceTable [32]byte
}

func (g *Triangle) Init() {
	g.sequenceTable = [32]byte{
		15, 14, 13, 12, 11, 10, 9, 8,
		7, 6, 5, 4, 3, 2, 1, 0,
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
	}
}

func (g *Triangle) DriveTimer() {
	if g.timerCounter > 0 {
		g.timerCounter--
	} else {
		g.timerCounter = g.getTimer()

		if g.lengthCounter > 0 && g.linearCounter > 0 {
			g.timerSequence++

			if g.timerSequence == 32 {
				g.timerSequence = 0
			}
		}
	}
}

func (g *Triangle) DriveEnvelope() {}

func (g *Triangle) DriveLinear() {
	if g.linearReloadFlag {
		g.linearCounter = g.getLinearCounter()
	} else if g.linearCounter > 0 {
		g.linearCounter--
	}

	if !g.isDisabledLengthCounter() {
		g.linearReloadFlag = true
	}
}

func (g *Triangle) DriveLength() {
	if !g.isDisabledLengthCounter() && g.lengthCounter > 0 {
		g.lengthCounter--
	}
}

func (g *Triangle) DriveSweep() {}

func (g *Triangle) SetValue(index byte, value byte) {
	g.registers[index] = value

	if index == 3 {
		if g.enabled {
			g.lengthCounter = lengthTable[g.getLengthCounterIndex()]
		}

		g.timerSequence = 0
		g.linearReloadFlag = true
	}
}

func (g *Triangle) SetEnabled(enabled bool) {
	g.enabled = enabled

	if !enabled {
		g.lengthCounter = 0
	}
}

func (g *Triangle) GetOutput() byte {
	if !g.enabled || g.lengthCounter == 0 || g.linearCounter == 0 || g.getTimer() < 2 {
		return 0
	}

	return g.sequenceTable[g.timerSequence] & 0xf
}

func (g *Triangle) GetRemainingBytes() uint16 { return 0 }

func (g *Triangle) GetLengthCounter() byte { return g.lengthCounter }

// this.register0.loadBits(0, 7)
func (g *Triangle) getLinearCounter() byte {
	return g.registers[0] & 0x7F
}

// this.register0.isBitSet(7)
func (g *Triangle) isDisabledLengthCounter() bool {
	return g.registers[0]&0x40 > 0
}

// this.register3.loadBits(3, 5)
func (g *Triangle) getLengthCounterIndex() byte {
	return g.registers[3] >> 3
}

// this.register3.loadBits(0, 3) | this.register2
func (g *Triangle) getTimer() uint16 {
	return (uint16((g.registers[3])&0x07) << 8) | uint16(g.registers[2])
}
