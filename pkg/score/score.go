package score

import (
	"catRock/pkg/core"
	io "catRock/pkg/io"
	"fmt"
	"sort"
	"time"
)

type Score struct {
	// 元数据
	Title    string
	Composer string
	Year     int

	// 全局设置
	BPM           float64
	TimeSignature TimeSignature
	KeySignature  KeySignature
	Volume        int

	// 核心：音乐内容
	RootElement Playable // 根级播放元素

	// 预定义的可复用元素
	Definitions map[string]Playable // let定义的元素
}

type TimeSignature struct {
	Numerator   int // 分子，表示每小节的拍数
	Denominator int // 分母，表示每拍的时值（如4表示四分音符）
}

type KeySignature struct {
	Tonic core.BaseNoteName // 主音，表示调性的主音
	Mode  core.Accidental   // 模式，表示调性的模式（大调、小调等）
}

type Playable interface {
	GenerateEvents(startTime float64, context PlayContext) []Event
	Duration(context PlayContext) float64
	GetID() string
	GetType() PlayableType
}

type PlayableType int

const (
	NOTE_TYPE PlayableType = iota
	CHORD_TYPE
	REST_TYPE
	TRACK_TYPE
	SECTION_TYPE
)

type PlayContext struct {
	// 全局设置（来自Score）
	GlobalBPM           float64
	GlobalTimeSignature TimeSignature
	GlobalKeySignature  KeySignature
	GlobalVolume        int

	// 当前生效设置（经过层层覆盖后的最终值）
	CurrentBPM        float64
	CurrentVolume     int
	CurrentInstrument string
	CurrentChannel    int

	// 播放状态
	StartTime   float64 // 播放开始时间
	CurrentTime float64 // 当前播放时间

	// 循环检测（避免无限递归）
	ElementStack []string // 正在播放的元素ID栈
}
type Event struct {
	// 时间信息
	Time     float64 // 绝对播放时间（拍数）
	Duration float64 // 事件持续时间（用于NOTE_OFF计算）

	// 事件类型和动作
	Type   EventType   // 事件类型
	Action EventAction // 具体动作

	// 音乐数据
	Data interface{} // 具体的音乐数据

	// 播放属性
	Channel  int   // MIDI通道
	Velocity uint8 // 力度/音量

	// 来源信息
	SourceTrack   string // 来源轨道名称
	SourceElement string // 来源元素ID
}

type EventType int

const (
	NOTE_EVENT    EventType = iota // 音符事件
	CHORD_EVENT                    // 和弦事件
	CONTROL_EVENT                  // 控制事件
	PROGRAM_EVENT                  // 乐器切换事件
	TEMPO_EVENT                    // 速度变化事件
)

type EventAction int

const (
	// 音符动作
	NOTE_ON EventAction = iota
	NOTE_OFF

	// 控制动作
	VOLUME_CHANGE
	PAN_CHANGE
	PITCH_BEND

	// 程序动作
	PROGRAM_CHANGE
	TEMPO_CHANGE
)

func NewPlayContext(score *Score) PlayContext {
    return PlayContext{
        // 全局设置
        GlobalBPM:           score.BPM,
        GlobalTimeSignature: score.TimeSignature,
        GlobalKeySignature:  score.KeySignature,
        GlobalVolume:        score.Volume,
        
        // 当前设置初始化为全局设置
        CurrentBPM:        score.BPM,
        CurrentVolume:     score.Volume,
        CurrentInstrument: "piano", // 默认乐器
        CurrentChannel:    1,       // 默认通道
        
        // 播放状态
        StartTime:   0.0,
        CurrentTime: 0.0,
        
        // 空的元素栈
        ElementStack: []string{},
    }
}

// PlayContext 的上下文方法
func (pc PlayContext) WithTrackSettings(track *Track) PlayContext {
    newContext := pc
    newContext.CurrentInstrument = string(track.Instrument)
    newContext.CurrentChannel = track.Channel
    newContext.CurrentVolume = track.Volume
    if track.BPMOverride != nil {
        newContext.CurrentBPM = *track.BPMOverride
    }
    return newContext
}

func (pc PlayContext) WithSectionSettings(section *Section) PlayContext {
    newContext := pc
    if section.BPMOverride != nil {
        newContext.CurrentBPM = *section.BPMOverride
    }
    if section.VolumeOverride != nil {
        newContext.CurrentVolume = *section.VolumeOverride
    }
    if section.InstrumentOverride != nil {
        newContext.CurrentInstrument = string(*section.InstrumentOverride)
    }
    if section.ChannelOverride != nil {
        newContext.CurrentChannel = *section.ChannelOverride
    }
    return newContext
}

func (s *Score) Play() ([]Event, error) {
    if s.RootElement == nil {
        return nil, fmt.Errorf("score has no root element to play")
    }
    
    // 1. 创建播放上下文
    context := NewPlayContext(s)
    
    // 2. 生成事件序列
    events := s.RootElement.GenerateEvents(0.0, context)
    
    // 3. 按时间排序
    sort.Slice(events, func(i, j int) bool {
        return events[i].Time < events[j].Time
    })
    
    return events, nil
}

// 获取总播放时长
func (s *Score) Duration() float64 {
    if s.RootElement == nil {
        return 0.0
    }
    
    context := NewPlayContext(s)
    return s.RootElement.Duration(context)
}


type PlayEngine struct {
    midiPlayer io.IO
    bpm        float64
}

func NewPlayEngine(midiPlayer io.IO, bpm float64) *PlayEngine {
    return &PlayEngine{
        midiPlayer: midiPlayer,
        bpm:        bpm,
    }
}

func (pe *PlayEngine) PlayEvents(events []Event) error {
    if len(events) == 0 {
        return nil
    }
    
    startTime := time.Now()
    
    for _, event := range events {
        // 等待到正确时间
        eventTime := pe.beatsToTime(event.Time)
        targetTime := startTime.Add(eventTime)
        time.Sleep(time.Until(targetTime))
        
        // 直接调用io包的低级方法
        err := pe.executeEvent(event)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func (pe *PlayEngine) executeEvent(event Event) error {
    switch event.Action {
    case NOTE_ON:
        midiNote := event.Data.(uint8)
        return pe.midiPlayer.SendNoteOn(
            uint8(event.Channel), 
            midiNote, 
            event.Velocity,
        )
    case NOTE_OFF:
        midiNote := event.Data.(uint8)
        return pe.midiPlayer.SendNoteOff(
            uint8(event.Channel), 
            midiNote, 
            0,
        )
    default:
        return nil
    }
}

func (pe *PlayEngine) beatsToTime(beats float64) time.Duration {
    secondsPerBeat := 60.0 / pe.bpm
    seconds := beats * secondsPerBeat
    return time.Duration(seconds * float64(time.Second))
}