package null

type Audio struct {
}

func (a *Audio) Init() error {
	return nil
}

func (a *Audio) GetSampleRate() uint32 {
	return 1
}

func (a *Audio) PlaySample(sample float32) {

}

func NewAudio() *Audio {
	return &Audio{}
}
