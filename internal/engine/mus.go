package engine

import (
	"bufio"
	"bytes"
	_ "embed"
	"github.com/gopxl/beep/speaker"
	"github.com/sinshu/go-meltysynth/meltysynth"
	"unsafe"
)

var synthesizer *meltysynth.Synthesizer

func InitMus() {
	soundFont, err := meltysynth.NewSoundFont(bytes.NewReader(RolandSC55Presets))
	if err != nil {
		FatalError("InitMus: Failed to load soundfont", "err", err)
	}

	settings := meltysynth.NewSynthesizerSettings(SampleRate)
	settings.EnableReverbAndChorus = false // TODO: Reverb and chorus sound great but don't sound authentic to the original experience; would be great to provide user setting to toggle this.
	synthesizer, err = meltysynth.NewSynthesizer(soundFont, settings)
	if err != nil {
		FatalError("InitMus: Failed to create synthesizer", "err", err)
	}
}

// TODO: Provide a way to change the soundfont

//go:embed "soundfont/Roland SC-55 Presets.sf2"
var RolandSC55Presets []byte

type musEventType byte

const (
	musReleaseNote      = 0
	musPlayNote         = 1
	musPitchBend        = 2
	musSystemEvent      = 3
	musChangeController = 4
	musEndOfMeasure     = 5
	musScoreEnd         = 6
	musUnused           = 7
)

type musEvent struct {
	eventType  musEventType
	data1      byte
	data2      byte
	channel    byte
	delayTicks int32 // 140 Hz ticks
}

func musEvents(r *bufio.Reader) []musEvent {
	var events []musEvent
	for {
		event := musNextEvent(r)
		events = append(events, event)
		if event.eventType == musScoreEnd {
			break
		}
	}
	return events
}

func musNextEvent(r *bufio.Reader) musEvent {
	b, err := r.ReadByte()
	if err != nil {
		FatalError("musNextEvent: Failed to read event descriptor byte", "err", err)
	}

	last := (b & 0b10000000) != 0
	eventType := musEventType((b & 0b01110000) >> 4)
	channel := b & 0b00001111

	var data1, data2 byte
	if eventType <= musChangeController {
		data1, err = r.ReadByte()
		if err != nil {
			FatalError("musNextEvent: Failed to read data1 byte", "err", err)
		}
	}

	if eventType == musPlayNote && data1&0b10000000 != 0 {
		//data1 &= 0b01111111
		data2, err = r.ReadByte()
		if err != nil {
			FatalError("musNextEvent: Failed to read data2 byte", "err", err)
		}
	}

	if eventType == musChangeController {
		data2, err = r.ReadByte()
		if err != nil {
			FatalError("musNextEvent: Failed to read data2 byte", "err", err)
		}
	}

	var delayTicks int32
	for last {
		b, err = r.ReadByte()
		if err != nil {
			FatalError("musNextEvent: Failed to read follow byte", "err", err)
		}
		last = (b & 0b10000000) != 0

		delayTicBits := b & 0b01111111
		delayTicks = int32(delayTicBits) + (delayTicks << 7)
	}

	return musEvent{
		eventType:  eventType,
		channel:    channel,
		data1:      data1,
		data2:      data2,
		delayTicks: delayTicks,
	}
}

type musStreamer struct {
	musLump    *MusLump
	events     []musEvent
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
		delay = synthesizer.SampleRate * event.delayTicks / 140
		m.processMusMessage(mapMusChannel(event.channel), event.eventType, int32(event.data1), int32(event.data2))
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

func (m *musStreamer) processMusMessage(channel int32, event musEventType, data1 int32, data2 int32) {
	switch event {
	case musReleaseNote:
		synthesizer.NoteOff(channel, data1)
	case musPlayNote:
		note := data1 & 0x7f
		hasVolume := data1&0x80 != 0
		if hasVolume {
			m.chVel[channel] = byte(data2 & 0x7f)
		}
		synthesizer.NoteOn(channel, note, int32(m.chVel[channel]))
	case musPitchBend:
		synthesizer.ProcessMidiMessage(channel, 0xE0, data1*64, 0)
	case musSystemEvent:
		if data1 >= 10 && data1 <= 14 {
			var val int32
			if data1 == 12 {
				val = int32(m.musLump.channelsPrimary + m.musLump.channelsSecondary)
			}
			synthesizer.ProcessMidiMessage(channel, 0xB0, controllerMap[data1], val)
		}
	case musChangeController:
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
	case musScoreEnd, musEndOfMeasure, musUnused:
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

type MusLump struct {
	id                [4]byte // Valid data starts with "MUS" 0x1A
	scoreLength       uint16  // num bytes of score
	scoreStart        uint16  // offset to score data
	channelsPrimary   uint16
	channelsSecondary uint16
	instrumentCount   uint16
	dummy             uint16
	instruments       *uint16

	//data [4]byte
}

func PlaySong(song string) {
	songLump := LumpByName[MusLump]("d_" + song)
	if songLump.id != [4]byte{'M', 'U', 'S', 0x1A} {
		FatalError("PlaySong: Invalid MUS lump does not have MUS\\x1A id.", "song", song, "id", songLump.id)
	}

	eventBytes := (*byte)(unsafe.Add(unsafe.Pointer(songLump), songLump.scoreStart))
	reader := bufio.NewReader(bytes.NewReader(unsafe.Slice(eventBytes, int(songLump.scoreLength))))

	events := musEvents(reader)

	streamer := musStreamer{
		musLump: songLump,
		events:  events,
		loop:    false,
		playing: true,
	}

	speaker.Play(&streamer)
}

func StopSong() {

}
