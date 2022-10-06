package lexer

import (
	"testing"

	"github.com/NilCent/eval/token"
)

func TestLexer(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	expectResult := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectErr error
	}{
		{token.LET, "let", nil},
		{token.IDENT, "five", nil},
		{token.ASSIGN, "=", nil},
		{token.INT, "5", nil},
		{token.SEMICOLON, ";", nil},
		{token.LET, "let", nil},
		{token.IDENT, "ten", nil},
		{token.ASSIGN, "=", nil},
		{token.INT, "10", nil},
		{token.SEMICOLON, ";", nil},
		{token.LET, "let", nil},
		{token.IDENT, "add", nil},
		{token.ASSIGN, "=", nil},
		{token.FUNCTION, "fn", nil},
		{token.LPAREN, "(", nil},
		{token.IDENT, "x", nil},
		{token.COMMA, ",", nil},
		{token.IDENT, "y", nil},
		{token.RPAREN, ")", nil},
		{token.LBRACE, "{", nil},
		{token.IDENT, "x", nil},
		{token.PLUS, "+", nil},
		{token.IDENT, "y", nil},
		{token.SEMICOLON, ";", nil},
		{token.RBRACE, "}", nil},
		{token.SEMICOLON, ";", nil},
		{token.LET, "let", nil},
		{token.IDENT, "result", nil},
		{token.ASSIGN, "=", nil},
		{token.IDENT, "add", nil},
		{token.LPAREN, "(", nil},
		{token.IDENT, "five", nil},
		{token.COMMA, ",", nil},
		{token.IDENT, "ten", nil},
		{token.RPAREN, ")", nil},
		{token.SEMICOLON, ";", nil},
		{token.BANG, "!", nil},
		{token.MINUS, "-", nil},
		{token.SLASH, "/", nil},
		{token.ASTERISK, "*", nil},
		{token.INT, "5", nil},
		{token.SEMICOLON, ";", nil},
		{token.INT, "5", nil},
		{token.LT, "<", nil},
		{token.INT, "10", nil},
		{token.GT, ">", nil},
		{token.INT, "5", nil},
		{token.SEMICOLON, ";", nil},
		{token.IF, "if", nil},
		{token.LPAREN, "(", nil},
		{token.INT, "5", nil},
		{token.LT, "<", nil},
		{token.INT, "10", nil},
		{token.RPAREN, ")", nil},
		{token.LBRACE, "{", nil},
		{token.RETURN, "return", nil},
		{token.TRUE, "true", nil},
		{token.SEMICOLON, ";", nil},
		{token.RBRACE, "}", nil},
		{token.ELSE, "else", nil},
		{token.LBRACE, "{", nil},
		{token.RETURN, "return", nil},
		{token.FALSE, "false", nil},
		{token.SEMICOLON, ";", nil},
		{token.RBRACE, "}", nil},
		{token.INT, "10", nil},
		{token.EQ, "==", nil},
		{token.INT, "10", nil},
		{token.SEMICOLON, ";", nil},
		{token.INT, "10", nil},
		{token.NOT_EQ, "!=", nil},
		{token.INT, "9", nil},
		{token.SEMICOLON, ";", nil},
		{token.EOF, "", nil},
	}

	l := New(input)

	for _, res := range expectResult {
		tok, _ := l.NextToken()
		if res.expectedType != tok.Type || 
		res.expectedLiteral != tok.Literal {
			t.Errorf("expected %v, got %v", res, tok)
		}
	}
}
