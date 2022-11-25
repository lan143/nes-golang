package ram

type Ram struct {
	memory    []byte
	memoryLen uint16
}

func (r *Ram) Init(size int) {
	r.memory = make([]byte, size)
	r.memoryLen = uint16(len(r.memory))
}

func (r *Ram) GetByte(address uint16) byte {
	address &= r.memoryLen - 1

	return r.memory[address]
}

func (r *Ram) SetByte(address uint16, value byte) {
	address &= r.memoryLen - 1

	r.memory[address] = value
}
