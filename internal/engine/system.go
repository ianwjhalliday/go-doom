package engine

import (
	"log/slog"
	"os"
)

func FatalError(msg string, args ...any) {
	slog.Error(msg, args...)
	Terminate(1)
}

func Terminate(code int) {
	slog.Debug("Terminate", "code", code)
	ShutdownGraphics()
	os.Exit(code)
}

// GameLoop is a placeholder until the proper D_DoomLoop with tic updates etc is implemented.
func GameLoop() {
	// Temporarily show title screen to prove OpenGL rendering works
	pal := LumpByName[[3 * 256]byte]("PLAYPAL")
	SetPalette(pal, 0)

	titleScreen := LumpByName[Patch]("TITLEPIC")
	DrawPatch(0, 0, 0, titleScreen)

	for !WindowShouldClose() {
		FinishUpdate(&screens[0])
		PollInputEvents()
	}
}
