package metadata

import "strings"

type Token int

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	EQ:  "=",
	LT:  "<",
	LTE: "<=",
	GT:  ">",
	GTE: ">=",
	PGT: "~>",

	IDENT:     "IDENT",
	STRING:    "STRING",
	BADSTRING: "BADSTRING",
	BADESCAPE: "BADESCAPE",
	TRUE:      "TRUE",
	FALSE:     "FALSE",

	COMMENT: "#",
	COMMA:   ",",

	VERSION: "version",
	NAME:    "name",
	DEPENDS: "depends",
}

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	operator_beg
	// operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	// Verion Specifiers
	EQ  // =
	LT  // <
	LTE // <=
	GT  // >
	GTE // >=
	PGT // ~> the pessimistic gt

	operator_end

	literal_beg
	// Literals
	IDENT
	STRING
	BADSTRING
	BADESCAPE
	TRUE
	FALSE
	literal_end

	// Misc Characters
	COMMENT
	COMMA // ,

	keyword_beg
	//Keywords
	VERSION
	NAME
	DEPENDS
	keyword_end
)

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok := keyword_beg + 1; tok < keyword_end; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	keywords["true"] = TRUE
	keywords["false"] = FALSE
}

// Pos specifies the line and character position of a token.
// The Char and Line are both zero-based indexes.
type Pos struct {
	Line int
	Char int
}

// String returns the string representation of the token.
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return ""
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
