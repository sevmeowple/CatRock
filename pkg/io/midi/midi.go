package midi

import (
	"catRock/pkg/core"
	"catRock/pkg/io"
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // 导入 MIDI 驱动
)

type MIDIPlayer struct {
	Status    io.ConnectStatus
	driverOut drivers.Out
	sender    func(msg midi.Message) error
	bpm       int
	volume    int
	Channel   uint8             // MIDI 通道
	Program   core.InstrumentID // MIDI 程序号
}

var _ io.IO = (*MIDIPlayer)(nil) // 确保 MIDIPlayer 实现了 io.IO 接口

func NewMIDIPlayer() *MIDIPlayer {
	return &MIDIPlayer{
		Status:  io.Disconnected,
		Channel: 0, // 默认 MIDI 通道
		Program: 1, // 默认 MIDI 程序号
	}
}

func (p *MIDIPlayer) GetStatus() io.ConnectStatus {
	return p.Status
}

func (p *MIDIPlayer) Connect() (io.ConnectStatus, error) {
	// 获取可用的 MIDI 输出端口
	outports := midi.GetOutPorts()

	// 检查是否有可用端口
	if len(outports) == 0 {
		return io.Disconnected, fmt.Errorf("no MIDI output ports available")
	}

	// 选择第一个可用端口
	p.driverOut = outports[0]

	// 创建发送器
	sender, err := midi.SendTo(p.driverOut)
	if err != nil {
		return io.Disconnected, fmt.Errorf("failed to create MIDI sender: %v", err)
	}

	p.sender = sender
	p.Status = io.Connected

	fmt.Printf("Connected to MIDI device: %s\n", p.driverOut.String())
	return p.Status, nil
}

func (p *MIDIPlayer) Disconnect() (io.ConnectStatus, error) {
	if p.driverOut != nil {
		err := p.driverOut.Close()
		if err != nil {
			return io.Disconnected, err
		}
		p.driverOut = nil
		p.sender = nil
	}
	p.Status = io.Disconnected
	return io.Disconnected, nil
}

func (p *MIDIPlayer) PlayNote(params io.PlayNoteParams) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	// 智能切换通道和乐器
	targetChannel := p.Channel
	targetInstrument := p.Program
	// 检查是否需要切换通道
	if params.Note.Channel != 0 && params.Note.Channel != p.Channel {
		targetChannel = params.Note.Channel
	}

	// 检查是否需要切换乐器
	if params.Note.Instrument != 0 && params.Note.Instrument != int(p.Program) {
		targetInstrument = core.InstrumentID(params.Note.Instrument)
	}

	// 执行切换（如果需要）
	if targetChannel != p.Channel {
		err := p.SetChannel(targetChannel)
		if err != nil {
			return fmt.Errorf("failed to switch channel: %v", err)
		}
	}

	if targetInstrument != p.Program {
		err := p.SetProgram(targetInstrument)
		if err != nil {
			return fmt.Errorf("failed to switch instrument: %v", err)
		}
	}

	// 使用note的velocity，如果没有设置则使用player默认值
	velocity := params.Note.Velocity
	if velocity == 0 {
		velocity = uint8(p.volume)
	}

	// 发送 Note On
	err := p.sender(midi.NoteOn(p.Channel, uint8(params.Note.MIDINote[0]), velocity))
	if err != nil {
		return fmt.Errorf("failed to send note on: %v", err)
	}

	// 计算音符持续时间
	bpm := params.Bpm
	if bpm == 0 {
		bpm = p.bpm
	}

	duration := params.Note.Duration(float64(bpm))
	time.Sleep(duration)

	fmt.Printf("Playing note %s for duration %v at BPM %d on channel %d\n",
		params.Note.Name, duration, bpm, p.Channel)

	// 发送 Note Off
	err = p.sender(midi.NoteOff(p.Channel, uint8(params.Note.MIDINote[0])))
	if err != nil {
		return fmt.Errorf("failed to send note off: %v", err)
	}

	return nil
}

