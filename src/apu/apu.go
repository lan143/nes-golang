package apu

import (
	"main/src/audio"
)

type APU interface {
	Init(sampleRate uint32, audio audio.Audio)
	RunCycle()
}
