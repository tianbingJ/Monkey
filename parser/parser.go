package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)
//Pratt Parser
type (
	prefixParseFn func() ast.Expression
	infixParseFn func(expression ast.Expression)  ast.Expression
)

type Parser struct {
	l *lexer.Lexer
	errors []string

	curToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}


func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() * ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.peekTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()  //跳过';'
	}
	return program
}

func (p *Parser)  parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser)  parseLetStatement() ast.Statement {
	letStmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil;
	}

	letStmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	//parse expressions

	//跳过，直到遇到'；'
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return letStmt
}

func (p *Parser) parseReturnStatement() ast.Statement{

	returnStmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	//跳过，直到遇到'；'
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return returnStmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false;
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
