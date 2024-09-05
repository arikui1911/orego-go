package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"

	"github.com/arikui1911/orego-go/token"
)

type Lexer struct {
	Line         int
	Column       int
	src          *bufio.Reader
	isEOF        bool
	savedRune    rune
	hasSavedRune bool
	prevLine     int
	prevCol      int
	nextLine     int
	nextCol      int
	lastTag      token.Tag
}

func New(src io.Reader) *Lexer {
	return &Lexer{
		Line:     1,
		Column:   1,
		src:      bufio.NewReader(src),
		prevLine: 1,
		prevCol:  1,
		nextLine: 1,
		nextCol:  1,
	}
}

func (l *Lexer) NextToken() (t token.Token, err error) {
	defer func() { l.lastTag = t.Tag }()

	err = l.skipSpacesAndComments()
	if err != nil {
		return
	}

	c, err := l.getc()
	if err == io.EOF {
		l.Column++
		t.Tag = token.EOF
		if l.isNewlineRequired() {
			t.Tag = token.NEWLINE
			t.Value = "\n"
		}
		l.copyLocation(&t.Location, &t.Location)
		err = nil
		return
	}
	if err != nil {
		return
	}
	l.copyLocation(&t.Location, nil)

	switch {
	case c == '\n':
		err = l.scanNewline(&t)
	case c == '}':
		if l.isNewlineRequired() {
			l.ungetc(c)
			t.Tag = token.NEWLINE
			t.Value = "\n"
			l.copyLocation(nil, &t.Location)
			return
		}
		t.Tag = token.RBRACE
		t.Value = "}"
		l.copyLocation(nil, &t.Location)
		return
	case c == '"':
		err = l.scanString(&t)
	case c == '0':
		err = l.scanZero(&t)
	case isDigit(c):
		err = l.scanInt(&t, c)
	case isSymbol(c):
		err = l.scanIdent(&t, c)
	default:
		err = l.scanOperator(&t, c)
	}
	return
}

func (l *Lexer) skipSpacesAndComments() error {
	inComment := false
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch c {
		case '\n':
			inComment = false
			if l.isNewlineRequired() {
				l.ungetc(c)
				return nil
			}
		case '#':
			inComment = true
		default:
			if !inComment && !isSpace(c) {
				l.ungetc(c)
				return nil
			}
		}
	}
	return nil
}

func (l *Lexer) scanNewline(t *token.Token) error {
	t.Tag = token.NEWLINE
	t.Value = "\n"
	l.copyLocation(nil, &t.Location)
	return nil
}

func (l *Lexer) scanString(t *token.Token) error {
	buf := []rune{}
	for {
		c, err := l.getc()
		if err == io.EOF {
			return fmt.Errorf("%v: unterminated string literal", t.Location)
		}
		if err != nil {
			return err
		}
		if c == '"' {
			t.Tag = token.LITERAL_STRING
			t.Value = string(buf)
			l.copyLocation(nil, &t.Location)
			break
		}
		buf = append(buf, c)
	}
	return nil
}

func (l *Lexer) scanZero(t *token.Token) error {
	c, err := l.getc()
	switch {
	case err == io.EOF:
		// do nothing
	case err != nil:
		return err
	default:
		if c == '.' {
			return l.scanFloat(t, []rune{'0', '.'})
		}
		l.ungetc(c)
	}
	t.Tag = token.LITERAL_INT
	t.Value = "0"
	l.copyLocation(nil, &t.Location)
	return nil
}

func (l *Lexer) scanInt(t *token.Token, c rune) error {
	buf := []rune{c}
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if c == '.' {
			buf = append(buf, c)
			return l.scanFloat(t, buf)
		}
		if !isDigit(c) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	t.Tag = token.LITERAL_INT
	t.Value = string(buf)
	l.copyLocation(nil, &t.Location)
	return nil
}

func (l *Lexer) scanFloat(t *token.Token, buf []rune) error {
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !isDigit(c) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	t.Tag = token.LITERAL_FLOAT
	t.Value = string(buf)
	l.copyLocation(nil, &t.Location)
	return nil
}

