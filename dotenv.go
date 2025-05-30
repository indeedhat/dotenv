package dotenv

import "os"

// Load loads the provided list of .env files into the os.environment.
// If no files are provided it will default to loading .env from the current working directory.
//
// Looad operoations will not replace any existing variables already in the environment.
func Load(filepaths ...string) error {
	for _, filepath := range pathFallback(filepaths) {
		p, err := ParseFile(filepath)
		if err != nil {
			return err
		}

		assignEnvars(p.Parse(), false)
	}

	return nil
}

// LoadStrict loads the provided list of .env files into the os.environment.
// If no files are provided it will default to loading .env from the current working directory.
//
// Looad operoations will not replace any existing variables already in the environment.
//
// Strict operations will break at the first .env file that contains invalid syntax, all previous
// files in the list will still be loaded into the environment but any files that appear in the list
// after the first invalid file will be skipped
func LoadStrict(filepaths ...string) error {
	for _, filepath := range pathFallback(filepaths) {
		p, err := ParseFile(filepath)
		if err != nil {
			return err
		}

		pairs, err := p.ParseStrict()
		if err != nil {
			return err
		}

		assignEnvars(pairs, false)
	}

	return nil
}

// Overload loads the provided list of .env files into the os.environment.
// If no files are provided it will default to loading .env from the current working directory.
//
// Unlike with the Load operation any existing environment variables will be overloaded with the
// present in the provided env files.
func Overload(filepaths ...string) error {
	for _, filepath := range pathFallback(filepaths) {
		p, err := ParseFile(filepath)
		if err != nil {
			return err
		}

		assignEnvars(p.Parse(), true)
	}

	return nil
}

// OverloadStrict loads the provided list of .env files into the os.environment.
// If no files are provided it will default to loading .env from the current working directory.
//
// Unlike with the Load operation any existing environment variables will be overloaded with the
// present in the provided env files.
//
// Strict operations will break at the first .env file that contains invalid syntax, all previous
// files in the list will still be loaded into the environment but any files that appear in the list
// after the first invalid file will be skipped
func OverloadStrict(filepaths ...string) error {
	for _, filepath := range pathFallback(filepaths) {
		p, err := ParseFile(filepath)
		if err != nil {
			return err
		}

		pairs, err := p.ParseStrict()
		if err != nil {
			return err
		}

		assignEnvars(pairs, true)
	}

	return nil
}

// ParseFile returns the underlying Parser instance representing the provided env file
//
// This will not load anything into the environment but allow you to handle the found values manually
func ParseFile(filepath string) (*Parser, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return newParser(newLexer(string(data))), nil
}

func pathFallback(filepaths []string) []string {
	if len(filepaths) == 0 {
		return []string{".env"}
	}

	return filepaths
}

func assignEnvars(pairs []ParseEntry, overwrite bool) {
	for _, v := range pairs {
		if !overwrite {
			if _, ok := os.LookupEnv(v.Key); ok {
				continue
			}
		}

		if !v.Raw && v.Value != "" {
			os.Setenv(v.Key, Expand(v.Value, os.Getenv))
		} else {
			os.Setenv(v.Key, v.Value)
		}
	}
}
