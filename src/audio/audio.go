package audio

type Audio interface {
	Init() error
	GetSampleRate() uint32
	PlaySample(sample float32)
}
