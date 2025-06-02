package dsl

import (
	"catRock/pkg/dsl/ast"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// 读取两个token，初始化currentToken和peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	errorMsg := fmt.Sprintf("解析错误 %d:%d - %s",
		p.currentToken.Position.Line,
		p.currentToken.Position.Column,
		msg)
	p.errors = append(p.errors, errorMsg)
}

// 主解析入口
func (p *Parser) ParseScore() *ast.ScoreNode {
	score := &ast.ScoreNode{
		GlobalSets: []*ast.SetNode{},
		Elements:   []ast.PlayableNode{},
		Position:   p.currentToken.Position,
	}

	// 解析顶层元素
	for p.currentToken.Type != EOF {
		element := p.parseTopLevelElement()
		if element != nil {
			switch elem := element.(type) {
			case *ast.SetNode:
				score.GlobalSets = append(score.GlobalSets, elem)
			case ast.PlayableNode:
				score.Elements = append(score.Elements, elem)
			}
		}
	}

	return score
}

// 解析顶层元素
func (p *Parser) parseTopLevelElement() ast.ASTNode {
	switch p.currentToken.Type {
	case SET:
		return p.parseSetBlock(ast.GlobalContext)
	case TRACK:
		return p.parseTrack()
	case SECTION:
		return p.parseSection()
	case NOTE_C, NOTE_D, NOTE_E, NOTE_F, NOTE_G, NOTE_A, NOTE_B:
		return p.parseNote()
	case LBRACKET:
		return p.parseChord()
	case IDENTIFIER:
		// 可能是休止符或其他标识符
		if p.currentToken.Literal == "rest" {
			return p.parseRest()
		}
		p.addError(fmt.Sprintf("未知标识符: %s", p.currentToken.Literal))
		p.nextToken()
		return nil
	case NEWLINE:
		p.nextToken() // 跳过空行
		return nil
	default:
		p.addError(fmt.Sprintf("期望顶层元素，得到 %s", p.currentToken.Literal))
		p.nextToken() // 跳过错误token
		return nil
	}
}

// 解析Set块
func (p *Parser) parseSetBlock(context ast.ParameterContext) *ast.SetNode {
	position := p.currentToken.Position

	if !p.expectToken(SET) {
		return nil
	}

	if !p.expectToken(LBRACE) {
		return nil
	}

	parameters := make(map[string]interface{})

	// 解析参数列表
	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		if p.currentToken.Type == NEWLINE {
			p.nextToken()
			continue
		}

		// 解析参数名
		if p.currentToken.Type != IDENTIFIER {
			p.addError(fmt.Sprintf("期望参数名，得到 %s", p.currentToken.Literal))
			p.nextToken()
			continue
		}

		paramName := p.currentToken.Literal
		p.nextToken()

		if !p.expectToken(COLON) {
			continue
		}

		// 解析参数值
		paramValue := p.parseParameterValue()
		if paramValue != nil {
			parameters[paramName] = paramValue
		}
	}

	if !p.expectToken(RBRACE) {
		return nil
	}

	return &ast.SetNode{
		Parameters: parameters,
		Context:    context,
		Position:   position,
	}
}

