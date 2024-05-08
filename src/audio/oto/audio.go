package oto

import (
	"github.com/hajimehoshi/oto"
	"log"
	"main/src/audio/filter"
)

const (
	AudioBufferSize int = 8192
)

type Audio struct {
	otoCtx *oto.Context
	player *oto.Player

	playChan chan byte

	f1, f2 *filter.HighPassFilter
	f3     *filter.LowPassFilter
}

func (a *Audio) Init() error {
	a.playChan = make(chan byte, AudioBufferSize)
	a.f1 = filter.NewHighPassFilter(90, a.GetSampleRate())
	a.f2 = filter.NewHighPassFilter(440, a.GetSampleRate())
	a.f3 = filter.NewLowPassFilter(14000, a.GetSampleRate())

	go a.playInDevice()

	return nil
}

func (a *Audio) GetSampleRate() uint32 {
	return 44100
}

func (a *Audio) PlaySample(sample float32) {
	sample = a.f1.Filter(sample)
	sample = a.f2.Filter(sample)
	sample = a.f3.Filter(sample)

	a.playChan <- byte(sample * ((1 << 7) - 1))
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
