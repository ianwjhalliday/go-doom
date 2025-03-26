package engine

import (
	"github.com/ianwjhalliday/go-doom/internal/engine/glfwimpl"
)

// GameLoop is a placeholder until the proper D_DoomLoop with tic updates etc is implemented.
func GameLoop() {
	// Temporarily show title screen to prove OpenGL rendering works
	pal := LumpByName[[3 * 256]byte]("PLAYPAL")
	SetPalette(pal, 0)

	titleScreen := LumpByName[Patch]("TITLEPIC")
	DrawPatch(0, 0, 0, titleScreen)

	for !glfwimpl.WindowShouldClose() {
		glfwimpl.FinishUpdate(&screens[0])
		glfwimpl.PollInputEvents()
	}
}
