package dotenv

import "os"

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

func assignEnvars(pairs map[string]string, overwrite bool) {
	for k, v := range pairs {
		if !overwrite {
			if _, ok := os.LookupEnv(k); !ok {
				continue
			}
		}

		os.Setenv(k, v)
	}
}
