package dsl

import (
	"catRock/pkg/dsl/ast"
	"catRock/pkg/score"
	"fmt"
	"strconv"
	"strings"
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

func (g *Generator) addError(msg string) {
    g.errors = append(g.errors, msg)
}

// 更新：使用新的AST结构
func (g *Generator) GenerateScore(scoreNode *ast.ScoreNode) (*score.Score, error) {
    if scoreNode == nil {
        return nil, fmt.Errorf("AST为空")
    }
    
    // 创建Score对象
    scoreObj := &score.Score{
        Title:    "CatRock DSL Song",
        Composer: "CatRock DSL",
        Year:     2024,
        BPM:      120.0,  // 默认BPM，会被全局设置覆盖
        Volume:   100,
    }
    
    // 应用全局设置
    if err := g.applyGlobalSettings(scoreObj, scoreNode.GlobalSets); err != nil {
        g.addError(fmt.Sprintf("应用全局设置失败: %v", err))
    }
    
    // 转换所有顶层元素
    elements := []score.Playable{}
    for _, elementNode := range scoreNode.Elements {
        if elementNode != nil {
            playable := elementNode.ToPlayable()
            if playable != nil {
                elements = append(elements, playable)
            }
        }
    }
    
    if len(elements) == 0 {
        return nil, fmt.Errorf("没有音乐元素")
    }
    
    // 创建根容器
    var rootElement score.Playable
    if len(elements) == 1 {
        // 如果只有一个元素，直接使用它
        rootElement = elements[0]
    } else {
        // 多个元素，创建一个Section来顺序播放
        rootSection := score.NewSection("main")
        for _, element := range elements {
            rootSection.AddElement(element)
        }
        rootElement = rootSection
    }
    
    scoreObj.RootElement = rootElement
    
    return scoreObj, nil
}

// 应用全局设置
func (g *Generator) applyGlobalSettings(scoreObj *score.Score, globalSets []*ast.SetNode) error {
    // 合并所有全局设置
    globalParams := make(map[string]interface{})
    
    for _, setNode := range globalSets {
        setNode.Context = ast.GlobalContext
        resolved, err := setNode.ResolveParameters()
        if err != nil {
            return fmt.Errorf("解析全局参数失败: %v", err)
        }
        
        // 后面的设置覆盖前面的
        for key, value := range resolved {
            globalParams[key] = value
        }
    }
    
    // 应用设置到Score对象
    if bpmValue, ok := globalParams["BPM"]; ok {
        bpm := g.parseNumericValue(bpmValue)
        if bpm > 0 {
            scoreObj.SetBPM(bpm)
        }
    }
    
    if volumeValue, ok := globalParams["volume"]; ok {
        volume := int(g.parseNumericValue(volumeValue))
        if volume >= 0 && volume <= 127 {
            scoreObj.SetVolume(volume)
        }
    }
    
    // 可以添加更多全局设置的处理...
    
    return nil
}

// 新增：统一的数值解析方法
func (g *Generator) parseNumericValue(value interface{}) float64 {
    switch v := value.(type) {
    case int:
        return float64(v)
    case float64:
        return v
    case string:
        // 尝试解析普通数字
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            return f
        }
        
        // 尝试解析分数
        if parts := strings.Split(v, "/"); len(parts) == 2 {
            if num, err1 := strconv.ParseFloat(parts[0], 64); err1 == nil {
                if den, err2 := strconv.ParseFloat(parts[1], 64); err2 == nil && den != 0 {
                    return num / den
                }
            }
        }
        
        g.addError(fmt.Sprintf("无法解析数值: %v", v))
        return 0
    default:
        g.addError(fmt.Sprintf("不支持的数值类型: %T", v))
        return 0
    }
}

// 便捷方法：从文本生成Score
func (g *Generator) GenerateFromText(input string) (*score.Score, error) {
    // 创建词法分析器
    lexer := NewLexer(input)
    
    // 创建语法分析器
    parser := NewParser(lexer)
    
    // 解析AST
    scoreAST := parser.ParseScore()
    if len(parser.Errors()) > 0 {
        return nil, fmt.Errorf("解析错误: %v", parser.Errors())
    }
    
    // 生成Score
    return g.GenerateScore(scoreAST)
}

// 验证AST的辅助方法
func (g *Generator) validateAST(scoreNode *ast.ScoreNode) error {
    if scoreNode == nil {
        return fmt.Errorf("ScoreNode为空")
    }
    
    // 检查是否有可播放元素
    if len(scoreNode.Elements) == 0 {
        return fmt.Errorf("没有可播放元素")
    }
    
    // 可以添加更多验证逻辑...
    
    return nil
}