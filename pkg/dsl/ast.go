package dsl

import (
    "catRock/pkg/core"
    "catRock/pkg/score"
    "fmt"
)

// AST节点基础接口
type ASTNode interface {
    String() string
}

// 顶层Score节点
type ScoreNode struct {
    Header *HeaderNode
    Body   []ElementNode
}

func (s *ScoreNode) String() string {
    return fmt.Sprintf("Score{Header: %s, Elements: %d}", s.Header, len(s.Body))
}

// Header节点
type HeaderNode struct {
    BPM int
}

func (h *HeaderNode) String() string {
    return fmt.Sprintf("Header{BPM: %d}", h.BPM)
}

// 元素节点接口
type ElementNode interface {
    ASTNode
    ToPlayable() score.Playable
}

// 音符节点
type NoteNode struct {
    Name     string // C, D, E, F, G, A, B
    Octave   int    // 0-9
    Duration string // quarter, half, whole, eighth
}

func (n *NoteNode) String() string {
    return fmt.Sprintf("Note{%s%d %s}", n.Name, n.Octave, n.Duration)
}

// 转换为Score系统的Playable
func (n *NoteNode) ToPlayable() score.Playable {
    // 转换音符名称
    noteName := n.stringToNoteName(n.Name)
    
    // 转换时值
    beatValue := n.stringToBeatValue(n.Duration)
    
    // 创建core.Note
    note := core.NewNote(core.NewNoteParams{
        Name:       noteName,
        Octave:     n.Octave,
        Accidental: core.Natural,
        Beat:       beatValue,
    })
    
    // 创建NoteElement
    return score.NewNoteElement(note)
}

func (n *NoteNode) stringToNoteName(name string) core.BaseNoteName {
    switch name {
    case "C": return core.C
    case "D": return core.D
    case "E": return core.E
    case "F": return core.F
    case "G": return core.G
    case "A": return core.A
    case "B": return core.B
    default:  return core.C // 默认值
    }
}

func (n *NoteNode) stringToBeatValue(duration string) core.BeatValue {
    switch duration {
    case "whole":   return core.Whole
    case "half":    return core.Half
    case "quarter": return core.Quarter
    case "eighth":  return core.Eighth
    default:        return core.Quarter // 默认值
    }
}