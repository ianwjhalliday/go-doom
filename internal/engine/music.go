package engine

import (
	"github.com/ianwjhalliday/go-doom/internal/engine/mus"
)

func InitMusic() {
	InitMusicMeltySynth()
}

func PlaySong(song string) {
	songLump := LumpByName[mus.Header]("d_" + song)

	reader, err := mus.NewReader(songLump)
	if err != nil {
		FatalError("PlaySong: Failed to create MUS reader", "song", song, "err", err)
	}

	events, err := mus.ReadEvents(reader)
	if err != nil {
		FatalError("PlaySong: Failed to read MUS events", "song", song, "err", err)
	}

	PlaySongMeltySynth(songLump, events)
}
