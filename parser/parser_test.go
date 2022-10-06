package parser

import (
	"reflect"
	"testing"

	"github.com/NilCent/eval/ast"
	"github.com/NilCent/eval/lexer"
)

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
func TestInteger(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue int
	}{
		{"5", 5},
		{"5;", 5},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Error(err)
			continue
		}

		if len(program.Statements) != 1 {
			t.Errorf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		testLiteralExpression(t, stmt.Expression, tc.expectedValue)
	}
}


func TestIdentifier(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue string
	}{
		{"asd_1","asd_1"},
		{"asd_1;", "asd_1"},
		{"_", "_"},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Error(err)
			continue
		}

		if len(program.Statements) != 1 {
			t.Errorf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		testLiteralExpression(t, stmt.Expression, tc.expectedValue)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expectedOperator string	
		expectedValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Error(err)
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tc.expectedOperator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tc.expectedOperator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tc.expectedValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Error(err)
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tc.leftValue,
			tc.operator, tc.rightValue) {
			return
		}
	}
}

func TestErrorCase(t *testing.T) {
	testCases := []struct {
		input         string
		expectedError error
	}{
		{"@", &lexer.ErrUnexpectedChar{}},
		{"1111111111111111111111111111111111111111111", &ErrInteger{}},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		_, err := p.ParseProgram()
		if err == nil {
			t.Error("error can not be nil")
		} else {
			if tc.expectedError != nil && reflect.TypeOf(err) == reflect.TypeOf(tc.expectedError) {
				t.Log(err)
			} else {
				t.Errorf("expected %#v, got %#v", tc.expectedError, err)
			}
		}
	}
}