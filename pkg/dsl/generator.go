package dsl

import (
    "catRock/pkg/score"
    "fmt"
)

type Generator struct {
    errors []string
}

func NewGenerator() *Generator {
    return &Generator{
        errors: []string{},
    }
}

func (g *Generator) Errors() []string {
    return g.errors
}

func (g *Generator) GenerateScore(ast *ScoreNode) (*score.Score, error) {
    if ast == nil {
        return nil, fmt.Errorf("AST为空")
    }
    
    if ast.Header == nil {
        return nil, fmt.Errorf("缺少Header")
    }
    
    // 创建Score对象
    scoreObj := &score.Score{
        Title:    "Simple DSL Song",
        Composer: "CatRock DSL",
        Year:     2024,
        BPM:      float64(ast.Header.BPM),
        Volume:   100,
    }
    
    // 转换body元素
    elements := []score.Playable{}
    for _, elementNode := range ast.Body {
        if elementNode != nil {
            playable := elementNode.ToPlayable()
            elements = append(elements, playable)
        }
    }
    
    if len(elements) == 0 {
        return nil, fmt.Errorf("没有音乐元素")
    }
    
    // 创建一个Section来包含所有元素（顺序播放）
    rootSection := score.NewSection("main", elements)
    scoreObj.RootElement = rootSection
    
    return scoreObj, nil
}