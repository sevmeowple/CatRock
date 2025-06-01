package score

import (
	"catRock/pkg/core"
	"fmt"
	"sort"
)

// NoteElement 表示单个音符的播放元素
type NoteElement struct {
    ID   string    // 唯一标识符
    Note core.Note // 底层音符数据
    
    // 可选的覆盖设置
    VolumeOverride     *int     // 音量覆盖
    InstrumentOverride *core.InstrumentID // 乐器覆盖
    ChannelOverride    *int     // 通道覆盖
}

var _ Playable = (*NoteElement)(nil) // 确保 NoteElement 实现了 Playable 接口

// 实现 Playable 接口
func (ne *NoteElement) GetID() string {
    if ne.ID != "" {
        return ne.ID
    }
    // 自动生成ID：音符名_八度_时值
    return fmt.Sprintf("note_%s_%d_%.2f", 
        ne.getNoteName(), ne.Note.Octave, float64(ne.Note.Beat))
}

func (ne *NoteElement) GetType() PlayableType {
    return NOTE_TYPE
}

func (ne *NoteElement) Duration(context PlayContext) float64 {
    // 直接返回音符的节拍值（以四分音符为单位）
    return float64(ne.Note.Beat)
}

func (ne *NoteElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    // 计算实际播放参数
    velocity := ne.calculateVelocity(context)
    channel := ne.calculateChannel(context)
    midiNote := ne.Note.MIDINote[0]
    duration := ne.Duration(context)
    
    // 生成 NOTE_ON 和 NOTE_OFF 事件
    events := []Event{
        // NOTE_ON 事件
        {
            Time:     startTime,
            Duration: duration,
            Type:     NOTE_EVENT,
            Action:   NOTE_ON,
            Data:     midiNote,
            Channel:  channel,
            Velocity: velocity,
            SourceElement: ne.GetID(),
        },
        // NOTE_OFF 事件
        {
            Time:     startTime + duration,
            Duration: 0,
            Type:     NOTE_EVENT,
            Action:   NOTE_OFF,
            Data:     midiNote,
            Channel:  channel,
            Velocity: 0, // NOTE_OFF 通常速度为0
            SourceElement: ne.GetID(),
        },
    }
    
    return events
}

// 辅助方法：计算实际音量
func (ne *NoteElement) calculateVelocity(context PlayContext) uint8 {
    // 优先级：元素覆盖 > 音符自身 > 上下文
    if ne.VolumeOverride != nil {
        return uint8(*ne.VolumeOverride)
    }
    
    if ne.Note.Velocity > 0 {
        return ne.Note.Velocity
    }
    
    return uint8(context.CurrentVolume)
}

// 辅助方法：计算实际通道
func (ne *NoteElement) calculateChannel(context PlayContext) int {
    if ne.ChannelOverride != nil {
        return *ne.ChannelOverride
    }
    
    if ne.Note.Channel > 0 {
        return int(ne.Note.Channel)
    }
    
    return context.CurrentChannel
}

// 辅助方法：获取音符名称字符串
func (ne *NoteElement) getNoteName() string {
    noteNames := []string{"C", "Cs", "D", "Ds", "E", "F", "Fs", "G", "Gs", "A", "As", "B"}
    return noteNames[int(ne.Note.Name)]
}

// 构造函数
func NewNoteElement(note core.Note) *NoteElement {
    return &NoteElement{
        Note: note,
    }
}

// 带配置的构造函数
func NewNoteElementWithConfig(note core.Note, config NoteElementConfig) *NoteElement {
    ne := &NoteElement{
        ID:   config.ID,
        Note: note,
    }
    
    if config.VolumeOverride != nil {
        ne.VolumeOverride = config.VolumeOverride
    }
    
    if config.InstrumentOverride != nil {
        ne.InstrumentOverride = config.InstrumentOverride
    }
    
    if config.ChannelOverride != nil {
        ne.ChannelOverride = config.ChannelOverride
    }
    
    return ne
}

// 配置结构体
type NoteElementConfig struct {
    ID                 string
    VolumeOverride     *int
    InstrumentOverride *core.InstrumentID
    ChannelOverride    *int
}

// ...existing code...

// ChordElement 表示和弦的播放元素
type ChordElement struct {
    ID    string     // 唯一标识符
    Chord core.Chord // 底层和弦数据
    
    // 可选的覆盖设置
    VolumeOverride     *int               // 音量覆盖
    InstrumentOverride *core.InstrumentID // 乐器覆盖
    ChannelOverride    *int               // 通道覆盖
}

var _ Playable = (*ChordElement)(nil)

