# CatRock RoadMap

CatRock is a concise music description language. This document outlines the project's development roadmap and implementation priorities.

[English](roadmap_en.md) | [中文](roadmap.md)

## 📋 Version Planning

### 🚀 v0.2.0 - MVP (Current Version)

- ✅ Basic DSL syntax parsing
- ✅ MIDI event generation and playback
- ✅ Multi-track, multi-section support
- ✅ Basic command-line tools
- ✅ Simple example files

### 🎵 v0.3.0 - Stable Version

- 🎹 3-5 built-in sound presets
- 📚 Complete syntax and API documentation

### 🎨 v0.4.0 - Enhanced Version

- 📦 SoundFont (.sf2) sound library support
- 🔊 Native Audio output (bypass MIDI limitations)
- 🎵 Extended chord syntax
- 🎛️ Expression controls (velocity, pitch bend, etc.)
- 🎪 Real-time preview functionality

### 🌟 v1.0.0 - Official Release

- 🔌 VST3 plugin support
- 📤 Multi-format export (MIDI/WAV/MP3)
- 📖 Complete user manual

## 🎼 Development Phases

### Phase 1: Core Feature Enhancement 🔧

#### ✅ Completed

- Basic syntax parser
- MIDI event generation engine
- Command-line playback tool
- Basic example files

#### 🔄 In Progress

- **Error handling optimization** - Fix `<nil>` playback errors
- **Pitch sensitivity fixes** - Optimize sound performance across different octaves
- **Syntax documentation completion** - Provide complete DSL reference

#### 📝 To Be Implemented

- More user-friendly error messages
- Playback progress visualization
- Cross-platform MIDI device compatibility
- More example music files

---

### Phase 2: Sound System Extension 🎵

#### 🎯 RouteMap Sound Routing (v0.3.0)

```bash
# User experience goals
catrock play song.crock --route=default    # MIDI output
catrock play song.crock --route=enhanced   # Mixed output
catrock play song.crock --route=hq         # Custom sounds
```

**Core Features:**

- Configuration-driven sound mapping
- Namespaced sound syntax (`midi.89`, `builtin.piano`)
- Backward-compatible routing system
- User-defined routing configurations

#### 🎼 Built-in Sound Engine (v0.3.0)

**Goal:** Solve current MIDI sound quality issues

**Planned Implementation:**

- `builtin.warm_piano` - Full-range friendly piano
- `builtin.soft_flute` - Specifically addresses harshness issues
- `builtin.smooth_strings` - Warm string ensemble
- `builtin.ambient_pad` - Ambient background sounds
- `builtin.clean_guitar` - Clean guitar sound

**Technical Architecture:**

- Basic waveform synthesizer (sine, sawtooth, triangle waves)
- ADSR envelope control
- Real-time filter processing
- Native Audio output

#### 🔊 Hybrid Output Architecture (v0.4.0)

**Design Goal:** Same DSL file supports multiple sound engines

```bash
DSL Parser → Event Generator → Route Mapper → Audio Router
                                               ├─ MIDI Output
                                               ├─ Audio Synth  
                                               └─ SoundFont Player
```

---

### Phase 3: Advanced Features 🎨

#### 📦 External Sound Support (v0.3.0)

- **SoundFont (.sf2) Integration**
  - Open-source sound library support
  - Automatic recommended sound downloads
  - Sound library management tools
  
- **VST3 Plugin Support** (v1.0.0)
  - Third-party sound plugin loading
  - Professional-grade sound processing
  - Plugin parameter control

#### 📝 Syntax Extensions (v0.4.0 - v1.0.0)

```groovy
// Enhanced chord support
[Cmaj]/4 [Am]/4 [F]/4 [G]/4

// Expression controls
C4/4~120         // Specify velocity
C4/4^            // Accent mark
C4/4 bend(+200)  // Pitch bend control

// Template system
template verse {
    C4/4 D4/4 E4/4 F4/4
}

track melody {
    use verse transpose(+12)  // Transposed usage
}
```

#### 🛠️ Composition Tools (v0.4.0+)

- **Real-time preview** - Instant playback while editing
- **Syntax highlighting** - VSCode/Vim plugins
- **Music theory assistant** - Chord recommendations, scale generation
- **MIDI keyboard input** - Record and convert to DSL

---

### Phase 4: Ecosystem Building 🌍

#### 👥 Community Features (v1.0.0+)

- **Sharing platform** - Sound library and composition sharing
- **Collaborative editing** - Multi-user real-time composition
- **Version control** - Git-style music version management

#### 🔗 Integration Extensions (v1.0.0+)

- **Export functionality**
  - Standard MIDI file export
  - High-quality WAV/MP3 rendering
  - Sheet music PDF generation
  
- **DAW Integration**
  - Major DAW plugins
  - Project file import/export

**Last Updated**: June 3, 2025  
**Document Version**: v0.1.0

> All the ambitious plans in this document are written by AI (laughs)
> Welcome to participate in CatRock's development and discussion! Please submit issues or PRs on GitHub.
