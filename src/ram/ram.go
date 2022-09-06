package ram

type Ram struct {
	memory []byte
}

func (r *Ram) Init(size int) {
	r.memory = make([]byte, size)
}

func (r *Ram) GetByte(address uint16) byte {
	address &= uint16(len(r.memory)) - 1

	return r.memory[address]
}

func (r *Ram) SetByte(address uint16, value byte) {
	address &= uint16(len(r.memory)) - 1

	r.memory[address] = value
}
