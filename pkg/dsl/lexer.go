package dsl


type Lexer struct {
    input        string
    position     int  // 当前字符位置
    readPosition int  // 下一个字符位置
    ch           byte // 当前字符
    line         int  // 当前行号
    column       int  // 当前列号
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		line:         1,
		column:       0,
	}
	l.readChar() // 初始化第一个字符
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    }
    return l.input[l.readPosition]
}


func (l *Lexer) NextToken() Token {
    var tok Token
    
    l.skipWhitespace()
    
    // 记录当前位置
    pos := Position{Line: l.line, Column: l.column}
    
    switch l.ch {
    case ':':
        tok = Token{Type: COLON, Literal: string(l.ch), Position: pos}
    case '\n':
        tok = Token{Type: NEWLINE, Literal: "\\n", Position: pos}
    case '\r':
        if l.peekChar() == '\n' {
            l.readChar() // 跳过\r
            tok = Token{Type: NEWLINE, Literal: "\\r\\n", Position: pos}
        } else {
            tok = Token{Type: NEWLINE, Literal: "\\r", Position: pos}
        }
    case 0:
        tok = Token{Type: EOF, Literal: "", Position: pos}
    default:
        if isLetter(l.ch) {
            literal := l.readIdentifier()
            tokenType := lookupIdentifierType(literal)
            return Token{Type: tokenType, Literal: literal, Position: pos}
        } else if isDigit(l.ch) {
            literal := l.readNumber()
            return Token{Type: NUMBER, Literal: literal, Position: pos}
        } else {
            tok = Token{Type: ILLEGAL, Literal: string(l.ch), Position: pos}
        }
    }
    
    l.readChar()
    return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t'  {
		l.readChar()
	}
}


// 字面量识别
func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

// 关键字和标识符映射
var keywords = map[string]TokenType{
    "BPM":     BPM,
    "C":       NOTE_C,
    "D":       NOTE_D,
    "E":       NOTE_E,
    "F":       NOTE_F,
    "G":       NOTE_G,
    "A":       NOTE_A,
    "B":       NOTE_B,
    "quarter": IDENTIFIER,
    "half":    IDENTIFIER,
    "whole":   IDENTIFIER,
    "eighth":  IDENTIFIER,
}

func lookupIdentifierType(identifier string) TokenType {
    if tok, ok := keywords[identifier]; ok {
        return tok
    }
    return IDENTIFIER
}