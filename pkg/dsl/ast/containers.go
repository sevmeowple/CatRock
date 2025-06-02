package ast

import (
	"catRock/pkg/dsl/mytype"
	"catRock/pkg/score"
	"fmt"
)

// Track节点 - 并行播放容器
type TrackNode struct {
    Name     string
    Sets     []*SetNode     // track内的设置
    Elements []PlayableNode // track内的可播放元素
    Position mytype.Position
}

var _ ContainerNode = (*TrackNode)(nil)

func (t *TrackNode) String() string {
    return fmt.Sprintf("Track{Name: %s, Sets: %d, Elements: %d}", 
        t.Name, len(t.Sets), len(t.Elements))
}

func (t *TrackNode) ToPlayable() score.Playable {
    // 解析所有Set块，合并参数
    params := t.MergeSetParameters()

    track := score.NewTrack(t.Name)

    // 应用参数设置
    applyContainerParams(track, params)

    // 添加所有子元素
    for _, element := range t.Elements {
        track.AddElement(element.ToPlayable())
    }

    return track
}

func (t *TrackNode) AddElement(element PlayableNode) {
    t.Elements = append(t.Elements, element)
}

// 合并所有Set块的参数
func (t *TrackNode) MergeSetParameters() map[string]interface{} {
    merged := make(map[string]interface{})

    // 按顺序应用所有Set块
    for _, setBlock := range t.Sets {
        setBlock.Context = TrackContext // 设置上下文
        resolved, err := setBlock.ResolveParameters()
        if err != nil {
            // 处理错误，可以记录日志或继续
            continue
        }

        // 后面的参数覆盖前面的
        for key, value := range resolved {
            merged[key] = value
        }
    }

    return merged
}

// Section节点 - 顺序播放容器
type SectionNode struct {
    Name     string
    Sets     []*SetNode     // section内的设置
    Elements []PlayableNode // section内的可播放元素
    Position mytype.Position
}

var _ ContainerNode = (*SectionNode)(nil)

func (s *SectionNode) String() string {
    return fmt.Sprintf("Section{Name: %s, Sets: %d, Elements: %d}", 
        s.Name, len(s.Sets), len(s.Elements))
}

func (s *SectionNode) ToPlayable() score.Playable {
    // 解析所有Set块，合并参数
    params := s.MergeSetParameters()

    section := score.NewSection(s.Name)

    // 应用参数设置
    applyContainerParams(section, params)

    // 添加所有子元素
    for _, element := range s.Elements {
        section.AddElement(element.ToPlayable())
    }

    return section
}

func (s *SectionNode) AddElement(element PlayableNode) {
    s.Elements = append(s.Elements, element)
}

// 合并所有Set块的参数
func (s *SectionNode) MergeSetParameters() map[string]interface{} {
    merged := make(map[string]interface{})

    // 按顺序应用所有Set块
    for _, setBlock := range s.Sets {
        setBlock.Context = SectionContext // 设置上下文
        resolved, err := setBlock.ResolveParameters()
        if err != nil {
            // 处理错误
            continue
        }

        // 后面的参数覆盖前面的
        for key, value := range resolved {
            merged[key] = value
        }
    }

    return merged
}

// 统一的容器参数应用函数
func applyContainerParams(container interface{}, params map[string]interface{}) {
    // 使用类型断言来应用参数
    if c, ok := container.(interface{ SetVolume(int) }); ok {
        if volume := getVolume(params); volume > 0 {
            c.SetVolume(volume)
        }
    }
    
    if c, ok := container.(interface{ SetChannel(int) }); ok {
        if channel := getChannel(params); channel > 0 {
            c.SetChannel(channel)
        }
    }
    
    if c, ok := container.(interface{ SetBPM(float64) }); ok {
        if bpm, ok := params["BPM"]; ok {
            if bpmFloat, ok := bpm.(float64); ok {
                c.SetBPM(bpmFloat)
            }
        }
    }
}