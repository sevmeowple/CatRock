package score

import (
    "fmt"
)

// GroupElement 表示一个播放元素组
type GroupElement struct {
	ID string // 可选的组ID
    elements []Playable
    duration *float64 // 可选的组整体时值
}

// 实现 Playable 接口
var _ Playable = (*GroupElement)(nil)

func (g *GroupElement) GetID() string {
	if g.ID != "" {
		return g.ID
	}
	// 如果没有指定ID，生成一个基于元素数量的默认ID
	return fmt.Sprintf("group_%d_elements", len(g.elements))
}

func (g *GroupElement) GetType() PlayableType {
	return GROUP_TYPE
}

// 构造函数
func NewGroupElement() *GroupElement {
    return &GroupElement{
        elements: make([]Playable, 0),
        duration: nil,
    }
}

func NewGroupElementWithDuration(duration float64) *GroupElement {
    return &GroupElement{
        elements: make([]Playable, 0),
        duration: &duration,
    }
}

// 添加元素
func (g *GroupElement) AddElement(element Playable) {
    if element != nil {
        g.elements = append(g.elements, element)
    }
}

func (g *GroupElement) AddElements(elements ...Playable) {
    for _, element := range elements {
        g.AddElement(element)
    }
}

// 获取元素
func (g *GroupElement) GetElements() []Playable {
    return g.elements
}

func (g *GroupElement) ElementCount() int {
    return len(g.elements)
}

// 实现 Playable 接口
func (g *GroupElement) Duration(context PlayContext) float64 {
    if g.duration != nil {
        // 如果指定了组时值，返回指定值
        return *g.duration
    }
    
    // 否则计算所有元素的总时长（串行播放）
    totalDuration := 0.0
    for _, element := range g.elements {
        totalDuration += element.Duration(context)
    }
    return totalDuration
}

func (g *GroupElement) GenerateEvents(startTime float64, context PlayContext) []Event {
    var events []Event
    currentTime := startTime
    
    if g.duration != nil {
        // 有指定组时值 - 需要分配时间
        return g.generateEventsWithGroupDuration(startTime, context)
    } else {
        // 无指定组时值 - 简单串行播放
        for _, element := range g.elements {
            elementEvents := element.GenerateEvents(currentTime, context)
            events = append(events, elementEvents...)
            currentTime += element.Duration(context)
        }
    }
    
    return events
}

// 带组时值的事件生成
func (g *GroupElement) generateEventsWithGroupDuration(startTime float64, context PlayContext) []Event {
    if len(g.elements) == 0 || g.duration == nil {
        return []Event{}
    }
    
    var events []Event
    groupDuration := *g.duration
    
    // 计算每个元素的原始时长
    totalOriginalDuration := 0.0
    elementDurations := make([]float64, len(g.elements))
    
    for i, element := range g.elements {
        elementDurations[i] = element.Duration(context)
        totalOriginalDuration += elementDurations[i]
    }
    
    // 按比例分配组时值
    currentTime := startTime
    for i, element := range g.elements {
        if totalOriginalDuration > 0 {
            // 按原始时长比例分配
            scaledDuration := (elementDurations[i] / totalOriginalDuration) * groupDuration
            
            // 创建新的上下文，调整时值
            newContext := context
            // 这里可能需要调整上下文来影响元素的实际播放时长
            
            elementEvents := element.GenerateEvents(currentTime, newContext)
            events = append(events, elementEvents...)
            currentTime += scaledDuration
        }
    }
    
    return events
}

// 调试和显示
func (g *GroupElement) String() string {
    if g.duration != nil {
        return fmt.Sprintf("Group(%d elements, duration: %.2f)", len(g.elements), *g.duration)
    }
    return fmt.Sprintf("Group(%d elements)", len(g.elements))
}

// 验证组完整性
func (g *GroupElement) Validate() error {
    if len(g.elements) == 0 {
        return fmt.Errorf("空组")
    }
    
    for i, element := range g.elements {
        if element == nil {
            return fmt.Errorf("组元素 %d 为 nil", i)
        }
    }
    
    if g.duration != nil && *g.duration <= 0 {
        return fmt.Errorf("组时值必须大于0，当前值: %.2f", *g.duration)
    }
    
    return nil
}

// 设置组时值
func (g *GroupElement) SetDuration(duration float64) {
    if duration > 0 {
        g.duration = &duration
    }
}

func (g *GroupElement) ClearDuration() {
    g.duration = nil
}

func (g *GroupElement) HasDuration() bool {
    return g.duration != nil
}

func (g *GroupElement) DetailedString(indent string) string {
    result := fmt.Sprintf("Group {\n")
    result += fmt.Sprintf("%s  ID: %s\n", indent, g.GetID())
    result += fmt.Sprintf("%s  元素数量: %d\n", indent, len(g.elements))
    
    if g.duration != nil {
        result += fmt.Sprintf("%s  指定时长: %.3f拍\n", indent, *g.duration)
    } else {
        result += fmt.Sprintf("%s  计算时长: %.3f拍\n", indent, g.Duration(PlayContext{}))
    }
    
    if len(g.elements) > 0 {
        result += fmt.Sprintf("%s  元素:\n", indent)
        for i, element := range g.elements {
            result += fmt.Sprintf("%s    [%d] %s", indent, i, element.DetailedString(indent+"      "))
        }
    }
    
    result += fmt.Sprintf("%s}\n", indent)
    return result
}