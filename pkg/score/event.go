package score

// 播放事件
type Event struct {
    Time     float64
    Duration float64
    Type     EventType
    Action   EventAction
    Data     interface{}
    Channel  int
    Velocity uint8
    SourceElement string
}

// 事件类型
type EventType int

const (
    NOTE_EVENT EventType = iota
    CHORD_EVENT
    CONTROL_EVENT
)

// 事件动作
type EventAction int

const (
    NOTE_ON EventAction = iota
    NOTE_OFF
    VOLUME_CHANGE
    PROGRAM_CHANGE
)
