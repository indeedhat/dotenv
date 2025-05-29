package dotenv

import (
	"os"
	"strconv"
)

type EnVar[T any] interface {
	Get(...T) T
	Lookup(...T) T
}

type String string

// Get the value of the String envar
func (k String) Get(fallback ...string) string {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return val
}

func (k String) Lookup(fallback ...string) string {
	val, ok := os.LookupEnv(string(k))

	if !ok && len(fallback) > 0 {
		return fallback[0]
	}

	return val
}

var _ EnVar[string] = (*String)(nil)

type Int string

// Get the value of the Int envar
func (k Int) Get(fallback ...int) int {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.ParseInt(val, 0, 0)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return int(parsed)
}

func (k Int) Lookup(fallback ...int) int {
	val, ok := os.LookupEnv(string(k))

	if !ok && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, _ := strconv.ParseInt(val, 0, 0)
	return int(parsed)
}

var _ EnVar[int] = (*Int)(nil)

type Float string

// Get the value of the Float envar
func (k Float) Get(fallback ...float64) float64 {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return parsed
}

func (k Float) Lookup(fallback ...float64) float64 {
	val, ok := os.LookupEnv(string(k))

	if !ok && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, _ := strconv.ParseFloat(val, 64)
	return parsed
}

var _ EnVar[float64] = (*Float)(nil)

type Bool string

// Get the value of the Bool envar
func (k Bool) Get(fallback ...bool) bool {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return parsed
}

func (k Bool) Lookup(fallback ...bool) bool {
	val, ok := os.LookupEnv(string(k))

	if !ok && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, _ := strconv.ParseBool(val)
	return parsed
}

var _ EnVar[bool] = (*Bool)(nil)
