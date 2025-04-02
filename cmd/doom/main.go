package main

import (
	"github.com/ianwjhalliday/go-doom/internal/cmd"
	"github.com/ianwjhalliday/go-doom/internal/engine"
	"github.com/ianwjhalliday/go-doom/internal/game"
	"runtime"
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

	// TODO: All the game package functions can become methods on the Game struct
	g := game.New(&game.Config)

	engine.GameLoop(g)
}
