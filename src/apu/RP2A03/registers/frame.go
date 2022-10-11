package registers

type APUFrameFlag byte

const (
	ApuFrameMode APUFrameFlag = 0x80
	ApuFrameIrq               = 0x40
)

type FrameRegister struct {
	value APUFrameFlag
}

func (r *FrameRegister) GetValue() byte {
	return byte(r.value)
}

func (r *FrameRegister) SetValue(value byte) {
	r.value = APUFrameFlag(value)
}

func (r *FrameRegister) IsFiveStepMode() bool {
	return r.isFlag(ApuFrameMode)
}

func (r *FrameRegister) IsDisabledIrq() bool {
	return r.isFlag(ApuFrameIrq)
}

func (r *FrameRegister) setFlag(flag APUFrameFlag) {
	r.value |= flag
}

func (r *FrameRegister) isFlag(flag APUFrameFlag) bool {
	return r.value&flag > 0
}

func (r *FrameRegister) clearFlag(flag APUFrameFlag) {
	r.value &= ^flag
}