// 解析参数值
func (p *Parser) parseParameterValue() interface{} {
	switch p.currentToken.Type {
	case NUMBER:
		if p.peekToken.Type == SLASH {
			return p.parseFraction()
		}
		value, err := strconv.Atoi(p.currentToken.Literal)
		if err != nil {
			p.addError(fmt.Sprintf("无效的数字: %s", p.currentToken.Literal))
			return nil
		}
		p.nextToken()
		return value

	case IDENTIFIER:
		value := p.currentToken.Literal
		p.nextToken()
		return value

	default:
		p.addError(fmt.Sprintf("期望参数值，得到 %s", p.currentToken.Literal))
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseFraction() float64 {
	numerator := p.currentToken.Literal
	p.nextToken()
	p.nextToken() // 跳过 SLASH

	if p.currentToken.Type != NUMBER {
		p.addError(fmt.Sprintf("期望分母数字，得到 %s", p.currentToken.Literal))
		return 0.0
	}

	denominator := p.currentToken.Literal
	p.nextToken()

	// 转换为浮点数
	num, err := strconv.ParseFloat(numerator, 64)
	if err != nil {
		p.addError(fmt.Sprintf("无效的分子: %s", numerator))
		return 0.0
	}
	
	den, err := strconv.ParseFloat(denominator, 64)
	if err != nil {
		p.addError(fmt.Sprintf("无效的分母: %s", denominator))
		return 0.0
	}
	
	if den == 0 {
		p.addError("分母不能为零")
		return 0.0
	}
	
	return num / den
}


// 解析Track
func (p *Parser) parseTrack() *ast.TrackNode {
	position := p.currentToken.Position

	if !p.expectToken(TRACK) {
		return nil
	}

	if p.currentToken.Type != IDENTIFIER {
		p.addError(fmt.Sprintf("期望轨道名称，得到 %s", p.currentToken.Literal))
		return nil
	}

	name := p.currentToken.Literal
	p.nextToken()

	if !p.expectToken(LBRACE) {
		return nil
	}

	track := &ast.TrackNode{
		Name:     name,
		Sets:     []*ast.SetNode{},
		Elements: []ast.PlayableNode{},
		Position: position,
	}

	// 解析Track内容
	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		element := p.parseContainerElement()
		if element != nil {
			switch elem := element.(type) {
			case *ast.SetNode:
				elem.Context = ast.TrackContext
				track.Sets = append(track.Sets, elem)
			case ast.PlayableNode:
				track.Elements = append(track.Elements, elem)
			}
		}
	}

	if !p.expectToken(RBRACE) {
		return nil
	}

	return track
}

// 解析Section
func (p *Parser) parseSection() *ast.SectionNode {
	position := p.currentToken.Position

	if !p.expectToken(SECTION) {
		return nil
	}

	if p.currentToken.Type != IDENTIFIER {
		p.addError(fmt.Sprintf("期望段落名称，得到 %s", p.currentToken.Literal))
		return nil
	}

	name := p.currentToken.Literal
	p.nextToken()

	if !p.expectToken(LBRACE) {
		return nil
	}

	section := &ast.SectionNode{
		Name:     name,
		Sets:     []*ast.SetNode{},
		Elements: []ast.PlayableNode{},
		Position: position,
	}

	// 解析Section内容
	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		element := p.parseContainerElement()
		if element != nil {
			switch elem := element.(type) {
			case *ast.SetNode:
				elem.Context = ast.SectionContext
				section.Sets = append(section.Sets, elem)
			case ast.PlayableNode:
				section.Elements = append(section.Elements, elem)
			}
		}
	}

	if !p.expectToken(RBRACE) {
		return nil
	}

	return section
}


// 完全重写parseNote方法
func (p *Parser) parseNote() *ast.NoteNode {
    position := p.currentToken.Position

    if !p.isNoteToken(p.currentToken.Type) {
        p.addError(fmt.Sprintf("期望音符名称，得到 %s", p.currentToken.Literal))
        return nil
    }

    noteName := p.currentToken.Literal
    p.nextToken()

    if p.currentToken.Type != NUMBER {
        p.addError(fmt.Sprintf("期望八度数字，得到 %s", p.currentToken.Literal))
        return nil
    }

    octave, err := strconv.Atoi(p.currentToken.Literal)
    if err != nil || octave < 0 || octave > 9 {
        p.addError(fmt.Sprintf("无效的八度值: %s", p.currentToken.Literal))
        return nil
    }

    p.nextToken()

    // 解析时值 - 支持 /分数表示法
    duration := p.parseNoteDuration()

    return &ast.NoteNode{
        Name:     noteName,
        Octave:   octave,
        Duration: duration,
        Position: position,
    }
}

// 新的音符时值解析方法
func (p *Parser) parseNoteDuration() string {
    // 默认时值
    defaultDuration := "1/4" // 四分音符

    // 检查是否有时值修饰符
    if p.currentToken.Type == SLASH {
        p.nextToken() // 跳过斜杠
        
        if p.currentToken.Type == NUMBER {
            denominator := p.currentToken.Literal
            p.nextToken()
            
            // 检查附点
            dotted := ""
            if p.currentToken.Type == DOT {
                dotted = "."
                p.nextToken()
            }
            
            // 返回分数形式的时值
            return fmt.Sprintf("1/%s%s", denominator, dotted)
        } else {
            p.addError("期望时值分母")
            return defaultDuration
        }
    } else if p.currentToken.Type == IDENTIFIER && p.isValidDuration(p.currentToken.Literal) {
        // 传统时值表示 如 quarter, half（保持向后兼容）
        duration := p.currentToken.Literal
        p.nextToken()
        
        // 检查附点
        if p.currentToken.Type == DOT {
            duration = duration + "."
            p.nextToken()
        }
        
        return duration
    }

    // 没有指定时值，使用默认
    return defaultDuration
}

// 完全重写parseChord方法
func (p *Parser) parseChord() *ast.ChordNode {
    position := p.currentToken.Position

    if !p.expectToken(LBRACKET) {
        return nil
    }

    var content interface{}

    // 检查是和弦名还是音符列表
    if p.currentToken.Type == IDENTIFIER {
        // 和弦名 如 Am, Cmaj
        chordName := p.currentToken.Literal
        content = chordName
        p.nextToken()
    } else {
        // 音符列表 如 [C4 E4 G4]
        notes := []*ast.NoteNode{}
        for p.currentToken.Type != RBRACKET && p.currentToken.Type != EOF {
            if p.isNoteToken(p.currentToken.Type) {
                note := p.parseNote()
                if note != nil {
                    notes = append(notes, note)
                }
            } else {
                p.addError(fmt.Sprintf("期望音符，得到 %s", p.currentToken.Literal))
                p.nextToken()
            }
        }
        content = notes
    }

    if !p.expectToken(RBRACKET) {
        return nil
    }

    // 解析和弦时值
    duration := p.parseNoteDuration()

    return &ast.ChordNode{
        Content:  content,
        Duration: duration,
        Position: position,
    }
}

