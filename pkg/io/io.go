package io

import (
	"catRock/pkg/core"
)

type ConnectStatus int // 连接状态码
const (
	Connected    ConnectStatus = iota // 已连接
	Disconnected                      // 已断开
	Connecting                        // 正在连接中
)

type PlayNoteParams struct {
	Note core.Note // 要播放的音符
	Bpm  int       // 节拍
}

type PlayChordParams struct {
	Chord core.Chord // 要播放的和弦
	Bpm   int        // 节拍
}

type IO interface {
	// Read reads data from the input source
	Connect() (ConnectStatus, error)
	Disconnect() (ConnectStatus, error)
	PlayNote(params PlayNoteParams) error
	PlayChord(params PlayChordParams) error
	SetBPM(bpm int) error
	SetVolume(volume int) error
	SetProgram(program core.InstrumentID) error // 设置乐器程序
	SetChannel(channel uint8) error // 设置MIDI通道

	MIDIEventSender // 新增：低级MIDI事件接口
}

// 新增：低级MIDI事件接口
type MIDIEventSender interface {
    // 原始MIDI消息发送
    SendNoteOn(channel uint8, note uint8, velocity uint8) error
    SendNoteOff(channel uint8, note uint8, velocity uint8) error
    SendProgramChange(channel uint8, program uint8) error
    SendControlChange(channel uint8, controller uint8, value uint8) error
}

type PlayEventsParams struct {
    Events []Event  // 事件序列
    BPM    float64  // 播放BPM
}

// 新增：MIDI事件结构（简化版）
type Event struct {
    Time     float64    // 绝对时间（拍数）
    Type     EventType  // 事件类型
    Channel  uint8      // MIDI通道
    Data1    uint8      // 数据1（如音符编号）
    Data2    uint8      // 数据2（如力度）
}

type EventType int
const (
    NOTE_ON_EVENT EventType = iota
    NOTE_OFF_EVENT
    PROGRAM_CHANGE_EVENT
    CONTROL_CHANGE_EVENT
)