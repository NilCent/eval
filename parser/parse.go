package parser

import (
	"github.com/NilCent/eval/ast"
	"github.com/NilCent/eval/token"
	"strconv"
)

func (p *Parser) ParseProgram() (*ast.Program, error) {
	err := p.preload()
	if err != nil {
		return nil, err
	}
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curToken.Is(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)

		err = p.advance()
		if err != nil {
			return nil, err
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.curToken}

	err := p.expectPeek(token.IDENT)
	if err != nil {
		return nil, err
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	err = p.expectPeek(token.ASSIGN)
	if err != nil {
		return nil, err
	}

	err = p.advance()
	if err != nil {
		return nil, err
	}

	stmt.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Is(token.SEMICOLON) {
		err = p.advance()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	err := p.advance()
	if err != nil {
		return nil, err
	}

	stmt.ReturnValue, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Is(token.SEMICOLON) {
		err = p.advance()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = expr

	if p.peekToken.Is(token.SEMICOLON) {
		err = p.advance()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil, &ErrNoPrefixParseFn{Line: p.curToken.Line, TokenType: string(p.curToken.Type)}
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekToken.Is(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp, nil
		}

		err = p.advance()
		if err != nil {
			return nil, err
		}

		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Is(token.TRUE)}, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, &ErrInteger{Err: err, Line: p.curToken.Line, TokenLiteral: p.curToken.Literal}
	}

	lit.Value = value

	return lit, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	err := p.advance()
	if err != nil {
		return nil, err
	}

	expression.Right, err = p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	err := p.advance()
	if err != nil {
		return nil, err
	}

	expression.Right, err = p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	err := p.advance()
	if err != nil {
		return nil, err
	}

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	err = p.expectPeek(token.RPAREN)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	expression := &ast.IfExpression{Token: p.curToken}

	err := p.expectPeek(token.LPAREN)
	if err != nil {
		return nil, err
	}

	err = p.advance()
	if err != nil {
		return nil, err
	}

	expression.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(token.RPAREN)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(token.LBRACE)
	if err != nil {
		return nil, err
	}

	expression.Consequence, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	if p.peekToken.Is(token.ELSE) {
		err = p.advance()
		if err != nil {
			return nil, err
		}

		err = p.expectPeek(token.LBRACE)
		if err != nil {
			return nil, err
		}

		expression.Alternative, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	return expression, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	err := p.advance()
	if err != nil {
		return nil, err
	}

	for !p.curToken.Is(token.RBRACE) && !p.curToken.Is(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)

		err = p.advance()
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	err := p.expectPeek(token.LPAREN)
	if err != nil {
		return nil, err
	}

	lit.Parameters, err = p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	
	err = p.expectPeek(token.LBRACE)
	if err != nil {
		return nil, err
	}

	lit.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, error) {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Is(token.RPAREN) {
		err := p.advance()
		if err != nil {
			return nil, err
		}
		return identifiers, nil
	}

	err := p.advance()
	if err != nil {
		return nil, err
	}

	//todo 不是identifier怎么处理
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Is(token.COMMA) {
		err = p.advance()
		if err != nil {
			return nil, err
		}
		err = p.advance()
		if err != nil {
			return nil, err
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	err = p.expectPeek(token.RPAREN)
	if err != nil {
		return nil, err
	}

	return identifiers, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	args, err := p.parseCallArguments()
	if err != nil {
		return nil, err
	}
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = args
	return exp, nil
}

func (p *Parser) parseCallArguments() ([]ast.Expression, error) {
	args := []ast.Expression{}

	if p.peekToken.Is(token.RPAREN) {
		err := p.advance()
		if err != nil {
			return nil, err
		}
		return args, nil
	}

	err := p.advance()
	if err != nil {
		return nil, err
	}
	arg, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	args = append(args, arg)

	for p.peekToken.Is(token.COMMA) {
		err = p.advance()
		if err != nil {
			return nil, err
		}
		err = p.advance()
		if err != nil {
			return nil, err
		}
		arg, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	err = p.expectPeek(token.RPAREN)
	if err != nil {
		return nil, err
	}

	return args, nil
}