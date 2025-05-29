package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nextTokenTestCases = []struct {
	file     string
	expected []token
}{
	{
		"fixtures/basic.env",
		[]token{
			{Line: 0, Pos: 1, Type: "EXPORT", Literal: "export"},
			{Line: 0, Pos: 8, Type: "IDENT", Literal: "EXPORTED"},
			{Line: 0, Pos: 16, Type: "EQUALS", Literal: "="},
			{Line: 0, Pos: 17, Type: "VALUE", Literal: "data"},
			{Line: 0, Pos: 21, Type: "EOL", Literal: ""},
			{Line: 1, Pos: 1, Type: "IDENT", Literal: "UNEXPORTED"},
			{Line: 1, Pos: 11, Type: "EQUALS", Literal: "="},
			{Line: 1, Pos: 12, Type: "VALUE", Literal: "data"},
			{Line: 1, Pos: 16, Type: "EOL", Literal: ""},
			{Line: 2, Pos: 1, Type: "IDENT", Literal: "SINGLE_QUOTE"},
			{Line: 2, Pos: 13, Type: "EQUALS", Literal: "="},
			{Line: 2, Pos: 14, Type: "VALUE", Literal: "single quote"},
			{Line: 2, Pos: 28, Type: "EOL", Literal: ""},
			{Line: 3, Pos: 1, Type: "IDENT", Literal: "DOUBLE_QUOTE"},
			{Line: 3, Pos: 13, Type: "EQUALS", Literal: "="},
			{Line: 3, Pos: 14, Type: "VALUE", Literal: "double quote"},
			{Line: 3, Pos: 28, Type: "EOL", Literal: ""},
			{Line: 4, Pos: 1, Type: "IDENT", Literal: "UNQUOTED"},
			{Line: 4, Pos: 9, Type: "EQUALS", Literal: "="},
			{Line: 4, Pos: 10, Type: "VALUE", Literal: "unquoted data"},
			{Line: 4, Pos: 23, Type: "EOL", Literal: ""},
			{Line: 5, Pos: 1, Type: "IDENT", Literal: "WITH_COMMENT"},
			{Line: 5, Pos: 13, Type: "EQUALS", Literal: "="},
			{Line: 5, Pos: 14, Type: "VALUE", Literal: "some data"},
			{Line: 5, Pos: 24, Type: "COMMENT", Literal: "with a comment"},
			{Line: 5, Pos: 40, Type: "EOL", Literal: ""},
			{Line: 6, Pos: 1, Type: "IDENT", Literal: "HASH_WITH_COMMENT"},
			{Line: 6, Pos: 18, Type: "EQUALS", Literal: "="},
			{Line: 6, Pos: 19, Type: "VALUE", Literal: "some#data"},
			{Line: 6, Pos: 29, Type: "COMMENT", Literal: "with a comment"},
			{Line: 6, Pos: 45, Type: "EOL", Literal: ""},
			{Line: 7, Pos: 1, Type: "EOF", Literal: ""},
		},
	},
	{
		"fixtures/broken.env",
		[]token{
			{Line: 0, Pos: 1, Type: "IDENT", Literal: "just"},
			{Line: 0, Pos: 6, Type: "IDENT", Literal: "some"},
			{Line: 0, Pos: 11, Type: "IDENT", Literal: "words"},
			{Line: 0, Pos: 16, Type: "EOL", Literal: ""},
			{Line: 1, Pos: 1, Type: "COMMENT", Literal: "comment"},
			{Line: 1, Pos: 10, Type: "EOL", Literal: ""},
			{Line: 2, Pos: 1, Type: "EXPORT", Literal: "export"},
			{Line: 2, Pos: 8, Type: "IDENT", Literal: "EXPORTED"},
			{Line: 2, Pos: 16, Type: "EQUALS", Literal: "="},
			{Line: 2, Pos: 17, Type: "VALUE", Literal: "data"},
			{Line: 2, Pos: 21, Type: "EOL", Literal: ""},
			{Line: 3, Pos: 1, Type: "IDENT", Literal: "UNEXPORTED"},
			{Line: 3, Pos: 11, Type: "EQUALS", Literal: "="},
			{Line: 3, Pos: 12, Type: "EOL", Literal: ""},
			{Line: 4, Pos: 1, Type: "EQUALS", Literal: "="},
			{Line: 4, Pos: 2, Type: "VALUE", Literal: "single quote"},
			{Line: 4, Pos: 16, Type: "EOL", Literal: ""},
			{Line: 5, Pos: 1, Type: "IDENT", Literal: "DOUBLE_QUOTE"},
			{Line: 5, Pos: 13, Type: "EQUALS", Literal: "="},
			{Line: 5, Pos: 15, Type: "COMMENT", Literal: "data"},
			{Line: 5, Pos: 21, Type: "EOL", Literal: ""},
			{Line: 6, Pos: 1, Type: "IDENT", Literal: "FINAL"},
			{Line: 6, Pos: 6, Type: "EQUALS", Literal: "="},
			{Line: 6, Pos: 7, Type: "VALUE", Literal: "valid"},
			{Line: 6, Pos: 12, Type: "EOL", Literal: ""},
			{Line: 7, Pos: 1, Type: "EOF", Literal: ""},
		},
	},
}

func TestNextToken(t *testing.T) {
	for _, tc := range nextTokenTestCases {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			require.Nil(t, err, "failed to load fixture: %s", err)

			l := newLexer(string(data))

			var tkns []token
			for {
				tkn := l.NextToken()
				tkns = append(tkns, tkn)
				if tkn.Type == tknEOF {
					break
				}
			}

			assert.Equal(t, tc.expected, tkns)
		})
	}
}
