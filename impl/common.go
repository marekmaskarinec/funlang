package main

import (
	"fmt"
	"strconv"
)

type Location struct {
	line, column int
	index        int
}

const (
	kindNull = iota
	kindArrayStart
	kindArrayEnd
	kindArraySep
	kindArgSubst
	kindLF
	kindChainOperator
	kindEOF
	kindNumber
	kindIdent
	kindRef
	kindParOpen
	kindParClose
)

type TokenKind int

type Token struct {
	kind  TokenKind
	loc   Location
	value []rune
}

func (t *Token) errorf(format string, args ...interface{}) error {
	argss := make([]interface{}, len(args)+2)
	argss[0] = t.loc.line
	argss[1] = t.loc.column
	for i, arg := range args {
		argss[i+2] = arg
	}

	return fmt.Errorf("(%d %d): "+format, argss...)
}

func (t *Token) toInt() (int, error) {
	val, err := strconv.Atoi(string(t.value))
	if err != nil {
		return val, t.errorf("expected a number")
	}

	return val, nil
}

type Node struct {
	value    Token
	children []Node
}

func (n *Node) toString(indent int) string {
	out := ""

	for i := 0; i < indent; i++ {
		out += "  "
	}

	out += string(n.value.value) + "\n"

	for i := 0; i < len(n.children); i++ {
		out += n.children[i].toString(indent + 1)
	}

	return out
}

func isBuiltin(name string) bool {
	var builtins = []string{
		"print", "neg", "branch", "greater", "sub", "div", "add", "mul", "call",
	}

	for i := 0; i < len(builtins); i++ {
		if name == builtins[i] {
			return true
		}
	}

	return false
}
