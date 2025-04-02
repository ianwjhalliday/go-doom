package engine

type Game interface {
	Display()
	Shutdown()
	Tick()
}
