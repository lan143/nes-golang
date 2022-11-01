package bus

type Interrupt byte

const (
	NMI Interrupt = 1
	IRQ           = 2
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
	cpuWrites      map[uint16]func(value byte)
	cpuReads       map[uint16]func() byte
	interrupts     map[Interrupt]func()
	keyEvent       func(button JoyPadButton, pressed bool)
	readFromCPU    func(address uint16) byte
	apuDMCActivate func()
	cpuSkipCycles  func(cycles uint16)
	ppuScanline    func()
}

func (b *Bus) Init() {
	b.cpuWrites = make(map[uint16]func(value byte))
	b.cpuReads = make(map[uint16]func() byte)
	b.interrupts = make(map[Interrupt]func())
}

func (b *Bus) DrivePPUScanline() {
	if b.ppuScanline != nil {
		b.ppuScanline()
	}
}

func (b *Bus) OnPPUScanline(fn func()) {
	b.ppuScanline = fn
}

func (b *Bus) CPUSkipCycles(cycles uint16) {
	b.cpuSkipCycles(cycles)
}

func (b *Bus) OnCPUSkipCycles(fn func(cycles uint16)) {
	b.cpuSkipCycles = fn
}

func (b *Bus) ApuDMCActivate() {
	b.apuDMCActivate()
}

func (b *Bus) OnApuDMCActivate(fn func()) {
	b.apuDMCActivate = fn
}

func (b *Bus) ReadFromCPU(address uint16) byte {
	if b.readFromCPU != nil {
		return b.readFromCPU(address)
	}

	return 0
}

func (b *Bus) OnReadFromCPU(fn func(address uint16) byte) {
	b.readFromCPU = fn
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
