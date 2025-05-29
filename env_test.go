package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type envarTestCase[S EnVar[T], T any] struct {
	subject        S
	env            map[string]string
	fallback       []T
	getExpected    T
	lookupExpected T
}

var stringTestCases = []envarTestCase[String, string]{
	{
		subject:        String("VALID"),
		env:            map[string]string{"VALID": "valid"},
		fallback:       []string{"fallback value"},
		getExpected:    "valid",
		lookupExpected: "valid",
	},
	{
		subject:     String("FALLBACK_EMPTY"),
		env:         map[string]string{"FALLBACK_EMPTY": ""},
		fallback:    []string{"fallback value"},
		getExpected: "fallback value",
	},
	{
		subject:        String("FALLBACK_UNSET"),
		fallback:       []string{"fallback value"},
		getExpected:    "fallback value",
		lookupExpected: "fallback value",
	},
	{
		subject: String("NO_FALLBACK_EMPTY"),
		env:     map[string]string{"NO_FALLBACK_EMPTY": ""},
	},
	{
		subject: String("NO_FALLBACK_UNSET"),
	},
}

func TestStringGet(t *testing.T) {
	for _, tc := range stringTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.getExpected, tc.subject.Get(tc.fallback...))
		})
	}
}

func TestStringLookup(t *testing.T) {
	for _, tc := range stringTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.lookupExpected, tc.subject.Lookup(tc.fallback...))
		})
	}
}

var intTestCases = []envarTestCase[Int, int]{
	{
		subject:        Int("VALID"),
		env:            map[string]string{"VALID": "100"},
		fallback:       []int{200},
		getExpected:    100,
		lookupExpected: 100,
	},
	{
		subject:     Int("FALLBACK_EMPTY"),
		env:         map[string]string{"FALLBACK_EMPTY": ""},
		fallback:    []int{200},
		getExpected: 200,
	},
	{
		subject:        Int("FALLBACK_UNSET"),
		fallback:       []int{200},
		getExpected:    200,
		lookupExpected: 200,
	},
	{
		subject: Int("NO_FALLBACK_EMYTY"),
		env:     map[string]string{"NO_FALLBACK_EMYTY": ""},
	},
	{
		subject: Int("NO_FALLBACK_UNSET"),
	},
	{
		subject: Int("NOT_INT"),
		env:     map[string]string{"NOT_INT": "is a string"},
	},
	{
		subject:     Int("NOT_INT_FALLBACK"),
		env:         map[string]string{"NOT_INT_FALLBACK": "invalid"},
		fallback:    []int{123},
		getExpected: 123,
	},
	{
		subject:        Int("HEX"),
		env:            map[string]string{"HEX": "0xF"},
		fallback:       []int{123},
		getExpected:    15,
		lookupExpected: 15,
	},
	{
		subject:        Int("BIN"),
		env:            map[string]string{"BIN": "0b10000000"},
		fallback:       []int{123},
		getExpected:    128,
		lookupExpected: 128,
	},
}

func TestIntGet(t *testing.T) {
	for _, tc := range intTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.getExpected, tc.subject.Get(tc.fallback...))
		})
	}
}

func TestIntLookup(t *testing.T) {
	for _, tc := range intTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.lookupExpected, tc.subject.Lookup(tc.fallback...))
		})
	}
}

var floatTestCases = []envarTestCase[Float, float64]{
	{
		subject:        Float("VALID"),
		env:            map[string]string{"VALID": "100.5"},
		fallback:       []float64{200},
		getExpected:    100.5,
		lookupExpected: 100.5,
	},
	{
		subject:     Float("FALLBACK_EMPTY"),
		env:         map[string]string{"FALLBACK_EMPTY": ""},
		fallback:    []float64{200},
		getExpected: 200,
	},
	{
		subject:        Float("FALLBACK_UNSET"),
		fallback:       []float64{200},
		getExpected:    200,
		lookupExpected: 200,
	},
	{
		subject: Float("NO_FALLBACK_EMPTY"),
		env:     map[string]string{"NO_FALLBACK_EMPTY": ""},
	},
	{
		subject: Float("NO_FALLBACK_UNSET"),
	},
	{
		subject: Float("NOT_FLOAT"),
		env:     map[string]string{"NOT_FLOAT": "is a string"},
	},
	{
		subject:     Float("NOT_FLOAT_FALLBACK"),
		env:         map[string]string{"NOT_FLOAT_FALLBACK": "invalid"},
		fallback:    []float64{123.33},
		getExpected: 123.33,
	},
	{
		subject:        Float("EXPONENT"),
		env:            map[string]string{"EXPONENT": "1e6"},
		fallback:       []float64{123},
		getExpected:    1_000_000,
		lookupExpected: 1_000_000,
	},
}

func TestFloatGet(t *testing.T) {
	for _, tc := range floatTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.getExpected, tc.subject.Get(tc.fallback...))
		})
	}
}

func TestFloatLookup(t *testing.T) {
	for _, tc := range floatTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.lookupExpected, tc.subject.Lookup(tc.fallback...))
		})
	}
}

var boolTestCases = []envarTestCase[Bool, bool]{
	{
		subject:        Bool("TRUE"),
		env:            map[string]string{"TRUE": "true"},
		fallback:       []bool{false},
		getExpected:    true,
		lookupExpected: true,
	},
	{
		subject:        Bool("FALSE"),
		env:            map[string]string{"FALSE": "false"},
		fallback:       []bool{true},
		getExpected:    false,
		lookupExpected: false,
	},
	{
		subject:        Bool("FALLBACK_UNSET"),
		env:            map[string]string{},
		fallback:       []bool{true},
		getExpected:    true,
		lookupExpected: true,
	},
	{
		subject:     Bool("FALLBACK_EMPTY"),
		env:         map[string]string{"FALLBACK_EMPTY": ""},
		fallback:    []bool{true},
		getExpected: true,
	},
	{
		subject: Bool("NO_FALLBACK_UNSET"),
	},
	{
		subject: Bool("NO_FALLBACK_EMPTY"),
		env:     map[string]string{"NO_FALLBACK_EMPTY": ""},
	},
	{
		subject:        Bool("INT_TRUE"),
		env:            map[string]string{"INT_TRUE": "1"},
		getExpected:    true,
		lookupExpected: true,
	},
	{
		subject:        Bool("INT_FALSE"),
		env:            map[string]string{"INT_FALSE": "0"},
		getExpected:    false,
		lookupExpected: false,
	},
	{
		subject: Bool("NOT_A_BOOL"),
		env:     map[string]string{"NOT_A_BOOL": "something"},
	},
	{
		subject:     Bool("NOT_A_BOOL_FALLBACK"),
		env:         map[string]string{"NOT_A_BOOL_FALLBACK": "something"},
		fallback:    []bool{true},
		getExpected: true,
	},
}

func TestBoolGet(t *testing.T) {
	for _, tc := range boolTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.getExpected, tc.subject.Get(tc.fallback...))
		})
	}
}

func TestBoolLookup(t *testing.T) {
	for _, tc := range boolTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.lookupExpected, tc.subject.Lookup(tc.fallback...))
		})
	}
}
