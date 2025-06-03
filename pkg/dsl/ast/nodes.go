package ast

import (
	"catRock/pkg/dsl/mytype"
	"fmt"
)

// 顶层Score节点
type ScoreNode struct {
    GlobalSets []*SetNode     // 全局设置
    Elements   []PlayableNode // 顶层可播放元素
    Position   mytype.Position
}

func (s *ScoreNode) String() string {
    return fmt.Sprintf("Score{GlobalSets: %d, Elements: %d, Position: %s}",
        len(s.GlobalSets), len(s.Elements), s.Position)
}

func (s *ScoreNode) DetailedString(indent string) string {
    result := fmt.Sprintf("%sScoreNode {\n", indent)
    result += fmt.Sprintf("%s  位置: %s\n", indent, s.Position)
    
    if len(s.GlobalSets) > 0 {
        result += fmt.Sprintf("%s  全局设置 (%d个):\n", indent, len(s.GlobalSets))
        for i, setNode := range s.GlobalSets {
            result += fmt.Sprintf("%s    [%d] %s", indent, i, setNode.DetailedString(indent+"      "))
        }
    }
    
    if len(s.Elements) > 0 {
        result += fmt.Sprintf("%s  顶层元素 (%d个):\n", indent, len(s.Elements))
        for i, element := range s.Elements {
            result += fmt.Sprintf("%s    [%d] %s", indent, i, element.DetailedString(indent+"      "))
        }
    }
    
    result += fmt.Sprintf("%s}\n", indent)
    return result
}

// Header节点
type HeaderNode struct {
    BPM      int
    Position mytype.Position
}

func (h *HeaderNode) String() string {
    return fmt.Sprintf("Header{BPM: %d}", h.BPM)
}