// 修改parseRest方法
func (p *Parser) parseRest() *ast.RestNode {
    position := p.currentToken.Position

    if p.currentToken.Literal != "rest" {
        p.addError(fmt.Sprintf("期望 'rest'，得到 %s", p.currentToken.Literal))
        return nil
    }

    p.nextToken()

    // 解析休止符时值
    duration := p.parseNoteDuration()

    return &ast.RestNode{
        Duration: duration,
        Position: position,
    }
}

// 新增：解析分组音符（处理括号）
func (p *Parser) parseGroup() *ast.GroupNode {
    position := p.currentToken.Position

    if !p.expectToken(LPAREN) {
        return nil
    }

    group := &ast.GroupNode{
        Elements: []ast.ElementNode{},
        Position: position,
    }

    // 解析组内元素
    for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
        element := p.parsePlayableElement()
        if element != nil {
            group.Elements = append(group.Elements, element)
        }
    }

    if !p.expectToken(RPAREN) {
        return nil
    }

    // 组可以有整体时值修饰符
    group.Duration = p.parseNoteDuration()

    return group
}


// 新增：解析可播放元素的通用方法
func (p *Parser) parsePlayableElement() ast.PlayableNode {
    switch p.currentToken.Type {
    case NOTE_C, NOTE_D, NOTE_E, NOTE_F, NOTE_G, NOTE_A, NOTE_B:
        return p.parseNote()
    case LBRACKET:
        return p.parseChord()
    case LPAREN:
        return p.parseGroup()
    case IDENTIFIER:
        if p.currentToken.Literal == "rest" {
            return p.parseRest()
        }
        p.addError(fmt.Sprintf("未知标识符: %s", p.currentToken.Literal))
        p.nextToken()
        return nil
    case NEWLINE:
        p.nextToken() // 跳过空行
        return nil
    default:
        p.addError(fmt.Sprintf("期望可播放元素，得到 %s", p.currentToken.Literal))
        p.nextToken()
        return nil
    }
}

// 更新parseContainerElement方法，支持分组
func (p *Parser) parseContainerElement() ast.ASTNode {
    switch p.currentToken.Type {
    case SET:
        return p.parseSetBlock(ast.SectionContext)
    case SECTION:
        return p.parseSection()
    case TRACK:
        return p.parseTrack()
    case NOTE_C, NOTE_D, NOTE_E, NOTE_F, NOTE_G, NOTE_A, NOTE_B:
        return p.parseNote()
    case LBRACKET:
        return p.parseChord()
    case LPAREN: // 新增：支持分组
        return p.parseGroup()
    case IDENTIFIER:
        if p.currentToken.Literal == "rest" {
            return p.parseRest()
        }
        p.addError(fmt.Sprintf("未知标识符: %s", p.currentToken.Literal))
        p.nextToken()
        return nil
    case NEWLINE:
        p.nextToken() // 跳过空行
        return nil
    default:
        p.addError(fmt.Sprintf("期望容器元素，得到 %s", p.currentToken.Literal))
        p.nextToken()
        return nil
    }
}

// 辅助方法
func (p *Parser) expectToken(expectedType TokenType) bool {
	if p.currentToken.Type == expectedType {
		p.nextToken()
		return true
	}

	p.addError(fmt.Sprintf("期望 %s，得到 %s",
		tokenTypeToString(expectedType), tokenTypeToString(p.currentToken.Type)))
	return false
}

func (p *Parser) isNoteToken(tokenType TokenType) bool {
	return tokenType == NOTE_C || tokenType == NOTE_D ||
		tokenType == NOTE_E || tokenType == NOTE_F ||
		tokenType == NOTE_G || tokenType == NOTE_A ||
		tokenType == NOTE_B
}

// 更新时值验证方法
func (p *Parser) isValidDuration(duration string) bool {
    // 传统时值名称（保持向后兼容）
    validDurations := []string{"whole", "half", "quarter", "eighth", "sixteenth"}
    for _, valid := range validDurations {
        if duration == valid {
            return true
        }
    }
    
    // 分数形式时值验证（可选）
    if strings.HasPrefix(duration, "1/") {
        return true
    }
    
    return false
}

// Token类型转字符串（用于错误信息）
func tokenTypeToString(tokenType TokenType) string {
	switch tokenType {
	case SET:
		return "SET"
	case TRACK:
		return "TRACK"
	case SECTION:
		return "SECTION"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case COLON:
		return ":"
	case SLASH:
		return "/"
	case DOT:
		return "."
	case NUMBER:
		return "NUMBER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NEWLINE:
		return "NEWLINE"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("Unknown(%d)", int(tokenType))
	}
}
