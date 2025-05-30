package dotenv

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var parseTestCases = []struct {
	file     string
	expected []ParseEntry
}{
	{
		"fixtures/basic.env",
		[]ParseEntry{
			{Key: "EXPORTED", Value: "data", Raw: false},
			{Key: "UNEXPORTED", Value: "data", Raw: false},
			{Key: "SINGLE_QUOTE", Value: "single quote", Raw: true},
			{Key: "DOUBLE_QUOTE", Value: "double quote", Raw: false},
			{Key: "UNQUOTED", Value: "unquoted data", Raw: false},
			{Key: "WITH_COMMENT", Value: "some data", Raw: false},
			{Key: "HASH_WITH_COMMENT", Value: "some#data", Raw: false},
		},
	},
	{
		"fixtures/broken.env",
		[]ParseEntry{
			{Key: "EXPORTED", Value: "exported data", Raw: false},
			{Key: "EMPTY", Value: "", Raw: false},
			{Key: "EMPTY_WITH_COMMENT", Value: "", Raw: false},
			{Key: "FINAL", Value: "valid", Raw: false},
		},
	},
	{
		"fixtures/replacement.env",
		[]ParseEntry{
			{Key: "VALUE", Value: "inserted", Raw: false},
			{Key: "REPLACE", Value: "${VALUE}", Raw: false},
			{Key: "REPLACE_SINGLE", Value: "${VALUE}", Raw: true},
			{Key: "REPLACE_DOUBLE", Value: "${VALUE}", Raw: false},
			{Key: "REPLACE_PARTIAL", Value: "partialy ${VALUE} value", Raw: false},
			{Key: "REPLACE_ESCAPED", Value: "partialy \\${VALUE} value", Raw: false},
			{Key: "REPLACE_FROM_BASIC", Value: "${HASH_WITH_COMMENT}", Raw: false},
			{Key: "REPLACE_FROM_BROKEN", Value: "${EMPTY}", Raw: false},
		},
	},
}

func TestParse(t *testing.T) {
	for _, tc := range parseTestCases {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			require.Nil(t, err, "failed to load fixture: %s", err)

			p := newParser(newLexer(string(data)))
			require.Equal(t, tc.expected, p.Parse())
		})
	}
}

var parseStrictTestCases = []struct {
	file          string
	expected      []ParseEntry
	expectedError error
}{
	{
		"fixtures/basic.env",
		[]ParseEntry{
			{Key: "EXPORTED", Value: "data", Raw: false},
			{Key: "UNEXPORTED", Value: "data", Raw: false},
			{Key: "SINGLE_QUOTE", Value: "single quote", Raw: true},
			{Key: "DOUBLE_QUOTE", Value: "double quote", Raw: false},
			{Key: "UNQUOTED", Value: "unquoted data", Raw: false},
			{Key: "WITH_COMMENT", Value: "some data", Raw: false},
			{Key: "HASH_WITH_COMMENT", Value: "some#data", Raw: false},
		},
		nil,
	},
	{
		"fixtures/broken.env",
		nil,
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
	},
	{
		"fixtures/replacement.env",
		[]ParseEntry{
			{Key: "VALUE", Value: "inserted", Raw: false},
			{Key: "REPLACE", Value: "${VALUE}", Raw: false},
			{Key: "REPLACE_SINGLE", Value: "${VALUE}", Raw: true},
			{Key: "REPLACE_DOUBLE", Value: "${VALUE}", Raw: false},
			{Key: "REPLACE_PARTIAL", Value: "partialy ${VALUE} value", Raw: false},
			{Key: "REPLACE_ESCAPED", Value: "partialy \\${VALUE} value", Raw: false},
			{Key: "REPLACE_FROM_BASIC", Value: "${HASH_WITH_COMMENT}", Raw: false},
			{Key: "REPLACE_FROM_BROKEN", Value: "${EMPTY}", Raw: false},
		},
		nil,
	},
}

func TestParseStrict(t *testing.T) {
	for _, tc := range parseStrictTestCases {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			require.Nil(t, err, "failed to load fixture: %s", err)

			p := newParser(newLexer(string(data)))
			pairs, err := p.ParseStrict()
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
				require.Empty(t, pairs)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, pairs)
			}
		})
	}
}
