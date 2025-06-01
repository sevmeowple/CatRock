package test

import (
    "catRock/pkg/core"
    "catRock/pkg/io"
    "catRock/pkg/io/midi"
    "catRock/pkg/score"
    "testing"
    "time"
)

func TestComplexScorePlayback(t *testing.T) {
    // 1. 创建基础音符
    noteC4 := core.NewNote(core.NewNoteParams{
        Name:       core.C,
        Octave:     4,
        Accidental: core.Natural,
        Beat:       core.Quarter, // 四分音符
    })
    
    noteE4 := core.NewNote(core.NewNoteParams{
        Name:       core.E,
        Octave:     4,
        Accidental: core.Natural,
        Beat:       core.Quarter,
    })
    
    noteG4 := core.NewNote(core.NewNoteParams{
        Name:       core.G,
        Octave:     4,
        Accidental: core.Natural,
        Beat:       core.Quarter,
    })
    
    // 2. 创建和弦（C大三和弦）
    cMajorChord := core.NewChordFromQuality(noteC4, core.Major)
    
    // 3. 创建 Playable 元素
    noteElementC := score.NewNoteElement(noteC4)
    noteElementE := score.NewNoteElement(noteE4)
    noteElementG := score.NewNoteElement(noteG4)
    chordElement := score.NewChordElement(cMajorChord)
    
    // 创建休止符（四分休止符）
    restElement := score.NewRestElement(core.Quarter)
    
    // 4. 创建旋律段落（顺序播放：C4 -> E4 -> G4 -> 休止符）
    melodySection := score.NewSection("melody", []score.Playable{
        noteElementC,
        noteElementE,
        noteElementG,
        restElement,
    })
    
    // 5. 创建和声轨道（并行播放：和弦重复2次）
    harmonyTrack := score.NewTrack("harmony", core.AcousticGuitarSteel, 1)
    harmonyTrack.AddElement(chordElement)
    harmonyTrack.AddElement(chordElement) // 重复播放同一和弦
    
    // 6. 创建节奏轨道（简单节拍）
    kickNote := core.NewNote(core.NewNoteParams{
        Name:       core.C,
        Octave:     2,
        Accidental: core.Natural,
        Beat:       core.Eighth, // 八分音符
    })
    kickElement := score.NewNoteElement(kickNote)
    
    rhythmTrack := score.NewTrack("rhythm", core.RoomDrumKit, 10) // 通道10通常是鼓
    rhythmTrack.AddElement(kickElement)
    rhythmTrack.AddElement(kickElement)
    rhythmTrack.AddElement(kickElement)
    rhythmTrack.AddElement(kickElement)
    
       // 7. 直接创建主轨道，不要嵌套太多层
    masterTrack := score.NewTrack("master", core.AcousticGrandPiano, 1)
    masterTrack.AddElement(melodySection)  // 旋律段落（顺序播放）
    
    // 8. 和声轨道设置正确的通道
    harmonyTrack.Channel = 2  // 🔧 设置不同的通道
    masterTrack.AddElement(harmonyTrack)   
    
    // 9. 节奏轨道保持通道10
    masterTrack.AddElement(rhythmTrack) 
    // 9. 创建完整的 Score
    complexScore := &score.Score{
        Title:    "Complex Test Song",
        Composer: "AI Test Suite",
        Year:     2024,
        BPM:      120,
        TimeSignature: score.TimeSignature{
            Numerator:   4,
            Denominator: 4,
        },
        KeySignature: score.KeySignature{
            Tonic: core.C,
            Mode:  core.Natural,
        },
        Volume:      100,
        RootElement: masterTrack,
    }
    
    // 10. 生成事件序列
    events, err := complexScore.Play()
    if err != nil {
        t.Fatalf("生成事件失败: %v", err)
    }
    
    // 11. 验证事件数量和结构
    t.Logf("生成了 %d 个事件", len(events))
    t.Logf("预期播放时长: %.2f 拍", complexScore.Duration())
    
    // 分析事件分布
    eventsByType := make(map[score.EventType]int)
    eventsByAction := make(map[score.EventAction]int)
    eventsByChannel := make(map[int]int)
    
    for _, event := range events {
        eventsByType[event.Type]++
        eventsByAction[event.Action]++
        eventsByChannel[event.Channel]++
        
        // 打印前10个事件的详细信息
        if len(events) <= 10 || eventsByType[event.Type] <= 5 {
            t.Logf("事件: 时间=%.2f, 类型=%v, 动作=%v, 通道=%d, 数据=%v",
                event.Time, event.Type, event.Action, event.Channel, event.Data)
        }
    }
    
    // 12. 验证事件正确性
    if eventsByAction[score.NOTE_ON] == 0 {
        t.Error("没有生成 NOTE_ON 事件")
    }
    
    if eventsByAction[score.NOTE_OFF] == 0 {
        t.Error("没有生成 NOTE_OFF 事件")
    }
    
    // NOTE_ON 和 NOTE_OFF 数量应该相等
    if eventsByAction[score.NOTE_ON] != eventsByAction[score.NOTE_OFF] {
        t.Errorf("NOTE_ON (%d) 和 NOTE_OFF (%d) 数量不匹配",
            eventsByAction[score.NOTE_ON], eventsByAction[score.NOTE_OFF])
    }
    
    // 应该有多个通道的事件
    if len(eventsByChannel) < 2 {
        t.Logf("警告: 只使用了 %d 个通道，预期至少2个", len(eventsByChannel))
    }
    
    t.Logf("事件分析:")
    t.Logf("  - NOTE_ON: %d", eventsByAction[score.NOTE_ON])
    t.Logf("  - NOTE_OFF: %d", eventsByAction[score.NOTE_OFF])
    t.Logf("  - 使用的通道: %v", getKeys(eventsByChannel))
    
    // 13. 尝试实际MIDI播放（可选）
    if testing.Short() {
        t.Skip("跳过MIDI播放测试（使用 -test.short）")
        return
    }
    
    midiPlayer := midi.NewMIDIPlayer()
    status, err := midiPlayer.Connect()
    if err != nil || status != io.Connected {
        t.Skip("MIDI不可用，跳过实际播放测试")
        return
    }
    defer midiPlayer.Disconnect()
    
    t.Log("开始MIDI播放测试...")
    
    // 创建播放引擎
    playEngine := score.NewPlayEngine(midiPlayer, complexScore.BPM)
    
    // 播放事件（实际播放可能需要几秒钟）
    start := time.Now()
    err = playEngine.PlayEvents(events)
    duration := time.Since(start)
    
    if err != nil {
        t.Errorf("播放事件失败: %v", err)
    } else {
        t.Logf("播放完成，实际用时: %v", duration)
    }

	    // 🔧 添加：详细的事件时间分析
    t.Logf("详细事件分析:")
    eventsByTime := make(map[float64][]score.Event)
    for _, event := range events {
        eventsByTime[event.Time] = append(eventsByTime[event.Time], event)
    }
    
    for time, timeEvents := range eventsByTime {
        t.Logf("时间 %.2f: %d 个事件", time, len(timeEvents))
        for _, event := range timeEvents {
            if len(timeEvents) <= 5 { // 只显示少量事件的详情
                t.Logf("  - 通道%d, 动作%v, 数据%v, 来源%s", 
                    event.Channel, event.Action, event.Data, event.SourceTrack)
            }
        }
    }
}

