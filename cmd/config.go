package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strconv"
)

const Version = 110

const (
	devdata = "devdata"
	devmaps = "devmaps"
)

const (
	GameModeShareware    = iota // DOOM 1 shareware, E1 M9
	GameModeRegistered          // DOOM 1 registered, E3 M27
	GameModeRetail              // DOOM 1 retail, E4 M36
	GameModeCommercial          // DOOM 2 retail, E1 M34
	GameModeIndetermined        // Well, no IWAD found.
)

const (
	LanguageEnglish = iota
	LanguageFrench
)

const (
	StartActionTitle = iota
	StartActionAutoStart
	StartActionPlayDemo
	StartActionTimeDemo
	StartActionLoadGame
)

const (
	DeathmatchModeNone = iota
	DeathmatchMode1
	DeathmatchMode2
)

const (
	SkillBaby = iota
	SkillEasy
	SkillMedium
	SkillHard
	SkillNightmare
)

type GameConfig struct {
	GameMode int
	Language int

	BaseDefaultConfigFile string
	ConfigFile            string
	WadFiles              []string // IWAD, PWADs, and lumps
	ModifiedGame          bool

	DevParm bool
	NoDraw  bool

	// Game start settings
	StartAction       int
	RecordDemo        string
	RecordDemoMaxSize int
	DemoFile          string
	LoadGameSlot      int

	// Auto-start settings
	StartEpisode int
	StartMap     int
	StartSkill   int

	// Game play config
	NoMonsters bool
	Respawn    bool
	Fast       bool
	TurboScale int

	DoomPort       int
	DeathmatchMode int
	LevelTimer     int
}

var Config = GameConfig{Language: LanguageEnglish}

func Fatal(msg string, args ...interface{}) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func initConfig() {
	identifyVersion()
	initConfigFromArgs()

	// TODO: validate args and set Config accordingly

	slog.Debug("initConfig", "Config", Config)
}

func identifyVersion() {
	// Check for devparm modes first
	homeDir, err := os.UserHomeDir()
	if err != nil {
		Fatal(fmt.Sprintf("Failed to get user home directory: %s", err))
	}
	Config.BaseDefaultConfigFile = path.Join(homeDir, ".doomrc")

	if args.ShDev {
		Config.GameMode = GameModeShareware
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, devdata+"doom1.wad")
		Config.WadFiles = append(Config.WadFiles, devmaps+"data_se/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, devmaps+"data_se/pnames.lmp")
		Config.BaseDefaultConfigFile = devdata + "default.cfg"
		return
	}

	if args.RegDev {
		Config.GameMode = GameModeRegistered
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, devdata+"doom.wad")
		Config.WadFiles = append(Config.WadFiles, devmaps+"data_se/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, devmaps+"data_se/texture2.lmp")
		Config.WadFiles = append(Config.WadFiles, devmaps+"data_se/pnames.lmp")
		Config.BaseDefaultConfigFile = devdata + "default.cfg"
		return
	}

	if args.ComDev {
		Config.GameMode = GameModeCommercial
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, devdata+"doom2.wad")
		Config.WadFiles = append(Config.WadFiles, devmaps+"cdata/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, devmaps+"cdata/pnames.lmp")
		Config.BaseDefaultConfigFile = devdata + "default.cfg"
		return
	}

	// Check for IWADs
	wadPath := "."
	if wp := os.Getenv("DOOMWADDIR"); wp != "" {
		wadPath = wp
	}

	wadDoom2 := path.Join(wadPath, "doom2.wad")
	wadDoomU := path.Join(wadPath, "doomu.wad")
	wadDoom := path.Join(wadPath, "doom.wad")
	wadDoom1 := path.Join(wadPath, "doom1.wad")
	wadPlutonia := path.Join(wadPath, "plutonia.wad")
	wadTNT := path.Join(wadPath, "tnt.wad")
	wadDoom2F := path.Join(wadPath, "doom2f.wad")

	// Check for IWADs
	exists := func(filePath string) bool {
		_, err := os.Stat(filePath)
		return err == nil
	}

	switch {
	case exists(wadDoom2F):
		Config.GameMode = GameModeCommercial
		// C'est ridicule!
		// Let's handle languages in config files, okay?
		Config.Language = LanguageFrench
		Config.WadFiles = append(Config.WadFiles, wadDoom2F)
		slog.Info("French version.")
	case exists(wadDoom2):
		Config.GameMode = GameModeCommercial
		Config.WadFiles = append(Config.WadFiles, wadDoom2)
	case exists(wadPlutonia):
		Config.GameMode = GameModeCommercial
		Config.WadFiles = append(Config.WadFiles, wadPlutonia)
	case exists(wadTNT):
		Config.GameMode = GameModeCommercial
		Config.WadFiles = append(Config.WadFiles, wadTNT)
	case exists(wadDoomU):
		Config.GameMode = GameModeRetail
		Config.WadFiles = append(Config.WadFiles, wadDoomU)
	case exists(wadDoom):
		Config.GameMode = GameModeRegistered
		Config.WadFiles = append(Config.WadFiles, wadDoom)
	case exists(wadDoom1):
		Config.GameMode = GameModeShareware
		Config.WadFiles = append(Config.WadFiles, wadDoom1)
	default:
		Config.GameMode = GameModeIndetermined
		slog.Error("Game mode indeterminate. IWAD not found.")
		// We don't abort. Let's see what the PWAD contains.
	}
}

