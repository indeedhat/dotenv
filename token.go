package dotenv

const (
	TknExport     = "EXPORT"
	TknIdentifier = "IDENT"
	TknEquals     = "EQUALS"
	TknValue      = "VALUE"
	TknComment    = "COMMENT"
	TknEOL        = "EOL"
	TknEOF        = "EOF"
	TknIllegal    = "ILLEGAL"
)

type Token struct {
	Type    string
	Literal string
}
