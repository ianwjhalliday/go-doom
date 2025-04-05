package dstrings

var DOSY = "(press y to quit)"

var EndMessages = []string{
	// DOOM1
	"are you sure you want to\nquit this great game?",
	"please don't leave, there's more\ndemons to toast!",
	"let's beat it -- this is turning\ninto a bloodbath!",
	"i wouldn't leave if i were you.\ndos is much worse.",
	"you're trying to say you like dos\nbetter than me, right?",
	"don't leave yet -- there's a\ndemon around that corner!",
	"ya know, next time you come in here\ni'm gonna toast ya.",
	"go ahead and leave. see if i care.",

	// QuitDOOM II messages
	"you want to quit?\nthen, thou hast lost an eighth!",
	"don't go now, there's a \ndimensional shambler waiting\nat the dos prompt!",
	"get outta here and go back\nto your boring programs.",
	"if i were your boss, i'd \n deathmatch ya in a minute!",
	"look, bud. you leave now\nand you forfeit your body count!",
	"just leave. when you come\nback, i'll be waiting with a bat.",
	"you're lucky i don't smack\nyou for thinking about leaving.",

	// Internal debug. Different style, too.
	"THIS IS NO MESSAGE!\nPage intentionally left blank.",
}

// InitFrench overwrites all string globals with French counterparts
func InitFrench() {
	DOSY = frenchDOSY

	EndMessages[0] = frenchQuitMessage
}
