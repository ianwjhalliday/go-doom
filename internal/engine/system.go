package engine

import (
	"log/slog"
	"os"
	"time"
)

func Quit() {
	// TODO: Game needs some callback hook here to save game config to file before quitting
	game.Shutdown()
	ShutdownSound()
	Terminate(0)
}

func FatalError(msg string, args ...any) {
	slog.Error(msg, args...)
	Terminate(1)
}

func Terminate(code int) {
	slog.Debug("Terminate", "code", code)
	ShutdownGraphics()
	os.Exit(code)
}

var game Game

// GameLoop is a placeholder until the proper D_DoomLoop with tic updates etc is implemented.
func GameLoop(g Game) {
	// TODO: Create an engine struct to hold game interface and other engine state
	game = g

	// Temporarily show title screen to prove OpenGL rendering works
	// TODO: Where does the palette actually get loaded
	pal := LumpByName[[3 * 256]byte]("PLAYPAL")
	SetPalette(pal, 0)

	var lastTicks int

	for {
		// non-timedemo loop:
		nextTicks := TryRunTicks(lastTicks)

		if nextTicks != lastTicks {
			lastTicks = nextTicks
			// S_UpdateSounds() -- update sound channels
			Display()
		} else {
			// TODO: Should we actually sleep? OS can hold us up for a long while.
			//  If so, how long? Should we modify GetTicks() to identify time remaining
			//  to the next tick to make a smart request for sleep duration?
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// ticRate of 35 ticks per second per original Doom
const ticRate = 35

var baseTime = time.Now()
var gametick int

// GetTicks returns the number of tics since the game started.
func GetTicks() int {
	t := time.Now()
	d := t.Sub(baseTime)
	return int(d.Microseconds() * ticRate / 1000000)
}

func TryRunTicks(lastTicks int) (enterTicks int) {
	enterTicks = GetTicks()
	realTicks := enterTicks - lastTicks

	if realTicks == 0 {
		// Not enough time has passed since last tick
		return
	}

	if realTicks > 1 {
		slog.Debug("Missed tick timing", "realTicks", realTicks)
	}

	StartTick()
	ProcessEvents(game)
	//BuildTiccmd()

	game.Tick()

	gametick += 1


	return
}

func StartTick() {
	if WindowShouldClose() {
		Quit()
	}
		PollInputEvents()
	}

func Display() {
	game.Display()
	FinishUpdate(&screens[0])
}
