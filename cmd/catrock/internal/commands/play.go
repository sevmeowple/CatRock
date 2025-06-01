package commands

import (
    "catRock/pkg/dsl"
    "catRock/pkg/io"
    "catRock/pkg/io/midi"
    "catRock/pkg/score"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/fatih/color"
    "github.com/schollz/progressbar/v3"
    "github.com/spf13/cobra"
)

type PlayOptions struct {
    Tempo    float64
    Volume   int
    DryRun   bool
    ShowAST  bool
    ShowEvents bool
}

func newPlayCmd() *cobra.Command {
    var opts PlayOptions
    
    playCmd := &cobra.Command{
        Use:   "play [file.crock]",
        Short: "🎵 播放CatRock音乐文件", 
        Long: `播放指定的.crock音乐文件。

文件必须是.crock格式，包含有效的CatRock DSL语法。
播放时会自动连接系统MIDI设备进行音频输出。`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return runPlay(args[0], &opts)
        },
    }

    // 播放选项
    playCmd.Flags().Float64Var(&opts.Tempo, "tempo", 0, "覆盖文件中的BPM设置")
    playCmd.Flags().IntVar(&opts.Volume, "volume", 100, "播放音量 (0-127)")
    playCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "只解析验证，不实际播放")
    
    // 调试选项
    playCmd.Flags().BoolVar(&opts.ShowAST, "show-ast", false, "显示抽象语法树")
    playCmd.Flags().BoolVar(&opts.ShowEvents, "show-events", false, "显示生成的MIDI事件")

    return playCmd
}

func runPlay(filename string, opts *PlayOptions) error {
    // 颜色定义
    green := color.New(color.FgGreen, color.Bold)
    red := color.New(color.FgRed, color.Bold)
    yellow := color.New(color.FgYellow)
    cyan := color.New(color.FgCyan)
    white := color.New(color.FgWhite)

    // 1. 验证文件
    if err := validateFile(filename); err != nil {
        red.Printf("❌ 文件错误: %v\n", err)
        return err
    }

    green.Printf("🎵 开始处理: %s\n", filepath.Base(filename))

    // 2. 读取文件
    content, err := os.ReadFile(filename)
    if err != nil {
        red.Printf("❌ 读取失败: %v\n", err)
        return err
    }

    if verbose {
        cyan.Printf("📄 文件大小: %d 字节\n", len(content))
    }

    // 3. 解析过程
    yellow.Println("🔍 正在解析...")
    
    // 词法分析
    lexer := dsl.NewLexer(string(content))
    parser := dsl.NewParser(lexer)
    ast := parser.ParseScore()

    if len(parser.Errors()) > 0 {
        red.Println("❌ 解析错误:")
        for _, err := range parser.Errors() {
            fmt.Printf("   %s\n", err)
        }
        return fmt.Errorf("解析失败")
    }

    green.Println("✅ 解析成功")

    if opts.ShowAST {
        white.Println("\n🌳 抽象语法树:")
        fmt.Printf("   %s\n", ast)
    }

    // 4. 代码生成
    yellow.Println("⚙️  正在生成音乐...")
    
    generator := dsl.NewGenerator()
    scoreObj, err := generator.GenerateScore(ast)
    if err != nil {
        red.Printf("❌ 生成失败: %v\n", err)
        return err
    }

    // 应用选项
    if opts.Tempo > 0 {
        scoreObj.BPM = opts.Tempo
        cyan.Printf("🎼 BPM设置为: %.0f\n", opts.Tempo)
    }

    if opts.Volume != 100 {
        scoreObj.Volume = opts.Volume
        cyan.Printf("🔊 音量设置为: %d\n", opts.Volume)
    }

    // 5. 生成事件
    events, err := scoreObj.Play()
    if err != nil {
        red.Printf("❌ 事件生成失败: %v\n", err)
        return err
    }

    green.Println("✅ 音乐生成完成")

    // 显示音乐信息
    white.Println("\n📊 音乐信息:")
    fmt.Printf("   🎼 BPM: %.0f\n", scoreObj.BPM)
    fmt.Printf("   ⏱️  时长: %.2f拍 (约%.1f秒)\n", 
        scoreObj.Duration(), 
        scoreObj.Duration()*60/scoreObj.BPM)
    fmt.Printf("   🎹 事件数: %d\n", len(events))

    if opts.ShowEvents {
        white.Println("\n🎹 MIDI事件:")
        showEvents(events)
    }

    // 6. 播放
    if opts.DryRun {
        yellow.Println("\n🚫 Dry-run模式，跳过播放")
        return nil
    }

    return playMusic(scoreObj, events)
}

func validateFile(filename string) error {
    if !strings.HasSuffix(filename, ".crock") {
        return fmt.Errorf("文件必须是.crock格式")
    }

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return fmt.Errorf("文件不存在: %s", filename)
    }

    return nil
}

func showEvents(events []score.Event) {
    for i, event := range events {
        if i >= 10 { // 只显示前10个事件
            fmt.Printf("   ... 还有 %d 个事件\n", len(events)-10)
            break
        }
        fmt.Printf("   [%.2f] %s Ch%d Data%v\n", 
            event.Time, event.Action, event.Channel, event.Data)
    }
}

func playMusic(scoreObj *score.Score, events []score.Event) error {
    green := color.New(color.FgGreen, color.Bold)
    red := color.New(color.FgRed, color.Bold)
    yellow := color.New(color.FgYellow)

    // 连接MIDI
    yellow.Println("\n🎹 正在连接MIDI设备...")
    
    midiPlayer := midi.NewMIDIPlayer()
    status, err := midiPlayer.Connect()
    if err != nil {
        red.Printf("❌ MIDI连接失败: %v\n", err)
        return err
    }

    if status != io.Connected {
        red.Println("❌ MIDI设备未连接")
        return fmt.Errorf("MIDI设备不可用")
    }

    defer midiPlayer.Disconnect()
    green.Println("✅ MIDI设备连接成功")

    // 播放进度条
    playDuration := scoreObj.Duration() * 60 / scoreObj.BPM // 转换为秒
    bar := progressbar.NewOptions(int(playDuration*10), // 精度0.1秒
        progressbar.OptionSetDescription("🎵 播放中"),
        progressbar.OptionSetTheme(progressbar.Theme{
            Saucer:        "█",
            SaucerHead:    "█",
            SaucerPadding: "░",
            BarStart:      "╢",
            BarEnd:        "╟",
        }),
        progressbar.OptionShowCount(),
        progressbar.OptionShowElapsedTimeOnFinish(),
        progressbar.OptionSetWidth(50),
    )

    yellow.Printf("\n🎵 开始播放... (按Ctrl+C停止)\n\n")

    // 创建播放引擎
    playEngine := score.NewPlayEngine(midiPlayer, scoreObj.BPM)

    // 启动进度条更新
    done := make(chan bool)
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        start := time.Now()
        for {
            select {
            case <-ticker.C:
                elapsed := time.Since(start).Seconds()
                if elapsed >= playDuration {
                    bar.Finish()
                    return
                }
                bar.Set(int(elapsed * 10))
            case <-done:
                bar.Finish()
                return
            }
        }
    }()

    // 执行播放
    start := time.Now()
    err = playEngine.PlayEvents(events)
    duration := time.Since(start)
    
    done <- true

    if err != nil {
        red.Printf("\n❌ 播放失败: %v\n", err)
        return err
    }

    green.Printf("\n✅ 播放完成! (用时: %v)\n", duration)
    return nil
}