// 测试单独的Section和Track行为
func TestSectionAndTrackBehavior(t *testing.T) {
    // 测试Section的顺序播放
    note1 := score.NewNoteElement(core.NewNote(core.NewNoteParams{
        Name: core.C, Octave: 4, Beat: core.Quarter,
    }))
    note2 := score.NewNoteElement(core.NewNote(core.NewNoteParams{
        Name: core.D, Octave: 4, Beat: core.Quarter,
    }))
    
    section := score.NewSection("test_section", []score.Playable{note1, note2})
    
    // 测试Section的持续时间（应该是两个音符的时间之和）
    context := score.NewPlayContext(&score.Score{BPM: 120, Volume: 100})
    sectionDuration := section.Duration(context)
    expectedDuration := 2.0 // 两个四分音符
    
    if sectionDuration != expectedDuration {
        t.Errorf("Section持续时间错误: 得到 %.2f, 期望 %.2f", sectionDuration, expectedDuration)
    }
    
    // 测试事件生成
    events := section.GenerateEvents(0.0, context)
    
    // 验证事件时间（第二个音符应该在第一个音符结束后开始）
    var note1Events, note2Events []score.Event
    for _, event := range events {
        if event.Data == note1.Note.MIDINote[0] {
            note1Events = append(note1Events, event)
        } else if event.Data == note2.Note.MIDINote[0] {
            note2Events = append(note2Events, event)
        }
    }
    
    if len(note1Events) == 0 || len(note2Events) == 0 {
        t.Fatal("未生成预期的音符事件")
    }
    
    // 第一个音符从时间0开始
    note1StartTime := findNoteOnTime(note1Events)
    if note1StartTime != 0.0 {
        t.Errorf("第一个音符开始时间错误: %.2f, 期望 0.0", note1StartTime)
    }
    
    // 第二个音符应该在第一个音符结束后开始
    note2StartTime := findNoteOnTime(note2Events)
    if note2StartTime != 1.0 { // 四分音符 = 1.0拍
        t.Errorf("第二个音符开始时间错误: %.2f, 期望 1.0", note2StartTime)
    }
    
    t.Logf("Section顺序播放验证通过")
}

