package audio

import (
	"main/src/audio/null"
	"main/src/audio/oto"
)

type Factory struct {
	null *null.Audio
	oto  *oto.Audio
}

func (f *Factory) GetAudio() Audio {
	return f.oto
}

func NewFactory(
	null *null.Audio,
	oto *oto.Audio,
) *Factory {
	return &Factory{
		null: null,
		oto:  oto,
	}
}
