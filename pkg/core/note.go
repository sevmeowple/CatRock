package core

import (
	"fmt"
	"time"
)



type BaseNoteName int // 基础音符类型

const (
	C BaseNoteName = iota // C音符
	Cs
	D
	Ds
	E
	F
	Fs
	G
	Gs
	A
	As
	B // B音符
) //12

type Accidental int // 音符的变音类型
const (
	DoubleFlat  Accidental = iota - 2 // 双降音
	Flat                              // 降音
	Natural                           // 自然音
	Sharp                             // 升音
	DoubleSharp                       // 双升音
)

// 添加节拍类型定义
type BeatValue float64

const (
	Whole        BeatValue = 4.0   // 全音符
	Half         BeatValue = 2.0   // 二分音符
	Quarter      BeatValue = 1.0   // 四分音符
	Eighth       BeatValue = 0.5   // 八分音符
	Sixteenth    BeatValue = 0.25  // 十六分音符
	ThirtySecond BeatValue = 0.125 // 三十二分音符
)

type Note struct {
	Name       BaseNoteName
	Octave     int // 八度
	Accidental Accidental
	MIDINote   []byte
	// 节拍数
	Beat BeatValue

	TrackID    int   // Track ID for MIDI
	Channel    uint8 // MIDI channel
	Instrument int   // MIDI instrument number
	Velocity   uint8 // MIDI velocity
}

type NewNoteParams struct {
	Name       BaseNoteName // 音符名称
	Octave     int          // 八度
	Accidental Accidental   // 变音
	Beat       BeatValue    // 节拍

	TrackID	int   // Track ID for MIDI
	Channel	uint8 // MIDI channel
	Instrument int   // MIDI instrument number
	Velocity	uint8 // MIDI velocity
}

func NewNote(params NewNoteParams) Note {
	note := Note{
		Name:       params.Name,
		Octave:     params.Octave,
		Accidental: params.Accidental,
		Beat:       params.Beat,
		TrackID:    params.TrackID,
		Channel:    params.Channel,
		Instrument: params.Instrument,
		Velocity:   params.Velocity,
	}
	switch params.Name {
	case C:
		note.MIDINote = []byte{0x3C} // MIDI note number for C
	case Cs:
		note.MIDINote = []byte{0x3C + 1} // MIDI note number for C#
	case D:
		note.MIDINote = []byte{0x3E} // MIDI note number for D
	case Ds:
		note.MIDINote = []byte{0x3E + 1} // MIDI note number for D#
	case E:
		note.MIDINote = []byte{0x40} // MIDI note number for E
	case F:
		note.MIDINote = []byte{0x41} // MIDI note number for F
	case Fs:
		note.MIDINote = []byte{0x41 + 1} // MIDI note number for F#
	case G:
		note.MIDINote = []byte{0x43} // MIDI note number for G
	case Gs:
		note.MIDINote = []byte{0x43 + 1} // MIDI note number for G#
	case A:
		note.MIDINote = []byte{0x45} // MIDI note number for A
	case As:
		note.MIDINote = []byte{0x45 + 1} // MIDI note number for A#
	case B:
		note.MIDINote = []byte{0x47} // MIDI note number for B
	default:
		note.MIDINote = []byte{0x00} // 默认值，可能需要处理错误
	}
	if params.Octave < 0 || params.Octave > 10 {
		params.Octave = 0 // 限制八度范围
	}
	note.MIDINote[0] += byte(params.Octave * 12) // MIDI音符编号计算
	if params.Accidental == DoubleFlat {
		note.MIDINote[0] -= 2 // 双降音
	} else if params.Accidental == Flat {
		note.MIDINote[0] -= 1 // 降音
	} else if params.Accidental == Sharp {
		note.MIDINote[0] += 1 // 升音
	} else if params.Accidental == DoubleSharp {
		note.MIDINote[0] += 2 // 双升音
	} // 自然音不做任何变动

	if note.Beat == 0 {
		note.Beat = Quarter // 默认四分音符
	}

	if note.Velocity == 0 {
		note.Velocity = 100 // 默认音量
	}

	if note.MIDINote[0] < 0 || note.MIDINote[0] > 127 {
		// 处理错误
	}

	return note
}

func (n Note) Duration(bpm float64) time.Duration {
	// BPM 是每分钟的四分音符数
	// 一个四分音符的时长 = 60秒 / BPM
	quarterNoteDuration := 60.0 / bpm

	// 当前音符的时长 = 四分音符时长 * 节拍比例
	noteDurationSeconds := quarterNoteDuration * float64(n.Beat)

	return time.Duration(noteDurationSeconds * float64(time.Second))
}

func (note *Note) SetBeat(beat BeatValue) {
	// 检查是否为有效的节拍值
	validBeats := []BeatValue{Whole, Half, Quarter, Eighth, Sixteenth, ThirtySecond}

	for _, validBeat := range validBeats {
		if beat == validBeat {
			note.Beat = beat
			return
		}
	}

	// 如果不是有效值，使用默认值
	note.Beat = Quarter
}
func (noteName BaseNoteName) String() string {
	names := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	if noteName >= 0 && int(noteName) < len(names) {
		return names[noteName]
	}
	return ""
}

func (acc Accidental) String() string {
	switch acc {
	case DoubleFlat:
		return "bb"
	case Flat:
		return "b"
	case Natural:
		return ""
	case Sharp:
		return "#"
	case DoubleSharp:
		return "x"
	default:
		return ""
	}
}

func (n Note) String() string {
	return fmt.Sprintf("%s%d%s", n.Name, n.Octave, n.Accidental)
}
