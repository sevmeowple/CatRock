package ast

import (
	"catRock/pkg/dsl/mytype"
	"fmt"
)

// 参数规范
type ParameterSpec struct {
    Name         string
    Type         ParameterType
    DefaultValue interface{}
    Required     bool
    Description  string
}

type ParameterType int

const (
    ParamInt ParameterType = iota
    ParamFloat
    ParamString
    ParamBool
)

// 参数上下文
type ParameterContext int

const (
    GlobalContext ParameterContext = iota
    TrackContext
    SectionContext
)

// 全局参数规范
var GlobalParameters = map[string]ParameterSpec{
    "BPM": {
        Name:         "BPM",
        Type:         ParamInt,
        DefaultValue: 120,
        Required:     false,
        Description:  "每分钟节拍数",
    },
    "base_duration": {
        Name:         "base_duration",
        Type:         ParamString,
        DefaultValue: "1/4",
        Required:     false,
        Description:  "默认音符时值",
    },
}

// Track参数规范
var TrackParameters = map[string]ParameterSpec{
    "channel": {
        Name:         "channel",
        Type:         ParamInt,
        DefaultValue: 1,
        Required:     false,
        Description:  "MIDI通道",
    },
    "instrument": {
        Name:         "instrument",
        Type:         ParamInt,
        DefaultValue: 0, // 默认乐器ID
        Required:     false,
        Description:  "乐器类型",
    },
    "volume": {
        Name:         "volume",
        Type:         ParamInt,
        DefaultValue: 100,
        Required:     false,
        Description:  "音量 (0-127)",
    },
}

// Section参数规范
var SectionParameters = map[string]ParameterSpec{
    "channel": {
        Name:         "channel",
        Type:         ParamInt,
        DefaultValue: 1,
        Required:     false,
        Description:  "MIDI通道",
    },
    "instrument": {
        Name:         "instrument",
        Type:         ParamInt,
        DefaultValue: 0, // 默认乐器ID
        Required:     false,
        Description:  "乐器类型",
    },
    "volume": {
        Name:         "volume",
        Type:         ParamInt,
        DefaultValue: 100,
        Required:     false,
        Description:  "音量 (0-127)",
    },
}

// Set设置节点
type SetNode struct {
    Parameters map[string]interface{}
    Context    ParameterContext
    Position   mytype.Position
}

var _ ASTNode = (*SetNode)(nil)

func (s *SetNode) String() string {
    return fmt.Sprintf("Set{Parameters: %v, Position: %s}", s.Parameters, s.Position)
}

func (s *SetNode) DetailedString(indent string) string {
    result := fmt.Sprintf("SetNode {\n")
    result += fmt.Sprintf("%s  上下文: %s\n", indent, s.getContextName())
    result += fmt.Sprintf("%s  位置: %s\n", indent, s.Position)
    
    if len(s.Parameters) > 0 {
        result += fmt.Sprintf("%s  参数:\n", indent)
        for key, value := range s.Parameters {
            result += fmt.Sprintf("%s    %s: %v (%T)\n", indent, key, value, value)
        }
    }
    
    // 显示解析后的参数
    if resolved, err := s.ResolveParameters(); err == nil {
        result += fmt.Sprintf("%s  解析后参数:\n", indent)
        for key, value := range resolved {
            result += fmt.Sprintf("%s    %s: %v (%T)\n", indent, key, value, value)
        }
    }
    
    result += fmt.Sprintf("%s}\n", indent)
    return result
}

func (s *SetNode) getContextName() string {
    switch s.Context {
    case GlobalContext:
        return "Global"
    case TrackContext:
        return "Track"
    case SectionContext:
        return "Section"
    default:
        return "Unknown"
    }
}
// 参数解析和验证
func (s *SetNode) ResolveParameters() (map[string]interface{}, error) {
    resolved := make(map[string]interface{})

    // 根据上下文选择参数规范
    var specs map[string]ParameterSpec
    switch s.Context {
    case GlobalContext:
        specs = GlobalParameters
    case TrackContext:
        specs = TrackParameters
    case SectionContext:
        specs = SectionParameters
    }

    // 验证和转换参数
    for key, value := range s.Parameters {
        spec, exists := specs[key]
        if !exists {
            return nil, fmt.Errorf("未知参数: %s", key)
        }

        // 类型验证和转换
        converted, err := s.convertParameter(value, spec.Type)
        if err != nil {
            return nil, fmt.Errorf("参数 %s 类型错误: %v", key, err)
        }

        resolved[key] = converted
    }

    // 添加默认值
    for key, spec := range specs {
        if _, exists := resolved[key]; !exists {
            resolved[key] = spec.DefaultValue
        }
    }

    return resolved, nil
}

func (s *SetNode) convertParameter(value interface{}, paramType ParameterType) (interface{}, error) {
    switch paramType {
    case ParamInt:
        if v, ok := value.(int); ok {
            return v, nil
        }
        return nil, fmt.Errorf("期望整数类型")
    case ParamFloat:
        if v, ok := value.(float64); ok {
            return v, nil
        }
        return nil, fmt.Errorf("期望浮点数类型")
    case ParamString:
        if v, ok := value.(string); ok {
            return v, nil
        }
        return nil, fmt.Errorf("期望字符串类型")
    case ParamBool:
        if v, ok := value.(bool); ok {
            return v, nil
        }
        return nil, fmt.Errorf("期望布尔类型")
    }
    return nil, fmt.Errorf("未知参数类型")
}