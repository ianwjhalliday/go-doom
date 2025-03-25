package system

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log/slog"
	"os"
	"runtime"
)

const ScreenWidth = 320
const ScreenHeight = 200

var window *glfw.Window

func InitGraphics() {
	initGlfwWindow()
	initGLScreen()
}

func ShutdownGraphics() {
	glfw.Terminate()
}

func initGlfwWindow() {
	err := glfw.Init()
	if err != nil {
		slog.Error("InitGraphics: Could not initialize GLFW:", "err", err)
		os.Exit(1)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	window, err = glfw.CreateWindow(3*ScreenWidth, 3*ScreenWidth*3/4, "go-doom", nil, nil)
	if err != nil {
		slog.Error("InitGraphics: Could not create window", "err", err)
		os.Exit(1)
	}

	window.SetSizeLimits(ScreenWidth, ScreenWidth*3/4, glfw.DontCare, glfw.DontCare)
	window.SetAspectRatio(4, 3)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// TODO: Input callbacks
	//window.SetKeyCallback()
	//window.SetMouseButtonCallback()
	//window.SetCursorPosCallback()

	window.MakeContextCurrent()
}

func initGLScreen() {
	// TODO
}

// TempGlfwLoop is a placeholder until the proper D_DoomLoop with tic updates etc is implemented.
func TempGlfwLoop() {
	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
