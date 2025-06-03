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

	// 新增：生成乐器切换事件
	if s.Instrument != nil && *s.Instrument != context.CurrentInstrument {

		targetChannel := sectionContext.CurrentChannel

		if core.IsDrumKit(*s.Instrument) {
			targetChannel = 9
		}

		programChangeEvent := Event{
			Time:          startTime - 0.001,
			Duration:      0,
			Type:          CONTROL_EVENT,
			Action:        PROGRAM_CHANGE,
			Data:          *s.Instrument,
			Channel:       targetChannel,
			Velocity:      0,
			SourceElement: s.GetID(),
		}
		events = append(events, programChangeEvent)

		sectionContext.CurrentInstrument = *s.Instrument
		sectionContext.CurrentChannel = targetChannel // 也要更新通道
	}

	// 新增：生成音量切换事件（如果需要）
	if s.Volume != nil && *s.Volume != context.CurrentVolume {
		volumeChangeEvent := Event{
			Time:          startTime,
			Duration:      0,
			Type:          CONTROL_EVENT,
			Action:        VOLUME_CHANGE,
			Data:          uint8(*s.Volume),
			Channel:       sectionContext.CurrentChannel,
			Velocity:      0,
			SourceElement: s.GetID(),
		}
		events = append(events, volumeChangeEvent)
	}

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

func (s *Section) SetInstrument(instrument int) {
	id := core.InstrumentID(instrument)
	s.Instrument = &id
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

func (s *Section) DetailedString(indent string) string {
	result := fmt.Sprintf("Section '%s' {\n", s.Name)
	result += fmt.Sprintf("%s  ID: %s\n", indent, s.GetID())
	result += fmt.Sprintf("%s  类型: %s\n", indent, s.getTypeName())
	result += fmt.Sprintf("%s  时长: %.3f拍\n", indent, s.Duration(PlayContext{}))

	// 显示容器参数
	if s.BPM != nil || s.Volume != nil || s.Instrument != nil || s.Channel != nil {
		result += fmt.Sprintf("%s  容器参数:\n", indent)
		if s.BPM != nil {
			result += fmt.Sprintf("%s    BPM: %.1f\n", indent, *s.BPM)
		}
		if s.Volume != nil {
			result += fmt.Sprintf("%s    音量: %d\n", indent, *s.Volume)
		}
		if s.Instrument != nil {
			result += fmt.Sprintf("%s    乐器: %d\n", indent, int(*s.Instrument))
		}
		if s.Channel != nil {
			result += fmt.Sprintf("%s    通道: %d\n", indent, *s.Channel)
		}
	}

	if len(s.Elements) > 0 {
		result += fmt.Sprintf("%s  元素 (%d个):\n", indent, len(s.Elements))
		for i, element := range s.Elements {
			result += fmt.Sprintf("%s    [%d] %s", indent, i, element.DetailedString(indent+"      "))
		}
	}

	result += fmt.Sprintf("%s}\n", indent)
	return result
}

func (s *Section) getTypeName() string {
	switch s.GetType() {
	case SECTION_TYPE:
		return "段落(顺序播放)"
	default:
		return "未知"
	}
}
