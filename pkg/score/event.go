package score

import "fmt"

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

func (e *Event) String() string {
    var actionName string
    switch e.Action {
    case NOTE_ON:
        actionName = "NOTE_ON"
    case NOTE_OFF:
        actionName = "NOTE_OFF"
    case PROGRAM_CHANGE:
        actionName = "PROGRAM_CHANGE"
    case VOLUME_CHANGE:
        actionName = "VOLUME_CHANGE"
    default:
        actionName = fmt.Sprintf("ACTION_%d", e.Action)
    }
    
    return fmt.Sprintf("[%.3fs] %s Ch:%d Data:%v Vel:%d Dur:%.3f Src:%s",
        e.Time, actionName, e.Channel, e.Data, e.Velocity, e.Duration, e.SourceElement)
}

// 添加详细字符串方法
func (e *Event) DetailedString() string {
    result := fmt.Sprintf("Event at %.3fs:\n", e.Time)
    result += fmt.Sprintf("  Action: %s\n", e.getActionName())
    result += fmt.Sprintf("  Channel: %d\n", e.Channel)
    result += fmt.Sprintf("  Data: %v\n", e.Data)
    
    if e.Velocity > 0 {
        result += fmt.Sprintf("  Velocity: %d\n", e.Velocity)
    }
    if e.Duration > 0 {
        result += fmt.Sprintf("  Duration: %.3f\n", e.Duration)
    }
    if e.SourceElement != "" {
        result += fmt.Sprintf("  Source: %s\n", e.SourceElement)
    }
    
    return result
}

func (e *Event) getActionName() string {
    switch e.Action {
    case NOTE_ON:
        return "NOTE_ON"
    case NOTE_OFF:
        return "NOTE_OFF"
    case PROGRAM_CHANGE:
        return "PROGRAM_CHANGE"
    case VOLUME_CHANGE:
        return "VOLUME_CHANGE"
    default:
        return fmt.Sprintf("UNKNOWN_%d", e.Action)
    }
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
