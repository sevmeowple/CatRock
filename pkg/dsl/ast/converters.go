package ast

import "catRock/pkg/core"

// 字符串转音符名称
func stringToNoteName(name string) core.BaseNoteName {
	switch name {
	// 基本音符
	case "C", "c":
		return core.C
	case "D", "d":
		return core.D
	case "E", "e":
		return core.E
	case "F", "f":
		return core.F
	case "G", "g":
		return core.G
	case "A", "a":
		return core.A
	case "B", "b":
		return core.B

	// 升号音符 (s后缀)
	case "Cs", "cs":
		return core.Cs
	case "Ds", "ds":
		return core.Ds
	case "Fs", "fs":
		return core.Fs
	case "Gs", "gs":
		return core.Gs
	case "As", "as":
		return core.As

	// 降号音符 (b后缀) - 转换为对应的升号
	case "Db", "db":
		return core.Cs // Db = C#
	case "Eb", "eb":
		return core.Ds // Eb = D#
	case "Gb", "gb":
		return core.Fs // Gb = F#
	case "Ab", "ab":
		return core.Gs // Ab = G#
	case "Bb", "bb":
		return core.As // Bb = A#

	default:
		return core.C // 默认返回C
	}
}

// 字符串转节拍值
func stringToBeatValue(duration string) core.BeatValue {
	switch duration {
	case "1/1":
		return core.Whole
	case "1/2":
		return core.Half
	case "1/4":
		return core.Quarter
	case "1/8":
		return core.Eighth
	case "1/16":
		return core.Sixteenth
	default:
		return core.Quarter
	}
}

// 从参数中提取乐器
func getInstrument(params map[string]interface{}) core.InstrumentID {
	if inst, ok := params["instrument"]; ok {
		if instInt, ok := inst.(int); ok {
			return core.InstrumentID(instInt)
		}
	}
	return core.InstrumentID(0) // 默认值
}

// 从参数中提取通道
func getChannel(params map[string]interface{}) int {
	if ch, ok := params["channel"]; ok {
		if chInt, ok := ch.(int); ok {
			return chInt
		}
	}
	return 1 // 默认值
}

// 从参数中提取音量
func getVolume(params map[string]interface{}) int {
	if vol, ok := params["volume"]; ok {
		if volInt, ok := vol.(int); ok {
			return volInt
		}
	}
	return 100 // 默认值
}

// 从和弦名创建和弦
func createChordFromName(chordName string, duration string) core.Chord {
	// 这里需要解析和弦名称，如 "Am", "Cmaj", "F#dim" 等
	// 暂时简化实现
	beatValue := stringToBeatValue(duration)

	// 创建一个简单的C大三和弦作为示例
	notes := []core.Note{
		core.NewNote(core.NewNoteParams{Name: core.C, Octave: 4, Beat: beatValue}),
		core.NewNote(core.NewNoteParams{Name: core.E, Octave: 4, Beat: beatValue}),
		core.NewNote(core.NewNoteParams{Name: core.G, Octave: 4, Beat: beatValue}),
	}

	return core.NewChord(notes)
}
