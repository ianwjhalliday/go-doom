package mus

import (
	"bufio"
	"bytes"
	"fmt"
	"unsafe"
)

type Header struct {
	ID                [4]byte // Valid data starts with "MUS" 0x1A
	ScoreLength       uint16  // num bytes of score
	ScoreOffset       uint16  // offset from header base to score data
	ChannelsPrimary   uint16  // count of required channels, 0-8
	ChannelsSecondary uint16  // count of optional channels, 9-14 -- channel 15 is reserved for drums
	InstrumentCount   uint16
	dummy             uint16
	// instruments follow header at end of dummy, for InstrumentCount length
	// Score MUS events follow header at ScoreOffset, for ScoreLength length
}

type Reader struct {
	*bufio.Reader
}

func NewReader(data *Header) (*Reader, error) {
	if data.ID != [4]byte{'M', 'U', 'S', 0x1A} {
		return nil, fmt.Errorf("invalid MUS, header ID is not 'MUS\\x1A'. Got '%v' instead", data.ID)
	}

	eventBytes := (*byte)(unsafe.Add(unsafe.Pointer(data), data.ScoreOffset))
	reader := bufio.NewReader(bytes.NewReader(unsafe.Slice(eventBytes, int(data.ScoreLength))))

	return &Reader{Reader: reader}, nil
}

type EventType byte

const (
	ReleaseNote      = 0
	PlayNote         = 1
	PitchBend        = 2
	SystemEvent      = 3
	ChangeController = 4
	EndOfMeasure     = 5
	ScoreEnd         = 6
	Unused           = 7
)

type Event struct {
	Type       EventType
	Data1      byte
	Data2      byte
	Channel    byte
	DelayTicks int32 // 140 Hz ticks
}

func ReadEvents(r *Reader) ([]Event, error) {
	var events []Event
	for {
		event, err := ReadEvent(r)
		if err != nil {
			return events, err
		}
		events = append(events, event)
		if event.Type == ScoreEnd {
			break
		}
	}
	return events, nil
}

func ReadEvent(r *Reader) (Event, error) {
	b, err := r.ReadByte()
	if err != nil {
		return Event{}, fmt.Errorf("ReadEvent: Failed to read event descriptor byte: %v", err)
	}

	last := (b & 0b10000000) != 0
	eventType := EventType((b & 0b01110000) >> 4)
	channel := b & 0b00001111

	var data1, data2 byte
	if eventType <= ChangeController {
		data1, err = r.ReadByte()
		if err != nil {
			return Event{}, fmt.Errorf("ReadEvent: Failed to read data1 byte: %v", err)
		}
	}

	if eventType == PlayNote && data1&0b10000000 != 0 {
		data2, err = r.ReadByte()
		if err != nil {
			return Event{}, fmt.Errorf("ReadEvent: Failed to read data2 byte: %v", err)
		}
	}

	if eventType == ChangeController {
		data2, err = r.ReadByte()
		if err != nil {
			return Event{}, fmt.Errorf("ReadEvent: Failed to read data2 byte: %v", err)
		}
	}

	var delayTicks int32
	for last {
		b, err = r.ReadByte()
		if err != nil {
			return Event{}, fmt.Errorf("ReadEvent: Failed to read follow byte: %v", err)
		}
		last = (b & 0b10000000) != 0

		delayTicBits := b & 0b01111111
		delayTicks = int32(delayTicBits) + (delayTicks << 7)
	}

	return Event{
		Type:       eventType,
		Channel:    channel,
		Data1:      data1,
		Data2:      data2,
		DelayTicks: delayTicks,
	}, nil
}