func initConfigFromArgs() {
	// Config file from default home location or shdev/regdev/comdev mode unless overridden by --config flag
	Config.ConfigFile = Config.BaseDefaultConfigFile
	if args.Config != "" {
		Config.ConfigFile = args.Config
	}

	// DevParm mode from explicit --devparm flag or implicit from shdev/regdev/comdev mode
	Config.DevParm = Config.DevParm || args.DevParm

	// Defunct flags
	if args.CDROM {
		slog.Warn("Ignoring defunct flag --cdrom")
	}
	if args.DebugFile {
		slog.Warn("Ignoring defunct flag --debugfile (use --log-level=debug instead)")
	}
	if args.Dup {
		slog.Warn("Ignoring defunct flag --dup")
	}
	if !args.ExtraTic {
		slog.Warn("Ignoring defunct flag --extratic (cannot be disabled)")
	}
	if args.Net {
		slog.Warn("Ignoring defunct flag --net")
	}
	if args.StatCopy {
		slog.Warn("Ignoring defunct flag --statcopy")
	}

	// Auto-start settings
	Config.StartEpisode = 1
	Config.StartMap = 1
	Config.StartSkill = SkillMedium

	if args.Skill != 0 {
		if args.Skill < 1 || args.Skill > 5 {
			Fatal(fmt.Sprintf("Invalid skill level: %d. Must be 1 through 5.", args.Skill))
		}
		Config.StartAction = StartActionAutoStart
		Config.StartSkill = args.Skill - 1
	}

	if args.Episode != 0 {
		if Config.GameMode == GameModeCommercial || Config.GameMode == GameModeIndetermined {
			Fatal("Episode specified without Doom IWAD. Use --warp instead for Doom II.")
		}
		if args.Episode < 1 || args.Episode > 4 {
			Fatal(fmt.Sprintf("Invalid episode number: %d. Must be 1 through 4.", args.Episode))
		}
		if args.Wart != "" {
			slog.Warn("Both --wart and --episode specified. Ignoring --episode.")
		} else if args.Warp != "" {
			slog.Warn("Both --warp and --episode specified. Ignoring --episode.")
		} else {
			Config.StartAction = StartActionAutoStart
			Config.StartEpisode = args.Episode
			slog.Info(fmt.Sprintf("Starting episode %d", args.Episode))
		}
	}

	// --wart and --warp. Wart is like warp but adds a dev map WAD automatically to the WadFiles list.
	if args.Wart != "" {
		if args.Warp != "" {
			slog.Warn("Both --wart and --warp specified. Ignoring --warp.")
		}
		args.Warp = args.Wart
	}

	if args.Warp != "" {
		Config.StartAction = StartActionAutoStart

		switch Config.GameMode {
		case GameModeShareware, GameModeRegistered, GameModeRetail:
			var e, m int
			var err error
			// Expect two single digits separated by optional whitespace
			_, err = fmt.Sscanf(args.Warp, "%1d%1d", &e, &m)
			if err != nil {
				Fatal(fmt.Sprintf("Invalid warp map: %s.", err))
			}

			if args.Wart != "" {
				wad := devmaps + fmt.Sprintf("E%dM%d.wad", e, m)
				Config.WadFiles = append(Config.WadFiles, wad)
			}

			Config.StartEpisode = e
			Config.StartMap = m
			slog.Info(fmt.Sprintf("Starting episode %d map %d.", e, m))

		case GameModeCommercial, GameModeIndetermined:
			m, err := strconv.ParseInt(args.Warp, 10, 0)
			if err != nil {
				Fatal(fmt.Sprintf("Invalid warp map: %s.", err))
			}

			if args.Wart != "" {
				wad := devmaps + fmt.Sprintf("cdata/map%02d.wad", m)
				Config.WadFiles = append(Config.WadFiles, wad)
			}
			slog.Info(fmt.Sprintf("Starting map %d.", m))
		}
	}

	// Additional WAD and lump files from --file flags
	Config.WadFiles = append(Config.WadFiles, args.Files...)
	Config.ModifiedGame = len(args.Files) > 0

	// Demo playback and recording

	// Demo recording is a state independent of StartAction
	Config.RecordDemo = args.Record
	Config.RecordDemoMaxSize = args.MaxDemo * 1024
	if args.Record != "" {
		slog.Info(fmt.Sprintf("Recording demo to %s. Max recording file size %dkb.", Config.RecordDemo, Config.RecordDemoMaxSize/1024))
	}

	if args.NoDraw {
		if args.TimeDemo == "" {
			slog.Warn("Ignoring --nodraw. Only has effect with --timedemo.")
		} else {
			Config.NoDraw = true
		}
	}

	if args.PlayDemo != "" || args.TimeDemo != "" {
		if args.PlayDemo != "" && args.TimeDemo != "" {
			slog.Warn("Both --playdemo and --timedemo specified. Ignoring --timedemo.")
		}

		if args.PlayDemo != "" {
			Config.StartAction = StartActionPlayDemo
			Config.DemoFile = args.PlayDemo
		} else if args.TimeDemo != "" {
			Config.StartAction = StartActionTimeDemo
			Config.DemoFile = args.TimeDemo
		}

		lmp := Config.DemoFile + ".lmp"
		Config.WadFiles = append(Config.WadFiles, lmp)

		slog.Info(fmt.Sprintf("Playing demo %s.lmp", Config.DemoFile))
	}

	if args.LoadGame != -1 {
		if args.LoadGame < 0 || args.LoadGame > 5 {
			Fatal(fmt.Sprintf("Invalid load game slot: %d. Must be 0 through 5.", args.LoadGame))
		}
		Config.StartAction = StartActionLoadGame
		Config.LoadGameSlot = args.LoadGame
	}

	// Net
	Config.DoomPort = 5029
	if args.Port != 0 {
		Config.DoomPort = args.Port
	}

	// Game play options
	Config.NoMonsters = args.NoMonsters
	Config.Respawn = args.Respawn
	Config.Fast = args.Fast

	Config.DeathmatchMode = DeathmatchModeNone
	if args.AltDeath {
		if args.Deathmatch {
			slog.Warn("Both --altdeath and --deathmatch specified. Ignoring --deathmatch.")
		}
		Config.DeathmatchMode = DeathmatchMode2
	} else if args.Deathmatch {
		Config.DeathmatchMode = DeathmatchMode1
	}

	if args.Turbo != 0 {
		if args.Turbo < 10 {
			slog.Warn("Turbo cannot be less than 10%. Limiting to 10%.")
			args.Turbo = 10
		} else if args.Turbo > 400 {
			slog.Warn("Turbo cannot be more than 400%. Limiting to 400%.")
			args.Turbo = 400
		}
		Config.TurboScale = args.Turbo

		slog.Info(fmt.Sprintf("Turbo mode enabled. Movement speed scaled to %d%%.", Config.TurboScale))
	}

	if args.Timer != 0 {
		if args.Avg {
			slog.Warn("Both --timer and --avg specified. Ignoring --avg.")
		}
		if Config.DeathmatchMode == DeathmatchModeNone {
			slog.Warn("Timer specified without deathmatch mode. Ignoring --timer.")
		} else {
			Config.LevelTimer = args.Timer
			plural := ""
			if Config.LevelTimer > 1 {
				plural = "s"
			}
			slog.Info(fmt.Sprintf("Levels will end after %d minute%s.", Config.LevelTimer, plural))
		}
	} else if args.Avg {
		if Config.DeathmatchMode == DeathmatchModeNone {
			slog.Warn("AVG mode specified without deathmatch mode. Ignoring --avg.")
		} else {
			Config.LevelTimer = 20
			slog.Info("Austin Virtual Gaming: Levels will end after 20 minutes.")
		}
	}
}
