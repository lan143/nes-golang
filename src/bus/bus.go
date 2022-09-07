package bus

type Interrupt byte

const (
	NMI Interrupt = 1
)

type Bus struct {
	cpuWrites  map[uint16]func(value byte)
	cpuReads   map[uint16]func() byte
	interrupts map[Interrupt]func()
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

func NewBus() *Bus {
	return &Bus{}
}
