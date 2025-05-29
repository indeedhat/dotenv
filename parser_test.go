package dotenv

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var parseTestCases = []struct {
	file     string
	expected map[string]string
}{
	{
		"fixtures/basic.env",
		map[string]string{
			"DOUBLE_QUOTE":      "double quote",
			"EXPORTED":          "data",
			"HASH_WITH_COMMENT": "some#data",
			"SINGLE_QUOTE":      "single quote",
			"UNEXPORTED":        "data",
			"UNQUOTED":          "unquoted data",
			"WITH_COMMENT":      "some data",
		},
	},
	{
		"fixtures/broken.env",
		map[string]string{
			"EXPORTED": "data",
			"FINAL":    "valid",
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
	expected      map[string]string
	expectedError error
}{
	{
		"fixtures/basic.env",
		map[string]string{
			"DOUBLE_QUOTE":      "double quote",
			"EXPORTED":          "data",
			"HASH_WITH_COMMENT": "some#data",
			"SINGLE_QUOTE":      "single quote",
			"UNEXPORTED":        "data",
			"UNQUOTED":          "unquoted data",
			"WITH_COMMENT":      "some data",
		},
		nil,
	},
	{
		"fixtures/broken.env",
		map[string]string{},
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
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
