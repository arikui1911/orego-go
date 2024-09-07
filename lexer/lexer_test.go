package lexer_test

import (
	"strings"
	"testing"

	"github.com/arikui1911/orego-go/lexer"
	"github.com/arikui1911/orego-go/token"

	"github.com/google/go-cmp/cmp"
)

func TestLexer(t *testing.T) {
	data := []struct {
		name string
		src  string
		tag  token.Tag
		val  string
		ec   int
	}{
		{"left paren", `(`, token.LPAREN, "(", 1},
		{"right paren", `)`, token.RPAREN, ")", 1},
		{"left bracket", `[`, token.LBRACKET, "[", 1},
		{"right bracket", `]`, token.RBRACKET, "]", 1},
		{"left brrace", `{`, token.LBRACE, "{", 1},
		{"right brace", `}`, token.RBRACE, "}", 1},
		{"arrow", `->`, token.ARROW, "->", 2},
		{"comma", `,`, token.COMMA, ",", 1},
		{"semicolon", `;`, token.SEMICOLON, ";", 1},
		{"colon", `:`, token.COLON, ":", 1},
		{"let", `=`, token.LET, "=", 1},
		{"eq", `==`, token.EQ, "==", 2},
		{"ne", `!=`, token.NE, "!=", 2},
		{"ge", `>=`, token.GE, ">=", 2},
		{"le", `<=`, token.LE, "<=", 2},
		{"gt", `>`, token.GT, ">", 1},
		{"lt", `<`, token.LT, "<", 1},
		{"add", `+`, token.ADD, "+", 1},
		{"sub", `-`, token.SUB, "-", 1},
		{"mul", `*`, token.MUL, "*", 1},
		{"div", `/`, token.DIV, "/", 1},
		{"mod", `%`, token.MOD, "%", 1},
		{"not", `!`, token.BANG, "!", 1},
		{"let add", `+=`, token.LET_ADD, "+=", 2},
		{"let sub", `-=`, token.LET_SUB, "-=", 2},
		{"let mul", `*=`, token.LET_MUL, "*=", 2},
		{"let div", `/=`, token.LET_DIV, "/=", 2},
		{"let mod", `%=`, token.LET_MOD, "%=", 2},
		{"not", `!`, token.BANG, "!", 1},
		{"def", `def`, token.KW_DEF, "def", 3},
		{"if", `if`, token.KW_IF, "if", 2},
		{"else", `else`, token.KW_ELSE, "else", 4},
		{"elsif", `elsif`, token.KW_ELSIF, "elsif", 5},
		{"while", `while`, token.KW_WHILE, "while", 5},
		{"break", `break`, token.KW_BREAK, "break", 5},
		{"continue", `continue`, token.KW_CONTINUE, "continue", 8},
		{"return", `return`, token.KW_RETURN, "return", 6},
		{"true", `true`, token.KW_TRUE, "true", 4},
		{"false", `false`, token.KW_FALSE, "false", 5},
		{"nil", `nil`, token.KW_NIL, "nil", 3},
		{"ident", `hoge_123`, token.IDENTIFIER, "hoge_123", 8},
		{"int literal", `123`, token.LITERAL_INT, "123", 3},
		{"zero literal", `0`, token.LITERAL_INT, "0", 1},
		{"float literal", `12.3`, token.LITERAL_FLOAT, "12.3", 4},
		{"float literal 2", `0.12`, token.LITERAL_FLOAT, "0.12", 4},
		{"string literal", `"Hello"`, token.LITERAL_STRING, "Hello", 7},
	}

	for _, d := range data {
		expected := token.Token{
			Tag:   d.tag,
			Value: d.val,
			Location: token.Location{
				StartLine:   1,
				StartColumn: 1,
				EndLine:     1,
				EndColumn:   d.ec,
			},
		}
		t.Run(d.name, func(t *testing.T) { testSrcNextToken(t, d.src, expected) })
		t.Run(d.name+" (delim sp)", func(t *testing.T) { testSrcNextToken(t, d.src+" ", expected) })
	}
}

func TestNewline(t *testing.T) {
	data := []struct {
		name string
		src  string
		tag  token.Tag
		val  string
		ec   int
	}{
		{"right paren", `)`, token.RPAREN, ")", 1},
		{"right bracket", `]`, token.RBRACKET, "]", 1},
		{"right brace", `}`, token.RBRACE, "}", 1},
		{"break", `break`, token.KW_BREAK, "break", 5},
		{"continue", `continue`, token.KW_CONTINUE, "continue", 8},
		{"return", `return`, token.KW_RETURN, "return", 6},
		{"true", `true`, token.KW_TRUE, "true", 4},
		{"false", `false`, token.KW_FALSE, "false", 5},
		{"nil", `nil`, token.KW_NIL, "nil", 3},
		{"ident", `hoge`, token.IDENTIFIER, "hoge", 4},
		{"int literal", `123`, token.LITERAL_INT, "123", 3},
		{"float literal", `12.3`, token.LITERAL_FLOAT, "12.3", 4},
		{"string literal", `"Hello"`, token.LITERAL_STRING, "Hello", 7},
	}

	for _, d := range data {
		expected := token.Token{
			Tag:   d.tag,
			Value: d.val,
			Location: token.Location{
				StartLine:   1,
				StartColumn: 1,
				EndLine:     1,
				EndColumn:   d.ec,
			},
		}
		t.Run(d.name, func(t *testing.T) {
			l := lexer.New(strings.NewReader(d.src))
			testNextToken(t, l, expected)
			testNextToken(t, l, token.Token{
				Tag:   token.NEWLINE,
				Value: "\n",
				Location: token.Location{
					StartLine:   1,
					StartColumn: d.ec + 1,
					EndLine:     1,
					EndColumn:   d.ec + 1,
				},
			})
			testNextToken(t, l, token.Token{
				Tag:   token.EOF,
				Value: "",
				Location: token.Location{
					StartLine:   1,
					StartColumn: d.ec + 1,
					EndLine:     1,
					EndColumn:   d.ec + 1,
				},
			})
		})
	}
}

func testSrcNextToken(t *testing.T, src string, expected token.Token) {
	testNextToken(t, lexer.New(strings.NewReader(src)), expected)
}

func testNextToken(t *testing.T, l *lexer.Lexer, expected token.Token) {
	actual, err := l.NextToken()
	if err != nil {
		t.Error(err)
		return
	}
	diff := cmp.Diff(expected, actual)
	if diff != "" {
		t.Errorf("want <-> got <+>:\n%s", diff)
	}
}

func TestComment(t *testing.T) {
	src := `
# comment
hoge
	`
	testSrcNextToken(t, src, token.Token{
		Tag:      token.IDENTIFIER,
		Value:    "hoge",
		Location: token.Location{StartLine: 3, StartColumn: 1, EndLine: 3, EndColumn: 4},
	})
}

func TestPostComment(t *testing.T) {
	src := `hoge  # comment`
	testSrcNextToken(t, src, token.Token{
		Tag:      token.IDENTIFIER,
		Value:    "hoge",
		Location: token.Location{StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 4},
	})
}

func TestIncompleteString(t *testing.T) {
	src := `"Hello`
	tok, err := lexer.New(strings.NewReader(src)).NextToken()
	if err == nil {
		t.Errorf("want error, got %v", tok)
	}
}

func TestInvalidCharacter(t *testing.T) {
	src := `@`
	tok, err := lexer.New(strings.NewReader(src)).NextToken()
	if err == nil {
		t.Errorf("want error, got %v", tok)
	}
}
