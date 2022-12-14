package eval

import (
	"errors"
	"github.com/NilCent/eval/evaluator"
	"github.com/NilCent/eval/lexer"
	"github.com/NilCent/eval/object"
	"github.com/NilCent/eval/parser"
	"fmt"
)

type interpreter struct {
	env *object.Environment
}

func New() *interpreter {
	return &interpreter{
		env: object.NewEnvironment(),
	}
}

func (i *interpreter) EvalInt(input string) (int, error) {

	l := lexer.New(input)
	p := parser.New(l)

	program, err := p.ParseProgram()
	if err != nil {
		return 0, err
	}

	evaluated := evaluator.Eval(program, i.env)
	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			return 0, errors.New(evaluated.Inspect())
		} else if evaluated.Type() != object.INTEGER_OBJ {
			return 0, errors.New(fmt.Sprintf("expect %s, got %s", object.INTEGER_OBJ, evaluated.Type()))
		} else {
			return int(evaluated.(*object.Integer).Value), nil
		}
	}

	return 0, nil
}

func (i *interpreter) Do(input string) error {
	l := lexer.New(input)
	p := parser.New(l)

	program, err := p.ParseProgram()
	if err != nil {
		return err
	}

	evaluator.Eval(program, i.env)

	return nil
}
