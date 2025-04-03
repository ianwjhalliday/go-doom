package game

import "github.com/ianwjhalliday/go-doom/internal/engine"

type GameState int

const (
	GameStateLevel GameState = iota
	GameStateIntermission
	GameStateFinale
	GameStateTitleScreen
	GameStateForceWipe
	GameStateForceRedraw
)

type GameAction int

const (
	GameActionNothing GameAction = iota
	GameActionLoadLevel
	GameActionNewGame
	GameActionLoadGame
	GameActionSaveGame
	GameActionPlayDemo
	GameActionCompleted
	GameActionVictory
	GameActionWorldDone
	GameActionScreenshot
)

type Game struct {
	config *GameConfig
	state  GameState
	action GameAction
	menu   GameMenu

	userGame bool
	paused   bool

	titleAdvance  bool
	titleSequence int
	titlePageTics int
	titlePageName string
}

func New(config *GameConfig) *Game {
	return &Game{
		config:       config,
		state:        GameStateTitleScreen,
		action:       GameActionNothing,
		titleAdvance: true,
	}
}

func (g *Game) Display() {
	switch g.state {
	case GameStateLevel:
		// TODO
	case GameStateIntermission:
		// TODO
	case GameStateFinale:
		// TODO
	case GameStateTitleScreen:
		g.displayTitleScreen()
	case GameStateForceWipe, GameStateForceRedraw:
		// Nothing
	}

	// Menu always draws on top of everything else
	g.displayMenu()
}

func (g *Game) HandleEvent(ev *engine.Event) {
	if g.eventMenu(ev) {
		return // menu ate the event
	}
	g.eventGame(ev)
}

func (g *Game) eventGame(ev *engine.Event) bool {
	return false
}

func (g *Game) Shutdown() {
	SaveSettings(g.config.ConfigFile)
}

func (g *Game) Tick() {
	if g.titleAdvance {
		g.advanceTitleScreen()
	}

	g.tickMenu()
	g.tickGame()
}

func (g *Game) tickGame() {
	switch g.state {
	case GameStateLevel:
		// TODO
	case GameStateIntermission:
		// TODO
	case GameStateFinale:
		// TODO
	case GameStateTitleScreen:
		g.tickTitleScreen()
	case GameStateForceWipe, GameStateForceRedraw:
		// Nothing
	}
}

func (g *Game) advanceTitleScreen() {
	//g.players[g.consolePlayer].playerState = PlayerStateLive // TODO
	g.titleAdvance = false
	g.userGame = false
	g.paused = false
	g.action = GameActionNothing

	switch g.titleSequence {
	case 0:
		g.state = GameStateTitleScreen
		g.titlePageName = "titlepic"
		g.titlePageTics = 170
		music := "intro"
		if g.config.GameMode == GameModeCommercial {
			g.titlePageTics = 35 * 11
			music = "dm2ttl"
		}
		engine.StartMusic(music)
	case 1:
		g.titleAdvance = true // TODO: "demo1"
	case 2:
		g.state = GameStateTitleScreen
		g.titlePageTics = 200
		g.titlePageName = "credit"
	case 3:
		g.titleAdvance = true // TODO: "demo2"
	case 4:
		g.state = GameStateTitleScreen
		g.titlePageTics = 200
		g.titlePageName = "help2"
		if g.config.GameMode == GameModeRetail {
			g.titlePageName = "credit"
		} else if g.config.GameMode == GameModeCommercial {
			g.titlePageTics = 35 * 11
			g.titlePageName = "titlepic"
			engine.StartMusic("dm2ttl")
		}
	case 5:
		g.titleAdvance = true // TODO: "demo3"
	case 6:
		g.titleAdvance = true // TODO: "demo4"
	}

	g.titleSequence += 1
	if g.config.GameMode == GameModeRetail {
		g.titleSequence %= 7
	} else {
		g.titleSequence %= 6
	}
}

func (g *Game) tickTitleScreen() {
	g.titlePageTics -= 1
	if g.titlePageTics == 0 {
		g.titleAdvance = true
	}
}

func (g *Game) displayTitleScreen() {
	patch := engine.LumpByName[engine.Patch](g.titlePageName)
	engine.DrawPatch(0, 0, 0, patch)
}
