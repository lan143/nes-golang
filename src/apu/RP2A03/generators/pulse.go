package generators

type Pulse struct {
	enabled bool

	timerCounter  uint16
	timerPeriod   uint16
	timerSequence byte

	envelopeStartFlag         bool
	envelopeCounter           byte
	envelopeDecayLevelCounter byte

	lengthCounter byte

	sweepReloadFlag bool
	sweepCycle      uint32
	sweepCounter    byte

	registers [4]byte

	Channel byte

	dutyTable [4][8]byte
}

func (g *Pulse) Init() {
	g.dutyTable = [4][8]byte{
		{0, 1, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0},
		{1, 0, 0, 1, 1, 1, 1, 1},
	}
}

func (g *Pulse) DriveTimer() {
	if g.timerCounter > 0 {
		g.timerCounter--
	} else {
		g.timerCounter = g.timerPeriod
		g.timerSequence++

		if g.timerSequence == 8 {
			g.timerSequence = 0
		}
	}
}

func (g *Pulse) DriveEnvelope() {
	if g.envelopeStartFlag {
		g.envelopeCounter = g.getEnvelopePeriod()
		g.envelopeDecayLevelCounter = 0xf
		g.envelopeStartFlag = false
		return
	}

	if g.envelopeCounter > 0 {
		g.envelopeCounter--
	} else {
		g.envelopeCounter = g.getEnvelopePeriod()

		if g.envelopeDecayLevelCounter > 0 {
			g.envelopeDecayLevelCounter--
		} else if g.envelopeDecayLevelCounter == 0 && g.isEnabledEnvelopeLoop() {
			g.envelopeDecayLevelCounter = 0xf
		}
	}
}

func (g *Pulse) DriveLinear() {}

func (g *Pulse) DriveLength() {
	if !g.isDisabledEnvelope() && g.lengthCounter > 0 {
		g.lengthCounter--
	}
}

func (g *Pulse) DriveSweep() {
	if g.sweepCounter > 0 && g.isEnabledSweep() && g.getSweepShiftAmount() != 0 && g.timerPeriod >= 8 && g.timerPeriod <= 0x7FF {
		change := g.timerPeriod >> g.getSweepShiftAmount()

		if g.isNegatedSweep() {
			g.timerPeriod -= change

			if g.Channel == 1 {
				g.timerPeriod--
			}
		} else {
			g.timerPeriod += change
		}
	}

	if g.sweepReloadFlag || g.sweepCounter == 0 {
		g.sweepReloadFlag = false
		g.sweepCounter = g.getSweepPeriod()
	} else {
		g.sweepCounter--
	}
}

func (g *Pulse) SetValue(index byte, value byte) {
	g.registers[index] = value

	if index == 1 {
		g.sweepReloadFlag = true
	} else if index == 2 {
		g.timerPeriod = g.getTimer()
	} else if index == 3 {
		if g.enabled {
			g.lengthCounter = lengthTable[g.getLengthCounterIndex()]
		}

		g.timerPeriod = g.getTimer()
		g.timerSequence = 0
		g.envelopeStartFlag = true
	}
}

func (g *Pulse) SetEnabled(enabled bool) {
	g.enabled = enabled

	if !enabled {
		g.lengthCounter = 0
	}
}

func (g *Pulse) GetOutput() byte {
	if g.lengthCounter == 0 || g.timerPeriod < 8 || g.timerPeriod > 0x7FF || g.dutyTable[g.getDuty()][g.timerSequence] == 0 {
		return 0
	}

	if g.isDisabledEnvelope() {
		return g.getEnvelopePeriod() & 0xF
	} else {
		return g.envelopeDecayLevelCounter & 0xF
	}
}

func (g *Pulse) GetRemainingBytes() uint16 { return 0 }

func (g *Pulse) GetLengthCounter() byte { return g.lengthCounter }

func (g *Pulse) getTimer() uint16 {
	return (uint16((g.registers[3])&((1<<3)-1)) << 8) | uint16(g.registers[2])
}

func (g *Pulse) getLengthCounterIndex() byte {
	return (g.registers[3] >> 3) & ((1 << 5) - 1)
}

func (g *Pulse) getDuty() byte {
	return (g.registers[0] >> 6) & ((1 << 2) - 1)
}

func (g *Pulse) isEnabledSweep() bool {
	return g.registers[1]&0x80 > 0
}

func (g *Pulse) getSweepPeriod() byte {
	return g.registers[1] >> 4 & ((1 << 3) - 1)
}

func (g *Pulse) isNegatedSweep() bool {
	return g.registers[1]&0x08 > 0
}

func (g *Pulse) getSweepShiftAmount() byte {
	return g.registers[1] & ((1 << 3) - 1)
}

func (g *Pulse) isEnabledEnvelopeLoop() bool {
	return g.registers[0]&0x20 > 0
}

func (g *Pulse) isDisabledEnvelope() bool {
	return g.registers[0]&0x10 > 0
}

func (g *Pulse) getEnvelopePeriod() byte {
	return g.registers[0] & ((1 << 4) - 1)
}
