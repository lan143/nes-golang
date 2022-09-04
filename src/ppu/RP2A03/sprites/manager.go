package sprites

import "log"

type Manager struct {
	sprites []Sprite
	memory  []byte
}

func (m *Manager) Init(memory []byte) {
	m.memory = memory
	var i uint16
	spritesLen := uint16(len(m.memory)) / 4

	for i = 0; i < spritesLen; i++ {
		m.sprites = append(m.sprites, NewSprite(i, i, m.memory[i*4:(i+1)*4]))
	}

	log.Printf("Sprites len: %d", len(m.sprites))
}

func (m *Manager) GetCount() int {
	return len(m.sprites)
}

func (m *Manager) Copy(index int, sprite Sprite) {
	m.sprites[index].Copy(sprite)
}

func (m *Manager) GetSprite(index int) Sprite {
	return m.sprites[index]
}