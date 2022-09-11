package bus

type Interrupt byte

const (
	NMI Interrupt = 1
)

type JoyPadButton byte

const (
	JoyPadButtonUp     JoyPadButton = 1
	JoyPadButtonDown                = 2
	JoyPadButtonLeft                = 3
	JoyPadButtonRight               = 4
	JoyPadButtonA                   = 5
	JoyPadButtonB                   = 6
	JoyPadButtonSelect              = 7
	JoyPadButtonStart               = 8
)

type Bus struct {
	cpuWrites  map[uint16]func(value byte)
	cpuReads   map[uint16]func() byte
	interrupts map[Interrupt]func()
	keyEvent   func(button JoyPadButton, pressed bool)
}

func (b *Bus) Init() {
	b.cpuWrites = make(map[uint16]func(value byte))
	b.cpuReads = make(map[uint16]func() byte)
	b.interrupts = make(map[Interrupt]func())
}

func (b *Bus) WriteByCPU(address uint16, value byte) {
	if fn, ok := b.cpuWrites[address]; ok {
		fn(value)
	}
}

func (b *Bus) OnCPUWrite(address uint16, fn func(value byte)) {
	b.cpuWrites[address] = fn
}

func (b *Bus) ReadByCPU(address uint16) byte {
	if fn, ok := b.cpuReads[address]; ok {
		return fn()
	}

	return 0
}

func (b *Bus) OnCPURead(address uint16, fn func() byte) {
	b.cpuReads[address] = fn
}

func (b *Bus) Interrupt(interrupt Interrupt) {
	if fn, ok := b.interrupts[interrupt]; ok {
		fn()
	}
}

func (b *Bus) OnInterrupt(interrupt Interrupt, fn func()) {
	b.interrupts[interrupt] = fn
}

func (b *Bus) KeyEvent(button JoyPadButton, pressed bool) {
	if b.keyEvent != nil {
		b.keyEvent(button, pressed)
	}
}

func (b *Bus) OnKeyEvent(fn func(button JoyPadButton, pressed bool)) {
	b.keyEvent = fn
}

func NewBus() *Bus {
	return &Bus{}
}
