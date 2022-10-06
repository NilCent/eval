package parser

import (
	"fmt"
)

type ErrUnexpectedToken struct {
	Line     int
	Expected string
	Got      string
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("Line %d Unexpected Token: expected %s, got %s", e.Line, e.Expected, e.Got)
}

type ErrNoPrefixParseFn struct {
	Line      int
	TokenType string
}

func (e *ErrNoPrefixParseFn) Error() string {
	return fmt.Sprintf("Line %d Unrecognized Token: no prefix parse function for %s", e.Line, e.TokenType)
}

type ErrInteger struct {
	Err          error
	Line         int
	TokenLiteral string
}

func (e *ErrInteger) Error() string {
	return fmt.Sprintf("Line %d Could not parse %s as integer: %s", e.Line, e.TokenLiteral, e.Err.Error())
}
