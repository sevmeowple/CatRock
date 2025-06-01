package test

import (
	"catRock/pkg/dsl"
	"catRock/pkg/io"
	"catRock/pkg/io/midi"
	"catRock/pkg/score"
	"strings"
	"testing"
)

func TestSimpleDSLParsing(t *testing.T) {
    input := `BPM: 120
C4 quarter
D4 quarter
E4 quarter
F4 quarter
G4 half
`
    
    // 1. 词法分析测试
    t.Log("=== 词法分析测试 ===")
    lexer := dsl.NewLexer(input)
    
    expectedTokens := []dsl.TokenType{
        dsl.BPM, dsl.COLON, dsl.NUMBER, dsl.NEWLINE,
        dsl.NOTE_C, dsl.NUMBER, dsl.IDENTIFIER, dsl.NEWLINE,
        dsl.NOTE_D, dsl.NUMBER, dsl.IDENTIFIER, dsl.NEWLINE,
        dsl.NOTE_E, dsl.NUMBER, dsl.IDENTIFIER, dsl.NEWLINE,
        dsl.NOTE_F, dsl.NUMBER, dsl.IDENTIFIER, dsl.NEWLINE,
        dsl.NOTE_G, dsl.NUMBER, dsl.IDENTIFIER, dsl.NEWLINE,
        dsl.EOF,
    }
    
    for i, expectedType := range expectedTokens {
        token := lexer.NextToken()
        if token.Type != expectedType {
            t.Errorf("Token %d: 期望 %s, 得到 %s (%s)", 
                i, expectedType, token.Type, token.Literal)
        } else {
            t.Logf("Token %d: ✅ %s", i, token)
        }
    }
    
    // 2. 语法分析测试
    t.Log("\n=== 语法分析测试 ===")
    lexer2 := dsl.NewLexer(input)
    parser := dsl.NewParser(lexer2)
    
    ast := parser.ParseScore()
    
    if len(parser.Errors()) > 0 {
        t.Fatalf("解析错误: %v", parser.Errors())
    }
    
    if ast == nil {
        t.Fatal("AST为空")
    }
    
    // 验证AST结构
    if ast.Header.BPM != 120 {
        t.Errorf("BPM错误: 期望120，得到%d", ast.Header.BPM)
    }
    
    if len(ast.Body) != 5 {
        t.Errorf("音符数量错误: 期望5，得到%d", len(ast.Body))
    }
    
    t.Logf("AST: %s", ast)
    
    // 3. 代码生成测试
    t.Log("\n=== 代码生成测试 ===")
    generator := dsl.NewGenerator()
    scoreObj, err := generator.GenerateScore(ast)
    
    if err != nil {
        t.Fatalf("代码生成错误: %v", err)
    }
    
    if len(generator.Errors()) > 0 {
        t.Fatalf("生成器错误: %v", generator.Errors())
    }
    
    // 验证Score对象
    if scoreObj.BPM != 120 {
        t.Errorf("Score BPM错误: 期望120，得到%.0f", scoreObj.BPM)
    }
    
    if scoreObj.RootElement == nil {
        t.Fatal("RootElement为空")
    }
    
    t.Logf("生成的Score: BPM=%.0f, Title=%s", scoreObj.BPM, scoreObj.Title)
    
    // 4. 事件生成测试
    t.Log("\n=== 事件生成测试 ===")
    events, err := scoreObj.Play()
    if err != nil {
        t.Fatalf("事件生成失败: %v", err)
    }
    
    t.Logf("生成了 %d 个事件", len(events))
    t.Logf("预期播放时长: %.2f 拍", scoreObj.Duration())
    
    // 验证事件
    noteOnCount := 0
    noteOffCount := 0
    for _, event := range events {
        if event.Action == score.NOTE_ON {
            noteOnCount++
        } else if event.Action == score.NOTE_OFF {
            noteOffCount++
        }
        
        if noteOnCount <= 5 { // 只显示前5个事件
            t.Logf("事件: 时间%.2f, 动作%v, 数据%v", 
                event.Time, event.Action, event.Data)
        }
    }
    
    if noteOnCount != 5 {
        t.Errorf("NOTE_ON数量错误: 期望5，得到%d", noteOnCount)
    }
    
    if noteOffCount != 5 {
        t.Errorf("NOTE_OFF数量错误: 期望5，得到%d", noteOffCount)
    }
    
    // 5. 实际播放测试（可选）
    if !testing.Short() {
        t.Log("\n=== MIDI播放测试 ===")
        
        midiPlayer := midi.NewMIDIPlayer()
        status, err := midiPlayer.Connect()
        if err == nil && status == io.Connected { // Connected
            defer midiPlayer.Disconnect()
            
            playEngine := score.NewPlayEngine(midiPlayer, scoreObj.BPM)
            err = playEngine.PlayEvents(events)
            
            if err != nil {
                t.Errorf("播放失败: %v", err)
            } else {
                t.Log("✅ 播放成功")
            }
        } else {
            t.Skip("MIDI不可用，跳过播放测试")
        }
    }
}

func TestDSLErrorHandling(t *testing.T) {
    testCases := []struct{
        name     string
        input    string
        expectError bool
    }{
        {
            name:  "缺少BPM",
            input: `C4 quarter`,
            expectError: true,
        },
        {
            name:  "无效八度",
            input: `BPM: 120\nC99 quarter\n`,
            expectError: true,
        },
        {
            name:  "无效时值",
            input: `BPM: 120\nC4 invalid\n`,
            expectError: true,
        },
        {
            name:  "缺少换行",
            input: `BPM: 120 C4 quarter`,
            expectError: true,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            input := strings.ReplaceAll(tc.input, "\\n", "\n")
            lexer := dsl.NewLexer(input)
            parser := dsl.NewParser(lexer)
            
            parser.ParseScore()
            hasErrors := len(parser.Errors()) > 0
            
            if tc.expectError && !hasErrors {
                t.Errorf("期望错误，但解析成功")
            } else if !tc.expectError && hasErrors {
                t.Errorf("期望成功，但解析错误: %v", parser.Errors())
            } else {
                t.Logf("✅ 错误处理正确")
                if hasErrors {
                    t.Logf("错误信息: %v", parser.Errors())
                }
            }
        })
    }
}