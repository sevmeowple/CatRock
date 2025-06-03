package core

type InstrumentFamily int

const (
	Piano               InstrumentFamily = iota // 0-7: 钢琴类
	ChromaticPercussion                         // 8-15: 色彩打击乐
	Organ                                       // 16-23: 风琴类
	Guitar                                      // 24-31: 吉他类
	Bass                                        // 32-39: 贝斯类
	Strings                                     // 40-47: 弦乐类
	Ensemble                                    // 48-55: 合奏类
	Brass                                       // 56-63: 铜管类
	Reed                                        // 64-71: 簧片类
	Pipe                                        // 72-79: 管乐类
	SynthLead                                   // 80-87: 合成主音
	SynthPad                                    // 88-95: 合成铺底
	SynthEffects                                // 96-103: 合成效果
	Ethnic                                      // 104-111: 民族乐器
	Percussive                                  // 112-119: 打击乐器
	SoundEffects                                // 120-127: 音效
	DrumKits                                    // 128+: 鼓组（自定义范围）
)

type InstrumentType int

const (
	MelodicInstrument InstrumentType = iota
	DrumKit
)

type InstrumentID int

type Instrument struct {
	ID          InstrumentID // 唯一标识符
	MIDIProgram int          // MIDI Program Number (0-127)
	Name        string       // 英文名
	ChineseName string       // 中文名
	Family      InstrumentFamily
	Type        InstrumentType
	Channel     uint8
	Range       Range
	Polyphonic  bool
	Expressive  bool
	Tags        []string
}

type Range struct {
	Lowest  Note // 最低音
	Highest Note // 最高音
}

// 只定义最基本的判断函数
func IsDrumKit(instrumentID InstrumentID) bool {
	return instrumentID >= 128 // 128+ 表示鼓组
}

func GetMIDIProgram(instrumentID InstrumentID) uint8 {
	if IsDrumKit(instrumentID) {
		return 0 // 鼓组 Program 无关紧要
	}

	if instrumentID > 127 {
		return 0 // 超出范围的ID直接返回0
	}

	return uint8(instrumentID) // 直接映射
}

func GetDefaultChannel(instrumentID InstrumentID) uint8 {
	if IsDrumKit(instrumentID) {
		return 9 // 鼓组固定通道 9
	}
	return 0 // 默认通道 0（用户可以自己设置）
}
