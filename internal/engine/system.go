package engine

import (
	"github.com/ianwjhalliday/go-doom/internal/engine/system"
)

// GameLoop is a placeholder until the proper D_DoomLoop with tic updates etc is implemented.
func GameLoop() {
	var pal [3 * 256]uint8
	SetPalette(&pal, 0)

	for !system.WindowShouldClose() {
		system.FinishUpdate(&screens[0])
		system.PollInputEvents()
	}
}
