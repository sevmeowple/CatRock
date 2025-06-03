# CatRock 音乐播放器

一个简单的音乐编程语言
使用 MIDI 作为后端播放

[中文](README.md) | [English](README_EN.md)

## 快速开始

### 安装

```bash
git clone <repository-url>
cd catRock
go mod tidy
go build -o bin/catrock ./cmd/catrock
```

### 使用

```bash
# 显示帮助
./bin/catrock

# 播放音乐文件
./bin/catrock play example/simple.crock
```

### 语法示例

```crock
set {
    BPM: 120
    base_duration: 1/4      
}

track melody {
    section intro {
        C4 D4 E4 F4
        
        G4/2 A4/1               
        C4/8 D4/8 E4/8 F4/8     
  
        G4/2. A4/4.            
       
        (C4/8 D4/8 E4/8)       
    }
    
    section verse {

        C4 D4/8 E4/8 F4/2
        [C4 E4 G4]/1        
    }
}
```

## 项目结构

```go
catRock/
├── cmd/catrock/          # CLI工具
├── pkg/
│   ├── dsl/             # DSL解析器
│   ├── core/            # 音乐核心
│   ├── score/           # 乐谱系统
│   └── io/              # MIDI接口
├── example/             # 示例文件
└── test/                # 测试
```

## 当前功能

- 基本音符播放
- MIDI 设备连接
- 简单的 CLI 界面

## 开发中 - WIP

- 和弦支持
- 多轨道功能
- Section 组织结构
- 更多语法特性

## 开发

```bash
# 测试
go test ./...
# ps: 其实根本没啥测试(╯°□°）╯︵ ┻━┻

# 构建
go build -o bin/catrock ./cmd/catrock
```

## 许可证

[MIT License](LICENSE)