var keywords = map[string]token.Tag{
	"def":      token.KW_DEF,
	"if":       token.KW_IF,
	"else":     token.KW_ELSE,
	"elsif":    token.KW_ELSIF,
	"while":    token.KW_WHILE,
	"break":    token.KW_BREAK,
	"continue": token.KW_CONTINUE,
	"return":   token.KW_RETURN,
	"true":     token.KW_TRUE,
	"false":    token.KW_FALSE,
	"nil":      token.KW_NIL,
}

func (l *Lexer) scanIdent(t *token.Token, c rune) error {
	buf := []rune{c}
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !isSymbol(c) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	t.Tag = token.IDENTIFIER
	t.Value = string(buf)
	l.copyLocation(nil, &t.Location)

	if tt, ok := keywords[t.Value]; ok {
		t.Tag = tt
	}
	return nil
}

var operators = map[string]token.Tag{
	"(": token.LPAREN,
	")": token.RPAREN,
	"[": token.LBRACKET,
	"]": token.RBRACKET,
	"{": token.LBRACE,
	// "}":  token.RBRACE,
	"->": token.ARROW,
	",":  token.COMMA,
	";":  token.SEMICOLON,
	":":  token.COLON,
	"=":  token.LET,
	"==": token.EQ,
	"!=": token.NE,
	">=": token.GE,
	"<=": token.LE,
	">":  token.GT,
	"<":  token.LT,
	"+":  token.ADD,
	"-":  token.SUB,
	"*":  token.MUL,
	"/":  token.DIV,
	"%":  token.MOD,
	"+=": token.LET_ADD,
	"-=": token.LET_SUB,
	"*=": token.LET_MUL,
	"/=": token.LET_DIV,
	"%=": token.LET_MOD,
	"!":  token.BANG,
}

func (l *Lexer) scanOperator(t *token.Token, c rune) error {
	buf := []rune{c}
	for {
		if _, ok := operators[string(buf)]; !ok {
			break
		}
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		buf = append(buf, c)
	}
	last := buf[len(buf)-1]
	var s string
	if l.isEOF {
		s = string(buf)
	} else {
		s = string(buf[:len(buf)-1])
	}
	if tt, ok := operators[s]; ok {
		if !l.isEOF {
			l.ungetc(last)
		}
		t.Tag = tt
		t.Value = s
		l.copyLocation(nil, &t.Location)
		return nil
	}
	return fmt.Errorf("%v: invalid character - '%c'", t.Location, last)
}

func (l *Lexer) isNewlineRequired() bool {
	switch l.lastTag {
	case token.RPAREN,
		token.RBRACKET,
		token.RBRACE,
		token.KW_BREAK,
		token.KW_CONTINUE,
		token.KW_RETURN,
		token.KW_TRUE,
		token.KW_FALSE,
		token.KW_NIL,
		token.IDENTIFIER,
		token.LITERAL_INT,
		token.LITERAL_FLOAT,
		token.LITERAL_STRING:
		return true
	}
	return false
}

func isSpace(c rune) bool {
	return unicode.IsSpace(c)
}

func isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

func isSymbol(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || isDigit(c)
}

func (l *Lexer) copyLocation(start *token.Location, end *token.Location) {
	if start != nil {
		start.StartLine = l.Line
		start.StartColumn = l.Column
	}
	if end != nil {
		end.EndLine = l.Line
		end.EndColumn = l.Column
	}
}

func (l *Lexer) getc() (rune, error) {
	c, err := l.getcCore()
	if err != nil {
		return 0, err
	}
	l.prevLine = l.Line
	l.prevCol = l.Column
	l.Line = l.nextLine
	l.Column = l.nextCol
	l.nextCol++
	if c == '\n' {
		l.nextLine++
		l.nextCol = 1
	}
	return c, nil
}

func (l *Lexer) getcCore() (rune, error) {
	if l.hasSavedRune {
		l.hasSavedRune = false
		return l.savedRune, nil
	}
	c, _, err := l.src.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.isEOF = true
		}
		return 0, err
	}
	return c, nil
}

func (l *Lexer) ungetc(c rune) {
	l.savedRune = c
	l.hasSavedRune = true
	l.nextLine = l.Line
	l.nextCol = l.Column
	l.Line = l.prevLine
	l.Column = l.prevCol
}
