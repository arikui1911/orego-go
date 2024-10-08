// Code generated by "stringer -type=Tag token.go"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INVALID-0]
	_ = x[EOF-1]
	_ = x[NEWLINE-2]
	_ = x[LPAREN-3]
	_ = x[RPAREN-4]
	_ = x[LBRACKET-5]
	_ = x[RBRACKET-6]
	_ = x[LBRACE-7]
	_ = x[RBRACE-8]
	_ = x[ARROW-9]
	_ = x[COMMA-10]
	_ = x[SEMICOLON-11]
	_ = x[COLON-12]
	_ = x[BANG-13]
	_ = x[LET-14]
	_ = x[EQ-15]
	_ = x[NE-16]
	_ = x[GE-17]
	_ = x[LE-18]
	_ = x[GT-19]
	_ = x[LT-20]
	_ = x[ADD-21]
	_ = x[SUB-22]
	_ = x[MUL-23]
	_ = x[DIV-24]
	_ = x[MOD-25]
	_ = x[LET_ADD-26]
	_ = x[LET_SUB-27]
	_ = x[LET_MUL-28]
	_ = x[LET_DIV-29]
	_ = x[LET_MOD-30]
	_ = x[KW_DEF-31]
	_ = x[KW_IF-32]
	_ = x[KW_ELSE-33]
	_ = x[KW_ELSIF-34]
	_ = x[KW_WHILE-35]
	_ = x[KW_BREAK-36]
	_ = x[KW_CONTINUE-37]
	_ = x[KW_RETURN-38]
	_ = x[KW_TRUE-39]
	_ = x[KW_FALSE-40]
	_ = x[KW_NIL-41]
	_ = x[IDENTIFIER-42]
	_ = x[LITERAL_INT-43]
	_ = x[LITERAL_FLOAT-44]
	_ = x[LITERAL_STRING-45]
}

const _Tag_name = "INVALIDEOFNEWLINELPARENRPARENLBRACKETRBRACKETLBRACERBRACEARROWCOMMASEMICOLONCOLONBANGLETEQNEGELEGTLTADDSUBMULDIVMODLET_ADDLET_SUBLET_MULLET_DIVLET_MODKW_DEFKW_IFKW_ELSEKW_ELSIFKW_WHILEKW_BREAKKW_CONTINUEKW_RETURNKW_TRUEKW_FALSEKW_NILIDENTIFIERLITERAL_INTLITERAL_FLOATLITERAL_STRING"

var _Tag_index = [...]uint16{0, 7, 10, 17, 23, 29, 37, 45, 51, 57, 62, 67, 76, 81, 85, 88, 90, 92, 94, 96, 98, 100, 103, 106, 109, 112, 115, 122, 129, 136, 143, 150, 156, 161, 168, 176, 184, 192, 203, 212, 219, 227, 233, 243, 254, 267, 281}

func (i Tag) String() string {
	if i < 0 || i >= Tag(len(_Tag_index)-1) {
		return "Tag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Tag_name[_Tag_index[i]:_Tag_index[i+1]]
}
