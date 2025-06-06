package score

import (
	"catRock/pkg/core"
	"catRock/pkg/io"
	"fmt"
	"sort"
	"time"
)

// Score 表示一个完整的音乐作品
type Score struct {
	// 元数据
	Title    string
	Composer string
	Year     int

	// 播放设置
	BPM    float64
	Volume int

	// 根元素 - 整个作品的入口
	RootElement Playable

	// 可选的全局设置
	GlobalSettings map[string]interface{}
}

// 播放引擎
type PlayEngine struct {
	score   *Score
	context PlayContext
	events  []Event
}

// 构造函数
func NewScore(title string) *Score {
	return &Score{
		Title:          title,
		Composer:       "",
		Year:           time.Now().Year(),
		BPM:            120.0,
		Volume:         100,
		RootElement:    nil,
		GlobalSettings: make(map[string]interface{}),
	}
}

// Score的基本方法
func (s *Score) SetBPM(bpm float64) {
	if bpm > 0 {
		s.BPM = bpm
	}
}

func (s *Score) SetVolume(volume int) {
	if volume >= 0 && volume <= 127 {
		s.Volume = volume
	}
}

func (s *Score) SetRootElement(element Playable) {
	s.RootElement = element
}

func (s *Score) GetDuration() float64 {
	if s.RootElement == nil {
		return 0.0
	}

	context := s.createPlayContext()
	return s.RootElement.Duration(context)
}

func (s *Score) String() string {
	return fmt.Sprintf("Score{Title: %s, Composer: %s, BPM: %.1f, Duration: %.2fs}",
		s.Title, s.Composer, s.BPM, s.GetDuration())
}

// 新增：详细字符串方法
func (s *Score) DetailedString(indent string) string {
	result := fmt.Sprintf("%sScore {\n", indent)
	result += fmt.Sprintf("%s  标题: %s\n", indent, s.Title)
	result += fmt.Sprintf("%s  作曲: %s\n", indent, s.Composer)
	result += fmt.Sprintf("%s  年份: %d\n", indent, s.Year)
	result += fmt.Sprintf("%s  BPM: %.1f\n", indent, s.BPM)
	result += fmt.Sprintf("%s  音量: %d\n", indent, s.Volume)
	result += fmt.Sprintf("%s  时长: %.2f秒\n", indent, s.GetDuration())

	if len(s.GlobalSettings) > 0 {
		result += fmt.Sprintf("%s  全局设置:\n", indent)
		for key, value := range s.GlobalSettings {
			result += fmt.Sprintf("%s    %s: %v (%T)\n", indent, key, value, value)
		}
	}

	if s.RootElement != nil {
		result += fmt.Sprintf("%s  根元素:\n", indent)
		result += s.RootElement.DetailedString(indent + "    ")
	} else {
		result += fmt.Sprintf("%s  根元素: (无)\n", indent)
	}

	result += fmt.Sprintf("%s}\n", indent)
	return result
}

// 创建播放上下文
func (s *Score) createPlayContext() PlayContext {
	return NewPlayContext(s.BPM, s.Volume)
}

// 播放

// 使用指定的IO设备播放Score
func (s *Score) PlayWithIO(ioDevice io.IO) error {
	// 生成事件
	engine := NewPlayEngine(s)
	events, err := engine.GenerateEvents()
	if err != nil {
		return fmt.Errorf("生成事件失败: %v", err)
	}

	// 执行播放
	return engine.PlayEventsWithIO(ioDevice, events)
}


// PlayEngine的播放方法
func (pe *PlayEngine) PlayEventsWithIO(ioDevice io.IO, events []Event) error {
	if len(events) == 0 {
		return fmt.Errorf("没有事件可播放")
	}

	// 连接检查
	if status := ioDevice.GetStatus(); status != io.Connected {
		return fmt.Errorf("IO设备未连接")
	}

	// 计算每拍的毫秒数
	millisecondsPerBeat := 60000.0 / pe.score.BPM

	startTime := time.Now()

	for _, event := range events {
		// 计算事件的实际播放时间
		eventTime := time.Duration(event.Time*millisecondsPerBeat) * time.Millisecond
		targetTime := startTime.Add(eventTime)

		// 等待到事件时间
		now := time.Now()
		if targetTime.After(now) {
			time.Sleep(targetTime.Sub(now))
		}

		// 执行事件
		if err := pe.executeEventWithIO(ioDevice, event); err != nil {
			return fmt.Errorf("执行事件失败: %v", err)
		}
	}

	return nil
}

// 异步播放事件
func (pe *PlayEngine) PlayEventsWithIOAsync(ioDevice io.IO, events []Event) (<-chan error,error) {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		if err := pe.PlayEventsWithIO(ioDevice, events); err != nil {
			errChan <- err
		}
	}()

	return errChan, nil
}

