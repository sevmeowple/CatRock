package score

import (
    "catRock/pkg/core"
    "fmt"
)

// 和弦播放元素
type ChordElement struct {
    ID    string
    Chord core.Chord
    
    // 可选覆盖设置
    VolumeOverride     *int
    InstrumentOverride *core.InstrumentID
    ChannelOverride    *int
}

var _ Element = (*ChordElement)(nil)

func (ce *ChordElement) GetID() string {
    if ce.ID != "" {
        return ce.ID
    }
    // 修正：根据和弦内容生成ID
    if len(ce.Chord.Notes) > 0 {
        firstNote := ce.Chord.Notes[0]
        return fmt.Sprintf("chord_%s%d_%.2f", 
            firstNote.Name.String(), firstNote.Octave, float64(firstNote.Beat))
    }
    return "chord_empty"
}

func (ce *ChordElement) GetType() PlayableType {
    return CHORD_TYPE
}

func (ce *ChordElement) Duration(context PlayContext) float64 {
    // 修正：使用和弦中第一个音符的节拍，或默认值
    if len(ce.Chord.Notes) > 0 {
        return float64(ce.Chord.Notes[0].Beat)
    }
    return 1.0 // 默认一拍
}

func (ce *ChordElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    velocity := ce.calculateVelocity(context)
    channel := ce.calculateChannel(context)
    duration := ce.Duration(context)
    
    // 修正：为和弦中每个音符生成事件
    for _, note := range ce.Chord.Notes {
        if len(note.MIDINote) == 0 {
            continue // 跳过无效音符
        }
        
        midiNote := note.MIDINote[0]
        
        // NOTE_ON 事件
        events = append(events, Event{
            Time:          startTime,
            Duration:      duration,
            Type:          CHORD_EVENT, // 修正：使用CHORD_EVENT
            Action:        NOTE_ON,
            Data:          midiNote,
            Channel:       channel,
            Velocity:      velocity,
            SourceElement: ce.GetID(), // 修正：使用ce而不是ne
        })
        
        // NOTE_OFF 事件
        events = append(events, Event{
            Time:          startTime + duration,
            Duration:      0,
            Type:          CHORD_EVENT,
            Action:        NOTE_OFF,
            Data:          midiNote,
            Channel:       channel,
            Velocity:      0,
            SourceElement: ce.GetID(), // 修正：使用ce而不是ne
        })
    }
    
    return events
}

// 实现Element接口
func (ce *ChordElement) SetVolumeOverride(volume int) {
    ce.VolumeOverride = &volume
}

func (ce *ChordElement) SetInstrumentOverride(instrument core.InstrumentID) {
    ce.InstrumentOverride = &instrument
}

func (ce *ChordElement) SetChannelOverride(channel int) {
    ce.ChannelOverride = &channel
}

// 辅助方法
func (ce *ChordElement) calculateVelocity(context PlayContext) uint8 {
    if ce.VolumeOverride != nil {
        return uint8(*ce.VolumeOverride)
    }
    
    // 从和弦中第一个音符获取velocity
    if len(ce.Chord.Notes) > 0 && ce.Chord.Notes[0].Velocity > 0 {
        return ce.Chord.Notes[0].Velocity
    }
    
    return uint8(context.CurrentVolume)
}

func (ce *ChordElement) calculateChannel(context PlayContext) int {
    if ce.ChannelOverride != nil {
        return *ce.ChannelOverride
    }
    
    // 从和弦中第一个音符获取channel
    if len(ce.Chord.Notes) > 0 && ce.Chord.Notes[0].Channel > 0 {
        return int(ce.Chord.Notes[0].Channel)
    }
    
    return context.CurrentChannel
}

// 构造函数
func NewChordElement(chord core.Chord) *ChordElement {
    return &ChordElement{Chord: chord}
}

// 便捷构造函数 - 从根音和和弦质量创建
func NewChordElementFromQuality(root core.Note, quality core.ChordQuality) *ChordElement {
    chord := core.NewChordFromQuality(root, quality)
    return &ChordElement{Chord: chord}
}