package registers

type APUStatusFlag byte

const (
	EnablePulse1   APUStatusFlag = 0x01
	EnablePulse2                 = 0x02
	EnableTriangle               = 0x04
	EnableNoise                  = 0x08
	EnableDMC                    = 0x10
)

type StatusRegister struct {
	value APUStatusFlag
}

func (r *StatusRegister) GetValue() byte {
	return byte(r.value)
}

func (r *StatusRegister) SetValue(value byte) {
	r.value = APUStatusFlag(value)
}

func (r *StatusRegister) IsEnabledPulse1() bool {
	return r.isFlag(EnablePulse1)
}

func (r *StatusRegister) IsEnabledPulse2() bool {
	return r.isFlag(EnablePulse2)
}

func (r *StatusRegister) IsEnabledTriangle() bool {
	return r.isFlag(EnableTriangle)
}

func (r *StatusRegister) IsEnabledNoise() bool {
	return r.isFlag(EnableNoise)
}

func (r *StatusRegister) IsEnabledDMC() bool {
	return r.isFlag(EnableDMC)
}

func (r *StatusRegister) setFlag(flag APUStatusFlag) {
	r.value |= flag
}

func (r *StatusRegister) isFlag(flag APUStatusFlag) bool {
	return r.value&flag > 0
}

func (r *StatusRegister) clearFlag(flag APUStatusFlag) {
	r.value &= ^flag
}
