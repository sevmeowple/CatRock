package test

import (
    "catRock/pkg/core"
    "catRock/pkg/score"
    "testing"
)

// TestChannelAssignment 专门测试通道分配是否正确传递
func TestChannelAssignment(t *testing.T) {
    t.Log("=== 测试MIDI通道分配 ===")
    
    // 1. 创建测试音符（都是相同的C4，方便对比）
    testNote := core.NewNote(core.NewNoteParams{
        Name:       core.C,
        Octave:     4,
        Accidental: core.Natural,
        Beat:       core.Quarter,
    })
    
    // 2. 创建三个不同通道的轨道
    track1 := score.NewTrack("piano", core.AcousticGrandPiano, 1)    // 通道1 - 钢琴
    track1.AddElement(score.NewNoteElement(testNote))
    
    track2 := score.NewTrack("guitar", core.AcousticGuitarSteel, 2) // 通道2 - 吉他  
    track2.AddElement(score.NewNoteElement(testNote))
    
    track10 := score.NewTrack("drums", core.StandardDrumKit, 10)    // 通道10 - 鼓
    track10.AddElement(score.NewNoteElement(testNote))
    
    t.Logf("创建了3个轨道:")
    t.Logf("  - %s: 通道%d, 乐器%v", track1.Name, track1.Channel, track1.Instrument)
    t.Logf("  - %s: 通道%d, 乐器%v", track2.Name, track2.Channel, track2.Instrument)
    t.Logf("  - %s: 通道%d, 乐器%v", track10.Name, track10.Channel, track10.Instrument)
    
    // 3. 方案A：使用Section并行组织（推荐方案）
    t.Log("\n--- 方案A：Section并行组织 ---")
    sectionRoot := score.NewSection("all_tracks", []score.Playable{
        track1,
        track2, 
        track10,
    })
    
    scoreA := &score.Score{
        BPM:         120,
        Volume:      100,
        RootElement: sectionRoot,
    }
    
    eventsA, err := scoreA.Play()
    if err != nil {
        t.Fatalf("方案A生成事件失败: %v", err)
    }
    
    // 分析方案A的通道分布
    channelsA := analyzeChannels(eventsA)
    t.Logf("方案A结果: 使用通道 %v", getChannelNumbers(channelsA))
    
    // 验证方案A
    expectedChannels := []int{1, 2, 10}
    for _, expectedCh := range expectedChannels {
        if count, exists := channelsA[expectedCh]; !exists || count == 0 {
            t.Errorf("❌ 方案A缺少通道%d的事件", expectedCh)
        } else {
            t.Logf("✅ 方案A通道%d: %d个事件", expectedCh, count)
        }
    }
    
    // 4. 方案B：使用Track嵌套（当前有问题的方案）
    t.Log("\n--- 方案B：Track嵌套组织 ---")
    masterTrack := score.NewTrack("master", core.AcousticGrandPiano, 1)
    masterTrack.AddElement(track1)   // 嵌套track1
    masterTrack.AddElement(track2)   // 嵌套track2
    masterTrack.AddElement(track10)  // 嵌套track10
    
    scoreB := &score.Score{
        BPM:         120,
        Volume:      100,
        RootElement: masterTrack,
    }
    
    eventsB, err := scoreB.Play()
    if err != nil {
        t.Fatalf("方案B生成事件失败: %v", err)
    }
    
    // 分析方案B的通道分布
    channelsB := analyzeChannels(eventsB)
    t.Logf("方案B结果: 使用通道 %v", getChannelNumbers(channelsB))
    
    // 验证方案B（预期会失败，显示问题所在）
    for _, expectedCh := range expectedChannels {
        if count, exists := channelsB[expectedCh]; !exists || count == 0 {
            t.Logf("❌ 方案B缺少通道%d的事件 (这是预期的问题)", expectedCh)
        } else {
            t.Logf("✅ 方案B通道%d: %d个事件", expectedCh, count)
        }
    }
    
    // 5. 详细事件分析
    t.Log("\n--- 详细事件分析 ---")
    
    t.Log("方案A事件详情:")
    printEventDetails(t, eventsA, "A")
    
    t.Log("方案B事件详情:")
    printEventDetails(t, eventsB, "B")
    
    // 6. 结论
    t.Log("\n--- 测试结论 ---")
    if len(channelsA) >= 3 && len(channelsB) < 3 {
        t.Log("✅ 通道问题确认：方案A正确，方案B有问题")
        t.Log("📋 建议：使用Section而不是嵌套Track来组织并行播放")
    } else if len(channelsA) >= 3 && len(channelsB) >= 3 {
        t.Log("✅ 两种方案都正确，通道传递已修复")
    } else {
        t.Log("❌ 两种方案都有问题，需要进一步调试")
    }
}

