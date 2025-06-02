package ast

import (
	"catRock/pkg/score"
)

// AST节点基础接口
type ASTNode interface {
    String() string
}

// 可播放元素接口 - 能转换为score.Playable
type PlayableNode interface {
    ASTNode
    ToPlayable() score.Playable
}

// 元素节点接口
type ElementNode interface {
    ASTNode
    ToPlayable() score.Playable
}

// 容器节点接口
type ContainerNode interface {
    PlayableNode
    AddElement(element PlayableNode)
    MergeSetParameters() map[string]interface{}
}
