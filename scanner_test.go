package metadata_test

import (
	"strings"
	"testing"

	"github.com/go-chef/metadata-parser"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok metadata.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: metadata.EOF},
		{s: `/`, tok: metadata.ILLEGAL, lit: `/`},
		{s: ` `, tok: metadata.WS, lit: " "},
		{s: "\t", tok: metadata.WS, lit: "\t"},
		{s: "\n", tok: metadata.WS, lit: "\n"},
		{s: `#`, tok: metadata.COMMENT},

		// Identifiers
		{s: `foo`, tok: metadata.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: metadata.IDENT, lit: `Zx12_3U_`},
		{s: ",", tok: metadata.COMMA, lit: ","},

		// Keywords
		{s: `version`, tok: metadata.VERSION},
		{s: `depends`, tok: metadata.DEPENDS},
		{s: `name`, tok: metadata.NAME},
	}

	for i, tt := range tests {
		s := metadata.NewScanner(strings.NewReader(tt.s))
		tok, _, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
