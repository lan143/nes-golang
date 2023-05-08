package oto

import (
	"github.com/hajimehoshi/oto"
	"log"
)

const (
	AudioBufferSize int = 8192
)

type Audio struct {
	otoCtx *oto.Context
	player *oto.Player

	playChan chan byte
}

func (a *Audio) Init() error {
	a.playChan = make(chan byte, AudioBufferSize)

	go a.playInDevice()

	return nil
}

func (a *Audio) GetSampleRate() uint32 {
	return 44100
}

func (a *Audio) PlaySample(sample float32) {
	if len(a.playChan) < AudioBufferSize {
		a.playChan <- byte(sample * ((1 << 7) - 1))
	}
}

func (a *Audio) playInDevice() {
	var err error
	a.otoCtx, err = oto.NewContext(int(a.GetSampleRate()), 1, 1, AudioBufferSize)
	if err != nil {
		log.Printf("OTO init context err: %s", err.Error())
		return
	}

	a.player = a.otoCtx.NewPlayer()

	for {
		_, err := a.player.Write([]byte{<-a.playChan})
		if err != nil {
			log.Printf("playInDevice error: %s", err.Error())
		}
	}
}

func NewAudio() *Audio {
	return &Audio{}
}
