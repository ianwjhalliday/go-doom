package main

import (
	"bufio"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/pflag"
	"log"
	"log/slog"
	"math"
	"os"
	"path"
)

type Arguments struct {
	AltDeath   bool
	Avg        bool
	CDROM      bool
	ComDev     bool
	Config     string
	Deathmatch bool
	DebugFile  bool
	DevParm    bool
	Dup        bool
	Episode    int
	ExtraTic   bool
	Fast       bool
	Files      []string
	LoadGame   int
	MaxDemo    int
	Net        bool
	NoDraw     bool
	NoMonsters bool
	PlayDemo   string
	Port       int
	Record     string
	RegDev     bool
	Respawn    bool
	Skill      int
	ShDev      bool
	StatCopy   bool
	TimeDemo   string
	Timer      int
	Turbo      int
	Warp       string
	Wart       string

	LogLevel string
}

var args = Arguments{}

func initArgs() {
	preprocessArgs()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("FATAL Failed to get user home directory: %s", err)
	}
	defaultConfigPath := path.Join(homeDir, ".doomrc")

	pflag.BoolVar(&args.AltDeath, "altdeath", false, "enable deathmatch 2.0 mode (item respawn and normal weapon ammo)")
	pflag.BoolVar(&args.Avg, "avg", false, "enable Austin Virtual Gaming mode (same as --timer 20)")
	pflag.BoolVar(&args.CDROM, "cdrom", false, "defunct: enabled running game directly from CD-ROM")
	pflag.BoolVar(&args.ComDev, "comdev", false, "enable commercial game (Doom II) development mode")
	pflag.StringVar(&args.Config, "config", defaultConfigPath, "use configuration file `<path>`")
	pflag.BoolVar(&args.Deathmatch, "deathmatch", false, "enable deathmatch mode")
	pflag.BoolVar(&args.DebugFile, "debugfile", false, "output tic debug info to debug0.txt")
	pflag.BoolVar(&args.DevParm, "devparm", false, "enable developer mode")
	pflag.BoolVar(&args.Dup, "dup", false, "defunct: enabled tic duplication to save network bandwidth at cost of movement accuracy")
	pflag.IntVar(&args.Episode, "episode", 0, "start episode `<num>`. Equivalent to -warp \"1 num\"")
	pflag.BoolVar(&args.ExtraTic, "extratic", true, "defunct: enabled extra tic data over network for resilience at cost of bandwidth")
	pflag.BoolVar(&args.Fast, "fast", false, "enable fast monsters mode. Requires -warp")
	pflag.StringArrayVar(&args.Files, "file", []string{}, "load additional WADs and lump `<file>`s. Specify multiple times for multiple files")
	pflag.IntVar(&args.LoadGame, "loadgame", -1, "load the savegame at slot `<num>` (0-5)")
	pflag.IntVar(&args.MaxDemo, "maxdemo", 0, "set demo recording maximum file `<size>` in KiB")
	pflag.BoolVar(&args.Net, "net", false, "defunct: enabled networking in Doom 1.0 and 1.1")
	pflag.BoolVar(&args.NoDraw, "nodraw", false, "disable graphics rendering for -timedemo performance testing")
	pflag.BoolVar(&args.NoMonsters, "nomonsters", false, "disable monster spawning. Requires -warp")
	pflag.StringVar(&args.PlayDemo, "playdemo", "", "play demo recording `<lump>`")
	pflag.IntVar(&args.Port, "port", 0, "specify network `port` for UDP connection")
	pflag.StringVar(&args.Record, "record", "", "record demo to a lump `<file>`")
	pflag.BoolVar(&args.RegDev, "regdev", false, "enable registered game (Doom) development mode")
	pflag.BoolVar(&args.Respawn, "respawn", false, "enable monster respawning. Requires -warp")
	pflag.BoolVar(&args.ShDev, "shdev", false, "enable shareware game (Doom Shareware) development mode")
	pflag.IntVar(&args.Skill, "skill", 0, "start game immediately at skill `<level>` (1-5). Can use with -warp")
	pflag.BoolVar(&args.StatCopy, "statcopy", false, "defunct: enabled statistics copying to a stats driver program")
	pflag.StringVar(&args.TimeDemo, "timedemo", "", "play demo `<lump>` and measure avg frame rate with 35 fps cap removed")
	pflag.IntVar(&args.Timer, "timer", 0, "exit levels after `<num>` minutes. Requires -deathmatch or -altdeath")
	pflag.IntVar(&args.Turbo, "turbo", 0, "scale player movement speed by `<num>` percent (10% - 400%)")
	pflag.StringVar(&args.Warp, "warp", "", "warp to level `\"<episode> <map>\"` (Doom) or <map> (Doom II)")
	pflag.StringVar(&args.Wart, "wart", "", "defunct: developer mode warped to `<map>` loaded from a WAD file")

	pflag.StringVar(&args.LogLevel, "log-level", "info", "Set the log `<level>` (debug, info, warn, error, none)")

	pflag.Parse()

	slog.SetLogLoggerLevel(logLevel(args.LogLevel))
	slog.Debug("initArgs", "args", args)
}

func logLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "none":
		return math.MaxInt
	default:
		log.Fatalf("FATAL Invalid log level: %s", level)
		return 0
	}
}

func preprocessArgs() {
	// Expand response file arguments
	newArgs := make([]string, 0, len(os.Args))
	replaceOSArgs := false
	for _, arg := range os.Args {
		if arg[0] == '@' {
			inlinedArgs := readResponseFile(arg[1:])
			newArgs = append(newArgs, inlinedArgs...)
			replaceOSArgs = true
			continue
		}
		newArgs = append(newArgs, arg)
	}

	if replaceOSArgs {
		os.Args = newArgs
	}
}

func readResponseFile(responsePath string) []string {
	file, err := os.Open(responsePath)
	if err != nil {
		log.Fatalf("FATAL Failed to open response file: %s", err)
	}
	defer file.Close()

	args := make([]string, 0)
	scanner := bufio.NewScanner(file)

	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		lineArgs, err := shellwords.Parse(line)
		if err != nil {
			log.Fatalf("FATAL Failed to parse response file line %d '%s': %s", lineNumber, line, err)
		}
		args = append(args, lineArgs...)
	}

	return args
}
