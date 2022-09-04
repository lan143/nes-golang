package registers

type PPUCtrlFlag byte

const (
	NameTableAddress       PPUCtrlFlag = 0x1  // Адрес активной экранной страницы (00 – $2000; 01 – $2400; 10 – $2800; 11 - $2C00)
	IncrementAddress                   = 0x4  // Выбор режима инкремента адреса при обращении к видеопамяти (0 – увеличение на единицу «горизонтальная запись»; 1 - увеличение на 32 «вертикальная запись»)
	SpritesPatternTable                = 0x8  // Выбор знакогенератора спрайтов (0/1)
	BackgroundPatternTable             = 0x10 // Выбор знакогенератора фона (0/1)
	SpritesSize                        = 0x20 // Размер спрайтов (0 - 8x8; 1 - 8x16)
	MasterSlave                        = 0x40 // Не используется (должен быть 0)
	NMIVBlank                          = 0x80 // Формирование запроса прерывания NMI при кадровом синхроимпульсе (0 - запрещено; 1 - разрешено)
)

type PPUCtrlRegister struct {
	value PPUCtrlFlag
}

func (r *PPUCtrlRegister) GetValue() byte {
	return byte(r.value)
}

func (r *PPUCtrlRegister) SetValue(value byte) {
	r.value = PPUCtrlFlag(value)
}

func (r *PPUCtrlRegister) IsNMIVBlank() bool {
	return r.isFlag(NMIVBlank)
}

func (r *PPUCtrlRegister) IsMasterSlave() bool {
	return r.isFlag(MasterSlave)
}

func (r *PPUCtrlRegister) IsSpritesSize() bool {
	return r.isFlag(SpritesSize)
}

func (r *PPUCtrlRegister) GetNameTableAddress() byte {
	return byte(r.value&NameTableAddress | r.value&(NameTableAddress+1))
}

func (r *PPUCtrlRegister) GetBackgroundPatternTable() byte {
	if r.isFlag(BackgroundPatternTable) {
		return 1
	} else {
		return 0
	}
}

func (r *PPUCtrlRegister) IsSpritesPatternTable() bool {
	return r.isFlag(SpritesPatternTable)
}

func (r *PPUCtrlRegister) IsIncrementAddress() bool {
	return r.isFlag(IncrementAddress)
}

func (r *PPUCtrlRegister) setFlag(flag PPUCtrlFlag) {
	r.value |= flag
}

func (r *PPUCtrlRegister) isFlag(flag PPUCtrlFlag) bool {
	return r.value&flag > 0
}

func (r *PPUCtrlRegister) clearFlag(flag PPUCtrlFlag) {
	r.value &= ^flag
}
