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

	// 新增：生成乐器切换事件
	if t.Instrument != nil && *t.Instrument != context.CurrentInstrument {
		
		targetChannel := trackContext.CurrentChannel
		if core.IsDrumKit(*t.Instrument) {
			targetChannel = 9 // 鼓组通常使用通道9
		}
		
		programChangeEvent := Event{
			Time:          startTime,
			Duration:      0,
			Type:          CONTROL_EVENT,
			Action:        PROGRAM_CHANGE,
			Data:          *t.Instrument,
			Channel:       targetChannel,
			Velocity:      0,
			SourceElement: t.GetID(),
		}
		events = append(events, programChangeEvent)

		trackContext.CurrentInstrument = *t.Instrument
		trackContext.CurrentChannel = targetChannel // 更新通道
	}

	// 新增：生成音量切换事件（如果需要）
	if t.Volume != nil && *t.Volume != context.CurrentVolume {
		volumeChangeEvent := Event{
			Time:          startTime,
			Duration:      0,
			Type:          CONTROL_EVENT,
			Action:        VOLUME_CHANGE,
			Data:          uint8(*t.Volume),
			Channel:       trackContext.CurrentChannel,
			Velocity:      0,
			SourceElement: t.GetID(),
		}
		events = append(events, volumeChangeEvent)
	}

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

func (t *Track) SetInstrument(instrument int) {
	id := core.InstrumentID(instrument)
	t.Instrument = &id
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

func (t *Track) DetailedString(indent string) string {
	result := fmt.Sprintf("Track '%s' {\n", t.Name)
	result += fmt.Sprintf("%s  ID: %s\n", indent, t.GetID())
	result += fmt.Sprintf("%s  类型: %s\n", indent, t.getTypeName())
	result += fmt.Sprintf("%s  时长: %.3f拍\n", indent, t.Duration(PlayContext{}))

	// 显示容器参数
	if t.BPM != nil || t.Volume != nil || t.Instrument != nil || t.Channel != nil {
		result += fmt.Sprintf("%s  容器参数:\n", indent)
		if t.BPM != nil {
			result += fmt.Sprintf("%s    BPM: %.1f\n", indent, *t.BPM)
		}
		if t.Volume != nil {
			result += fmt.Sprintf("%s    音量: %d\n", indent, *t.Volume)
		}
		if t.Instrument != nil {
			result += fmt.Sprintf("%s    乐器: %d\n", indent, int(*t.Instrument))
		}
		if t.Channel != nil {
			result += fmt.Sprintf("%s    通道: %d\n", indent, *t.Channel)
		}
	}

	if len(t.Elements) > 0 {
		result += fmt.Sprintf("%s  元素 (%d个):\n", indent, len(t.Elements))
		for i, element := range t.Elements {
			result += fmt.Sprintf("%s    [%d] %s", indent, i, element.DetailedString(indent+"      "))
		}
	}

	result += fmt.Sprintf("%s}\n", indent)
	return result
}

func (t *Track) getTypeName() string {
	switch t.GetType() {
	case TRACK_TYPE:
		return "轨道(并行播放)"
	default:
		return "未知"
	}
}
