package sprites

type Manager struct {
	sprites []Sprite
	memory  []byte
}

func (m *Manager) Init(memory []byte) {
	m.memory = memory
	m.sprites = make([]Sprite, 0, len(m.memory)/4)
	var i uint16
	spritesLen := uint16(len(m.memory)) / 4

	for i = 0; i < spritesLen; i++ {
		m.sprites = append(m.sprites, NewSprite(i, i, m.memory[i*4:(i+1)*4]))
	}
}

func (m *Manager) GetMemory() []byte {
	return m.memory
}

func (m *Manager) GetSprites() []Sprite {
	return m.sprites
}

func (m *Manager) GetCount() int {
	return len(m.sprites)
}

func (m *Manager) Copy(index int, sprite Sprite) {
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

func (m *Manager) GetSprite(index int) Sprite {
	return m.sprites[index]
}
