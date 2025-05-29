package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var lexNextTokenTestCases = []struct {
	file     string
	expected []token
}{
	{
		"fixtures/basic.env",
		[]token{
			{Line: 0, Pos: 1, Type: tknExport, Literal: "export"},
			{Line: 0, Pos: 8, Type: tknIdentifier, Literal: "EXPORTED"},
			{Line: 0, Pos: 16, Type: tknEquals, Literal: "="},
			{Line: 0, Pos: 17, Type: tknValue, Literal: "data"},
			{Line: 0, Pos: 21, Type: tknEOL, Literal: ""},
			{Line: 1, Pos: 1, Type: tknIdentifier, Literal: "UNEXPORTED"},
			{Line: 1, Pos: 11, Type: tknEquals, Literal: "="},
			{Line: 1, Pos: 12, Type: tknValue, Literal: "data"},
			{Line: 1, Pos: 16, Type: tknEOL, Literal: ""},
			{Line: 2, Pos: 1, Type: tknIdentifier, Literal: "SINGLE_QUOTE"},
			{Line: 2, Pos: 13, Type: tknEquals, Literal: "="},
			{Line: 2, Pos: 14, Type: tknRawValue, Literal: "single quote"},
			{Line: 2, Pos: 28, Type: tknEOL, Literal: ""},
			{Line: 3, Pos: 1, Type: tknIdentifier, Literal: "DOUBLE_QUOTE"},
			{Line: 3, Pos: 13, Type: tknEquals, Literal: "="},
			{Line: 3, Pos: 14, Type: tknValue, Literal: "double quote"},
			{Line: 3, Pos: 28, Type: tknEOL, Literal: ""},
			{Line: 4, Pos: 1, Type: tknIdentifier, Literal: "UNQUOTED"},
			{Line: 4, Pos: 9, Type: tknEquals, Literal: "="},
			{Line: 4, Pos: 10, Type: tknValue, Literal: "unquoted data"},
			{Line: 4, Pos: 23, Type: tknEOL, Literal: ""},
			{Line: 5, Pos: 1, Type: tknIdentifier, Literal: "WITH_COMMENT"},
			{Line: 5, Pos: 13, Type: tknEquals, Literal: "="},
			{Line: 5, Pos: 14, Type: tknValue, Literal: "some data"},
			{Line: 5, Pos: 24, Type: tknComment, Literal: "with a comment"},
			{Line: 5, Pos: 40, Type: tknEOL, Literal: ""},
			{Line: 6, Pos: 1, Type: tknIdentifier, Literal: "HASH_WITH_COMMENT"},
			{Line: 6, Pos: 18, Type: tknEquals, Literal: "="},
			{Line: 6, Pos: 19, Type: tknValue, Literal: "some#data"},
			{Line: 6, Pos: 29, Type: tknComment, Literal: "with a comment"},
			{Line: 6, Pos: 45, Type: tknEOL, Literal: ""},
			{Line: 7, Pos: 1, Type: tknEOF, Literal: ""},
		},
	},
	{
		"fixtures/broken.env",
		[]token{
			{Line: 0, Pos: 1, Type: tknIdentifier, Literal: "just"},
			{Line: 0, Pos: 6, Type: tknIdentifier, Literal: "some"},
			{Line: 0, Pos: 11, Type: tknIdentifier, Literal: "words"},
			{Line: 0, Pos: 16, Type: tknEOL, Literal: ""},
			{Line: 1, Pos: 1, Type: tknComment, Literal: "comment"},
			{Line: 1, Pos: 10, Type: tknEOL, Literal: ""},
			{Line: 2, Pos: 1, Type: tknExport, Literal: "export"},
			{Line: 2, Pos: 8, Type: tknIdentifier, Literal: "EXPORTED"},
			{Line: 2, Pos: 16, Type: tknEquals, Literal: "="},
			{Line: 2, Pos: 17, Type: tknValue, Literal: "data"},
			{Line: 2, Pos: 21, Type: tknEOL, Literal: ""},
			{Line: 3, Pos: 1, Type: tknIdentifier, Literal: "UNEXPORTED"},
			{Line: 3, Pos: 11, Type: tknEquals, Literal: "="},
			{Line: 3, Pos: 12, Type: tknEOL, Literal: ""},
			{Line: 4, Pos: 1, Type: tknEquals, Literal: "="},
			{Line: 4, Pos: 2, Type: tknRawValue, Literal: "single quote"},
			{Line: 4, Pos: 16, Type: tknEOL, Literal: ""},
			{Line: 5, Pos: 1, Type: tknIdentifier, Literal: "DOUBLE_QUOTE"},
			{Line: 5, Pos: 13, Type: tknEquals, Literal: "="},
			{Line: 5, Pos: 15, Type: tknComment, Literal: "data"},
			{Line: 5, Pos: 21, Type: tknEOL, Literal: ""},
			{Line: 6, Pos: 1, Type: tknIdentifier, Literal: "FINAL"},
			{Line: 6, Pos: 6, Type: tknEquals, Literal: "="},
			{Line: 6, Pos: 7, Type: tknValue, Literal: "valid"},
			{Line: 6, Pos: 12, Type: tknEOL, Literal: ""},
			{Line: 7, Pos: 1, Type: tknEOF, Literal: ""},
		},
	},
}

func TestLexerNextToken(t *testing.T) {
	for _, tc := range lexNextTokenTestCases {
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
