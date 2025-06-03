package dsl

import (
	"catRock/pkg/dsl/mytype"
	"fmt"
)

type TokenType int

// 关键字（Keyword）
// 标识符（Identifier）
// 字面量（Literal / Constant）
// 运算符（Operator）
// 分隔符/界符（Separator / Delimiter）
// 注释（Comment）（一般被丢弃或单独标记）
// 空白字符（Whitespace）（通常被跳过）
const (
	// 特殊关键字
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// 字面量
	NUMBER
	IDENTIFIER // 字符串

	// 关键字
	SET     //set
	TRACK   //track
	SECTION //section

	// 音符名称
	NOTE_C
	NOTE_D
	NOTE_E
	NOTE_F
	NOTE_G
	NOTE_A
	NOTE_B

	// 半音 #号
	NOTE_CS
	NOTE_DS
	NOTE_FS
	NOTE_GS
	NOTE_AS

	// 半音 b号
	NOTE_DB  // Db
    NOTE_EB  // Eb
    NOTE_GB  // Gb
    NOTE_AB  // Ab
    NOTE_BB  // Bb

	// 符号
	COLON    //:
	SLASH    // /
	DOT      // .
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]
	LPAREN   // (
	RPAREN   // )
)

type Token struct {
	Type     TokenType // 令牌类型
	Literal  string    // 令牌的文本内容
	Position mytype.Position  // 令牌在源代码中的位置
}


// Token类型名称映射
var tokenNames = map[TokenType]string{
    ILLEGAL:    "ILLEGAL",
    EOF:        "EOF", 
    NEWLINE:    "NEWLINE",
    NUMBER:     "NUMBER",
    IDENTIFIER: "IDENTIFIER",
    SET:        "SET",
    TRACK:      "TRACK",
    SECTION:    "SECTION",
    NOTE_C:     "NOTE_C",
    NOTE_D:     "NOTE_D",
    NOTE_E:     "NOTE_E",
    NOTE_F:     "NOTE_F",
    NOTE_G:     "NOTE_G",
    NOTE_A:     "NOTE_A",
    NOTE_B:     "NOTE_B",

	// 半音符号
	    // 升号音符 (s后缀)
    NOTE_CS:    "NOTE_CS",
    NOTE_DS:    "NOTE_DS", 
    NOTE_FS:    "NOTE_FS",
    NOTE_GS:    "NOTE_GS",
    NOTE_AS:    "NOTE_AS",
    
    // 降号音符 (b后缀)
    NOTE_DB:    "NOTE_DB",
    NOTE_EB:    "NOTE_EB",
    NOTE_GB:    "NOTE_GB", 
    NOTE_AB:    "NOTE_AB",
    NOTE_BB:    "NOTE_BB",


    COLON:      "COLON",
    SLASH:      "SLASH",
    DOT:        "DOT",
    LBRACE:     "LBRACE",
    RBRACE:     "RBRACE",
    LBRACKET:   "LBRACKET",
    RBRACKET:   "RBRACKET",
    LPAREN:     "LPAREN",
    RPAREN:     "RPAREN",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", int(t))
}
