package bus

type Event byte

const (
	NMIInterrupt Event = 1
	Write2000          = 2
	Write2001          = 3
	Write2003          = 3
	Write2005          = 4
	Write2006          = 5
	Write2007          = 6
	Write4014          = 7
)

type Bus struct {
	events map[Event][]func()
}

func (b *Bus) Init() {
	b.events = make(map[Event][]func())
}

func (b *Bus) PushEvent(event Event) {
	if events, ok := b.events[event]; ok {
		for _, e := range events {
			e()
		}
	}
}

func (b *Bus) Subscribe(event Event, fn func()) {
	b.events[event] = append(b.events[event], fn)
}

func NewBus() *Bus {
	return &Bus{}
}
