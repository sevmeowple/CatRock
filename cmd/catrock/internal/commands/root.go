package commands

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

var (
    verbose bool
    rootCmd *cobra.Command
)

func Execute(version string) error {
    rootCmd = &cobra.Command{
        Use:     "catrock",
        Short:   "🎵 CatRock DSL音乐播放器",
        Long:    getBanner(version),
        Version: version,
        Run:     showInfo,
    }

    // 全局标志
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

    // 添加子命令
    rootCmd.AddCommand(newPlayCmd())
    rootCmd.AddCommand(newDebugCmd())
    return rootCmd.Execute()
}

func showInfo(cmd *cobra.Command, args []string) {
    printBanner()
    printUsage()
}

func getBanner(version string) string {
    return fmt.Sprintf(`
🎵 CatRock DSL音乐播放器 v%s

一个简单而强大的音乐编程语言播放器
支持实时解析.crock文件并通过MIDI播放音乐
`, version)
}

func printBanner() {
    cyan := color.New(color.FgCyan, color.Bold)
    
    cyan.Println(`
 ██████╗ █████╗ ████████╗██████╗  ██████╗  ██████╗██╗  ██╗
██╔════╝██╔══██╗╚══██╔══╝██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝
██║     ███████║   ██║   ██████╔╝██║   ██║██║     █████╔╝ 
██║     ██╔══██║   ██║   ██╔══██╗██║   ██║██║     ██╔═██╗ 
╚██████╗██║  ██║   ██║   ██║  ██║╚██████╔╝╚██████╗██║  ██╗
 ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝`)
    
}

func printUsage() {
    green := color.New(color.FgGreen, color.Bold)
    white := color.New(color.FgWhite)
    blue := color.New(color.FgBlue)
    
    green.Println("\n📖 使用方法:")
    
    white.Println("\n基础命令:")
    blue.Println("  catrock                    # 显示此帮助信息")
    blue.Println("  catrock play <file.crock>  # 播放音乐文件")
    blue.Println("  catrock --version          # 显示版本信息")
    blue.Println("  catrock --help             # 显示详细帮助")
    
    white.Println("\n播放选项:")
    blue.Println("  catrock play song.crock -v        # 详细模式播放")
    blue.Println("  catrock play song.crock --tempo 140  # 自定义BPM")
    blue.Println("  catrock play song.crock --dry-run    # 只解析不播放")
    
    white.Println("\n示例文件格式 (.crock):")
    fmt.Println("  BPM: 120")
    fmt.Println("  C4 quarter")
    fmt.Println("  D4 quarter") 
    fmt.Println("  E4 quarter")
    fmt.Println("  F4 half")
    
    white.Println("\n📁 示例文件:")
    blue.Println("  example/simple.crock       # 简单旋律示例")
    blue.Println("  example/demo.crock         # 演示文件")
    
    color.New(color.FgMagenta).Println("\n✨ 开始你的音乐编程之旅吧！")
}