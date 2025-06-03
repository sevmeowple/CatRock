# CatRock RoadMap 项目路线图

CatRock是一个简洁的音乐描述语言，本文档描述了项目的发展规划和实现优先级。

[English](roadmap_en.md) | [中文](roadmap.md)

## 📋 版本规划

### 🚀 v0.2.0 - MVP (当前版本)

- ✅ 基础DSL语法解析
- ✅ MIDI事件生成和播放
- ✅ 多音轨、多段落支持
- ✅ 基础命令行工具
- ✅ 简单示例文件

### 🎵 v0.3.0 - 稳定版

- 🎹 3-5个内置音色
- 📚 完整语法和API文档

### 🎨 v0.4.0 - 增强版

- 📦 SoundFont (.sf2) 音色库支持
- 🔊 原生Audio输出 (绕过MIDI限制)
- 🎵 和弦语法扩展
- 🎛️ 表达控制 (力度、弯音等)
- 🎪 实时预览功能

### 🌟 v1.0.0 - 正式版

- 🔌 VST3插件支持
- 📤 多格式导出 (MIDI/WAV/MP3)
- 📖 完整用户手册

## 🎼 功能开发阶段

### Phase 1: 核心功能完善 🔧

#### ✅ 已完成

- 基础语法解析器
- MIDI事件生成引擎
- 命令行播放工具
- 基础示例文件

#### 🔄 进行中

- **错误处理优化** - 修复`<nil>`播放错误
- **音高敏感性修复** - 优化不同音区的音色表现
- **语法文档完善** - 提供完整的DSL参考

#### 📝 待实现

- 更友好的错误信息提示
- 播放进度可视化
- 跨平台MIDI设备兼容性
- 更多示例音乐文件

---

### Phase 2: 音色系统扩展 🎵

#### 🎯 RouteMap音色路由 (v0.3.0)

```bash
# 用户体验目标
catrock play song.crock --route=default    # MIDI输出
catrock play song.crock --route=enhanced   # 混合输出
catrock play song.crock --route=hq         # 自定义音色
```

**核心功能：**

- 配置驱动的音色映射
- 命名空间音色语法 (`midi.89`, `builtin.piano`)
- 向后兼容的路由系统
- 用户自定义路由配置

#### 🎼 内置音色引擎 (v0.3.0)

**目标：** 解决当前MIDI音色的质量问题

**计划实现：**

- `builtin.warm_piano` - 全音区友好的钢琴
- `builtin.soft_flute` - 专门解决刺耳问题
- `builtin.smooth_strings` - 温暖弦乐合奏
- `builtin.ambient_pad` - 环境铺底音色
- `builtin.clean_guitar` - 清音吉他

**技术架构：**

- 基础波形合成器 (正弦、锯齿、三角波)
- ADSR包络控制
- 实时滤波器处理
- 原生Audio输出

#### 🔊 混合输出架构 (v0.4.0)

**设计目标：** 同一DSL文件支持多种音色引擎

```bash
DSL Parser → Event Generator → Route Mapper → Audio Router
                                               ├─ MIDI Output
                                               ├─ Audio Synth  
                                               └─ SoundFont Player
```

---

### Phase 3: 高级功能 🎨

#### 📦 外部音色支持 (v0.3.0)

- **SoundFont (.sf2) 集成**
  - 开源音色库支持
  - 自动下载推荐音色
  - 音色库管理工具
  
- **VST3插件支持** (v1.0.0)
  - 第三方音色插件加载
  - 专业级音色处理
  - 插件参数控制

#### 📝 语法扩展 (v0.4.0 - v1.0.0)

```groovy
// 增强和弦支持
[Cmaj]/4 [Am]/4 [F]/4 [G]/4

// 表达控制
C4/4~120         // 指定力度
C4/4^            // 重音标记
C4/4 bend(+200)  // 弯音控制

// 模板系统
template verse {
    C4/4 D4/4 E4/4 F4/4
}

track melody {
    use verse transpose(+12)  // 移调使用
}
```

#### 🛠️ 创作工具 (v0.4.0+)

- **实时预览** - 编辑时即时播放
- **语法高亮** - VSCode/Vim插件
- **音乐理论助手** - 和弦推荐、音阶生成
- **MIDI键盘输入** - 录制转换为DSL

---

### Phase 4: 生态建设 🌍

#### 👥 社区功能 (v1.0.0+)

- **分享平台** - 音色库和作品分享
- **协作编辑** - 多人实时创作
- **版本控制** - Git风格的音乐版本管理

#### 🔗 集成扩展 (v1.0.0+)

- **导出功能**
  - 标准MIDI文件导出
  - 高质量WAV/MP3渲染
  - 乐谱PDF生成
  
- **DAW集成**
  - 主流DAW插件
  - 工程文件导入导出

**上次更新**: 2025年6月3日  
**文档版本**: v0.1.0

> 本文档所有的大饼都是AI写的(笑)
> 欢迎参与CatRock的开发和讨论！请在GitHub上提交issue或PR。
