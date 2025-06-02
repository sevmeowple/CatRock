package core

import "fmt"

// Rest 休止符结构
type Rest struct {
    Beat     BeatValue // 休止符时值
    Position int       // 在乐谱中的位置（可选）
}

// NewRestParams 创建休止符的参数
type NewRestParams struct {
    Beat     BeatValue // 休止符时值
    Position int       // 位置（可选）
}

// NewRest 创建一个新的休止符
func NewRest(beat BeatValue) Rest {
    return Rest{
        Beat:     beat,
        Position: 0,
    }
}

// NewRestWithParams 使用参数创建休止符
func NewRestWithParams(params NewRestParams) Rest {
    rest := Rest{
        Beat:     params.Beat,
        Position: params.Position,
    }
    
    // 验证时值
    if rest.Beat == 0 {
        rest.Beat = Quarter // 默认四分休止符
    }
    
    return rest
}

// Duration 返回休止符的持续时间（以拍为单位）
func (r Rest) Duration() float64 {
    return float64(r.Beat)
}

// String 返回休止符的字符串表示
func (r Rest) String() string {
    beatName := r.getBeatName()
    if r.Position > 0 {
        return fmt.Sprintf("Rest{%s, Pos: %d}", beatName, r.Position)
    }
    return fmt.Sprintf("Rest{%s}", beatName)
}

// getBeatName 返回节拍值的名称
func (r Rest) getBeatName() string {
    switch r.Beat {
    case Whole:
        return "whole"
    case Half:
        return "half"
    case Quarter:
        return "quarter"
    case Eighth:
        return "eighth"
    case Sixteenth:
        return "sixteenth"
    case ThirtySecond:
        return "thirty-second"
    default:
        return fmt.Sprintf("unknown(%d)", int(r.Beat))
    }
}

// IsEqual 比较两个休止符是否相等
func (r Rest) IsEqual(other Rest) bool {
    return r.Beat == other.Beat && r.Position == other.Position
}

// Clone 克隆休止符
func (r Rest) Clone() Rest {
    return Rest{
        Beat:     r.Beat,
        Position: r.Position,
    }
}

// WithBeat 返回修改时值后的新休止符
func (r Rest) WithBeat(beat BeatValue) Rest {
    newRest := r.Clone()
    newRest.Beat = beat
    return newRest
}

// WithPosition 返回修改位置后的新休止符  
func (r Rest) WithPosition(position int) Rest {
    newRest := r.Clone()
    newRest.Position = position
    return newRest
}

// 便捷构造函数

// NewWholeRest 创建全休止符
func NewWholeRest() Rest {
    return NewRest(Whole)
}

// NewHalfRest 创建二分休止符
func NewHalfRest() Rest {
    return NewRest(Half)
}

// NewQuarterRest 创建四分休止符
func NewQuarterRest() Rest {
    return NewRest(Quarter)
}

// NewEighthRest 创建八分休止符
func NewEighthRest() Rest {
    return NewRest(Eighth)
}

// NewSixteenthRest 创建十六分休止符
func NewSixteenthRest() Rest {
    return NewRest(Sixteenth)
}