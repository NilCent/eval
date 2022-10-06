package lexer

import (
	"github.com/NilCent/eval/token"
)

type Lexer struct {
	input   string
	line    int
	start   int
	current int
}

func (l *Lexer) advance() byte {
	var ch byte
	if l.current < len(l.input) {
		ch = l.input[l.current]
		l.current++
	}

	return ch
}

func (l *Lexer) peek() byte {
	if l.current >= len(l.input) {
		return 0
	} else {
		return l.input[l.current]
	}
}
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
	}
	return l
}

func (l *Lexer) newToken(tokenType token.TokenType, str string) token.Token {
	l.start = l.current
	return token.New(tokenType, str, l.line)
}

func (l *Lexer) readIdentifier() string {
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.advance()
	}
	return l.input[l.start:l.current]
}
func (l *Lexer) readNumber() string {
	for isDigit(l.peek()) {
		l.advance()
	}
	return l.input[l.start:l.current]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) NextToken() (token.Token, error) {
	var tok token.Token

Re:
	ch := l.advance()
	switch ch {
	case '\n':
		l.line++
		fallthrough
	case '\t':
		fallthrough
	case '\r':
		fallthrough
	case ' ':
		l.start = l.current
		goto Re
	case '=':
		if l.peek() == '=' {
			next := l.advance()
			literal := string(ch) + string(next)
			tok = l.newToken(token.EQ, literal)
		} else {
			tok = l.newToken(token.ASSIGN, string(ch))
		}
	case '+':
		tok = l.newToken(token.PLUS, string(ch))
	case '-':
		tok = l.newToken(token.MINUS, string(ch))
	case '!':
		if l.peek() == '=' {
			next := l.advance()
			literal := string(ch) + string(next)
			tok = l.newToken(token.NOT_EQ, literal)
		} else {
			tok = l.newToken(token.BANG, string(ch))
		}
	case '/':
		tok = l.newToken(token.SLASH, string(ch))
	case '*':
		tok = l.newToken(token.ASTERISK, string(ch))
	case '<':
		tok = l.newToken(token.LT, string(ch))
	case '>':
		tok = l.newToken(token.GT, string(ch))
	case ';':
		tok = l.newToken(token.SEMICOLON, string(ch))
	case ',':
		tok = l.newToken(token.COMMA, string(ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(ch))
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
	default:
		if isLetter(ch) {
			literal := l.readIdentifier()
			tok = l.newToken(token.LookupIdent(literal), literal)
		} else if isDigit(ch) {
			tok = l.newToken(token.INT, l.readNumber())
		} else {
			return tok, &ErrUnexpectedChar{Line: l.line, Char: ch}
		}
	}
	return tok, nil
}
