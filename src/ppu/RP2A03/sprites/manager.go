package sprites

type Manager struct {
	memory     []byte
	memorySize uint16

	sprites     []Sprite
	spritesSize byte
}

func (m *Manager) Init(size uint16) {
	m.memorySize = size
	m.spritesSize = byte(m.memorySize / 4)

	m.memory = make([]byte, m.memorySize)
	m.sprites = make([]Sprite, 0, m.spritesSize)

	var i uint16
	for i = 0; i < uint16(m.spritesSize); i++ {
		m.sprites = append(m.sprites, NewSprite(byte(i), byte(i), m.memory[i*4:(i+1)*4]))
	}
}

func (m *Manager) SetByte(address byte, value byte) {
	m.memory[address] = value
}

func (m *Manager) GetByte(address byte) byte {
	return m.memory[address]
}

func (m *Manager) Reset() {
	var i uint16
	for i = 0; i < m.memorySize; i++ {
		m.memory[i] = 0xFF
	}
}

func (m *Manager) GetMemory() []byte {
	return m.memory
}

func (m *Manager) GetSprites() []Sprite {
	return m.sprites
}

func (m *Manager) GetCount() byte {
	return m.spritesSize
}

func (m *Manager) Copy(index byte, sprite Sprite) {
	var data [4]byte

	for i := 0; i < 4; i++ {
		data[i] = sprite.data[i]
	}

	j := 0

	for i := index * 4; i < (index+1)*4; i++ {
		m.memory[i] = data[j]
		j++
	}

	m.sprites[index] = NewSprite(sprite.id, sprite.index, m.memory[index*4:(index+1)*4])
}

func (m *Manager) GetSprite(index byte) Sprite {
	return m.sprites[index]
}
