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

// Header节点
type HeaderNode struct {
    BPM      int
    Position mytype.Position
}

func (h *HeaderNode) String() string {
    return fmt.Sprintf("Header{BPM: %d}", h.BPM)
}