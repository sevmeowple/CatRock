package dsl

import (
    "fmt"
    "strconv"
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
func (p *Parser) ParseScore() *ScoreNode {
    score := &ScoreNode{}
    
    // 解析header
    header := p.parseHeader()
    if header == nil {
        return nil
    }
    score.Header = header
    
    // 解析body
    score.Body = p.parseBody()
    
    return score
}

func (p *Parser) parseHeader() *HeaderNode {
    // 期望: BPM : NUMBER NEWLINE
    if !p.expectToken(BPM) {
        return nil
    }
    
    if !p.expectToken(COLON) {
        return nil
    }
    
    if p.currentToken.Type != NUMBER {
        p.addError(fmt.Sprintf("期望数字，得到 %s", p.currentToken.Literal))
        return nil
    }
    
    bpm, err := strconv.Atoi(p.currentToken.Literal)
    if err != nil {
        p.addError(fmt.Sprintf("无效的BPM值: %s", p.currentToken.Literal))
        return nil
    }
    
    p.nextToken() // 消费NUMBER
    
    if !p.expectToken(NEWLINE) {
        return nil
    }
    
    return &HeaderNode{BPM: bpm}
}

func (p *Parser) parseBody() []ElementNode {
    elements := []ElementNode{}
    
    for p.currentToken.Type != EOF {
        element := p.parseElement()
        if element != nil {
            elements = append(elements, element)
        }
    }
    
    return elements
}

func (p *Parser) parseElement() ElementNode {
    switch p.currentToken.Type {
    case NOTE_C, NOTE_D, NOTE_E, NOTE_F, NOTE_G, NOTE_A, NOTE_B:
        return p.parseNote()
    case NEWLINE:
        p.nextToken() // 跳过空行
        return nil
    default:
        p.addError(fmt.Sprintf("期望音符，得到 %s", p.currentToken.Literal))
        p.nextToken() // 跳过错误token
        return nil
    }
}

func (p *Parser) parseNote() *NoteNode {
    // 期望: NOTE_NAME NUMBER IDENTIFIER NEWLINE
    
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
    
    if p.currentToken.Type != IDENTIFIER {
        p.addError(fmt.Sprintf("期望时值，得到 %s", p.currentToken.Literal))
        return nil
    }
    
    duration := p.currentToken.Literal
    if !p.isValidDuration(duration) {
        p.addError(fmt.Sprintf("无效的时值: %s", duration))
        return nil
    }
    
    p.nextToken()
    
    if !p.expectToken(NEWLINE) {
        return nil
    }
    
    return &NoteNode{
        Name:     noteName,
        Octave:   octave,
        Duration: duration,
    }
}

func (p *Parser) expectToken(expectedType TokenType) bool {
    if p.currentToken.Type == expectedType {
        p.nextToken()
        return true
    }
    
    p.addError(fmt.Sprintf("期望 %s，得到 %s", 
        expectedType, p.currentToken.Type))
    return false
}

func (p *Parser) isNoteToken(tokenType TokenType) bool {
    return tokenType == NOTE_C || tokenType == NOTE_D || 
           tokenType == NOTE_E || tokenType == NOTE_F ||
           tokenType == NOTE_G || tokenType == NOTE_A || 
           tokenType == NOTE_B
}

func (p *Parser) isValidDuration(duration string) bool {
    validDurations := []string{"whole", "half", "quarter", "eighth"}
    for _, valid := range validDurations {
        if duration == valid {
            return true
        }
    }
    return false
}