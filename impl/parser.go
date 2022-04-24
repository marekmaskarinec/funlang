package main

type Parser struct {
	lexer Lexer
	token Token
	defs  map[string]Node
}

func (p *Parser) isDefined(name string) bool {
	_, ok := p.defs[name]
	return ok
}

func (p *Parser) accept() error {
	var err error
	p.token, err = p.lexer.next()
	return err
}

func (p *Parser) arrayLiteral() (Node, error) {
	ast := Node{value: p.token}

	for p.token.kind != kindArrayEnd {
		node, err := p.expression()
		if err != nil {
			return ast, err
		}
		ast.children = append(ast.children, node)

		p.accept()
		if p.token.kind == kindEOF {
			return ast, p.token.errorf("Unended array literal.")
		}

		if val, _ := p.lexer.peek(); p.token.kind == kindLF && val.kind == kindArrayEnd {
			p.accept()
		}

		if p.token.kind != kindArraySep &&
			p.token.kind != kindLF &&
			p.token.kind != kindArrayEnd {

			return ast, p.token.errorf(
				"Expected an array delimeter or an array end, got %s.", p.token.value)
		}
	}

	return ast, nil
}

func (p *Parser) parenthesis() (Node, error) {
	ast := Node{value: p.token}

	for {
		if token, _ := p.lexer.peek(); token.kind == kindParClose {
			p.accept()
			break
		}

		node, err := p.expression()
		if err != nil {
			return ast, err
		}

		ast.children = append(ast.children, node)
	}

	return ast, nil
}

func (p *Parser) expression() (Node, error) {
	// skip LF tokens. they are only needed for array literals
	for token, _ := p.lexer.peek(); token.kind == kindLF; token, _ = p.lexer.peek() {
		p.accept()
	}

	ast := Node{}
	p.accept()
	if p.token.value == "fun" {
		ast.value = p.token
		node, err := p.expression()
		if err != nil {
			return ast, err
		}

		ast.children = append(ast.children, node)

	} else if p.token.value == "def" {
		if err := p.accept(); err != nil {
			return ast, p.token.errorf("%w", err)
		}

		if p.token.kind != kindIdent {
			return ast, p.token.errorf("Expected an identifier, got %d", p.token.kind)
		}

		name := p.token.value
		// make sure the identifier exists when evaluating
		p.defs[name] = Node{}
		var err error
		p.defs[name], err = p.expression()
		if err != nil {
			return ast, err
		}

	} else if p.token.value == "[" {
		var err error
		ast, err = p.arrayLiteral()
		if err != nil {
			return ast, err
		}

	} else if p.token.value == "(" {
		var err error
		ast, err = p.parenthesis()
		if err != nil {
			return ast, err
		}

	} else if p.token.kind == kindIdent {
		if !p.isDefined(p.token.value) && !isBuiltin(p.token.value) {
			return ast, p.token.errorf("Unknown identifier %s.", p.token.value)
		}

		ast.value = p.token

	} else if p.token.kind == kindNumber ||
		p.token.kind == kindArgSubst {
		ast.value = p.token
	}

	if token, _ := p.lexer.peek(); token.value == "-" {
		p.accept()
		return ast, nil
	}

	for {
		token, _ := p.lexer.peek()
		if token.value != "->" && token.value != "!" {
			break
		}

		p.accept()
		node, err := p.expression()
		if err != nil {
			return ast, err
		}

		ast = Node{
			value:    token,
			children: []Node{ast, node},
		}
	}

	return ast, nil
}

func (p *Parser) topLevel() (Node, error) {
	ast := Node{}
	for {
		node, err := p.expression()
		if err != nil {
			return ast, err
		}

		if p.token.kind == kindEOF {
			return ast, nil
		}

		ast.children = append(ast.children, node)
	}
}
