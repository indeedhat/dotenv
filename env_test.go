package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var stringTestCases = []struct {
	subject  String
	env      map[string]string
	fallback []string
	expected string
}{
	{
		String("VALID"),
		map[string]string{"VALID": "valid"},
		[]string{"fallback value"},
		"valid",
	},
	{
		String("FALLBACK"),
		map[string]string{"VALID": "valid"},
		[]string{"fallback value"},
		"fallback value",
	},
	{
		String("NO_FALLBACK"),
		map[string]string{"VALID": "valid"},
		nil,
		"",
	},
}

func TestString(t *testing.T) {
	for _, tc := range stringTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.expected, tc.subject.Get(tc.fallback...))
		})
	}
}

var intTestCases = []struct {
	subject  Int
	env      map[string]string
	fallback []int
	expected int
}{
	{
		Int("VALID"),
		map[string]string{"VALID": "100"},
		[]int{200},
		100,
	},
	{
		Int("FALLBACK"),
		map[string]string{"VALID": "valid"},
		[]int{200},
		200,
	},
	{
		Int("NO_FALLBACK"),
		map[string]string{"VALID": "valid"},
		nil,
		0,
	},
	{
		Int("NOT_INT"),
		map[string]string{"NOT_INT": "is a string"},
		nil,
		0,
	},
	{
		Int("NOT_INT_FALLBACK"),
		map[string]string{"NOT_INT_FALLBACK": "invalid"},
		[]int{123},
		123,
	},
	{
		Int("HEX"),
		map[string]string{"HEX": "0xF"},
		[]int{123},
		15,
	},
	{
		Int("BIN"),
		map[string]string{"BIN": "0b10000000"},
		[]int{123},
		128,
	},
}

func TestInt(t *testing.T) {
	for _, tc := range intTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.expected, tc.subject.Get(tc.fallback...))
		})
	}
}

var floatTestCases = []struct {
	subject  Float
	env      map[string]string
	fallback []float64
	expected float64
}{
	{
		Float("VALID"),
		map[string]string{"VALID": "100.5"},
		[]float64{200},
		100.5,
	},
	{
		Float("FALLBACK"),
		map[string]string{"VALID": "valid"},
		[]float64{200},
		200,
	},
	{
		Float("NO_FALLBACK"),
		map[string]string{"VALID": "valid"},
		nil,
		0,
	},
	{
		Float("NOT_FLOAT"),
		map[string]string{"NOT_FLOAT": "is a string"},
		nil,
		0,
	},
	{
		Float("NOT_FLOAT_FALLBACK"),
		map[string]string{"NOT_FLOAT_FALLBACK": "invalid"},
		[]float64{123.33},
		123.33,
	},
	{
		Float("EXPONENT"),
		map[string]string{"EXPONENT": "1e6"},
		[]float64{123},
		1_000_000,
	},
}

func TestFloat(t *testing.T) {
	for _, tc := range floatTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.expected, tc.subject.Get(tc.fallback...))
		})
	}
}

var boolTestCases = []struct {
	subject  Bool
	env      map[string]string
	fallback []bool
	expected bool
}{
	{
		Bool("TRUE"),
		map[string]string{"TRUE": "true"},
		[]bool{false},
		true,
	},
	{
		Bool("FALSE"),
		map[string]string{"FALSE": "false"},
		[]bool{true},
		false,
	},
	{
		Bool("FALLBACK"),
		map[string]string{},
		[]bool{true},
		true,
	},
	{
		Bool("NO_FALLBACK"),
		map[string]string{},
		nil,
		false,
	},
	{
		Bool("INT_TRUE"),
		map[string]string{"INT_TRUE": "1"},
		nil,
		true,
	},
	{
		Bool("INT_FALSE"),
		map[string]string{"INT_FALSE": "0"},
		nil,
		false,
	},
}

func TestBool(t *testing.T) {
	for _, tc := range boolTestCases {
		t.Run(string(tc.subject), func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			require.Equal(t, tc.expected, tc.subject.Get(tc.fallback...))
		})
	}
}
