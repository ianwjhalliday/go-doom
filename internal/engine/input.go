package engine

const (
	KeyRightArrow = 0xae
	KeyLeftArrow  = 0xac
	KeyUpArrow    = 0xad
	KeyDownArrow  = 0xaf
	KeyEscape     = 27
	KeyEnter      = 13
	KeyTab        = 9
	KeyF1         = 0x80 + 0x3b
	KeyF2         = 0x80 + 0x3c
	KeyF3         = 0x80 + 0x3d
	KeyF4         = 0x80 + 0x3e
	KeyF5         = 0x80 + 0x3f
	KeyF6         = 0x80 + 0x40
	KeyF7         = 0x80 + 0x41
	KeyF8         = 0x80 + 0x42
	KeyF9         = 0x80 + 0x43
	KeyF10        = 0x80 + 0x44
	KeyF11        = 0x80 + 0x57
	KeyF12        = 0x80 + 0x58
	KeyBackspace  = 127
	KeyPause      = 0xff
	KeyEquals     = 0x3d
	KeyMinus      = 0x2d
	KeyRShift     = 0x80 + 0x36
	KeyRCtrl      = 0x80 + 0x1d
	KeyRAlt       = 0x80 + 0x38
	KeyLAlt       = KeyRAlt
)

var events [64]Event
var eventHead int
var eventTail int

type EventType int

const (
	KeyDown EventType = iota
	KeyUp
	Mouse
	Joystick
)

type Event struct {
	Type  EventType
	Data1 int // keys / mouse/joystick buttons
	Data2 int // mouse/joystick x move
	Data3 int // mouse/joystick y move
}

func PostEvent(ev *Event) {
	events[eventHead] = *ev
	eventHead = (1 + eventHead) % len(events)
}

type EventHandler interface {
	HandleEvent(ev *Event)
}

func ProcessEvents(eh EventHandler) {
	for eventHead != eventTail {
		ev := &events[eventTail]
		eventTail = (1 + eventTail) % len(events)
		eh.HandleEvent(ev)
	}
}

func PauseMouseCapture() {
	pauseMouseCaptureGLFW()
}

func ResumeMouseCapture() {
	resumeMouseCaptureGLFW()
}
