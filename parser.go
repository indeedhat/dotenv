package dotenv

import (
	"fmt"
)

type Parser struct {
	lex *lexer
}
type ParseEntry struct {
	Key   string
	Value string
	Raw   bool
}

func newParser(l *lexer) *Parser {
	return &Parser{
		lex: l,
	}
}

func (p *Parser) Parse() []ParseEntry {
	var pairs []ParseEntry

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
		case tknValue, tknRawValue, tknComment, tknEOL, tknEOF:
			if prev[0] != nil && prev[1] != nil {
				if tkn.Type == tknValue || tkn.Type == tknRawValue {
					pairs = append(pairs, ParseEntry{prev[0].Literal, tkn.Literal, tkn.Type == tknRawValue})
				} else {
					pairs = append(pairs, ParseEntry{prev[0].Literal, "", false})
				}
			}
			fallthrough
		default:
			prev[0] = nil
			prev[1] = nil
		}
	}

	return pairs
}

func (p *Parser) ParseStrict() ([]ParseEntry, error) {
	var pairs []ParseEntry

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
			switch valTkn.Type {
			case tknValue, tknRawValue:
				pairs = append(pairs, ParseEntry{tkn.Literal, valTkn.Literal, valTkn.Type == tknRawValue})
			case tknComment, tknEOL, tknEOF:
				pairs = append(pairs, ParseEntry{tkn.Literal, "", false})
			default:
				return nil, fmt.Errorf("Unexpected token %s", valTkn)
			}
		default:
			return nil, fmt.Errorf("Unexpected token %s", tkn)
		}
	}

	return pairs, nil
}
