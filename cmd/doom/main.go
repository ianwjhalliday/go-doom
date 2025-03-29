package main

import (
	"runtime"

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

engine.InitSound()
	defer engine.ShutdownSound()
	engine.InitGraphics()
	defer engine.ShutdownGraphics()

	engine.GameLoop()
}
