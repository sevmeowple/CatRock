package score

import (
	"catRock/pkg/core"
	"fmt"
)

// 休止符播放元素
type RestElement struct {
	ID   string
	Rest core.Rest
}

var _ Element = (*RestElement)(nil)

func (re *RestElement) GetID() string {
	if re.ID != "" {
		return re.ID
	}
	return fmt.Sprintf("rest_%.2f", re.Rest.Duration())
}

func (re *RestElement) GetType() PlayableType {
	return REST_TYPE
}

func (re *RestElement) Duration(context PlayContext) float64 {
	return re.Rest.Duration()
}

func (re *RestElement) GenerateEvents(startTime float64, context PlayContext) []Event {
	// 休止符不产生任何MIDI事件，只占用时间
	return []Event{}
}

// 实现Element接口（休止符不需要这些覆盖）
func (re *RestElement) SetVolumeOverride(volume int) {
	// 休止符无音量
}

func (re *RestElement) SetInstrumentOverride(instrument core.InstrumentID) {
	// 休止符无乐器
}

func (re *RestElement) SetChannelOverride(channel int) {
	// 休止符无通道
}

// 构造函数
func NewRestElement(rest core.Rest) *RestElement {
	return &RestElement{Rest: rest}
}

func (re *RestElement) String() string {
	return fmt.Sprintf("RestElement(ID: %s, Duration: %.2f)", re.GetID(), re.Rest.Duration())
}

func (re *RestElement) DetailedString(indent string) string {

	return fmt.Sprintf("%sRestElement(ID: %s, Duration: %.2f)",
		indent, re.GetID(), re.Rest.Duration())
}
