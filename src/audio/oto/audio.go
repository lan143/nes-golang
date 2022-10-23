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
	a.otoCtx, err = oto.NewContext(int(a.GetSampleRate()), 1, 1, AudioBufferSize*4)
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
		/*sample = (sample - 0.5) * 0.5
		int8Sample := int8(sample * ((1 << 7) - 1))
		a.buffer[a.bufferIndex] = byte(int8Sample)*/

		a.buffer[a.bufferIndex] = byte(sample * ((1 << 7) - 1))
		a.bufferIndex++
	} else {
		a.bufferIndex = 0
		a.playChan <- a.buffer[:]
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
