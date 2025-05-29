package dotenv

const (
	prefixComment     = `#`
	prefixSingleQuote = `'`
	prefixDoubleQuote = `"`
	prefixExport      = `export`
)

type parser struct {
	data []byte
	pos  int
}

func newParser(data []byte) *parser {
	return &parser{
		data: data,
		pos:  0,
	}
}

func (p *parser) parse() (map[string]string, error) {
	return nil, nil
}
