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

	playChan chan []byte
}

func (a *Audio) Init() error {
	var err error
	a.otoCtx, err = oto.NewContext(int(a.GetSampleRate()), 1, 1, AudioBufferSize)
	if err != nil {
		return err
	}

	a.player = a.otoCtx.NewPlayer()

	a.playChan = make(chan []byte)
	go a.playInDevice()

	return nil
}

func (a *Audio) GetSampleRate() uint32 {
	return 44100
}

func (a *Audio) PlaySample(sample float32) {
	if a.bufferIndex < AudioBufferSize {
		a.buffer[a.bufferIndex] = byte(sample * ((1 << 7) - 1))
		a.bufferIndex++
	} else {
		a.bufferIndex = 0

		playArray := make([]byte, len(a.buffer), len(a.buffer))

		for i, val := range a.buffer {
			playArray[i] = val
		}

		a.playChan <- playArray
	}
}

func (a *Audio) playInDevice() {
	for {
		data := <-a.playChan

		_, err := a.player.Write(data)
		if err != nil {
			log.Printf("playInDevice error: %s", err.Error())
		}
	}
}

func NewAudio() *Audio {
	return &Audio{}
}
