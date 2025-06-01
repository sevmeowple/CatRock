# CatRock DSL Music Player

A simple music programming language that converts text descriptions to MIDI audio playback.

[中文](README.md) | [English](README_EN.md)

## Quick Start

### Installation

```bash
git clone <repository-url>
cd catRock
go mod tidy
go build -o bin/catrock ./cmd/catrock
```

### Usage

```bash
# Show help
./bin/catrock

# Play music file
./bin/catrock play example/simple.crock
```

### Syntax Example

```crock
BPM: 120
C4 quarter
D4 quarter
E4 quarter
F4 quarter
G4 half
```

## Project Structure

```go
catRock/
├── cmd/catrock/          # CLI tool
├── pkg/
│   ├── dsl/             # DSL parser
│   ├── core/            # Music core
│   ├── score/           # Score system
│   └── io/              # MIDI interface
├── example/             # Example files
└── test/                # Tests
```

## Current Features

- Basic note playback
- MIDI device connection
- Simple CLI interface

## Work in Progress

- Chord support
- Multi-track functionality
- Section organization
- More syntax features

## Development

```bash
# Test
go test ./...
# ps: actually there are barely any tests (╯°□°）╯︵ ┻━┻

# Build
go build -o bin/catrock ./cmd/catrock
```

## License

[MIT License](LICENSE)
