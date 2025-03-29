package engine

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"time"
	"unsafe"
)

func InitSoundBeep() {
	sr := beep.SampleRate(SampleRate)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		FatalError("InitSound: Failed to initialize beep.speaker", "err", err)
	}
}

func ShutdownSoundBeep() {
	// NOTE: beep uses oto and oto cannot be closed once initialized which apparently can prevent
	// sleep of a computer. So e.g. if we run go-doom and then switch to another app, pausing doom,
	// we wouldn't be able to let the computer sleep? There may also be GL rendering preventing that.
	// Something to consider and investigate later.
	// See https://github.com/ebitengine/oto/issues/149#issuecomment-2506803480.
	// So speaker.Close as oto is, is unnecessary, but we'll keep it here for now.
	speaker.Close()
}

type playingSound struct {
	samples    []byte
	sampleRate uint16
	position   uint16
	volLeft    int16
	volRight   int16
}

func (s *playingSound) Stream(samples [][2]float64) (int, bool) {
	// Doom sound samples are mono 8-bit PCM 0-255
	// beep expects float64 samples -1 to 1
	n := int(s.position)

	if n >= len(s.samples) {
		return 0, false
	}

	increment := 1
	if s.sampleRate == 22050 {
		// Hack downsample to 11025 by skipping every other sample
		increment = 2
	}

	// Doom volumes are integer 0-15
	volLeft := float64(s.volLeft) / 15
	volRight := float64(s.volRight) / 15

	for i := range samples {
		if n >= len(s.samples) {
			break
		}

		val := float64(s.samples[n])/(1<<8-1)*2 - 1
		samples[i][0] = val * volLeft
		samples[i][1] = val * volRight
		n += increment
	}

	c := (n - int(s.position)) / increment
	s.position = uint16(n)
	return c, true
}

func (s *playingSound) Err() error {
	return nil
}

func PlaySoundBeep(sound *SoundLump) {
	speaker.Play(&playingSound{
		samples:    unsafe.Slice(&sound.Samples[0], sound.SampleCount),
		sampleRate: sound.SampleRate,
		volLeft:    15,
		volRight:   15,
	})
}
