package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ianwjhalliday/go-doom/internal/cmd"
	"github.com/ianwjhalliday/go-doom/internal/engine"
	"github.com/ianwjhalliday/go-doom/internal/game"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	args := cmd.InitArguments()
	game.InitConfig(args)
	game.InitSettings(game.Config.ConfigFile)
	defer game.SaveSettings(game.Config.ConfigFile) // TODO: This won't save settings when we use os.Exit() to exit the program

	engine.InitFiles(game.Config.WadFiles)
	game.VetFiles(&game.Config)
	game.PrintBanner(&game.Config)

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
