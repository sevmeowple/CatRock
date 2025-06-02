package dsl

import "catRock/pkg/dsl/mytype"

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
	pos := mytype.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch), Position: pos}
	case '\n':
		tok = Token{Type: NEWLINE, Literal: "\\n", Position: pos}
	case '/':
		// 检查是否是注释
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken() // 递归获取下一个有效token
		} else {
			// 普通斜杠（用于时值）
			tok = Token{Type: SLASH, Literal: string(l.ch), Position: pos}
		}
	case '.': // 新增 - 附点音符
		tok = Token{Type: DOT, Literal: string(l.ch), Position: pos}
	case '{': // 新增 - 左花括号
		tok = Token{Type: LBRACE, Literal: string(l.ch), Position: pos}
	case '}': // 新增 - 右花括号
		tok = Token{Type: RBRACE, Literal: string(l.ch), Position: pos}
	case '[': // 新增 - 左方括号
		tok = Token{Type: LBRACKET, Literal: string(l.ch), Position: pos}
	case ']': // 新增 - 右方括号
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Position: pos}
	case '(': // 新增 - 左圆括号
		tok = Token{Type: LPAREN, Literal: string(l.ch), Position: pos}
	case ')': // 新增 - 右圆括号
		tok = Token{Type: RPAREN, Literal: string(l.ch), Position: pos}
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
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	l.readChar() // 跳过第一个斜杠
	l.readChar() // 跳过第二个斜杠
	for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
		l.readChar()
	}
}

// 修改 readIdentifier 方法，支持标识符中的数字
func (l *Lexer) readIdentifier() string {
	position := l.position

	// 第一个字符必须是字母或下划线
	if !isLetter(l.ch) {
		return ""
	}

	// 后续字符可以是字母、数字或下划线
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
	"set":     SET,
	"track":   TRACK,
	"section": SECTION,

	// 大小写都映射
	"c":       NOTE_C,
	"d":       NOTE_D,
	"e":       NOTE_E,
	"f":       NOTE_F,
	"g":       NOTE_G,
	"a":       NOTE_A,
	"b":       NOTE_B,
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
