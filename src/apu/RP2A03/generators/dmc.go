package generators

import "main/src/bus"

type DMC struct {
	enabled bool

	timerCounter uint16
	timerPeriod  uint16

	deltaCounter          byte
	addressCounter        uint16
	remainingBytesCounter uint16

	sampleAddress       uint16
	sampleBuffer        byte
	sampleBufferIsEmpty bool

	shiftRegister        byte
	remainingBitsCounter byte

	silenceFlag bool

	registers [4]byte

	timerTable [16]uint16

	Bus *bus.Bus
}

func (g *DMC) Init() {
	g.timerTable = [16]uint16{
		0x1AC, 0x17C, 0x154, 0x140,
		0x11E, 0x0FE, 0x0E2, 0x0D6,
		0x0BE, 0x0A0, 0x08E, 0x080,
		0x06A, 0x054, 0x048, 0x036,
	}
}

func (g *DMC) DriveTimer() {
	if g.timerCounter > 0 {
		g.timerCounter--
	} else {
		g.timerCounter = g.timerPeriod

		if g.remainingBytesCounter > 0 && g.sampleBufferIsEmpty {
			g.sampleBuffer = g.Bus.ReadFromCPU(g.sampleAddress)
			g.sampleAddress++
			g.sampleBufferIsEmpty = false

			if g.sampleAddress == 0xFFFF {
				g.sampleAddress = 0x8000
			}

			g.remainingBytesCounter--

			if g.remainingBytesCounter == 0 {
				if g.isLoop() {
					g.start()
				} else if g.isEnabledIrq() {
					g.Bus.ApuDMCActivate()
				}
			}

			g.Bus.CPUSkipCycles(4)
		}

		if g.remainingBitsCounter == 0 {
			g.remainingBitsCounter = 8

			if g.sampleBufferIsEmpty {
				g.silenceFlag = true
			} else {
				g.silenceFlag = false
				g.sampleBufferIsEmpty = true
				g.shiftRegister = g.sampleBuffer
				g.sampleBuffer = 0
			}
		}

		if !g.silenceFlag {
			if g.shiftRegister&0x1 == 0 {
				if g.deltaCounter > 1 {
					g.deltaCounter -= 2
				} else if g.deltaCounter < 126 {
					g.deltaCounter += 2
				}
			}
		}

		g.shiftRegister >>= 1
		g.remainingBitsCounter--
	}
}

func (g *DMC) DriveEnvelope() {}

func (g *DMC) DriveLinear() {}

func (g *DMC) DriveLength() {}

func (g *DMC) DriveSweep() {}

func (g *DMC) SetValue(index byte, value byte) {
	g.registers[index] = value

	if index == 0 {
		g.timerPeriod = g.timerTable[g.getTimerIndex()] >> 1
	} else {
		g.start()
	}
}

func (g *DMC) SetEnabled(enabled bool) {
	g.enabled = enabled

	if enabled {
		if g.remainingBytesCounter == 0 {
			g.start()
		}
	} else {
		g.remainingBytesCounter = 0
	}
}

func (g *DMC) GetOutput() byte {
	if g.silenceFlag {
		return 0
	}

	return g.deltaCounter & 0x7F
}

func (g *DMC) GetRemainingBytes() uint16 { return g.remainingBytesCounter }

func (g *DMC) GetLengthCounter() byte { return 0 }

// this.register0.isBitSet(7)
func (g *DMC) isEnabledIrq() bool {
	return g.registers[0]&0x40 > 0
}

// this.register0.isBitSet(6)
func (g *DMC) isLoop() bool {
	return g.registers[0]&0x20 > 0
}

// this.register0.loadBits(0, 4)
func (g *DMC) getTimerIndex() byte {
	return g.registers[0] & 0x0F
}

// this.register1.loadBits(0, 7)
func (g *DMC) getDeltaCounter() byte {
	return g.registers[1] & 0x7F
}

// this.register2
func (g *DMC) getSampleAddress() byte {
	return g.registers[2]
}

// this.register3
func (g *DMC) getSampleLength() byte {
	return g.registers[3]
}

func (g *DMC) start() {
	g.deltaCounter = g.getDeltaCounter()
	g.addressCounter = uint16(g.getSampleAddress())*0x40 + 0xC000
	g.remainingBytesCounter = uint16(g.getSampleLength()*0x10) + 1
}
