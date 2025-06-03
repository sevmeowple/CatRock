package score

import "catRock/pkg/core"

// 播放上下文
type PlayContext struct {
	// 全局设置
	GlobalBPM    float64
	GlobalVolume int

	// 当前生效设置
	CurrentBPM        float64
	CurrentVolume     int
	CurrentInstrument core.InstrumentID
	CurrentChannel    int

	// 循环检测
	ElementStack []string
}

// 创建新的播放上下文
func NewPlayContext(bpm float64, volume int) PlayContext {
	return PlayContext{
		GlobalBPM:         bpm,
		GlobalVolume:      volume,
		CurrentBPM:        bpm,
		CurrentVolume:     volume,
		CurrentInstrument: 1, // 默认乐器ID
		CurrentChannel:    1,
		ElementStack:      []string{},
	}
}

// 应用容器设置
func (pc PlayContext) WithContainerSettings(params ContainerParams) PlayContext {
	context := pc

	if params.BPM != nil {
		context.CurrentBPM = *params.BPM
	}
	if params.Volume != nil {
		context.CurrentVolume = *params.Volume
	}
	if params.Instrument != nil {
		context.CurrentInstrument = *params.Instrument
	}
	if params.Channel != nil {
		context.CurrentChannel = *params.Channel
	}

	return context
}

// 通用容器参数
type ContainerParams struct {
	BPM        *float64
	Volume     *int
	Instrument *core.InstrumentID
	Channel    *int
}
