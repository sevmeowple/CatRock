package commands

import (
	"catRock/pkg/dsl"
	"catRock/pkg/dsl/ast"
	"catRock/pkg/score"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type debugOpts struct {
	ShowTokens bool
	ShowAST    bool
	ShowEvents bool
	ShowScore  bool
}

func newDebugCmd() *cobra.Command {
    var opts debugOpts

    debugCmd := &cobra.Command{
        Use:   "debug <file.crock>",
        Short: "🔧 调试 CatRock 文件解析过程",
        Long: `
显示 CatRock 文件的详细解析信息：
- 词法分析结果 (tokens)
- 抽象语法树 (AST)
- 生成的事件序列
- Score 对象信息
`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return runDebug(args[0], &opts)
        },
    }
    
    debugCmd.Flags().BoolVar(&opts.ShowTokens, "tokens", false, "显示词法分析结果")
    debugCmd.Flags().BoolVar(&opts.ShowAST, "ast", true, "显示抽象语法树")
    debugCmd.Flags().BoolVar(&opts.ShowEvents, "events", true, "显示生成的事件")
    debugCmd.Flags().BoolVar(&opts.ShowScore, "score", false, "显示Score对象详情")
    
    return debugCmd
}

func runDebug(filename string, opts *debugOpts) error {
    green := color.New(color.FgGreen, color.Bold)
    yellow := color.New(color.FgYellow)
    
    // 读取文件
    content, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    green.Printf("🔧 调试文件: %s\n", filename)
    
    // 1. 词法分析
    if opts.ShowTokens {
        yellow.Println("\n🔤 词法分析结果:")
        lexer := dsl.NewLexer(string(content))
        for {
            token := lexer.NextToken()
            if token.Type == dsl.EOF {
                break
            }
            fmt.Printf("   %s: '%s' (行:%d 列:%d)\n", 
                token.Type, token.Literal, token.Position.Line, token.Position.Column)
        }
    }
    
    // 2. 语法分析
    lexer := dsl.NewLexer(string(content))
    parser := dsl.NewParser(lexer)
    ast := parser.ParseScore()
    
    if len(parser.Errors()) > 0 {
        color.Red("❌ 解析错误:")
        for _, err := range parser.Errors() {
            fmt.Printf("   %s\n", err)
        }
        return fmt.Errorf("解析失败")
    }
    
    if opts.ShowAST {
        yellow.Println("\n🌳 抽象语法树:")
        printDetailedAST(ast)
    }
    
    // 3. 代码生成
    generator := dsl.NewGenerator()
    scoreObj, err := generator.GenerateScore(ast)
    if err != nil {
        return err
    }
    
     if opts.ShowScore {
        yellow.Println("\n🎼 生成的Score对象:")
        fmt.Print(scoreObj.DetailedString("   "))
    }
    
    // 4. 事件生成
    if opts.ShowEvents {
        engine := score.NewPlayEngine(scoreObj)
        events, err := engine.GenerateEvents()
        if err != nil {
            return err
        }
        
        yellow.Println("\n🎹 生成的事件:")
        showDetailedEvents(events)
    }
    
    return nil
}

func printDetailedAST(scoreNode *ast.ScoreNode) {
    if scoreNode == nil {
        fmt.Println("   (空的AST)")
        return
    }
    
    // 直接使用详细字符串方法
    fmt.Print(scoreNode.DetailedString("   "))
}
func showDetailedEvents(events []score.Event) {
    if len(events) == 0 {
        fmt.Println("   (没有生成事件)")
        return
    }

    // 添加颜色
    greenColor := color.New(color.FgGreen)
    redColor := color.New(color.FgRed)
    yellowColor := color.New(color.FgYellow, color.Bold)
    cyanColor := color.New(color.FgCyan)

    for _ , event := range events {
        // 根据事件动作选择颜色
        var eventColor *color.Color
        var actionName string
        
        switch event.Action {
        case score.NOTE_ON:
            eventColor = greenColor
            actionName = "NOTE_ON"
        case score.NOTE_OFF:
            eventColor = redColor
            actionName = "NOTE_OFF"
        case score.PROGRAM_CHANGE:
            eventColor = yellowColor
            actionName = "PROGRAM_CHANGE"
        case score.VOLUME_CHANGE:
            eventColor = cyanColor
            actionName = "VOLUME_CHANGE"
        default:
            eventColor = color.New(color.FgWhite)
            actionName = fmt.Sprintf("UNKNOWN_%d", event.Action)
        }

        // 显示基本事件信息
        eventColor.Printf("   [%6.3fs] %s", event.Time, actionName)
        fmt.Printf(" Ch:%d", event.Channel)
        
        // 根据事件类型显示详细数据
        switch event.Action {
        case score.NOTE_ON, score.NOTE_OFF:
            if event.Velocity > 0 {
                fmt.Printf(" Vel:%d", event.Velocity)
            }
            // 显示音符数据
            fmt.Printf(" Data:%v", event.Data)
            
        case score.PROGRAM_CHANGE:
            fmt.Printf(" Instrument:%v", event.Data)
            
        case score.VOLUME_CHANGE:
            fmt.Printf(" Volume:%v", event.Data)
        }
        
        if event.Duration > 0 {
            fmt.Printf(" Dur:%.3f", event.Duration)
        }
        
        if event.SourceElement != "" {
            fmt.Printf(" Src:%s", event.SourceElement)
        }
        
        fmt.Println()

    }
    
    // 显示统计信息
    showEventStatistics(events)
}

// 显示事件统计信息
func showEventStatistics(events []score.Event) {
    stats := make(map[score.EventAction]int)
    channels := make(map[int]int)
    
    for _, event := range events {
        stats[event.Action]++
        channels[event.Channel]++
    }
    
    fmt.Println("\n   📊 事件统计:")
    for action, count := range stats {
        var actionName string
        switch action {
        case score.NOTE_ON:
            actionName = "音符开始"
        case score.NOTE_OFF:
            actionName = "音符结束"
        case score.PROGRAM_CHANGE:
            actionName = "乐器切换"
        case score.VOLUME_CHANGE:
            actionName = "音量变化"
        default:
            actionName = fmt.Sprintf("未知(%d)", action)
        }
        fmt.Printf("      %s: %d个\n", actionName, count)
    }
    
    fmt.Println("   📡 使用的通道:")
    for channel, count := range channels {
        fmt.Printf("      通道%d: %d个事件\n", channel, count)
    }
}