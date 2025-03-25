package game

import (
	"github.com/ianwjhalliday/go-doom/internal/engine"
	"log/slog"
	"os"
)

// VetFiles checks the game files to ensure they are valid. However, this is legacy holdover from
// the original linux Doom source code. It only checks Doom 1 shareware and registered versions. It
// does not check for Ultimate Doom or Doom 2, or the expansions. It is also merely based on the
// name of the IWAD so is easily circumvented. Could remove this if desired.
func VetFiles(config *GameConfig) {
	if config.ModifiedGame {
		if config.GameMode == GameModeShareware {
			slog.Error("You cannot --file with the shareware version. Register!")
			os.Exit(1)
		}
	}

	if config.GameMode == GameModeRegistered {
		names := []string{"e2m1", "e2m2", "e2m3", "e2m4", "e2m5", "e2m6", "e2m7", "e2m8", "e2m9", "e3m1", "e3m3", "e3m3", "e3m4", "e3m5", "e3m6", "e3m7", "e3m8", "e3m9", "dphoof", "bfgga0", "heada1", "cybra1", "spida1d1"}
		for _, name := range names {
			if _, found := engine.LookupLumpNum(name); !found {
				slog.Error("This is not the registered version. Missing lumps.", "lump", name)
				os.Exit(1)
			}
		}
	}
}

// PrintBanner prints out the banners from the original linux Doom source code. This is for fun but
// probably should be changed or removed if releasing go-doom under an OSS license.
func PrintBanner(config *GameConfig) {
	// If additional PWAD files are used, print modified banner
	if config.ModifiedGame {
		slog.Info("===========================================================================")
		slog.Info("ATTENTION:  This version of DOOM has been modified.  If you would like to")
		slog.Info("get a copy of the original game, call 1-800-IDGAMES or see the readme file.")
		slog.Info("        You will not receive technical support for modified games.")
		slog.Info("                      press enter to continue")
		slog.Info("===========================================================================")
		//_, _ = bufio.NewReader(os.Stdin).ReadString('\n') // skip press enter to continue
	}

	switch config.GameMode {
	case GameModeShareware, GameModeIndetermined:
		slog.Info("===========================================================================")
		slog.Info("                                Shareware!")
		slog.Info("===========================================================================")
	case GameModeRegistered, GameModeRetail, GameModeCommercial:
		slog.Info("===========================================================================")
		slog.Info("                 Commercial product - do not distribute!")
		slog.Info("         Please report software piracy to the SPA: 1-800-388-PIR8")
		slog.Info("===========================================================================")
	}
}
