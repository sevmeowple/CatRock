package test

import (
	"catRock/pkg/core"
	"catRock/pkg/io"
	"catRock/pkg/io/midi"
	"catRock/pkg/score"
	"fmt"
	"testing"
)

func TestNotePlay(t *testing.T) {
	note := core.NewNote(core.NewNoteParams{
		Name:       core.C,
		Octave:     4,
		Accidental: core.Natural,
	})
	player := midi.NewMIDIPlayer()

	// 尝试连接，如果失败则跳过测试
	status, err := player.Connect()
	if err != nil {
		t.Skipf("MIDI connection failed: %v", err)
		return
	}

	if status != io.Connected {
		t.Skipf("MIDI not connected, status: %v", status)
		return
	}

	defer player.Disconnect()


	note.SetBeat(core.Eighth)
	// 测试播放音符
	err = player.PlayNote(io.PlayNoteParams{
		Note: note,
		Bpm:  120, // 设置节拍为120 BPM
	})
	note_2 := core.NewNote(core.NewNoteParams{
		Name:       core.Ds,
		Octave:     4,
		Accidental: core.Natural,
		Beat: 	 core.Quarter,
	})

	err = player.PlayNote(io.PlayNoteParams{
		Note: note_2,
		Bpm:  120, // 设置节拍为120 BPM
	})

	if err != nil {
		t.Errorf("PlayNote failed: %v", err)
	}

	t.Logf("Successfully played note: %v", note)
}


func TestPlayChord(t *testing.T) {
	note := core.NewNote(core.NewNoteParams{
		Name:       core.C,
		Octave:    4,
		Accidental: core.Natural,
	})
	player := midi.NewMIDIPlayer()
	player.Connect()
	chord := core.NewChordFromQuality(note, core.Major)
	err := player.PlayChord(io.PlayChordParams{
		Chord: chord,
		Bpm:   120,
	})
	if err != nil {
		t.Errorf("PlayChord failed: %v", err)
	}

	t.Logf("Successfully played chord: %v", chord)
}

// 测试代码（可以放在单独的测试文件中）
func TestBasicPlayback(t *testing.T) {
    // 创建一个简单的音符
    note := core.Note{
        Name:     core.C,
        Octave:   4,
        Beat:     1.0, // 四分音符
        Velocity: 80,
        MIDINote: []uint8{60}, // C4
    }
    
    // 创建 NoteElement
    noteElement := score.NewNoteElement(note)

    // 创建 Score
    score := &score.Score{
        BPM:         120,
        Volume:      100,
        RootElement: noteElement,
    }
    
    // 播放测试
    events, err := score.Play()
    if err != nil {
        fmt.Printf("播放错误: %v\n", err)
        return
    }
    
    fmt.Printf("生成了 %d 个事件\n", len(events))
    for _, event := range events {
        fmt.Printf("时间: %.2f, 动作: %v, 数据: %v\n", 
            event.Time, event.Action, event.Data)
    }
}