package game

import "github.com/ianwjhalliday/go-doom/internal/engine"

type Menu struct {
	open bool
}

func (g *Game) openMenu() {
	g.menu.open = true
	// TODO: Pause mouse capture
	// TODO: Set which menu is active
	engine.PlaySound("swtchn")
}

func (g *Game) closeMenu() {
	g.menu.open = false
	engine.PlaySound("swtchx")
}

func (g *Game) eventMenu(ev *engine.Event) bool {
	var ch int
	if ev.Type == engine.KeyDown {
		ch = ev.Data1
	}

	// When menu isn't open the only response is to open the menu if ESC is pressed
	if !g.menu.open {
		if ch == engine.KeyEscape {
			g.openMenu()
			return true
		}
		return false
	}

	// Otherwise, navigate the menus
	if ch == engine.KeyEscape {
		g.closeMenu()
	}
	return false
}

func (g *Game) displayMenu() {
	if !g.menu.open {
		return
	}

	engine.DrawPatch(94, 2, 0, engine.LumpByName[engine.Patch]("M_DOOM"))
}
