package game

import (
	"github.com/ianwjhalliday/go-doom/internal/engine"
	"math/rand"
	"time"
)

type GameMenu struct {
	open             bool
	curMenu          *menu
	whichSkull       int
	skullAnimCounter int
}

type menu struct {
	items    []menuItem
	curItem  int
	prevMenu *menu
	drawer   func()
	x, y     int
}

type menuItemKind int

const (
	menuItemKindNoAction menuItemKind = iota
	menuItemKindAction
	menuItemKindSlider
	menuItemKindSpacer menuItemKind = -1
)

type menuItem struct {
	kind      menuItemKind
	patchName string
	handler   func(g *Game, choice int)
	alphaKey  int
}

var mainMenuItems = []menuItem{
	{menuItemKindAction, "M_NGAME", menuActionTODO /*M_NewGame*/, 'n'},
	{menuItemKindAction, "M_OPTION", menuActionTODO /*M_Options*/, 'o'},
	{menuItemKindAction, "M_LOADG", menuActionTODO /*M_LoadGame*/, 'l'},
	{menuItemKindAction, "M_SAVEG", menuActionTODO /*M_SaveGame*/, 's'},
	{menuItemKindAction, "M_RDTHIS", menuActionTODO /*M_ReadThis*/, 'r'},
	{menuItemKindAction, "M_QUITG", menuActionQuit, 'q'},
}

var mainMenu = menu{
	items:  mainMenuItems,
	drawer: drawMainMenu,
	x:      97, y: 64,
}

func (g *Game) openMenu() {
	g.menu.open = true
	// TODO: Pause mouse capture
	// TODO: Set which menu is active
	g.menu.curMenu = &mainMenu
	engine.PlaySound("swtchn")
}

func (g *Game) closeMenu() {
	g.menu.open = false
	// TODO: Resume mouse capture
	engine.PlaySound("swtchx")
}

func (g *Game) eventMenu(ev *engine.Event) bool {
	ch := -1
	if ev.Type == engine.KeyDown {
		ch = ev.Data1
	}
	// TODO: ch from other event types
	if ch == -1 {
		return false
	}

	// TODO: handle save game slot name input
	// TODO: handle message prompt input

	// TODO: F-keys and global hotkeys like -/+

	// When menu isn't open the only response is to open the menu if ESC is pressed
	if !g.menu.open {
		if ch == engine.KeyEscape {
			g.openMenu()
			return true
		}
		return false
	}

	// Otherwise, navigate the menus
	switch ch {
	case engine.KeyDownArrow:
		g.selectMenuItem(1)
	case engine.KeyUpArrow:
		g.selectMenuItem(-1)
	case engine.KeyLeftArrow:
		g.slideMenuItem(-1)
	case engine.KeyRightArrow:
		g.slideMenuItem(1)
	case engine.KeyEnter:
		g.activateMenuItem()
	case engine.KeyEscape:
		g.closeMenu()
	case engine.KeyBackspace:
		g.previousMenu()
	default:
		g.jumpToMenuItem(ch)
	}

	// always eat the event if the menu was open
	return true
}

func (g *Game) activateMenuItem() {
	mi := g.menu.curMenu.items[g.menu.curMenu.curItem]
	switch mi.kind {
	case menuItemKindAction:
		mi.handler(g, g.menu.curMenu.curItem)
		engine.PlaySound("pistol")
	case menuItemKindSlider:
		// Behave like right arrow
		mi.handler(g, 1)
		engine.PlaySound("stnmov")
	default:
		// ignore
	}
}

func (g *Game) jumpToMenuItem(ch int) {
	for i, item := range g.menu.curMenu.items {
		if ch == item.alphaKey && i != g.menu.curMenu.curItem {
			g.menu.curMenu.curItem = i
			engine.PlaySound("pstop")
		}
	}
}

func (g *Game) previousMenu() {
	if g.menu.curMenu.prevMenu != nil {
		g.menu.curMenu = g.menu.curMenu.prevMenu
		engine.PlaySound("swtchn")
	}
}

func (g *Game) selectMenuItem(dir int) {
	for {
		g.menu.curMenu.curItem += dir + len(g.menu.curMenu.items)
		g.menu.curMenu.curItem %= len(g.menu.curMenu.items)
		if g.menu.curMenu.items[g.menu.curMenu.curItem].kind != menuItemKindSpacer {
			break
		}
	}
	engine.PlaySound("pstop")
}

func (g *Game) slideMenuItem(dir int) {
	mi := g.menu.curMenu.items[g.menu.curMenu.curItem]
	if mi.kind == menuItemKindSlider {
		mi.handler(g, dir)
		engine.PlaySound("stnmov")
	}
}

func (g *Game) tickMenu() {
	g.menu.skullAnimCounter--
	if g.menu.skullAnimCounter < 0 {
		g.menu.skullAnimCounter = 8
		g.menu.whichSkull ^= 1
	}
}

const lineHeight = 16

func (g *Game) displayMenu() {
	if !g.menu.open {
		return
	}

	// draw menu title
	g.menu.curMenu.drawer()

	// draw menu items
	x := g.menu.curMenu.x
	y := g.menu.curMenu.y
	for _, item := range g.menu.curMenu.items {
		if item.patchName != "" {
			engine.DrawPatch(x, y, 0, engine.LumpByName[engine.Patch](item.patchName))
		}
		y += lineHeight
	}

	// draw skull cursor
	skullX := x - 32
	skullY := g.menu.curMenu.y - 5 + g.menu.curMenu.curItem*lineHeight
	skullPatch := "M_SKULL1"
	if g.menu.whichSkull == 1 {
		skullPatch = "M_SKULL2"
	}
	engine.DrawPatch(skullX, skullY, 0, engine.LumpByName[engine.Patch](skullPatch))
}

func drawMainMenu() {
	engine.DrawPatch(94, 2, 0, engine.LumpByName[engine.Patch]("M_DOOM"))
}

func menuActionTODO(g *Game, choice int) {}

func menuActionQuit(g *Game, choice int) {
	// TODO: This is actually the 'y' response to the quit menu "are you sure?" message
	//  Still need to implement the message prompt.

	quitSounds := []string{
		"pldeth",
		"dmpain",
		"popain",
		"slop",
		"telept",
		"posit1",
		"posit3",
		"sgtatk",
	}

	if g.config.GameMode == GameModeCommercial {
		quitSounds = []string{
			"vilact",
			"getpow",
			"boscub",
			"slop",
			"skeswg",
			"kntdth",
			"bspact",
			"sgtatk",
		}
	}

	engine.PlaySound(quitSounds[rand.Intn(len(quitSounds))])
	time.Sleep(1500 * time.Millisecond)

	engine.Quit()
}
