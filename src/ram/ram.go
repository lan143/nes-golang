package ram

type Ram struct {
	memory    []byte
	memoryLen int
}

func (r *Ram) Init(size int) {
	r.memory = make([]byte, size)
	r.memoryLen = len(r.memory)
}

func (r *Ram) GetByte(address uint16) byte {
	address &= uint16(r.memoryLen) - 1

	return r.memory[address]
}

func (r *Ram) SetByte(address uint16, value byte) {
	address &= uint16(len(r.memory)) - 1

	r.memory[address] = value
}
