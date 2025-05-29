package dotenv

import (
	"bytes"
	"strings"
	"unicode"
)

const runeEOF = 0

type Lexer struct {
	data []rune

	pos     int
	readPos int
	char    rune

	line      int
	linePos   int
	prevToken Token
}

func NewLexer(data string) *Lexer {
	l := &Lexer{data: []rune(data)}
	l.readRune()

	return l
}

func (l *Lexer) NextToken() Token {
	var skipRead bool
	defer func() {
		if !skipRead {
			l.readRune()
		}
	}()

	tkn := (func() Token {
	skip:
		switch l.char {
		case '\r':
			if l.peekRune() == '\n' {
				l.readRune()
				goto skip
			}
			l.line++
			l.linePos = 0
			return Token{TknEOL, ""}
		case '\n':
			l.line++
			l.linePos = 0
			return Token{TknEOL, ""}
		case '\t', '\v', '\f', ' ', 0x85, 0xA0:
			l.readRune()
			goto skip
		case runeEOF:
			return Token{TknEOF, ""}
		case '#':
			skipRead = true
			return Token{TknComment, l.readComment()}
		case '=':
			return Token{TknEquals, "="}
		case '\'', '"':
			str, ok := l.readQuotedString(l.char)
			if !ok {
				return Token{TknIllegal, str}
			}
			return Token{TknValue, str}
		case 'e':
			if l.peekIdentifier() == "export" {
				l.readIdentifier()
				return Token{TknExport, "export"}
			}
			fallthrough
		default:
			if l.prevToken.Type != TknEquals {
				if !isIdentRune(l.char) {
					return Token{TknIllegal, string(l.char)}
				}

				return Token{TknIdentifier, l.readIdentifier()}
			}

			skipRead = true
			return Token{TknValue, l.readUnquotedString()}
		}
	})()

	l.prevToken = tkn
	return tkn
}

func (l *Lexer) readRune() {
	if l.readPos >= len(l.data) {
		l.char = 0
	} else {
		l.char = l.data[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
	l.linePos++
}

func (l *Lexer) peekRune() rune {
	if l.readPos >= len(l.data) {
		return 0
	}

	return l.data[l.readPos]
}

func (l *Lexer) readComment() string {
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

func (l *Lexer) readQuotedString(terminator rune) (string, bool) {
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

func (l *Lexer) readUnquotedString() string {
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

func (l *Lexer) readIdentifier() string {
	pos := l.pos

	for isIdentRune(l.peekRune(), true) {
		l.readRune()
	}

	return string(l.data[pos:l.readPos])
}

func (l *Lexer) peekIdentifier() string {
	pos := l.pos
	readPos := l.readPos

	ident := l.readIdentifier()

	l.pos = pos
	l.readPos = readPos

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
