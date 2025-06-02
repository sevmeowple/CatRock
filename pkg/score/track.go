package score

import (
    "catRock/pkg/core"
    "fmt"
    "sort"
)

// 轨道 - 并行播放容器
type Track struct {
    ID       string
    Name     string
    Elements []Playable
    
    // 容器参数
    ContainerParams
}

var _ Container = (*Track)(nil)

func (t *Track) GetID() string {
    if t.ID != "" {
        return t.ID
    }
    return fmt.Sprintf("track_%s", t.Name)
}

func (t *Track) GetType() PlayableType {
    return TRACK_TYPE
}

func (t *Track) Duration(context PlayContext) float64 {
    // 并行播放：最长元素的时间
    maxDuration := 0.0
    trackContext := context.WithContainerSettings(t.ContainerParams)
    
    for _, element := range t.Elements {
        duration := element.Duration(trackContext)
        if duration > maxDuration {
            maxDuration = duration
        }
    }
    
    return maxDuration
}

func (t *Track) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    trackContext := context.WithContainerSettings(t.ContainerParams)
    
    // 循环检测
    for _, stackID := range context.ElementStack {
        if stackID == t.GetID() {
            return events
        }
    }
    
    trackContext.ElementStack = append(context.ElementStack, t.GetID())
    
    // 并行播放：所有元素同时开始
    for _, element := range t.Elements {
        elementEvents := element.GenerateEvents(startTime, trackContext)
        events = append(events, elementEvents...)
    }
    
    return t.sortEventsByTime(events)
}

// 实现Container接口
func (t *Track) AddElement(element Playable) {
    t.Elements = append(t.Elements, element)
}

func (t *Track) SetBPM(bpm float64) {
    t.BPM = &bpm
}

func (t *Track) SetVolume(volume int) {
    t.Volume = &volume
}

func (t *Track) SetInstrument(instrument core.InstrumentID) {
    t.Instrument = &instrument
}

func (t *Track) SetChannel(channel int) {
    t.Channel = &channel
}

// 辅助方法
func (t *Track) sortEventsByTime(events []Event) []Event {
    sort.Slice(events, func(i, j int) bool {
        if events[i].Time == events[j].Time {
            return events[i].Action == NOTE_ON && events[j].Action == NOTE_OFF
        }
        return events[i].Time < events[j].Time
    })
    return events
}

// 构造函数
func NewTrack(name string) *Track {
    return &Track{
        Name:     name,
        Elements: []Playable{},
    }
}