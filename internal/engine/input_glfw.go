package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"unicode"
)

func pauseMouseCaptureGLFW() {
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	window.SetMouseButtonCallback(nil)
	window.SetCursorPosCallback(nil)
}

func resumeMouseCaptureGLFW() {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetMouseButtonCallback(glfwMouseButtonCallback)
	window.SetCursorPosCallback(glfwCursorPosCallback)
}

func glfwKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	eventType := KeyUp
	if action == glfw.Press || action == glfw.Repeat {
		eventType = KeyDown
	}

	PostEvent(&Event{
		Type:  eventType,
		Data1: translateKey(key),
	})
}

func glfwMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	// TODO
}

func glfwCursorPosCallback(w *glfw.Window, xPos float64, yPos float64) {
	// TODO
}

func translateKey(glfwKey glfw.Key) int {
	switch glfwKey {
	// Function keys
	case glfw.KeyEscape:
		return KeyEscape
	case glfw.KeyEnter:
		return KeyEnter
	case glfw.KeyTab:
		return KeyTab
	case glfw.KeyBackspace:
		return KeyBackspace
	case glfw.KeyInsert:
		return 0 // Unsupported key
	case glfw.KeyDelete:
		return 0 // Unsupported key
	case glfw.KeyRight:
		return KeyRightArrow
	case glfw.KeyLeft:
		return KeyLeftArrow
	case glfw.KeyDown:
		return KeyDownArrow
	case glfw.KeyUp:
		return KeyUpArrow
	case glfw.KeyPageUp:
		return 0 // Unsupported key
	case glfw.KeyPageDown:
		return 0 // Unsupported key
	case glfw.KeyHome:
		return 0 // Unsupported key
	case glfw.KeyEnd:
		return 0 // Unsupported key
	case glfw.KeyCapsLock:
		return 0 // Unsupported key
	case glfw.KeyScrollLock:
		return 0 // Unsupported key
	case glfw.KeyNumLock:
		return 0 // Unsupported key
	case glfw.KeyPrintScreen:
		return 0 // Unsupported key
	case glfw.KeyPause:
		return KeyPause
	case glfw.KeyF1:
		return KeyF1
	case glfw.KeyF2:
		return KeyF2
	case glfw.KeyF3:
		return KeyF3
	case glfw.KeyF4:
		return KeyF4
	case glfw.KeyF5:
		return KeyF5
	case glfw.KeyF6:
		return KeyF6
	case glfw.KeyF7:
		return KeyF7
	case glfw.KeyF8:
		return KeyF8
	case glfw.KeyF9:
		return KeyF9
	case glfw.KeyF10:
		return KeyF10
	case glfw.KeyF11:
		return KeyF11
	case glfw.KeyF12:
		return KeyF12
	case glfw.KeyF13:
		return 0 // Unsupported key
	case glfw.KeyF14:
		return 0 // Unsupported key
	case glfw.KeyF15:
		return 0 // Unsupported key
	case glfw.KeyF16:
		return 0 // Unsupported key
	case glfw.KeyF17:
		return 0 // Unsupported key
	case glfw.KeyF18:
		return 0 // Unsupported key
	case glfw.KeyF19:
		return 0 // Unsupported key
	case glfw.KeyF20:
		return 0 // Unsupported key
	case glfw.KeyF21:
		return 0 // Unsupported key
	case glfw.KeyF22:
		return 0 // Unsupported key
	case glfw.KeyF23:
		return 0 // Unsupported key
	case glfw.KeyF24:
		return 0 // Unsupported key
	case glfw.KeyF25:
		return 0 // Unsupported key
	case glfw.KeyKP0:
		return '0'
	case glfw.KeyKP1:
		return '1'
	case glfw.KeyKP2:
		return '2'
	case glfw.KeyKP3:
		return '3'
	case glfw.KeyKP4:
		return '4'
	case glfw.KeyKP5:
		return '5'
	case glfw.KeyKP6:
		return '6'
	case glfw.KeyKP7:
		return '7'
	case glfw.KeyKP8:
		return '8'
	case glfw.KeyKP9:
		return '9'
	case glfw.KeyKPDecimal:
		return '.'
	case glfw.KeyKPDivide:
		return '/'
	case glfw.KeyKPMultiply:
		return '*'
	case glfw.KeyKPSubtract:
		return KeyMinus
	case glfw.KeyKPAdd:
		return '+'
	case glfw.KeyKPEnter:
		return '\n'
	case glfw.KeyKPEqual:
		return KeyEquals
	case glfw.KeyLeftShift:
		return KeyRShift // Doom uses right shift key code for both shifts
	case glfw.KeyLeftControl:
		return KeyRCtrl // Doom uses right control key code for both controls
	case glfw.KeyLeftAlt:
		return KeyRAlt // Doom uses right alt key code for both alts
	case glfw.KeyLeftSuper:
		return 0 // Unsupported key
	case glfw.KeyRightShift:
		return KeyRShift
	case glfw.KeyRightControl:
		return KeyRCtrl
	case glfw.KeyRightAlt:
		return KeyLAlt
	case glfw.KeyRightSuper:
		return 0 // Unsupported key
	case glfw.KeyMenu:
		return 0 // Unsupported key

		// Printable keys -- Doom uses the same value
	case
		glfw.KeySpace,
		glfw.KeyApostrophe,
		glfw.KeyComma,
		glfw.KeyMinus,
		glfw.KeyPeriod,
		glfw.KeySlash,
		glfw.Key0,
		glfw.Key1,
		glfw.Key2,
		glfw.Key3,
		glfw.Key4,
		glfw.Key5,
		glfw.Key6,
		glfw.Key7,
		glfw.Key8,
		glfw.Key9,
		glfw.KeySemicolon,
		glfw.KeyEqual,
		glfw.KeyA,
		glfw.KeyB,
		glfw.KeyC,
		glfw.KeyD,
		glfw.KeyE,
		glfw.KeyF,
		glfw.KeyG,
		glfw.KeyH,
		glfw.KeyI,
		glfw.KeyJ,
		glfw.KeyK,
		glfw.KeyL,
		glfw.KeyM,
		glfw.KeyN,
		glfw.KeyO,
		glfw.KeyP,
		glfw.KeyQ,
		glfw.KeyR,
		glfw.KeyS,
		glfw.KeyT,
		glfw.KeyU,
		glfw.KeyV,
		glfw.KeyW,
		glfw.KeyX,
		glfw.KeyY,
		glfw.KeyZ,
		glfw.KeyLeftBracket,
		glfw.KeyBackslash,
		glfw.KeyRightBracket,
		glfw.KeyGraveAccent,
		glfw.KeyWorld1,
		glfw.KeyWorld2:
		return int(unicode.ToLower(rune(glfwKey)))
	}
	return 0 // Unknown or unsupported key
}