func (p *MIDIPlayer) PlayChord(params io.PlayChordParams) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	// 对和弦中的每个音符发送 Note On
	for _, note := range params.Chord.Notes {
		// 检查每个音符是否需要切换设置
		if note.Channel != 0 && note.Channel != p.Channel {
			p.SetChannel(note.Channel)
		}

		if note.Instrument != 0 && note.Instrument != int(p.Program) {
			p.SetProgram(core.InstrumentID(note.Instrument))
		}

		velocity := note.Velocity
		if velocity == 0 {
			velocity = uint8(p.volume)
		}

		err := p.sender(midi.NoteOn(p.Channel, uint8(note.MIDINote[0]), velocity))
		if err != nil {
			return fmt.Errorf("failed to send note on for chord: %v", err)
		}
	}

	// 计算和弦持续时间（使用第一个音符的时长）
	bpm := params.Bpm
	if bpm == 0 {
		bpm = p.bpm
	}

	duration := params.Chord.Notes[0].Duration(float64(bpm))
	time.Sleep(duration)

	fmt.Printf("Playing chord for duration %v at BPM %d on channel %d\n",
		duration, bpm, p.Channel)

	// 对和弦中的每个音符发送 Note Off
	for _, note := range params.Chord.Notes {
		err := p.sender(midi.NoteOff(p.Channel, uint8(note.MIDINote[0])))
		if err != nil {
			return fmt.Errorf("failed to send note off for chord: %v", err)
		}
	}

	return nil
}

func (p *MIDIPlayer) SetBPM(bpm int) error {
	p.bpm = bpm
	return nil
}

func (p *MIDIPlayer) SetVolume(volume int) error {
	p.volume = volume
	return nil
}

// 设置MIDI通道
func (p *MIDIPlayer) SetChannel(channel uint8) error {
	if channel > 15 {
		return fmt.Errorf("MIDI channel must be 0-15, got %d", channel)
	}
	p.Channel = channel
	fmt.Printf("Switched to MIDI channel %d\n", channel)
	return nil
}

// 设置乐器程序
func (p *MIDIPlayer) SetProgram(instrumentID core.InstrumentID) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	// 获取MIDI程序号
	midiProgram := core.GetMIDIProgram(instrumentID)

	// 发送Program Change消息
	err := p.sender(midi.ProgramChange(p.Channel, uint8(midiProgram)))
	if err != nil {
		return fmt.Errorf("failed to change program: %v", err)
	}

	p.Program = instrumentID

	return nil
}

// 实现低级MIDI事件发送
func (p *MIDIPlayer) SendNoteOn(channel uint8, note uint8, velocity uint8) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	return p.sender(midi.NoteOn(channel, note, velocity))
}

func (p *MIDIPlayer) SendNoteOff(channel uint8, note uint8, velocity uint8) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	return p.sender(midi.NoteOff(channel, note))
}

func (p *MIDIPlayer) SendProgramChange(channel uint8, program uint8) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	if channel == 9 {
		fmt.Printf("DEBUG: 跳过鼓组通道的 PROGRAM_CHANGE - 通道:%d\n", channel)
		return nil // 跳过鼓组通道的 Program Change
	}

	// 🔥 检查 program 是否为鼓组标识符
	if program == 0 && channel != 9 {
		fmt.Printf("DEBUG: Program 0 可能是鼓组，但通道不是 9 - 通道:%d, program:%d\n", channel, program)
	}

	return p.sender(midi.ProgramChange(channel, program))
}

func (p *MIDIPlayer) SendControlChange(channel uint8, controller uint8, value uint8) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	if p.sender == nil {
		return fmt.Errorf("MIDI sender not initialized")
	}

	return p.sender(midi.ControlChange(channel, controller, value))
}

// 新增：批量事件播放（可选的高级方法）
func (p *MIDIPlayer) PlayEvents(params io.PlayEventsParams) error {
	if p.Status != io.Connected {
		return fmt.Errorf("MIDI player not connected")
	}

	// 这里可以实现更高效的批量播放逻辑
	// 比如预处理、优化、或者异步播放
	for _, event := range params.Events {
		err := p.executeEvent(event, params.BPM)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *MIDIPlayer) executeEvent(event io.Event, bpm float64) error {
	switch event.Type {
	case io.NOTE_ON_EVENT:
		return p.SendNoteOn(event.Channel, event.Data1, event.Data2)
	case io.NOTE_OFF_EVENT:
		return p.SendNoteOff(event.Channel, event.Data1, event.Data2)
	case io.PROGRAM_CHANGE_EVENT:
		return p.SendProgramChange(event.Channel, event.Data1)
	case io.CONTROL_CHANGE_EVENT:
		return p.SendControlChange(event.Channel, event.Data1, event.Data2)
	default:
		return fmt.Errorf("unsupported event type: %v", event.Type)
	}
}
