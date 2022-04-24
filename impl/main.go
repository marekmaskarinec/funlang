package main

import (
	"fmt"
	"os"
)

func main() {
	lex := Lexer{
		buf: []rune(
`[2  0]->div`),
	}

	fmt.Println("input:", string(lex.buf))

	par := Parser{
		lexer: lex,
		defs:  map[string]Node{},
	}

	ast, err := par.topLevel()
	fmt.Println("ast:")
	fmt.Println(ast.toString(0))
	if err != nil {
		fmt.Println(err)
	}

	ev := Eval{
		defs: par.defs,
	}

	fmt.Println("output:")
	_, err = ev.expression(ast, 2, 2)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
