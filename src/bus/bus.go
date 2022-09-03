package bus

type Event byte

const (
	VBlink Event = 1
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
