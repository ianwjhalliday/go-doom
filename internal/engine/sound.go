package engine

const SampleRate = 11025

func InitSound() {
	InitSoundBeep()
}

func ShutdownSound() {
	ShutdownSoundBeep()
}

type SoundLump struct {
	FormatNum   uint16 // always 3
	SampleRate  uint16 // 11025 except a couple that are 22050
	SampleCount uint32
	Samples     [8]byte // actually [SampleCount]byte, need unsafe pointer
}

func PlaySound(name string) {
	sound := LumpByName[SoundLump]("ds" + name)
	PlaySoundBeep(sound)
}
