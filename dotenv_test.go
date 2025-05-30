package dotenv

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var loadTestCases = []struct {
	name     string
	files    []string
	expected []string
}{
	{
		"basic",
		[]string{"fixtures/basic.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
	},
	{
		"broken",
		[]string{"fixtures/broken.env"},
		[]string{
			"EMPTY=",
			"EMPTY_WITH_COMMENT=",
			"EXPORTED=exported data",
			"FINAL=valid",
		},
	},
	{
		"replacement",
		[]string{"fixtures/replacement.env"},
		[]string{
			"VALUE=inserted",
			"REPLACE=inserted",
			"REPLACE_SINGLE=${VALUE}",
			"REPLACE_DOUBLE=inserted",
			"REPLACE_PARTIAL=partialy inserted value",
			"REPLACE_ESCAPED=partialy ${VALUE} value",
			"REPLACE_FROM_BASIC=",
			"REPLACE_FROM_BROKEN=",
		},
	},
	{
		"multi file",
		[]string{"fixtures/basic.env", "fixtures/broken.env", "fixtures/replacement.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
			"EMPTY=",
			"EMPTY_WITH_COMMENT=",
			"FINAL=valid",
			"VALUE=inserted",
			"REPLACE=inserted",
			"REPLACE_SINGLE=${VALUE}",
			"REPLACE_DOUBLE=inserted",
			"REPLACE_PARTIAL=partialy inserted value",
			"REPLACE_ESCAPED=partialy ${VALUE} value",
			"REPLACE_FROM_BASIC=some#data",
			"REPLACE_FROM_BROKEN=",
		},
	},
}

func TestLoad(t *testing.T) {
	for _, tc := range loadTestCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			err := Load(tc.files...)
			require.Nil(t, err)

			require.ElementsMatch(t, tc.expected, os.Environ())
		})
	}
}

var loadStrictTestCases = []struct {
	name          string
	files         []string
	expected      []string
	expectedError error
}{
	{
		"basic",
		[]string{"fixtures/basic.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
		nil,
	},
	{
		"broken",
		[]string{"fixtures/broken.env"},
		[]string{},
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
	},
	{
		"multi file",
		[]string{"fixtures/basic.env", "fixtures/broken.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
	},
}

func TestLoadStrict(t *testing.T) {
	for _, tc := range loadStrictTestCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			err := LoadStrict(tc.files...)
			require.Equal(t, tc.expectedError, err)

			require.ElementsMatch(t, tc.expected, os.Environ())
		})
	}
}

var overloadTestCases = []struct {
	name     string
	files    []string
	expected []string
}{
	{
		"basic",
		[]string{"fixtures/basic.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
	},
	{
		"broken",
		[]string{"fixtures/broken.env"},
		[]string{
			"EMPTY=",
			"EMPTY_WITH_COMMENT=",
			"EXPORTED=exported data",
			"FINAL=valid",
		},
	},
	{
		"replacement",
		[]string{"fixtures/replacement.env"},
		[]string{

			"VALUE=inserted",
			"REPLACE=inserted",
			"REPLACE_SINGLE=${VALUE}",
			"REPLACE_DOUBLE=inserted",
			"REPLACE_PARTIAL=partialy inserted value",
			"REPLACE_ESCAPED=partialy ${VALUE} value",
			"REPLACE_FROM_BASIC=",
			"REPLACE_FROM_BROKEN=",
		},
	},
	{
		"multi file",
		[]string{"fixtures/basic.env", "fixtures/broken.env", "fixtures/replacement.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=exported data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
			"EMPTY=",
			"EMPTY_WITH_COMMENT=",
			"FINAL=valid",
			"VALUE=inserted",
			"REPLACE=inserted",
			"REPLACE_SINGLE=${VALUE}",
			"REPLACE_DOUBLE=inserted",
			"REPLACE_PARTIAL=partialy inserted value",
			"REPLACE_ESCAPED=partialy ${VALUE} value",
			"REPLACE_FROM_BASIC=some#data",
			"REPLACE_FROM_BROKEN=",
		},
	},
}

func TestOverload(t *testing.T) {
	for _, tc := range overloadTestCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			err := Overload(tc.files...)
			require.Nil(t, err)

			require.ElementsMatch(t, tc.expected, os.Environ())
		})
	}
}

var overloadStrictTestCases = []struct {
	name          string
	files         []string
	expected      []string
	expectedError error
}{
	{
		"basic",
		[]string{"fixtures/basic.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
		nil,
	},
	{
		"broken",
		[]string{"fixtures/broken.env"},
		[]string{},
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
	},
	{
		"multi file",
		[]string{"fixtures/basic.env", "fixtures/broken.env"},
		[]string{
			"DOUBLE_QUOTE=double quote",
			"EXPORTED=data",
			"HASH_WITH_COMMENT=some#data",
			"SINGLE_QUOTE=single quote",
			"UNEXPORTED=data",
			"UNQUOTED=unquoted data",
			"WITH_COMMENT=some data",
		},
		errors.New("Unexpected token IDENT value=some line=0 pos=6"),
	},
}

func TestOverloadStrict(t *testing.T) {
	for _, tc := range overloadStrictTestCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			err := OverloadStrict(tc.files...)
			require.Equal(t, tc.expectedError, err)

			require.ElementsMatch(t, tc.expected, os.Environ())
		})
	}
}

func TestParseFile(t *testing.T) {
	// NB: test cases come from parse_test, this is basically the same test just using the
	//     ParseFile api function
	for _, tc := range parseTestCases {
		t.Run(tc.file, func(t *testing.T) {
			p, err := ParseFile(tc.file)

			require.Nil(t, err)
			require.Equal(t, tc.expected, p.Parse())
		})
	}
}
