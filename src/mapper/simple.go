package mapper

type SimpleMapper struct {
	memory [0xffff]byte
}

func (m *SimpleMapper) LoadRom(data []byte) {
	var i, j uint16
	j = 0

	// Prg ROM
	for i = 0x8000; i <= 0xBFFF; i++ {
		m.memory[i] = data[j]
		j++
	}

	// Prg ROM
	for i = 0xC000; i <= 0xFFFE; i++ {
		m.memory[i] = data[j]
		j++
	}

	// Chr ROM
	for i = 0x0000; i <= 0x1FFF; i++ {
		m.memory[i] = data[j]
		j++
	}
}

func (m *SimpleMapper) GetByte(address uint16) byte {
	return m.memory[address]
}

func (m *SimpleMapper) PutByte(address uint16, value byte) {
	m.memory[address] = value
}