// 测试Track的并行播放
func TestTrackParallelBehavior(t *testing.T) {
    note1 := score.NewNoteElement(core.NewNote(core.NewNoteParams{
        Name: core.C, Octave: 4, Beat: core.Quarter,
    }))
    note2 := score.NewNoteElement(core.NewNote(core.NewNoteParams{
        Name: core.E, Octave: 4, Beat: core.Half, // 二分音符，更长的持续时间
    }))
    
    track := score.NewTrack("test_track", core.AcousticGrandPiano, 1)
    track.AddElement(note1)
    track.AddElement(note2)
    
    context := score.NewPlayContext(&score.Score{BPM: 120, Volume: 100})
    
    // Track的持续时间应该是最长元素的时间
    trackDuration := track.Duration(context)
    expectedDuration := 2.0 // 二分音符的长度
    
    if trackDuration != expectedDuration {
        t.Errorf("Track持续时间错误: 得到 %.2f, 期望 %.2f", trackDuration, expectedDuration)
    }
    
    // 测试并行播放：两个音符应该同时开始
    events := track.GenerateEvents(0.0, context)
    
    var note1Events, note2Events []score.Event
    for _, event := range events {
        if event.Data == note1.Note.MIDINote[0] {
            note1Events = append(note1Events, event)
        } else if event.Data == note2.Note.MIDINote[0] {
            note2Events = append(note2Events, event)
        }
    }
    
    note1StartTime := findNoteOnTime(note1Events)
    note2StartTime := findNoteOnTime(note2Events)
    
    if note1StartTime != note2StartTime {
        t.Errorf("Track中的音符未并行播放: note1=%.2f, note2=%.2f", 
            note1StartTime, note2StartTime)
    }
    
    if note1StartTime != 0.0 {
        t.Errorf("Track开始时间错误: %.2f, 期望 0.0", note1StartTime)
    }
    
    t.Logf("Track并行播放验证通过")
}

// 辅助函数：获取map的键
func getKeys(m map[int]int) []int {
    keys := make([]int, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

// 辅助函数：找到NOTE_ON事件的时间
func findNoteOnTime(events []score.Event) float64 {
    for _, event := range events {
        if event.Action == score.NOTE_ON {
            return event.Time
        }
    }
    return -1 // 未找到
}

func TestAvailableInstruments(t *testing.T) {
    midiPlayer := midi.NewMIDIPlayer()
    status, err := midiPlayer.Connect()
    if err != nil || status != io.Connected {
        t.Skip("MIDI不可用")
        return
    }
    defer midiPlayer.Disconnect()
    
    // 测试一些常见乐器
    testInstruments := []struct{
        name string
        program uint8
    }{
        {"钢琴", 0},
        {"电钢琴", 4},
        {"吉他", 24},
        {"贝斯", 32},
        {"小提琴", 40},
        {"长笛", 73},
        {"合成器", 80},
		{"drum kit", 118}, // 通常是鼓组
    }
    
    for _, inst := range testInstruments {
        t.Logf("测试乐器: %s (程序号 %d)", inst.name, inst.program)
        
        // 设置乐器
        err = midiPlayer.SendProgramChange(1, inst.program)
        if err == nil {
            // 播放测试音符
            midiPlayer.SendNoteOn(1, 60, 80)
            time.Sleep(300 * time.Millisecond)
            midiPlayer.SendNoteOff(1, 60, 0)
            time.Sleep(100 * time.Millisecond)
        } else {
            t.Logf("  ❌ 设置失败: %v", err)
        }
    }
}

