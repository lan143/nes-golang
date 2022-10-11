package oto

import (
	"github.com/hajimehoshi/oto"
	"log"
)

const (
	AudioBufferSize int = 4096
)

type Audio struct {
	otoCtx *oto.Context
	player *oto.Player

	buffer      [AudioBufferSize]byte
	bufferIndex int
}

func (a *Audio) Init() error {
	var err error
	a.otoCtx, err = oto.NewContext(int(a.GetSampleRate()), 1, 1, 8912)
	if err != nil {
		return err
	}

	a.player = a.otoCtx.NewPlayer()

	return nil
}

func (a *Audio) GetSampleRate() uint32 {
	return 48000
}

func (a *Audio) PlaySample(sample float32) {
	if a.bufferIndex < AudioBufferSize {
		a.buffer[a.bufferIndex] = byte(sample / 1.0 * 255)
		a.bufferIndex++
	} else {
		a.bufferIndex = 0

		go func(data []byte) {
			_, err := a.player.Write(data)
			if err != nil {
				log.Printf("PlaySample error: %s", err.Error())
			}
		}(a.buffer[:])
	}
}

func NewAudio() *Audio {
	return &Audio{}
}
