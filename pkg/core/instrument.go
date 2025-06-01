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

type InstrumentID int

// 普通乐器定义（按 GM 标准）
const (
	// Piano Family (0-7)
	AcousticGrandPiano  InstrumentID = 0 // 大钢琴
	BrightAcousticPiano InstrumentID = 1 // 明亮大钢琴
	ElectricGrandPiano  InstrumentID = 2 // 电子大钢琴
	HonkyTonkPiano      InstrumentID = 3 // 酒吧钢琴
	ElectricPiano1      InstrumentID = 4 // 电钢琴1
	ElectricPiano2      InstrumentID = 5 // 电钢琴2
	Harpsichord         InstrumentID = 6 // 羽管键琴
	Clavinet            InstrumentID = 7 // 古钢琴

	// Chromatic Percussion (8-15)
	Celesta      InstrumentID = 8  // 钢片琴
	Glockenspiel InstrumentID = 9  // 钟琴
	MusicBox     InstrumentID = 10 // 音乐盒
	Vibraphone   InstrumentID = 11 // 颤音琴
	Marimba      InstrumentID = 12 // 马林巴
	Xylophone    InstrumentID = 13 // 木琴
	TubularBells InstrumentID = 14 // 管钟
	Dulcimer     InstrumentID = 15 // 洋琴

	// Guitar Family (24-31)
	AcousticGuitarNylon InstrumentID = 24 // 古典吉他
	AcousticGuitarSteel InstrumentID = 25 // 民谣吉他
	ElectricGuitarJazz  InstrumentID = 26 // 爵士电吉他
	ElectricGuitarClean InstrumentID = 27 // 清音电吉他
	ElectricGuitarMuted InstrumentID = 28 // 闷音电吉他
	OverdrivenGuitar    InstrumentID = 29 // 过载吉他
	DistortionGuitar    InstrumentID = 30 // 失真吉他
	GuitarHarmonics     InstrumentID = 31 // 吉他泛音

	// String Family (40-47)
	Violin                       InstrumentID = 40 // 小提琴
	Viola                        InstrumentID = 41 // 中提琴
	Cello                        InstrumentID = 42 // 大提琴
	ContrabassInstrumentID                    = 43 // 低音提琴
	TremoloStringsInstrumentID                = 44 // 颤音弦乐
	PizzicatoStringsInstrumentID              = 45 // 拨弦
	OrchestralHarpInstrumentID                = 46 // 竖琴
	Timpani                      InstrumentID = 47 // 定音鼓

	// Brass Family (56-63)
	Trumpet      InstrumentID = 56 // 小号
	Trombone     InstrumentID = 57 // 长号
	Tuba         InstrumentID = 58 // 大号
	MutedTrumpet InstrumentID = 59 // 弱音小号
	FrenchHorn   InstrumentID = 60 // 圆号
	BrassSection InstrumentID = 61 // 铜管组
	SynthBrass1  InstrumentID = 62 // 合成铜管1
	SynthBrass2  InstrumentID = 63 // 合成铜管2
)

// 鼓组定义（使用 128+ 范围避免冲突）
const (
	StandardDrumKit   InstrumentID = 128 // 标准鼓组
	RoomDrumKit       InstrumentID = 129 // 房间鼓组
	PowerDrumKit      InstrumentID = 130 // 强力鼓组
	ElectronicDrumKit InstrumentID = 131 // 电子鼓组
	TR808DrumKit      InstrumentID = 132 // TR-808鼓组
)

type InstrumentType int

const (
	MelodicInstrument InstrumentType = iota
	DrumKit
)

type Instrument struct {
	ID          InstrumentID     // 唯一标识符
	MIDIProgram int              // MIDI Program Number (0-127)
	Name        string           // 英文名
	ChineseName string           // 中文名
	Family      InstrumentFamily // 乐器族
	Type        InstrumentType   // 乐器类型
	Channel     uint8            // 默认 MIDI 通道
	Range       Range            // 音域
	Polyphonic  bool             // 是否支持和声
	Expressive  bool             // 是否有表情控制
	Tags        []string         // 标签（古典、流行、民族等）
}

type Range struct {
	Lowest  Note // 最低音
	Highest Note // 最高音
}

// 预定义的常用乐器
var Instruments = map[InstrumentID]Instrument{
	AcousticGrandPiano: {
		ID:          AcousticGrandPiano,
		MIDIProgram: 0, // MIDI Program 0
		Name:        "Acoustic Grand Piano",
		ChineseName: "三角钢琴",
		Family:      Piano,
		Type:        MelodicInstrument,
		Channel:     0,
		Range: Range{
			Lowest:  NewNote(NewNoteParams{Name: A, Octave: 0, Accidental: Natural}),
			Highest: NewNote(NewNoteParams{Name: C, Octave: 8, Accidental: Natural}),
		},
		Polyphonic: true,
		Expressive: true,
		Tags:       []string{"classical", "popular", "versatile"},
	},

	Violin: {
		ID:          Violin,
		MIDIProgram: 40, // MIDI Program 40
		Name:        "Violin",
		ChineseName: "小提琴",
		Family:      Strings,
		Type:        MelodicInstrument,
		Channel:     1,
		Range: Range{
			Lowest:  NewNote(NewNoteParams{Name: G, Octave: 3, Accidental: Natural}),
			Highest: NewNote(NewNoteParams{Name: A, Octave: 7, Accidental: Natural}),
		},
		Polyphonic: false, // 主要是单音
		Expressive: true,
		Tags:       []string{"classical", "orchestral", "expressive"},
	},

	StandardDrumKit: {
		ID:          StandardDrumKit,
		MIDIProgram: 0, // 鼓组在通道9，Program无关紧要
		Name:        "Standard Drum Kit",
		ChineseName: "标准鼓组",
		Family:      DrumKits,
		Type:        DrumKit,
		Channel:     9,    // 鼓组通常在通道9
		Polyphonic:  true, // 可以同时敲击多个鼓
		Expressive:  false,
		Tags:        []string{"rhythm", "popular", "rock"},
	},
}

// 辅助函数：判断是否为鼓组
func IsDrumKit(instrumentID InstrumentID) bool {
	if instrument, exists := Instruments[instrumentID]; exists {
		return instrument.Type == DrumKit
	}
	return false
}

// 辅助函数：获取MIDI Program Number
func GetMIDIProgram(instrumentID InstrumentID) int {
	if instrument, exists := Instruments[instrumentID]; exists {
		return instrument.MIDIProgram
	}
	return 0 // 默认钢琴
}

// 辅助函数：获取默认通道
func GetDefaultChannel(instrumentID InstrumentID) uint8 {
	if instrument, exists := Instruments[instrumentID]; exists {
		return instrument.Channel
	}
	return 0 // 默认通道0
}
