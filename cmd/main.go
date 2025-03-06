package main

import (
	"log"
	"runtime"
)
import "github.com/go-gl/glfw/v3.3/glfw"

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	InitArgs()

	err := glfw.Init()
	if err != nil {
		log.Fatalf("FATAL Failed to initialize GLFW: %s", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		log.Fatalf("FATAL Failed to create window: %s", err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
