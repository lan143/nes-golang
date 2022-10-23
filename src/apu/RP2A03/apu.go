package RP2A03

import (
	"main/src/apu/RP2A03/generators"
	"main/src/apu/RP2A03/registers"
	"main/src/audio"
	"main/src/bus"
)

type APU struct {
	dmc      generators.Generator
	noise    generators.Generator
	pulse1   generators.Generator
	pulse2   generators.Generator
	triangle generators.Generator

	status registers.StatusRegister // 0x4015
	frame  registers.FrameRegister  // 0x4017

	cycle        uint64
	step         byte
	samplePeriod uint32

	frameIrqActive bool
	dmcIrqActive   bool

	pulseTable [32]float32
	tndTable   [203]float32

	b     *bus.Bus
	audio audio.Audio
}

func (a *APU) Init(sampleRate uint32, audio audio.Audio) {
	for i := 0; i < 31; i++ {
		a.pulseTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}

	for i := 0; i < 203; i++ {
		a.tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}

	a.audio = audio

	a.dmc = &generators.DMC{Bus: a.b}
	a.dmc.Init()

	a.noise = &generators.Noise{}
	a.noise.Init()

	a.pulse1 = &generators.Pulse{Channel: 1}
	a.pulse1.Init()

	a.pulse2 = &generators.Pulse{Channel: 2}
	a.pulse2.Init()

	a.triangle = &generators.Triangle{}
	a.triangle.Init()

	a.status.SetValue(0x00)

	a.samplePeriod = 1789773 / sampleRate

	a.b.OnCPUWrite(0x4000, func(value byte) {
		a.pulse1.SetValue(0, value)
	})
	a.b.OnCPUWrite(0x4001, func(value byte) {
		a.pulse1.SetValue(1, value)
	})
	a.b.OnCPUWrite(0x4002, func(value byte) {
		a.pulse1.SetValue(2, value)
	})
	a.b.OnCPUWrite(0x4003, func(value byte) {
		a.pulse1.SetValue(3, value)
	})

	a.b.OnCPUWrite(0x4004, func(value byte) {
		a.pulse2.SetValue(0, value)
	})
	a.b.OnCPUWrite(0x4005, func(value byte) {
		a.pulse2.SetValue(1, value)
	})
	a.b.OnCPUWrite(0x4006, func(value byte) {
		a.pulse2.SetValue(2, value)
	})
	a.b.OnCPUWrite(0x4007, func(value byte) {
		a.pulse2.SetValue(3, value)
	})

	a.b.OnCPUWrite(0x4008, func(value byte) {
		a.triangle.SetValue(0, value)
	})
	a.b.OnCPUWrite(0x4009, func(value byte) {
		a.triangle.SetValue(1, value)
	})
	a.b.OnCPUWrite(0x400A, func(value byte) {
		a.triangle.SetValue(2, value)
	})
	a.b.OnCPUWrite(0x400B, func(value byte) {
		a.triangle.SetValue(3, value)
	})

	a.b.OnCPUWrite(0x400C, func(value byte) {
		a.noise.SetValue(0, value)
	})
	a.b.OnCPUWrite(0x400D, func(value byte) {
		a.noise.SetValue(1, value)
	})
	a.b.OnCPUWrite(0x400E, func(value byte) {
		a.noise.SetValue(2, value)
	})
	a.b.OnCPUWrite(0x400F, func(value byte) {
		a.noise.SetValue(3, value)
	})

	a.b.OnCPUWrite(0x4010, func(value byte) {
		a.dmc.SetValue(0, value)
	})
	a.b.OnCPUWrite(0x4011, func(value byte) {
		a.dmc.SetValue(1, value)
	})
	a.b.OnCPUWrite(0x4012, func(value byte) {
		a.dmc.SetValue(2, value)
	})
	a.b.OnCPUWrite(0x4013, func(value byte) {
		a.dmc.SetValue(3, value)
	})

	a.b.OnCPUWrite(0x4015, func(value byte) {
		a.status.SetValue(value)

		a.pulse1.SetEnabled(a.status.IsEnabledPulse1())
		a.pulse2.SetEnabled(a.status.IsEnabledPulse2())
		a.triangle.SetEnabled(a.status.IsEnabledTriangle())
		a.noise.SetEnabled(a.status.IsEnabledNoise())
		a.dmc.SetEnabled(a.status.IsEnabledDMC())

		a.dmcIrqActive = true
	})

	a.b.OnCPUWrite(0x4017, func(value byte) {
		a.frame.SetValue(value)

		if a.frame.IsDisabledIrq() {
			a.frameIrqActive = false
		}
	})

	a.b.OnCPURead(0x4015, func() byte {
		var value byte = 0

		if a.dmcIrqActive {
			value |= 1 << 7
		}

		if a.frameIrqActive && !a.frame.IsDisabledIrq() {
			value |= 1 << 6
		}

		if a.dmc.GetRemainingBytes() > 0 {
			value |= 1 << 4
		}

		if a.noise.GetLengthCounter() > 0 {
			value |= 1 << 3
		}

		if a.triangle.GetLengthCounter() > 0 {
			value |= 1 << 2
		}

		if a.pulse2.GetLengthCounter() > 0 {
			value |= 1 << 1
		}

		if a.pulse1.GetLengthCounter() > 0 {
			value |= 1
		}

		a.frameIrqActive = false

		return value
	})

	a.b.OnApuDMCActivate(func() {
		a.dmcIrqActive = true
	})
}

func (a *APU) Run() {
	a.runCycle()
}

func (a *APU) runCycle() {
	a.cycle++

	if a.cycle%uint64(a.samplePeriod) == 0 {
		a.sample()
	}

	if a.cycle%2 == 0 {
		a.pulse1.DriveTimer()
		a.pulse2.DriveTimer()
		a.noise.DriveTimer()
		a.dmc.DriveTimer()
	}

	a.triangle.DriveTimer()

	if a.cycle%7457 == 0 {
		if a.frame.IsFiveStepMode() {
			if a.step < 4 {
				a.pulse1.DriveEnvelope()
				a.pulse2.DriveEnvelope()
				a.triangle.DriveLinear()
				a.noise.DriveEnvelope()
			}

			if a.step == 0 || a.step == 2 {
				a.pulse1.DriveLength()
				a.pulse1.DriveSweep()

				a.pulse2.DriveLength()
				a.pulse2.DriveSweep()

				a.triangle.DriveLength()
				a.noise.DriveLength()
			}

			a.step = (a.step + 1) % 5
		} else {
			a.pulse1.DriveEnvelope()
			a.pulse2.DriveEnvelope()
			a.triangle.DriveLinear()
			a.noise.DriveEnvelope()

			if a.step == 1 || a.step == 3 {
				a.pulse1.DriveLength()
				a.pulse1.DriveSweep()

				a.pulse2.DriveLength()
				a.pulse2.DriveSweep()

				a.triangle.DriveLength()
				a.noise.DriveLength()
			}

			if a.step == 3 && !a.frame.IsDisabledIrq() {
				a.frameIrqActive = true
				a.b.Interrupt(bus.IRQ)
			}

			a.step = (a.step + 1) % 4
		}

		if a.dmcIrqActive {
			a.b.Interrupt(bus.IRQ)
		}
	}
}

func (a *APU) sample() {
	a.audio.PlaySample(
		a.pulseTable[a.pulse1.GetOutput()+a.pulse2.GetOutput()] +
			a.tndTable[3*a.triangle.GetOutput()+2*a.noise.GetOutput()+a.dmc.GetOutput()],
	)
}

func NewApu(b *bus.Bus) *APU {
	return &APU{
		b: b,
	}
}
