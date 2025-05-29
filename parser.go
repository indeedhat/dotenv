package dotenv

import (
	"fmt"
)

type Parser struct {
	lex *lexer
}

func newParser(l *lexer) *Parser {
	return &Parser{
		lex: l,
	}
}

func (p *Parser) Parse() map[string]string {
	pairs := make(map[string]string)

	prev := make([]*token, 2)

	for {
		tkn := p.lex.NextToken()
		if tkn.Type == tknEOF {
			break
		}

		switch tkn.Type {
		case tknIdentifier:
			prev[0] = &tkn
			prev[1] = nil
		case tknEquals:
			if prev[0] != nil {
				prev[1] = &tkn
			}
		case tknValue, tknRawValue:
			if prev[0] != nil && prev[1] != nil {
				pairs[prev[0].Literal] = tkn.Literal
				break
			}
			fallthrough
		default:
			prev[0] = nil
			prev[1] = nil
		}
	}

	return pairs
}

func (p *Parser) ParseStrict() (map[string]string, error) {
	pairs := make(map[string]string)

loop:
	for {
		tkn := p.lex.NextToken()
		switch tkn.Type {
		case tknEOF:
			break loop
		case tknComment, tknEOL:
			// NB: these are all allowed on their own line and will be skipped
		case tknExport:
			tkn = p.lex.NextToken()
			if tkn.Type != tknIdentifier {
				return nil, fmt.Errorf("Unexpected token %s", tkn)
			}
			fallthrough
		case tknIdentifier:
			eqTkn := p.lex.NextToken()
			if eqTkn.Type != tknEquals {
				return nil, fmt.Errorf("Unexpected token %s", eqTkn)
			}

			valTkn := p.lex.NextToken()
			if valTkn.Type != tknValue && valTkn.Type != tknRawValue {
				return nil, fmt.Errorf("Unexpected token %s", valTkn)
			}

			pairs[tkn.Literal] = valTkn.Literal
		default:
			return nil, fmt.Errorf("Unexpected token %s", tkn)
		}
	}

	return pairs, nil
}