func (ce *ChordElement) GetID() string {
    if ce.ID != "" {
        return ce.ID
    }
    // 自动生成ID：chord_根音_音符数_时值
    if len(ce.Chord.Notes) > 0 {
        root := ce.Chord.Notes[0]
        return fmt.Sprintf("chord_%s_%d_%d_%.2f", 
            ce.getNoteName(root.Name), len(ce.Chord.Notes), root.Octave, float64(root.Beat))
    }
    return "chord_empty"
}

func (ce *ChordElement) GetType() PlayableType {
    return CHORD_TYPE
}

func (ce *ChordElement) Duration(context PlayContext) float64 {
    if len(ce.Chord.Notes) > 0 {
        return float64(ce.Chord.Notes[0].Beat)
    }
    return 1.0 // 默认四分音符
}

func (ce *ChordElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    duration := ce.Duration(context)
    
    // 为每个音符生成事件
    for _, note := range ce.Chord.Notes {
        velocity := ce.calculateVelocity(note, context)
        channel := ce.calculateChannel(note, context)
        midiNote := note.MIDINote[0]
        
        // NOTE_ON 事件
        events = append(events, Event{
            Time:          startTime,
            Duration:      duration,
            Type:          CHORD_EVENT,
            Action:        NOTE_ON,
            Data:          midiNote,
            Channel:       channel,
            Velocity:      velocity,
            SourceElement: ce.GetID(),
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
            SourceElement: ce.GetID(),
        })
    }
    
    return events
}

func (ce *ChordElement) calculateVelocity(note core.Note, context PlayContext) uint8 {
    if ce.VolumeOverride != nil {
        return uint8(*ce.VolumeOverride)
    }
    if note.Velocity > 0 {
        return note.Velocity
    }
    return uint8(context.CurrentVolume)
}

func (ce *ChordElement) calculateChannel(note core.Note, context PlayContext) int {
    if ce.ChannelOverride != nil {
        return *ce.ChannelOverride
    }
    if note.Channel > 0 {
        return int(note.Channel)
    }
    return context.CurrentChannel
}

func (ce *ChordElement) getNoteName(name core.BaseNoteName) string {
    noteNames := []string{"C", "Cs", "D", "Ds", "E", "F", "Fs", "G", "Gs", "A", "As", "B"}
    return noteNames[int(name)]
}

// RestElement 表示休止符的播放元素
type RestElement struct {
    ID       string          // 唯一标识符
    Duration_ core.BeatValue  // 休止时长
}

var _ Playable = (*RestElement)(nil)

func (re *RestElement) GetID() string {
    if re.ID != "" {
        return re.ID
    }
    return fmt.Sprintf("rest_%.2f", float64(re.Duration_))
}

func (re *RestElement) GetType() PlayableType {
    return REST_TYPE
}

func (re *RestElement) Duration(context PlayContext) float64 {
    return float64(re.Duration_)
}

func (re *RestElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    // 休止符不生成任何事件，只消耗时间
    return []Event{}
}

// Section 表示顺序播放的音乐段落
type Section struct {
    ID       string     // 唯一标识符
    Name     string     // 段落名称
    Elements []Playable // 子元素（顺序播放）
    
    // 段落级别设置
    BPMOverride        *float64 // BPM覆盖
    VolumeOverride     *int     // 音量覆盖
    InstrumentOverride *core.InstrumentID // 乐器覆盖
    ChannelOverride    *int     // 通道覆盖
}

var _ Playable = (*Section)(nil)

func (s *Section) GetID() string {
    if s.ID != "" {
        return s.ID
    }
    if s.Name != "" {
        return fmt.Sprintf("section_%s", s.Name)
    }
    return fmt.Sprintf("section_%d_elements", len(s.Elements))
}

func (s *Section) GetType() PlayableType {
    return SECTION_TYPE
}

func (s *Section) Duration(context PlayContext) float64 {
    // Section持续时间 = 所有元素时间的累加
    totalDuration := 0.0
    sectionContext := s.createSectionContext(context)
    
    for _, element := range s.Elements {
        totalDuration += element.Duration(sectionContext)
    }
    
    return totalDuration
}

func (s *Section) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    currentTime := startTime
    sectionContext := s.createSectionContext(context)
    
    // 防止循环引用
    for _, stackID := range context.ElementStack {
        if stackID == s.GetID() {
            return events // 检测到循环，返回空事件
        }
    }
    
    // 添加当前元素到栈
    sectionContext.ElementStack = append(context.ElementStack, s.GetID())
    
    // 顺序播放所有子元素
    for _, element := range s.Elements {
        elementEvents := element.GenerateEvents(currentTime, sectionContext)
        events = append(events, elementEvents...)
        
        // 时间推进：下一个元素在当前元素结束后开始
        elementDuration := element.Duration(sectionContext)
        currentTime += elementDuration
    }
    
    return events
}

