package sprites

type Sprite struct {
	id    byte
	index byte
	data  []byte
}

func (s *Sprite) GetId() byte {
	return s.id
}

func (s *Sprite) GetIndex() byte {
	return s.index
}

func (s *Sprite) GetData() []byte {
	return s.data
}

func (s *Sprite) IsEmpty() bool {
	return s.data[0] == 0xFF && s.data[1] == 0xFF && s.data[2] == 0xFF && s.data[3] == 0xFF
}

func (s *Sprite) IsVisible() bool {
	return s.data[0] < 0xEF
}

func (s *Sprite) GetXPosition() byte {
	return s.data[3]
}

func (s *Sprite) GetYPosition() byte {
	return s.data[0] - 1
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

func (s *Sprite) On(y uint16, length byte) bool {
	return (y >= uint16(s.GetYPosition())) && (y < uint16(s.GetYPosition()+length))
}

func NewSprite(id, index byte, data []byte) Sprite {
	return Sprite{
		id:    id,
		index: index,
		data:  data,
	}
}
