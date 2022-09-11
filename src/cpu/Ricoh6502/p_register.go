package Ricoh6502

type PFlag byte

const (
	C PFlag = 0x1  // перенос
	Z       = 0x2  // ноль
	I       = 0x4  // запрет внешних прерываний — IRQ (I=0 — прерывания разрешены)
	D       = 0x8  // режим BCD для инструкций сложения и вычитания с переносом;
	B       = 0x10 // обработка прерывания (B=1 после выполнения команды BRK);
	// V 0x20 не используется, равен 1;
	V = 0x40 // переполнение;
	N = 0x80 // знак. Равен старшему биту значения, загруженного в A, X или Y в результате выполнения операции (кроме TXS).
)

type PRegister struct {
	value byte
}

func (p *PRegister) Init() {
	p.value = 0x24
}

func (p *PRegister) GetValue() byte {
	return p.value
}

func (p *PRegister) SetValue(value byte) {
	value |= 0x20
	p.value = value
}

func (p *PRegister) IsC() bool {
	return p.isFlag(C)
}

func (p *PRegister) IsZ() bool {
	return p.isFlag(Z)
}

func (p *PRegister) IsI() bool {
	return p.isFlag(I)
}

func (p *PRegister) IsD() bool {
	return p.isFlag(D)
}

func (p *PRegister) IsB() bool {
	return p.isFlag(B)
}

func (p *PRegister) IsV() bool {
	return p.isFlag(V)
}

func (p *PRegister) IsN() bool {
	return p.isFlag(N)
}

func (p *PRegister) SetC() {
	p.setFlag(C)
}

func (p *PRegister) SetZ() {
	p.setFlag(Z)
}

func (p *PRegister) SetI() {
	p.setFlag(I)
}

func (p *PRegister) SetD() {
	p.setFlag(D)
}

func (p *PRegister) SetB() {
	p.setFlag(B)
}

func (p *PRegister) SetV() {
	p.setFlag(V)
}

func (p *PRegister) SetN() {
	p.setFlag(N)
}

func (p *PRegister) ClearC() {
	p.clearFlag(C)
}

func (p *PRegister) ClearZ() {
	p.clearFlag(Z)
}

func (p *PRegister) ClearI() {
	p.clearFlag(I)
}

func (p *PRegister) ClearD() {
	p.clearFlag(D)
}

func (p *PRegister) ClearB() {
	p.clearFlag(B)
}

func (p *PRegister) ClearV() {
	p.clearFlag(V)
}

func (p *PRegister) ClearN() {
	p.clearFlag(N)
}

func (p *PRegister) UpdateN(value byte) {
	if (value & 0x80) == 0 {
		p.ClearN()
	} else {
		p.SetN()
	}
}

func (p *PRegister) UpdateZ(value byte) {
	if value&0xFF == 0 {
		p.SetZ()
	} else {
		p.ClearZ()
	}
}

func (p *PRegister) UpdateC(value uint16) {
	if (value & 0x100) == 0 {
		p.ClearC()
	} else {
		p.SetC()
	}
}

func (p *PRegister) setFlag(flag PFlag) {
	p.value |= byte(flag)
}

func (p *PRegister) isFlag(flag PFlag) bool {
	return p.value&byte(flag) > 0
}

func (p *PRegister) clearFlag(flag PFlag) {
	p.value &= ^byte(flag)
}
