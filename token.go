package dotenv

import "fmt"

const (
	tknExport     = "EXPORT"
	tknIdentifier = "IDENT"
	tknEquals     = "EQUALS"
	tknValue      = "VALUE"
	tknComment    = "COMMENT"
	tknEOL        = "EOL"
	tknEOF        = "EOF"
	tknIllegal    = "ILLEGAL"
)

type token struct {
	Line    int
	Pos     int
	Type    string
	Literal string
}

func (t token) With(typ, value string) token {
	return token{
		Line:    t.Line,
		Pos:     t.Pos,
		Type:    typ,
		Literal: value,
	}
}

func (t token) String() string {
	return fmt.Sprintf("%s value=%s line=%d pos=%d",
		t.Type,
		t.Literal,
		t.Line,
		t.Pos,
	)
}
