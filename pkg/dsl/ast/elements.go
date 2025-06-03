package ast

import (
	"catRock/pkg/core"
	"catRock/pkg/dsl/mytype"
	"catRock/pkg/score"
	"fmt"
)

// 音符节点
type NoteNode struct {
    Name     string // C, D, E, F, G, A, B
    Octave   int    // 0-9
    Duration string // quarter, half, whole, eighth
    Position mytype.Position
}

var _ ElementNode = (*NoteNode)(nil)

func (n *NoteNode) String() string {
    return fmt.Sprintf("Note{%s%d %s}", n.Name, n.Octave, n.Duration)
}

// 转换为Score系统的Playable
func (n *NoteNode) ToPlayable() score.Playable {
    // 转换音符名称
    noteName := stringToNoteName(n.Name)

    // 转换时值
    beatValue := stringToBeatValue(n.Duration)

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

// 和弦节点
type ChordNode struct {
    Content  interface{} // string(和弦名) 或 []*NoteNode(手动构建)
    Duration string      // quarter, half, whole, eighth
    Position mytype.Position
}

var _ ElementNode = (*ChordNode)(nil)

func (c *ChordNode) String() string {
    return fmt.Sprintf("Chord{%v %s}", c.Content, c.Duration)
}

func (c *ChordNode) ToPlayable() score.Playable {
    // 根据Content类型创建和弦
    var chord core.Chord
    
    switch content := c.Content.(type) {
    case string:
        // 从和弦名创建，如 "Am", "Cmaj"
        chord = createChordFromName(content, c.Duration)
    case []*NoteNode:
        // 从音符列表手动构建
        notes := make([]core.Note, len(content))
        for i, noteNode := range content {
            noteElement := noteNode.ToPlayable().(*score.NoteElement)
            notes[i] = noteElement.Note
        }
        chord = core.NewChord(notes)
    default:
        // 默认创建C大三和弦
        chord = createChordFromName("C", c.Duration)
    }

    return score.NewChordElement(chord)
}

// 休止符节点
type RestNode struct {
    Duration string
    Position mytype.Position
}

var _ ElementNode = (*RestNode)(nil)

func (n *NoteNode) DetailedString(indent string) string {
    return fmt.Sprintf("NoteNode { 音符:%s%d, 时值:%s, 位置:%s }\n", 
        n.Name, n.Octave, n.Duration, n.Position)
}

func (c *ChordNode) DetailedString(indent string) string {
    result := fmt.Sprintf("ChordNode {\n")
    result += fmt.Sprintf("%s  内容: %v (%T)\n", indent, c.Content, c.Content)
    result += fmt.Sprintf("%s  时值: %s\n", indent, c.Duration)
    result += fmt.Sprintf("%s  位置: %s\n", indent, c.Position)
    result += fmt.Sprintf("%s}\n", indent)
    return result
}

func (r *RestNode) DetailedString(indent string) string {
    return fmt.Sprintf("RestNode { 时值:%s, 位置:%s }\n", r.Duration, r.Position)
}

func (g *GroupNode) DetailedString(indent string) string {
    result := fmt.Sprintf("GroupNode {\n")
    result += fmt.Sprintf("%s  时值: %s\n", indent, g.Duration)
    result += fmt.Sprintf("%s  位置: %s\n", indent, g.Position)
    
    if len(g.Elements) > 0 {
        result += fmt.Sprintf("%s  元素 (%d个):\n", indent, len(g.Elements))
        for i, element := range g.Elements {
            result += fmt.Sprintf("%s    [%d] %s", indent, i, element.DetailedString(indent+"      "))
        }
    }
    
    result += fmt.Sprintf("%s}\n", indent)
    return result
}

func (r *RestNode) String() string {
    return fmt.Sprintf("Rest{%s}", r.Duration)
}

func (r *RestNode) ToPlayable() score.Playable {
    beatValue := stringToBeatValue(r.Duration)
    rest := core.NewRest(beatValue)
    return score.NewRestElement(rest)
}

type GroupNode struct {
    Elements []ElementNode // 包含的元素
    Duration string // 可选的持续时间
    Position mytype.Position
}

var _ ElementNode = (*GroupNode)(nil)

func (g *GroupNode) String() string {
    return fmt.Sprintf("Group{%v %s}", g.Elements, g.Duration)
}

// 修改 GroupNode.ToPlayable() 方法
func (g *GroupNode) ToPlayable() score.Playable {
    // 创建 GroupElement
    var group *score.GroupElement
    
    if g.Duration != "" && g.Duration != "1/4" { // 如果有指定时值
        duration := stringToBeatValue(g.Duration)
        group = score.NewGroupElementWithDuration(float64(duration))
    } else {
        group = score.NewGroupElement()
    }
    
    // 添加所有子元素
    for _, elem := range g.Elements {
        if playable := elem.ToPlayable(); playable != nil {
            group.AddElement(playable)
        }
    }
    
    return group
}