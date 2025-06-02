package score

import (
    "catRock/pkg/core"
    "fmt"
)

// 段落 - 顺序播放容器
type Section struct {
    ID       string
    Name     string
    Elements []Playable
    
    // 容器参数
    ContainerParams
}

var _ Container = (*Section)(nil)

func (s *Section) GetID() string {
    if s.ID != "" {
        return s.ID
    }
    return fmt.Sprintf("section_%s", s.Name)
}

func (s *Section) GetType() PlayableType {
    return SECTION_TYPE
}

func (s *Section) Duration(context PlayContext) float64 {
    // 顺序播放：所有元素时间累加
    totalDuration := 0.0
    sectionContext := context.WithContainerSettings(s.ContainerParams)
    
    for _, element := range s.Elements {
        totalDuration += element.Duration(sectionContext)
    }
    
    return totalDuration
}

func (s *Section) GenerateEvents(startTime float64, context PlayContext) []Event {
    events := []Event{}
    currentTime := startTime
    sectionContext := context.WithContainerSettings(s.ContainerParams)
    
    // 循环检测
    for _, stackID := range context.ElementStack {
        if stackID == s.GetID() {
            return events
        }
    }
    
    sectionContext.ElementStack = append(context.ElementStack, s.GetID())
    
    // 顺序播放：每个元素依次开始
    for _, element := range s.Elements {
        elementEvents := element.GenerateEvents(currentTime, sectionContext)
        events = append(events, elementEvents...)
        currentTime += element.Duration(sectionContext)
    }
    
    return events
}

// 实现Container接口
func (s *Section) AddElement(element Playable) {
    s.Elements = append(s.Elements, element)
}

func (s *Section) SetBPM(bpm float64) {
    s.BPM = &bpm
}

func (s *Section) SetVolume(volume int) {
    s.Volume = &volume
}

func (s *Section) SetInstrument(instrument core.InstrumentID) {
    s.Instrument = &instrument
}

func (s *Section) SetChannel(channel int) {
    s.Channel = &channel
}

// 构造函数
func NewSection(name string) *Section {
    return &Section{
        Name:     name,
        Elements: []Playable{},
    }
}