package dotenv

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

const runeEOF = 0

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

type lexer struct {
	data []rune

	pos     int
	readPos int
	char    rune

	line      int
	linePos   int
	prevToken token
}

func newLexer(data string) *lexer {
	l := &lexer{data: []rune(data)}
	l.readRune()

	return l
}

func (l *lexer) NextToken() token {
	var skipRead bool
	defer func() {
		if !skipRead {
			l.readRune()
		}
	}()

	tkn := (func() token {
	skip:
		switch l.char {
		case '\r':
			if l.peekRune() == '\n' {
				l.readRune()
				goto skip
			}

			tkn := l.tkn().With(tknEOL, "")
			l.line++
			l.linePos = 0

			return tkn
		case '\n':
			tkn := l.tkn().With(tknEOL, "")
			l.line++
			l.linePos = 0

			return tkn
		case '\t', '\v', '\f', ' ', 0x85, 0xA0:
			l.readRune()
			goto skip
		case runeEOF:
			return l.tkn().With(tknEOF, "")
		case '#':
			skipRead = true
			return l.tkn().With(tknComment, l.readComment())
		case '=':
			return l.tkn().With(tknEquals, "=")
		case '\'', '"':
			tkn := l.tkn()
			str, ok := l.readQuotedString(l.char)
			if !ok {
				return tkn.With(tknIllegal, str)
			}
			return tkn.With(tknValue, str)
		case 'e':
			if l.peekIdentifier() == "export" {
				defer l.readIdentifier()
				return l.tkn().With(tknExport, "export")
			}
			fallthrough
		default:
			if l.prevToken.Type != tknEquals {
				if !isIdentRune(l.char) {
					return l.tkn().With(tknIllegal, string(l.char))
				}

				return l.tkn().With(tknIdentifier, l.readIdentifier())
			}

			skipRead = true
			return l.tkn().With(tknValue, l.readUnquotedString())
		}
	})()

	l.prevToken = tkn
	return tkn
}

func (l *lexer) tkn() token {
	return token{
		Line: l.line,
		Pos:  l.linePos,
	}
}

func (l *lexer) readRune() {
	if l.readPos >= len(l.data) {
		l.char = 0
	} else {
		l.char = l.data[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
	l.linePos++
}

func (l *lexer) peekRune() rune {
	if l.readPos >= len(l.data) {
		return 0
	}

	return l.data[l.readPos]
}

func (l *lexer) readComment() string {
	// skip the #
	l.readRune()

	var buf bytes.Buffer

	for {
		curRune := l.data[l.pos]
		peekRune := l.peekRune()

		if peekRune == runeEOF {
			break
		} else if curRune == '\r' && peekRune != '\n' {
			break
		} else if curRune == '\n' {
			break
		}

		buf.WriteRune(curRune)
		l.readRune()
	}

	return strings.TrimSpace(buf.String())
}

func (l *lexer) readQuotedString(terminator rune) (string, bool) {
	var buf bytes.Buffer

	for {
		curRune := l.data[l.pos]
		peekRune := l.peekRune()
		// if we dont find a closing " then its an invalid string
		if peekRune == runeEOF {
			return "", false
		}

		if curRune != '\\' || peekRune != '"' {
			buf.WriteRune(curRune)
		}

		l.readRune()
		if peekRune == terminator && curRune != '\\' {
			break
		}
	}

	// skip final "
	if l.peekRune() == terminator {
		l.readRune()
	}

	if buf.Len() < 2 {
		return "", true
	}

	// dont include the " on each side in the strings value
	return buf.String()[1:], true
}

func (l *lexer) readUnquotedString() string {
	var buf bytes.Buffer

	for {
		curRune := l.data[l.pos]
		peekRune := l.peekRune()

		if peekRune == runeEOF {
			break
		}

		if curRune == '\n' {
			break
		}

		if unicode.IsSpace(curRune) && peekRune == '#' {
			break
		}

		buf.WriteRune(curRune)
		l.readRune()
	}

	return strings.TrimSpace(buf.String())
}

func (l *lexer) readIdentifier() string {
	pos := l.pos

	for isIdentRune(l.peekRune(), true) {
		l.readRune()
	}

	return string(l.data[pos:l.readPos])
}

func (l *lexer) peekIdentifier() string {
	pos := l.pos
	readPos := l.readPos
	linePos := l.linePos
	line := l.line

	ident := l.readIdentifier()

	l.pos = pos
	l.readPos = readPos
	l.line = line
	l.linePos = linePos

	return ident
}

func isIdentRune(char rune, subsequent ...bool) bool {
	// numbers are allowed in idents but not as the first character
	if len(subsequent) > 0 && subsequent[0] && char >= '0' && char <= '9' {
		return true
	}

	return char >= 'a' && char <= 'z' ||
		char >= 'A' && char <= 'Z' ||
		char == '_'
}
