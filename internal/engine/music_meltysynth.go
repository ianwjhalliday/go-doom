package engine

import (
	"bytes"
	_ "embed"
	"github.com/gopxl/beep/speaker"
	"github.com/ianwjhalliday/go-doom/internal/engine/mus"
	"github.com/sinshu/go-meltysynth/meltysynth"
)

//go:embed "soundfont/Roland SC-55 Presets.sf2"
var RolandSC55Presets []byte

var synthesizer *meltysynth.Synthesizer

func InitMusicMeltySynth() {
	// TODO: Provide a way to change the soundfont
	soundFont, err := meltysynth.NewSoundFont(bytes.NewReader(RolandSC55Presets))
	if err != nil {
		FatalError("InitMusicMeltySynth: Failed to load soundfont", "err", err)
	}

	settings := meltysynth.NewSynthesizerSettings(SampleRate)
	// TODO: Reverb and chorus sound great but don't sound authentic to the original experience.
	//  It would be great to provide user setting to toggle this.
	settings.EnableReverbAndChorus = false
	synthesizer, err = meltysynth.NewSynthesizer(soundFont, settings)
	if err != nil {
		FatalError("InitMusicMeltySynth: Failed to create synthesizer", "err", err)
	}
}

func PlaySongMeltySynth(header *mus.Header, events []mus.Event) {
	streamer := musStreamer{
		header:  header,
		events:  events,
		loop:    false,
		playing: true,
	}

	speaker.Play(&streamer)
}

type musStreamer struct {
	header     *mus.Header
	events     []mus.Event
	loop       bool
	playing    bool
	delayRem   int32
	eventIndex int32
	chVel      [16]byte
}

func (m *musStreamer) Err() error {
	return nil
}

func (m *musStreamer) Stream(samples [][2]float64) (int, bool) {
	var wrote int32
	length := int32(len(samples))

	for m.playing && wrote < length {
		if m.delayRem == 0 {
			m.delayRem = m.processEvents()
		}
		dstRem := length - wrote
		rem := min(m.delayRem, dstRem)

		synthesizer.Render2xF64(samples[wrote : wrote+rem])

		m.delayRem -= rem
		wrote += rem
	}

	// TODO: Should ok = m.playing instead? false will remove the streamer from the speaker.
	return int(wrote), wrote > 0
}

func (m *musStreamer) processEvents() int32 {
	// process events until one has a non-zero delay or we reach the end
	delay := int32(0)
	eventsLength := int32(len(m.events))

	for delay == 0 && m.eventIndex < eventsLength {
		event := m.events[m.eventIndex]
		delay = synthesizer.SampleRate * event.DelayTicks / 140
		m.processMusMessage(mapMusChannel(event.Channel), event.Type, int32(event.Data1), int32(event.Data2))
		m.eventIndex++
	}

	if m.eventIndex == eventsLength {
		if m.loop {
			m.eventIndex = 0
			synthesizer.NoteOffAll(false)
		} else {
			m.playing = false
		}
	}

	return delay
}

func (m *musStreamer) processMusMessage(channel int32, event mus.EventType, data1 int32, data2 int32) {
	switch event {
	case mus.ReleaseNote:
		synthesizer.NoteOff(channel, data1)
	case mus.PlayNote:
		note := data1 & 0x7f
		hasVolume := data1&0x80 != 0
		if hasVolume {
			m.chVel[channel] = byte(data2 & 0x7f)
		}
		synthesizer.NoteOn(channel, note, int32(m.chVel[channel]))
	case mus.PitchBend:
		synthesizer.ProcessMidiMessage(channel, 0xE0, data1*64, 0)
	case mus.SystemEvent:
		if data1 >= 10 && data1 <= 14 {
			var val int32
			if data1 == 12 {
				val = int32(m.header.ChannelsPrimary + m.header.ChannelsSecondary)
			}
			synthesizer.ProcessMidiMessage(channel, 0xB0, controllerMap[data1], val)
		}
	case mus.ChangeController:
		if data1 == 0 {
			// TODO: What did this hack do and do we need to replicate with meltysynth?
			// HACK: when channels are setup, clear all sounds to turn off notes still playing
			// from previous looping track. Do this for every channel since not every song uses
			// the same number of channels.
			for ch := range 16 {
				synthesizer.ProcessMidiMessage(int32(ch), 0x78, 0, 0)
			}
			if channel == 9 {
				synthesizer.ProcessMidiMessage(channel, 0xB0, 0, 0)
			}
			synthesizer.ProcessMidiMessage(channel, 0xC0, data2, 0)
		} else if data1 >= 1 && data1 <= 9 {
			synthesizer.ProcessMidiMessage(channel, 0xB0, controllerMap[data1], data2)
		}
	case mus.ScoreEnd, mus.EndOfMeasure, mus.Unused:
		// ignored events
	}
}

func mapMusChannel(channel byte) int32 {
	if channel == 15 {
		return 9
	}
	if channel >= 9 {
		return int32(channel) + 1
	}
	return int32(channel)
}

var controllerMap = [15]int32{
	// MUS Change controller events (0-9)
	0,  // Instrument (patch, program) number
	0,  // (or 32?)    Bank select: 0 by default
	1,  // Modulation pot (frequency vibrato depth)
	7,  // Volume: 0-silent, ~100-normal, 127-loud
	10, // Pan (balance) pot: 0-left, 64-center (default), 127-right
	11, // Expression pot
	91, // Reverb depth
	93, // Chorus depth
	64, // Sustain pedal (hold)
	67, // Soft pedal
	// MUS System Events (10-14)
	120, // All sounds off
	123, // All notes off
	126, // Mono
	127, // Poly
	121, // Reset all controllers
}
