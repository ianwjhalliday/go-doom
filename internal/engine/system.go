package engine

import "github.com/ianwjhalliday/go-doom/internal/engine/system"

func InitGraphics() {
	system.InitGraphics()
}

func ShutdownGraphics() {
	system.ShutdownGraphics()
}

func GameLoop() {
	system.TempGlfwLoop()
}
