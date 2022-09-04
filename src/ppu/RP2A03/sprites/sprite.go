package sprites

type Sprite struct {
	id    uint16
	index uint16
	data  []byte
}

func (s *Sprite) GetId() uint16 {
	return s.id
}

func (s *Sprite) GetIndex() uint16 {
	return s.index
}

func (s *Sprite) GetData() []byte {
	return s.data
}

func (s *Sprite) Copy(sprite Sprite) {
	s.id = sprite.id
	s.data = sprite.data
}

func (s *Sprite) IsEmpty() bool {
	return s.data[0] == 0xFF && s.data[1] == 0xFF && s.data[2] == 0xFF && s.data[3] == 0xFF
}

func (s *Sprite) IsVisible() bool {
	return s.data[0] < 0xEF
}

func (s *Sprite) GetYPosition() byte {
	return s.data[0] - 1
}

func (s *Sprite) GetXPosition() byte {
	return s.data[3]
}

func (s *Sprite) GetTileIndex() byte {
	return s.data[1]
}

func (s *Sprite) GetTileIndexForSize16() uint16 {
	return ((uint16(s.data[1]) & 1) * 0x1000) + (uint16(s.data[1])>>1)*0x20
}

func (s *Sprite) GetPalletNum() byte {
	return s.data[2] & 0x3
}

func (s *Sprite) GetPriority() byte {
	return (s.data[2] >> 5) & 1
}

func (s *Sprite) DoFlipHorizontally() bool {
	if (s.data[2]>>6)&1 > 0 {
		return true
	} else {
		return false
	}
}

func (s *Sprite) DoFlipVertically() bool {
	if (s.data[2]>>7)&1 > 0 {
		return true
	} else {
		return false
	}
}

func (s *Sprite) On(y byte, length byte) bool {
	return (y >= s.GetYPosition()) && (y < s.GetYPosition()+length)
}

func NewSprite(id, index uint16, data []byte) Sprite {
	return Sprite{
		id:    id,
		index: index,
		data:  data,
	}
}
