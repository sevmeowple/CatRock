package core

type Chord struct {
	Notes []Note // 和弦包含的音符
}

type ChordQuality int

const (
    Major ChordQuality = iota
    Minor
    Diminished
    Augmented
    Major7
    Minor7
    Dominant7
    Sus2
    Sus4
)

// NewChord 创建一个新的和弦
func NewChord(notes []Note) Chord {
	return Chord{Notes: notes}
}

// getChordIntervals 返回和弦的音程
func getChordIntervals(quality ChordQuality) []int {
	switch quality {
	case Major:
		return []int{0, 4, 7}
	case Minor:
		return []int{0, 3, 7}
	case Diminished:
		return []int{0, 3, 6}
	case Augmented:
		return []int{0, 4, 8}
	case Major7:
		return []int{0, 4, 7, 11}
	case Minor7:
		return []int{0, 3, 7, 10}
	case Dominant7:
		return []int{0, 4, 7, 10}
	case Sus2:
		return []int{0, 2, 7}
	case Sus4:
		return []int{0, 5, 7}
	default:
		return []int{0, 4, 7} // 默认大三和弦
	}
}

func NewChordFromQuality(root Note, quality ChordQuality) Chord {
	intervals := getChordIntervals(quality)
	notes := make([]Note, len(intervals))

	for i, interval := range intervals {
		// 计算音符的MIDI音符值
		midiValue := root.MIDINote[0] + byte(interval)
		notes[i] = Note{
			Name:       root.Name,
			Octave:     root.Octave,
			Accidental: root.Accidental,
			MIDINote:   []byte{midiValue},
			Beat:       root.Beat,
		}
	}

	return NewChord(notes)
}