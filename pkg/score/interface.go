package score

import "catRock/pkg/core"

// 基础播放接口
type Playable interface {
    GenerateEvents(startTime float64, context PlayContext) []Event
    Duration(context PlayContext) float64
    GetID() string
    GetType() PlayableType
}

// 播放元素类型
type PlayableType int

const (
    NOTE_TYPE PlayableType = iota
    CHORD_TYPE
    REST_TYPE

    GROUP_TYPE

    SECTION_TYPE
    TRACK_TYPE
)

// 容器接口 - Section和Track的共同接口
type Container interface {
    Playable
    AddElement(element Playable)
    SetBPM(bpm float64)
    SetVolume(volume int)
    SetInstrument(instrument core.InstrumentID)
    SetChannel(channel int)
}

// 元素接口 - Note、Chord、Rest的共同接口
type Element interface {
    Playable
    SetVolumeOverride(volume int)
    SetInstrumentOverride(instrument core.InstrumentID)
    SetChannelOverride(channel int)
}