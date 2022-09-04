package registers

type PPUMaskFlag byte

const (
	Greyscale                 PPUMaskFlag = 0x1  // Тип дисплея: Color/Monochrome (в Денди не используется)
	LeftmostBackgroundVisible             = 0x2  // 0 – Рисунок фона невиден в крайнем левом столбце; 1- Весь фон виден
	LeftmostSpritesVisible                = 0x4  // 0 – Спрайты невидны в крайнем левом столбце; 1- Все спрайты видны
	BackgroundVisible                     = 0x8  // 0 – Фон не отображается; 1 – Фон отображается
	SpritesVisible                        = 0x10 // 0 – Спрайты не отображаются; 1 – Спрайты отображаются
	EmphasizeRed                          = 0x20 // Яркость экрана/интенсивность цвета в RGB (в Денди не используется)
	EmphasizeGreen                        = 0x40 // Яркость экрана/интенсивность цвета в RGB (в Денди не используется)
	EmphasizeBlue                         = 0x80 // Яркость экрана/интенсивность цвета в RGB (в Денди не используется)
)

type PPUMaskRegister struct {
	value PPUMaskFlag
}

func (r *PPUMaskRegister) SetValue(value byte) {
	r.value = PPUMaskFlag(value)
}

func (r *PPUMaskRegister) IsBackgroundVisible() bool {
	return r.isFlag(BackgroundVisible)
}

func (r *PPUMaskRegister) IsSpritesVisible() bool {
	return r.isFlag(SpritesVisible)
}

func (r *PPUMaskRegister) IsGreyscale() bool {
	return r.isFlag(Greyscale)
}

func (r *PPUMaskRegister) IsEmphasizeRed() bool {
	return r.isFlag(EmphasizeRed)
}

func (r *PPUMaskRegister) IsEmphasizeGreen() bool {
	return r.isFlag(EmphasizeGreen)
}

func (r *PPUMaskRegister) IsEmphasizeBlue() bool {
	return r.isFlag(EmphasizeBlue)
}

func (r *PPUMaskRegister) setFlag(flag PPUMaskFlag) {
	r.value |= flag
}

func (r *PPUMaskRegister) isFlag(flag PPUMaskFlag) bool {
	return r.value&flag > 0
}

func (r *PPUMaskRegister) clearFlag(flag PPUMaskFlag) {
	r.value &= ^flag
}
