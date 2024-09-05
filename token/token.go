package token

import "fmt"

//go:generate stringer -type=Tag token.go
type Tag int

const (
	INVALID Tag = iota
	EOF
	NEWLINE

	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	ARROW
	COMMA
	SEMICOLON
	COLON
	BANG

	LET
	EQ
	NE
	GE
	LE
	GT
	LT
	ADD
	SUB
	MUL
	DIV
	MOD
	LET_ADD
	LET_SUB
	LET_MUL
	LET_DIV
	LET_MOD

	KW_DEF
	KW_IF
	KW_ELSE
	KW_ELSIF
	KW_WHILE
	KW_BREAK
	KW_CONTINUE
	KW_RETURN
	KW_TRUE
	KW_FALSE
	KW_NIL

	IDENTIFIER
	LITERAL_INT
	LITERAL_FLOAT
	LITERAL_STRING
)

type Location struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

func (l Location) String() string {
	return fmt.Sprintf("(%d:%d):(%d:%d)", l.StartLine, l.StartColumn, l.EndLine, l.EndColumn)
}

type Token struct {
	Tag      Tag
	Value    string
	Location Location
}
