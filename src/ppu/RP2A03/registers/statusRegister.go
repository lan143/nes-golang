package registers

type PPUStatusFlag byte

const (
	RecordInVRAMAllowed PPUStatusFlag = 0x10 // 1 – Запись в видеопамять разрешена; 0 - запрещена
	SpriteOverflow                    = 0x20 // 1 – На линии больше 8-и спрайтов; 0 - меньше
	SpriteZeroHit                     = 0x40 // Устанавливается в 1 после вывода спрайта с номером 0. Сбрасывается при чтении или при кадровом синхроимпульсе.
	VBlank                            = 0x80 // 1 – PPU генерирует обратный кадровый импульс; 0 – PPU рисует картинку на экране. Сбрасывается при чтении.
)

type PPUStatusRegister struct {
	value PPUStatusFlag
}

func (r *PPUStatusRegister) GetValue() byte {
	return byte(r.value)
}

func (r *PPUStatusRegister) SetValue(value byte) {
	r.value = PPUStatusFlag(value)
}

func (r *PPUStatusRegister) IsSpriteZeroHit() bool {
	return r.isFlag(SpriteZeroHit)
}

func (r *PPUStatusRegister) IsSpriteOverflow() bool {
	return r.isFlag(SpriteOverflow)
}

func (r *PPUStatusRegister) IsVBlank() bool {
	return r.isFlag(VBlank)
}

func (r *PPUStatusRegister) IsRecordInVRAMAllowed() bool {
	return r.isFlag(RecordInVRAMAllowed)
}

func (r *PPUStatusRegister) SetSpriteZeroHit() {
	r.setFlag(SpriteZeroHit)
}

func (r *PPUStatusRegister) SetSpriteOverflow() {
	r.setFlag(SpriteOverflow)
}

func (r *PPUStatusRegister) SetVBlank() {
	r.setFlag(VBlank)
}

func (r *PPUStatusRegister) SetRecordInVRAMAllowed() {
	r.setFlag(RecordInVRAMAllowed)
}

func (r *PPUStatusRegister) ClearSpriteZeroHit() {
	r.clearFlag(SpriteZeroHit)
}

func (r *PPUStatusRegister) ClearSpriteOverflow() {
	r.clearFlag(SpriteOverflow)
}

func (r *PPUStatusRegister) ClearVBlank() {
	r.clearFlag(VBlank)
}

func (r *PPUStatusRegister) ClearRecordInVRAMAllowed() {
	r.clearFlag(RecordInVRAMAllowed)
}

func (r *PPUStatusRegister) setFlag(flag PPUStatusFlag) {
	r.value |= flag
}

func (r *PPUStatusRegister) isFlag(flag PPUStatusFlag) bool {
	return r.value&flag > 0
}

func (r *PPUStatusRegister) clearFlag(flag PPUStatusFlag) {
	r.value &= ^flag
}
