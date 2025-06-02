package score

import (
    "catRock/pkg/core"
    "fmt"
)

// 音符播放元素
type NoteElement struct {
    ID   string
    Note core.Note
    
    // 可选覆盖设置
    VolumeOverride     *int
    InstrumentOverride *core.InstrumentID
    ChannelOverride    *int
}

var _ Element = (*NoteElement)(nil)

func (ne *NoteElement) GetID() string {
    if ne.ID != "" {
        return ne.ID
    }
    return fmt.Sprintf("note_%s_%d_%.2f", 
        ne.getNoteName(), ne.Note.Octave, float64(ne.Note.Beat))
}

func (ne *NoteElement) GetType() PlayableType {
    return NOTE_TYPE
}

func (ne *NoteElement) Duration(context PlayContext) float64 {
    return float64(ne.Note.Beat)
}

func (ne *NoteElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    velocity := ne.calculateVelocity(context)
    channel := ne.calculateChannel(context)
    midiNote := ne.Note.MIDINote[0]
    duration := ne.Duration(context)
    
    return []Event{
        {
            Time:          startTime,
            Duration:      duration,
            Type:          NOTE_EVENT,
            Action:        NOTE_ON,
            Data:          midiNote,
            Channel:       channel,
            Velocity:      velocity,
            SourceElement: ne.GetID(),
        },
        {
            Time:          startTime + duration,
            Duration:      0,
            Type:          NOTE_EVENT,
            Action:        NOTE_OFF,
            Data:          midiNote,
            Channel:       channel,
            Velocity:      0,
            SourceElement: ne.GetID(),
        },
    }
}

// 实现Element接口
func (ne *NoteElement) SetVolumeOverride(volume int) {
    ne.VolumeOverride = &volume
}

func (ne *NoteElement) SetInstrumentOverride(instrument core.InstrumentID) {
    ne.InstrumentOverride = &instrument
}

func (ne *NoteElement) SetChannelOverride(channel int) {
    ne.ChannelOverride = &channel
}

// 辅助方法
func (ne *NoteElement) calculateVelocity(context PlayContext) uint8 {
    if ne.VolumeOverride != nil {
        return uint8(*ne.VolumeOverride)
    }
    if ne.Note.Velocity > 0 {
        return ne.Note.Velocity
    }
    return uint8(context.CurrentVolume)
}

func (ne *NoteElement) calculateChannel(context PlayContext) int {
    if ne.ChannelOverride != nil {
        return *ne.ChannelOverride
    }
    if ne.Note.Channel > 0 {
        return int(ne.Note.Channel)
    }
    return context.CurrentChannel
}

func (ne *NoteElement) getNoteName() string {
    noteNames := []string{"C", "Cs", "D", "Ds", "E", "F", "Fs", "G", "Gs", "A", "As", "B"}
    return noteNames[int(ne.Note.Name)]
}

// 构造函数
func NewNoteElement(note core.Note) *NoteElement {
    return &NoteElement{Note: note}
}