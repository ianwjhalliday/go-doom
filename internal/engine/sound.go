package engine

const SampleRate = 11025

func InitSound() {
	// TODO: Either oto, beep, or meltysynth are burning through idle wake ups (and maybe CPU).
	//  Need to investigate and find out if normal. Doesn't seem like it would be normal based on
	//  magnitude of wake ups, ~3000-4000 per second.
	//  doomzig only does around 150, albeit sound is broken in doomzig so unclear if miniaudio is
	//  actually running or not.
	//
	//  A brief look with Xcode's Instruments app suggests it could be use of or misuse of
	//  AudioToolkit, CoreAudio, or AQ something. At a glance miniaudio uses CoreAudio directly
	//  so perhaps it has to do with Oto using AudioToolkit. Unclear.
	InitSoundBeep()
	InitMusic()
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
