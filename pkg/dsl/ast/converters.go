package ast

import "catRock/pkg/core"

// 字符串转音符名称
func stringToNoteName(name string) core.BaseNoteName {
    switch name {
    case "C":
        return core.C
    case "D":
        return core.D
    case "E":
        return core.E
    case "F":
        return core.F
    case "G":
        return core.G
    case "A":
        return core.A
    case "B":
        return core.B
    default:
        return core.C // 默认值
    }
}

// 字符串转节拍值
func stringToBeatValue(duration string) core.BeatValue {
    switch duration {
    case "whole":
        return core.Whole
    case "half":
        return core.Half
    case "quarter":
        return core.Quarter
    case "eighth":
        return core.Eighth
    default:
        return core.Quarter // 默认值
    }
}

// 字符串转乐器
func stringToInstrument(inst string) core.InstrumentID {
    switch inst {
    case "piano":
        return core.AcousticGuitarSteel
    case "guitar":
        return core.AcousticGuitarNylon
    case "violin":
        return core.Violin
    case "drums":
        return core.ElectronicDrumKit
    default:
        return core.AcousticGuitarSteel // 默认值
    }
}

// 从参数中提取乐器
func getInstrument(params map[string]interface{}) core.InstrumentID {
    if inst, ok := params["instrument"]; ok {
        if instStr, ok := inst.(string); ok {
            return stringToInstrument(instStr)
        }
    }
    return core.AcousticGuitarSteel // 默认值
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