package engine

type Game interface {
	Display()
	HandleEvent(ev *Event)
	Shutdown()
	Tick()
}
