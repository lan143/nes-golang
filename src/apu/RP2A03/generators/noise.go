package generators

type Noise struct {
	enabled bool

	timerCounter uint16
	timerPeriod  uint16

	envelopeStartFlag         bool
	envelopeCounter           byte
	envelopeDecayLevelCounter byte

	lengthCounter byte

	registers     [4]byte
	shiftRegister uint16

	timerTable [16]uint16
}

func (g *Noise) Init() {
	g.shiftRegister = 0x01

	g.timerTable = [16]uint16{
		0x004, 0x008, 0x010, 0x020,
		0x040, 0x060, 0x080, 0x0A0,
		0x0CA, 0x0FE, 0x17C, 0x1FC,
		0x2FA, 0x3F8, 0x7F2, 0xFE4,
	}
}

func (g *Noise) DriveTimer() {
	if g.timerCounter > 0 {
		g.timerCounter--
	} else {
		g.timerCounter = g.timerPeriod

		var shiftCount byte
		if g.isRandom() {
			shiftCount = 6
		} else {
			shiftCount = 1
		}

		feedback := (g.shiftRegister & 0x1) ^ ((g.shiftRegister >> shiftCount) & 0x1)
		g.shiftRegister = (feedback << 14) | (g.shiftRegister >> 1)
	}
}

func (g *Noise) DriveEnvelope() {
	if g.envelopeStartFlag {
		g.envelopeCounter = g.getEnvelopePeriod()
		g.envelopeDecayLevelCounter = 0xf
		g.envelopeStartFlag = false
	}

	if g.envelopeCounter > 0 {
		g.envelopeCounter--
	} else {
		g.envelopeCounter = g.getEnvelopePeriod()

		if g.envelopeDecayLevelCounter > 0 {
			g.envelopeDecayLevelCounter--
		} else if g.envelopeDecayLevelCounter == 0 && g.isDisabledLengthCounter() {
			g.envelopeDecayLevelCounter = 0xf
		}
	}
}

func (g *Noise) DriveLinear() {}

func (g *Noise) DriveLength() {
	if !g.isDisabledLengthCounter() && g.lengthCounter > 0 {
		g.lengthCounter--
	}
}

func (g *Noise) DriveSweep() {}

func (g *Noise) SetValue(index byte, value byte) {
	g.registers[index] = value

	if index == 2 {
		g.timerPeriod = g.timerTable[g.getTimerIndex()]
	} else if index == 3 {
		if g.enabled {
			g.lengthCounter = lengthTable[g.getLengthCounterIndex()]
		}

		g.envelopeStartFlag = true
	}
}

func (g *Noise) SetEnabled(enabled bool) {
	g.enabled = enabled

	if !enabled {
		g.lengthCounter = 0
	}
}

func (g *Noise) GetOutput() byte {
	if g.lengthCounter == 0 || g.shiftRegister&0x1 == 1 {
		return 0
	}

	if g.isDisabledEnvelope() {
		return g.getEnvelopePeriod() & 0xf
	} else {
		return g.envelopeDecayLevelCounter & 0xf
	}
}

func (g *Noise) GetRemainingBytes() uint16 { return 0 }

func (g *Noise) GetLengthCounter() byte { return g.lengthCounter }

func (g *Noise) isDisabledLengthCounter() bool {
	return g.registers[0]&0x20 > 0
}

func (g *Noise) isRandom() bool {
	return g.registers[2]&0x80 > 0
}

func (g *Noise) getTimerIndex() byte {
	return g.registers[2] & 0xF
}

func (g *Noise) getLengthCounterIndex() byte {
	return g.registers[3] >> 3
}

func (g *Noise) isDisabledEnvelope() bool {
	return g.registers[0]&0x10 > 0
}

func (g *Noise) getEnvelopePeriod() byte {
	return g.registers[0] & 0xF
}