func (s *Section) createSectionContext(parentContext PlayContext) PlayContext {
    context := parentContext
    
    // 应用Section级别的覆盖
    if s.BPMOverride != nil {
        context.CurrentBPM = *s.BPMOverride
    }
    if s.VolumeOverride != nil {
        context.CurrentVolume = *s.VolumeOverride
    }
    if s.InstrumentOverride != nil {
        context.CurrentInstrument = string(*s.InstrumentOverride)
    }
    if s.ChannelOverride != nil {
        context.CurrentChannel = *s.ChannelOverride
    }
    
    return context
}

// Track 表示并行播放的音乐轨道
type Track struct {
    ID       string     // 唯一标识符
    Name     string     // 轨道名称
    Elements []Playable // 子元素（并行播放）
    
    // 轨道级别设置
    Instrument    core.InstrumentID // 轨道乐器
    Channel       int              // MIDI通道
    Volume        int              // 轨道音量
    StartOffset   float64          // 起始偏移时间（@4的实现）
    
    // 可选覆盖
    BPMOverride *float64 // BPM覆盖
}

var _ Playable = (*Track)(nil)

func (t *Track) GetID() string {
    if t.ID != "" {
        return t.ID
    }
    if t.Name != "" {
        return fmt.Sprintf("track_%s", t.Name)
    }
    return fmt.Sprintf("track_%d_ch%d", int(t.Instrument), t.Channel)
}

func (t *Track) GetType() PlayableType {
    return TRACK_TYPE
}

func (t *Track) Duration(context PlayContext) float64 {
    // Track持续时间 = 最长元素的持续时间
    maxDuration := 0.0
    trackContext := t.createTrackContext(context)
    
    for _, element := range t.Elements {
        elementDuration := element.Duration(trackContext)
        if elementDuration > maxDuration {
            maxDuration = elementDuration
        }
    }
    
    return maxDuration
}

func (t *Track) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    actualStartTime := startTime + t.StartOffset
    trackContext := t.createTrackContext(context)
    
    // 防止循环引用
    for _, stackID := range context.ElementStack {
        if stackID == t.GetID() {
            return events
        }
    }
    
    // 添加当前元素到栈
    trackContext.ElementStack = append(context.ElementStack, t.GetID())
    
    // 并行播放：所有元素从同一时间开始
    for _, element := range t.Elements {
        elementEvents := element.GenerateEvents(actualStartTime, trackContext)
        
        // 🔧 关键修复：正确处理子Track的通道
        for i := range elementEvents {
            // 如果子元素是Track，保持它自己的通道设置
            if childTrack, ok := element.(*Track); ok {
                elementEvents[i].Channel = childTrack.Channel  // 使用子Track的通道
            } else {
                elementEvents[i].Channel = t.Channel  // 使用当前Track的通道
            }
            elementEvents[i].SourceTrack = t.Name
        }
        
        events = append(events, elementEvents...)
    }
    
    return t.sortEventsByTime(events)
}

func (t *Track) createTrackContext(parentContext PlayContext) PlayContext {
    context := parentContext
    
    // 🔧 确保通道设置正确应用
    context.CurrentInstrument = string(t.Instrument)
    context.CurrentChannel = t.Channel  // 这是关键！
    context.CurrentVolume = t.Volume
    
    if t.BPMOverride != nil {
        context.CurrentBPM = *t.BPMOverride
    }
    
    return context
}

func (t *Track) sortEventsByTime(events []Event) []Event {
    // 使用标准库的高效排序
    sort.Slice(events, func(i, j int) bool {
        if events[i].Time == events[j].Time {
            // 同一时间：NOTE_ON 优先于 NOTE_OFF
            if events[i].Action == NOTE_ON && events[j].Action == NOTE_OFF {
                return true
            }
            if events[i].Action == NOTE_OFF && events[j].Action == NOTE_ON {
                return false
            }
        }
        return events[i].Time < events[j].Time
    })
    return events
}

// 构造函数们
func NewChordElement(chord core.Chord) *ChordElement {
    return &ChordElement{Chord: chord}
}

func NewRestElement(duration core.BeatValue) *RestElement {
    return &RestElement{Duration_: duration}
}

func NewSection(name string, elements []Playable) *Section {
    return &Section{
        Name:     name,
        Elements: elements,
    }
}

func NewTrack(name string, instrument core.InstrumentID, channel int) *Track {
    return &Track{
        Name:       name,
        Instrument: instrument,
        Channel:    channel,
        Volume:     100, // 默认音量
        Elements:   []Playable{},
    }
}

// Track的便捷方法
func (t *Track) AddElement(element Playable) {
    t.Elements = append(t.Elements, element)
}

func (t *Track) SetStartOffset(offset float64) {
    t.StartOffset = offset
}

// Section的便捷方法
func (s *Section) AddElement(element Playable) {
    s.Elements = append(s.Elements, element)
}