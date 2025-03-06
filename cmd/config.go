package main

import (
	"log"
	"log/slog"
	"os"
	"path"
)

const Version = 110

const (
	GameModeShareware    = iota // DOOM 1 shareware, E1 M9
	GameModeRegistered          // DOOM 1 registered, E1 M27
	GameModeCommercial          // DOOM 2 retail, E1 M34
	GameModeRetail              // DOOM 1 retail, E4 M36
	GameModeIndetermined        // Well, no IWAD found.
)

const (
	LanguageEnglish = iota
	LanguageFrench
)

type GameConfig struct {
	GameMode int
	Language int

	WadFiles []string // IWAD, PWADs, and lumps

	DevParm               bool
	BaseDefaultConfigFile string
}

var Config = GameConfig{Language: LanguageEnglish}

func initConfig() {
	identifyVersion()

	// TODO: validate args and set Config accordingly

	slog.Debug("initConfig", "Config", Config)
}

func identifyVersion() {
	// Check for devparm modes first
	const DEVDATA = "devdata"
	const DEVMAPS = "devmaps"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("FATAL Failed to get user home directory: %s", err)
	}
	Config.BaseDefaultConfigFile = path.Join(homeDir, ".doomrc")

	if args.ShDev {
		Config.GameMode = GameModeShareware
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, DEVDATA+"doom1.wad")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"data_se/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"data_se/pnames.lmp")
		Config.BaseDefaultConfigFile = DEVDATA + "default.cfg"
		return
	}

	if args.RegDev {
		Config.GameMode = GameModeRegistered
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, DEVDATA+"doom.wad")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"data_se/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"data_se/texture2.lmp")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"data_se/pnames.lmp")
		Config.BaseDefaultConfigFile = DEVDATA + "default.cfg"
		return
	}

	if args.ComDev {
		Config.GameMode = GameModeCommercial
		Config.DevParm = true
		Config.WadFiles = append(Config.WadFiles, DEVDATA+"doom2.wad")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"cdata/texture1.lmp")
		Config.WadFiles = append(Config.WadFiles, DEVMAPS+"cdata/pnames.lmp")
		Config.BaseDefaultConfigFile = DEVDATA + "default.cfg"
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
		// TODO: Print "French version."
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
		// TODO: Print "Game mode indeterminate."
		// We don't abort. Let's see what the PWAD contains.
	}
}
