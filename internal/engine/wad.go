package engine

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Lump struct {
	Name string
	Data []byte
}

// Lumps lists all the lumps in the WAD and lump files specified in the command line.
// Order matters and the name search will search backwards to give override priority to the last WAD or lump specified
// for each given lump name. Position of a lump in this slice is its handle, lump ID or lump num.
//
// Note: Because some lump names are not unique, namely those used for maps, and depend on order, we cannot use a map.
// Care could be taken to preprocesses the order specific lumps into special data structures, but for now we stick
// with the simple list.
var Lumps []Lump

func InitFiles(filePaths []string) {
	Lumps = make([]Lump, 0, 3072) // DOOM2.WAD has 2919 lumps

	slog.Info("Reading files...")

	// TODO: Support tilde for file reloading
	for _, filePath := range filePaths {
		slog.Info("  " + filePath)
		if strings.HasSuffix(strings.ToLower(filePath), ".wad") {
			addWad(filePath)
		} else {
			addLump(filePath)
		}
	}

	if len(Lumps) == 0 {
		slog.Error("No lumps found in WADs or lump files", "files", filePaths)
		os.Exit(1)
	}
	slog.Debug("Loaded files", "totalFiles", len(filePaths), "totalLumps", len(Lumps))
}

type wadInfo struct {
	Identification  [4]byte
	NumLumps        int32
	InfoTableOffset int32
}

type wadFileLump struct {
	FilePosition int32
	Size         int32
	Name         [8]byte
}

func addWad(filePath string) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Failed to read WAD file", "file", filePath, "err", err)
		os.Exit(1)
	}

	var wadInfo wadInfo
	_, err = binary.Decode(fileBytes, binary.LittleEndian, &wadInfo)
	if err != nil {
		slog.Error("Failed to read WAD header", "file", filePath, "err", err)
		os.Exit(1)
	}

	id := string(wadInfo.Identification[:])
	if id != "IWAD" && id != "PWAD" {
		slog.Error("WAD file does not have IWAD or PWAD id", "file", filePath, "id", id)
		os.Exit(1)
	}

	lumpTable := make([]wadFileLump, wadInfo.NumLumps)
	_, err = binary.Decode(fileBytes[wadInfo.InfoTableOffset:], binary.LittleEndian, &lumpTable)
	if err != nil {
		slog.Error("Failed to read WAD lump table", "file", filePath, "err", err)
		os.Exit(1)
	}

	for _, lump := range lumpTable {
		nameEnd := bytes.IndexByte(lump.Name[:], 0)
		if nameEnd == -1 {
			nameEnd = len(lump.Name)
		}
		Lumps = append(Lumps, Lump{
			Name: string(lump.Name[:nameEnd]),
			Data: fileBytes[lump.FilePosition : lump.FilePosition+lump.Size],
		})
	}
	slog.Debug("  Loaded WAD", "file", filePath, "lumps", wadInfo.NumLumps)
}

func addLump(filePath string) {
	name, _ := strings.CutSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Failed to read lump file", "file", filePath, "err", err)
		os.Exit(1)
	}

	Lumps = append(Lumps, Lump{
		Name: name,
		Data: fileBytes,
	})
	slog.Debug("  Loaded lump", "file", filePath)
}

func LookupLumpNum(name string) (int, bool) {
	name = strings.ToUpper(name)

	// scan backwards so patch lump files take precedence
	for i := len(Lumps) - 1; i >= 0; i-- {
		if Lumps[i].Name == name {
			return i, true
		}
	}

	return -1, false
}