// 执行单个事件
func (pe *PlayEngine) executeEventWithIO(ioDevice io.IO, event Event) error {
	switch event.Action {
	case NOTE_ON:
		if midiNote, ok := event.Data.(uint8); ok {
			return ioDevice.SendNoteOn(uint8(event.Channel), midiNote, event.Velocity)
		}
		return fmt.Errorf("NOTE_ON事件数据类型错误")

	case NOTE_OFF:
		if midiNote, ok := event.Data.(uint8); ok {
			return ioDevice.SendNoteOff(uint8(event.Channel), midiNote, 0)
		}
		return fmt.Errorf("NOTE_OFF事件数据类型错误")

	case VOLUME_CHANGE:
		if volume, ok := event.Data.(uint8); ok {
			return ioDevice.SendControlChange(uint8(event.Channel), 7, volume) // CC7 = 主音量
		}
		return fmt.Errorf("VOLUME_CHANGE事件数据类型错误")

	case PROGRAM_CHANGE:
		if program, ok := event.Data.(core.InstrumentID); ok {
			if core.IsDrumKit(program) {
				if event.Channel != 9 {
					fmt.Printf("WARNING: 鼓组程序变更事件不应在非鼓组通道发送 - 通道:%d, 程序:%d\n", event.Channel, program)
					return nil // 跳过非鼓组通道的程序变更
				}
				return nil
			}

			midiProgram := core.GetMIDIProgram(program)
			return ioDevice.SendProgramChange(uint8(event.Channel), midiProgram)
		}
		return fmt.Errorf("PROGRAM_CHANGE事件数据类型错误")

	default:
		// 忽略未知事件类型
		return nil
	}
}

// 播放引擎方法
func NewPlayEngine(score *Score) *PlayEngine {
	return &PlayEngine{
		score:   score,
		context: score.createPlayContext(),
		events:  []Event{},
	}
}

func (pe *PlayEngine) GenerateEvents() ([]Event, error) {
	if pe.score.RootElement == nil {
		return nil, fmt.Errorf("没有可播放的音乐元素")
	}

	// 生成所有事件
	events := pe.score.RootElement.GenerateEvents(0.0, pe.context)

	// 排序事件
	pe.events = pe.sortEvents(events)

	return pe.events, nil
}

func (pe *PlayEngine) GetEvents() []Event {
	return pe.events
}

func (pe *PlayEngine) sortEvents(events []Event) []Event {
	sort.Slice(events, func(i, j int) bool {
		// 按时间排序
		if events[i].Time != events[j].Time {
			return events[i].Time < events[j].Time
		}

		// 同一时间：NOTE_ON优先于NOTE_OFF
		if events[i].Action != events[j].Action {
			return events[i].Action == NOTE_ON
		}

		// 其他情况按源元素排序
		return events[i].SourceElement < events[j].SourceElement
	})

	return events
}

// 播放统计信息
type PlayStats struct {
	TotalDuration float64
	TotalNotes    int
	TotalChords   int
	TotalRests    int
	TotalTracks   int
	TotalSections int
	ChannelsUsed  []int
}

// 验证Score完整性
func (s *Score) Validate() error {
	if s.Title == "" {
		return fmt.Errorf("缺少作品标题")
	}

	if s.BPM <= 0 {
		return fmt.Errorf("BPM必须大于0，当前值: %.1f", s.BPM)
	}

	if s.Volume < 0 || s.Volume > 127 {
		return fmt.Errorf("音量必须在0-127范围内，当前值: %d", s.Volume)
	}

	if s.RootElement == nil {
		return fmt.Errorf("没有音乐内容")
	}

	// 验证根元素是否能正常生成事件
	context := s.createPlayContext()
	duration := s.RootElement.Duration(context)
	if duration <= 0 {
		return fmt.Errorf("音乐内容时长为0")
	}

	return nil
}

// 导出功能（为后续扩展预留）
type ExportFormat int

const (
	MIDI ExportFormat = iota
	JSON
	XML
)

type ExportOptions struct {
	Format   ExportFormat
	FileName string
	Options  map[string]interface{}
}

func (s *Score) Export(options ExportOptions) error {
	if err := s.Validate(); err != nil {
		return fmt.Errorf("导出前验证失败: %v", err)
	}

	switch options.Format {
	case MIDI:
		return s.exportMIDI(options)
	case JSON:
		return s.exportJSON(options)
	case XML:
		return s.exportXML(options)
	default:
		return fmt.Errorf("不支持的导出格式: %v", options.Format)
	}
}

// 导出实现（占位符）
func (s *Score) exportMIDI(options ExportOptions) error {
	// TODO: 实现MIDI导出
	return fmt.Errorf("MIDI导出尚未实现")
}

func (s *Score) exportJSON(options ExportOptions) error {
	// TODO: 实现JSON导出
	return fmt.Errorf("JSON导出尚未实现")
}

func (s *Score) exportXML(options ExportOptions) error {
	// TODO: 实现XML导出
	return fmt.Errorf("XML导出尚未实现")
}