// analyzeChannels 分析事件的通道分布
func analyzeChannels(events []score.Event) map[int]int {
    channels := make(map[int]int)
    for _, event := range events {
        if event.Action == score.NOTE_ON { // 只统计NOTE_ON事件避免重复
            channels[event.Channel]++
        }
    }
    return channels
}

// getChannelNumbers 获取使用的通道号列表
func getChannelNumbers(channelMap map[int]int) []int {
    channels := make([]int, 0, len(channelMap))
    for ch := range channelMap {
        channels = append(channels, ch)
    }
    
    // 简单排序
    for i := 0; i < len(channels)-1; i++ {
        for j := 0; j < len(channels)-i-1; j++ {
            if channels[j] > channels[j+1] {
                channels[j], channels[j+1] = channels[j+1], channels[j]
            }
        }
    }
    
    return channels
}

// printEventDetails 打印事件详情
func printEventDetails(t *testing.T, events []score.Event, label string) {
    noteOnCount := 0
    for _, event := range events {
        if event.Action == score.NOTE_ON {
            noteOnCount++
            if noteOnCount <= 5 { // 只显示前5个NOTE_ON事件
                t.Logf("  事件%s: 时间%.2f, 通道%d, 音符%v, 来源:%s", 
                    label, event.Time, event.Channel, event.Data, event.SourceTrack)
            }
        }
    }
    if noteOnCount > 5 {
        t.Logf("  ... 还有 %d 个NOTE_ON事件", noteOnCount-5)
    }
}

// TestSingleTrackChannel 测试单个轨道的通道设置
func TestSingleTrackChannel(t *testing.T) {
    t.Log("=== 测试单个轨道通道设置 ===")
    
    // 创建一个音符
    note := core.NewNote(core.NewNoteParams{
        Name: core.C, Octave: 4, Beat: core.Quarter,
    })
    
    // 测试不同通道设置
    testChannels := []int{1, 2, 5, 10, 16}
    
    for _, channel := range testChannels {
        t.Logf("\n--- 测试通道 %d ---", channel)
        
        track := score.NewTrack("test", core.AcousticGrandPiano, channel)
        track.AddElement(score.NewNoteElement(note))
        
        score := &score.Score{
            BPM: 120,
            Volume: 100,
            RootElement: track,
        }
        
        events, err := score.Play()
        if err != nil {
            t.Errorf("通道%d生成事件失败: %v", channel, err)
            continue
        }
        
        // 检查事件通道
        actualChannels := analyzeChannels(events)
        if len(actualChannels) != 1 {
            t.Errorf("❌ 通道%d: 预期1个通道，实际%d个: %v", 
                channel, len(actualChannels), getChannelNumbers(actualChannels))
        } else {
            actualChannel := getChannelNumbers(actualChannels)[0]
            if actualChannel == channel {
                t.Logf("✅ 通道%d: 正确", channel)
            } else {
                t.Errorf("❌ 通道%d: 预期通道%d，实际通道%d", 
                    channel, channel, actualChannel)
            }
        }
    }
}