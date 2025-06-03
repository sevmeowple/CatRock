# CatRock Music Player

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
