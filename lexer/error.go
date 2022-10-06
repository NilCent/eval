package lexer

import "fmt"

type ErrUnexpectedChar struct {
	Line int
	Char byte
}

func (e *ErrUnexpectedChar) Error() string {
	return fmt.Sprintf("Line %d Unexpected character: %s", e.Line, string(e.Char))
}
