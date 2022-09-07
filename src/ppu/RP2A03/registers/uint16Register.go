package registers

import "sync"

type Uint16Register struct {
	value uint16
	mx    sync.RWMutex
}

func (r *Uint16Register) Get() uint16 {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.value
}

func (r *Uint16Register) Set(value uint16) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.value = value
}

func (r *Uint16Register) Shift(value uint16) uint16 {
	r.mx.Lock()
	defer r.mx.Unlock()

	value &= 0x1
	carry := value & (1 << 15)
	r.value = (r.value << 1) | value

	return carry
}

func (r *Uint16Register) StoreLowerByte(value byte) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.value = uint16(value)
}
